package compiler

import (
	"strings"
	"testing"
)

func TestParsePolicy(t *testing.T) {
	data := `{
		"schema_version": "0.1",
		"config_id": "01JARZ3NDEKTSV4RRFFQ69G5FAV",
		"created_at": "2026-08-28T10:00:00Z",
		"created_by": "01JARZ3NDEKTSV4RRFFQ69G5FAW",
		"settings": {"log_level": "info", "event_retention_days": 30},
		"virtual_servers": [{
			"id": "01JARZ3NDEKTSV4RRFFQ69G5FAX",
			"name": "test-vs",
			"listen_addr": "0.0.0.0",
			"listen_port": 443,
			"tls": {
				"enabled": true,
				"certificate_ref": "cert-01",
				"min_version": "1.2",
				"protocols": ["http/1.1", "h2"]
			},
			"default_backend_pool_id": "01JARZ3NDEKTSV4RRFFQ69G5FAY"
		}],
		"backend_pools": [{
			"id": "01JARZ3NDEKTSV4RRFFQ69G5FAY",
			"name": "api-pool",
			"application_id": "01JARZ3NDEKTSV4RRFFQ69G5FAZ",
			"lb_algorithm": "round_robin",
			"nodes": [{
				"id": "01JARZ3NDEKTSV4RRFFQ69G5FBA",
				"host": "10.0.0.11",
				"port": 8080,
				"weight": 1,
				"active": true
			}]
		}]
	}`

	doc, err := ParsePolicy([]byte(data))
	if err != nil {
		t.Fatalf("ParsePolicy: %v", err)
	}
	if doc.ConfigID != "01JARZ3NDEKTSV4RRFFQ69G5FAV" {
		t.Fatalf("expected config_id, got %s", doc.ConfigID)
	}
	if len(doc.VirtualServers) != 1 {
		t.Fatalf("expected 1 virtual server, got %d", len(doc.VirtualServers))
	}
}

func TestBundleHash(t *testing.T) {
	data := `{"schema_version":"0.1","config_id":"test-01","created_at":"2026-08-28T10:00:00Z","created_by":"user-01","settings":{"log_level":"info","event_retention_days":30},"virtual_servers":[{"id":"vs-01","name":"test-vs","listen_addr":"0.0.0.0","listen_port":443,"tls":{"enabled":true,"certificate_ref":"cert-01","min_version":"1.2","protocols":["http/1.1","h2"]},"default_backend_pool_id":"pool-01"}],"backend_pools":[{"id":"pool-01","name":"api-pool","application_id":"app-01","lb_algorithm":"round_robin","nodes":[{"id":"node-01","host":"10.0.0.11","port":8080,"weight":1,"active":true}]}]}`

	doc, err := ParsePolicy([]byte(data))
	if err != nil {
		t.Fatalf("ParsePolicy: %v", err)
	}

	h1, err := doc.BundleHash()
	if err != nil {
		t.Fatalf("BundleHash: %v", err)
	}
	h2, err := doc.BundleHash()
	if err != nil {
		t.Fatalf("BundleHash: %v", err)
	}
	if h1 != h2 {
		t.Fatalf("deterministic hash failed: %s != %s", h1, h2)
	}
	if len(h1) != 64 {
		t.Fatalf("expected 64 hex chars, got %d in %s", len(h1), h1)
	}
}

func TestRenderNginxConfig(t *testing.T) {
	doc := &PolicyDocument{
		SchemaVersion:  "0.1",
		ConfigID:       "test-01",
		CreatedAt:      "2026-08-28T10:00:00Z",
		CreatedBy:      "user-01",
		GatewayTargets: []string{"gw-01"},
		Settings:       Settings{LogLevel: "info", EventRetentionDays: 30},
		VirtualServers: []VirtualServer{
			{
				ID:                  "vs-01",
				Name:                "app-prod",
				ListenAddr:          "0.0.0.0",
				ListenPort:          443,
				TLS:                 TLSProfile{Enabled: true, CertificateRef: "cert-01", MinVersion: "1.2", Protocols: []string{"http/1.1", "h2"}},
				DefaultBackendPoolID: "pool-01",
			},
		},
		BackendPools: []BackendPool{
			{
				ID: "pool-01", Name: "api-pool", ApplicationID: "app-01", LBAlgorithm: "round_robin",
				Nodes: []BackendNode{{ID: "node-01", Host: "10.0.0.11", Port: 8080, Weight: 1, Active: true}},
			},
		},
	}

	config, err := RenderNginxConfig(doc, map[string]string{"cert-01": "/etc/ssl/shield"})
	if err != nil {
		t.Fatalf("RenderNginxConfig: %v", err)
	}

	if !strings.Contains(config, "listen 0.0.0.0:443 ssl") {
		t.Fatalf("expected 'listen 0.0.0.0:443 ssl' in generated config:\n%s", config)
	}
	if !strings.Contains(config, "server 10.0.0.11:8080") {
		t.Fatalf("expected backend server in generated config:\n%s", config)
	}
	if !strings.Contains(config, "upstream pool_api_pool") {
		t.Fatalf("expected upstream pool_api_pool in generated config:\n%s", config)
	}
	if !strings.Contains(config, "ssl_certificate /etc/ssl/shield/cert.pem") {
		t.Fatalf("expected ssl certificate path in generated config:\n%s", config)
	}
}

func TestRenderNoActiveNodes(t *testing.T) {
	doc := &PolicyDocument{
		SchemaVersion:  "0.1",
		ConfigID:       "test-02",
		CreatedAt:      "2026-08-28T10:00:00Z",
		CreatedBy:      "user-01",
		GatewayTargets: []string{"gw-01"},
		Settings:       Settings{LogLevel: "info", EventRetentionDays: 30},
		VirtualServers: []VirtualServer{
			{ID: "vs-01", Name: "test-vs", ListenAddr: "0.0.0.0", ListenPort: 80, TLS: TLSProfile{Enabled: false}, DefaultBackendPoolID: "pool-01"},
		},
		BackendPools: []BackendPool{
			{ID: "pool-01", Name: "empty-pool", ApplicationID: "app-01", LBAlgorithm: "round_robin", Nodes: []BackendNode{{ID: "node-01", Host: "10.0.0.11", Port: 8080, Active: false}}},
		},
	}

	_, err := RenderNginxConfig(doc, map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty pool, got nil")
	}
	if !strings.Contains(err.Error(), "no active nodes") {
		t.Fatalf("expected 'no active nodes' error, got: %v", err)
	}
}