package roll

import (
	"reflect"
	"testing"
)

func TestSeedProducesRepeatableSequences(t *testing.T) {
	// Seed twice with same value and ensure sequences repeat
	Seed(2024)
	a := Roll(0)
	b := Roll(1)

	Seed(2024)
	a2 := Roll(0)
	b2 := Roll(1)

	if !reflect.DeepEqual(a, a2) || !reflect.DeepEqual(b, b2) {
		t.Fatalf("Seed did not produce repeatable sequences: got %v,%v vs %v,%v", a, b, a2, b2)
	}

	// Reset to default RNG and ensure calls still succeed
	ResetRand()
	_ = Roll(0)
}
