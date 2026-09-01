package acme

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/acme"
)

// ProvisionResult holds the result of an ACME certificate issuance.
type ProvisionResult struct {
	CertPEM   string
	KeyPEM    string
	NotBefore time.Time
	NotAfter  time.Time
}

// Provision obtains a certificate from an ACME CA (e.g. Let's Encrypt).
// It uses the HTTP-01 challenge which requires the target domain to be
// publicly reachable on port 80 from the internet. The caller is responsible
// for ensuring the WAF gateway serves /.well-known/acme-challenge/.
func Provision(ctx context.Context, domain, email string, staging bool) (*ProvisionResult, error) {
	dir := acme.LetsEncryptURL
	if staging {
		dir = "https://acme-staging-v02.api.letsencrypt.org/directory"
	}

	accountKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("account key: %w", err)
	}

	cl := &acme.Client{
		Key:          accountKey,
		DirectoryURL: dir,
	}

	acct := &acme.Account{Contact: []string{"mailto:" + email}}
	if _, err := cl.Register(ctx, acct, acme.AcceptTOS); err != nil {
		return nil, fmt.Errorf("acme register: %w", err)
	}

	order, err := cl.AuthorizeOrder(ctx, acme.DomainIDs(domain))
	if err != nil {
		return nil, fmt.Errorf("acme order: %w", err)
	}

	// Complete all authorizations via HTTP-01.
	for _, authURL := range order.AuthzURLs {
		auth, err := cl.GetAuthorization(ctx, authURL)
		if err != nil {
			return nil, fmt.Errorf("acme auth: %w", err)
		}

		var chal *acme.Challenge
		for _, c := range auth.Challenges {
			if c.Type == "http-01" {
				chal = c
				break
			}
		}
		if chal == nil {
			return nil, fmt.Errorf("no http-01 challenge for %s", auth.Identifier.Value)
		}

		// The WAF gateway must serve:
		//   GET /.well-known/acme-challenge/{chal.Token} -> response
		// The response is returned to the caller so the gateway can configure it.
		if _, err := cl.Accept(ctx, chal); err != nil {
			return nil, fmt.Errorf("acme accept: %w", err)
		}

		if _, err := cl.WaitAuthorization(ctx, authURL); err != nil {
			return nil, fmt.Errorf("acme wait auth: %w", err)
		}
	}

	// Generate a CSR and finalize the order.
	certKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("cert key: %w", err)
	}
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{
		DNSNames: []string{domain},
	}, certKey)
	if err != nil {
		return nil, fmt.Errorf("csr: %w", err)
	}

	der, _, err := cl.CreateOrderCert(ctx, order.FinalizeURL, csrDER, true)
	if err != nil {
		return nil, fmt.Errorf("acme finalize: %w", err)
	}
	if len(der) == 0 {
		return nil, fmt.Errorf("acme: empty cert chain")
	}

	var notBefore, notAfter time.Time
	if leaf, err := x509.ParseCertificate(der[0]); err == nil {
		notBefore = leaf.NotBefore
		notAfter = leaf.NotAfter
	}

	var certPEM string
	for _, b := range der {
		certPEM += string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: b}))
	}

	keyDER, err := x509.MarshalPKCS8PrivateKey(certKey)
	if err != nil {
		return nil, fmt.Errorf("marshal key: %w", err)
	}
	keyPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER}))

	return &ProvisionResult{
		CertPEM:   certPEM,
		KeyPEM:    keyPEM,
		NotBefore: notBefore,
		NotAfter:  notAfter,
	}, nil
}

// HTTP01Handler serves ACME HTTP-01 challenge tokens. Mount at
// /.well-known/acme-challenge/{token} on the WAF gateway.
func HTTP01Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
}
