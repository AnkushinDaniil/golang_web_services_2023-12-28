package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	var write func(prefix, path string)

	write = func(prefix, dir string) {
		file, _ := os.Open(dir)
		files, _ := file.Readdir(0)

		if !printFiles {
			n := 0
			for _, file := range files {
				if file.IsDir() {
					files[n] = file
					n++
				}
			}

			files = files[:n]
		}

		if len(files) == 0 {
			return
		}

		sort.Slice(files, func(i, j int) bool {
			return files[i].Name() < files[j].Name()
		})

		for i := 0; i < len(files)-1; i++ {
			if files[i].IsDir() {
				fmt.Fprintf(out, "%v├───%v\n", prefix, files[i].Name())
				write("│	"+prefix, fmt.Sprintf("%v/%v", dir, files[i].Name()))
			} else if printFiles {
				sizeStr := ""
				if sizeInt64 := files[i].Size(); sizeInt64 == 0 {
					sizeStr = "empty"
				} else {
					sizeStr = strconv.FormatInt(files[i].Size(), 10) + "b"
				}
				fmt.Fprintf(out, "%v├───%v (%v)\n", prefix, files[i].Name(), sizeStr)
			}
		}
		lastFile := files[len(files)-1]
		if lastFile.IsDir() {
			fmt.Fprintf(out, "%v└───%v\n", prefix, lastFile.Name())
			write(prefix+"	", fmt.Sprintf("%v/%v", dir, lastFile.Name()))
		} else if printFiles {
			sizeStr := ""
			if sizeInt64 := lastFile.Size(); sizeInt64 == 0 {
				sizeStr = "empty"
			} else {
				sizeStr = strconv.FormatInt(lastFile.Size(), 10) + "b"
			}
			fmt.Fprintf(out, "%v└───%v (%v)\n", prefix, lastFile.Name(), sizeStr)
		}
	}

	write("", path)

	return nil
}
