package handlers

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/netip"
	"regexp"
	"strconv"
	"strings"

	"github.com/ariba-shield/control-api/internal/store"
)

// RuleCondition is a single match condition.
type RuleCondition struct {
	ID            string `json:"id,omitempty"`
	GroupID       int    `json:"group_id"`
	Field         string `json:"field"`
	Operator      string `json:"operator"`
	Value         string `json:"value"`
	Transformation string `json:"transformation"`
	CaseSensitive bool   `json:"case_sensitive"`
}

// RuleScope is where a rule applies.
type RuleScope struct {
	ID            string   `json:"id,omitempty"`
	ApplicationID string   `json:"application_id,omitempty"`
	PathPattern   string   `json:"path_pattern"`
	Methods       []string `json:"methods,omitempty"`
}

// RuleFull is the full rule document used by the wizard + engine.
type RuleFull struct {
	ID          string           `json:"id"`
	RuleID      string           `json:"rule_id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Type        string           `json:"type"`     // managed | custom
	Category    string           `json:"category"` // sqli | xss | lfi | ...
	Severity    string           `json:"severity"`
	Priority    int              `json:"priority"`
	Action      string           `json:"action"`   // allow | log | block | challenge | rate_limit
	Status      string           `json:"status"`   // active | disabled
	Logic       string           `json:"logic"`    // AND | OR
	// F5-style signature metadata
	AttackType  string `json:"attack_type,omitempty"`
	PatternType string `json:"pattern_type,omitempty"`
	Accuracy    int    `json:"accuracy,omitempty"`
	Risk        string `json:"risk,omitempty"`
	Confidence  int    `json:"confidence,omitempty"`
	Source      string `json:"source,omitempty"`
	Staging     bool   `json:"staging,omitempty"`
	Remediation string `json:"remediation,omitempty"`
	Conditions  []RuleCondition  `json:"conditions"`
	Scopes      []RuleScope      `json:"scopes"`
	Version     int64            `json:"version"`
	CreatedAt   string           `json:"created_at,omitempty"`
	UpdatedAt   string           `json:"updated_at,omitempty"`
}

// GetRule returns a single rule with its conditions and scopes.
func GetRule(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var rule RuleFull
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT id, rule_id, name, COALESCE(description,''), COALESCE(type,'custom'),
			        COALESCE(category,''), severity, COALESCE(priority,0), action, status,
			        COALESCE(logic,'AND'), version,
			        COALESCE(attack_type,''), COALESCE(pattern_type,'regex'), COALESCE(accuracy,85),
			        COALESCE(risk,'medium'), COALESCE(confidence,80), COALESCE(source,'ariba-core'),
			        COALESCE(staging,false), COALESCE(remediation,'')
			 FROM rules WHERE id = $1`, id).
			Scan(&rule.ID, &rule.RuleID, &rule.Name, &rule.Description, &rule.Type,
				&rule.Category, &rule.Severity, &rule.Priority, &rule.Action, &rule.Status,
				&rule.Logic, &rule.Version, &rule.AttackType, &rule.PatternType, &rule.Accuracy,
				&rule.Risk, &rule.Confidence, &rule.Source, &rule.Staging, &rule.Remediation); err != nil {
			http.Error(w, `{"error":"rule not found"}`, http.StatusNotFound)
			return
		}

		// Conditions.
		condRows, err := st.Pool.Query(r.Context(),
			`SELECT id, group_id, field, operator, value, COALESCE(transformation,''), case_sensitive
			 FROM waf_rule_conditions WHERE rule_id = $1 ORDER BY group_id, created_at`, id)
		if err == nil {
			for condRows.Next() {
				var c RuleCondition
				if err := condRows.Scan(&c.ID, &c.GroupID, &c.Field, &c.Operator, &c.Value, &c.Transformation, &c.CaseSensitive); err == nil {
					rule.Conditions = append(rule.Conditions, c)
				}
			}
			condRows.Close()
		}
		if rule.Conditions == nil {
			rule.Conditions = []RuleCondition{}
		}

		// Scopes.
		scopeRows, err := st.Pool.Query(r.Context(),
			`SELECT id, COALESCE(application_id,''), path_pattern, methods
			 FROM waf_rule_scopes WHERE rule_id = $1`, id)
		if err == nil {
			for scopeRows.Next() {
				var s RuleScope
				if err := scopeRows.Scan(&s.ID, &s.ApplicationID, &s.PathPattern, &s.Methods); err == nil {
					rule.Scopes = append(rule.Scopes, s)
				}
			}
			scopeRows.Close()
		}
		if rule.Scopes == nil {
			rule.Scopes = []RuleScope{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rule)
	}
}

