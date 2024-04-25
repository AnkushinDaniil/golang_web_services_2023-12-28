package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
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
			defer wg.Done()
			defer close(out)
			j(in, out)
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
	quotaChannel := make(chan struct{}, quotaLimit)

	for value := range in {
		wg.Add(1)
		go func(val interface{}, quotaCh chan struct{}) {
			data := fmt.Sprintf("%v", val)
			md5 := make(chan string)

			go func(m chan<- string, qCh chan struct{}, d string) {
				qCh <- struct{}{}
				m <- DataSignerMd5(d)
				<-qCh
			}(md5, quotaCh, data)

			go func(d string, m <-chan string) {
				crc32 := make(chan string)
				go func(c chan<- string, s string) {
					c <- DataSignerCrc32(d)
				}(crc32, d)

				crc32_md5 := make(chan string)
				go func(c chan<- string, mCh <-chan string) {
					c <- DataSignerCrc32(<-mCh)
				}(crc32_md5, m)

				res := fmt.Sprintf("%v~%v", <-crc32, <-crc32_md5)

				// fmt.Println("SingleHash", d)
				// fmt.Println("SingleHash md5(data)", m)
				// fmt.Println("SingleHash crc32(md5(data))", crc32_md5)
				// fmt.Println("SingleHash crc32(data)", crc32)
				// fmt.Println("SingleHash result", res)

				out <- res
				wg.Done()
			}(data, md5)
		}(value, quotaChannel)
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
	wgMain := &sync.WaitGroup{}
	defer wgMain.Wait()
	for value := range in {
		wgMain.Add(1)
		go func(val interface{}) {
			defer wgMain.Done()
			data, ok := (val).(string)
			if !ok {
				fmt.Println("cant convert data to string")
			}

			arr := make([]string, 6)
			wg := &sync.WaitGroup{}
			mu := &sync.Mutex{}

			for i := 0; i < 6; i++ {
				wg.Add(1)
				go func(j int, d string) {
					defer wg.Done()
					t := DataSignerCrc32(strconv.Itoa(j) + d)
					mu.Lock()
					arr[j] = t
					mu.Unlock()
					// fmt.Println(data, "MultiHash: crc32(th+step1)", j, t)
				}(i, data)
			}
			wg.Wait()
			res := strings.Join(arr, "")
			// fmt.Println(data, "MultiHash result:", res)
			out <- res
		}(value)
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
