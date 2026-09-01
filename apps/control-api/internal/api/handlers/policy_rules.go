package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ariba-shield/control-api/internal/store"
)

// ===== Policy ⇄ Rule Junction =====

// AddPolicyRule attaches a rule to a policy with optional action override.
func AddPolicyRule(st *store.Store) http.HandlerFunc {
	type req struct {
		RuleID         string `json:"rule_id"`
		ActionOverride string `json:"action_override"` // inherit|allow|log|block|challenge
		Enabled        *bool  `json:"enabled"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		policyID := r.PathValue("id")
		var body req
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.RuleID == "" {
			http.Error(w, `{"error":"rule_id required"}`, http.StatusBadRequest)
			return
		}
		enabled := true
		if body.Enabled != nil {
			enabled = *body.Enabled
		}
		id, _ := st.NewID()
		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO policy_rules (id, policy_id, rule_id, action_override, enabled)
			 VALUES ($1,$2,$3,$4,$5)
			 ON CONFLICT (policy_id, rule_id) DO UPDATE SET action_override=EXCLUDED.action_override, enabled=EXCLUDED.enabled`,
			id, policyID, body.RuleID, body.ActionOverride, enabled); err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "policy_id": policyID, "rule_id": body.RuleID})
	}
}

// RemovePolicyRule detaches a rule from a policy.
func RemovePolicyRule(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		policyID := r.PathValue("id")
		ruleID := r.PathValue("ruleId")
		ct, err := st.Pool.Exec(r.Context(),
			`DELETE FROM policy_rules WHERE policy_id = $1 AND rule_id = $2`, policyID, ruleID)
		if err != nil {
			http.Error(w, `{"error":"delete failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"policy-rule not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "detached"})
	}
}

// ListPolicyRules returns rules attached to a policy.
func ListPolicyRules(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		policyID := r.PathValue("id")
		rows, err := st.Pool.Query(r.Context(),
			`SELECT pr.id, pr.rule_id, pr.action_override, pr.enabled, r.name, r.category, r.severity, r.priority, r.status
			 FROM policy_rules pr JOIN rules r ON r.id = pr.rule_id
			 WHERE pr.policy_id = $1 ORDER BY r.priority ASC`, policyID)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		list := []map[string]any{}
		for rows.Next() {
			var prID, ruleID, override, name, cat, sev, status string
			var enabled bool
			var priority int
			if rows.Scan(&prID, &ruleID, &override, &enabled, &name, &cat, &sev, &priority, &status) == nil {
				list = append(list, map[string]any{
					"id": prID, "rule_id": ruleID, "action_override": override, "enabled": enabled,
					"name": name, "category": cat, "severity": sev, "priority": priority, "status": status,
				})
			}
		}
		if list == nil {
			list = []map[string]any{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
	}
}

// ===== Rule Import (JSON + CSV) =====

// ImportRules imports rules from JSON or CSV payload. Detects duplicates by rule_id.
func ImportRules(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		total, imported, updated, failed := 0, 0, 0, 0
		failures := []string{}

		if strings.Contains(contentType, "csv") {
			reader := csv.NewReader(r.Body)
			rows, err := reader.ReadAll()
			if err != nil || len(rows) < 2 {
				http.Error(w, `{"error":"invalid CSV: header + data required"}`, http.StatusBadRequest)
				return
			}
			headers := rows[0]
			for _, row := range rows[1:] {
				if len(row) < 3 {
					continue
				}
				total++
				rec := map[string]string{}
				for i, h := range headers {
					if i < len(row) {
						rec[h] = row[i]
					}
				}
				ruleID := rec["rule_id"]
				if ruleID == "" {
					ruleID = rec["id"]
				}
				if ruleID == "" {
					failed++; failures = append(failures, fmt.Sprintf("row %d: no rule_id", total))
					continue
				}
				name := rec["name"]
				if name == "" {
					name = ruleID
				}
				// Duplicate check
				var exists int
				st.Pool.QueryRow(r.Context(),
					`SELECT COUNT(*) FROM rules WHERE rule_id = $1`, ruleID).Scan(&exists)
				if exists > 0 {
					// Update existing
					_, err := st.Pool.Exec(r.Context(),
						`UPDATE rules SET name=$1, description=$2, category=$3, severity=$4, action=$5, status='active', updated_at=now()
						 WHERE rule_id=$6`,
						name, rec["description"], rec["category"], rec["severity"], rec["action"], ruleID)
					if err != nil {
						failed++; failures = append(failures, fmt.Sprintf("%s: update failed", ruleID))
					} else {
						updated++
					}
				} else {
					// Insert new
					id, err := st.NewID()
					if err != nil {
						failed++; continue
					}
					_, err = st.Pool.Exec(r.Context(),
						`INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, lifecycle_status)
						 VALUES ($1,$2,$3,$4,$5,'imported',$7,$8,0,$10,'active','AND',1,'imported')`,
						id, orgID, ruleID, name, rec["description"], rec["category"], rec["severity"], rec["action"])
					if err != nil {
						failed++; failures = append(failures, fmt.Sprintf("%s: insert failed", ruleID))
					} else {
						imported++
					}
				}
			}
		} else {
			// JSON
			var body struct {
				Rules []map[string]any `json:"rules"`
			}
			if json.NewDecoder(r.Body).Decode(&body); len(body.Rules) == 0 {
				http.Error(w, `{"error":"rules array required"}`, http.StatusBadRequest)
				return
			}
			for _, rec := range body.Rules {
				total++
				ruleID, _ := rec["rule_id"].(string)
				name, _ := rec["name"].(string)
				if ruleID == "" {
					failed++; continue
				}
				var exists int
				st.Pool.QueryRow(r.Context(),
					`SELECT COUNT(*) FROM rules WHERE rule_id = $1`, ruleID).Scan(&exists)
				if exists > 0 {
					_, err := st.Pool.Exec(r.Context(),
						`UPDATE rules SET name=$1, category=$2, severity=$3, action=$4, updated_at=now() WHERE rule_id=$5`,
						name, rec["category"], rec["severity"], rec["action"], ruleID)
					if err == nil {
						updated++
					} else {
						failed++
					}
				} else {
					id, err := st.NewID()
					if err != nil {
						failed++; continue
					}
					cat := fmt.Sprintf("%v", rec["category"])
					sev := fmt.Sprintf("%v", rec["severity"])
					act := fmt.Sprintf("%v", rec["action"])
					_, err = st.Pool.Exec(r.Context(),
						`INSERT INTO rules (id, organization_id, rule_id, name, type, category, severity, action, status, logic, version, lifecycle_status)
						 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,'active','AND',1,'imported')`,
						id, orgID, ruleID, name, "imported", cat, sev, act)
					if err == nil {
						imported++
					} else {
						failed++
						failures = append(failures, fmt.Sprintf("%s: insert failed: %v", ruleID, err))
					}
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"total": total, "imported": imported, "updated": updated,
			"skipped": total - imported - updated - failed, "failed": failed,
			"failures": failures,
		})
	}
}

// ===== Rule Lifecycle =====

// UpdateRuleLifecycle transitions a rule's lifecycle status.
// Allowed transitions: imported→draft→validated→active→deprecated→archived
func UpdateRuleLifecycle(st *store.Store) http.HandlerFunc {
	type req struct {
		Status string `json:"status"`
	}
	valid := map[string]bool{"imported": true, "draft": true, "validated": true, "active": true, "deprecated": true, "archived": true}
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var body req
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || !valid[body.Status] {
			http.Error(w, `{"error":"valid status required: imported|draft|validated|active|deprecated|archived"}`, http.StatusBadRequest)
			return
		}
		ct, err := st.Pool.Exec(r.Context(),
			`UPDATE rules SET lifecycle_status = $1, updated_at = now() WHERE id = $2`, body.Status, id)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"rule not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "lifecycle_status": body.Status})
	}
}

// ===== Policy Activation Validation =====

// ValidatePolicyActivation checks if a policy can be activated.
func ValidatePolicyActivation(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var ruleCount int
		st.Pool.QueryRow(r.Context(),
			`SELECT COUNT(*) FROM policy_rules WHERE policy_id = $1 AND enabled = true`, id).Scan(&ruleCount)

		errors := []string{}
		if ruleCount == 0 {
			errors = append(errors, "At least one enabled security rule is required before activating this policy.")
		}

		valid := len(errors) == 0
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"policy_id":  id,
			"valid":      valid,
			"rule_count": ruleCount,
			"errors":     errors,
		})
	}
}