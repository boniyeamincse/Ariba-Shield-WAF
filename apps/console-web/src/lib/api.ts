export type Application = {
  id: string;
  name: string;
  description: string;
  status: string;
  owner_user_id?: string;
  version: number;
};

export type Origin = {
  id: string;
  name: string;
  protocol: string;
  host: string;
  port: number;
  weight: number;
  enabled: boolean;
};

export type Domain = {
  id: string;
  hostname: string;
  enabled: boolean;
};

export type Gateway = {
  id: string;
  hostname: string;
  ip: string;
  version: string;
  status: string;
  last_seen_at?: string;
  applied_hash?: string;
};

export type SecurityPolicy = {
  id: string;
  name: string;
  description: string;
  enforcement_mode: string;
  application_id?: string;
  version: number;
};

export type User = {
  id: string;
  email: string;
  role: string;
  status: string;
  created_at: string;
};

export type SecurityEvent = {
  id: string;
  event_id: string;
  request_id: string;
  timestamp: string;
  severity: string;
  decision_action: string;
  reason: string;
  rule_ids: string[];
  client_ip: string;
  method: string;
  path: string;
  status: number;
  created_at: string;
};

export type Incident = {
  id: string;
  title: string;
  severity: string;
  status: string;
  owner_user_id?: string;
  notes: string;
  related_events: string[];
  created_at: string;
  updated_at: string;
};

export type AuditEvent = {
  id: string;
  action: string;
  resource: string;
  resource_id: string;
  actor_user_id: string;
  ip: string;
  request_id: string;
  created_at: string;
};

export type DashboardOverview = {
  period_days: number;
  total_events: number;
  blocked_events: number;
  total_requests: number;
  applications: number;
  gateways: number;
  active_incidents: number;
};

export type Report = {
  id: string;
  name: string;
  kind: string;
  status: string;
  created_by: string;
  created_at: string;
  updated_at: string;
};

export type NotificationChannel = {
  id: string;
  name: string;
  kind: string;
  is_default: boolean;
  enabled: boolean;
  created_at: string;
  updated_at: string;
};

export type License = {
  id: string;
  license_key: string;
  product: string;
  edition: string;
  status: string;
  seats: number;
  max_gateways: number;
  max_applications: number;
  issued_at: string;
  expires_at: string;
  activated_at: string;
};

export type Group = {
  id: string;
  name: string;
  organization_id: string;
  created_at: string;
};

export type Role = {
  id: string;
  name: string;
  permissions: string[];
};

export type IPList = {
  id: string;
  name: string;
  description: string;
  type: string;
  status: string;
};

export type RateLimit = {
  id: string;
  name: string;
  limit_count: number;
  window_seconds: number;
  action: string;
  status: string;
};

export type LearningSession = {
  id: string;
  organization_id: string;
  name: string;
  source: string;
  description: string;
  confidence_threshold: string;
  status: string;
  created_by: string;
  created_at: string;
  updated_at: string;
};

export type LearningSuggestion = {
  id: string;
  session_id: string;
  application_id: string;
  rule_id: string;
  severity: string;
  confidence: string;
  rationale: string;
  status: string;
  applied_at: string;
  created_at: string;
  updated_at: string;
};

export const API_BASE = process.env.NEXT_PUBLIC_API_BASE ?? "http://127.0.0.1:8443";

async function request<T>(path: string, init?: RequestInit, role?: string): Promise<T> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(init?.headers as Record<string, string> ?? {}),
  };
  // Forward role for mock auth RBAC enforcement
  if (role) headers["X-User-Role"] = role;

  const res = await fetch(`${API_BASE}/api/v1${path}`, {
    ...init,
    headers,
    credentials: "include",
    cache: "no-store",
  });
  if (!res.ok) {
    throw new Error(`API ${res.status}: ${await res.text()}`);
  }
  return res.json() as Promise<T>;
}

