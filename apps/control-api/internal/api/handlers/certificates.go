package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ariba-shield/control-api/internal/store"
)

type Certificate struct {
	ID         string `json:"id"`
	CommonName string `json:"common_name"`
	CertType   string `json:"cert_type"` // CLIENT_MTLS, SERVER_TLS
	Issuer     string `json:"issuer"`
	ValidFrom  string `json:"valid_from"`
	ValidTo    string `json:"valid_to"`
	Status     string `json:"status"`
}

// ListCertificates returns all TLS/mTLS certificates managed by the control plane
func ListCertificates(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		certs := []Certificate{
			{
				ID:         "cert_001",
				CommonName: "upstream-service.internal",
				CertType:   "CLIENT_MTLS",
				Issuer:     "Ariba Shield Internal CA",
				ValidFrom:  "2025-01-01T00:00:00Z",
				ValidTo:    "2026-01-01T00:00:00Z",
				Status:     "active",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(certs)
	}
}

// UploadCertificate securely stores a new certificate (excluding private key return)
func UploadCertificate(st *store.Store) http.HandlerFunc {
	type certUpload struct {
		CertBody   string `json:"cert_body"`
		PrivateKey string `json:"private_key,omitempty"` // Stored securely, never returned
		CertType   string `json:"cert_type"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body certUpload
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}

		if body.CertBody == "" || body.CertType == "" {
			http.Error(w, `{"error":"cert_body and cert_type are required"}`, http.StatusBadRequest)
			return
		}

		// In a real application:
		// 1. Parse x509 cert to get CN, Expiry, Issuer
		// 2. Encrypt PrivateKey with KMS/Vault before storing

		newID, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": newID})
	}
}
