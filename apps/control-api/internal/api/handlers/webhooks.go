package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

type webhook struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	URL              string            `json:"url"`
	Severity         []string          `json:"severity"`
	EventTypes       []string          `json:"event_types"`
	CustomHeaders    map[string]string `json:"custom_headers"`
	MaxRetryAttempts int               `json:"max_retry_attempts"`
	Status           string            `json:"status"`
}

// ListWebhooks returns all configured webhooks.
func ListWebhooks(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, name, url, severity, event_types, COALESCE(custom_headers::text,'{}'), max_retry_attempts, status
			 FROM webhooks ORDER BY name`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		webhooks := []webhook{}
		for rows.Next() {
			var w webhook
			var headers string
			if err := rows.Scan(&w.ID, &w.Name, &w.URL, &w.Severity, &w.EventTypes, &headers, &w.MaxRetryAttempts, &w.Status); err != nil {
				continue
			}
			_ = json.Unmarshal([]byte(headers), &w.CustomHeaders)
			webhooks = append(webhooks, w)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(webhooks)
	}
}

// CreateWebhook creates a webhook (no plaintext secret stored in the event
// table; webhook tokens are out of scope until the secrets module).
func CreateWebhook(st *store.Store) http.HandlerFunc {
	type create struct {
		Name             string            `json:"name"`
		URL              string            `json:"url"`
		Severity         []string          `json:"severity"`
		EventTypes       []string          `json:"event_types"`
		CustomHeaders    map[string]string `json:"custom_headers"`
		MaxRetryAttempts int               `json:"max_retry_attempts"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body create
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.Name == "" || body.URL == "" {
			http.Error(w, `{"error":"name and url required"}`, http.StatusBadRequest)
			return
		}

		id, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}

		headersJSON, _ := json.Marshal(body.CustomHeaders)
		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO webhooks (id, organization_id, name, url, severity, event_types, custom_headers, max_retry_attempts)
			 VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb, $8)`,
			id, orgID, body.Name, body.URL, body.Severity, body.EventTypes, headersJSON, body.MaxRetryAttempts); err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}