// ===== Applications =====

export function listApplications(): Promise<Application[]> {
  return request<Application[]>("/applications");
}

export function createApplication(name: string, description: string): Promise<{ id: string }> {
  return request<{ id: string }>("/applications", {
    method: "POST",
    body: JSON.stringify({ name, description }),
  });
}

export function getApplication(id: string): Promise<Application> {
  return request<Application>(`/applications/${id}`);
}

export function updateApplication(id: string, patch: Partial<Application>): Promise<{ id: string }> {
  return request<{ id: string }>(`/applications/${id}`, {
    method: "PATCH",
    body: JSON.stringify(patch),
  });
}

export function deleteApplication(id: string): Promise<{ status: string }> {
  return request<{ status: string }>(`/applications/${id}`, { method: "DELETE" });
}

export function listApplicationTraffic(id: string): Promise<unknown[]> {
  return request<unknown[]>(`/applications/${id}/traffic`);
}

export function listApplicationEvents(id: string): Promise<SecurityEvent[]> {
  return request<SecurityEvent[]>(`/applications/${id}/events`);
}

export function listApplicationIncidents(id: string): Promise<Incident[]> {
  return request<Incident[]>(`/applications/${id}/incidents`);
}

export function listApplicationPolicies(id: string): Promise<SecurityPolicy[]> {
  return request<SecurityPolicy[]>(`/applications/${id}/policies`);
}

export function getApplicationHealth(id: string): Promise<Record<string, unknown>> {
  return request<Record<string, unknown>>(`/applications/${id}/health`);
}

// ===== Origins =====

export function listOrigins(appId: string): Promise<Origin[]> {
  return request<Origin[]>(`/applications/${appId}/origins`);
}

export function createOrigin(appId: string, origin: Omit<Origin, "id" | "enabled">): Promise<{ id: string }> {
  return request<{ id: string }>(`/applications/${appId}/origins`, {
    method: "POST",
    body: JSON.stringify(origin),
  });
}

// ===== Domains =====

export function listDomains(appId: string): Promise<Domain[]> {
  return request<Domain[]>(`/applications/${appId}/domains`);
}

export function createDomain(appId: string, hostname: string): Promise<{ id: string }> {
  return request<{ id: string }>(`/applications/${appId}/domains`, {
    method: "POST",
    body: JSON.stringify({ hostname }),
  });
}

// ===== Gateways =====

export function listGateways(): Promise<Gateway[]> {
  return request<Gateway[]>("/gateways");
}

export function getGateway(id: string): Promise<Gateway> {
  return request<Gateway>(`/gateways/${id}`);
}

export function updateGateway(id: string, patch: Partial<Gateway>): Promise<{ id: string }> {
  return request<{ id: string }>(`/gateways/${id}`, { method: "PATCH", body: JSON.stringify(patch) });
}

export function getGatewayStatus(id: string): Promise<Record<string, unknown>> {
  return request<Record<string, unknown>>(`/gateways/${id}/status`);
}

export function getGatewayMetrics(id: string): Promise<Record<string, unknown>> {
  return request<Record<string, unknown>>(`/gateways/${id}/metrics`);
}

export function getGatewayConfig(id: string): Promise<Record<string, unknown>> {
  return request<Record<string, unknown>>(`/gateways/${id}/config`);
}

export function deleteGateway(id: string): Promise<{ status: string }> {
  return request<{ status: string }>(`/gateways/${id}`, { method: "DELETE" });
}

// ===== Security Policies =====

export function listSecurityPolicies(): Promise<SecurityPolicy[]> {
  return request<SecurityPolicy[]>("/security-policies");
}

export function getSecurityPolicy(id: string): Promise<SecurityPolicy> {
  return request<SecurityPolicy>(`/security-policies/${id}`);
}

