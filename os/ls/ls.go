package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: ls directory_name")
		os.Exit(1)
	}

	files, err := ioutil.ReadDir(os.Args[1])
	checkErr(err)

	for _, f := range files {
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}

		if f.IsDir() {
			fmt.Printf("\033[34m%s\033[0m\n", f.Name())
		} else {
			fmt.Println(f.Name())
		}
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}