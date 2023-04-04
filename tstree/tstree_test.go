package tstree

import (
	"testing"
)

func TestSingleKey(t *testing.T) {
	const text string = "the quick brown fox jumps over the lazy dog"

	lut := make(LookupTable)

	lut.AppendLine(text)

	has := lut.Has(text)

	if !has {
		t.Error("expected the key to exist, but it did, in fact, not")
	}
}

func TestTwoCompletelyDifferentKeys(t *testing.T) {
	const text1 string = "the quick brown fox jumps over the lazy dog"
	const text2 string = "ala ma kota i pies"

	lut := make(LookupTable)

	lut.AppendLine(text1)
	lut.AppendLine(text2)

	has := lut.Has(text1) && lut.Has(text2)

	if !has {
		t.Error("expected the keys to exist, but some of them did, in fact, not")
	}
}

func TestTwoOverlappingKeys(t *testing.T) {
	const text1 string = "the quick brown fox jumps over the lazy dog"
	const text2 string = "the quick brown fox jumped over the lazy dog and ate the bone"

	lut := make(LookupTable)

	lut.AppendLine(text1)
	lut.AppendLine(text2)

	has := lut.Has(text1) && lut.Has(text2)

	if !has {
		t.Error("expected the keys to exist, but some of them did, in fact, not")
	}
}

// func BenchmarkLarge(b *testing.B) {
// 	chars := "abcdefghijklmnoprstuvwxyz012345678"

// 	const samples int = 10000

// 	lut := make(LookupTable)
// 	for i := 0; i < samples; i++ {
// 		var c string
// 		for j := 0; j < 128; j++ {
// 			c = fmt.Sprintf("%s%c", c, chars[rand.Intn(len(chars))])
// 		}
// 		lut.AppendLine(c)
// 	}
// }

// func TestSerializes(t *testing.T) {
// 	file, _ := os.Open("/home/kx/src/prir/proj/hashes2")
// 	defer file.Close()

// 	var lines []string
// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		lines = append(lines, scanner.Text())
// 	}

// 	lut := make(LookupTable)
// 	for _, line := range lines {
// 		lut.AppendLine(line)
// 	}

// 	jsonData, err := json.Marshal(lut)
// 	if err != nil {
// 		t.Errorf("marshal error: %s", err)
// 	}

// 	os.WriteFile("/home/kx/src/prir/proj/json.json", jsonData, os.ModeAppend)
// }
