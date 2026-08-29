package engine

import "strings"

// sensitiveKeys are parameter names whose values are never stored in events
// (master plan §22 rule 5: never store unmasked secrets for analyst convenience).
var sensitiveKeys = []string{
	"password", "passwd", "pwd", "secret", "token", "accesstoken",
	"authorization", "apikey", "api_key", "client_secret", "session",
	"sessionid", "cookie", "set-cookie", "x-api-key", "x-auth-token",
}

// MaskSensitiveValue replaces a parameter value with a placeholder when the
// key matches a sensitive name. Used on ARGS/headers before storing events.
func MaskSensitiveValue(key, value string) string {
	lower := strings.ToLower(key)
	for _, s := range sensitiveKeys {
		if strings.Contains(lower, s) {
			return "***MASKED***"
		}
	}
	return value
}

// MaskMatchData applies MaskSensitiveValue to the "data" field of match details
// which are formatted as "VARIABLE:key=value".
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
	// Strip the leading "ARGS:", "HEADER:", etc. prefix to match the key name.
	if idx := strings.LastIndex(key, ":"); idx != -1 {
		key = key[idx+1:]
	}
	value := MaskSensitiveValue(key, parts[1])
	return parts[0] + "=" + value
}