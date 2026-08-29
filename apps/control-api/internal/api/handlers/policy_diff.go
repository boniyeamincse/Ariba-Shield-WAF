package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

// DiffPolicyVersions returns a field-level diff between two policy versions.
func DiffPolicyVersions(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fromID := r.URL.Query().Get("from")
		toID := r.URL.Query().Get("to")
		if fromID == "" || toID == "" {
			http.Error(w, `{"error":"from and to required"}`, http.StatusBadRequest)
			return
		}

		var fromDoc, toDoc map[string]any
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT document FROM policy_versions WHERE id = $1`, fromID).Scan(&fromDoc); err != nil {
			http.Error(w, `{"error":"from version not found"}`, http.StatusNotFound)
			return
		}
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT document FROM policy_versions WHERE id = $1`, toID).Scan(&toDoc); err != nil {
			http.Error(w, `{"error":"to version not found"}`, http.StatusNotFound)
			return
		}

		fromJSON, _ := json.Marshal(fromDoc)
		toJSON, _ := json.Marshal(toDoc)
		diff := simpleDiff(fromJSON, toJSON)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(diff)
	}
}

// simpleDiff produces a shallow field-level diff of two JSON documents.
func simpleDiff(from, to []byte) map[string]any {
	var f, t map[string]any
	_ = json.Unmarshal(from, &f)
	_ = json.Unmarshal(to, &t)

	out := map[string]any{
		"added":   []string{},
		"removed": []string{},
		"changed": []string{},
	}

	// Detect added/changed fields in t.
	for k, tv := range t {
		fv, exists := f[k]
		if !exists {
			out["added"] = append(out["added"].([]string), k)
			continue
		}
		fjs, _ := json.Marshal(fv)
		tjs, _ := json.Marshal(tv)
		if string(fjs) != string(tjs) {
			out["changed"] = append(out["changed"].([]string), k)
		}
	}
	// Detect removed fields (in f but not t).
	for k := range f {
		if _, exists := t[k]; !exists {
			out["removed"] = append(out["removed"].([]string), k)
		}
	}

	if len(out["added"].([]string)) == 0 {
		out["added"] = []string{}
	}
	if len(out["removed"].([]string)) == 0 {
		out["removed"] = []string{}
	}
	if len(out["changed"].([]string)) == 0 {
		out["changed"] = []string{}
	}

	return out
}