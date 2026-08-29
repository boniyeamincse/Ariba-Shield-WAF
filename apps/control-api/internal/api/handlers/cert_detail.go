package handlers

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"time"

	"github.com/ariba-shield/control-api/internal/store"
)

// GetCertificate returns a single certificate by id.
func GetCertificate(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var b []byte
		err := st.Pool.QueryRow(r.Context(),
			"SELECT row_to_json(t) FROM (SELECT id, name, domain, issuer, serial, not_before, not_after, status FROM certificates WHERE id = $1) t", id).Scan(&b)
		if err != nil {
			http.Error(w, `{"error":"certificate not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

// DeleteCertificate removes a certificate.
func DeleteCertificate(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ct, err := st.Pool.Exec(r.Context(), `DELETE FROM certificates WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error":"delete failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"certificate not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	}
}

// ImportCertificate imports a PEM certificate chain (parses + stores metadata).
func ImportCertificate(st *store.Store) http.HandlerFunc {
	type importReq struct {
		Name        string `json:"name"`
		Domain      string `json:"domain"`
		Certificate string `json:"certificate"` // PEM chain
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body importReq
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.Name == "" || body.Domain == "" || body.Certificate == "" {
			http.Error(w, `{"error":"name, domain and certificate required"}`, http.StatusBadRequest)
			return
		}

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

		status := "active"
		if time.Now().After(cert.NotAfter) {
			status = "expired"
		} else if time.Until(cert.NotAfter) < 30*24*time.Hour {
			status = "expiring"
		}

		id, err := st.NewID()
		if err != nil {
			http.Error(w, `{"error":"id generation failed"}`, http.StatusInternalServerError)
			return
		}
		orgID := "01ARZ3NDEKTSV4RRFFQ69G5FAV"

		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO certificates
			   (id, organization_id, name, domain, issuer, serial, not_before, not_after, status, chain_pem)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			id, orgID, body.Name, body.Domain, cert.Issuer.CommonName, cert.SerialNumber.String(),
			cert.NotBefore, cert.NotAfter, status, body.Certificate); err != nil {
			http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": id, "status": status, "serial": cert.SerialNumber.String()})
	}
}

// RenewCertificate marks a certificate as renewed (updates metadata + extends expiry).
// In production this would trigger ACME or upload a new chain; here it accepts
// a new PEM chain.
func RenewCertificate(st *store.Store) http.HandlerFunc {
	type renewReq struct {
		Certificate string `json:"certificate"` // new PEM chain
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var body renewReq
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Certificate == "" {
			http.Error(w, `{"error":"certificate required"}`, http.StatusBadRequest)
			return
		}

		block, _ := pem.Decode([]byte(body.Certificate))
		if block == nil {
			http.Error(w, `{"error":"invalid PEM"}`, http.StatusBadRequest)
			return
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			http.Error(w, `{"error":"invalid certificate"}`, http.StatusBadRequest)
			return
		}

		status := "active"
		if time.Now().After(cert.NotAfter) {
			status = "expired"
		} else if time.Until(cert.NotAfter) < 30*24*time.Hour {
			status = "expiring"
		}

		ct, err := st.Pool.Exec(r.Context(),
			`UPDATE certificates SET
			   issuer = $1, serial = $2, not_before = $3, not_after = $4,
			   status = $5, chain_pem = $6, version = version + 1, updated_at = now()
			 WHERE id = $7`,
			cert.Issuer.CommonName, cert.SerialNumber.String(), cert.NotBefore, cert.NotAfter,
			status, body.Certificate, id)
		if err != nil {
			http.Error(w, `{"error":"renew failed"}`, http.StatusInternalServerError)
			return
		}
		if ct.RowsAffected() == 0 {
			http.Error(w, `{"error":"certificate not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "status": status, "serial": cert.SerialNumber.String()})
	}
}

// CertExpiry returns the expiry info for a certificate.
func CertExpiry(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var name, domain, status string
		var notAfter time.Time
		err := st.Pool.QueryRow(r.Context(),
			`SELECT name, domain, not_after, status FROM certificates WHERE id = $1`, id).
			Scan(&name, &domain, &notAfter, &status)
		if err != nil {
			http.Error(w, `{"error":"certificate not found"}`, http.StatusNotFound)
			return
		}
		now := time.Now()
		daysLeft := int(notAfter.Sub(now).Hours() / 24)
		if daysLeft < 0 {
			daysLeft = 0
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id": id, "name": name, "domain": domain,
			"not_after": notAfter.Format(time.RFC3339),
			"days_left": daysLeft, "status": status,
		})
	}
}