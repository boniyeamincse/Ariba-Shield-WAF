package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

// ListVirtualServers returns all virtual servers.
func ListVirtualServers(st *store.Store) http.HandlerFunc {
	type vs struct {
		ID                  string `json:"id"`
		Name                string `json:"name"`
		ListenAddr          string `json:"listen_addr"`
		ListenPort          int    `json:"listen_port"`
		TLSEnabled          bool   `json:"tls_enabled"`
		CertificateRef      string `json:"certificate_ref,omitempty"`
		DefaultBackendPoolID string `json:"default_backend_pool_id,omitempty"`
		Version             int64  `json:"version"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, name, listen_addr, listen_port, tls_enabled,
			        COALESCE(certificate_ref,''), COALESCE(default_backend_pool_id,''), version
			 FROM virtual_servers ORDER BY name`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var items []vs
		for rows.Next() {
			var v vs
			if err := rows.Scan(&v.ID, &v.Name, &v.ListenAddr, &v.ListenPort, &v.TLSEnabled,
				&v.CertificateRef, &v.DefaultBackendPoolID, &v.Version); err != nil {
				http.Error(w, `{"error":"scan failed"}`, http.StatusInternalServerError)
				return
			}
			items = append(items, v)
		}
		if items == nil {
			items = []vs{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}

// CreateVirtualServer creates a virtual server.
func CreateVirtualServer(st *store.Store) http.HandlerFunc {
	type create struct {
		Name                string `json:"name"`
		ListenAddr          string `json:"listen_addr"`
		ListenPort          int    `json:"listen_port"`
		TLSEnabled          bool   `json:"tls_enabled"`
		CertificateRef      string `json:"certificate_ref"`
		DefaultBackendPoolID string `json:"default_backend_pool_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body create
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.Name == "" || body.ListenPort == 0 {
			http.Error(w, `{"error":"name and listen_port required"}`, http.StatusBadRequest)
			return
		}
		if body.ListenAddr == "" {
			body.ListenAddr = "0.0.0.0"
		}

		id, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}
		orgID := "01ARZ3NDEKTSV4RRFFQ69G5FAV"

		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO virtual_servers
			   (id, organization_id, name, listen_addr, listen_port, tls_enabled, certificate_ref, default_backend_pool_id)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			id, orgID, body.Name, body.ListenAddr, body.ListenPort, body.TLSEnabled,
			nullIfEmpty(body.CertificateRef), nullIfEmpty(body.DefaultBackendPoolID)); err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

// DeleteVirtualServer deletes a virtual server.
func DeleteVirtualServer(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(), `DELETE FROM virtual_servers WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"delete failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"virtual server not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	}
}