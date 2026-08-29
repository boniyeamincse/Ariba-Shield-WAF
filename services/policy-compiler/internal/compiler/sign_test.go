package compiler

import (
	"crypto/ed25519"
	"testing"
)

func newTestDoc() *PolicyDocument {
	return &PolicyDocument{
		SchemaVersion:  "0.1",
		ConfigID:       "test-sign",
		CreatedAt:      "2026-08-29T00:00:00Z",
		CreatedBy:      "user-01",
		GatewayTargets: []string{"gw-01"},
		Settings:       Settings{LogLevel: "info", EventRetentionDays: 30},
		VirtualServers: []VirtualServer{
			{ID: "vs-01", Name: "app", ListenAddr: "0.0.0.0", ListenPort: 443, TLS: TLSProfile{Enabled: false}, DefaultBackendPoolID: "pool-01"},
		},
		BackendPools: []BackendPool{
			{ID: "pool-01", Name: "pool", ApplicationID: "app-01", LBAlgorithm: "round_robin", Nodes: []BackendNode{{ID: "n-01", Host: "10.0.0.1", Port: 8080, Weight: 1, Active: true}}},
		},
	}
}

func TestSignAndVerifyBundle(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}

	doc := newTestDoc()
	if err := SignBundle(doc, priv, "key-01"); err != nil {
		t.Fatalf("SignBundle: %v", err)
	}
	if doc.Signature == nil {
		t.Fatal("expected signature to be attached")
	}
	if doc.Signature.KeyID != "key-01" {
		t.Fatalf("expected key_id key-01, got %s", doc.Signature.KeyID)
	}

	if err := VerifyBundleSignature(doc, pub); err != nil {
		t.Fatalf("VerifyBundleSignature: %v", err)
	}
}

func TestVerifyBundleTampered(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)

	doc := newTestDoc()
	_ = SignBundle(doc, priv, "key-01")

	// Tamper with the document content AFTER signing.
	doc.Settings.LogLevel = "debug"

	if err := VerifyBundleSignature(doc, pub); err == nil {
		t.Fatal("expected verification to fail for tampered document")
	}
}

func TestVerifyBundleWrongKey(t *testing.T) {
	_, priv, _ := ed25519.GenerateKey(nil)
	otherPub, _, _ := ed25519.GenerateKey(nil)

	doc := newTestDoc()
	_ = SignBundle(doc, priv, "key-01")

	if err := VerifyBundleSignature(doc, otherPub); err == nil {
		t.Fatal("expected verification to fail with wrong public key")
	}
}

func TestVerifyBundleMissingSignature(t *testing.T) {
	pub, _, _ := ed25519.GenerateKey(nil)
	doc := newTestDoc()
	if err := VerifyBundleSignature(doc, pub); err == nil {
		t.Fatal("expected error for missing signature")
	}
}

func TestSignChangesHashStably(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)

	doc1 := newTestDoc()
	h1, _ := doc1.BundleHash()
	_ = SignBundle(doc1, priv, "key-01")
	h1signed, _ := doc1.BundleHash()
	if h1 != h1signed {
		t.Fatalf("bundle hash should not change after signing: %s != %s", h1, h1signed)
	}

	// A different doc signed with the same key has a different hash.
	doc2 := newTestDoc()
	doc2.Settings.LogLevel = "warn"
	h2, _ := doc2.BundleHash()
	if h1 == h2 {
		t.Fatal("different documents should have different hashes")
	}

	_ = pub
}