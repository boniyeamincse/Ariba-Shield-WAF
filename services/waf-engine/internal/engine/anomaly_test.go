package engine

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAnomalyThreshold(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	// Threshold 1: any single rule match (SQLi) sets anomaly_score=+1 => blocked.
	eng, err := New(Config{
		BackendURL:       backend.URL,
		RulesPath:        "../../../../rules/core/baseline.conf",
		DetectOnly:       false,
		BlockStatus:      403,
		RequestIDHdr:     "X-Request-ID",
		AnomalyThreshold: 1,
		BlockTitle:       "Blocked",
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	// Attack request must be blocked by the anomaly gate.
	req := httptest.NewRequest(http.MethodGet, "/?q=1%27%20UNION%20SELECT%20*%20FROM%20users--", nil)
	rec := httptest.NewRecorder()
	eng.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 from anomaly threshold, got %d", rec.Code)
	}

	// Normal request must pass (score stays 0).
	req2 := httptest.NewRequest(http.MethodGet, "/?q=hello", nil)
	rec2 := httptest.NewRecorder()
	eng.Handler().ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("expected 200 for clean request, got %d", rec2.Code)
	}
}

func TestAnomalyThresholdDetectOnly(t *testing.T) {
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
		AnomalyThreshold: 1,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	// Detect-only: anomaly gate must NOT block, traffic passes.
	req := httptest.NewRequest(http.MethodGet, "/?q=1%27%20UNION%20SELECT%20*%20FROM%20users--", nil)
	rec := httptest.NewRecorder()
	eng.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 in detect-only despite threshold, got %d", rec.Code)
	}
}

func TestAnomalyThresholdEventTagged(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	var buf strings.Builder
	eng, err := New(Config{
		BackendURL:       backend.URL,
		RulesPath:        "../../../../rules/core/baseline.conf",
		DetectOnly:       false,
		EventSink:        &buf,
		BlockStatus:      403,
		RequestIDHdr:     "X-Request-ID",
		AnomalyThreshold: 1,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/?q=%3Cscript%3Ealert(1)%3C%2Fscript%3E", nil)
	rec := httptest.NewRecorder()
	eng.Handler().ServeHTTP(rec, req)
	if !strings.Contains(buf.String(), "anomaly-threshold") {
		t.Fatalf("expected anomaly-threshold event, got: %s", buf.String())
	}
}