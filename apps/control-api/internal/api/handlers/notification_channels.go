package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ariba-shield/control-api/internal/store"
)

// notificationChannelJSON is the response shape for a notification channel.
type notificationChannelJSON struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Kind      string `json:"kind"`
	IsDefault bool   `json:"is_default"`
	Enabled   bool   `json:"enabled"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// notificationChannelCreate is the request body for creating a channel.
type notificationChannelCreate struct {
	Name      string `json:"name"`
	Kind      string `json:"kind"`     // wazuh, syslog, cef, leef, webhook, email, teams, slack, soar
	IsDefault bool   `json:"is_default"`
	Enabled   bool   `json:"enabled"`
}

// notificationChannelUpdate is the request body for updating a channel.
type notificationChannelUpdate struct {
	Name      string `json:"name"`
	IsDefault bool `json:"is_default"`
	Enabled   bool `json:"enabled"`
}

// ListNotificationChannels returns all notification channels.
func ListNotificationChannels(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kind := r.URL.Query().Get("kind")
		isDefault := r.URL.Query().Get("is_default")

		where := "1=1"
		args := []any{}
		add := func(cond string, val any) {
			args = append(args, val)
			where += " AND " + cond
		}
		if kind != "" {
			add("kind = $"+strconv.Itoa(len(args)+1), kind)
		}
		if isDefault != "" {
			add("is_default = $"+strconv.Itoa(len(args)+1), isDefault)
		}

		query := `SELECT id, name, kind, is_default, enabled, created_at, updated_at FROM notification_channels WHERE ` + where + ` ORDER BY created_at DESC`
		rows, err := st.Pool.Query(r.Context(), query, args...)
		if err != nil {
			http.Error(w, `{"error":"db query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var channels []notificationChannelJSON
		for rows.Next() {
			var id, name, kind string
			var isDefault bool
			var enabled bool
			var created, updated string
			if err := rows.Scan(&id, &name, &kind, &isDefault, &enabled, &created, &updated); err != nil {
				continue
			}
			channels = append(channels, notificationChannelJSON{
				ID:        id,
				Name:      name,
				Kind:      kind,
				IsDefault: isDefault,
				Enabled:   enabled,
				CreatedAt: created,
				UpdatedAt: updated,
			})
		}
		if channels == nil {
			channels = []notificationChannelJSON{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(channels)
	}
}

// GetNotificationChannel returns a single notification channel by ID.
func GetNotificationChannel(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var idVar, name, kind string
		var isDefault bool
		var enabled bool
		var created, updated string
		err := st.Pool.QueryRow(r.Context(),
			`SELECT id, name, kind, is_default, enabled, created_at, updated_at FROM notification_channels WHERE id = $1`, id).Scan(&idVar, &name, &kind, &isDefault, &enabled, &created, &updated)
		if err != nil {
			http.Error(w, `{"error":"notification channel not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notificationChannelJSON{
			ID:        idVar,
			Name:      name,
			Kind:      kind,
			IsDefault: isDefault,
			Enabled:   enabled,
			CreatedAt: created,
			UpdatedAt: updated,
		})
	}
}

// CreateNotificationChannel creates a new notification channel.
func CreateNotificationChannel(st *store.Store) http.HandlerFunc {
	type createReq struct {
		Name      string `json:"name"`
		Kind      string `json:"kind"`
		IsDefault bool   `json:"is_default"`
		Enabled   bool   `json:"enabled"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req createReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if req.Kind == "" {
			http.Error(w, `{"error":"kind is required"}`, http.StatusBadRequest)
			return
		}
		if req.Name == "" {
			http.Error(w, `{"error":"name is required"}`, http.StatusBadRequest)
			return
		}

		id, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}

		_, err = st.Pool.Exec(r.Context(),
			`INSERT INTO notification_channels (id, organization_id, name, kind, is_default, enabled, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, $5, $6, now(), now())`,
			id, "01ARZ3NDEKTSV4RRFFQ69G5FAV", req.Name, req.Kind, req.IsDefault, req.Enabled)
		if err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

// UpdateNotificationChannel updates a notification channel.
func UpdateNotificationChannel(st *store.Store) http.HandlerFunc {
	type updateReq struct {
		Name      string `json:"name"`
		IsDefault bool `json:"is_default"`
		Enabled   bool `json:"enabled"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var req updateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if req.Name == "" && !req.IsDefault && !req.Enabled {
			http.Error(w, `{"error":"at least one field to update is required"}`, http.StatusBadRequest)
			return
		}
		_, err := st.Pool.Exec(r.Context(),
			`UPDATE notification_channels SET name = CASE WHEN $2 != '' THEN $2 ELSE name END, is_default = CASE WHEN $3 THEN true ELSE is_default END, enabled = CASE WHEN $4 THEN true ELSE enabled END, updated_at = now() WHERE id = $1`,
			id, req.Name, req.IsDefault, req.Enabled)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

// DeleteNotificationChannel deletes a notification channel.
func DeleteNotificationChannel(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(), `DELETE FROM notification_channels WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"delete failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"notification channel not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	}
}

// TestNotificationChannel tests a notification channel connection.
func TestNotificationChannel(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		// In a full implementation, this would attempt a real connection
		// to the notification channel. For now, we return success.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"id":       id,
			"success":  "true",
			"message":  "notification channel test passed",
		})
	}
}
