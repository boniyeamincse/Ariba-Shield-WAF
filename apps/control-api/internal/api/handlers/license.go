package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ariba-shield/control-api/internal/store"
)

// licenseResponse is the response shape for license info.
type licenseResponse struct {
	ID              string `json:"id"`
	LicenseKey      string `json:"license_key"`
	Product         string `json:"product"`
	Edition         string `json:"edition"`
	Status          string `json:"status"`
	Seats           int    `json:"seats"`
	MaxGateways     int    `json:"max_gateways"`
	MaxApplications int    `json:"max_applications"`
	IssuedAt        string `json:"issued_at"`
	ExpiresAt       string `json:"expires_at"`
	ActivatedAt     string `json:"activated_at"`
}

// GetLicense returns the current active license.
func GetLicense(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var id, key, product, edition, status, issued, expires, activated string
		var seats, maxGateways, maxApps int
		err := st.Pool.QueryRow(r.Context(),
			`SELECT id, license_key, product, edition, status, seats, max_gateways, max_applications, issued_at, expires_at, activated_at
			 FROM licenses WHERE organization_id = $1 ORDER BY created_at DESC LIMIT 1`,
			"01ARZ3NDEKTSV4RRFFQ69G5FAV").Scan(&id, &key, &product, &edition, &status, &seats, &maxGateways, &maxApps, &issued, &expires, &activated)
		if err != nil {
			// No license found — return a default community license.
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"product":          "ariba-shield-waf",
				"edition":          "community",
				"status":           "active",
				"seats":            1,
				"max_gateways":     1,
				"max_applications": 10,
				"message":          "community edition — no license key required",
			})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(licenseResponse{
			ID: id, LicenseKey: key, Product: product, Edition: edition,
			Status: status, Seats: seats, MaxGateways: maxGateways,
			MaxApplications: maxApps, IssuedAt: issued, ExpiresAt: expires, ActivatedAt: activated,
		})
	}
}

// ActivateLicense activates a license key.
func ActivateLicense(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type activateReq struct {
			LicenseKey string `json:"license_key"`
		}
		var body activateReq
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.LicenseKey == "" {
			http.Error(w, `{"error":"license_key is required"}`, http.StatusBadRequest)
			return
		}
		// In a full implementation, this would validate the key against a licensing server.
		// For now, we accept any key and mark it active.
		id, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}
		now := time.Now()
		expires := now.AddDate(1, 0, 0) // 1 year validity
		_, err = st.Pool.Exec(r.Context(),
			`INSERT INTO licenses (id, organization_id, license_key, product, edition, status, seats, max_gateways, max_applications, issued_at, expires_at, activated_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			id, "01ARZ3NDEKTSV4RRFFQ69G5FAV", body.LicenseKey, "ariba-shield-waf", "pro", "active", 10, 5, 50, now, expires, now)
		if err != nil {
			http.Error(w, `{"error":"activation failed"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":               id,
			"license_key":      body.LicenseKey,
			"product":          "ariba-shield-waf",
			"edition":          "pro",
			"status":           "active",
			"seats":            10,
			"max_gateways":     5,
			"max_applications": 50,
			"expires_at":       expires.Format(time.RFC3339),
			"message":          "license activated — validation simulated",
		})
	}
}

// DeactivateLicense deactivates the current license.
func DeactivateLicense(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := st.Pool.Exec(r.Context(),
			`UPDATE licenses SET status = 'inactive', updated_at = now()
			 WHERE organization_id = $1 AND status = 'active'`,
			"01ARZ3NDEKTSV4RRFFQ69G5FAV")
		if err != nil {
			http.Error(w, `{"error":"deactivation failed"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "inactive",
			"message": "license deactivated",
		})
	}
}

// GetLicenseUsage returns usage stats against license limits.
func GetLicenseUsage(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var activeGateways, activeApps int
		_ = st.Pool.QueryRow(r.Context(), `SELECT COUNT(*) FROM gateways WHERE status IN ('active','starting')`).Scan(&activeGateways)
		_ = st.Pool.QueryRow(r.Context(), `SELECT COUNT(*) FROM applications WHERE status = 'active'`).Scan(&activeApps)

		var maxGateways, maxApps int
		_ = st.Pool.QueryRow(r.Context(),
			`SELECT COALESCE(max_gateways, 1), COALESCE(max_applications, 10)
			 FROM licenses WHERE organization_id = $1 AND status = 'active' ORDER BY created_at DESC LIMIT 1`,
			"01ARZ3NDEKTSV4RRFFQ69G5FAV").Scan(&maxGateways, &maxApps)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"gateways": map[string]any{
				"current": activeGateways,
				"limit":   maxGateways,
				"percent": safePercent(activeGateways, maxGateways),
			},
			"applications": map[string]any{
				"current": activeApps,
				"limit":   maxApps,
				"percent": safePercent(activeApps, maxApps),
			},
		})
	}
}

// GetLicenseEntitlements returns the feature set granted by the license.
func GetLicenseEntitlements(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var edition string
		_ = st.Pool.QueryRow(r.Context(),
			`SELECT COALESCE(edition, 'community') FROM licenses
			 WHERE organization_id = $1 AND status = 'active' ORDER BY created_at DESC LIMIT 1`,
			"01ARZ3NDEKTSV4RRFFQ69G5FAV").Scan(&edition)

		entitlements := map[string]bool{
			"core_proxy":         true,
			"basic_waf":          true,
			"rate_limiting":      edition == "pro" || edition == "enterprise",
			"bot_management":     edition == "pro" || edition == "enterprise",
			"api_security":       edition == "enterprise",
			"dlp":                edition == "enterprise",
			"ml_learning":        edition == "enterprise",
			"advanced_analytics": edition == "enterprise",
			"multi_tenant":       edition == "enterprise",
			"sso_saml":           edition == "pro" || edition == "enterprise",
			"audit_trail":        true,
			"rbac":               true,
			"support":            edition != "community",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"edition":      edition,
			"entitlements": entitlements,
		})
	}
}

// safePercent computes (current/limit)*100, returning 0 when limit is 0.
func safePercent(current, limit int) float64 {
	if limit <= 0 {
		return 0
	}
	return float64(current) / float64(limit) * 100
}
