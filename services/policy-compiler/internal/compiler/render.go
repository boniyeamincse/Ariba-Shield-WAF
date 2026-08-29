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
	Name       string
	LeastConn  bool
	Servers    []UpstreamServer
}

type UpstreamServer struct {
	Host       string
	Port       int
	Weight     int
	MaxFails   int
	FailTimeout int
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
}

// RenderNginxConfig produces the nginx configuration from a validated policy
// document. Deterministic: identical input yields identical bytes.
func RenderNginxConfig(doc *PolicyDocument, certPaths map[string]string) (string, error) {
	if err := validate(doc); err != nil {
		return "", err
	}

	// Build upstreams from backend pools.
	var upstreams []Upstream
	poolByName := make(map[string]string) // poolID -> upstream name
	for _, pool := range doc.BackendPools {
		name := "pool_" + sanitize(pool.Name)
		poolByName[pool.ID] = name

		up := Upstream{Name: name, LeastConn: pool.LBAlgorithm == "least_conn"}
		for _, node := range pool.Nodes {
			if !node.Active {
				continue
			}
			weight := node.Weight
			if weight == 0 {
				weight = 1
			}
			up.Servers = append(up.Servers, UpstreamServer{
				Host:        node.Host,
				Port:        node.Port,
				Weight:      weight,
				MaxFails:    3,
				FailTimeout: 30,
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
			ServerName:    vs.Name + ".shield.local",
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
			s.Locations = append(s.Locations, Location{Path: "/", Upstream: upstream})
		} else {
			for _, r := range vs.Routes {
				u, ok := poolByName[r.BackendPoolID]
				if !ok {
					return "", fmt.Errorf("virtual server %q: route %q references unknown pool %q", vs.Name, r.Path, r.BackendPoolID)
				}
				path := r.Path
				if r.Match == "prefix" && !strings.HasSuffix(path, "/") {
					path = path + "/"
				}
				s.Locations = append(s.Locations, Location{Path: path, Upstream: u})
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
		VirtualServerID: "vs",  // placeholder; one server block each
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