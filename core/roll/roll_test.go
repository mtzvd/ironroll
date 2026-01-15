package roll

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestDetermineOutcomeTable(t *testing.T) {
	cases := []struct {
		name      string
		total     int
		challenge [2]int
		want      Outcome
	}{
		{"CriticalSuccess", 3, [2]int{2, 2}, CriticalSuccess},
		{"Success", 3, [2]int{1, 2}, Success},
		{"PartialSuccess", 6, [2]int{5, 7}, PartialSuccess},
		{"CriticalFailure", 5, [2]int{8, 8}, CriticalFailure},
		{"Failure", 5, [2]int{9, 10}, Failure},
		{"EqualityCountsAsLoss", 5, [2]int{5, 6}, Failure},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := determineOutcome(c.total, c.challenge)
			if got != c.want {
				t.Fatalf("%s: determineOutcome(%d, %v) = %q, want %q", c.name, c.total, c.challenge, got, c.want)
			}
		})
	}
}

func TestRollProducesConsistentResult(t *testing.T) {
	// Use the package-level controllable RNG to make this deterministic.
	rsrc := rand.New(rand.NewSource(42))
	SetRand(rsrc)
	defer ResetRand()

	modifier := 2
	r := Roll(modifier)

	if r.Total != r.ActionDie+r.Modifier {
		t.Fatalf("Total mismatch: got %d, want %d", r.Total, r.ActionDie+r.Modifier)
	}

	if r.ActionDie < 1 || r.ActionDie > 6 {
		t.Fatalf("ActionDie out of range: %d", r.ActionDie)
	}

	for i, c := range r.ChallengeDice {
		if c < 1 || c > 10 {
			t.Fatalf("ChallengeDice[%d] out of range: %d", i, c)
		}
	}

	// Outcome must equal the pure deterministic calculation.
	want := determineOutcome(r.Total, r.ChallengeDice)
	if r.Outcome != want {
		t.Fatalf("Outcome mismatch: got %q, want %q", r.Outcome, want)
	}
}

func TestRollDeterministicRepeatable(t *testing.T) {
	// Create two RNGs with the same seed and ensure sequences match
	r1 := rand.New(rand.NewSource(12345))
	r2 := rand.New(rand.NewSource(12345))

	SetRand(r1)
	defer ResetRand()

	var seq1 []Result
	for i := 0; i < 20; i++ {
		seq1 = append(seq1, Roll(i%4-2)) // some varying modifiers
	}

	SetRand(r2)
	var seq2 []Result
	for i := 0; i < 20; i++ {
		seq2 = append(seq2, Roll(i%4-2))
	}

	if !reflect.DeepEqual(seq1, seq2) {
		t.Fatalf("deterministic sequences do not match")
	}
}
