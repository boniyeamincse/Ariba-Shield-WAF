package engine

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/corazawaf/coraza/v3"
	"github.com/corazawaf/coraza/v3/types"
)

// Config holds runtime settings for the WAF sidecar.
type Config struct {
	ListenAddr   string
	BackendURL   string
	RulesPath    string // path to the CRS/rule file
	DetectOnly   bool   // transparent mode: log events, never block
	EventSink    io.Writer
	BlockStatus  int
	RequestIDHdr string
}

// Engine wraps Coraza and the reverse proxy.
type Engine struct {
	waf         coraza.WAF
	backend     *url.URL
	detectOnly  bool
	sink        io.Writer
	blockStatus int
	reqIDHdr    string
}

// New builds the engine from config.
func New(cfg Config) (*Engine, error) {
	waf, err := coraza.NewWAF(
		coraza.NewWAFConfig().WithDirectivesFromFile(cfg.RulesPath),
	)
	if err != nil {
		return nil, fmt.Errorf("load rules: %w", err)
	}

	backend, err := url.Parse(cfg.BackendURL)
	if err != nil {
		return nil, fmt.Errorf("parse backend url: %w", err)
	}

	sink := cfg.EventSink
	if sink == nil {
		sink = os.Stdout
	}
	blockStatus := cfg.BlockStatus
	if blockStatus == 0 {
		blockStatus = 403
	}

	return &Engine{
		waf:         waf,
		backend:     backend,
		detectOnly:  cfg.DetectOnly,
		sink:        sink,
		blockStatus: blockStatus,
		reqIDHdr:    cfg.RequestIDHdr,
	}, nil
}

// Handler returns the HTTP handler that inspects and forwards.
func (e *Engine) Handler() http.Handler {
	proxy := httputil.NewSingleHostReverseProxy(e.backend)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Build a Coraza transaction.
		tx := e.waf.NewTransaction()
		defer tx.Close()

		reqID := r.Header.Get(e.reqIDHdr)
		event := &SecurityEvent{
			SchemaVersion:  "0.1",
			EventID:        newID(),
			RequestID:      reqID,
			Timestamp:      time.Now().UTC().Format(time.RFC3339),
			EventType:      "security",
			GatewayID:      os.Getenv("GATEWAY_ID"),
			Method:         r.Method,
			Path:           r.URL.Path,
			Host:           r.Host,
			ClientIP:       clientIP(r.RemoteAddr),
			DecisionAction: "pass",
		}

		// Seed the transaction from the request.
		tx.ProcessConnection(event.ClientIP, 0, r.Host, 0)
		tx.ProcessURI(r.URL.RequestURI(), r.Method, r.Proto)
		for k, vs := range r.Header {
			for _, v := range vs {
				tx.AddRequestHeader(k, v)
			}
		}

		// Buffer the body so Coraza can inspect it AND the proxy can forward it.
		var bodyBytes []byte
		if tx.IsRequestBodyAccessible() && r.Body != nil {
			var err error
			bodyBytes, err = io.ReadAll(r.Body)
			if err != nil {
				// Engine failure: fail open in detect-only, never block traffic.
				event.Reason = "engine_error"
				event.DecisionAction = "pass"
				e.logEvent(event)
				proxy.ServeHTTP(w, r)
				return
			}
			r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			if len(bodyBytes) > 0 {
				it, _, err := tx.ReadRequestBodyFrom(bytes.NewReader(bodyBytes))
				if err != nil {
					event.Reason = "engine_error"
					event.DecisionAction = "pass"
					e.logEvent(event)
					proxy.ServeHTTP(w, r)
					return
				}
				if it != nil {
					e.handleInterruption(w, event, tx, it)
					return
				}
			}
		}

		it := tx.ProcessRequestHeaders()
		if it != nil {
			e.handleInterruption(w, event, tx, it)
			return
		}

		// Phase 2 rules evaluate ARGS (query string + body). Must run even for
		// GET requests to evaluate query-string args.
		if tx.IsRequestBodyAccessible() {
			it, err := tx.ProcessRequestBody()
			if err != nil {
				event.Reason = "engine_error"
				event.DecisionAction = "pass"
				e.logEvent(event)
				proxy.ServeHTTP(w, r)
				return
			}
			if it != nil {
				e.handleInterruption(w, event, tx, it)
				return
			}
		}

		// Forward to backend, preserving the original body.
		proxy.ServeHTTP(w, r)
		event.Status = 200
		event.LatencyMs = time.Since(start).Milliseconds()
		event.DecisionAction = "pass"
	})
}

func (e *Engine) handleInterruption(w http.ResponseWriter, event *SecurityEvent, tx types.Transaction, it *types.Interruption) {
	matched := tx.MatchedRules()
	event.DecisionAction = "log"
	event.Severity = "high"
	event.RuleIDs = matchedIDs(matched)
	event.MatchDetails = MaskMatchData(matchedDetails(matched))
	event.Status = it.Status
	event.Reason = fmt.Sprintf("rule_%d_%s", it.RuleID, it.Action)
	e.logEvent(event)

	if e.detectOnly {
		// Transparent: serve a synthesized 200 so detection never blocks.
		w.WriteHeader(http.StatusOK)
		return
	}
	status := e.blockStatus
	if it.Status != 0 {
		status = it.Status
	}
	w.WriteHeader(status)
}

// logEvent writes a security event as JSON-lines (non-blocking upstream).
func (e *Engine) logEvent(ev *SecurityEvent) {
	line, err := ev.JSON()
	if err != nil {
		return
	}
	if _, err := e.sink.Write(line); err != nil {
		return
	}
	_, _ = fmt.Fprintln(e.sink)
}

func matchedIDs(ms []types.MatchedRule) []string {
	out := make([]string, 0, len(ms))
	for _, m := range ms {
		out = append(out, fmt.Sprintf("%d", m.Rule().ID()))
	}
	return out
}

func matchedDetails(ms []types.MatchedRule) []map[string]string {
	out := make([]map[string]string, 0, len(ms))
	for _, m := range ms {
		datas := m.MatchedDatas()
		var vars []string
		for _, d := range datas {
			vars = append(vars, fmt.Sprintf("%s:%s=%s", d.Variable().Name(), d.Key(), truncate(d.Value(), 128)))
		}
		tags := m.Rule().Tags()
		out = append(out, map[string]string{
			"rule_id": fmt.Sprintf("%d", m.Rule().ID()),
			"message": m.Message(),
			"data":    strings.Join(vars, "; "),
			"tags":    strings.Join(tags, ", "),
		})
	}
	return out
}

func clientIP(remote string) string {
	if i := strings.LastIndex(remote, ":"); i != -1 {
		return remote[:i]
	}
	return remote
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// newID is a small ULID-like generator; in production uses the real ULID lib.
func newID() string {
	return fmt.Sprintf("%d-%x", time.Now().UnixNano(), os.Getpid())
}
