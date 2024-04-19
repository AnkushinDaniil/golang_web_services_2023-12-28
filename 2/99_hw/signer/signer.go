package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
)

const quotaLimit = 1

func ExecutePipeline(hashSignJobs ...job) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	var input, output chan interface{}

	for _, hashSignJob := range hashSignJobs {
		input = output
		output = make(chan interface{})

		wg.Add(1)
		go func(j job, in, out chan interface{}) {
			j(in, out)
			wg.Done()
			close(out)
		}(hashSignJob, input, output)
	}
}

/*
SingleHash считает значение DataSignerCrc32(data)+"~"
+DataSignerCrc32(DataSignerMd5(data))
( конкатенация двух строк через ~), где data - то что пришло
на вход (по сути - числа из первой функции)
*/
func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	var md5Counter uint32

	for val := range in {
		data := fmt.Sprintf("%v", val)
		md5 := make(chan string)

		wg.Add(1)
		go func(m chan<- string) {
			defer wg.Done()
			for curCounter := atomic.LoadUint32(&md5Counter); curCounter >= quotaLimit; {
			}
			atomic.AddUint32(&md5Counter, 1)
			defer atomic.SwapUint32(&md5Counter, 0)
			md5 <- DataSignerMd5(data)
		}(md5)

		wg.Add(1)
		go func(d string, m <-chan string) {
			md5String := <-m
			crc32_md5 := DataSignerCrc32(md5String)
			crc32 := DataSignerCrc32(d)
			res := fmt.Sprintf("%v~%v", crc32, crc32_md5)

			// fmt.Println("SingleHash", d)
			// fmt.Println("SingleHash md5(data)", m)
			// fmt.Println("SingleHash crc32(md5(data))", crc32_md5)
			// fmt.Println("SingleHash crc32(data)", crc32)
			// fmt.Println("SingleHash result", res)

			out <- res
			wg.Done()
		}(data, md5)

	}
}

/*
MultiHash считает значение crc32(th+data))
(конкатенация цифры, приведённой к строке и строки),
где th=0..5 ( т.е. 6 хешей на каждое входящее значение ),
потом берёт конкатенацию результатов в порядке расчета (0..5),
где data - то что пришло на вход
(и ушло на выход из SingleHash)
*/
func MultiHash(in, out chan interface{}) {
	for val := range in {
		data, ok := (val).(string)
		if !ok {
			fmt.Println("cant convert data to string")
		}
		res := ""
		var i byte = '0'
		for ; i < '6'; i++ {
			t := DataSignerCrc32(string(i) + data)
			// fmt.Println(data, "MultiHash: crc32(th+step1)", string(i), t)
			res += t
		}
		// fmt.Println(data, "MultiHash result:", string(i), res)
		out <- res
	}
}

/*
CombineResults получает все результаты, сортирует
(https://golang.org/pkg/sort/),
объединяет отсортированный результат через _
(символ подчеркивания) в одну строку
*/
func CombineResults(in, out chan interface{}) {
	res := make([]string, 0)
	for val := range in {
		data, ok := (val).(string)
		if !ok {
			fmt.Println("cant convert data to string")
		}
		res = append(res, data)
	}
	sort.Strings(res)
	s := strings.Join(res, "_")
	// fmt.Println("CombineResults", s)
	out <- s
}
