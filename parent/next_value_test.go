package parent

import (
	"amogus/config"
	"testing"
)

type testCase struct {
	input  string
	expect string
}

var testCases = []testCase{
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

	for _, testCase := range testCases {
		got := GetNextValue(cfg, testCase.input)
		if got != testCase.expect {
			t.Errorf("got %s, wanted %s", got, testCase.expect)
		}
	}
}
