package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/ariba-shield/control-api/internal/store"
)

// GetSecurityPolicy returns a single policy with its lifecycle state.
func GetSecurityPolicy(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var p struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			Description     string `json:"description"`
			EnforcementMode string `json:"enforcement_mode"`
			ApplicationID   string `json:"application_id,omitempty"`
			LifecycleStatus string `json:"lifecycle_status"`
			ActiveVersionID string `json:"active_version_id,omitempty"`
			Version         int64  `json:"version"`
		}
		err := st.Pool.QueryRow(r.Context(),
			`SELECT id, name, COALESCE(description,''), enforcement_mode,
			        COALESCE(application_id,''), lifecycle_status,
			        COALESCE(active_version_id,''), version
			 FROM security_policies WHERE id = $1`, id).
			Scan(&p.ID, &p.Name, &p.Description, &p.EnforcementMode, &p.ApplicationID,
				&p.LifecycleStatus, &p.ActiveVersionID, &p.Version)
		if err != nil {
			http.Error(w, `{"error":"policy not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	}
}

// ValidatePolicy runs dry-run validation on the policy's latest document.
func ValidatePolicy(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		// Get the latest version document.
		var document json.RawMessage
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT document FROM policy_versions WHERE policy_id = $1 ORDER BY version DESC LIMIT 1`, id).Scan(&document); err != nil {
			http.Error(w, `{"error":"no versions to validate"}`, http.StatusNotFound)
			return
		}
		doc, err := parsePolicyDocument(document)
		errors := []string{}
		valid := true
		if err != nil {
			valid = false
			errors = append(errors, err.Error())
		} else {
			if doc.SchemaVersion == "" {
				valid = false
				errors = append(errors, "schema_version is required")
			}
			if len(doc.VirtualServers) == 0 {
				valid = false
				errors = append(errors, "at least one virtual server required")
			}
			if len(doc.BackendPools) == 0 {
				valid = false
				errors = append(errors, "at least one backend pool required")
			}
		}

		// Move lifecycle to VALIDATING, or APPROVAL_REQUIRED if valid.
		next := "validating"
		if valid {
			next = "approval_required"
		}
		_, _ = st.Pool.Exec(r.Context(),
			`UPDATE security_policies SET lifecycle_status = $1, updated_at = now() WHERE id = $2`, next, id)

		w.Header().Set("Content-Type", "application/json")
		if valid {
			json.NewEncoder(w).Encode(map[string]any{"valid": true, "lifecycle_status": next, "errors": []string{}})
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(map[string]any{"valid": false, "lifecycle_status": next, "errors": errors})
		}
	}
}

	// ActivatePolicy activates the latest APPROVED/CANARY version (atomic switch).
func ActivatePolicy(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		// Find the latest version in canary/approved state.
		var versionID string
		err := st.Pool.QueryRow(r.Context(),
			`SELECT id FROM policy_versions
			 WHERE policy_id = $1 AND status IN ('canary','approved')
			 ORDER BY version DESC LIMIT 1`, id).Scan(&versionID)
		if err != nil {
			http.Error(w, `{"error":"no approved/canary version to activate"}`, http.StatusNotFound)
			return
		}

		// Demote current active.
		_, _ = st.Pool.Exec(r.Context(),
			`UPDATE policy_versions SET status = 'superseded' WHERE policy_id = $1 AND status = 'active'`, id)
		// Activate the target.
		_, _ = st.Pool.Exec(r.Context(),
			`UPDATE policy_versions SET status = 'active' WHERE id = $1`, versionID)
		// Update policy lifecycle.
		_, _ = st.Pool.Exec(r.Context(),
			`UPDATE security_policies SET lifecycle_status = 'active', active_version_id = $1, version = version + 1, updated_at = now() WHERE id = $2`,
			versionID, id)

		// MVP PUSH AUTOMATION
		// Read document to generate WAF rules locally (since pull architecture is Phase 2)
		var doc []byte
		if err := st.Pool.QueryRow(r.Context(), `SELECT document FROM policy_versions WHERE id = $1`, versionID).Scan(&doc); err == nil {
			go pushToWAFEngine(doc)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "active_version_id": versionID, "lifecycle_status": "active"})
	}
}

