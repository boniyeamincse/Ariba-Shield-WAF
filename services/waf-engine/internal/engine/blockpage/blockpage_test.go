package blockpage

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRenderEscapesHTML(t *testing.T) {
	// Title/message must be HTML-escaped to prevent stored XSS on the block page.
	body := Render(`<script>alert(1)</script>`, `x < y & "z"`, "evt-1")
	if strings.Contains(body, "<script>") {
		t.Fatal("block page must escape injected HTML in title")
	}
	if !strings.Contains(body, "evt-1") {
		t.Fatal("event ID should appear in the block page")
	}
	if !strings.Contains(body, "&lt;script&gt;") {
		t.Fatal("expected escaped script tag")
	}
}

func TestWriteHeaders(t *testing.T) {
	rec := httptest.NewRecorder()
	Write(rec, http.StatusForbidden, "Blocked", "Security policy", "evt-abc")

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}
	if rec.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Fatalf("expected html content type, got %q", rec.Header().Get("Content-Type"))
	}
	if rec.Header().Get("X-Shield-Blocked") != "1" {
		t.Fatal("expected X-Shield-Blocked header")
	}
}

func TestRenderDefaults(t *testing.T) {
	body := Render("", "", "")
	if !strings.Contains(body, "Request Blocked") {
		t.Fatal("default title should be used when empty")
	}
}