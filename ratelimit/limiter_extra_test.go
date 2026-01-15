package ratelimit

import (
	"net"
	"testing"
	"time"
)

func TestTemporaryBlockBranch(t *testing.T) {
	l := New(1, 1*time.Second, 500*time.Millisecond)
	ip := net.ParseIP("10.0.0.1")

	// First allowed
	if !l.Allow(ip) {
		t.Fatalf("expected first allowed")
	}

	// Second causes block (sets blockedAt)
	if l.Allow(ip) {
		t.Fatalf("expected second to set block and be denied")
	}

	// Immediate third should hit the temporary block check
	if l.Allow(ip) {
		t.Fatalf("expected immediate third to be denied by temporary block")
	}

	// After blockTime but before window expiry, still denied (count still > limit)
	time.Sleep(600 * time.Millisecond)
	if l.Allow(ip) {
		t.Fatalf("expected still denied after blockTime before window expiry")
	}

	// After window expiry, should be allowed
	time.Sleep(500 * time.Millisecond)
	if !l.Allow(ip) {
		t.Fatalf("expected allowed after window expiry")
	}
}

func TestAllowWhenCountEqualsLimit(t *testing.T) {
	l := New(2, 1*time.Second, 1*time.Second)
	ip := net.ParseIP("192.0.2.1")

	// First request allowed
	if !l.Allow(ip) {
		t.Fatalf("expected first allowed")
	}

	// Second request should be allowed (count == limit)
	if !l.Allow(ip) {
		t.Fatalf("expected second allowed when count equals limit")
	}

	// Third request should exceed limit and be blocked
	if l.Allow(ip) {
		t.Fatalf("expected third to be blocked")
	}
}
