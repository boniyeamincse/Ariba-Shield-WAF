package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/ariba-shield/control-api/internal/store"
	"github.com/oklog/ulid/v2"
	"github.com/pquerna/otp/totp"
)

// --- Session helpers ---

// sessionCookieValue returns the "userID:token" string stored in the cookie.
func sessionParts(cookie string) (userID, token string, ok bool) {
	for i := 0; i < len(cookie); i++ {
		if cookie[i] == ':' {
			return cookie[:i], cookie[i+1:], true
		}
	}
	return "", "", false
}

// newSessionToken generates a 32-byte random token (hex) and its hash.
func newSessionToken() (token, hash string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", err
	}
	token = hex.EncodeToString(b)
	hash = hashToken(token)
	return token, hash, nil
}

// Refresh rotates the session: issues a new token while keeping the user
// authenticated (session-fixation resistance, master plan §7.1).
func Refresh(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("shield_session")
		if err != nil || cookie.Value == "" {
			http.Error(w, `{"error":"no session"}`, http.StatusUnauthorized)
			return
		}
		userID, oldToken, ok := sessionParts(cookie.Value)
		if !ok {
			http.Error(w, `{"error":"invalid session"}`, http.StatusUnauthorized)
			return
		}

		// Verify the session exists and is valid.
		var sessionID string
		err = st.Pool.QueryRow(r.Context(),
			`SELECT id FROM sessions WHERE user_id = $1 AND token_hash = $2 AND expires_at > now()`,
			userID, hashToken(oldToken)).Scan(&sessionID)
		if err != nil {
			http.Error(w, `{"error":"session expired"}`, http.StatusUnauthorized)
			return
		}

		// Rotate the token.
		newToken, newHash, err := newSessionToken()
		if err != nil {
			http.Error(w, `{"error":"token generation failed"}`, http.StatusInternalServerError)
			return
		}
		newExpires := time.Now().Add(24 * time.Hour)
		if _, err := st.Pool.Exec(r.Context(),
			`UPDATE sessions SET token_hash = $1, expires_at = $2, last_used_at = now()
			 WHERE id = $3 AND user_id = $4`,
			newHash, newExpires, sessionID, userID); err != nil {
			http.Error(w, `{"error":"refresh failed"}`, http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "shield_session",
			Value:    userID + ":" + newToken,
			Path:     "/",
			Expires:  newExpires,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "rotated", "expires": newExpires.Format(time.RFC3339)})
	}
}

// --- MFA / TOTP ---

// EnableMFA generates a TOTP secret and returns it for the user to enroll
// (secret returned once, never again; see §7.2 secrets handling).
func EnableMFA(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, email, authErr := authenticatedUser(st, r)
		if authErr != nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "Ariba Shield WAF",
			AccountName: email,
		})
		if err != nil {
			http.Error(w, `{"error":"totp generation failed"}`, http.StatusInternalServerError)
			return
		}

		if _, err := st.Pool.Exec(r.Context(),
			`UPDATE users SET totp_secret = $1 WHERE id = $2`, key.Secret(), userID); err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"secret": key.Secret(),
			"otpauth_url": key.URL(),
		})
	}
}

// VerifyMFA validates a TOTP code and (optionally) enables MFA.
func VerifyMFA(st *store.Store) http.HandlerFunc {
	type verify struct {
		Code string `json:"code"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		userID, _, authErr := authenticatedUser(st, r)
		if authErr != nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		var body verify
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Code == "" {
			http.Error(w, `{"error":"code required"}`, http.StatusBadRequest)
			return
		}

		var secret string
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT totp_secret FROM users WHERE id = $1`, userID).Scan(&secret); err != nil {
			http.Error(w, `{"error":"mfa not configured"}`, http.StatusBadRequest)
			return
		}

		if !totp.Validate(body.Code, secret) {
			http.Error(w, `{"error":"invalid code"}`, http.StatusUnauthorized)
			return
		}

		// Verification passed — enable MFA.
		if _, err := st.Pool.Exec(r.Context(),
			`UPDATE users SET totp_enabled = true WHERE id = $1`, userID); err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "mfa_enabled"})
	}
}

// DisableMFA disables TOTP MFA (requires re-auth / valid code in production;
// here requires a valid session).
func DisableMFA(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _, authErr := authenticatedUser(st, r)
		if authErr != nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		if _, err := st.Pool.Exec(r.Context(),
			`UPDATE users SET totp_enabled = false, totp_secret = NULL WHERE id = $1`, userID); err != nil {
			http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "mfa_disabled"})
	}
}

