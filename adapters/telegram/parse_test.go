package telegram

import "testing"

func TestParseModifierVariousInputs(t *testing.T) {
	cases := map[string]int{
		"":     0,
		"+2":   2,
		"2":    2,
		"-1":   -1,
		" +3 ": 3,
	}

	for input, want := range cases {
		got := parseModifier(input)
		if got != want {
			t.Fatalf("parseModifier(%q) = %d; want %d", input, got, want)
		}
	}

	// invalid numeric content -> treated as 0
	if parseModifier("bad") != 0 {
		t.Fatalf("parseModifier(" + "bad" + ") should return 0 on invalid input")
	}
}
