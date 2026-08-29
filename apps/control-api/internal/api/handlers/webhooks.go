package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

type webhook struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	URL               string            `json:"url"`
	Severity          []string          `json:"severity"`
	EventTypes        []string          `json:"event_types"`
	CustomHeaders     map[string]string `json:"custom_headers"`
	MaxRetryAttempts  int               `json:"max_retry_attempts"`
	Status            string            `json:"status"`
}

// ListWebhooks returns all configured SOAR/Incident webhooks
func ListWebhooks(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Mock response for now, to be replaced with PostgreSQL query
		webhooks := []webhook{
			{
				ID:               "wh_1",
				Name:             "SOC-Playbook",
				URL:              "https://soc.company.com/incident/create",
				Severity:         []string{"HIGH", "CRITICAL"},
				EventTypes:       []string{"SQL_INJECTION", "RCE"},
				CustomHeaders:    map[string]string{"X-API-Key": "encrypted_value"},
				MaxRetryAttempts: 3,
				Status:           "active",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(webhooks)
	}
}

// CreateWebhook registers a new webhook
func CreateWebhook(st *store.Store) http.HandlerFunc {
	type webhookCreate struct {
		Name          string            `json:"name"`
		URL           string            `json:"url"`
		SecretToken   string            `json:"secret_token"`
		Severity      []string          `json:"severity"`
		EventTypes    []string          `json:"event_types"`
		CustomHeaders map[string]string `json:"custom_headers"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body webhookCreate
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}

		if body.URL == "" || body.Name == "" {
			http.Error(w, `{"error":"name and url are required"}`, http.StatusBadRequest)
			return
		}

		// TODO: Encrypt SecretToken before storing into PostgreSQL
		newID, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": newID})
	}
}
