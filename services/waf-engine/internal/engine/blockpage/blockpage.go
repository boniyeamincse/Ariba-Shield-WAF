package blockpage

import (
	"html"
	"net/http"
	"strings"
)

// Render produces the HTML block page body. User-facing text is configurable;
// an optional support/event ID is embedded (master plan §6.8 unique event ID).
func Render(title, message, eventID string) string {
	if title == "" {
		title = "Request Blocked"
	}
	if message == "" {
		message = "This request was blocked by the Ariba Shield WAF security policy."
	}
	var b strings.Builder
	b.WriteString("<!doctype html><html lang=\"en\"><head><meta charset=\"utf-8\">")
	b.WriteString("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">")
	b.WriteString("<title>")
	b.WriteString(html.EscapeString(title))
	b.WriteString("</title><style>body{font-family:system-ui,sans-serif;display:flex;align-items:center;justify-content:center;height:100vh;margin:0;background:#f5f5f5;color:#1a1a1a}main{text-align:center;max-width:480px;padding:2rem}code{background:#eee;padding:.25rem .5rem;border-radius:4px}</style></head><body><main>")
	b.WriteString("<h1>")
	b.WriteString(html.EscapeString(title))
	b.WriteString("</h1><p>")
	b.WriteString(html.EscapeString(message))
	b.WriteString("</p>")
	if eventID != "" {
		b.WriteString("<p>Event ID: <code>")
		b.WriteString(html.EscapeString(eventID))
		b.WriteString("</code></p>")
	}
	b.WriteString("</main></body></html>")
	return b.String()
}

// Write serves the block page with the configured status code and content type.
func Write(w http.ResponseWriter, status int, title, message, eventID string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Shield-Blocked", "1")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(Render(title, message, eventID)))
}
