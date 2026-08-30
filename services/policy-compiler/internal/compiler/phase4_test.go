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

func TestRenderServerNameFromPolicy(t *testing.T) {
	doc := buildPoolDoc("round_robin", nil)
	cfg, err := RenderNginxConfig(doc, map[string]string{})
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if strings.Contains(cfg, ".shield.local") {
		t.Fatalf("hardcoded server_name must not appear:\n%s", cfg)
	}
	if !strings.Contains(cfg, "server_name app;") {
		t.Fatalf("expected server_name from policy (vs.Name), got:\n%s", cfg)
	}
}

func TestRenderServerNameFallbackCatchAll(t *testing.T) {
	doc := buildPoolDoc("round_robin", nil)
	doc.VirtualServers[0].Name = ""
	cfg, err := RenderNginxConfig(doc, map[string]string{})
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if !strings.Contains(cfg, "server_name _;") {
		t.Fatalf("expected catch-all server_name _ for empty name:\n%s", cfg)
	}
}

func TestRenderTLSProtocols(t *testing.T) {
	cases := []struct {
		minVersion string
		want       string
	}{
		{"1.2", "ssl_protocols TLSv1.2 TLSv1.3;"},
		{"1.3", "ssl_protocols TLSv1.3;"},
		{"1.1", "ssl_protocols TLSv1.1 TLSv1.2 TLSv1.3;"},
		{"", "ssl_protocols TLSv1.2 TLSv1.3;"},
	}
	for _, c := range cases {
		doc := buildPoolDoc("round_robin", nil)
		doc.VirtualServers[0].TLS = TLSProfile{Enabled: true, CertificateRef: "cert-01", MinVersion: c.minVersion}
		cfg, err := RenderNginxConfig(doc, map[string]string{"cert-01": "/etc/ssl/shield"})
		if err != nil {
			t.Fatalf("render(min=%q): %v", c.minVersion, err)
		}
		if !strings.Contains(cfg, c.want) {
			t.Fatalf("min_version %q: expected %q in config:\n%s", c.minVersion, c.want, cfg)
		}
	}
}

func TestRenderRoutePriorityOrder(t *testing.T) {
	doc := buildPoolDoc("round_robin", nil)
	doc.VirtualServers[0].Routes = []Route{
		{ID: "r-low", Path: "/api", Match: "prefix", BackendPoolID: "pool-01", Priority: 1},
		{ID: "r-high", Path: "/api/status", Match: "prefix", BackendPoolID: "pool-01", Priority: 100},
	}
	cfg, err := RenderNginxConfig(doc, map[string]string{})
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	low := strings.Index(cfg, "location /api/ {")
	high := strings.Index(cfg, "location /api/status/ {")
	if low == -1 || high == -1 {
		t.Fatalf("expected both locations in config:\n%s", cfg)
	}
	if low < high {
		t.Fatalf("higher-priority route must be declared before lower-priority route:\n%s", cfg)
	}
}

func TestRenderRoutePriorityConflict(t *testing.T) {
	doc := buildPoolDoc("round_robin", nil)
	doc.VirtualServers[0].Routes = []Route{
		{ID: "r-01", Path: "/api", Match: "prefix", BackendPoolID: "pool-01"},
		{ID: "r-02", Path: "/api/", Match: "prefix", BackendPoolID: "pool-01"},
	}
	_, err := RenderNginxConfig(doc, map[string]string{})
	if err == nil {
		t.Fatal("expected conflict error for duplicate paths, got nil")
	}
	if !strings.Contains(err.Error(), "conflicts with route") {
		t.Fatalf("expected conflict error message, got: %v", err)
	}
}

func TestRenderAccessLogRealIDs(t *testing.T) {
	doc := buildPoolDoc("round_robin", nil)
	cfg, err := RenderNginxConfig(doc, map[string]string{})
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if strings.Contains(cfg, `"virtual_server_id":"vs"`) || strings.Contains(cfg, `"application_id":"app"`) {
		t.Fatalf("hardcoded vs/app placeholders must not appear:\n%s", cfg)
	}
	if !strings.Contains(cfg, `"virtual_server_id":"$shield_vs_id"`) {
		t.Fatalf("log_format must use $shield_vs_id variable:\n%s", cfg)
	}
	if !strings.Contains(cfg, `"application_id":"$shield_app_id"`) {
		t.Fatalf("log_format must use $shield_app_id variable:\n%s", cfg)
	}
	if !strings.Contains(cfg, `set $shield_vs_id "vs-01";`) {
		t.Fatalf("server block must set $shield_vs_id from policy:\n%s", cfg)
	}
	if !strings.Contains(cfg, `set $shield_app_id "app-01";`) {
		t.Fatalf("server block must set $shield_app_id from policy:\n%s", cfg)
	}
}

func TestRenderHealthCheck(t *testing.T) {
	doc := buildPoolDoc("round_robin", func(p *BackendPool) {
		p.HealthMonitor = &HealthMonitor{
			Type: "http", IntervalMS: 5000, TimeoutMS: 2000,
			FailThreshold: 3, PassThreshold: 2,
			HTTPPath: "/healthz", HTTPExpectedStatus: []int{200, 204},
		}
	})
	cfg, err := RenderNginxConfig(doc, map[string]string{})
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if !strings.Contains(cfg, "health_check interval=5000ms timeout=2000ms fails=3 passes=2 uri=/healthz match=shield_health_200_204;") {
		t.Fatalf("expected health_check directive:\n%s", cfg)
	}
	if !strings.Contains(cfg, "match shield_health_200_204 {") {
		t.Fatalf("expected match block:\n%s", cfg)
	}
	if !strings.Contains(cfg, "status 200 204;") {
		t.Fatalf("expected status list in match block:\n%s", cfg)
	}
}

func TestFormatSize(t *testing.T) {
	cases := []struct {
		bytes int
		want  string
	}{
		{5*1024*1024 + 512*1024, "5.5m"},
		{5 * 1024 * 1024, "5m"},
		{8 * 1024, "8k"},
		{1536, "1.5k"},
		{512, "512"},
		{1024*1024 + 256*1024, "1.25m"},
	}
	for _, c := range cases {
		if got := formatSize(c.bytes); got != c.want {
			t.Errorf("formatSize(%d) = %q, want %q", c.bytes, got, c.want)
		}
	}
}
