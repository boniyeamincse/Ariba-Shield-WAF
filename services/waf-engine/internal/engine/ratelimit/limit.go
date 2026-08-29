package ratelimit

import (
	"sync"
	"time"
)

// SlidingWindow is a per-key fixed-window rate limiter with a sliding expiry,
// suitable for per-IP / per-route limits (master plan §6.6). Uses a mutex +
// map; good enough for a single gateway process in Release 0.1.
type SlidingWindow struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	buckets map[string]*bucket
}

type bucket struct {
	count    int
	windowStart time.Time
}

// New creates a rate limiter allowing `limit` events per `window`.
func New(limit int, window time.Duration) *SlidingWindow {
	if limit <= 0 {
		limit = 1
	}
	if window <= 0 {
		window = time.Minute
	}
	return &SlidingWindow{
		limit:   limit,
		window:  window,
		buckets: make(map[string]*bucket),
	}
}

// Allow reports whether a request for key is within the limit.
// It returns (allowed, remaining).
func (r *SlidingWindow) Allow(key string) (bool, int) {
	now := time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()

	b, ok := r.buckets[key]
	if !ok || now.Sub(b.windowStart) >= r.window {
		// New window: reset.
		b = &bucket{count: 0, windowStart: now}
		r.buckets[key] = b
	}
	b.count++
	remaining := r.limit - b.count
	if remaining < 0 {
		remaining = 0
	}
	return b.count <= r.limit, remaining
}
