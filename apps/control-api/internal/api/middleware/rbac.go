package middleware

import (
	"context"
	"strings"
)

type rbacKey string

const (
	rbacUserKey    rbacKey = "rbac_user"
	rbacRolesKey   rbacKey = "rbac_roles"
	rbacPermsKey   rbacKey = "rbac_permissions"
)

// RBACUser holds the authenticated user's identity.
type RBACUser struct {
	ID           string
	Email        string
	OrganizationID string
}

// RBACContext stores user + roles + permissions in the request context.
type RBACContext struct {
	User        RBACUser
	Roles       []string
	Permissions []string
}

// ContextWithRBAC stores the RBAC context.
func ContextWithRBAC(ctx context.Context, rbac *RBACContext) context.Context {
	ctx = context.WithValue(ctx, rbacUserKey, rbac.User)
	ctx = context.WithValue(ctx, rbacRolesKey, rbac.Roles)
	ctx = context.WithValue(ctx, rbacPermsKey, rbac.Permissions)
	return ctx
}

// UserFromContext returns the RBAC user from context.
func UserFromContext(ctx context.Context) (RBACUser, bool) {
	u, ok := ctx.Value(rbacUserKey).(RBACUser)
	return u, ok
}

// HasPermission checks if the request context contains a specific permission.
func HasPermission(ctx context.Context, perm string) bool {
	perms, ok := ctx.Value(rbacPermsKey).([]string)
	if !ok {
		return false
	}
	for _, p := range perms {
		if p == perm || p == "*" {
			return true
		}
	}
	return false
}

// HasRole checks if the request context contains a specific role.
func HasRole(ctx context.Context, role string) bool {
	roles, ok := ctx.Value(rbacRolesKey).([]string)
	if !ok {
		return false
	}
	for _, r := range roles {
		if strings.EqualFold(r, role) {
			return true
		}
	}
	return false
}

// Permissions for the Ariba Shield RBAC model (master plan §7.1).
// Scoped: read, write, admin scopes per resource domain.
const (
	PermAppRead          = "app:read"
	PermAppWrite         = "app:write"
	PermAppAdmin         = "app:admin"
	PermGatewayRead      = "gateway:read"
	PermGatewayWrite     = "gateway:write"
	PermPolicyRead       = "policy:read"
	PermPolicyWrite      = "policy:write"
	PermPolicyAdmin      = "policy:admin"
	PermEventRead        = "event:read"
	PermAuditRead        = "audit:read"
	PermUserRead         = "user:read"
	PermUserWrite        = "user:write"
	PermUserAdmin        = "user:admin"
	PermSystemAdmin      = "system:admin"
	PermCertificateRead  = "cert:read"
	PermCertificateWrite = "cert:write"
	PermIPListRead       = "ip:read"
	PermIPListWrite      = "ip:write"
	PermRateLimitRead    = "ratelimit:read"
	PermRateLimitWrite   = "ratelimit:write"
)

// RolePermissions maps roles to their permission set (master plan §7.1).
var RolePermissions = map[string][]string{
	"Super Admin":     {PermAppAdmin, PermGatewayWrite, PermPolicyAdmin, PermEventRead, PermAuditRead, PermUserAdmin, PermSystemAdmin, PermCertificateWrite, PermIPListWrite, PermRateLimitWrite},
	"Platform Admin":  {PermAppAdmin, PermGatewayWrite, PermPolicyAdmin, PermEventRead, PermAuditRead, PermUserRead, PermCertificateWrite, PermIPListWrite, PermRateLimitWrite},
	"Security Admin":  {PermAppWrite, PermPolicyWrite, PermEventRead, PermAuditRead, PermIPListWrite, PermRateLimitWrite, PermCertificateRead},
	"App Owner":       {PermAppRead, PermAppWrite, PermGatewayRead, PermPolicyRead, PermEventRead},
	"SOC Analyst":     {PermEventRead, PermAppRead, PermGatewayRead, PermPolicyRead},
	"Auditor":         {PermAuditRead, PermEventRead, PermAppRead, PermGatewayRead, PermPolicyRead},
	"Read Only":       {PermAppRead, PermGatewayRead, PermPolicyRead, PermEventRead},
}