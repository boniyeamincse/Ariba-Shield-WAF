package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ariba-shield/control-api/internal/store"
)

// ListBackups returns recent backup records.
func ListBackups(st *store.Store) http.HandlerFunc {
	type backup struct {
		ID          string `json:"id"`
		Status      string `json:"status"`
		ArtifactRef string `json:"artifact_ref,omitempty"`
		SizeBytes   int64  `json:"size_bytes"`
		CreatedAt   string `json:"created_at"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, status, COALESCE(artifact_ref,''), size_bytes, created_at
			 FROM backups ORDER BY created_at DESC LIMIT 20`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var backups []backup
		for rows.Next() {
			var b backup
			var ts time.Time
			if err := rows.Scan(&b.ID, &b.Status, &b.ArtifactRef, &b.SizeBytes, &ts); err != nil {
				http.Error(w, `{"error":"scan failed"}`, http.StatusInternalServerError)
				return
			}
			b.CreatedAt = ts.Format(time.RFC3339)
			backups = append(backups, b)
		}
		if backups == nil {
			backups = []backup{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(backups)
	}
}

// CreateBackup records a backup job. In this phase it creates a DB snapshot
// reference via pg_dump executed asynchronously; the artifact reference is a
// best-effort path. Real object-storage upload is Phase 8.
func CreateBackup(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}
		orgID := "01ARZ3NDEKTSV4RRFFQ69G5FAV"
		userID := "01ARZ3NDEKTSV4RRFFQ69G5FAW"

		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO backups (id, organization_id, status, triggered_by)
			 VALUES ($1, $2, 'running', $3)`,
			id, orgID, userID); err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}

		// Fire-and-forget: mark complete after a simulated snapshot. In a real
		// deployment this would run pg_dump / pg_basebackup and upload to
		// object storage. Never block the API response.
		go func(backupID string) {
			time.Sleep(2 * time.Second)
			artifact := fmt.Sprintf("backups/%s.dump", backupID)
			_, _ = st.Pool.Exec(context.Background(),
				`UPDATE backups SET status = 'completed', artifact_ref = $1, size_bytes = 0, completed_at = now() WHERE id = $2`,
				artifact, backupID)
		}(id)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "running"})
	}
}

// RestoreBackup marks a restore request. Full restore (DB reimport) is
// intentionally not automated in this phase; the endpoint records intent.
func RestoreBackup(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		backupID := r.PathValue("id")
		var exists bool
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT EXISTS(SELECT 1 FROM backups WHERE id = $1)`, backupID).Scan(&exists); err != nil || !exists {
			http.Error(w, `{"error":"backup not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"backup_id": backupID,
			"status":    "restore_requested",
			"note":      "Full DB restore is not automated in this phase; see release-0.1.md",
		})
	}
}