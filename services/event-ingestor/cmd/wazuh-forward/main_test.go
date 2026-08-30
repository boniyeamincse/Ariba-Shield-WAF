package main

import (
	"strings"
	"testing"
)

func TestRedactEvent(t *testing.T) {
	ev := map[string]any{
		"reason": "sqli",
		"path":   "/login",
		"cookie": "session=abc",
		"raw": map[string]any{
			"authorization": "Bearer topsecret",
			"password":      "hunter2",
			"headers":       []any{"x", "y"},
			"nested":        map[string]any{"token": "t-1", "keep": "val"},
		},
		"match_details": []any{map[string]any{"data": "ARGS:password=admin123"}},
	}
	got := redactEvent(ev).(map[string]any)

	if got["cookie"] != redactedPlaceholder {
		t.Errorf("cookie must be masked, got %v", got["cookie"])
	}
	raw := got["raw"].(map[string]any)
	if raw["authorization"] != redactedPlaceholder || raw["password"] != redactedPlaceholder {
		t.Errorf("raw credentials must be masked, got %v", raw)
	}
	nested := raw["nested"].(map[string]any)
	if nested["token"] != redactedPlaceholder {
		t.Errorf("nested token must be masked, got %v", nested["token"])
	}
	if nested["keep"] != "val" {
		t.Errorf("benign nested value must be unchanged, got %v", nested["keep"])
	}
	if got["path"] != "/login" || got["reason"] != "sqli" {
		t.Errorf("benign top-level fields must be unchanged, got %v", got)
	}
	md := got["match_details"].([]any)[0].(map[string]any)
	if md["data"] != redactedPlaceholder {
		t.Errorf("match detail with password must be masked, got %v", md["data"])
	}
}

func TestRedactStringValue(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Bearer abc", redactedPlaceholder},
		{"basic dXNlcjpwYXNz", redactedPlaceholder},
		{"/login?password=hunter2", redactedPlaceholder},
		{"jwt=eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjMifQ.sig", redactedPlaceholder},
		{"normal-value", "normal-value"},
		{"GET / 200", "GET / 200"},
	}
	for _, c := range cases {
		if got := redactStringValue(c.in); got != c.want {
			t.Errorf("redactStringValue(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestFormatSyslog(t *testing.T) {
	msg := []byte("hello")
	for _, rfc := range []int{3164, 5424} {
		line := formatSyslog(rfc, 13, "gw-01", "950001", msg)
		if !strings.HasPrefix(line, "<13>") {
			t.Errorf("rfc%d: expected PRI 13 prefix, got %q", rfc, line)
		}
		if !strings.Contains(line, "gw-01") || !strings.Contains(line, "hello") {
			t.Errorf("rfc%d: expected hostname and message in %q", rfc, line)
		}
	}
	if !strings.Contains(formatSyslog(5424, 13, "gw-01", "950001", msg), "1 ") {
		t.Error("rfc5424 must include the version '1'")
	}
}

func TestSyslogPRI(t *testing.T) {
	cases := []struct {
		level, want int
	}{
		{12, 10}, // LOG_CRIT
		{8, 11},  // LOG_ERR
		{5, 13},  // LOG_NOTICE
		{3, 14},  // LOG_INFO
	}
	for _, c := range cases {
		if got := syslogPRI(c.level); got != c.want {
			t.Errorf("syslogPRI(%d) = %d, want %d", c.level, got, c.want)
		}
	}
}
