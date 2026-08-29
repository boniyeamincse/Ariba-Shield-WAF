package engine

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestDetectOnlyForwardsToBackend verifies P0.2: in detect-only mode a matched
// request must be forwarded to the backend (transparent), not dropped with an
// empty 200.
func TestDetectOnlyForwardsToBackend(t *testing.T) {
	backendHit := false
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		backendHit = true
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	var buf strings.Builder
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

	// SQLi attack in detect-only: must reach the backend AND log an event.
	req := httptest.NewRequest(http.MethodGet, "/?q=1%27%20UNION%20SELECT%20*%20FROM%20users--", nil)
	rec := httptest.NewRecorder()
	eng.Handler().ServeHTTP(rec, req)

	if !backendHit {
		t.Fatal("P0.2: detect-only did not forward matched request to backend")
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if buf.Len() == 0 {
		t.Fatal("expected security event logged")
	}
	if !strings.Contains(buf.String(), "950001") {
		t.Fatalf("expected SQLi rule in event, got: %s", buf.String())
	}
}

// TestBlockingDoesNotForwardToBackend verifies blocking mode rejects without
// forwarding.
func TestBlockingDoesNotForwardToBackend(t *testing.T) {
	backendHit := false
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		backendHit = true
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	var buf strings.Builder
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

	req := httptest.NewRequest(http.MethodGet, "/?q=%3Cscript%3Ealert(1)%3C%2Fscript%3E", nil)
	rec := httptest.NewRecorder()
	eng.Handler().ServeHTTP(rec, req)

	if backendHit {
		t.Fatal("blocking mode must not forward matched request to backend")
	}
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}
	if !strings.Contains(buf.String(), "950002") {
		t.Fatalf("expected XSS rule in event, got: %s", buf.String())
	}
}

// TestTrustProxyXFF verifies X-Forwarded-For resolution when trustProxy is on.
func TestTrustProxyXFF(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	eng, err := New(Config{
		BackendURL:       backend.URL,
		RulesPath:        "../../../../rules/core/baseline.conf",
		DetectOnly:       true,
		BlockStatus:      403,
		RequestIDHdr:     "X-Request-ID",
		TrustProxyHeader: true,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	// Block the XFF client; with trustProxy on, the request must be blocked
	// even though RemoteAddr is a different (proxy) address.
	if err := eng.SetIPLists(nil, []string{"203.0.113.0/24"}); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.5:1234" // gateway address
	req.Header.Set("X-Forwarded-For", "203.0.113.9")
	rec := httptest.NewRecorder()
	eng.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for XFF client in blocked CIDR, got %d", rec.Code)
	}
}

// TestIPv6ClientParsing verifies IPv6 addresses parse correctly for reputation.
func TestIPv6ClientParsing(t *testing.T) {
	got := clientIP("[2001:db8::1]:8080")
	if got != "2001:db8::1" {
		t.Fatalf("expected 2001:db8::1, got %q", got)
	}
	if clientIP("192.0.2.1:443") != "192.0.2.1" {
		t.Fatal("IPv4 parsing failed")
	}
}