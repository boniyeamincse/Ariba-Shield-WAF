package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ariba-shield/control-api/internal/store"
)

// GetSecurityEvent returns a single event with its full record.
func GetSecurityEvent(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var b []byte
		err := st.Pool.QueryRow(r.Context(),
			`SELECT row_to_json(t) FROM (
			   SELECT id, event_id, request_id, gateway_id, application_id, virtual_server_id,
			          client_ip, method, path, host, status, severity, decision_action, reason,
			          rule_ids, match_details, masked, created_at
			   FROM security_events WHERE id = $1) t`, id).Scan(&b)
		if err != nil {
			http.Error(w, `{"error":"event not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

// EventMatches returns the matched rules + match details of an event.
func EventMatches(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var ruleIDs []string
		var matchDetails json.RawMessage
		var reason, decision, severity string
		err := st.Pool.QueryRow(r.Context(),
			`SELECT rule_ids, COALESCE(match_details,'{}'), COALESCE(reason,''),
			        COALESCE(decision_action,''), COALESCE(severity,'')
			 FROM security_events WHERE id = $1`, id).
			Scan(&ruleIDs, &matchDetails, &reason, &decision, &severity)
		if err != nil {
			http.Error(w, `{"error":"event not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"event_id": id, "rule_ids": ruleIDs, "match_details": json.RawMessage(matchDetails),
			"reason": reason, "decision_action": decision, "severity": severity,
		})
	}
}

// EventTimeline reconstructs the decision path timeline for an event
// (matched at → decided → correlated incidents) from event metadata.
func EventTimeline(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var requestID, clientIP, decision, reason, createdAt string
		var ruleIDs []string
		err := st.Pool.QueryRow(r.Context(),
			`SELECT request_id, COALESCE(client_ip,''), COALESCE(decision_action,''),
			        COALESCE(reason,''), rule_ids, created_at::text
			 FROM security_events WHERE id = $1`, id).
			Scan(&requestID, &clientIP, &decision, &reason, &ruleIDs, &createdAt)
		if err != nil {
			http.Error(w, `{"error":"event not found"}`, http.StatusNotFound)
			return
		}

		timeline := []map[string]any{
			{"step": "detected", "time": createdAt, "detail": fmt.Sprintf("Event %s detected for %s", id, clientIP)},
			{"step": "matched", "time": createdAt, "detail": fmt.Sprintf("Matched rules: %v", ruleIDs)},
			{"step": "decided", "time": createdAt, "detail": fmt.Sprintf("Decision: %s (%s)", decision, reason)},
			{"step": "correlated", "time": createdAt, "detail": fmt.Sprintf("Request ID: %s", requestID)},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"event_id": id, "timeline": timeline})
	}
}

// MaskSecurityEvent toggles the masked flag (payload retention control).
func MaskSecurityEvent(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(),
			`UPDATE security_events SET masked = true WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"event not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "masked": "true"})
	}
}

// ExportSecurityEvent exports an event as JSON (for SIEM/analyst export).
func ExportSecurityEvent(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var b []byte
		err := st.Pool.QueryRow(r.Context(),
			`SELECT row_to_json(t) FROM (
			   SELECT id, event_id, request_id, gateway_id, application_id, virtual_server_id,
			          client_ip, method, path, host, status, severity, decision_action, reason,
			          rule_ids, match_details, masked, created_at
			   FROM security_events WHERE id = $1) t`, id).Scan(&b)
		if err != nil {
			http.Error(w, `{"error":"event not found"}`, http.StatusNotFound)
			return
		}
		// Pretty-print for export.
		var pretty map[string]any
		_ = json.Unmarshal(b, &pretty)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="event-%s.json"`, id))
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		enc.Encode(map[string]any{"exported_at": time.Now().UTC().Format(time.RFC3339), "event": pretty})
	}
}