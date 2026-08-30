package compiler

import (
	"bytes"
	"embed"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

//go:embed templates/shield.conf.tmpl
var templatesFS embed.FS

// Upstream models an nginx upstream block derived from a backend pool.
type Upstream struct {
	Name            string
	LeastConn       bool
	IPHash          bool
	ConsistentHash  bool
	StickyCookie    string // non-empty => emit sticky cookie directive
	HealthCheck     bool   // active health checks enabled
	HealthInterval  string
	HealthTimeout   string
	HealthFailCount int
	HealthPassCount int
	HealthPath      string
	HealthStatus    string
	HealthMatchName string // name of the `match` block for expected statuses, "" if none
	Servers         []UpstreamServer
}

type UpstreamServer struct {
	Host        string
	Port        int
	Weight      int
	MaxFails    int
	FailTimeout int
	Protocol    string // "http" or "https" (P1.3 TLS re-encryption)
	SlowStart   string // e.g. "30s" or "" (P1.5)
}

// Server models an nginx server block derived from a virtual server.
type Server struct {
	ListenAddr      string
	ListenPort      int
	SSL             bool
	HTTP2           bool
	ServerName      string
	SSLProtocols    string
	SSLCertPath     string
	SSLKeyPath      string
	MaxBodySize     string
	MaxHeaderSize   string
	VirtualServerID string
	ApplicationID   string
	Locations       []Location
}

type Location struct {
	Path     string
	Upstream string
	Timeout  int
	// Scheme is the upstream protocol (http/https) for TLS re-encryption (P1.3).
	Scheme string
}

// HealthMatch models an nginx `match` block used by active health checks
// (P3.39). One block is emitted per distinct expected-status set.
type HealthMatch struct {
	Name   string
	Status string
}

// RenderNginxConfig produces the nginx configuration from a validated policy
// document. Deterministic: identical input yields identical bytes.
func RenderNginxConfig(doc *PolicyDocument, certPaths map[string]string) (string, error) {
	if err := validate(doc); err != nil {
		return "", err
	}

	// Build upstreams from backend pools.
	var upstreams []Upstream
	var healthMatches []HealthMatch
	healthMatchNames := make(map[string]string) // expected-status set -> match name
	poolByName := make(map[string]string)       // poolID -> upstream name
	poolScheme := make(map[string]string)       // poolID -> http|https (P1.3)
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
				name, ok := healthMatchNames[up.HealthStatus]
				if !ok {
					name = "shield_health_" + sanitize(up.HealthStatus)
					healthMatchNames[up.HealthStatus] = name
					healthMatches = append(healthMatches, HealthMatch{Name: name, Status: up.HealthStatus})
				}
				up.HealthMatchName = name
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
			ListenAddr:      addr,
			ListenPort:      vs.ListenPort,
			SSL:             vs.TLS.Enabled,
			HTTP2:           http2,
			ServerName:      serverNameFor(vs.Name),
			SSLProtocols:    tlsProtocols(vs.TLS.MinVersion),
			MaxBodySize:     "10m",
			MaxHeaderSize:   "8k",
			VirtualServerID: vs.ID,
			ApplicationID:   appIDForPool(doc, vs.DefaultBackendPoolID),
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
				s.MaxBodySize = formatSize(vs.Limits.MaxBodySize)
			}
			if vs.Limits.MaxHeaderSize > 0 {
				s.MaxHeaderSize = formatSize(vs.Limits.MaxHeaderSize)
			}
		}

		// Routes default to the default pool when no path match applies.
		if len(vs.Routes) == 0 {
			s.Locations = append(s.Locations, Location{Path: "/", Upstream: upstream, Scheme: poolScheme[vs.DefaultBackendPoolID]})
		} else {
			// P3.36: higher-priority routes are declared first so nginx picks
			// them first; duplicate/conflicting locations are rejected.
			routes := append([]Route(nil), vs.Routes...)
			sort.SliceStable(routes, func(i, j int) bool {
				return routes[i].Priority > routes[j].Priority
			})
			seen := make(map[string]string) // normalized location -> route id
			for _, r := range routes {
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
				if prev, dup := seen[path]; dup {
					return "", fmt.Errorf("virtual server %q: route %q conflicts with route %q at location %q", vs.Name, r.ID, prev, path)
				}
				seen[path] = r.ID
				s.Locations = append(s.Locations, Location{Path: path, Upstream: u, Scheme: poolScheme[r.BackendPoolID]})
			}
		}

		servers = append(servers, s)
	}

	// Render template.
	data := struct {
		LogLevel      string
		GatewayID     string
		Upstreams     []Upstream
		HealthMatches []HealthMatch
		Servers       []Server
	}{
		LogLevel:      doc.Settings.LogLevel,
		GatewayID:     first(doc.GatewayTargets, "local"),
		Upstreams:     upstreams,
		HealthMatches: healthMatches,
		Servers:       servers,
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

// appIDForPool returns the application_id of the first pool whose ID matches,
// or the pool's own application_id if found. Used to stamp X-Shield-App-ID.
func appIDForPool(doc *PolicyDocument, poolID string) string {
	for _, p := range doc.BackendPools {
		if p.ID == poolID {
			return p.ApplicationID
		}
	}
	return ""
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

// serverNameFor returns the hostname for the server block. The policy virtual
// server carries its hostname in the Name field; an empty name falls back to
// nginx's catch-all `_` (P3.34).
func serverNameFor(name string) string {
	if name == "" {
		return "_"
	}
	return name
}

// tlsProtocols maps the policy TLSProfile.MinVersion to the nginx
// ssl_protocols value, rendering only the configured TLS versions. The safe
// default is TLSv1.2 TLSv1.3 (P3.33).
func tlsProtocols(minVersion string) string {
	switch strings.ToLower(strings.TrimSpace(minVersion)) {
	case "1.0", "tls1.0", "tlsv1", "tlsv1.0":
		return "TLSv1 TLSv1.1 TLSv1.2 TLSv1.3"
	case "1.1", "tls1.1", "tlsv1.1":
		return "TLSv1.1 TLSv1.2 TLSv1.3"
	case "1.2", "tls1.2", "tlsv1.2":
		return "TLSv1.2 TLSv1.3"
	case "1.3", "tls1.3", "tlsv1.3":
		return "TLSv1.3"
	default:
		return "TLSv1.2 TLSv1.3"
	}
}

// formatSize renders a byte count as an nginx size directive, preserving
// sub-unit precision (5.5 MB becomes "5.5m", not "5m") (P3.38). Sizes at or
// above 1 MB are expressed in MB, sizes at or above 1 KB in KB, and anything
// smaller in raw bytes.
func formatSize(bytes int) string {
	if bytes <= 0 {
		return ""
	}
	switch {
	case bytes%(1024*1024) == 0:
		return fmt.Sprintf("%dm", bytes/(1024*1024))
	case bytes >= 1024*1024:
		return trimZeros(strconv.FormatFloat(float64(bytes)/(1024*1024), 'f', -1, 64)) + "m"
	case bytes%1024 == 0:
		return fmt.Sprintf("%dk", bytes/1024)
	case bytes < 1024:
		return fmt.Sprintf("%d", bytes)
	default:
		return trimZeros(strconv.FormatFloat(float64(bytes)/1024, 'f', -1, 64)) + "k"
	}
}

// trimZeros removes trailing ".0" from a decimal string produced by
// strconv.FormatFloat(..., 'f', -1, ...) when the value is a whole number.
func trimZeros(s string) string {
	if strings.HasSuffix(s, ".0") {
		return strings.TrimSuffix(s, ".0")
	}
	return s
}

func first(s []string, def string) string {
	if len(s) > 0 {
		return s[0]
	}
	return def
}
