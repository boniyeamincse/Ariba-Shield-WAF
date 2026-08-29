package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

// GetBotPolicy returns a single bot policy.
func GetBotPolicy(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var b []byte
		err := st.Pool.QueryRow(r.Context(),
			`SELECT row_to_json(t) FROM (
			   SELECT id, application_id, name, challenge_type, known_bots,
			          automation_signals, login_protection, scrape_protection, status,
			          created_at, updated_at
			   FROM bot_policies WHERE id = $1) t`, id).Scan(&b)
		if err != nil {
			http.Error(w, `{"error":"bot policy not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

// ListBotEvents returns recent bot detection events.
func ListBotEvents(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, COALESCE(application_id,''), client_ip, classification, reason, action, created_at::text
			 FROM bot_events ORDER BY created_at DESC LIMIT 100`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		events := []map[string]any{}
		for rows.Next() {
			var id, app, ip, class, reason, action, created string
			if err := rows.Scan(&id, &app, &ip, &class, &reason, &action, &created); err == nil {
				events = append(events, map[string]any{
					"id": id, "application_id": app, "client_ip": ip,
					"classification": class, "reason": reason, "action": action, "created_at": created,
				})
			}
		}
		if events == nil {
			events = []map[string]any{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"events": events})
	}
}

// ListBotClients returns tracked client classifications.
func ListBotClients(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, COALESCE(application_id,''), client_ip, classification, confidence,
			        user_agent, score, action, status, first_seen_at::text, last_seen_at::text
			 FROM bot_clients ORDER BY last_seen_at DESC LIMIT 100`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		clients := []map[string]any{}
		for rows.Next() {
			var id, app, ip, class, conf, ua, action, status, first, last string
			var score int
			if err := rows.Scan(&id, &app, &ip, &class, &conf, &ua, &score, &action, &status, &first, &last); err == nil {
				clients = append(clients, map[string]any{
					"id": id, "application_id": app, "client_ip": ip, "classification": class,
					"confidence": conf, "user_agent": ua, "score": score, "action": action,
					"status": status, "first_seen_at": first, "last_seen_at": last,
				})
			}
		}
		if clients == nil {
			clients = []map[string]any{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"clients": clients})
	}
}

// RevokeBotChallenge revokes an issued challenge.
func RevokeBotChallenge(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(),
			`UPDATE bot_challenges SET status = 'revoked' WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"challenge not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "revoked"})
	}
}