export function createSecurityPolicy(
  name: string,
  enforcementMode: string,
  applicationId?: string
): Promise<{ id: string }> {
  return request<{ id: string }>("/security-policies", {
    method: "POST",
    body: JSON.stringify({ name, enforcement_mode: enforcementMode, application_id: applicationId }),
  });
}

export function updateSecurityPolicy(id: string, patch: Partial<SecurityPolicy>): Promise<{ id: string }> {
  return request<{ id: string }>(`/security-policies/${id}`, {
    method: "PATCH",
    body: JSON.stringify(patch),
  });
}

export function deleteSecurityPolicy(id: string): Promise<{ status: string }> {
  return request<{ status: string }>(`/security-policies/${id}`, { method: "DELETE" });
}

export function activateSecurityPolicy(id: string): Promise<{ id: string; status: string }> {
  return request<{ id: string; status: string }>(`/security-policies/${id}/activate`, { method: "POST" });
}

export function disableSecurityPolicy(id: string): Promise<{ id: string; status: string }> {
  return request<{ id: string; status: string }>(`/security-policies/${id}/disable`, { method: "POST" });
}

export function rollbackSecurityPolicy(id: string): Promise<{ id: string; status: string }> {
  return request<{ id: string; status: string }>(`/security-policies/${id}/rollback`, { method: "POST" });
}

export function cloneSecurityPolicy(id: string): Promise<{ id: string }> {
  return request<{ id: string }>(`/security-policies/${id}/clone`, { method: "POST" });
}

export function listPolicyVersions(id: string): Promise<unknown[]> {
  return request<unknown[]>(`/security-policies/${id}/versions`);
}

// ===== Security Events =====

export function listSecurityEvents(params?: Record<string, string>): Promise<{ events: SecurityEvent[]; pagination: { limit: number; offset: number; count: number } }> {
  const qs = new URLSearchParams(params ?? {}).toString();
  return request<{ events: SecurityEvent[]; pagination: { limit: number; offset: number; count: number } }>(`/security-events${qs ? `?${qs}` : ""}`);
}

export function getSecurityEvent(id: string): Promise<SecurityEvent> {
  return request<SecurityEvent>(`/security-events/${id}`);
}

export function getEventMatches(id: string): Promise<unknown> {
  return request<unknown>(`/security-events/${id}/matches`);
}

export function getEventTimeline(id: string): Promise<unknown> {
  return request<unknown>(`/security-events/${id}/timeline`);
}

export function maskSecurityEvent(id: string): Promise<{ id: string }> {
  return request<{ id: string }>(`/security-events/${id}/mask`, { method: "POST" });
}

export function exportSecurityEvent(id: string): Promise<{ id: string }> {
  return request<{ id: string }>(`/security-events/${id}/export`, { method: "POST" });
}

// ===== Incidents =====

export function listIncidents(): Promise<Incident[]> {
  return request<Incident[]>("/incidents");
}

export function createIncident(payload: Partial<Incident>): Promise<{ id: string }> {
  return request<{ id: string }>("/incidents", { method: "POST", body: JSON.stringify(payload) });
}

export function getIncident(id: string): Promise<Incident> {
  return request<Incident>(`/incidents/${id}`);
}

export function updateIncident(id: string, patch: Partial<Incident>): Promise<{ id: string }> {
  return request<{ id: string }>(`/incidents/${id}`, { method: "PATCH", body: JSON.stringify(patch) });
}

export function deleteIncident(id: string): Promise<{ status: string }> {
  return request<{ status: string }>(`/incidents/${id}`, { method: "DELETE" });
}

export function assignIncident(id: string, ownerUserId: string): Promise<{ id: string }> {
  return request<{ id: string }>(`/incidents/${id}/assign`, { method: "POST", body: JSON.stringify({ owner_user_id: ownerUserId }) });
}

export function escalateIncident(id: string): Promise<{ id: string }> {
  return request<{ id: string }>(`/incidents/${id}/escalate`, { method: "POST" });
}

