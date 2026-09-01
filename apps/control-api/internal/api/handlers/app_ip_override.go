package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

// ListAppIPOverrides returns IP overrides for an application.
func ListAppIPOverrides(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appID := r.PathValue("id")
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, list_type, ip, COALESCE(reason,''), created_at::text
			 FROM app_ip_overrides WHERE application_id = $1 ORDER BY created_at DESC`, appID)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		list := []map[string]any{}
		for rows.Next() {
			var id, lt, ip, reason, created string
			if rows.Scan(&id, &lt, &ip, &reason, &created) == nil {
				list = append(list, map[string]any{"id": id, "list_type": lt, "ip": ip, "reason": reason, "created_at": created})
			}
		}
		if list == nil {
			list = []map[string]any{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
	}
}

// AddAppIPOverride adds a block/allow IP override for an application.
func AddAppIPOverride(st *store.Store) http.HandlerFunc {
	type req struct {
		ListType string `json:"list_type"` // block | allow
		IP       string `json:"ip"`
		Reason   string `json:"reason"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		appID := r.PathValue("id")
		var body req
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.IP == "" {
			http.Error(w, `{"error":"ip required"}`, http.StatusBadRequest)
			return
		}
		if body.ListType == "" {
			body.ListType = "block"
		}
		id, _ := st.NewID()
		orgID := "01ARZ3NDEKTSV4RRFFQ69G5FAV"
		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO app_ip_overrides (id, organization_id, application_id, list_type, ip, reason)
			 VALUES ($1,$2,$3,$4,$5,$6)
			 ON CONFLICT (application_id, ip) DO UPDATE SET list_type = EXCLUDED.list_type, reason = EXCLUDED.reason`,
			id, orgID, appID, body.ListType, body.IP, body.Reason); err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "application_id": appID, "ip": body.IP, "list_type": body.ListType})
	}
}

// DeleteAppIPOverride removes an IP override.
func DeleteAppIPOverride(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appID := r.PathValue("id")
		entryID := r.PathValue("entryId")
		ct, err := st.Pool.Exec(r.Context(),
			`DELETE FROM app_ip_overrides WHERE id = $1 AND application_id = $2`, entryID, appID)
		if err != nil {
			http.Error(w, `{"error":"delete failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"override not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	}
}
