package engine

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ariba-shield/waf-engine/internal/engine/blockpage"
	"github.com/ariba-shield/waf-engine/internal/engine/iplist"
	"github.com/ariba-shield/waf-engine/internal/engine/ratelimit"
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
	// AnomalyThreshold: if > 0, block when the transaction anomaly score
	// reaches this value even if no single rule interruption fired (Phase 3).
	AnomalyThreshold int
	// BlockTitle/BlockMessage/BlockEventID configure the custom block page.
	BlockTitle   string
	BlockMessage string
	BlockEventID string
	// TrustProxyHeader, when true, trusts X-Forwarded-For / X-Real-IP set by
	// the trusted gateway in front of the engine (ADR-004 mTLS chain). Only
	// enable when the engine is NOT directly reachable by clients.
	TrustProxyHeader bool
}

// Engine wraps Coraza and the reverse proxy.
type Engine struct {
	waf            coraza.WAF
	backend        *url.URL
	detectOnly     bool
	sink           io.Writer
	blockStatus    int
	reqIDHdr       string
	anomalyThresh  int
	blockTitle     string
	blockMessage   string
	blockEventID   string
	ipList         *iplist.List
	rateLimiter    *ratelimit.SlidingWindow
	rateLimitRoute string // path prefix to rate limit, "" = all
	rateLimitCount int
	rateLimitWin   time.Duration
	trustProxy     bool
}

// New builds the engine from config.
func New(cfg Config) (*Engine, error) {
	wafCfg := coraza.NewWAFConfig().WithDirectivesFromFile(cfg.RulesPath)

	// Phase 3: anomaly threshold override. The baseline.conf already has the
	// blocking rule (949110) with a default threshold. If the user requests a
	// different threshold, we override tx.blocking_anomaly_score.
	if cfg.AnomalyThreshold > 0 {
		wafCfg = wafCfg.WithDirectives(
			fmt.Sprintf("SecAction \"id:949101,phase:1,pass,nolog,setvar:tx.blocking_anomaly_score=%d\"", cfg.AnomalyThreshold),
		)
	}

	waf, err := coraza.NewWAF(wafCfg)
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
		waf:            waf,
		backend:        backend,
		detectOnly:     cfg.DetectOnly,
		sink:           sink,
		blockStatus:    blockStatus,
		reqIDHdr:       cfg.RequestIDHdr,
		anomalyThresh:  cfg.AnomalyThreshold,
		blockTitle:     cfg.BlockTitle,
		blockMessage:   cfg.BlockMessage,
		blockEventID:   cfg.BlockEventID,
		ipList:         iplist.New(),
		trustProxy:     cfg.TrustProxyHeader,
	}, nil
}

// clientAddress resolves the client IP for reputation/rate-limit decisions.
// If trustProxy is enabled, it honors X-Forwarded-For (leftmost entry set by
// the trusted gateway) or X-Real-IP. Otherwise it uses the connection peer.
func (e *Engine) clientAddress(r *http.Request) string {
	if e.trustProxy {
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			// Leftmost is the original client (gateway prepends).
			if i := strings.Index(xff, ","); i != -1 {
				xff = xff[:i]
			}
			return clientIP(xff)
		}
		if xri := r.Header.Get("X-Real-IP"); xri != "" {
			return clientIP(xri)
		}
	}
	return clientIP(r.RemoteAddr)
}

// SetIPLists replaces the allow/block IP prefix sets atomically.
func (e *Engine) SetIPLists(allowed, blocked []string) error {
	if err := e.ipList.SetAllowed(allowed); err != nil {
		return err
	}
	return e.ipList.SetBlocked(blocked)
}

// SetRateLimit configures per-IP rate limiting. limit<=0 disables.
func (e *Engine) SetRateLimit(route string, count int, window time.Duration) {
	if count <= 0 {
		e.rateLimiter = nil
		return
	}
	e.rateLimiter = ratelimit.New(count, window)
	e.rateLimitRoute = route
	e.rateLimitCount = count
	e.rateLimitWin = window
}

// EnableRateLimit creates a limiter when none exists (for tests).
func (e *Engine) EnableRateLimit(route string, count int, window time.Duration) {
	if e.rateLimiter == nil {
		e.rateLimiter = ratelimit.New(count, window)
	}
	e.rateLimitRoute = route
	e.rateLimitCount = count
	e.rateLimitWin = window
}

