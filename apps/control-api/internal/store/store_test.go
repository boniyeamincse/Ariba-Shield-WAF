package store

import (
	"testing"
)

func TestHashAndCheckPassword(t *testing.T) {
	hash, err := HashPassword("correct-horse-battery-staple")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if hash == "correct-horse-battery-staple" {
		t.Fatal("password must not be stored in plaintext")
	}
	if !CheckPassword("correct-horse-battery-staple", hash) {
		t.Fatal("correct password should verify")
	}
	if CheckPassword("wrong", hash) {
		t.Fatal("wrong password should not verify")
	}
}

func TestDefaultRolesComplete(t *testing.T) {
	required := []string{
		"Super Admin", "Platform Admin", "Security Admin", "App Owner",
		"SOC Analyst", "Auditor", "Read Only",
	}
	for _, r := range required {
		if _, ok := DefaultRoles[r]; !ok {
			t.Errorf("missing role %q in DefaultRoles", r)
		}
	}
	// Super Admin must have system:admin; Read Only must not have write perms.
	if !containsStr(DefaultRoles["Super Admin"], "system:admin") {
		t.Error("Super Admin should have system:admin")
	}
	for _, p := range []string{"app:write", "policy:write", "ip:write", "user:admin"} {
		if containsStr(DefaultRoles["Read Only"], p) {
			t.Errorf("Read Only should not have %q", p)
		}
	}
}

func containsStr(s []string, v string) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}