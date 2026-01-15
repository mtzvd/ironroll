package ratelimit

import (
	"net"
	"testing"
	"time"
)

func TestLimiterBasicFlow(t *testing.T) {
	l := New(1, 40*time.Millisecond, 20*time.Millisecond)
	ip := net.ParseIP("127.0.0.1")

	// First request allowed
	if !l.Allow(ip) {
		t.Fatalf("expected first request to be allowed")
	}

	// Second request within window should be blocked
	if l.Allow(ip) {
		t.Fatalf("expected second request to be blocked")
	}

	// After window expires, request should be allowed again
	time.Sleep(60 * time.Millisecond)
	if !l.Allow(ip) {
		t.Fatalf("expected request after window to be allowed")
	}
}
