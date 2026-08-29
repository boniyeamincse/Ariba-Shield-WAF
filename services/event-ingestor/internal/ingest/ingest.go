package ingest

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
)

// Event is the normalized shape written to PostgreSQL.
// Only fields with no sensitive payloads are stored (masking rule).
type Event struct {
	EventID      string            `json:"event_id"`
	RequestID    string            `json:"request_id"`
	Timestamp    string            `json:"timestamp"`
	EventType    string            `json:"event_type"`
	GatewayID    string            `json:"gateway_id"`
	VirtualSrvID string            `json:"virtual_server_id"`
	AppID        string            `json:"application_id"`
	ClientIP     string            `json:"client_ip"`
	Method       string            `json:"method"`
	Path         string            `json:"path"`
	Host         string            `json:"host"`
	Status       int               `json:"status"`
	Severity     string            `json:"severity"`
	Decision     string            `json:"decision_action"`
	Reason       string            `json:"reason"`
	RuleIDs      []string          `json:"rule_ids"`
	MatchDetails []json.RawMessage `json:"match_details"`
	Raw          json.RawMessage   `json:"raw"`
}

// Ingestor consumes JSON-lines events from an io.Reader and writes them to
// PostgreSQL in batches. Never blocks or crashes the producer: on DB failure
// it backs off and resumes; events in-flight may be dropped (best-effort).
type Ingestor struct {
	pool *pgxpool.Pool
	buf  []Event
}

// New creates an Ingestor bound to a DB pool.
func New(pool *pgxpool.Pool) *Ingestor {
	return &Ingestor{pool: pool}
}

// Run reads events from r, batched by maxBatch or maxWait, until ctx done or EOF.
func (in *Ingestor) Run(ctx context.Context, r io.Reader, maxBatch int, maxWait time.Duration) error {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)

	tick := time.NewTicker(maxWait)
	defer tick.Stop()

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var ev Event
		if err := json.Unmarshal(line, &ev); err != nil {
			// Malformed line: log and continue; never kill the pipeline.
			slog.Warn("malformed event line", "error", err)
			continue
		}
		in.buf = append(in.buf, ev)

		if len(in.buf) >= maxBatch {
			if err := in.Flush(ctx); err != nil {
				slog.Error("flush failed", "error", err)
			}
		}

		select {
		case <-tick.C:
			if err := in.Flush(ctx); err != nil {
				slog.Error("tick flush failed", "error", err)
			}
		default:
		}

		select {
		case <-ctx.Done():
			return in.Flush(ctx)
		default:
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return in.Flush(ctx)
}

// Flush writes the pending batch to PostgreSQL.
func (in *Ingestor) Flush(ctx context.Context) error {
	if len(in.buf) == 0 {
		return nil
	}
	batch := in.buf
	in.buf = nil

	tx, err := in.pool.Begin(ctx)
	if err != nil {
		// Re-queue on failure so data is not silently lost for a transient error.
		in.buf = append(in.buf, batch...)
		return err
	}
	defer tx.Rollback(ctx)

	for _, ev := range batch {
		id := ulid.Make().String()
		ts := ev.Timestamp
		if ts == "" {
			ts = time.Now().UTC().Format(time.RFC3339)
		}
		_, err := tx.Exec(ctx, `
			INSERT INTO security_events
			  (id, event_id, request_id, gateway_id, application_id, virtual_server_id,
			   client_ip, method, path, host, status, severity, decision_action, reason,
			   rule_ids, match_details, masked, raw, created_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,true,$17,$18::timestamptz)`,
			id, ev.EventID, ev.RequestID, ev.GatewayID, ev.AppID, ev.VirtualSrvID,
			ev.ClientIP, ev.Method, ev.Path, ev.Host, ev.Status, ev.Severity, ev.Decision, ev.Reason,
			ev.RuleIDs, ev.MatchDetails, ev.Raw, ts)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