func pushToWAFEngine(doc []byte) {
	var p struct {
		ConfigID string  `json:"config_id"`
		WAF      *wafDoc `json:"waf,omitempty"`
	}
	if err := json.Unmarshal(doc, &p); err != nil {
		return
	}

	var buf bytes.Buffer
	if p.WAF == nil || !p.WAF.Enabled {
		buf.WriteString("# Ariba Shield WAF — Minimal Baseline Configuration\n")
		buf.WriteString("Include @coraza.conf-recommended\n")
		buf.WriteString("SecRuleEngine DetectionOnly\n")
		buf.WriteString("SecRequestBodyAccess On\n")
		buf.WriteString("SecRequestBodyLimit 13107200\n")
		buf.WriteString("SecRequestBodyNoFilesLimit 131072\n")
		buf.WriteString("SecResponseBodyAccess Off\n")
		buf.WriteString("SecPcreMatchLimit 100000\n")
		buf.WriteString("SecPcreMatchLimitRecursion 100000\n\n")
		buf.WriteString("Include @crs-setup.conf.example\n")
		buf.WriteString("SecAction \"id:900000,phase:1,nolog,pass,t:none,setvar:tx.paranoia_level=1\"\n")
		buf.WriteString("Include @owasp_crs/*.conf\n")
		buf.WriteString("SecAction \"id:900001,phase:1,nolog,pass,t:none,setvar:tx.blocking_anomaly_score=5\"\n")
	} else {
		buf.WriteString("# Ariba Shield WAF — Dynamically Generated Configuration\n")
		buf.WriteString(fmt.Sprintf("# config_id=%s\n\n", p.ConfigID))
		buf.WriteString("Include @coraza.conf-recommended\n\n")
		buf.WriteString("SecRuleEngine On\n")
		buf.WriteString("SecRequestBodyAccess On\n")
		buf.WriteString("SecRequestBodyLimit 13107200\n")
		buf.WriteString("SecRequestBodyNoFilesLimit 131072\n")
		buf.WriteString("SecResponseBodyAccess Off\n")
		buf.WriteString("SecPcreMatchLimit 100000\n")
		buf.WriteString("SecPcreMatchLimitRecursion 100000\n\n")

		buf.WriteString("# --- CRS Setup ---\n")
		buf.WriteString("Include @crs-setup.conf.example\n\n")

		pl := p.WAF.ParanoiaLevel
		if pl < 1 || pl > 4 {
			pl = 1
		}
		buf.WriteString(fmt.Sprintf("SecAction \"id:900000,phase:1,nolog,pass,t:none,setvar:tx.paranoia_level=%d\"\n\n", pl))
		buf.WriteString("# --- OWASP CRS Managed Rules ---\n")
		buf.WriteString("Include @owasp_crs/*.conf\n\n")

		thresh := p.WAF.AnomalyThreshold
		if thresh <= 0 {
			thresh = 5
		}
		buf.WriteString("# --- Anomaly Score Threshold ---\n")
		buf.WriteString(fmt.Sprintf("SecAction \"id:900001,phase:1,nolog,pass,t:none,setvar:tx.blocking_anomaly_score=%d\"\n", thresh))
	}

	// Write to shared volume
	_ = os.WriteFile("/rules/crs.conf", buf.Bytes(), 0644)

	// Trigger WAF engine atomic reload
	resp, err := http.Post("http://waf-engine:8082/__reload", "application/json", nil)
	if err == nil {
		resp.Body.Close()
	}
}

