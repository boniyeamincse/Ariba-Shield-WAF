package compiler

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"
)

//go:embed templates/shield.conf.tmpl
var templatesFS embed.FS

// Upstream models an nginx upstream block derived from a backend pool.
type Upstream struct {
	Name          string
	LeastConn     bool
	IPHash        bool
	ConsistentHash bool
	StickyCookie  string // non-empty => emit sticky cookie directive
	HealthCheck   bool   // active health checks enabled
	HealthInterval  string
	HealthTimeout   string
	HealthFailCount int
	HealthPassCount int
	HealthPath      string
	HealthStatus    string
	Servers       []UpstreamServer
}

type UpstreamServer struct {
	Host       string
	Port       int
	Weight     int
	MaxFails   int
	FailTimeout int
	Protocol   string // "http" or "https" (P1.3 TLS re-encryption)
	SlowStart  string // e.g. "30s" or "" (P1.5)
}

// Server models an nginx server block derived from a virtual server.
type Server struct {
	ListenAddr    string
	ListenPort    int
	SSL           bool
	HTTP2         bool
	ServerName    string
	SSLCertPath   string
	SSLKeyPath    string
	MaxBodySize   string
	MaxHeaderSize string
	Locations     []Location
}

type Location struct {
	Path     string
	Upstream string
	Timeout  int
	// Scheme is the upstream protocol (http/https) for TLS re-encryption (P1.3).
	Scheme string
}

