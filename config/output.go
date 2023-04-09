package config

import (
	"fmt"
	"log"
	"os"
)

type OutputAppender func(line string)

func CreateOutputAppender(filename string) (OutputAppender, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return nil, err
	}

	file.Close()

	return func(line string) {
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
	}, nil
}

func CreateOutputAppenderWithTruncate(filename string) (OutputAppender, error) {
	truncate(filename, 0644)

	return func(line string) {
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
	}, nil
}

func truncate(filename string, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("could not open file %q for truncation: %v", filename, err)
	}
	if err = f.Close(); err != nil {
		return fmt.Errorf("could not close file handler for %q after truncation: %v", filename, err)
	}
	return nil
}
