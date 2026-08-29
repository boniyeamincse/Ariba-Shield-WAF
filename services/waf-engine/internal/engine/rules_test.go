package engine

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewRuleMatching(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
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

	cases := []struct {
		name string
		url  string
		want string
	}{
		{"cmdi-whoami", "/?cmd=%3B%20whoami", "950003"},
		{"cmdi-pwd", "/?y=%3B%20pwd", "950003"},
		{"sqli-or", "/?search=foo%27%20OR%201=1%20--", "950001"},
		{"lfi-file", "/?page=file%3A%2F%2F%2Fetc%2Fpasswd", "950005"},
		{"xss-script", "/?q=%3Cscript%3Ealert(1)%3C%2Fscript%3E", "950002"},
	}

	for _, c := range cases {
		buf.Reset()
		req := httptest.NewRequest(http.MethodGet, c.url, nil)
		rec := httptest.NewRecorder()
		eng.Handler().ServeHTTP(rec, req)
		if !strings.Contains(buf.String(), c.want) {
			t.Errorf("%s: expected rule %s, got events: %s", c.name, c.want, buf.String())
		} else {
			t.Logf("PASS %s: matched %s", c.name, c.want)
		}
	}
}