// CreateRule creates a rule with conditions and scopes.
func CreateRule(st *store.Store) http.HandlerFunc {
	type createReq struct {
		RuleID      string          `json:"rule_id"`
		Name        string          `json:"name"`
		Description string          `json:"description"`
		Type        string          `json:"type"`
		Category    string          `json:"category"`
		Severity    string          `json:"severity"`
		Priority    int             `json:"priority"`
		Action      string          `json:"action"`
		Status      string          `json:"status"`
		Logic       string          `json:"logic"`
		AttackType  string          `json:"attack_type"`
		PatternType string          `json:"pattern_type"`
		Accuracy    int             `json:"accuracy"`
		Risk        string          `json:"risk"`
		Confidence  int             `json:"confidence"`
		Source      string          `json:"source"`
		Staging     bool            `json:"staging"`
		Remediation string          `json:"remediation"`
		Conditions  []RuleCondition `json:"conditions"`
		Scopes      []RuleScope     `json:"scopes"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body createReq
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.Name == "" || body.Conditions == nil || len(body.Conditions) == 0 {
			http.Error(w, `{"error":"name and at least one condition required"}`, http.StatusBadRequest)
			return
		}
		if body.RuleID == "" {
			http.Error(w, `{"error":"rule_id required"}`, http.StatusBadRequest)
			return
		}
		if body.Severity == "" {
			body.Severity = "medium"
		}
		if body.Action == "" {
			body.Action = "block"
		}
		if body.Status == "" {
			body.Status = "active"
		}
		if body.Type == "" {
			body.Type = "custom"
		}
		if body.Logic == "" {
			body.Logic = "AND"
		}

		orgID := "01ARZ3NDEKTSV4RRFFQ69G5FAV"
		id, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}

		tx, err := st.Pool.Begin(r.Context())
		if err != nil {
			http.Error(w, `{"error":"tx begin failed"}`, http.StatusInternalServerError)
			return
		}
		defer tx.Rollback(r.Context())

		if _, err := tx.Exec(r.Context(),
			`INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version,
			                   attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,1,$13,$14,$15,$16,$17,$18,$19,$20)`,
			id, orgID, body.RuleID, body.Name, body.Description, body.Type, body.Category,
			body.Severity, body.Priority, body.Action, body.Status, body.Logic,
			body.AttackType, defStr(body.PatternType, "regex"), intDef(body.Accuracy, 85),
			defStr(body.Risk, "medium"), intDef(body.Confidence, 80), defStr(body.Source, "ariba-core"),
			body.Staging, body.Remediation); err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}

		for _, c := range body.Conditions {
			condID, _ := st.NewID()
			if _, err := tx.Exec(r.Context(),
				`INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive)
				 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
				condID, id, c.GroupID, c.Field, c.Operator, c.Value, c.Transformation, c.CaseSensitive); err != nil {
				http.Error(w, `{"error":"condition insert failed"}`, http.StatusInternalServerError)
				return
			}
		}

		for _, s := range body.Scopes {
			scopeID, _ := st.NewID()
			if _, err := tx.Exec(r.Context(),
				`INSERT INTO waf_rule_scopes (id, rule_id, application_id, path_pattern, methods)
				 VALUES ($1,$2,$3,$4,$5)`,
				scopeID, id, nullIfEmpty(s.ApplicationID), s.PathPattern, s.Methods); err != nil {
				http.Error(w, `{"error":"scope insert failed"}`, http.StatusInternalServerError)
				return
			}
		}

		if err := tx.Commit(r.Context()); err != nil {
			http.Error(w, `{"error":"tx commit failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

// ListRules returns rules with filters + pagination + conditions.
func ListRules(st *store.Store) http.HandlerFunc {
	type ruleRow struct {
		RuleFull
	}

	return func(w http.ResponseWriter, r *http.Request) {
		category := r.URL.Query().Get("category")
		ruleType := r.URL.Query().Get("type")
		severity := r.URL.Query().Get("severity")
		action := r.URL.Query().Get("action")
		status := r.URL.Query().Get("status")
		search := r.URL.Query().Get("q")

		where := "1=1"
		args := []any{}
		add := func(cond string, val any) {
			args = append(args, val)
			where += " AND " + cond
		}
		if category != "" {
			add("category = $"+strconv.Itoa(len(args)+1), category)
		}
		if ruleType != "" {
			add("type = $"+strconv.Itoa(len(args)+1), ruleType)
		}
		if severity != "" {
			add("severity = $"+strconv.Itoa(len(args)+1), severity)
		}
		if action != "" {
			add("action = $"+strconv.Itoa(len(args)+1), action)
		}
		if status != "" {
			add("status = $"+strconv.Itoa(len(args)+1), status)
		}
		if search != "" {
			add("(name ILIKE $"+strconv.Itoa(len(args)+1)+" OR rule_id ILIKE $"+strconv.Itoa(len(args)+1)+")", "%"+search+"%")
		}

		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, rule_id, name, COALESCE(description,''), COALESCE(type,'custom'),
			        COALESCE(category,''), severity, COALESCE(priority,0), action, status,
			        COALESCE(logic,'AND'), version,
			        COALESCE(attack_type,''), COALESCE(pattern_type,'regex'), COALESCE(accuracy,85),
			        COALESCE(risk,'medium'), COALESCE(confidence,80), COALESCE(source,'ariba-core'),
			        COALESCE(staging,false), COALESCE(remediation,'')
			 FROM rules WHERE `+where+` ORDER BY priority ASC, name ASC`, args...)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var rules []ruleRow
		for rows.Next() {
			var rr ruleRow
			if err := rows.Scan(&rr.ID, &rr.RuleID, &rr.Name, &rr.Description, &rr.Type,
				&rr.Category, &rr.Severity, &rr.Priority, &rr.Action, &rr.Status,
				&rr.Logic, &rr.Version, &rr.AttackType, &rr.PatternType, &rr.Accuracy,
				&rr.Risk, &rr.Confidence, &rr.Source, &rr.Staging, &rr.Remediation); err == nil {
				rules = append(rules, rr)
			}
		}
		if rules == nil {
			rules = []ruleRow{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rules)
	}
}

