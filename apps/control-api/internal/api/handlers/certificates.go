package handlers

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"time"

	"github.com/ariba-shield/control-api/internal/store"
)

type certificate struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Domain    string `json:"domain"`
	Issuer    string `json:"issuer"`
	Serial    string `json:"serial"`
	NotBefore string `json:"not_before"`
	NotAfter  string `json:"not_after"`
	Status    string `json:"status"`
}

// ListCertificates returns certificate metadata (never key material, §7.2).
func ListCertificates(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := st.Pool.Query(r.Context(),
			`SELECT id, name, domain, issuer, serial, not_before, not_after, status
			 FROM certificates ORDER BY domain`)
		if err != nil {
			http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		certs := []certificate{}
		for rows.Next() {
			var c certificate
			var nb, na time.Time
			if err := rows.Scan(&c.ID, &c.Name, &c.Domain, &c.Issuer, &c.Serial, &nb, &na, &c.Status); err != nil {
				continue
			}
			c.NotBefore = nb.Format(time.RFC3339)
			c.NotAfter = na.Format(time.RFC3339)
			certs = append(certs, c)
		}
		if certs == nil {
			certs = []certificate{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(certs)
	}
}

// UploadCertificate imports a certificate. Only the certificate chain is
// parsed for metadata; the private key is never stored or returned (§7.2).
func UploadCertificate(st *store.Store) http.HandlerFunc {
	type upload struct {
		Name        string `json:"name"`
		Domain      string `json:"domain"`
		Certificate string `json:"certificate"` // PEM
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body upload
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.Name == "" || body.Domain == "" || body.Certificate == "" {
			http.Error(w, `{"error":"name, domain and certificate required"}`, http.StatusBadRequest)
			return
		}

		// Parse the first PEM block to extract metadata.
		block, _ := pem.Decode([]byte(body.Certificate))
		if block == nil {
			http.Error(w, `{"error":"invalid PEM certificate"}`, http.StatusBadRequest)
			return
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			http.Error(w, `{"error":"invalid certificate"}`, http.StatusBadRequest)
			return
		}

		id, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}

		status := "active"
		if time.Now().After(cert.NotAfter) {
			status = "expired"
		} else if time.Until(cert.NotAfter) < 30*24*time.Hour {
			status = "expiring"
		}

		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO certificates
			   (id, organization_id, name, domain, issuer, serial, not_before, not_after, status)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			id, orgID, body.Name, body.Domain, cert.Issuer.CommonName, cert.SerialNumber.String(),
			cert.NotBefore, cert.NotAfter, status); err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": id, "status": status})
	}
}
