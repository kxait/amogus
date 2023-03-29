package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "test" {
		fmt.Printf("running pvm test program\n")
		TestPvm()
		os.Exit(0)
	}
}
