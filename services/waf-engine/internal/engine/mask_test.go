package engine

import "testing"

func TestMaskSensitiveValue(t *testing.T) {
	cases := []struct {
		key, value, want string
	}{
		{"q", "normal", "normal"},
		{"password", "supersecret", "***MASKED***"},
		{"PASSWORD", "x", "***MASKED***"},
		{"confirm_password", "x", "***MASKED***"},
		{"token", "abc123", "***MASKED***"},
		{"api_key", "k-123", "***MASKED***"},
		{"session_id", "sid-1", "***MASKED***"},
		{"page", "2", "2"},
		// Sensitive values are masked even when the key name is innocuous (P3.32).
		{"redirect", "https://x/login?password=admin123", "***MASKED***"},
		{"next", "jwt=eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjMifQ.sig", "***MASKED***"},
		{"user", "foo@example.com", "foo@example.com"},
	}
	for _, c := range cases {
		if got := MaskSensitiveValue(c.key, c.value); got != c.want {
			t.Errorf("MaskSensitiveValue(%q, %q) = %q, want %q", c.key, c.value, got, c.want)
		}
	}
}

func TestMaskHeaderValue(t *testing.T) {
	cases := []struct {
		key, value, want string
	}{
		{"Authorization", "Bearer xyz", "***MASKED***"},
		{"authorization", "Basic dXNlcjpwYXNz", "***MASKED***"},
		{"Cookie", "session=abc; user=bob", "***MASKED***"},
		{"Set-Cookie", "sid=123", "***MASKED***"},
		{"X-Api-Key", "k-123", "***MASKED***"},
		{"X-Auth-Token", "t-9", "***MASKED***"},
		{"X-Forwarded-For", "10.0.0.1", "10.0.0.1"},
		{"User-Agent", "curl/8.0", "curl/8.0"},
		{"X-Trace-Id", "abc123", "abc123"},
	}
	for _, c := range cases {
		if got := MaskHeaderValue(c.key, c.value); got != c.want {
			t.Errorf("MaskHeaderValue(%q, %q) = %q, want %q", c.key, c.value, got, c.want)
		}
	}
}

func TestMaskMatchData(t *testing.T) {
	details := []map[string]string{
		{"data": "ARGS:q=1 union select", "message": "SQLi"},
		{"data": "ARGS:password=admin123", "message": "x"},
		{"data": "HEADER:Authorization=Bearer tok", "message": "y"},
		{"data": "ARGS:redirect=/login?password=x", "message": "z"},
		{"data": "HEADER:Content-Type=application/json", "message": "w"},
	}
	out := MaskMatchData(details)

	if out[0]["data"] != "ARGS:q=1 union select" {
		t.Errorf("non-sensitive arg should be unchanged, got %q", out[0]["data"])
	}
	if out[1]["data"] != "ARGS:password=***MASKED***" {
		t.Errorf("password arg must be masked, got %q", out[1]["data"])
	}
	if out[2]["data"] != "HEADER:Authorization=***MASKED***" {
		t.Errorf("authorization header must be masked, got %q", out[2]["data"])
	}
	if out[3]["data"] != "ARGS:redirect=***MASKED***" {
		t.Errorf("arg value with embedded password must be masked, got %q", out[3]["data"])
	}
	if out[4]["data"] != "HEADER:Content-Type=application/json" {
		t.Errorf("non-sensitive header must be unchanged, got %q", out[4]["data"])
	}
}

func TestMaskMatchDataNoEquals(t *testing.T) {
	// A line without "=" should be returned unchanged.
	details := []map[string]string{{"data": "ARGS:justavalue"}}
	out := MaskMatchData(details)
	if out[0]["data"] != "ARGS:justavalue" {
		t.Errorf("expected unchanged, got %q", out[0]["data"])
	}
}