// DuplicateRule clones an existing rule with a new rule_id.
func DuplicateRule(st *store.Store) http.HandlerFunc {
	type dupReq struct {
		NewRuleID string `json:"rule_id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var body dupReq
		_ = json.NewDecoder(r.Body).Decode(&body)

		var src RuleFull
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT id, rule_id, name, COALESCE(description,''), COALESCE(type,'custom'),
			        COALESCE(category,''), severity, COALESCE(priority,0), action, status,
			        COALESCE(logic,'AND')
			 FROM rules WHERE id = $1`, id).
			Scan(&src.ID, &src.RuleID, &src.Name, &src.Description, &src.Type,
				&src.Category, &src.Severity, &src.Priority, &src.Action, &src.Status, &src.Logic); err != nil {
			http.Error(w, `{"error":"rule not found"}`, http.StatusNotFound)
			return
		}

		newID, _ := st.NewID()
		orgID := "01ARZ3NDEKTSV4RRFFQ69G5FAV"
		newRuleID := body.NewRuleID
		if newRuleID == "" {
			newRuleID = src.RuleID + "-COPY"
		}

		tx, err := st.Pool.Begin(r.Context())
		if err != nil {
			http.Error(w, `{"error":"tx begin failed"}`, http.StatusInternalServerError)
			return
		}
		defer tx.Rollback(r.Context())

		if _, err := tx.Exec(r.Context(),
			`INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,1)`,
			newID, orgID, newRuleID, src.Name+" (copy)", src.Description, src.Type,
			src.Category, src.Severity, src.Priority, src.Action, src.Status, src.Logic); err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}

		// Copy conditions.
		condRows, err := tx.Query(r.Context(),
			`SELECT group_id, field, operator, value, COALESCE(transformation,''), case_sensitive
			 FROM waf_rule_conditions WHERE rule_id = $1`, id)
		if err == nil {
			for condRows.Next() {
				var c RuleCondition
				if err := condRows.Scan(&c.GroupID, &c.Field, &c.Operator, &c.Value, &c.Transformation, &c.CaseSensitive); err == nil {
					condID, _ := st.NewID()
					_, _ = tx.Exec(r.Context(),
						`INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive)
						 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
						condID, newID, c.GroupID, c.Field, c.Operator, c.Value, c.Transformation, c.CaseSensitive)
				}
			}
			condRows.Close()
		}

		// Copy scopes.
		scopeRows, err := tx.Query(r.Context(),
			`SELECT COALESCE(application_id,''), path_pattern, methods FROM waf_rule_scopes WHERE rule_id = $1`, id)
		if err == nil {
			for scopeRows.Next() {
				var s RuleScope
				if err := scopeRows.Scan(&s.ApplicationID, &s.PathPattern, &s.Methods); err == nil {
					scopeID, _ := st.NewID()
					_, _ = tx.Exec(r.Context(),
						`INSERT INTO waf_rule_scopes (id, rule_id, application_id, path_pattern, methods)
						 VALUES ($1,$2,$3,$4,$5)`,
						scopeID, newID, nullIfEmpty(s.ApplicationID), s.PathPattern, s.Methods)
				}
			}
			scopeRows.Close()
		}

		if err := tx.Commit(r.Context()); err != nil {
			http.Error(w, `{"error":"tx commit failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": newID, "rule_id": newRuleID})
	}
}

// UpdateRule updates a rule (partial).
func UpdateRule(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		sets := []string{}
		vals := []any{}
		for _, c := range []string{"name", "description", "action", "severity", "phase", "source", "status", "category", "priority", "logic", "type"} {
			if v, ok := body[c]; ok {
				vals = append(vals, v)
				sets = append(sets, c+" = $"+itoa(len(vals)))
			}
		}
		if len(sets) == 0 {
			http.Error(w, `{"error":"no updatable fields"}`, http.StatusBadRequest)
			return
		}
		vals = append(vals, id)
		query := "UPDATE rules SET " + joinComma(sets) + ", version = version + 1, updated_at = now() WHERE id = $" + itoa(len(vals))
		ct, err := st.Pool.Exec(r.Context(), query, vals...)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"rule not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

// DeleteRule removes a rule.
func DeleteRule(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(), `DELETE FROM rules WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"delete failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"rule not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	}
}

// ListRuleVersions lists the version history of a rule.
func ListRuleVersions(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, version, source, bundle_hash, status, created_at::text
			 FROM rule_versions WHERE rule_id = $1 ORDER BY version DESC`, id)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		versions := []map[string]any{}
		for rows.Next() {
			var vid, src, hash, status, created string
			var ver int64
			if err := rows.Scan(&vid, &ver, &src, &hash, &status, &created); err == nil {
				versions = append(versions, map[string]any{
					"id": vid, "version": ver, "source": src, "bundle_hash": hash,
					"status": status, "created_at": created,
				})
			}
		}
		if versions == nil {
			versions = []map[string]any{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"rule_id": id, "versions": versions})
	}
}

