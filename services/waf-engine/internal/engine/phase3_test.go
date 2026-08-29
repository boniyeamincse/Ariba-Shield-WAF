package engine

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newTestEngine(t *testing.T, backendURL string, detectOnly bool) *Engine {
	t.Helper()
	eng, err := New(Config{
		BackendURL:   backendURL,
		RulesPath:    "../../../../rules/core/baseline.conf",
		DetectOnly:   detectOnly,
		BlockStatus:  403,
		RequestIDHdr: "X-Request-ID",
		BlockTitle:   "Blocked",
		BlockMessage: "Security policy",
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return eng
}

func TestIPBlockList(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	eng := newTestEngine(t, backend.URL, false)
	if err := eng.SetIPLists(nil, []string{"192.0.2.0/24"}); err != nil {
		t.Fatalf("SetIPLists: %v", err)
	}

	// Blocked IP must get 403 without reaching the backend.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.0.2.50:1234"
	rec := httptest.NewRecorder()
	eng.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for blocked IP, got %d", rec.Code)
	}

	// Allowed IP passes through.
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	req2.RemoteAddr = "198.51.100.10:1234"
	rec2 := httptest.NewRecorder()
	eng.Handler().ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("expected 200 for allowed IP, got %d", rec2.Code)
	}
}

func TestRateLimit(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	eng := newTestEngine(t, backend.URL, false)
	eng.EnableRateLimit("", 2, time.Minute)

	// First two requests pass, third is 429.
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "203.0.113.5:1234"
		rec := httptest.NewRecorder()
		eng.Handler().ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i+1, rec.Code)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "203.0.113.5:1234"
	rec := httptest.NewRecorder()
	eng.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 on 3rd request, got %d", rec.Code)
	}
	if rec.Header().Get("Retry-After") == "" {
		t.Fatal("expected Retry-After header on 429")
	}
}

func TestAllowListBypassesBlock(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	eng := newTestEngine(t, backend.URL, false)
	// Allow list and block list are separate CIDRs; the allowed range also
	// appears in the block list to prove allow wins.
	if err := eng.SetIPLists([]string{"198.51.100.0/24"}, []string{"198.51.100.0/24", "192.0.2.0/24"}); err != nil {
		t.Fatalf("SetIPLists: %v", err)
	}

	// Allowed IP must pass even though it is also in the block list.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "198.51.100.10:1234"
	rec := httptest.NewRecorder()
	eng.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for allowed IP despite blocklist, got %d", rec.Code)
	}

	// A non-allowed, blocked IP is still blocked.
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	req2.RemoteAddr = "192.0.2.10:1234"
	rec2 := httptest.NewRecorder()
	eng.Handler().ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for non-allowed blocked IP, got %d", rec2.Code)
	}
}

func TestBlockPageContent(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	eng := newTestEngine(t, backend.URL, false)
	if err := eng.SetIPLists(nil, []string{"192.0.2.0/24"}); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.0.2.50:1234"
	rec := httptest.NewRecorder()
	eng.Handler().ServeHTTP(rec, req)

	body := rec.Body.String()
	if !strings.Contains(body, "Blocked") {
		t.Fatalf("expected block page title in body, got: %s", body)
	}
	if rec.Header().Get("X-Shield-Blocked") != "1" {
		t.Fatal("expected X-Shield-Blocked header")
	}
	if rec.Header().Get("X-Shield-Event-ID") == "" {
		t.Fatal("expected X-Shield-Event-ID header")
	}
}
