package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
)

type auditKey string

const auditActorKey auditKey = "audit_actor"

// ContextWithActor stores the user ID for audit logging.
func ContextWithActor(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, auditActorKey, userID)
}

// AuditFromContext returns the actor user ID.
func AuditFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(auditActorKey).(string); ok {
		return v
	}
	return ""
}

// Audit logs all mutations (POST/PUT/PATCH/DELETE) to the audit_events table.
// Wraps the handler, captures the request body for the after_ref, and writes
// the audit event asynchronously (fire-and-forget with a goroutine).
func Audit(pool *pgxpool.Pool, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" || r.Method == "HEAD" || r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		// Wrap the response writer to capture the status code.
		sw := &statusWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(sw, r)

		// Fire-and-forget audit: never block the traffic path.
		go func() {
			actor := AuditFromContext(r.Context())
			reqID := RequestIDFromContext(r.Context())
			id := ulid.Make().String()
			// Use NULL actor for unauthenticated actions (e.g. login) so the
			// audit_events FK to users(id) is not violated by an empty string.
			var actorID any
			if actor != "" {
				actorID = actor
			}
			_, err := pool.Exec(context.Background(),
				`INSERT INTO audit_events (id, organization_id, actor_user_id, action, resource_type, resource_id, ip, request_id, created_at)
				 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
				id, "01ARZ3NDEKTSV4RRFFQ69G5FAV", actorID, r.Method,
				r.URL.Path, r.URL.Path, r.RemoteAddr, reqID, time.Now().UTC())
			if err != nil {
				slog.Warn("audit write failed", "error", err, "path", r.URL.Path)
			}
		}()
	})
}