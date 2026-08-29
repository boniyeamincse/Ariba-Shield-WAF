package compiler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// PolicyDocument is the declarative policy bundle (ADR-002 policy schema v0).
// Only the fields the compiler consumes are declared; unknown fields are
// preserved for additive forward-compatibility.
type PolicyDocument struct {
	SchemaVersion  string           `json:"schema_version"`
	ConfigID       string           `json:"config_id"`
	CreatedAt      string           `json:"created_at"`
	CreatedBy      string           `json:"created_by"`
	Replaces       string           `json:"replaces,omitempty"`
	GatewayTargets []string         `json:"gateway_targets"`
	Settings       Settings         `json:"settings"`
	VirtualServers []VirtualServer  `json:"virtual_servers"`
	BackendPools   []BackendPool    `json:"backend_pools"`
	Headers        *HeadersConfig   `json:"headers,omitempty"`
	Signature      *Signature       `json:"signature,omitempty"`
	Extensions     json.RawMessage  `json:"extensions,omitempty"`
}

type Settings struct {
	LogLevel           string `json:"log_level"`
	EventRetentionDays int    `json:"event_retention_days"`
}

type VirtualServer struct {
	ID                  string      `json:"id"`
	Name                string      `json:"name"`
	ListenAddr          string      `json:"listen_addr"`
	ListenPort          int         `json:"listen_port"`
	TLS                 TLSProfile  `json:"tls"`
	DefaultBackendPoolID string     `json:"default_backend_pool_id"`
	Routes              []Route     `json:"routes"`
	Limits              *Limits     `json:"limits,omitempty"`
}

type TLSProfile struct {
	Enabled        bool     `json:"enabled"`
	CertificateRef string   `json:"certificate_ref,omitempty"`
	MinVersion     string   `json:"min_version,omitempty"`
	Protocols      []string `json:"protocols"`
}

type Route struct {
	ID             string `json:"id"`
	Path           string `json:"path"`
	Match          string `json:"match"`
	BackendPoolID  string `json:"backend_pool_id"`
	Priority       int    `json:"priority,omitempty"`
}

type Limits struct {
	MaxRequestLine int `json:"max_request_line,omitempty"`
	MaxHeaderSize  int `json:"max_header_size,omitempty"`
	MaxBodySize    int `json:"max_body_size,omitempty"`
}

type BackendPool struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	ApplicationID string         `json:"application_id"`
	LBAlgorithm   string         `json:"lb_algorithm"`
	HealthMonitor *HealthMonitor `json:"health_monitor,omitempty"`
	Nodes         []BackendNode  `json:"nodes"`
}

type HealthMonitor struct {
	Type               string `json:"type"`
	IntervalMS         int    `json:"interval_ms"`
	TimeoutMS          int    `json:"timeout_ms"`
	FailThreshold      int    `json:"fail_threshold"`
	PassThreshold      int    `json:"pass_threshold"`
	HTTPPath           string `json:"http_path,omitempty"`
	HTTPExpectedStatus []int  `json:"http_expected_status,omitempty"`
}

type BackendNode struct {
	ID     string `json:"id"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Weight int    `json:"weight,omitempty"`
	Active bool   `json:"active,omitempty"`
}

type HeadersConfig struct {
	TrustedProxyHeaders []string `json:"trusted_proxy_headers"`
	RequestIDHeader     string   `json:"request_id_header"`
}

type Signature struct {
	Algorithm string `json:"algorithm"`
	KeyID     string `json:"key_id"`
	Value     string `json:"value"`
}

// ParsePolicy decodes a policy document from JSON bytes.
func ParsePolicy(data []byte) (*PolicyDocument, error) {
	var doc PolicyDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parse policy document: %w", err)
	}
	return &doc, nil
}

// CanonicalJSON produces the byte-stable serialization used for hashing and
// signing. Encoding a struct marshals fields in declaration order, so
// PolicyDocument is the single canonical ordering.
func (d *PolicyDocument) CanonicalJSON() ([]byte, error) {
	return json.Marshal(d)
}

// BundleHash returns the SHA-256 of the canonical JSON, hex-encoded.
func (d *PolicyDocument) BundleHash() (string, error) {
	canon, err := d.CanonicalJSON()
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(canon)
	return hex.EncodeToString(sum[:]), nil
}