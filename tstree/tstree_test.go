package tstree

import "testing"

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
