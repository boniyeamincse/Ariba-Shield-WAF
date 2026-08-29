package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ariba-shield/waf-engine/internal/engine"
)

func main() {
	listen := flag.String("listen", ":8082", "listening address (gateway forwards here)")
	backend := flag.String("backend", "", "upstream backend URL, e.g. http://backend-a:80")
	rules := flag.String("rules", "/rules/crs.conf", "path to the CRS/rule file")
	detectOnly := flag.Bool("detect-only", true, "transparent mode: log events, never block")
	blockStatus := flag.Int("block-status", 403, "HTTP status when blocking (ignored in detect-only)")
	blockTitle := flag.String("block-title", "", "custom block page title (empty = plain status)")
	blockMessage := flag.String("block-message", "", "custom block page body text")
	anomalyThreshold := flag.Int("anomaly-threshold", 0, "anomaly score threshold (0 = disabled)")
	rateLimit := flag.Int("rate-limit", 0, "per-IP rate limit req/min (0 = disabled)")
	rateLimitRoute := flag.String("rate-limit-route", "", "rate limit applies to this path prefix only (empty = all)")
	allowedIPs := flag.String("allowed-ips", "", "comma-separated list of allowed IP/CIDR")
	blockedIPs := flag.String("blocked-ips", "", "comma-separated list of blocked IP/CIDR")
	trustProxy := flag.Bool("trust-proxy", false, "trust X-Forwarded-For/X-Real-IP from the gateway (only when not directly reachable)")
	flag.Parse()

	if *backend == "" {
		log.Fatal("--backend is required")
	}

	cfg := engine.Config{
		ListenAddr:       *listen,
		BackendURL:       *backend,
		RulesPath:        *rules,
		DetectOnly:       *detectOnly,
		BlockStatus:      *blockStatus,
		RequestIDHdr:     "X-Request-ID",
		AnomalyThreshold: *anomalyThreshold,
		BlockTitle:       *blockTitle,
		BlockMessage:     *blockMessage,
		TrustProxyHeader: *trustProxy,
	}

	eng, err := engine.New(cfg)
	if err != nil {
		log.Fatalf("init engine: %v", err)
	}

	if *allowedIPs != "" || *blockedIPs != "" {
		allowed := splitCSV(*allowedIPs)
		blocked := splitCSV(*blockedIPs)
		if err := eng.SetIPLists(allowed, blocked); err != nil {
			log.Fatalf("ip lists: %v", err)
		}
	}

	if *rateLimit > 0 {
		eng.SetRateLimit(*rateLimitRoute, *rateLimit, time.Minute)
	}

	srv := &http.Server{
		Addr:              *listen,
		Handler:           eng.Handler(),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("waf-engine listening",
			"addr", *listen, "backend", *backend, "detect_only", *detectOnly,
			"rate_limit", *rateLimit, "anomaly_threshold", *anomalyThreshold)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
	os.Exit(0)
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	var out []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			if v := s[start:i]; v != "" {
				out = append(out, v)
			}
			start = i + 1
		}
	}
	return out
}