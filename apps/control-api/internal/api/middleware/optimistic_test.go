package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestIfMatchHeader verifies header parsing.
func TestIfMatchHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodPatch, "/x", nil)
	req.Header.Set("If-Match", `"3"`)
	v, ok := IfMatch(req)
	if !ok || v != 3 {
		t.Fatalf("expected version 3, got %d ok=%v", v, ok)
	}

	req2 := httptest.NewRequest(http.MethodPatch, "/x", nil)
	req2.Header.Set("X-Expected-Version", "7")
	v2, ok2 := IfMatch(req2)
	if !ok2 || v2 != 7 {
		t.Fatalf("expected version 7, got %d ok=%v", v2, ok2)
	}

	req3 := httptest.NewRequest(http.MethodPatch, "/x", nil)
	if _, ok3 := IfMatch(req3); ok3 {
		t.Fatal("no header should report not-present")
	}
}

// TestOptimisticConcurrencyNoHeaderPasses verifies mutations without a version
// header are not blocked by the middleware (delegated to the handler).
func TestOptimisticConcurrencyNoHeaderPasses(t *testing.T) {
	var hit bool
	h := OptimisticConcurrency(nil, "applications", "id", func(w http.ResponseWriter, r *http.Request) {
		hit = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/applications/abc", nil)
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if !hit {
		t.Fatal("handler should run when no version header present")
	}
}