// RenderNginxConfig produces the nginx configuration from a validated policy
// document. Deterministic: identical input yields identical bytes.
func RenderNginxConfig(doc *PolicyDocument, certPaths map[string]string) (string, error) {
	if err := validate(doc); err != nil {
		return "", err
	}

	// Build upstreams from backend pools.
	var upstreams []Upstream
	poolByName := make(map[string]string)   // poolID -> upstream name
	poolScheme := make(map[string]string)   // poolID -> http|https (P1.3)
	for _, pool := range doc.BackendPools {
		name := "pool_" + sanitize(pool.Name)
		poolByName[pool.ID] = name

		scheme := "http"
		for _, node := range pool.Nodes {
			if node.Protocol == "https" {
				scheme = "https"
				break
			}
		}
		poolScheme[pool.ID] = scheme

		up := Upstream{
			Name:           name,
			LeastConn:      pool.LBAlgorithm == "least_conn",
			IPHash:         pool.LBAlgorithm == "ip_hash" || (pool.Sticky && pool.StickyType == "ip_hash"),
			ConsistentHash: pool.LBAlgorithm == "consistent_hash",
		}
		if pool.Sticky && pool.StickyType == "cookie" {
			cookie := pool.CookieName
			if cookie == "" {
				cookie = "shield_sticky"
			}
			up.StickyCookie = cookie
		}

		// Active health checks (P1.6).
		if pool.HealthMonitor != nil && pool.HealthMonitor.Type == "http" {
			hm := pool.HealthMonitor
			interval := hm.IntervalMS
			if interval <= 0 {
				interval = 5000
			}
			timeout := hm.TimeoutMS
			if timeout <= 0 {
				timeout = 2000
			}
			fail := hm.FailThreshold
			if fail <= 0 {
				fail = 3
			}
			pass := hm.PassThreshold
			if pass <= 0 {
				pass = 2
			}
			path := hm.HTTPPath
			if path == "" {
				path = "/"
			}
			up.HealthCheck = true
			up.HealthInterval = fmt.Sprintf("%dms", interval)
			up.HealthTimeout = fmt.Sprintf("%dms", timeout)
			up.HealthFailCount = fail
			up.HealthPassCount = pass
			up.HealthPath = path
			if len(hm.HTTPExpectedStatus) > 0 {
				statuses := make([]string, 0, len(hm.HTTPExpectedStatus))
				for _, s := range hm.HTTPExpectedStatus {
					statuses = append(statuses, fmt.Sprintf("%d", s))
				}
				up.HealthStatus = strings.Join(statuses, " ")
			}
		}

		for _, node := range pool.Nodes {
			// Draining nodes keep serving in-flight but stop receiving new
			// requests: drop them from the upstream (P1.5).
			if !node.Active || node.Drain {
				continue
			}
			weight := node.Weight
			if weight == 0 {
				weight = 1
			}
			protocol := node.Protocol
			if protocol == "" {
				protocol = "http"
			}
			slowStart := ""
			if node.SlowStart > 0 {
				slowStart = fmt.Sprintf("%ds", node.SlowStart)
			}
			up.Servers = append(up.Servers, UpstreamServer{
				Host:        node.Host,
				Port:        node.Port,
				Weight:      weight,
				MaxFails:    3,
				FailTimeout: 30,
				Protocol:    protocol,
				SlowStart:   slowStart,
			})
		}
		if len(up.Servers) == 0 {
			return "", fmt.Errorf("pool %q has no active nodes", pool.Name)
		}
		upstreams = append(upstreams, up)
	}

	// Build server blocks from virtual servers.
	var servers []Server
	for _, vs := range doc.VirtualServers {
		upstream, ok := poolByName[vs.DefaultBackendPoolID]
		if !ok {
			return "", fmt.Errorf("virtual server %q references unknown default pool %q", vs.Name, vs.DefaultBackendPoolID)
		}

		addr := vs.ListenAddr
		if addr == "" {
			addr = "0.0.0.0"
		}

		http2 := false
		for _, p := range vs.TLS.Protocols {
			if p == "h2" {
				http2 = true
			}
		}

		s := Server{
			ListenAddr:    addr,
			ListenPort:    vs.ListenPort,
			SSL:           vs.TLS.Enabled,
			HTTP2:         http2,
			ServerName:    "_",
			MaxBodySize:   "10m",
			MaxHeaderSize: "8k",
		}

		if vs.TLS.Enabled {
			cert := certPaths[vs.TLS.CertificateRef]
			if cert == "" {
				return "", fmt.Errorf("virtual server %q: no cert path for ref %q", vs.Name, vs.TLS.CertificateRef)
			}
			s.SSLCertPath = cert + "/cert.pem"
			s.SSLKeyPath = cert + "/key.pem"
		}

		if vs.Limits != nil {
			if vs.Limits.MaxBodySize > 0 {
				s.MaxBodySize = fmt.Sprintf("%dm", vs.Limits.MaxBodySize/1024/1024)
			}
			if vs.Limits.MaxHeaderSize > 0 {
				s.MaxHeaderSize = fmt.Sprintf("%dk", vs.Limits.MaxHeaderSize/1024)
			}
		}

		// Routes default to the default pool when no path match applies.
		if len(vs.Routes) == 0 {
			s.Locations = append(s.Locations, Location{Path: "/", Upstream: upstream, Scheme: poolScheme[vs.DefaultBackendPoolID]})
		} else {
			for _, r := range vs.Routes {
				u, ok := poolByName[r.BackendPoolID]
				if !ok {
					return "", fmt.Errorf("virtual server %q: route %q references unknown pool %q", vs.Name, r.Path, r.BackendPoolID)
				}
				path := r.Path
				if r.Match == "exact" {
					// nginx `location = <path>` is the exact-match form.
					path = "= " + r.Path
				} else if !strings.HasSuffix(path, "/") {
					path = path + "/"
				}
				s.Locations = append(s.Locations, Location{Path: path, Upstream: u, Scheme: poolScheme[r.BackendPoolID]})
			}
		}

		servers = append(servers, s)
	}

	// Render template.
	data := struct {
		LogLevel        string
		GatewayID       string
		VirtualServerID string
		ApplicationID   string
		Upstreams       []Upstream
		Servers         []Server
	}{
		LogLevel:        doc.Settings.LogLevel,
		GatewayID:       first(doc.GatewayTargets, "local"),
		VirtualServerID: "vs",
		ApplicationID:   "app",
		Upstreams:       upstreams,
		Servers:         servers,
	}

	tmpl, err := template.ParseFS(templatesFS, "templates/shield.conf.tmpl")
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	return buf.String(), nil
}

func validate(doc *PolicyDocument) error {
	if doc.SchemaVersion != "0.1" {
		return fmt.Errorf("unsupported schema_version %q", doc.SchemaVersion)
	}
	if doc.ConfigID == "" {
		return fmt.Errorf("config_id is required")
	}
	if len(doc.VirtualServers) == 0 {
		return fmt.Errorf("at least one virtual server required")
	}
	if len(doc.BackendPools) == 0 {
		return fmt.Errorf("at least one backend pool required")
	}
	return nil
}

func sanitize(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_':
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}
	if b.Len() == 0 {
		return "pool"
	}
	return b.String()
}

func first(s []string, def string) string {
	if len(s) > 0 {
		return s[0]
	}
	return def
}