package handlers

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/ariba-shield/control-api/internal/store"
)

// Metrics exposes Prometheus-formatted baseline metrics (FR-0.1-034).
func Metrics(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		fmt.Fprintf(w, "# HELP shield_control_up 1 if the control plane is up.\n")
		fmt.Fprintf(w, "# TYPE shield_control_up gauge\n")
		fmt.Fprintf(w, "shield_control_up 1\n")

		fmt.Fprintf(w, "# HELP shield_control_process_start_time_unix Start time of the process.\n")
		fmt.Fprintf(w, "# TYPE shield_control_process_start_time_unix gauge\n")
		fmt.Fprintf(w, "shield_control_process_start_time_unix %d\n", start.Unix())

		fmt.Fprintf(w, "# HELP shield_control_goroutines Current number of goroutines.\n")
		fmt.Fprintf(w, "# TYPE shield_control_goroutines gauge\n")
		fmt.Fprintf(w, "shield_control_goroutines %d\n", runtime.NumGoroutine())

		fmt.Fprintf(w, "# HELP shield_control_go_mem_alloc_bytes Allocated bytes of heap objects.\n")
		fmt.Fprintf(w, "# TYPE shield_control_go_mem_alloc_bytes gauge\n")
		fmt.Fprintf(w, "shield_control_go_mem_alloc_bytes %d\n", m.Alloc)

		fmt.Fprintf(w, "# HELP shield_control_request_duration_ms Histogram of request durations.\n")
		fmt.Fprintf(w, "# TYPE shield_control_request_duration_ms histogram\n")

		// Gateway fleet count
		var gatewayCount int
		_ = st.Pool.QueryRow(r.Context(), `SELECT COUNT(*) FROM gateways`).Scan(&gatewayCount)
		fmt.Fprintf(w, "# HELP shield_gateways_registered Number of registered gateways.\n")
		fmt.Fprintf(w, "# TYPE shield_gateways_registered gauge\n")
		fmt.Fprintf(w, "shield_gateways_registered %d\n", gatewayCount)

		// Application count
		var appCount int
		_ = st.Pool.QueryRow(r.Context(), `SELECT COUNT(*) FROM applications`).Scan(&appCount)
		fmt.Fprintf(w, "# HELP shield_applications_total Number of applications.\n")
		fmt.Fprintf(w, "# TYPE shield_applications_total gauge\n")
		fmt.Fprintf(w, "shield_applications_total %d\n", appCount)
	}
}