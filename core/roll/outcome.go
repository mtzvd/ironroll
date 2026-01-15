package roll

// outcomeKey represents the normalized state used to determine
// the final outcome of a roll.
//
// Instead of mapping raw dice values directly, we first reduce
// the roll to two meaningful facts:
//
//   - wins: how many challenge dice were beaten (0, 1, or 2)
//   - isDouble: whether both challenge dice show the same value
//
// This mirrors the Ironsworn rules directly and keeps the logic
// easy to reason about and audit.
type outcomeKey struct {
	wins     int
	isDouble bool
}

// outcomeTable maps the normalized roll state to the final outcome.
//
// This table exactly matches the canonical Ironsworn outcome rules:
//
//	wins = 2, double => Critical Success
//	wins = 2         => Success
//	wins = 1         => Partial Success
//	wins = 0, double => Critical Failure
//	wins = 0         => Failure
var outcomeTable = map[outcomeKey]Outcome{
	{wins: 2, isDouble: true}:  CriticalSuccess,
	{wins: 2, isDouble: false}: Success,
	{wins: 1, isDouble: false}: PartialSuccess,
	{wins: 0, isDouble: true}:  CriticalFailure,
	{wins: 0, isDouble: false}: Failure,
}

// determineOutcome calculates the final outcome category
// based on the action score and the two challenge dice.
//
// This function contains no randomness and no side effects.
// It implements the Ironsworn comparison rules exactly.
func determineOutcome(total int, challenge [2]int) Outcome {
	wins := 0

	// Compare the action score against each challenge die.
	// Equality counts as a loss (canonical Ironsworn rule).
	for _, c := range challenge {
		if total > c {
			wins++
		}
	}

	isDouble := challenge[0] == challenge[1]

	key := outcomeKey{
		wins:     wins,
		isDouble: isDouble,
	}

	// The table lookup is guaranteed to succeed for all valid states.
	return outcomeTable[key]
}