export function closeIncident(id: string): Promise<{ id: string }> {
  return request<{ id: string }>(`/incidents/${id}/close`, { method: "POST" });
}

export function reopenIncident(id: string): Promise<{ id: string }> {
  return request<{ id: string }>(`/incidents/${id}/reopen`, { method: "POST" });
}

export function listIncidentEvents(id: string): Promise<unknown> {
  return request<unknown>(`/incidents/${id}/events`);
}

export function getIncidentTimeline(id: string): Promise<unknown> {
  return request<unknown>(`/incidents/${id}/timeline`);
}

// ===== Audit Events =====

export function listAuditEvents(): Promise<AuditEvent[]> {
  return request<AuditEvent[]>("/audit-events");
}

export function getAuditEvent(id: string): Promise<AuditEvent> {
  return request<AuditEvent>(`/audit-events/${id}`);
}

export async function exportAuditEvents(): Promise<Blob> {
  const res = await fetch(`${API_BASE}/api/v1/audit-events/export`, { credentials: "include", cache: "no-store" });
  if (!res.ok) throw new Error(`API ${res.status}`);
  return res.blob();
}

// ===== Users / Groups / Roles =====

export function listUsers(role?: string): Promise<User[]> {
  return request<{ users: User[] }>("/users", undefined, role).then((r) => r.users ?? []);
}

export function createUser(
  email: string,
  password: string,
  role: string,
  callerRole?: string
): Promise<{ id: string; email: string; role: string; status: string }> {
  return request<{ id: string; email: string; role: string; status: string }>(
    "/users",
    { method: "POST", body: JSON.stringify({ email, password, role }) },
    callerRole
  );
}

export function updateUser(id: string, patch: Partial<Pick<User, "role" | "status">>): Promise<{ id: string }> {
  return request<{ id: string }>(`/users/${id}`, { method: "PATCH", body: JSON.stringify(patch) });
}

export function deleteUser(id: string): Promise<{ status: string }> {
  return request<{ status: string }>(`/users/${id}`, { method: "DELETE" });
}

export function assignRoleToUser(id: string, roleId: string): Promise<{ id: string; role_id: string }> {
  return request<{ id: string; role_id: string }>(`/users/${id}/roles`, { method: "POST", body: JSON.stringify({ role_id: roleId }) });
}

export function removeRoleFromUser(id: string, roleId: string): Promise<{ id: string; role_id: string }> {
  return request<{ id: string; role_id: string }>(`/users/${id}/roles/${roleId}`, { method: "DELETE", body: JSON.stringify({ role_id: roleId }) });
}

export function listGroups(): Promise<{ groups: Group[] }> {
  return request<{ groups: Group[] }>("/groups");
}

export function createGroup(name: string): Promise<{ id: string }> {
  return request<{ id: string }>("/groups", { method: "POST", body: JSON.stringify({ name }) });
}

export function deleteGroup(id: string): Promise<{ status: string }> {
  return request<{ status: string }>(`/groups/${id}`, { method: "DELETE" });
}

export function listRoles(): Promise<{ roles: Role[] }> {
  return request<{ roles: Role[] }>("/roles");
}

// ===== Dashboard =====

export function getDashboardOverview(): Promise<DashboardOverview> {
  return request<DashboardOverview>("/dashboard/overview");
}

export type DashboardTraffic = {
  period_days: number;
  total_requests: number;
  avg_latency_ms: number;
  p99_latency_ms: number;
  by_status: { status: string; count: number }[];
};

export function getDashboardTraffic(): Promise<DashboardTraffic> {
  return request<DashboardTraffic>("/dashboard/traffic");
}

export type DashboardSecurity = {
  period_days: number;
  total_events: number;
  blocked_events: number;
  unique_ips: number;
  by_severity: { severity: string; count: number }[];
};

export function getDashboardSecurity(): Promise<DashboardSecurity> {
  return request<DashboardSecurity>("/dashboard/security");
}

