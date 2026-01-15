package roll

import "math/rand"

// Package-level `intn` function is used for randomness and can be
// replaced in tests via `SetRand`/`ResetRand` to make behavior
// deterministic. By default it delegates to the global `rand.Intn`.
var intn = rand.Intn

// SetRand replaces the package RNG with the provided *rand.Rand.
// This makes `Roll` deterministic when the same seed is used.
func SetRand(r *rand.Rand) {
	intn = r.Intn
}

// ResetRand restores the default global RNG behavior.
func ResetRand() {
	intn = rand.Intn
}

// Seed is a convenience that sets the RNG to a new source derived
// from the provided seed.
func Seed(seed int64) {
	SetRand(rand.New(rand.NewSource(seed)))
}

// Roll performs a single Ironsworn roll with the given modifier.
//
// This function is pure from the caller's perspective:
//   - it does not panic
//   - it does not log
//   - it does not allocate unnecessary resources
//   - it does not depend on any external context
//
// The returned Result contains the full dice breakdown
// and the final outcome category.
func Roll(modifier int) Result {
	// Roll the action die (1d6).
	actionDie := intn(6) + 1

	// Roll the two challenge dice (1d10 each).
	challenge := [2]int{
		intn(10) + 1,
		intn(10) + 1,
	}

	// Calculate the action score.
	total := actionDie + modifier

	// Determine the final outcome.
	outcome := determineOutcome(total, challenge)

	return Result{
		ActionDie:     actionDie,
		Modifier:      modifier,
		ChallengeDice: challenge,
		Total:         total,
		Outcome:       outcome,
	}
}
