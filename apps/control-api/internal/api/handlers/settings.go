package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

const orgID = "01ARZ3NDEKTSV4RRFFQ69G5FAV"

// GetSettings returns all settings for a category (or all if category empty).
func GetSettings(st *store.Store, category string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var catFilter string
		args := []any{orgID}
		if category != "" {
			catFilter = " AND category = $2"
			args = append(args, category)
		}
		q, err := st.Pool.Query(r.Context(),
			`SELECT category, key, value FROM system_settings WHERE organization_id = $1`+catFilter+` ORDER BY category, key`, args...)
		if err != nil {
			http.Error(w, `{"error":"db query failed"}`, http.StatusInternalServerError)
			return
		}
		defer q.Close()

		grouped := map[string]map[string]any{}
		for q.Next() {
			var cat, key string
			var val json.RawMessage
			if err := q.Scan(&cat, &key, &val); err != nil {
				continue
			}
			if grouped[cat] == nil {
				grouped[cat] = map[string]any{}
			}
			grouped[cat][key] = json.RawMessage(val)
		}
		if category != "" {
			// Return just the object for that category.
			w.Header().Set("Content-Type", "application/json")
			if grouped[category] == nil {
				grouped[category] = map[string]any{}
			}
			json.NewEncoder(w).Encode(grouped[category])
			return
		}
		if grouped == nil {
			grouped = map[string]map[string]any{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(grouped)
	}
}

// UpdateSettings upserts settings for a category from a JSON object body.
func UpdateSettings(st *store.Store, category string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		for k, v := range body {
			b, err := json.Marshal(v)
			if err != nil {
				continue
			}
			_, err = st.Pool.Exec(r.Context(),
				`INSERT INTO system_settings (id, organization_id, category, key, value, updated_at)
				 VALUES ($1, $2, $3, $4, $5, now())
				 ON CONFLICT (organization_id, category, key)
				 DO UPDATE SET value = EXCLUDED.value, updated_at = now()`,
				mustNewID(st), orgID, category, k, json.RawMessage(b))
			if err != nil {
				http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "updated", "category": category})
	}
}

// mustNewID generates a ULID; falls back to a timestamp string on failure.
func mustNewID(st *store.Store) string {
	id, err := st.NewID()
	if err != nil {
		return "unknown"
	}
	return id
}

// GetSecuritySettings returns security settings.
func GetSecuritySettings(st *store.Store) http.HandlerFunc {
	return GetSettings(st, "security")
}

// UpdateSecuritySettings updates security settings.
func UpdateSecuritySettings(st *store.Store) http.HandlerFunc {
	return UpdateSettings(st, "security")
}

// GetLocalizationSettings returns localization settings.
func GetLocalizationSettings(st *store.Store) http.HandlerFunc {
	return GetSettings(st, "localization")
}

// UpdateLocalizationSettings updates localization settings.
func UpdateLocalizationSettings(st *store.Store) http.HandlerFunc {
	return UpdateSettings(st, "localization")
}

// GetRetentionSettings returns retention settings.
func GetRetentionSettings(st *store.Store) http.HandlerFunc {
	return GetSettings(st, "retention")
}

// UpdateRetentionSettings updates retention settings.
func UpdateRetentionSettings(st *store.Store) http.HandlerFunc {
	return UpdateSettings(st, "retention")
}
