package middleware

import (
	"context"
	"testing"
)

func TestRolePermissionsMapping(t *testing.T) {
	// Every role must map to a non-empty permission set.
	for role, perms := range RolePermissions {
		if len(perms) == 0 {
			t.Errorf("role %q has empty permissions", role)
		}
	}

	// Super Admin should have system:admin.
	if !contains(RolePermissions["Super Admin"], PermSystemAdmin) {
		t.Error("Super Admin should have system:admin")
	}
	// Read Only should NOT have app:write.
	if contains(RolePermissions["Read Only"], PermAppWrite) {
		t.Error("Read Only should not have app:write")
	}
	// Security Admin should have policy:write but not user:admin.
	if !contains(RolePermissions["Security Admin"], PermPolicyWrite) {
		t.Error("Security Admin should have policy:write")
	}
	if contains(RolePermissions["Security Admin"], PermUserAdmin) {
		t.Error("Security Admin should not have user:admin")
	}
}

func TestHasPermission(t *testing.T) {
	ctx := context.Background()
	rbac := &RBACContext{
		User:        RBACUser{ID: "u1"},
		Roles:       []string{"Security Admin"},
		Permissions: RolePermissions["Security Admin"],
	}
	ctx = ContextWithRBAC(ctx, rbac)

	if !HasPermission(ctx, PermPolicyWrite) {
		t.Error("expected policy:write")
	}
	if HasPermission(ctx, PermUserAdmin) {
		t.Error("did not expect user:admin")
	}
}

func TestNormalizeRole(t *testing.T) {
	cases := map[string]string{
		"super_admin":     "Super Admin",
		"SECURITY-ADMIN":  "Security Admin",
		"soc_analyst":     "SOC Analyst",
		"auditor":         "Auditor",
		"read-only":       "Read Only",
		"unknown-role":    "unknown-role",
	}
	for in, want := range cases {
		if got := NormalizeRole(in); got != want {
			t.Errorf("NormalizeRole(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestCheckRoutePermission(t *testing.T) {
	ctx := context.Background()

	// Read Only can read applications but not write them.
	readOnly := ContextWithRBAC(ctx, &RBACContext{
		Roles:       []string{"Read Only"},
		Permissions: RolePermissions["Read Only"],
	})
	if !CheckRoutePermission(readOnly, "GET", "/api/v1/applications") {
		t.Error("Read Only should GET applications")
	}
	if CheckRoutePermission(readOnly, "POST", "/api/v1/applications") {
		t.Error("Read Only should NOT POST applications")
	}

	// Health is always open.
	if !CheckRoutePermission(readOnly, "GET", "/api/v1/health") {
		t.Error("health should be open")
	}

	// Security Admin can create policies.
	secAdmin := ContextWithRBAC(ctx, &RBACContext{
		Roles:       []string{"Security Admin"},
		Permissions: RolePermissions["Security Admin"],
	})
	if !CheckRoutePermission(secAdmin, "POST", "/api/v1/security-policies") {
		t.Error("Security Admin should POST policies")
	}
	if CheckRoutePermission(secAdmin, "GET", "/api/v1/users") {
		t.Error("Security Admin should NOT read users")
	}
}

func contains(s []string, v string) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}