package httpapi

import (
	"testing"

	"github.com/mtzvd/ironroll/core/roll"
)

func TestFormatResult(t *testing.T) {
	r := roll.Result{
		ActionDie:     3,
		Modifier:      2,
		ChallengeDice: [2]int{5, 5},
		Total:         5,
		Outcome:       roll.PartialSuccess,
	}

	got := formatResult(r)
	if got.ActionDie != r.ActionDie || got.Modifier != r.Modifier || got.Total != r.Total {
		t.Fatalf("formatResult mismatch: got %+v, want %+v", got, r)
	}
	if got.Outcome != string(r.Outcome) {
		t.Fatalf("formatResult outcome mismatch: got %q want %q", got.Outcome, r.Outcome)
	}
}
