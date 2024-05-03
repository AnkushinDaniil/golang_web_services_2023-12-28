package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"hw3/user"
)

type User = user.User

func FastSearch(out io.Writer) {
	out.Write([]byte("found users:\n"))

	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)

	var (
		notSeenBefore, isAndroid, isMSIE bool
		browser                          string
		data                             []byte
	)

	seenBrowsers := make([]string, 0, 128)
	uniqueBrowsers := 0

	help := func() {
		notSeenBefore = true
		for _, item := range seenBrowsers {
			if item == browser {
				notSeenBefore = false
			}
		}
		if notSeenBefore {
			seenBrowsers = append(seenBrowsers, browser)
			uniqueBrowsers++
		}
	}

	var user User

	i := 0
	for scanner.Scan() {
		data = scanner.Bytes()
		user.UnmarshalJSON(data)
		// err := json.Unmarshal(data, &user)
		isAndroid, isMSIE = false, false

		for _, browser = range user.Browsers {
			if strings.Contains(browser, "Android") {
				isAndroid = true
				help()
			}
			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				help()
			}
		}

		if !(isAndroid && isMSIE) {
			i++
			continue
		}

		out.Write([]byte{'['})
		out.Write([]byte(strconv.Itoa(i)))
		out.Write([]byte{']', ' '})
		out.Write([]byte(user.Name))
		out.Write([]byte{' ', '<'})
		out.Write([]byte(strings.ReplaceAll(user.Email, "@", " [at] ")))
		out.Write([]byte{'>', '\n'})

		i++
	}

	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}
