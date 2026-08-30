package handlers

import (
	"encoding/json"
	"fmt"
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

// simpleDiff produces a deep field-level diff of two JSON documents.
// Returns {added, changed, removed} where each value is a []string
// of dotted paths for nested fields (e.g. "rules.0.action").
func simpleDiff(from, to []byte) map[string]any {
	var f, t any
	_ = json.Unmarshal(from, &f)
	_ = json.Unmarshal(to, &t)

	var added, changed, removed []string
	diffValues("", f, t, &added, &changed, &removed)

	if added == nil {
		added = []string{}
	}
	if removed == nil {
		removed = []string{}
	}
	if changed == nil {
		changed = []string{}
	}

	return map[string]any{
		"added":   added,
		"removed": removed,
		"changed": changed,
	}
}

// diffValues recursively compares two JSON values and appends dotted paths.
func diffValues(path string, from, to any, added, changed, removed *[]string) {
	// If both are maps, recurse by key.
	fm, fIsMap := from.(map[string]any)
	tm, tIsMap := to.(map[string]any)
	if fIsMap && tIsMap {
		// added/changed in to
		for k, tv := range tm {
			key := joinPath(path, k)
			fv, ok := fm[k]
			if !ok {
				*added = append(*added, key)
				continue
			}
			if jsonEqual(fv, tv) {
				continue
			}
			if bothStructured(fv, tv) {
				diffValues(key, fv, tv, added, changed, removed)
			} else {
				*changed = append(*changed, key)
			}
		}
		// removed
		for k := range fm {
			if _, ok := tm[k]; !ok {
				*removed = append(*removed, joinPath(path, k))
			}
		}
		return
	}

	// If both are arrays, compare element by element.
	fa, fIsArr := from.([]any)
	ta, tIsArr := to.([]any)
	if fIsArr && tIsArr {
		max := len(fa)
		if len(ta) > max {
			max = len(ta)
		}
		for i := 0; i < max; i++ {
			key := fmt.Sprintf("%s[%d]", path, i)
			if i >= len(fa) {
				*added = append(*added, key)
				continue
			}
			if i >= len(ta) {
				*removed = append(*removed, key)
				continue
			}
			if jsonEqual(fa[i], ta[i]) {
				continue
			}
			if bothStructured(fa[i], ta[i]) {
				diffValues(key, fa[i], ta[i], added, changed, removed)
			} else {
				*changed = append(*changed, key)
			}
		}
		return
	}

	// One or both are scalars (or types differ) — mark as changed.
	if path != "" {
		*changed = append(*changed, path)
	}
}

func joinPath(parent, child string) string {
	if parent == "" {
		return child
	}
	return parent + "." + child
}

func bothStructured(a, b any) bool {
	_, aOk := a.(map[string]any)
	_, bOk := b.(map[string]any)
	if aOk && bOk {
		return true
	}
	_, aOk = a.([]any)
	_, bOk = b.([]any)
	return aOk && bOk
}

func jsonEqual(a, b any) bool {
	aj, _ := json.Marshal(a)
	bj, _ := json.Marshal(b)
	return string(aj) == string(bj)
}
