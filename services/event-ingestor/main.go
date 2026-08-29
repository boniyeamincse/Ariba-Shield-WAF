package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ariba-shield/event-ingestor/internal/ingest"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	databaseURL := flag.String("database", "", "PostgreSQL URL")
	maxBatch := flag.Int("batch", 100, "max events per batch")
	maxWait := flag.Duration("wait", 5*time.Second, "max time to hold a batch")
	flag.Parse()

	if *databaseURL == "" {
		log.Fatal("--database is required")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := pgxpool.New(ctx, *databaseURL)
	if err != nil {
		log.Fatalf("db pool: %v", err)
	}
	defer pool.Close()

	in := ingest.New(pool)
	slog.Info("event-ingestor consuming stdin", "batch", *maxBatch, "wait", *maxWait)
	if err := in.Run(ctx, os.Stdin, *maxBatch, *maxWait); err != nil && err != context.Canceled {
		slog.Error("ingest run failed", "error", err)
		os.Exit(1)
	}
}
