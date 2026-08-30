package main

import "strings"

// redactedPlaceholder replaces sensitive event content before forwarding
// (P2.30: never ship unmasked secrets to Wazuh).
const redactedPlaceholder = "***MASKED***"

// redactKeys are map keys whose values are always masked, matched by
// case-insensitive substring (cookie, authorization, credentials, tokens,
// passwords, API keys).
var redactKeys = []string{
	"password", "passwd", "pwd", "secret", "token", "accesstoken",
	"authorization", "apikey", "api_key", "client_secret", "credential",
	"session", "sessionid", "cookie", "set-cookie",
	"x-api-key", "x-auth-token", "proxy-authorization",
	"x-csrf-token", "x-xsrf-token",
}

// redactValuePatterns flag string values that carry sensitive material even
// under an innocuous key.
var redactValuePatterns = []string{
	"bearer ", "basic ", "digest ",
	"password=", "passwd=", "pwd=", "secret=", "token=",
	"apikey=", "api_key=", "client_secret=", "authorization=",
	"private_key", "-----begin",
}

func isRedactKey(key string) bool {
	lower := strings.ToLower(key)
	for _, s := range redactKeys {
		if strings.Contains(lower, s) {
			return true
		}
	}
	return false
}

func redactStringValue(s string) string {
	lower := strings.ToLower(s)
	for _, p := range redactValuePatterns {
		if strings.Contains(lower, p) {
			return redactedPlaceholder
		}
	}
	if looksLikeJWT(s) {
		return redactedPlaceholder
	}
	return s
}

// redactEvent walks an event map and masks sensitive keys and values in place
// (recursively through nested maps and arrays).
func redactEvent(v any) any {
	switch val := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(val))
		for k, x := range val {
			if isRedactKey(k) {
				out[k] = redactedPlaceholder
				continue
			}
			out[k] = redactEvent(x)
		}
		return out
	case []any:
		out := make([]any, len(val))
		for i, x := range val {
			out[i] = redactEvent(x)
		}
		return out
	case string:
		return redactStringValue(val)
	default:
		return v
	}
}

// looksLikeJWT reports whether s contains a three-dot JWT with a base64url
// header (eyJ...), e.g. "Bearer eyJ..." or "jwt=eyJ...".
func looksLikeJWT(s string) bool {
	if strings.Count(s, ".") < 2 {
		return false
	}
	for _, part := range strings.FieldsFunc(s, func(r rune) bool {
		return r == ' ' || r == '=' || r == '"' || r == '\''
	}) {
		if !strings.HasPrefix(part, "eyJ") && !strings.HasPrefix(strings.ToLower(part), "eyj") {
			continue
		}
		rest := part[strings.Index(part, ".")+1:]
		segs := strings.Split(rest, ".")
		if len(segs) < 2 || segs[0] == "" || segs[1] == "" {
			continue
		}
		return true
	}
	return false
}
