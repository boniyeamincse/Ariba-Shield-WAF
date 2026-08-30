package engine

import "strings"

// MaskedValue is the placeholder used in place of sensitive data.
const MaskedValue = "***MASKED***"

// sensitiveKeys are ARGS-style parameter names whose values are never stored
// in events (master plan §22 rule 5: never store unmasked secrets).
var sensitiveKeys = []string{
	"password", "passwd", "pwd", "secret", "token", "accesstoken",
	"apikey", "api_key", "client_secret", "session", "sessionid",
	"credential", "credentials", "auth", "passphrase",
}

// headerKeys are request-header names (and similar locations such as cookies)
// whose values are never stored. Authorization, cookies and API-key headers
// belong to this list, not the ARGS key list (P3.32).
var headerKeys = []string{
	"authorization", "cookie", "set-cookie", "x-api-key", "x-auth-token",
	"proxy-authorization", "x-csrf-token", "x-xsrf-token",
}

// sensitiveValuePatterns are substrings that indicate a *value* carries
// sensitive material even when the key name is innocuous (P3.32).
var sensitiveValuePatterns = []string{
	"bearer ", "basic ", "digest ",
	"password=", "passwd=", "pwd=", "secret=", "token=",
	"apikey=", "api_key=", "client_secret=", "authorization=",
	"credential", "private_key", "-----begin", "x-api-key", "x-auth-token",
}

// MaskSensitiveValue replaces a parameter value with a placeholder when the
// key matches a sensitive ARGS name or the value itself looks sensitive.
func MaskSensitiveValue(key, value string) string {
	lower := strings.ToLower(key)
	for _, s := range sensitiveKeys {
		if strings.Contains(lower, s) {
			return MaskedValue
		}
	}
	return maskSensitiveValueContent(value)
}

// MaskHeaderValue masks request-header values. Authorization, cookies and
// API-key style headers are treated as header keys (P3.32).
func MaskHeaderValue(key, value string) string {
	lower := strings.ToLower(key)
	for _, s := range headerKeys {
		if strings.Contains(lower, s) {
			return MaskedValue
		}
	}
	// A header can still be one of the generic sensitive key names.
	return MaskSensitiveValue(key, value)
}

// maskSensitiveValueContent masks a value when it carries sensitive material
// regardless of the associated key: bearer/basic auth payloads, embedded
// key=value credentials, or JWT-shaped tokens.
func maskSensitiveValueContent(value string) string {
	if value == "" {
		return value
	}
	lower := strings.ToLower(value)
	for _, s := range sensitiveValuePatterns {
		if strings.Contains(lower, s) {
			return MaskedValue
		}
	}
	// JWT-shaped token: header.payload.signature, header starts with eyJ.
	if looksLikeJWT(value) {
		return MaskedValue
	}
	return value
}

// looksLikeJWT reports whether value contains a three-dot JWT with a
// base64url header (eyJ...), e.g. "Bearer eyJ..." or "jwt=eyJ...".
func looksLikeJWT(value string) bool {
	if strings.Count(value, ".") < 2 {
		return false
	}
	for _, part := range strings.FieldsFunc(value, func(r rune) bool {
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

// MaskMatchData applies masking to the "data" field of match details which are
// formatted as "VARIABLE:key=value".
func MaskMatchData(details []map[string]string) []map[string]string {
	out := make([]map[string]string, 0, len(details))
	for _, d := range details {
		copy := make(map[string]string, len(d))
		for k, v := range d {
			copy[k] = v
		}
		if data, ok := copy["data"]; ok {
			copy["data"] = maskDataLine(data)
		}
		out = append(out, copy)
	}
	return out
}

func maskDataLine(line string) string {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return line
	}
	key := parts[0]
	// Strip the leading "ARGS:", "HEADER:", "COOKIE:", etc. prefix to match the
	// key name and decide whether this is a header/cookie context.
	isHeader := false
	if idx := strings.LastIndex(key, ":"); idx != -1 {
		prefix := key[:idx]
		isHeader = strings.Contains(prefix, "HEADER") || strings.Contains(prefix, "COOKIE")
		key = key[idx+1:]
	}
	var value string
	if isHeader {
		value = MaskHeaderValue(key, parts[1])
	} else {
		value = MaskSensitiveValue(key, parts[1])
	}
	return parts[0] + "=" + value
}
