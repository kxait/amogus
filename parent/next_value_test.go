package parent_test

import (
	"amogus/config"
	"amogus/parent"
	"testing"
)

type testCaseSingle struct {
	input  string
	expect string
}

var testCasesSingle = []testCaseSingle{
	{"a", "b"},
	{"b", "c"},
	{"z", "aa"},
	{"aa", "ab"},
	{"abc", "abd"},
	{"zzz", "aaaa"},
	{"az", "ba"},
	{"aza", "azb"},
	{"azz", "baa"},
	{"aaaaazzzzz", "aaaabaaaaa"},
	{"zzzzzzzzzzzzzzzzzzzzzzzz", ""},
}

func TestNextValue(t *testing.T) {
	cfg := &config.AmogusConfig{
		Characters: "abcdefghijklmnoprstuvwxyz",
	}

	for _, testCase := range testCasesSingle {
		got := parent.GetNextValue(cfg, testCase.input)
		if got != testCase.expect {
			t.Errorf("got %s, wanted %s", got, testCase.expect)
		}
	}
}

type testCaseOffset struct {
	input  string
	offset int64
	expect string
}

var testCasesOffset = []testCaseOffset{
	{"a", 1, "b"},
	{"a", 2, "c"},
	{"zzz", 1000, "aboz"},
}

func TestNextValueOffset(t *testing.T) {
	cfg := &config.AmogusConfig{
		Characters: "abcdefghijklmnoprstuvwxyz",
	}

	for _, testCase := range testCasesOffset {
		got := parent.GetNextValueOffset(cfg, testCase.input, testCase.offset)
		if got != testCase.expect {
			t.Errorf("got %s, wanted %s", got, testCase.expect)
		}
	}
}