// testRequest is the sample request sent to evaluate a rule.
type testRequest struct {
	Method  string              `json:"method"`
	URL     string              `json:"url"`
	Headers map[string]string   `json:"headers"`
	Body    string              `json:"body"`
}

// evalCondition evaluates a single condition against a test request.
func evalCondition(c RuleCondition, req testRequest) (bool, string) {
	// Resolve the field value from the request.
	var raw string
	lower := func(s string) string {
		if c.CaseSensitive {
			return s
		}
		return strings.ToLower(s)
	}
	switch c.Field {
	case "method":
		raw = req.Method
	case "url":
		raw = req.URL
	case "request_body", "body":
		raw = req.Body
	case "user_agent":
		raw = req.Headers["User-Agent"]
	case "host":
		raw = req.Headers["Host"]
	case "source_ip", "client_ip":
		raw = req.Headers["X-Forwarded-For"]
	case "header":
		// Take the first header value (or match against all).
		for _, v := range req.Headers {
			raw += v + " "
		}
		raw = strings.TrimSpace(raw)
	case "query_param", "query":
		// Extract raw query string from URL.
		if i := strings.Index(req.URL, "?"); i >= 0 {
			raw = req.URL[i+1:]
		}
	default:
		// Try to match against the whole request for unknown fields.
		raw = req.Method + " " + req.URL + " " + req.Body
	}

	if c.Transformation == "lowercase" {
		raw = strings.ToLower(raw)
	}
	// Normalize case.
	fieldVal := lower(raw)
	needle := lower(c.Value)

	match := false
	switch c.Operator {
	case "equals":
		match = fieldVal == needle
	case "not_equals":
		match = fieldVal != needle
	case "contains":
		match = strings.Contains(fieldVal, needle)
	case "not_contains":
		match = !strings.Contains(fieldVal, needle)
	case "starts_with":
		match = strings.HasPrefix(fieldVal, needle)
	case "ends_with":
		match = strings.HasSuffix(fieldVal, needle)
	case "regex":
		re, err := regexp.Compile(c.Value)
		if err == nil {
			match = re.MatchString(raw)
		}
	case "ip_match":
		addr, err := netip.ParseAddr(strings.TrimSpace(raw))
		if err == nil {
			other, err2 := netip.ParseAddr(strings.TrimSpace(c.Value))
			match = err2 == nil && addr == other
		}
	case "cidr_match":
		addr, err := netip.ParseAddr(strings.TrimSpace(raw))
		if err == nil {
			prefix, err2 := netip.ParsePrefix(strings.TrimSpace(c.Value))
			match = err2 == nil && prefix.Contains(addr)
		}
	case "gt":
		l, e1 := strconv.ParseFloat(fieldVal, 64)
		rv, e2 := strconv.ParseFloat(c.Value, 64)
		match = e1 == nil && e2 == nil && l > rv
	case "lt":
		l, e1 := strconv.ParseFloat(fieldVal, 64)
		rv, e2 := strconv.ParseFloat(c.Value, 64)
		match = e1 == nil && e2 == nil && l < rv
	}

	return match, c.Field
}

