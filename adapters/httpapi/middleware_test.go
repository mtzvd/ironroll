package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mtzvd/ironroll/ratelimit"
)

func TestRateLimitMiddleware_AllowsAndBlocks(t *testing.T) {
	limiter := ratelimit.New(1, 50*time.Millisecond, 20*time.Millisecond)

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	mw := RateLimitMiddleware(limiter, next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "127.0.0.1:1234"

	// First request should pass
	rw1 := httptest.NewRecorder()
	mw.ServeHTTP(rw1, req)
	if rw1.Code != http.StatusOK || !called {
		t.Fatalf("expected first request to be allowed")
	}

	// Second request should be blocked
	called = false
	rw2 := httptest.NewRecorder()
	mw.ServeHTTP(rw2, req)
	if rw2.Code != http.StatusTooManyRequests {
		t.Fatalf("expected second request to be rate limited, got %d", rw2.Code)
	}

	// After window expires, next request should be allowed
	time.Sleep(60 * time.Millisecond)
	called = false
	rw3 := httptest.NewRecorder()
	mw.ServeHTTP(rw3, req)
	if rw3.Code != http.StatusOK || !called {
		t.Fatalf("expected request after window to be allowed")
	}
}

func TestClientIPFallbackWhenRemoteAddrInvalid(t *testing.T) {
	limiter := ratelimit.New(1, 50*time.Millisecond, 20*time.Millisecond)

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	mw := RateLimitMiddleware(limiter, next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "invalid-remote-addr"

	// First request should pass (uses 0.0.0.0 as IP)
	rw := httptest.NewRecorder()
	mw.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK || !called {
		t.Fatalf("expected handler to run when RemoteAddr invalid")
	}
}