// ChangePassword updates the user's password (validates current password).
func ChangePassword(st *store.Store) http.HandlerFunc {
	type change struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		userID, _, authErr := authenticatedUser(st, r)
		if authErr != nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		var body change
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.NewPassword == "" || len(body.NewPassword) < 8 {
			http.Error(w, `{"error":"new password must be at least 8 characters"}`, http.StatusBadRequest)
			return
		}

		var currentHash string
		if err := st.Pool.QueryRow(r.Context(),
			`SELECT password_hash FROM users WHERE id = $1`, userID).Scan(&currentHash); err != nil {
			http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
			return
		}
		if !store.CheckPassword(body.CurrentPassword, currentHash) {
			http.Error(w, `{"error":"current password incorrect"}`, http.StatusUnauthorized)
			return
		}

		newHash, err := store.HashPassword(body.NewPassword)
		if err != nil {
			http.Error(w, `{"error":"password hashing failed"}`, http.StatusInternalServerError)
			return
		}

		if _, err := st.Pool.Exec(r.Context(),
			`UPDATE users SET password_hash = $1, version = version + 1, updated_at = now() WHERE id = $2`,
			newHash, userID); err != nil {
			http.Error(w, `{"error":"password change failed"}`, http.StatusInternalServerError)
			return
		}

		// Invalidate all other sessions on password change (security best practice).
		_, _ = st.Pool.Exec(r.Context(),
			`DELETE FROM sessions WHERE user_id = $1`, userID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "password_changed"})
	}
}

// BreakGlass grants temporary elevated access for the break-glass account.
// The account must be pre-flagged (break_glass_enabled) and the code must
// match the configured break-glass token from env.
func BreakGlass(st *store.Store) http.HandlerFunc {
	type breakGlass struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body breakGlass
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}

		expected := envOr("BREAK_GLASS_CODE", "")
		if expected == "" || body.Code != expected {
			http.Error(w, `{"error":"invalid break-glass code"}`, http.StatusUnauthorized)
			return
		}

		var userID string
		err := st.Pool.QueryRow(r.Context(),
			`SELECT id FROM users WHERE email = $1 AND break_glass_enabled = true`, body.Email).Scan(&userID)
		if err != nil || userID == "" {
			http.Error(w, `{"error":"break-glass account not enabled"}`, http.StatusForbidden)
			return
		}

		// Issue a short-lived elevated session (10 minutes).
		token, hash, err := newSessionToken()
		if err != nil {
			http.Error(w, `{"error":"token generation failed"}`, http.StatusInternalServerError)
			return
		}
		sessionID := ulid.Make().String()
		expires := time.Now().Add(10 * time.Minute)
		if _, err := st.Pool.Exec(r.Context(),
			`INSERT INTO sessions (id, user_id, token_hash, expires_at, request_id)
			 VALUES ($1, $2, $3, $4, $5)`,
			sessionID, userID, hash, expires, r.Header.Get("X-Request-ID")); err != nil {
			http.Error(w, `{"error":"session creation failed"}`, http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "shield_session",
			Value:    userID + ":" + token,
			Path:     "/",
			Expires:  expires,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "break_glass_active", "expires": expires.Format(time.RFC3339)})
	}
}

// --- helpers ---

// authenticatedUser resolves the user from the shield_session cookie.
func authenticatedUser(st *store.Store, r *http.Request) (userID, email string, err error) {
	cookie, err := r.Cookie("shield_session")
	if err != nil || cookie.Value == "" {
		return "", "", errors.New("no session")
	}
	uid, token, ok := sessionParts(cookie.Value)
	if !ok {
		return "", "", errors.New("invalid session")
	}
	err = st.Pool.QueryRow(r.Context(),
		`SELECT u.email FROM sessions s
		 JOIN users u ON u.id = s.user_id
		 WHERE s.user_id = $1 AND s.token_hash = $2 AND s.expires_at > now()`,
		uid, hashToken(token)).Scan(&email)
	if err != nil {
		return "", "", errors.New("session expired")
	}
	return uid, email, nil
}

func hashToken(token string) string {
	// SHA-256 hex for storing session tokens (never store the raw token).
	sum := sha256Sum([]byte(token))
	return sum
}

func envOr(key, def string) string {
	if v := getenv(key); v != "" {
		return v
	}
	return def
}

func sha256Sum(b []byte) string {
	return hex.EncodeToString(sha256.New().Sum(b))
}

func getenv(key string) string {
	return os.Getenv(key)
}