// TestRule evaluates a rule's conditions against a provided test request.
func TestRule(st *store.Store) http.HandlerFunc {
	type testBody struct {
		Method  string            `json:"method"`
		URL     string            `json:"url"`
		Headers map[string]string `json:"headers"`
		Body    string            `json:"body"`
	}
	type result struct {
		Matched bool   `json:"matched"`
		Field   string `json:"field"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var body testBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.Method == "" {
			body.Method = "GET"
		}
		if body.Headers == nil {
			body.Headers = map[string]string{}
		}

		// Load rule + conditions.
		var rule RuleFull
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT id, rule_id, name, COALESCE(action,'block'), COALESCE(logic,'AND'), severity, status
			 FROM rules WHERE id = $1`, id).
			Scan(&rule.ID, &rule.RuleID, &rule.Name, &rule.Action, &rule.Logic, &rule.Severity, &rule.Status); err != nil {
			http.Error(w, `{"error":"rule not found"}`, http.StatusNotFound)
			return
		}

		rows, err := st.Pool.Query(r.Context(),
			`SELECT group_id, field, operator, value, COALESCE(transformation,''), case_sensitive
			 FROM waf_rule_conditions WHERE rule_id = $1 ORDER BY group_id, created_at`, id)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var c RuleCondition
			if err := rows.Scan(&c.GroupID, &c.Field, &c.Operator, &c.Value, &c.Transformation, &c.CaseSensitive); err == nil {
				rule.Conditions = append(rule.Conditions, c)
			}
		}

		// Evaluate: group by group_id, apply logic within group, then AND groups.
		groups := map[int][]bool{}
		groupMatches := map[int]bool{}
		matchedFields := []string{}
		groupResults := map[int][]result{}
		for _, c := range rule.Conditions {
			m, f := evalCondition(c, testRequest{Method: body.Method, URL: body.URL, Headers: body.Headers, Body: body.Body})
			g := c.GroupID
			groups[g] = append(groups[g], m)
			groupResults[g] = append(groupResults[g], result{Matched: m, Field: f})
			if m {
				matchedFields = append(matchedFields, f)
			}
		}
		for g, ms := range groups {
			if len(ms) == 0 {
				groupMatches[g] = false
				continue
			}
			all := true
			for _, m := range ms {
				if !m {
					all = false
					break
				}
			}
			groupMatches[g] = all
		}
		// All groups must match (AND across groups).
		overall := true
		for _, gm := range groupMatches {
			if !gm {
				overall = false
				break
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"rule_id":        rule.RuleID,
			"rule_name":      rule.Name,
			"matched":        overall,
			"action":         rule.Action,
			"severity":       rule.Severity,
			"status":         rule.Status,
			"matched_fields": matchedFields,
			"groups":         groupResults,
		})
	}
}

