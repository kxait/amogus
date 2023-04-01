package config

import (
	"fmt"
	"log"
	"os"
)

func CreateOutputAppender(filename string) (f func(line string)) {
	f = func(line string) {
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := file.Write([]byte(fmt.Sprintf("%s\n", line))); err != nil {
			log.Fatal(err)
		}
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}

	return
}
