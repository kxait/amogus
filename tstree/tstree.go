package tstree

import (
	"bufio"
	"os"
)

type LookupKey byte

type LookupTable map[LookupKey]*LookupTable

func (lut *LookupTable) AppendLine(line string) {
	bytes := []byte(line)
	lut.appendBytes(bytes)
}

func (lut *LookupTable) Has(text string) bool {
	bytes := []byte(text)

	thisLut := lut
	for _, b := range bytes {
		val, ok := (*thisLut)[LookupKey(b)]

		if !ok {
			return false
		}

		thisLut = val
	}

	return len(*thisLut) == 0
}

func (lut *LookupTable) appendBytes(letters []byte) {
	thisLetter := letters[0]
	lut.appendLetter(thisLetter)

	if len(letters) > 1 {
		(*lut)[LookupKey(thisLetter)].appendBytes(letters[1:])
	}
}

func (lut *LookupTable) appendLetter(letter byte) {
	_, ok := (*lut)[LookupKey(letter)]

	if !ok {
		newLut := make(LookupTable)
		(*lut)[LookupKey(letter)] = &newLut
	}
}

func BuildLookupTableFromFile(filename string) (*LookupTable, error) {
	result := make(LookupTable)

	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanBytes)

	for scanner.Scan() {
		result.AppendLine(scanner.Text())
	}

	return &result, nil
}

func BuildLookupTableFromLines(lines []string) *LookupTable {
	result := make(LookupTable)

	for _, v := range lines {
		result.AppendLine(v)
	}

	return &result
}
