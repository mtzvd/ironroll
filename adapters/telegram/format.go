package telegram

import (
	"fmt"

	"github.com/mtzvd/ironroll/core/roll"
)

// formatResult renders a single-line, minimal Ironsworn roll result.
//
// Format:
// 🎲 (action + modifier) vs (challenge1 & challenge2) → Ironsworn result
func formatResult(r roll.Result) string {
	return fmt.Sprintf(
		"🎲 (%d %+d) vs (%d & %d) → %s",
		r.ActionDie,
		r.Modifier,
		r.ChallengeDice[0],
		r.ChallengeDice[1],
		ironswornOutcome(r),
	)
}

// ironswornOutcome maps internal outcome categories
// to canonical Ironsworn terminology.
func ironswornOutcome(r roll.Result) string {
	switch r.Outcome {
	case roll.CriticalSuccess:
		return "Strong Hit (Match)"
	case roll.Success:
		return "Strong Hit"
	case roll.PartialSuccess:
		return "Weak Hit"
	case roll.Failure:
		return "Miss"
	case roll.CriticalFailure:
		return "Miss (Match)"
	default:
		return string(r.Outcome)
	}
}
