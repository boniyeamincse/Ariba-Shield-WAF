package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

// GetIPList returns a single IP list.
func GetIPList(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var b []byte
		err := st.Pool.QueryRow(r.Context(),
			`SELECT row_to_json(t) FROM (
			   SELECT id, name, list_type, entries, description, version, created_at, updated_at
			   FROM ip_lists WHERE id = $1) t`, id).Scan(&b)
		if err != nil {
			http.Error(w, `{"error":"ip list not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

// ListIPEntries returns the entries of an IP list.
func ListIPEntries(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var name, listType string
		var entries []string
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT name, list_type, entries FROM ip_lists WHERE id = $1`, id).Scan(&name, &listType, &entries); err != nil {
			http.Error(w, `{"error":"ip list not found"}`, http.StatusNotFound)
			return
		}
		if entries == nil {
			entries = []string{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"id": id, "name": name, "list_type": listType, "entries": entries})
	}
}

// AddIPEntry appends an entry to an IP list (deduplicated).
func AddIPEntry(st *store.Store) http.HandlerFunc {
	type add struct {
		Entry string `json:"entry"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var body add
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Entry == "" {
			http.Error(w, `{"error":"entry required"}`, http.StatusBadRequest)
			return
		}

		// Validate it's a plausible IP/CIDR.
		if !looksLikeCIDR(body.Entry) {
			http.Error(w, `{"error":"invalid IP/CIDR entry"}`, http.StatusBadRequest)
			return
		}

		var entries []string
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT entries FROM ip_lists WHERE id = $1`, id).Scan(&entries); err != nil {
			http.Error(w, `{"error":"ip list not found"}`, http.StatusNotFound)
			return
		}

		// Deduplicate.
		exists := false
		for _, e := range entries {
			if e == body.Entry {
				exists = true
				break
			}
		}
		if !exists {
			entries = append(entries, body.Entry)
		}

		if _, err := st.Pool.Exec(r.Context(),
			`UPDATE ip_lists SET entries = $1, version = version + 1, updated_at = now() WHERE id = $2`,
			entries, id); err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"id": id, "entry": body.Entry, "added": !exists})
	}
}

// DeleteIPEntry removes an entry from an IP list.
func DeleteIPEntry(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		entryID := r.PathValue("entryId")

		var entries []string
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT entries FROM ip_lists WHERE id = $1`, id).Scan(&entries); err != nil {
			http.Error(w, `{"error":"ip list not found"}`, http.StatusNotFound)
			return
		}

		removed := false
		kept := []string{}
		for _, e := range entries {
			if e == entryID {
				removed = true
				continue
			}
			kept = append(kept, e)
		}
		if !removed {
			http.Error(w, `{"error":"entry not found"}`, http.StatusNotFound)
			return
		}

		if _, err := st.Pool.Exec(r.Context(),
			`UPDATE ip_lists SET entries = $1, version = version + 1, updated_at = now() WHERE id = $2`,
			kept, id); err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"id": id, "removed": entryID})
	}
}

// looksLikeCIDR does a light sanity check (IP or IP/CIDR form). Full parsing
// is done by the engine's iplist package at runtime.
func looksLikeCIDR(s string) bool {
	hasDot := false
	hasColon := false
	for _, c := range s {
		switch {
		case c == '.':
			hasDot = true
		case c == ':':
			hasColon = true
		case c == '/' || c == '[' || c == ']' || c == 'x' || c == 'X' || c == 'a' || c == 'b' || c == 'c' || c == 'd' || c == 'e' || c == 'f':
		case c >= '0' && c <= '9':
		default:
			return false
		}
	}
	return hasDot || hasColon
}