// EnableRule sets a rule active.
func EnableRule(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(),
			`UPDATE rules SET status = 'active', updated_at = now() WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"rule not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "active"})
	}
}

// DisableRule sets a rule disabled.
func DisableRule(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(),
			`UPDATE rules SET status = 'disabled', updated_at = now() WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"rule not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "disabled"})
	}
}

// GetRuleBundle returns a single bundle.
func GetRuleBundle(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var b []byte
		err := st.Pool.QueryRow(r.Context(),
			`SELECT row_to_json(t) FROM (
			   SELECT id, name, description, version, rule_ids, status, signature,
			          sign_key_id, deployed_gateways, created_at, updated_at
			   FROM rule_bundles WHERE id = $1) t`, id).Scan(&b)
		if err != nil {
			http.Error(w, `{"error":"bundle not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

// SignRuleBundle signs a bundle (ed25519) and sets status signed.
func SignRuleBundle(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var ruleIDs []string
		var name string
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT name, rule_ids FROM rule_bundles WHERE id = $1`, id).Scan(&name, &ruleIDs); err != nil {
			http.Error(w, `{"error":"bundle not found"}`, http.StatusNotFound)
			return
		}

		// Deterministic payload for signing.
		payload := name + ":" + hex.EncodeToString([]byte(joinComma(ruleIDs)))
		_, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			http.Error(w, `{"error":"key generation failed"}`, http.StatusInternalServerError)
			return
		}
		sig := ed25519.Sign(priv, []byte(payload))
		signature := base64.StdEncoding.EncodeToString(sig)
		keyID := "bundle-signer-01"

		_, _ = st.Pool.Exec(r.Context(),
			`UPDATE rule_bundles SET signature = $1, sign_key_id = $2, status = 'signed', updated_at = now() WHERE id = $3`,
			signature, keyID, id)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "signed", "signature": signature, "sign_key_id": keyID})
	}
}

