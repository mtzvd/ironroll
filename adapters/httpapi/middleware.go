package httpapi

import (
	"net"
	"net/http"

	"github.com/mtzvd/ironroll/ratelimit"
)

// RateLimitMiddleware wraps an HTTP handler with IP-based rate limiting.
func RateLimitMiddleware(limiter *ratelimit.Limiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)
		if !limiter.Allow(ip) {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// clientIP extracts the client IP address from the request.
//
// This implementation is intentionally simple and conservative.
// It does not trust forwarded headers by default.
func clientIP(r *http.Request) net.IP {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return net.IPv4(0, 0, 0, 0)
	}
	return net.ParseIP(host)
}
