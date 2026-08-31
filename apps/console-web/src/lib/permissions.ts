export type NavItem = {
  href: string;
  labelKey: string;
  /** Roles that can see this item. Empty = visible to all authenticated. */
  roles?: string[];
};

/** Permission-aware navigation model. Mirrors backend RBAC roles. */
export const ROLE_NAV: Record<string, NavItem[]> = {
  monitor: [
    { href: "/", labelKey: "overview" },
    { href: "/security-events", labelKey: "security_events", roles: ["Super Admin", "Platform Admin", "Security Admin", "App Owner", "SOC Analyst", "Auditor", "Read Only"] },
    { href: "/incidents", labelKey: "incidents", roles: ["Super Admin", "Platform Admin", "Security Admin", "SOC Analyst"] },
  ],
protect: [
    { href: "/applications", labelKey: "applications" },
    { href: "/policies", labelKey: "policies", roles: ["Super Admin", "Platform Admin", "Security Admin", "App Owner"] },
    { href: "/rules", labelKey: "rules", roles: ["Super Admin", "Platform Admin", "Security Admin", "App Owner"] },
    { href: "/traffic-control", labelKey: "traffic", roles: ["Super Admin", "Platform Admin", "Security Admin"] },
  ],
  platform: [
    { href: "/gateways", labelKey: "gateways_clusters" },
    { href: "/integrations", labelKey: "integrations", roles: ["Super Admin", "Platform Admin", "Security Admin"] },
    { href: "/notification-channels", labelKey: "notification_channels", roles: ["Super Admin", "Platform Admin", "Security Admin"] },
    { href: "/reports", labelKey: "reports", roles: ["Super Admin", "Platform Admin", "Security Admin", "Auditor"] },
    { href: "/users", labelKey: "users_access", roles: ["Super Admin", "Platform Admin"] },
    { href: "/audit", labelKey: "audit_log", roles: ["Super Admin", "Platform Admin", "Auditor"] },
    { href: "/certificates", labelKey: "certificates", roles: ["Super Admin", "Platform Admin", "Security Admin"] },
    { href: "/backups", labelKey: "backups", roles: ["Super Admin", "Platform Admin"] },
    { href: "/threat-intelligence", labelKey: "threat_intelligence", roles: ["Super Admin", "Platform Admin", "Security Admin"] },
    { href: "/settings", labelKey: "system_settings", roles: ["Super Admin", "Platform Admin"] },
    { href: "/license", labelKey: "license", roles: ["Super Admin", "Platform Admin"] },
  ],
};

/** Returns the nav items a role is allowed to see. */
export function navItemsForRole(role: string): Record<string, NavItem[]> {
  const result: Record<string, NavItem[]> = {};
  for (const [section, items] of Object.entries(ROLE_NAV)) {
    result[section] = items.filter((item) => !item.roles || item.roles.includes(role));
  }
  return result;
}

/** Map of backend permission strings to frontend action capabilities. */
export const ROLE_PERMISSIONS: Record<string, string[]> = {
  "Super Admin": ["create", "edit", "delete", "activate", "deploy", "admin"],
  "Platform Admin": ["create", "edit", "delete", "activate", "deploy", "admin"],
  "Security Admin": ["create", "edit", "activate"],
  "App Owner": ["create", "edit"],
  "SOC Analyst": ["read"],
  Auditor: ["read"],
  "Read Only": ["read"],
};

/** Check if a role has a given action capability. */
export function can(role: string | undefined, action: string): boolean {
  if (!role) return false;
  return (ROLE_PERMISSIONS[role] ?? []).includes(action);
}

/** Roles that can read this page; empty = all authenticated. */
export const PAGE_ROLES: Record<string, string[]> = {
  "/users": ["Super Admin", "Platform Admin"],
  "/audit": ["Super Admin", "Platform Admin", "Auditor"],
  "/settings": ["Super Admin", "Platform Admin"],
  "/backups": ["Super Admin", "Platform Admin"],
  "/license": ["Super Admin", "Platform Admin"],
  "/integrations": ["Super Admin", "Platform Admin", "Security Admin"],
  "/notification-channels": ["Super Admin", "Platform Admin", "Security Admin"],
  "/policies": ["Super Admin", "Platform Admin", "Security Admin", "App Owner"],
  "/rules": ["Super Admin", "Platform Admin", "Security Admin", "App Owner"],
  "/traffic-control": ["Super Admin", "Platform Admin", "Security Admin"],
  "/certificates": ["Super Admin", "Platform Admin", "Security Admin"],
  "/threat-intelligence": ["Super Admin", "Platform Admin", "Security Admin"],
  "/incidents": ["Super Admin", "Platform Admin", "Security Admin", "SOC Analyst"],
};

/** Check if a role can access a given page path. */
export function canAccessPage(role: string | undefined, path: string): boolean {
  if (!role) return false;
  const allowed = PAGE_ROLES[path];
  if (!allowed) return true;
  return allowed.includes(role);
}
