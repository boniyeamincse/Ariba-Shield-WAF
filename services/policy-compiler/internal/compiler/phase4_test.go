package compiler

import (
	"strings"
	"testing"
)

// buildPoolDoc returns a policy doc with a single pool and VS for compiler tests.
func buildPoolDoc(lb string, mutate func(pool *BackendPool)) *PolicyDocument {
	pool := BackendPool{
		ID: "pool-01", Name: "api-pool", ApplicationID: "app-01",
		LBAlgorithm: lb,
		Nodes:       []BackendNode{{ID: "n-01", Host: "10.0.0.1", Port: 8080, Weight: 1, Active: true}},
	}
	if mutate != nil {
		mutate(&pool)
	}
	return &PolicyDocument{
		SchemaVersion:  "0.1",
		ConfigID:       "test-lb",
		CreatedAt:      "2026-08-29T00:00:00Z",
		CreatedBy:      "user-01",
		GatewayTargets: []string{"gw-01"},
		Settings:       Settings{LogLevel: "info", EventRetentionDays: 30},
		VirtualServers: []VirtualServer{
			{ID: "vs-01", Name: "app", ListenAddr: "0.0.0.0", ListenPort: 443, TLS: TLSProfile{Enabled: false}, DefaultBackendPoolID: "pool-01"},
		},
		BackendPools: []BackendPool{pool},
	}
}

func TestRenderIPHash(t *testing.T) {
	doc := buildPoolDoc("ip_hash", nil)
	cfg, err := RenderNginxConfig(doc, map[string]string{})
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if !strings.Contains(cfg, "ip_hash;") {
		t.Fatalf("expected ip_hash in config:\n%s", cfg)
	}
}

func TestRenderConsistentHash(t *testing.T) {
	doc := buildPoolDoc("consistent_hash", nil)
	cfg, err := RenderNginxConfig(doc, map[string]string{})
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if !strings.Contains(cfg, "hash $request_uri consistent;") {
		t.Fatalf("expected consistent hash in config:\n%s", cfg)
	}
}

func TestRenderTLSReencryption(t *testing.T) {
	doc := buildPoolDoc("round_robin", func(p *BackendPool) {
		p.Nodes[0].Protocol = "https"
	})
	cfg, err := RenderNginxConfig(doc, map[string]string{})
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if !strings.Contains(cfg, "proxy_pass https://pool_api_pool;") {
		t.Fatalf("expected https proxy_pass:\n%s", cfg)
	}
	if !strings.Contains(cfg, "proxy_ssl_server_name on;") {
		t.Fatalf("expected proxy_ssl_server_name on:\n%s", cfg)
	}
}

func TestRenderDraining(t *testing.T) {
	doc := buildPoolDoc("round_robin", func(p *BackendPool) {
		p.Nodes = append(p.Nodes,
			BackendNode{ID: "n-02", Host: "10.0.0.2", Port: 8080, Weight: 1, Active: true, Drain: true},
			BackendNode{ID: "n-03", Host: "10.0.0.3", Port: 8080, Weight: 1, Active: false},
		)
	})
	cfg, err := RenderNginxConfig(doc, map[string]string{})
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if strings.Contains(cfg, "10.0.0.2:8080") {
		t.Fatalf("draining node must be excluded from upstream:\n%s", cfg)
	}
	if strings.Contains(cfg, "10.0.0.3:8080") {
		t.Fatalf("inactive node must be excluded from upstream:\n%s", cfg)
	}
	if !strings.Contains(cfg, "10.0.0.1:8080") {
		t.Fatalf("active node must be present:\n%s", cfg)
	}
}

func TestRenderExactRouteMatch(t *testing.T) {
	doc := buildPoolDoc("round_robin", nil)
	doc.VirtualServers[0].Routes = []Route{
		{ID: "r-01", Path: "/api/status", Match: "exact", BackendPoolID: "pool-01"},
	}
	cfg, err := RenderNginxConfig(doc, map[string]string{})
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if !strings.Contains(cfg, "location = /api/status {") {
		t.Fatalf("expected exact location match:\n%s", cfg)
	}
}

func TestRenderServerNameCatchAll(t *testing.T) {
	doc := buildPoolDoc("round_robin", nil)
	cfg, err := RenderNginxConfig(doc, map[string]string{})
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if strings.Contains(cfg, ".shield.local") {
		t.Fatalf("hardcoded server_name must not appear:\n%s", cfg)
	}
	if !strings.Contains(cfg, "server_name _;") {
		t.Fatalf("expected catch-all server_name _:\n%s", cfg)
	}
}