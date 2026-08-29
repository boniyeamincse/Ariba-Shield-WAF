package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestIdempotencyReplay verifies a retried POST with the same key returns the
// original response without re-invoking the handler (FR-0.1-043).
func TestIdempotencyReplay(t *testing.T) {
	store := NewIdempotencyStore(0)
	handlerCalls := 0

	h := Idempotency(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalls++
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"created"}`))
	}))

	do := func() *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/applications", strings.NewReader(`{}`))
		req.Header.Set("Idempotency-Key", "key-123")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		return rec
	}

	rec1 := do()
	if rec1.Code != http.StatusCreated {
		t.Fatalf("first call: expected 201, got %d", rec1.Code)
	}

	rec2 := do()
	if rec2.Code != http.StatusCreated {
		t.Fatalf("replay: expected 201, got %d", rec2.Code)
	}
	if handlerCalls != 1 {
		t.Fatalf("expected handler called once, got %d calls", handlerCalls)
	}
	if rec2.Header().Get("X-Idempotent-Replay") != "true" {
		t.Fatal("expected X-Idempotent-Replay header on replay")
	}
	body1, _ := io.ReadAll(rec1.Result().Body)
	body2, _ := io.ReadAll(rec2.Result().Body)
	if string(body1) != string(body2) {
		t.Fatal("replayed body must match original")
	}
}

// TestIdempotencyNoKeyPasses verifies requests without a key are not deduped.
func TestIdempotencyNoKeyPasses(t *testing.T) {
	store := NewIdempotencyStore(0)
	handlerCalls := 0

	h := Idempotency(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalls++
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodPost, "/x", strings.NewReader(`{}`))
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
	}
	if handlerCalls != 2 {
		t.Fatalf("expected 2 handler calls without key, got %d", handlerCalls)
	}
}

// TestIdempotencyGETNotAffected verifies GET requests bypass idempotency.
func TestIdempotencyGETNotAffected(t *testing.T) {
	store := NewIdempotencyStore(0)
	handlerCalls := 0

	h := Idempotency(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalls++
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Idempotency-Key", "k")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if handlerCalls != 1 {
		t.Fatalf("GET must not be deduped")
	}
}