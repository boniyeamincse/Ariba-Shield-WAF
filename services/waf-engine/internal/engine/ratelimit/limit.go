package ratelimit

import (
	"sync"
	"time"
)

// SlidingWindow is a per-key fixed-window rate limiter with a sliding expiry,
// suitable for per-IP / per-route limits (master plan §6.6). Uses a mutex +
// map; good enough for a single gateway process in Release 0.1.
//
// P2.20: buckets are evicted once their window has fully elapsed so the map
// does not grow unboundedly with unique client IPs.
type SlidingWindow struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	buckets map[string]*bucket
	// maxKeys bounds the map; the oldest entries are swept when exceeded.
	maxKeys int
	// lastSweep is used to amortise eviction rather than sweeping every call.
	lastSweep time.Time
}

type bucket struct {
	count       int
	windowStart time.Time
}

// New creates a rate limiter allowing `limit` events per `window`.
func New(limit int, window time.Duration) *SlidingWindow {
	return NewWithMax(limit, window, 100_000)
}

// NewWithMax creates a rate limiter with an explicit bucket-map upper bound.
func NewWithMax(limit int, window time.Duration, maxKeys int) *SlidingWindow {
	if limit <= 0 {
		limit = 1
	}
	if window <= 0 {
		window = time.Minute
	}
	if maxKeys <= 0 {
		maxKeys = 100_000
	}
	return &SlidingWindow{
		limit:     limit,
		window:    window,
		buckets:   make(map[string]*bucket),
		maxKeys:   maxKeys,
		lastSweep: time.Now(),
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

	// Amortised eviction: sweep at most once per window to bound the map.
	if now.Sub(r.lastSweep) >= r.window || len(r.buckets) > r.maxKeys {
		r.sweep(now)
	}

	return b.count <= r.limit, remaining
}

// Len returns the current number of tracked buckets (for tests).
func (r *SlidingWindow) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.buckets)
}

// sweep removes buckets whose window has elapsed. When the map exceeds
// maxKeys it also evicts the oldest entries (FIFO) to bound memory.
func (r *SlidingWindow) sweep(now time.Time) {
	r.lastSweep = now
	over := len(r.buckets) - r.maxKeys
	for k, b := range r.buckets {
		if over <= 0 && now.Sub(b.windowStart) < r.window {
			continue
		}
		delete(r.buckets, k)
		over--
	}
}