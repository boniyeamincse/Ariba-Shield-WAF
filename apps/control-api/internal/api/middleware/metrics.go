package middleware

import (
	"sync"
	"time"
)

// requestDurationHistogram is a Prometheus-style cumulative histogram that
// records HTTP request durations in milliseconds (P2.26). It is observed by the
// Logging middleware and rendered by the handlers.Metrics endpoint.
var requestDurationHistogram = newDurationHistogram([]float64{
	0.5, 1, 2.5, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000,
})

// durationHistogram holds cumulative per-bucket counts, total count and sum.
type durationHistogram struct {
	mu     sync.Mutex
	bounds []float64 // upper bounds in ms, ascending
	counts []uint64  // cumulative count per bucket (le = bounds[i])
	count  uint64
	sum    float64
}

func newDurationHistogram(bounds []float64) *durationHistogram {
	return &durationHistogram{bounds: append([]float64(nil), bounds...), counts: make([]uint64, len(bounds))}
}

func (h *durationHistogram) observe(ms float64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.count++
	h.sum += ms
	for i, b := range h.bounds {
		if ms <= b {
			h.counts[i]++
		}
	}
}

// ObserveRequestDuration records a single request duration for /metrics.
func ObserveRequestDuration(d time.Duration) {
	requestDurationHistogram.observe(float64(d) / float64(time.Millisecond))
}

// RequestDurationHistogram returns a snapshot of the histogram: bucket upper
// bounds (ms), cumulative bucket counts, total count, and total sum (ms).
func RequestDurationHistogram() (bounds []float64, counts []uint64, count uint64, sum float64) {
	h := requestDurationHistogram
	h.mu.Lock()
	defer h.mu.Unlock()
	return append([]float64(nil), h.bounds...), append([]uint64(nil), h.counts...), h.count, h.sum
}
