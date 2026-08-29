package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ariba-shield/control-api/internal/store"
)

// IfMatch extracts the expected version from the If-Match or X-Expected-Version
// header. Returns the version and whether it was present.
func IfMatch(r *http.Request) (int64, bool) {
	if v := r.Header.Get("If-Match"); v != "" {
		v = strings.Trim(v, `"`)
		n, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return n, true
		}
	}
	if v := r.Header.Get("X-Expected-Version"); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return n, true
		}
	}
	return 0, false
}

// OptimisticConcurrency returns a handler wrapper that enforces a version
// check before a mutation (SRS FR-0.1-044). If the client supplies an
// If-Match / X-Expected-Version header and it does not match the current row
// version, it responds 409 Conflict so the client can retry with the latest
// version. If no header is supplied, the mutation proceeds.
func OptimisticConcurrency(st *store.Store, table, idField string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		expected, hasHeader := IfMatch(r)
		if !hasHeader {
			next(w, r)
			return
		}

		resourceID := r.PathValue(idField)
		if resourceID == "" {
			next(w, r)
			return
		}

		var currentVersion int64
		err := st.Pool.QueryRow(r.Context(),
			fmt.Sprintf("SELECT version FROM %s WHERE id = $1", table), resourceID).Scan(&currentVersion)
		if err != nil {
			// Row not found or query error: let the handler produce the proper 404.
			next(w, r)
			return
		}

		if currentVersion != expected {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Current-Version", fmt.Sprintf("%d", currentVersion))
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(
				fmt.Sprintf(`{"error":"conflict","code":"version_mismatch","expected_version":%d,"current_version":%d}`,
					expected, currentVersion)))
			return
		}

		next(w, r)
	}
}