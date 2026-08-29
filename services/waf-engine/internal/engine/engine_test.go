package engine

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEngineDetectOnly(t *testing.T) {
	// Start a backend that returns 200.
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	var buf bytes.Buffer
	eng, err := New(Config{
		BackendURL:   backend.URL,
		RulesPath:    "../../../../rules/core/baseline.conf",
		DetectOnly:   true,
		EventSink:    &buf,
		BlockStatus:  403,
		RequestIDHdr: "X-Request-ID",
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	// Test 1: normal request passes through.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-ID", "test-01")
	rec := httptest.NewRecorder()
	eng.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for normal request, got %d", rec.Code)
	}
	// No security events logged for normal traffic.
	if buf.Len() > 0 {
		t.Logf("events logged: %s", buf.String())
	}

	// Test 2: SQL injection in args triggers a security event but passes (detect-only).
	buf.Reset()
	req2 := httptest.NewRequest(http.MethodGet, "/?q=1%27%20UNION%20SELECT%20*%20FROM%20users--", nil)
	req2.Header.Set("X-Request-ID", "test-02")
	rec2 := httptest.NewRecorder()
	eng.Handler().ServeHTTP(rec2, req2)
	// Detect-only: request should STILL get 200.
	if rec2.Code != http.StatusOK {
		t.Fatalf("expected 200 in detect-only for attack, got %d", rec2.Code)
	}
	// But a security event must have been logged.
	if buf.Len() == 0 {
		t.Fatal("expected security event logged for SQLi attack, got none")
	}
	body := buf.String()
	if !strings.Contains(body, "SQL Injection") {
		t.Fatalf("expected SQL Injection in event, got: %s", body)
	}
	if !strings.Contains(body, "attack-sqli") {
		t.Fatalf("expected attack-sqli tag in event, got: %s", body)
	}
}

func TestEngineBlocking(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	var buf bytes.Buffer
	eng, err := New(Config{
		BackendURL:   backend.URL,
		RulesPath:    "../../../../rules/core/baseline.conf",
		DetectOnly:   false,
		EventSink:    &buf,
		BlockStatus:  403,
		RequestIDHdr: "X-Request-ID",
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	// SQLi should be blocked.
	req := httptest.NewRequest(http.MethodGet, "/?q=1%27%20UNION%20SELECT%20*%20FROM%20users--", nil)
	req.Header.Set("X-Request-ID", "test-03")
	rec := httptest.NewRecorder()
	eng.Handler().ServeHTTP(rec, req)
	if rec.Code != 403 {
		t.Fatalf("expected 403 for XSS in blocking mode, got %d", rec.Code)
	}
	if buf.Len() == 0 {
		t.Fatal("expected security event log for blocked request")
	}
}

func TestEngineBodyProcessing(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	var buf bytes.Buffer
	eng, err := New(Config{
		BackendURL:   backend.URL,
		RulesPath:    "../../../../rules/core/baseline.conf",
		DetectOnly:   true,
		EventSink:    &buf,
		BlockStatus:  403,
		RequestIDHdr: "X-Request-ID",
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	// POST with SQLi in body.
	body := strings.NewReader("user=admin&q=1%27%20UNION%20SELECT%20*%20FROM%20users--")
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Request-ID", "test-04")
	rec := httptest.NewRecorder()
	eng.Handler().ServeHTTP(rec, req)
	if buf.Len() == 0 {
		t.Fatal("expected security event for POST body SQLi")
	}
	if !strings.Contains(buf.String(), "SQL Injection") {
		t.Fatalf("expected SQL Injection in event, got: %s", buf.String())
	}
}