// DisablePolicy sets the policy lifecycle to a disabled state (stops enforcement).
func DisablePolicy(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(),
			`UPDATE security_policies SET lifecycle_status = 'disabled', updated_at = now() WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"policy not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "lifecycle_status": "disabled"})
	}
}

// RollbackPolicy rolls the policy back to the last superseded (previously active) version.
func RollbackPolicy(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		// Mark current active rolled_back.
		_, _ = st.Pool.Exec(r.Context(),
			`UPDATE policy_versions SET status = 'rolled_back' WHERE policy_id = $1 AND status = 'active'`, id)
		// Find the newest superseded version (the previous active).
		var versionID string
		err := st.Pool.QueryRow(r.Context(),
			`SELECT id FROM policy_versions
			 WHERE policy_id = $1 AND status = 'superseded'
			 ORDER BY version DESC LIMIT 1`, id).Scan(&versionID)
		if err != nil {
			http.Error(w, `{"error":"no previous active version to roll back to"}`, http.StatusNotFound)
			return
		}
		_, _ = st.Pool.Exec(r.Context(),
			`UPDATE policy_versions SET status = 'active' WHERE id = $1`, versionID)
		_, _ = st.Pool.Exec(r.Context(),
			`UPDATE security_policies SET lifecycle_status = 'active', active_version_id = $1, version = version + 1, updated_at = now() WHERE id = $2`,
			versionID, id)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "active_version_id": versionID, "lifecycle_status": "active", "rolled_back": "true"})
	}
}

// ListPolicyVersions lists the version history of a policy.
func ListPolicyVersions(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, version, bundle_hash, status, created_at::text
			 FROM policy_versions WHERE policy_id = $1 ORDER BY version DESC`, id)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		versions := []map[string]any{}
		for rows.Next() {
			var vid, hash, status, created string
			var ver int64
			if err := rows.Scan(&vid, &ver, &hash, &status, &created); err == nil {
				versions = append(versions, map[string]any{
					"id": vid, "version": ver, "bundle_hash": hash, "status": status, "created_at": created,
				})
			}
		}
		if versions == nil {
			versions = []map[string]any{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"policy_id": id, "versions": versions})
	}
}

// DiffPolicy diffs two versions of a policy (from/to query).
func DiffPolicy(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")
		if from == "" || to == "" {
			http.Error(w, `{"error":"from and to version ids required"}`, http.StatusBadRequest)
			return
		}

		var fromDoc, toDoc map[string]any
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT document FROM policy_versions WHERE id = $1 AND policy_id = $2`, from, id).Scan(&fromDoc); err != nil {
			http.Error(w, `{"error":"from version not found"}`, http.StatusNotFound)
			return
		}
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT document FROM policy_versions WHERE id = $1 AND policy_id = $2`, to, id).Scan(&toDoc); err != nil {
			http.Error(w, `{"error":"to version not found"}`, http.StatusNotFound)
			return
		}

		fj, _ := json.Marshal(fromDoc)
		tj, _ := json.Marshal(toDoc)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(simpleDiff(fj, tj))
	}
}

// ClonePolicy creates a new policy from an existing one (copies name/schema, new id).
func ClonePolicy(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		srcID := r.PathValue("id")
		var name, desc, mode, appID string
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT name, COALESCE(description,''), enforcement_mode, COALESCE(application_id,'')
			 FROM security_policies WHERE id = $1`, srcID).
			Scan(&name, &desc, &mode, &appID); err != nil {
			http.Error(w, `{"error":"source policy not found"}`, http.StatusNotFound)
			return
		}

		newID, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}
		userID := "01ARZ3NDEKTSV4RRFFQ69G5FAW"
		orgID := "01ARZ3NDEKTSV4RRFFQ69G5FAV"

		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO security_policies
			   (id, organization_id, application_id, name, description, enforcement_mode, created_by, lifecycle_status)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, 'draft')`,
			newID, orgID, nullIfEmpty(appID), name+" (clone)", desc, mode, userID); err != nil {
			http.Error(w, `{"error":"clone failed"}`, http.StatusInternalServerError)
			return
		}

		// Copy the latest version document as the first draft version.
		var latestDoc json.RawMessage
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT document FROM policy_versions WHERE policy_id = $1 ORDER BY version DESC LIMIT 1`, srcID).Scan(&latestDoc); err == nil {
			versionID, _ := st.NewID()
			_, _ = st.Pool.Exec(r.Context(),
				`INSERT INTO policy_versions (id, policy_id, version, document, bundle_hash, status, created_by)
				 VALUES ($1, $2, 1, $3, '', 'draft', $4)`,
				versionID, newID, latestDoc, userID)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": newID, "name": name + " (clone)", "lifecycle_status": "draft"})
	}
}

