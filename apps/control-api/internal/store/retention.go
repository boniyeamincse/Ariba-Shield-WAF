package store

import (
	"context"
	"log/slog"
	"time"
)

// RetentionPruner runs on a schedule to delete old records per retention settings.
type RetentionPruner struct {
	Store *Store
}

// RunOnce prunes audit_events and security_events based on the retention settings.
func (r *RetentionPruner) RunOnce(ctx context.Context) {
	auditDays := 365 // default
	eventDays := 90  // default

	// Read retention settings (values are raw JSON primitives).
	if st, err := r.Store.LoadSettings(ctx, "retention"); err == nil {
		auditDays = st.Int("audit_log_days", 365)
		eventDays = st.Int("security_event_days", 90)
	}

	// Prune audit_events.
	if auditDays > 0 {
		res, err := r.Store.Pool.Exec(ctx,
			`DELETE FROM audit_events WHERE created_at < now() - make_interval(days => $1)`, auditDays)
		if err != nil {
			slog.Warn("retention: audit_events prune failed", "error", err)
		} else {
			slog.Info("retention: audit_events pruned", "deleted", res.RowsAffected(), "retention_days", auditDays)
		}
	}

	// Prune security_events.
	if eventDays > 0 {
		res, err := r.Store.Pool.Exec(ctx,
			`DELETE FROM security_events WHERE created_at < now() - make_interval(days => $1)`, eventDays)
		if err != nil {
			slog.Warn("retention: security_events prune failed", "error", err)
		} else {
			slog.Info("retention: security_events pruned", "deleted", res.RowsAffected(), "retention_days", eventDays)
		}
	}
}

// StartLoop runs the retention pruner every 24 hours.
func (r *RetentionPruner) StartLoop(ctx context.Context) {
	slog.Info("retention: pruner started (24h interval)")
	r.RunOnce(ctx)
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			r.RunOnce(ctx)
		case <-ctx.Done():
			slog.Info("retention: pruner stopped")
			return
		}
	}
}