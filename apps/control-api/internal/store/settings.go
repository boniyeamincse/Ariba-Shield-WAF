package store

import (
	"context"
	"encoding/json"
	"strconv"
)

// Settings is a typed, cached view of system_settings for a category.
type Settings struct {
	values map[string]json.RawMessage
}

// LoadSettings fetches all settings for an org+category into a typed map.
func (s *Store) LoadSettings(ctx context.Context, category string) (*Settings, error) {
	rows, err := s.Pool.Query(ctx,
		`SELECT key, value FROM system_settings
		 WHERE organization_id = $1 AND category = $2`,
		"01ARZ3NDEKTSV4RRFFQ69G5FAV", category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	st := &Settings{values: map[string]json.RawMessage{}}
	for rows.Next() {
		var key string
		var val json.RawMessage
		if err := rows.Scan(&key, &val); err != nil {
			continue
		}
		st.values[key] = val
	}
	return st, nil
}

// LoadAllSettings fetches all settings grouped by category.
func (s *Store) LoadAllSettings(ctx context.Context) (map[string]map[string]json.RawMessage, error) {
	rows, err := s.Pool.Query(ctx,
		`SELECT category, key, value FROM system_settings
		 WHERE organization_id = $1`,
		"01ARZ3NDEKTSV4RRFFQ69G5FAV")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := map[string]map[string]json.RawMessage{}
	for rows.Next() {
		var cat, key string
		var val json.RawMessage
		if err := rows.Scan(&cat, &key, &val); err != nil {
			continue
		}
		if out[cat] == nil {
			out[cat] = map[string]json.RawMessage{}
		}
		out[cat][key] = val
	}
	return out, nil
}

// Bool returns a boolean setting with a fallback.
func (st *Settings) Bool(key string, fallback bool) bool {
	if st == nil {
		return fallback
	}
	v, ok := st.values[key]
	if !ok {
		return fallback
	}
	var b bool
	if err := json.Unmarshal(v, &b); err != nil {
		return fallback
	}
	return b
}

// Int returns an integer setting with a fallback.
func (st *Settings) Int(key string, fallback int) int {
	if st == nil {
		return fallback
	}
	v, ok := st.values[key]
	if !ok {
		return fallback
	}
	var n int
	if err := json.Unmarshal(v, &n); err == nil {
		return n
	}
	// Try parsing a JSON string number.
	var s string
	if err := json.Unmarshal(v, &s); err == nil {
		if n, err := strconv.Atoi(s); err == nil {
			return n
		}
	}
	return fallback
}

// String returns a string setting with a fallback.
func (st *Settings) String(key, fallback string) string {
	if st == nil {
		return fallback
	}
	v, ok := st.values[key]
	if !ok {
		return fallback
	}
	var s string
	if err := json.Unmarshal(v, &s); err != nil {
		return fallback
	}
	return s
}
