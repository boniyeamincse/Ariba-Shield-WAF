package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"sync"
	"time"
)

// IdempotencyStore records completed requests by key so retries return the
// original response instead of repeating the mutation (SRS FR-0.1-043).
// In-memory TTL cache; a Redis-backed store is the Phase 8+ upgrade path.
type IdempotencyStore struct {
	mu    sync.Mutex
	items map[string]*idemEntry
	ttl   time.Duration
}

type idemEntry struct {
	status  int
	body    []byte
	created time.Time
}

// NewIdempotencyStore creates a store with a retention TTL.
func NewIdempotencyStore(ttl time.Duration) *IdempotencyStore {
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	return &IdempotencyStore{items: make(map[string]*idemEntry), ttl: ttl}
}

// Get returns a stored response for a key.
func (s *IdempotencyStore) Get(key string) (int, []byte, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.items[key]
	if !ok {
		return 0, nil, false
	}
	if time.Since(e.created) > s.ttl {
		delete(s.items, key)
		return 0, nil, false
	}
	return e.status, e.body, true
}

// Put stores a response for a key.
func (s *IdempotencyStore) Put(key string, status int, body []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[key] = &idemEntry{status: status, body: body, created: time.Now()}
}

// IdempotencyKey normalizes the Idempotency-Key header into a stable hash.
func IdempotencyKey(r *http.Request) (string, bool) {
	key := r.Header.Get("Idempotency-Key")
	if key == "" {
		return "", false
	}
	sum := sha256.Sum256([]byte(r.Method + "|" + r.URL.Path + "|" + key))
	return hex.EncodeToString(sum[:]), true
}

// idemWriter buffers the response so it can be replayed for retries.
type idemWriter struct {
	http.ResponseWriter
	status int
	body   []byte
}

func (w *idemWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *idemWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return w.ResponseWriter.Write(b)
}

// Idempotency returns middleware that deduplicates writes carrying an
// Idempotency-Key header. Must wrap handlers AFTER auth so actor context is set.
func Idempotency(store *IdempotencyStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only applies to mutating methods with a key present.
			switch r.Method {
			case "POST", "PUT", "PATCH", "DELETE":
			default:
				next.ServeHTTP(w, r)
				return
			}

			key, ok := IdempotencyKey(r)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			// Replay a completed request.
			if status, body, found := store.Get(key); found {
				w.Header().Set("X-Idempotent-Replay", "true")
				w.WriteHeader(status)
				_, _ = w.Write(body)
				return
			}

			// Run once, buffer the result.
			iw := &idemWriter{ResponseWriter: w, status: 200}
			next.ServeHTTP(iw, r)
			store.Put(key, iw.status, iw.body)
		})
	}
}