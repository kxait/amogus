package main

import (
	"amogus/child"
	"amogus/parent"
	"fmt"
	"os"
	"path/filepath"
)

const defaultConfigPath string = "amogus.yaml"
const defaultOutputPath string = "cracked"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		usage(false)
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "--test" {
		fmt.Printf("running pvm test program\n")
		TestPvm()
		os.Exit(0)
	}

	if len(os.Args) > 1 && os.Args[1] == "--child" {
		fmt.Printf("starting child\n")
		err := child.RunChild()
		if err != nil {
			panic(err)
			//fmt.Printf("ERROR: %s\n", err.Error())
			//os.Exit(1)
			//debug.PrintStack()
		}
		os.Exit(0)
	}

	var hashesPath string
	var configPath string
	var outputPath string

	// 1 arg - hashes list is #1, rest is default
	if len(os.Args) > 1 {
		hashesPath = os.Args[1]
	} else {
		usage(true)
		return
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

	err := parent.RunParent(hashesPath, configPath, outputPath)

	if err != nil {
		//fmt.Printf("ERROR: %s\n", err.Error())
		//debug.PrintStack()
		panic(err)
		//os.Exit(1)
	}

}

func usage(err bool) {
	progName := filepath.Base(os.Args[0])
	fmt.Printf("usage: %s hashes_path [config_path] [output_path]\n", progName)
	fmt.Println("github.com/kxait/amogus")

	var exitCode int
	if err {
		exitCode = -1
	} else {
		exitCode = 0
	}

	os.Exit(exitCode)
}
