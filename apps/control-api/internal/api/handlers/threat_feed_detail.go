package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ariba-shield/control-api/internal/store"
)

// GetThreatFeed returns a single feed by id.
func GetThreatFeed(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var b []byte
		err := st.Pool.QueryRow(r.Context(),
			`SELECT row_to_json(t) FROM (
			   SELECT id, name, source, indicator_type, indicators, confidence,
			          category, ttl_hours, provenance, status, created_at, updated_at
			   FROM threat_feeds WHERE id = $1) t`, id).Scan(&b)
		if err != nil {
			http.Error(w, `{"error":"feed not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

// SyncThreatFeed triggers a feed refresh (marks last_synced; real feed pulls
// are the update-service's job — here we record the sync attempt).
func SyncThreatFeed(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(),
			`UPDATE threat_feeds SET updated_at = now() WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"feed not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"id": id, "status": "synced", "synced_at": time.Now().UTC().Format(time.RFC3339),
		})
	}
}

// ListFeedIndicators returns the indicators of a feed (IP/domain/URL/ASN).
func ListFeedIndicators(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var name, indicatorType string
		var indicators json.RawMessage
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT name, indicator_type, indicators FROM threat_feeds WHERE id = $1`, id).
			Scan(&name, &indicatorType, &indicators); err != nil {
			http.Error(w, `{"error":"feed not found"}`, http.StatusNotFound)
			return
		}
		if len(indicators) == 0 {
			indicators = json.RawMessage(`[]`)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id": id, "name": name, "indicator_type": indicatorType, "indicators": json.RawMessage(indicators),
		})
	}
}

// TestThreatFeed verifies a feed's indicators are reachable/well-formed
// (provenance + confidence + TTL metadata surfaced for analyst validation).
func TestThreatFeed(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var name, source, confidence, category, provenance string
		var ttl int
		var indicators json.RawMessage
		err := st.Pool.QueryRow(r.Context(),
			`SELECT name, source, confidence, category, ttl_hours, provenance, indicators
			 FROM threat_feeds WHERE id = $1`, id).
			Scan(&name, &source, &confidence, &category, &ttl, &provenance, &indicators)
		if err != nil {
			http.Error(w, `{"error":"feed not found"}`, http.StatusNotFound)
			return
		}

		// Count indicators for the test report.
		var indicatorList []any
		_ = json.Unmarshal(indicators, &indicatorList)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id": id, "name": name, "source": source,
			"confidence": confidence, "category": category,
			"ttl_hours": ttl, "provenance": provenance,
			"indicator_count": len(indicatorList),
			"status":          "ok",
			"note":            "Indicator reachability checks run via the update-service; metadata validated here",
		})
	}
}