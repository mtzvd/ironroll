package roll

// Outcome represents the final category of an Ironsworn roll.
// The string values are intended to be human-readable and
// suitable for direct display in chat environments.
type Outcome string

const (
	CriticalFailure Outcome = "Critical Failure"
	Failure         Outcome = "Failure"
	PartialSuccess  Outcome = "Partial Success"
	Success         Outcome = "Success"
	CriticalSuccess Outcome = "Critical Success"
)

// Result is the complete, explicit outcome of a single roll.
//
// It contains the full dice breakdown so that:
//   - the roll is transparent
//   - the result can be audited by humans
//   - no hidden logic exists outside this structure
//
// No additional interpretation or consequence is applied.
type Result struct {
	ActionDie     int     // Result of the 1d6 action die
	Modifier      int     // Applied modifier (Z)
	ChallengeDice [2]int  // Results of the two 1d10 challenge dice
	Total         int     // ActionDie + Modifier
	Outcome       Outcome // Final outcome category
}