// Handler returns the HTTP handler that inspects and forwards.
func (e *Engine) Handler() http.Handler {
	proxy := httputil.NewSingleHostReverseProxy(e.backend)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ip := e.clientAddress(r)

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
			ClientIP:       ip,
			DecisionAction: "pass",
		}

		// Phase 3: IP reputation pre-check (step 7 of the pipeline).
		if e.ipList.IsBlocked(ip) {
			event.DecisionAction = "block"
			event.Severity = "medium"
			event.Reason = "ip_blocklist"
			event.RuleIDs = []string{"ip-list"}
			event.Status = e.blockStatus
			e.logEvent(event)
			e.serveBlock(w, event)
			return
		}

		// Phase 3: rate limit (step 14 of the pipeline).
		if e.rateLimiter != nil {
			routeOK := e.rateLimitRoute == "" || strings.HasPrefix(r.URL.Path, e.rateLimitRoute)
			if routeOK {
				allowed, remaining := e.rateLimiter.Allow(ip)
				if !allowed {
					event.DecisionAction = "rate_limit"
					event.Severity = "medium"
					event.Reason = "rate_limit_exceeded"
					event.RuleIDs = []string{"rate-limit"}
					event.Status = 429
					e.logEvent(event)
					w.Header().Set("Retry-After", fmt.Sprintf("%d", int(e.rateLimitWin.Seconds())))
					w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
					blockpage.Write(w, http.StatusTooManyRequests,
						e.blockTitle, "Too many requests. Please try again later.", event.EventID)
					return
				}
				w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
			}
		}

		// Build a Coraza transaction.
		tx := e.waf.NewTransaction()
		defer tx.Close()

		// Seed the transaction from the request.
		tx.ProcessConnection(event.ClientIP, 0, r.Host, 0)
		tx.ProcessURI(r.URL.RequestURI(), r.Method, r.Proto)
		for k, vs := range r.Header {
			for _, v := range vs {
				tx.AddRequestHeader(k, v)
			}
		}

		// Buffer the body so Coraza can inspect it AND the proxy can forward it.
		// Cap the body size to SecRequestBodyLimit (13 MB in baseline.conf) to
		// prevent unbounded memory growth (P0.7).
		const maxBody = 13 << 20 // 13 MB
		var bodyBytes []byte
		if tx.IsRequestBodyAccessible() && r.Body != nil {
			limited := io.LimitReader(r.Body, maxBody+1)
			var err error
			bodyBytes, err = io.ReadAll(limited)
			if err != nil {
				event.Reason = "engine_error"
				event.DecisionAction = "pass"
				e.logEvent(event)
				proxy.ServeHTTP(w, r)
				return
			}
			if len(bodyBytes) > maxBody {
				// Body exceeds limit — reject with 413 (fail open in detect-only).
				event.Reason = "body_too_large"
				event.DecisionAction = "block"
				event.Severity = "medium"
				event.RuleIDs = []string{"body-limit"}
				event.Status = 413
				e.logEvent(event)
				if !e.detectOnly {
					e.serveBlock(w, event)
					return
				}
				// In detect-only: log and continue.
				r.Body = io.NopCloser(bytes.NewReader(bodyBytes[:maxBody]))
			} else {
				r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			}
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
					if e.handleInterruption(w, event, tx, it) {
						proxy.ServeHTTP(w, r)
					}
					return
				}
			}
		}

		it := tx.ProcessRequestHeaders()
		if it != nil {
			if e.handleInterruption(w, event, tx, it) {
				proxy.ServeHTTP(w, r)
			}
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
				if e.handleInterruption(w, event, tx, it) {
					proxy.ServeHTTP(w, r)
				}
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

// handleInterruption logs the matched rules and either forwards to the backend
// (detect-only) or serves a block page. Returns true if the caller should
// continue to proxy.ServeHTTP (detect-only passthrough), false if blocked.
func (e *Engine) handleInterruption(w http.ResponseWriter, event *SecurityEvent, tx types.Transaction, it *types.Interruption) bool {
	matched := tx.MatchedRules()
	event.DecisionAction = "log"
	event.Severity = "high"
	event.RuleIDs = matchedIDs(matched)
	event.MatchDetails = MaskMatchData(matchedDetails(matched))
	event.Status = it.Status
	event.Reason = fmt.Sprintf("rule_%d_%s", it.RuleID, it.Action)
	e.logEvent(event)

	if e.detectOnly {
		// Transparent: log the event, then forward to the backend.
		// The caller must call proxy.ServeHTTP after this returns true.
		return true
	}
	// Blocking: serve the block page.
	e.serveBlock(w, event)
	return false
}

// serveBlock writes the custom block page or a plain status line.
func (e *Engine) serveBlock(w http.ResponseWriter, event *SecurityEvent) {
	w.Header().Set("X-Shield-Blocked", "1")
	w.Header().Set("X-Shield-Event-ID", event.EventID)
	if e.blockTitle != "" {
		blockpage.Write(w, e.blockStatus, e.blockTitle, e.blockMessage, event.EventID)
	} else {
		w.WriteHeader(e.blockStatus)
	}
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
	// Strip port from IPv4 (1.2.3.4:5678) or IPv6 ([::1]:8080).
	// net.SplitHostPort handles both correctly.
	if host, _, err := net.SplitHostPort(remote); err == nil {
		remote = host
	}
	// Strip brackets from IPv6 (should already be clean after SplitHostPort).
	remote = strings.Trim(remote, "[]")
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
