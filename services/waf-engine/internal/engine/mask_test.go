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
		{"authorization", "Bearer xyz", "***MASKED***"},
		{"session_id", "sid-1", "***MASKED***"},
		{"page", "2", "2"},
	}
	for _, c := range cases {
		if got := MaskSensitiveValue(c.key, c.value); got != c.want {
			t.Errorf("MaskSensitiveValue(%q, %q) = %q, want %q", c.key, c.value, got, c.want)
		}
	}
}

func TestMaskMatchData(t *testing.T) {
	details := []map[string]string{
		{"data": "ARGS:q=1 union select", "message": "SQLi"},
		{"data": "ARGS:password=admin123", "message": "x"},
		{"data": "HEADER:Authorization=Bearer tok", "message": "y"},
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
}

func TestMaskMatchDataNoEquals(t *testing.T) {
	// A line without "=" should be returned unchanged.
	details := []map[string]string{{"data": "ARGS:justavalue"}}
	out := MaskMatchData(details)
	if out[0]["data"] != "ARGS:justavalue" {
		t.Errorf("expected unchanged, got %q", out[0]["data"])
	}
}