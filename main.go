package main

import (
	"amogus/parent"
	"fmt"
	"os"
)

const defaultConfigPath string = "amogus.yaml"
const defaultOutputPath string = "cracked"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		// TODO
		fmt.Printf("help is on the way!")
		os.Exit(1)
	}

	if len(os.Args) > 1 && os.Args[1] == "--test" {
		fmt.Printf("running pvm test program\n")
		TestPvm()
		os.Exit(0)
	}

	var hashesPath string
	var configPath string
	var outputPath string

	// 1 arg - hashes list is #1, rest is default
	if len(os.Args) > 1 {
		hashesPath = os.Args[1]
	} else {
		fmt.Println("path to input file was not supplied")
		os.Exit(1)
	}

	// 2 args - hashes list is #1, config path is #2, output path is default
	if len(os.Args) > 2 {
		configPath = os.Args[2]
	} else {
		configPath = defaultConfigPath
	}

	// 3 args - hashes list is #1, config path is #2, output path is #3
	if len(os.Args) > 3 {
		outputPath = os.Args[3]
	} else {
		outputPath = defaultOutputPath
	}

	fmt.Printf("input: %s, config: %s, output: %s\n", hashesPath, configPath, outputPath)

	err := parent.RunParent(hashesPath, configPath, outputPath)

	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}

}
