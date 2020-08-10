package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	var input string
	for {
		fmt.Print("$: ")
		n, err := fmt.Scan(&input)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		if n < 0 {
			break
		}

		cmd := exec.Command(input)

		var out = new(bytes.Buffer)
		cmd.Stdout = out
		cmd.Stderr = out

		err = cmd.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		fmt.Print(out.String())
	}
}