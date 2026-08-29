package handlers

import (
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
)

func TestSessionParts(t *testing.T) {
	userID, token, ok := sessionParts("user123:abcToken")
	if !ok || userID != "user123" || token != "abcToken" {
		t.Fatalf("expected user123/abcToken, got %s/%s ok=%v", userID, token, ok)
	}
	if _, _, ok := sessionParts("no-colon"); ok {
		t.Fatal("expected no-colon to fail")
	}
	if _, _, ok := sessionParts(""); ok {
		t.Fatal("expected empty to fail")
	}
}

func TestNewSessionToken(t *testing.T) {
	t1, h1, err := newSessionToken()
	if err != nil {
		t.Fatal(err)
	}
	t2, h2, _ := newSessionToken()
	if t1 == t2 {
		t.Fatal("tokens must be unique")
	}
	if len(t1) != 64 {
		t.Fatalf("expected 64 hex chars, got %d", len(t1))
	}
	if h1 == t1 || h2 == t2 {
		t.Fatal("hash must differ from raw token")
	}
	if h1 == h2 {
		t.Fatal("hashes must differ")
	}
}

func TestHashTokenStable(t *testing.T) {
	if hashToken("abc") != hashToken("abc") {
		t.Fatal("hash must be deterministic")
	}
	if hashToken("abc") == hashToken("abd") {
		t.Fatal("different inputs must hash differently")
	}
}

func TestTOTPValidation(t *testing.T) {
	// Generate a key, get a valid code, and confirm validation passes.
	key, err := totp.Generate(totp.GenerateOpts{Issuer: "test", AccountName: "u@test.local"})
	if err != nil {
		t.Fatal(err)
	}
	code, err := totp.GenerateCode(key.Secret(), time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if !totp.Validate(code, key.Secret()) {
		t.Fatal("valid code should validate")
	}
	if totp.Validate("000000", key.Secret()) {
		t.Fatal("invalid code should not validate")
	}
}