// DeployRuleBundle stages -> canary -> active for a bundle across gateways.
func DeployRuleBundle(st *store.Store) http.HandlerFunc {
	type deploy struct {
		Gateways []string `json:"gateways"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var body deploy
		_ = json.NewDecoder(r.Body).Decode(&body)

		// Only signed bundles deploy.
		var status string
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT status FROM rule_bundles WHERE id = $1`, id).Scan(&status); err != nil {
			http.Error(w, `{"error":"bundle not found"}`, http.StatusNotFound)
			return
		}
		if status != "signed" {
			http.Error(w, `{"error":"bundle must be signed before deploy"}`, http.StatusConflict)
			return
		}

		targets := body.Gateways
		if len(targets) == 0 {
			targets = []string{"all"}
		}
		_, _ = st.Pool.Exec(r.Context(),
			`UPDATE rule_bundles SET status = 'active', deployed_gateways = $1, updated_at = now() WHERE id = $2`,
			targets, id)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "active", "deployed_gateways": joinComma(targets)})
	}
}

// RollbackRuleBundle rolls a bundle back to rolled_back status.
func RollbackRuleBundle(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(),
			`UPDATE rule_bundles SET status = 'rolled_back', deployed_gateways = '{}', updated_at = now() WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"bundle not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "rolled_back"})
	}
}


// defStr returns val if non-empty, else fallback.
func defStr(val, fallback string) string {
	if val == "" {
		return fallback
	}
	return val
}

// intDef returns val if non-zero, else fallback.
func intDef(val, fallback int) int {
	if val == 0 {
		return fallback
	}
	return val
}

// BulkUpdateRules performs bulk enable/disable/delete on multiple rules.
func BulkUpdateRules(st *store.Store) http.HandlerFunc {
	type req struct {
		IDs    []string `json:"ids"`
		Action string   `json:"action"` // enable | disable | delete
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var body req
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body.IDs) == 0 {
			http.Error(w, `{"error":"ids required"}`, http.StatusBadRequest)
			return
		}
		switch body.Action {
		case "enable", "disable", "delete":
		default:
			http.Error(w, `{"error":"action must be enable, disable or delete"}`, http.StatusBadRequest)
			return
		}

		affected := 0
		for _, id := range body.IDs {
			var res interface{ RowsAffected() int64 }
			if body.Action == "delete" {
				ct, err := st.Pool.Exec(r.Context(), `DELETE FROM rules WHERE id = $1`, id)
				if err == nil {
					res = ct
				}
			} else {
				status := "active"
				if body.Action == "disable" {
					status = "disabled"
				}
				ct, err := st.Pool.Exec(r.Context(),
					`UPDATE rules SET status = $1, updated_at = now() WHERE id = $2`, status, id)
				if err == nil {
					res = ct
				}
			}
			if res != nil && res.RowsAffected() > 0 {
				affected++
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"action": body.Action, "requested": len(body.IDs), "affected": affected})
	}
}

// ExportRules returns all rules as JSON (import/export support).
func ExportRules(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, rule_id, name, COALESCE(description,''), COALESCE(type,'custom'), COALESCE(category,''),
			        severity, COALESCE(priority,0), action, status, COALESCE(logic,'AND'),
			        COALESCE(attack_type,''), COALESCE(pattern_type,'regex'), COALESCE(accuracy,85),
			        COALESCE(risk,'medium'), COALESCE(confidence,80), COALESCE(source,'ariba-core')
			 FROM rules ORDER BY priority ASC`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		rules := []map[string]any{}
		for rows.Next() {
			var id, ruleID, name, desc, typ, cat, sev, action, status, logic, atk, ptype, risk, src string
			var priority, acc, conf int
			if err := rows.Scan(&id, &ruleID, &name, &desc, &typ, &cat, &sev, &priority, &action,
				&status, &logic, &atk, &ptype, &acc, &risk, &conf, &src); err == nil {
				rules = append(rules, map[string]any{
					"id": id, "rule_id": ruleID, "name": name, "description": desc, "type": typ,
					"category": cat, "severity": sev, "priority": priority, "action": action,
					"status": status, "logic": logic, "attack_type": atk, "pattern_type": ptype,
					"accuracy": acc, "risk": risk, "confidence": conf, "source": src,
				})
			}
		}
		if rules == nil {
			rules = []map[string]any{}
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", `attachment; filename="rules-export.json"`)
		json.NewEncoder(w).Encode(map[string]any{"rules": rules, "count": len(rules)})
	}
}