export function getDashboardAttacks(): Promise<Record<string, unknown>> {
  return request<Record<string, unknown>>("/dashboard/attacks");
}

export function getDashboardTopIPs(): Promise<{ period_days: number; top_ips: { client_ip: string; hits: number; blocked: number }[] }> {
  return request<{ period_days: number; top_ips: { client_ip: string; hits: number; blocked: number }[] }>("/dashboard/top-ips");
}

export function getDashboardTopRules(): Promise<{ period_days: number; top_rules: { rule_id: string; hits: number }[] }> {
  return request<{ period_days: number; top_rules: { rule_id: string; hits: number }[] }>("/dashboard/top-rules");
}

export function getDashboardGateways(): Promise<{ gateways: unknown[]; total: number; active: number; offline: number }> {
  return request<{ gateways: unknown[]; total: number; active: number; offline: number }>("/dashboard/gateways");
}

export function getDashboardApplications(): Promise<{ period_days: number; applications: { id: string; name: string; status: string; requests: number; events: number; blocked: number }[] }> {
  return request<{ period_days: number; applications: { id: string; name: string; status: string; requests: number; events: number; blocked: number }[] }>("/dashboard/applications");
}

// ===== Reports =====

export function listReports(): Promise<Report[]> {
  return request<Report[]>("/reports");
}

export function getReport(id: string): Promise<Report & { summary: unknown }> {
  return request<Report & { summary: unknown }>(`/reports/${id}`);
}

export function deleteReport(id: string): Promise<{ status: string }> {
  return request<{ status: string }>(`/reports/${id}`, { method: "DELETE" });
}

export function generateReport(kind: "security" | "traffic" | "incidents" | "compliance"): Promise<{ id: string; kind: string; status: string }> {
  return request<{ id: string; kind: string; status: string }>(`/reports/${kind}`, { method: "POST" });
}

// ===== Notification Channels =====

export function listNotificationChannels(): Promise<NotificationChannel[]> {
  return request<NotificationChannel[]>("/notification-channels");
}

export function createNotificationChannel(payload: { name: string; kind: string; is_default?: boolean }): Promise<{ id: string }> {
  return request<{ id: string }>("/notification-channels", { method: "POST", body: JSON.stringify(payload) });
}

export function deleteNotificationChannel(id: string): Promise<{ status: string }> {
  return request<{ status: string }>(`/notification-channels/${id}`, { method: "DELETE" });
}

export function testNotificationChannel(id: string): Promise<{ id: string; success: string }> {
  return request<{ id: string; success: string }>(`/notification-channels/${id}/test`, { method: "POST" });
}

// ===== Settings =====

export function getSettings(): Promise<Record<string, Record<string, unknown>>> {
  return request<Record<string, Record<string, unknown>>>("/settings");
}

export function updateSettings(category: string, values: Record<string, unknown>): Promise<{ status: string }> {
  return request<{ status: string }>(`/settings${category ? `/${category}` : ""}`, { method: "PATCH", body: JSON.stringify(values) });
}

// ===== License =====

export function getLicense(): Promise<License | Record<string, unknown>> {
  return request<License | Record<string, unknown>>("/license");
}

export function activateLicense(licenseKey: string): Promise<Record<string, unknown>> {
  return request<Record<string, unknown>>("/license/activate", { method: "POST", body: JSON.stringify({ license_key: licenseKey }) });
}

export function deactivateLicense(): Promise<{ status: string }> {
  return request<{ status: string }>("/license/deactivate", { method: "POST" });
}

export function getLicenseUsage(): Promise<Record<string, unknown>> {
  return request<Record<string, unknown>>("/license/usage");
}

export function getLicenseEntitlements(): Promise<Record<string, unknown>> {
  return request<Record<string, unknown>>("/license/entitlements");
}

// ===== IP Lists & Rate Limits =====

