package random

import (
	"math/rand"
	"time"
)

// RNG is the shared random number generator used by the application.
//
// It is initialized once at startup using a time-based seed.
// A local rand.Rand is used instead of the global generator
// to avoid deprecated rand.Seed usage and to ensure non-deterministic rolls.
var RNG = rand.New(rand.NewSource(time.Now().UnixNano()))
