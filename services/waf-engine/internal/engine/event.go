package engine

import "encoding/json"

// SecurityEvent is the security-event extension of the ADR-002 event schema v0.
// Additive-only; never removes base fields.
type SecurityEvent struct {
	SchemaVersion  string              `json:"schema_version"`
	EventID        string              `json:"event_id"`
	RequestID      string              `json:"request_id"`
	Timestamp      string              `json:"timestamp"`
	EventType      string              `json:"event_type"`
	GatewayID      string              `json:"gateway_id"`
	VirtualServerID string             `json:"virtual_server_id,omitempty"`
	ApplicationID  string              `json:"application_id,omitempty"`
	ClientIP       string              `json:"client_ip,omitempty"`
	Method         string              `json:"method,omitempty"`
	Path           string              `json:"path,omitempty"`
	Host           string              `json:"host,omitempty"`
	Status         int                 `json:"status,omitempty"`
	LatencyMs      int64               `json:"latency_ms,omitempty"`
	DecisionAction string              `json:"decision_action,omitempty"`
	Severity       string              `json:"severity,omitempty"`
	Reason         string              `json:"reason,omitempty"`
	RuleIDs        []string            `json:"rule_ids,omitempty"`
	MatchDetails   []map[string]string `json:"match_details,omitempty"`
	Raw            json.RawMessage     `json:"raw,omitempty"`
}

// JSON serializes the event for the JSON-lines sink.
func (e *SecurityEvent) JSON() ([]byte, error) {
	return json.Marshal(e)
}
