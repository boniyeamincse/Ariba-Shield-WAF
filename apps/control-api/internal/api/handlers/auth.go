package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ariba-shield/control-api/internal/store"
	"github.com/oklog/ulid/v2"
)

// mockAuthEnabled gates the dev mock-login bypass. It MUST be off by default
// (P0.4): the mock path accepts the hardcoded password "admin", which is not
// acceptable in any non-dev environment. Enable with AUTH_MOCK_ENABLED=true.
func mockAuthEnabled() bool {
	return os.Getenv("AUTH_MOCK_ENABLED") == "true"
}

// Login authenticates with email+password, creates a session, returns a cookie.
func Login(st *store.Store) http.HandlerFunc {
	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body loginRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.Email == "" || body.Password == "" {
			http.Error(w, `{"error":"email and password required"}`, http.StatusBadRequest)
			return
		}

		// DEV BYPASS (P0.4): only active when AUTH_MOCK_ENABLED=true. The mock
		// path accepts the hardcoded password "admin"; never enable in prod.
		if mockAuthEnabled() {
			mockRoles := map[string]string{
				"superadmin@aribashield.local": "SUPER_ADMIN",
				"platform@aribashield.local":   "PLATFORM_ADMIN",
				"security@aribashield.local":   "SECURITY_ADMIN",
				"appowner@aribashield.local":   "APP_OWNER",
				"soc@aribashield.local":        "SOC_ANALYST",
				"auditor@aribashield.local":    "AUDITOR",
				"readonly@aribashield.local":   "READ_ONLY",
			}

			if role, isMock := mockRoles[body.Email]; isMock && body.Password == "admin" {
				expiresAt := time.Now().Add(24 * time.Hour)
				http.SetCookie(w, &http.Cookie{
					Name:     "shield_session",
					Value:    "mock_session_token",
					Path:     "/",
					Expires:  expiresAt,
					HttpOnly: true,
					Secure:   false, // False for localhost dev
					SameSite: http.SameSiteLaxMode,
				})
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]any{
					"user_id":  "usr_mock_001",
					"email":    body.Email,
					"role":     role,
					"expires":  expiresAt.Format(time.RFC3339),
				})
				return
			}
		}

		// Look up the user and verify the password hash.
		var userID, passwordHash, orgID string
		var totpEnabled bool
		err := st.Pool.QueryRow(r.Context(),
			`SELECT u.id, u.password_hash, u.organization_id, u.totp_enabled
			 FROM users u WHERE u.email = $1 AND u.status = 'active'`,
			body.Email).Scan(&userID, &passwordHash, &orgID, &totpEnabled)
		if err != nil {
			http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
			return
		}

		if !store.CheckPassword(body.Password, passwordHash) {
			http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
			return
		}

		role, err := st.LookupUserRole(r.Context(), userID)
		if err != nil {
			role = "Read Only"
		}

		// Create a session.
		sessionID := ulid.Make().String()
		tokenBytes := make([]byte, 32)
		if _, err := rand.Read(tokenBytes); err != nil {
			http.Error(w, `{"error":"session creation failed"}`, http.StatusInternalServerError)
			return
		}
		token := hex.EncodeToString(tokenBytes)

		expiresAt := sessionExpiry(st, r.Context())

		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO sessions (id, user_id, token_hash, expires_at, request_id)
			 VALUES ($1, $2, $3, $4, $5)`,
			sessionID, userID, token, expiresAt, r.Header.Get("X-Request-ID")); err != nil {
			http.Error(w, `{"error":"session creation failed"}`, http.StatusInternalServerError)
			return
		}

		// Set the session cookie.
		http.SetCookie(w, &http.Cookie{
			Name:     "shield_session",
			Value:    fmt.Sprintf("%s:%s", userID, token),
			Path:     "/",
			Expires:  expiresAt,
			HttpOnly: true,
			Secure:   r.TLS != nil, // false for localhost HTTP dev
			SameSite: http.SameSiteStrictMode,
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"user_id":  userID,
			"email":    body.Email,
			"role":     role,
			"expires":  expiresAt.Format(time.RFC3339),
		})
	}
}

// Logout invalidates the session.
func Logout(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("shield_session")
		if err == nil && cookie.Value != "" {
			// Extract user ID from "userID:token" format.
			userID := cookie.Value
			if i := indexOf(cookie.Value, ":"); i != -1 {
				userID = cookie.Value[:i]
			}
			_, _ = st.Pool.Exec(r.Context(),
				`DELETE FROM sessions WHERE user_id = $1`, userID)
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "shield_session",
			Value:    "",
			Path:     "/",
			Expires:  time.Now().Add(-1 * time.Hour),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "logged_out"})
	}
}

// Me returns the current authenticated user's details.
func Me(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("shield_session")
		if err != nil || cookie.Value == "" {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		// DEV bypass (P0.4): only when AUTH_MOCK_ENABLED=true.
		if mockAuthEnabled() && cookie.Value == "mock_session_token" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"user": map[string]string{
					"id":    "usr_mock_001",
					"name":  "Mock User",
					"role":  "Mock Role",
				},
			})
			return
		}

		// Parse "userID:token".
		val := cookie.Value
		var userID, token string
		for i := 0; i < len(val); i++ {
			if val[i] == ':' {
				userID = val[:i]
				token = val[i+1:]
				break
			}
		}
		if userID == "" || token == "" {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		// Fetch the real user from the session.
		var email string
		err = st.Pool.QueryRow(r.Context(),
			`SELECT u.email FROM sessions s
			 JOIN users u ON u.id = s.user_id
			 WHERE s.token_hash = $1 AND s.user_id = $2 AND s.expires_at > now()`,
			token, userID).Scan(&email)
		if err != nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		role, err := st.LookupUserRole(r.Context(), userID)
		if err != nil {
			role = "Read Only"
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"user": map[string]string{
				"id":    userID,
				"email": email,
				"role":  role,
			},
		})
	}
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}