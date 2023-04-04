package config

import (
	"amogus/tstree"
	"bufio"
	"fmt"
	"os"
	"regexp"
)

type Hashes struct {
	//InputRawLines []string
	Lut *tstree.LookupTable
}

func ReadHashesFile(filename string, mode Mode) (*Hashes, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	result := &Hashes{
		Lut: tstree.BuildLookupTableFromLines(lines),
	}

	if compatible, incompatibleLine := isHashesFileCompatible(lines, mode); !compatible {
		return nil, fmt.Errorf("hash '%s' incompatible with mode %s or mode unsupported", incompatibleLine, mode)
	}

	return result, scanner.Err()
}

func isHashesFileCompatible(lines []string, mode Mode) (bool, string) {
	for _, i := range lines {
		if mode == Sha512 {
			if !isLineCompatibleSha512(i) {
				return false, i
			}
		} else {
			// no other algos supported
			return false, i
		}
	}

	return true, ""
}

func isLineCompatibleSha512(line string) bool {
	m, _ := regexp.MatchString("^\\w{128}$", line)

	return m
}