export function listIPLists(): Promise<IPList[]> {
  return request<IPList[]>("/ip-lists");
}

export function createIPList(payload: Partial<IPList>): Promise<{ id: string }> {
  return request<{ id: string }>("/ip-lists", { method: "POST", body: JSON.stringify(payload) });
}

export function deleteIPList(id: string): Promise<{ status: string }> {
  return request<{ status: string }>(`/ip-lists/${id}`, { method: "DELETE" });
}

export function listRateLimits(): Promise<RateLimit[]> {
  return request<RateLimit[]>("/rate-limits");
}

export function createRateLimit(payload: Partial<RateLimit>): Promise<{ id: string }> {
  return request<{ id: string }>("/rate-limits", { method: "POST", body: JSON.stringify(payload) });
}

export function deleteRateLimit(id: string): Promise<{ status: string }> {
  return request<{ status: string }>(`/rate-limits/${id}`, { method: "DELETE" });
}

// ===== Certificates =====

export function listCertificates(): Promise<unknown[]> {
  return request<unknown[]>("/certificates");
}

export function deleteCertificate(id: string): Promise<{ status: string }> {
  return request<{ status: string }>(`/certificates/${id}`, { method: "DELETE" });
}

// ===== Backups =====

export function listBackups(): Promise<unknown[]> {
  return request<unknown[]>("/backups");
}

export function createBackup(): Promise<{ id: string }> {
  return request<{ id: string }>("/backups", { method: "POST" });
}

export function restoreBackup(id: string): Promise<{ id: string }> {
  return request<{ id: string }>(`/backups/${id}/restore`, { method: "POST" });
}

// ===== Learning =====

export function listLearningSessions(): Promise<LearningSession[]> {
  return request<LearningSession[]>("/learning/sessions");
}

export function createLearningSession(payload: { name: string; source: string; description?: string; confidence_threshold?: string }): Promise<{ id: string }> {
  return request<{ id: string }>("/learning/sessions", { method: "POST", body: JSON.stringify(payload) });
}

export function startLearningSession(id: string): Promise<{ id: string }> {
  return request<{ id: string }>(`/learning/sessions/${id}/start`, { method: "POST", body: JSON.stringify({ status: "start" }) });
}

export function stopLearningSession(id: string): Promise<{ id: string }> {
  return request<{ id: string }>(`/learning/sessions/${id}/stop`, { method: "POST", body: JSON.stringify({ status: "stop" }) });
}

export function listLearningSuggestions(): Promise<LearningSuggestion[]> {
  return request<LearningSuggestion[]>("/learning/suggestions");
}

export function acceptSuggestion(id: string): Promise<{ id: string }> {
  return request<{ id: string }>(`/learning/suggestions/${id}/accept`, { method: "POST", body: JSON.stringify({ accepted: true }) });
}

export function rejectSuggestion(id: string): Promise<{ id: string }> {
  return request<{ id: string }>(`/learning/suggestions/${id}/reject`, { method: "POST", body: JSON.stringify({ accepted: false }) });
}

// ===== Integrations =====

export function listIntegrations(): Promise<unknown[]> {
  return request<unknown[]>("/integrations");
}

export function testIntegration(id: string): Promise<{ id: string; success: string }> {
  return request<{ id: string; success: string }>(`/integrations/${id}/test`, { method: "POST" });
}

export function setIntegrationEnabled(id: string, enabled: boolean): Promise<{ id: string }> {
  return request<{ id: string }>(`/integrations/${id}/${enabled ? "enable" : "disable"}`, { method: "POST", body: JSON.stringify({ enabled }) });
}

// ===== Analytics =====

export function getSecurityAnalytics(): Promise<Record<string, unknown>> {
  return request<Record<string, unknown>>("/analytics");
}

export function getRuleAnalytics(): Promise<Record<string, unknown>> {
  return request<Record<string, unknown>>("/rule-analytics");
}
