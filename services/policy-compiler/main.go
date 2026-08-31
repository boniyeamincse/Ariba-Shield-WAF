package main

import (
	"crypto/ed25519"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ariba-shield/policy-compiler/internal/compiler"
)

func main() {
	policyPath := flag.String("policy", "", "path to policy document (JSON, ADR-002 schema v0)")
	outputPath := flag.String("output", "-", "output path for generated nginx config (default stdout)")
	wafOutputPath := flag.String("waf-output", "", "output path for generated WAF rules (Coraza config)")
	certDir := flag.String("cert-dir", "/etc/shield-waf/certs", "base directory for certificate_ref lookups (<ref>/cert.pem, <ref>/key.pem inside)")
	signKey := flag.String("sign-key", "", "path to ed25519 private key PEM; if set, signs the bundle before rendering")
	signKeyID := flag.String("sign-key-id", "control-plane-01", "key ID attached to the signature")
	flag.Parse()

	if *policyPath == "" {
		log.Fatal("--policy is required")
	}

	data, err := os.ReadFile(*policyPath)
	if err != nil {
		log.Fatalf("read policy: %v", err)
	}

	doc, err := compiler.ParsePolicy(data)
	if err != nil {
		log.Fatalf("parse policy: %v", err)
	}

	hash, err := doc.BundleHash()
	if err != nil {
		log.Fatalf("bundle hash: %v", err)
	}

	// Build cert lookup: each certificate_ref maps to its own directory.
	certPaths := map[string]string{}
	for _, vs := range doc.VirtualServers {
		if vs.TLS.Enabled && vs.TLS.CertificateRef != "" {
			certPaths[vs.TLS.CertificateRef] = filepath.Join(*certDir, vs.TLS.CertificateRef)
		}
	}

	// Phase 4: sign the bundle if a signing key is provided.
	if *signKey != "" {
		keyBytes, err := os.ReadFile(*signKey)
		if err != nil {
			log.Fatalf("read signing key: %v", err)
		}
		priv := ed25519.PrivateKey(keyBytes)
		if err := compiler.SignBundle(doc, priv, *signKeyID); err != nil {
			log.Fatalf("sign bundle: %v", err)
		}
	}

	// 1. Render Nginx Config
	config, err := compiler.RenderNginxConfig(doc, certPaths)
	if err != nil {
		log.Fatalf("render nginx config: %v", err)
	}

	out := os.Stdout
	if *outputPath != "-" {
		f, err := os.Create(*outputPath)
		if err != nil {
			log.Fatalf("create nginx output: %v", err)
		}
		defer f.Close()
		out = f
	}

	if _, err := fmt.Fprintf(out, "# Ariba Shield WAF generated config\n# config_id=%s bundle_hash=%s\n%s", doc.ConfigID, hash, config); err != nil {
		log.Fatalf("write nginx output: %v", err)
	}
	log.Printf("generated nginx config for %s (bundle_hash=%s)", doc.ConfigID, hash)

	// 2. Render WAF Config (if path provided)
	if *wafOutputPath != "" {
		wafConfig, err := compiler.RenderWAFConfig(doc)
		if err != nil {
			log.Fatalf("render WAF config: %v", err)
		}

		wafFile, err := os.Create(*wafOutputPath)
		if err != nil {
			log.Fatalf("create waf output: %v", err)
		}
		defer wafFile.Close()

		if _, err := wafFile.WriteString(wafConfig); err != nil {
			log.Fatalf("write waf output: %v", err)
		}
		log.Printf("generated waf config for %s at %s", doc.ConfigID, *wafOutputPath)
	}
}