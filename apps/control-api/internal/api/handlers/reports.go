package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ariba-shield/control-api/internal/store"
)

// reportJSON is the response shape for a report.
type reportJSON struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Kind      string `json:"kind"`
	Status    string `json:"status"`
	CreatedBy string `json:"created_by"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// reportDetailJSON is the full report response including summary.
type reportDetailJSON struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Kind      string `json:"kind"`
	Status    string `json:"status"`
	Summary   any    `json:"summary"`
	CreatedBy string `json:"created_by"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ListReports returns all reports.
func ListReports(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kind := r.URL.Query().Get("kind")
		status := r.URL.Query().Get("status")

		where := "1=1"
		args := []any{}
		add := func(cond string, val any) {
			args = append(args, val)
			where += " AND " + cond
		}
		if kind != "" {
			add("kind = $"+strconv.Itoa(len(args)+1), kind)
		}
		if status != "" {
			add("status = $"+strconv.Itoa(len(args)+1), status)
		}

		query := `SELECT id, name, kind, status, COALESCE(created_by,''), created_at, updated_at FROM reports WHERE ` + where + ` ORDER BY created_at DESC`
		rows, err := st.Pool.Query(r.Context(), query, args...)
		if err != nil {
			http.Error(w, `{"error":"db query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var reports []reportJSON
		for rows.Next() {
			var id, name, kind, status, createdBy, created, updated string
			if err := rows.Scan(&id, &name, &kind, &status, &createdBy, &created, &updated); err != nil {
				continue
			}
			reports = append(reports, reportJSON{
				ID:        id,
				Name:      name,
				Kind:      kind,
				Status:    status,
				CreatedBy: createdBy,
				CreatedAt: created,
				UpdatedAt: updated,
			})
		}
		if reports == nil {
			reports = []reportJSON{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reports)
	}
}

// GetReport returns a single report by ID.
func GetReport(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var idVar, name, kind, status, createdBy, created, updated string
		var summary any
		var params any
		err := st.Pool.QueryRow(r.Context(),
			`SELECT id, name, kind, status, params, summary, COALESCE(created_by,''), created_at, updated_at FROM reports WHERE id = $1`, id).Scan(&idVar, &name, &kind, &status, &params, &summary, &createdBy, &created, &updated)
		if err != nil {
			http.Error(w, `{"error":"report not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reportDetailJSON{
			ID:        idVar,
			Name:      name,
			Kind:      kind,
			Status:    status,
			Summary:   summary,
			CreatedBy: createdBy,
			CreatedAt: created,
			UpdatedAt: updated,
		})
	}
}

// CreateReport creates a report from a kind + params in the body.
func CreateReport(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type createReq struct {
			Kind   string         `json:"kind"`
			Params map[string]any `json:"params"`
		}
		var req createReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if req.Kind == "" {
			http.Error(w, `{"error":"kind is required"}`, http.StatusBadRequest)
			return
		}
		if req.Params == nil {
			req.Params = map[string]any{}
		}
		generateReport(st, w, r, req.Kind)
	}
}

// DeleteReport removes a report.
func DeleteReport(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(), `DELETE FROM reports WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"delete failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"report not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	}
}

// generateReport runs a report generation query and writes the result.
func generateReport(st *store.Store, w http.ResponseWriter, r *http.Request, kind string) {
	id, err := st.NewID()
	if err != nil {
		http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
		return
	}

	// Determine the query and name based on kind.
	var query string
	var name string
	var summary map[string]any

	params := map[string]any{}
	switch kind {
	case "security":
		name = "Security Events Report"
		query = `SELECT COALESCE(severity,'unknown') as sev, COUNT(*) as cnt FROM security_events GROUP BY severity ORDER BY cnt DESC`
	case "traffic":
		name = "Traffic Summary Report"
		query = `SELECT COALESCE(method,'-') as method, COUNT(*) as cnt, AVG(latency_ms)::numeric(10,2) as avg_latency FROM access_events GROUP BY method ORDER BY cnt DESC`
	case "incidents":
		name = "Incidents Report"
		query = `SELECT COALESCE(severity,'unknown') as sev, COUNT(*) as cnt FROM incidents GROUP BY severity ORDER BY cnt DESC`
	case "compliance":
		name = "Compliance Summary Report"
		query = `SELECT 'audit_events' as source, COUNT(*) as total_events, MIN(created_at) as earliest, MAX(created_at) as latest FROM audit_events`
	default:
		http.Error(w, `{"error":"invalid report kind"}`, http.StatusBadRequest)
		return
	}

	rows, err := st.Pool.Query(r.Context(), query)
	if err != nil {
		http.Error(w, `{"error":"report generation failed"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	rowsData := []map[string]any{}
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			continue
		}
		row := map[string]any{}
		for i, col := range rows.FieldDescriptions() {
			row[string(col.Name)] = values[i]
		}
		rowsData = append(rowsData, row)
	}
	if rowsData == nil {
		rowsData = []map[string]any{}
	}

	summary = map[string]any{
		"query": query,
		"rows":  rowsData,
		"count": len(rowsData),
	}

	_, err = st.Pool.Exec(r.Context(),
		`INSERT INTO reports (id, organization_id, name, kind, status, params, summary, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now())`,
		id, "01ARZ3NDEKTSV4RRFFQ69G5FAV", name, kind, "ready", params, summary)
	if err != nil {
		http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id, "name": name, "kind": kind, "status": "ready"})
}

// GenerateSecurityReport generates a security events report.
func GenerateSecurityReport(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		generateReport(st, w, r, "security")
	}
}

// GenerateTrafficReport generates a traffic summary report.
func GenerateTrafficReport(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		generateReport(st, w, r, "traffic")
	}
}

// GenerateIncidentsReport generates an incidents report.
func GenerateIncidentsReport(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		generateReport(st, w, r, "incidents")
	}
}

// GenerateComplianceReport generates a compliance report.
func GenerateComplianceReport(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		generateReport(st, w, r, "compliance")
	}
}

// DownloadReport returns the full report data as a JSON download.
func DownloadReport(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var idVar, name, kind, status, createdBy, created, updated string
		var summary any
		err := st.Pool.QueryRow(r.Context(),
			`SELECT id, name, kind, status, summary, COALESCE(created_by,''), created_at, updated_at FROM reports WHERE id = $1`, id).Scan(&idVar, &name, &kind, &status, &summary, &createdBy, &created, &updated)
		if err != nil {
			http.Error(w, `{"error":"report not found"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", `attachment; filename="report-`+kind+`-`+id+`.json"`)
		json.NewEncoder(w).Encode(map[string]any{
			"id":         idVar,
			"name":       name,
			"kind":       kind,
			"status":     status,
			"summary":    summary,
			"created_by": createdBy,
			"created_at": created,
			"updated_at": updated,
		})
	}
}
