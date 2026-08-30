package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ariba-shield/control-api/internal/store"
)

type exception struct {
	ID            string  `json:"id"`
	PolicyID      string  `json:"policy_id,omitempty"`
	ApplicationID string  `json:"application_id,omitempty"`
	RuleID        string  `json:"rule_id,omitempty"`
	URLPattern    string  `json:"url_pattern,omitempty"`
	Parameter     string  `json:"parameter,omitempty"`
	Reason        string  `json:"reason"`
	ExpiresAt     *string `json:"expires_at,omitempty"`
	Status        string  `json:"status"`
}

// ListExceptions returns policy exceptions (false-positive exclusions).
func ListExceptions(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, COALESCE(policy_id,''), COALESCE(application_id,''), COALESCE(rule_id,''),
			        COALESCE(url_pattern,''), COALESCE(parameter,''), reason, expires_at, status
			 FROM exceptions ORDER BY created_at DESC LIMIT 100`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		exceptions := []exception{}
		for rows.Next() {
			var e exception
			var expiresAt *time.Time
			if err := rows.Scan(&e.ID, &e.PolicyID, &e.ApplicationID, &e.RuleID,
				&e.URLPattern, &e.Parameter, &e.Reason, &expiresAt, &e.Status); err != nil {
				continue
			}
			if expiresAt != nil {
				s := expiresAt.Format(time.RFC3339)
				e.ExpiresAt = &s
			}
			exceptions = append(exceptions, e)
		}
		if exceptions == nil {
			exceptions = []exception{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(exceptions)
	}
}

// CreateException creates a time-limited policy exception.
func CreateException(st *store.Store) http.HandlerFunc {
	type create struct {
		PolicyID      string `json:"policy_id"`
		ApplicationID string `json:"application_id"`
		RuleID        string `json:"rule_id"`
		URLPattern    string `json:"url_pattern"`
		Parameter     string `json:"parameter"`
		Reason        string `json:"reason"`
		ExpiresAt     string `json:"expires_at"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body create
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.Reason == "" {
			http.Error(w, `{"error":"reason (owner justification) required"}`, http.StatusBadRequest)
			return
		}

		id, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}

		var expiresAt any
		if body.ExpiresAt != "" {
			if t, err := time.Parse(time.RFC3339, body.ExpiresAt); err == nil {
				expiresAt = t
			}
		}

		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO exceptions
			   (id, organization_id, policy_id, application_id, rule_id, url_pattern, parameter, reason, expires_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			id, orgID, nullIfEmpty(body.PolicyID), nullIfEmpty(body.ApplicationID),
			nullIfEmpty(body.RuleID), nullIfEmpty(body.URLPattern), nullIfEmpty(body.Parameter),
			body.Reason, expiresAt); err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}
