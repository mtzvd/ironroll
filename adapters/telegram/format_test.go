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

	checks := []string{
		"Action Die",
		"Modifier",
		"Challenge Dice",
		"Total",
		"Outcome",
		"4",
		"-1",
		"7",
		"2",
		"3",
		"Failure",
	}

	for _, c := range checks {
		if !strings.Contains(s, c) {
			t.Fatalf("formatted result missing %q in %q", c, s)
		}
	}
}
