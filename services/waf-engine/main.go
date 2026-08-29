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
	flag.Parse()

	if *backend == "" {
		log.Fatal("--backend is required")
	}

	cfg := engine.Config{
		ListenAddr:   *listen,
		BackendURL:   *backend,
		RulesPath:    *rules,
		DetectOnly:   *detectOnly,
		BlockStatus:  *blockStatus,
		RequestIDHdr: "X-Request-ID",
	}

	eng, err := engine.New(cfg)
	if err != nil {
		log.Fatalf("init engine: %v", err)
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
		slog.Info("waf-engine listening", "addr", *listen, "backend", *backend, "detect_only", *detectOnly)
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