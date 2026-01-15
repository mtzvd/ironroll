package telegram

import (
	"strings"
	"testing"

	"github.com/mtzvd/ironroll/core/roll"
)

func TestFormatResultContainsFields(t *testing.T) {
	r := roll.Result{
		ActionDie:     4,
		Modifier:      -1,
		ChallengeDice: [2]int{7, 2},
		Total:         3,
		Outcome:       roll.Failure,
	}

	s := formatResult(r)

	// formatResult returns compact format: 🎲 (action +mod) vs (c1 & c2) → Outcome
	checks := []string{
		"🎲",
		"4",   // action die
		"-1",  // modifier
		"7",   // challenge die 1
		"2",   // challenge die 2
		"vs",
		"Miss", // outcome (Failure maps to Miss)
	}

	for _, c := range checks {
		if !strings.Contains(s, c) {
			t.Fatalf("formatted result missing %q in %q", c, s)
		}
	}
}
