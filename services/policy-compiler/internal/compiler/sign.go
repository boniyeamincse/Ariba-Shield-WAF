package compiler

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// canonicalUnsignedJSON returns the canonical JSON of the document with the
// signature field stripped, so the signature is computed over the payload only.
func (d *PolicyDocument) canonicalUnsignedJSON() ([]byte, error) {
	cp := *d
	cp.Signature = nil
	return json.Marshal(&cp)
}

// SignBundle signs the canonical (unsigned) JSON of a policy document with an
// ed25519 private key and attaches the signature to the document (ADR-002).
func SignBundle(doc *PolicyDocument, key ed25519.PrivateKey, keyID string) error {
	canon, err := doc.canonicalUnsignedJSON()
	if err != nil {
		return fmt.Errorf("canonical json: %w", err)
	}
	sig := ed25519.Sign(key, canon)
	doc.Signature = &Signature{
		Algorithm: "ed25519",
		KeyID:     keyID,
		Value:     base64.StdEncoding.EncodeToString(sig),
	}
	return nil
}

// VerifyBundleSignature verifies the ed25519 signature over the canonical
// unsigned JSON of the policy document.
func VerifyBundleSignature(doc *PolicyDocument, pub ed25519.PublicKey) error {
	if doc.Signature == nil {
		return fmt.Errorf("missing signature")
	}
	if doc.Signature.Algorithm != "ed25519" {
		return fmt.Errorf("unsupported signature algorithm %q", doc.Signature.Algorithm)
	}
	sig, err := base64.StdEncoding.DecodeString(doc.Signature.Value)
	if err != nil {
		return fmt.Errorf("decode signature: %w", err)
	}
	canon, err := doc.canonicalUnsignedJSON()
	if err != nil {
		return fmt.Errorf("canonical json: %w", err)
	}
	if !ed25519.Verify(pub, canon, sig) {
		return fmt.Errorf("signature verification failed")
	}
	return nil
}
