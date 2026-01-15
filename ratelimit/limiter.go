package ratelimit

import (
	"net"
	"sync"
	"time"
)

// Limiter implements a simple in-memory IP-based rate limiter.
//
// This limiter is best-effort and instance-local.
// It is designed to protect against casual abuse,
// not to provide strict global guarantees.
//
// In stateless cloud environments:
//   - counters reset on cold start
//   - limits apply per instance
//
// This tradeoff is intentional and documented.
type Limiter struct {
	mu        sync.Mutex
	visitors  map[string]*visitor
	limit     int
	window    time.Duration
	blockTime time.Duration
}

type visitor struct {
	count     int
	expiresAt time.Time
	blockedAt time.Time
}

// New creates a new Limiter.
//
// limit: number of requests allowed per window
// window: time window for counting requests
// blockTime: duration of temporary blocking after limit is exceeded
func New(limit int, window, blockTime time.Duration) *Limiter {
	return &Limiter{
		visitors:  make(map[string]*visitor),
		limit:     limit,
		window:    window,
		blockTime: blockTime,
	}
}

// Allow reports whether a request from the given IP should be allowed.
func (l *Limiter) Allow(ip net.IP) bool {
	now := time.Now()
	key := ip.String()

	l.mu.Lock()
	defer l.mu.Unlock()

	v, exists := l.visitors[key]
	if !exists {
		l.visitors[key] = &visitor{
			count:     1,
			expiresAt: now.Add(l.window),
		}
		return true
	}

	// Check temporary block
	if !v.blockedAt.IsZero() && now.Sub(v.blockedAt) < l.blockTime {
		return false
	}

	// Reset window if expired
	if now.After(v.expiresAt) {
		v.count = 1
		v.expiresAt = now.Add(l.window)
		v.blockedAt = time.Time{}
		return true
	}

	v.count++
	if v.count > l.limit {
		v.blockedAt = now
		return false
	}

	return true
}
