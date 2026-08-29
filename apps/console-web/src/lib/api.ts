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

export function listApplications(): Promise<Application[]> {
  return request<Application[]>("/applications");
}

export function createApplication(name: string, description: string): Promise<{ id: string }> {
  return request<{ id: string }>("/applications", {
    method: "POST",
    body: JSON.stringify({ name, description }),
  });
}

export function listOrigins(appId: string): Promise<Origin[]> {
  return request<Origin[]>(`/applications/${appId}/origins`);
}

export function createOrigin(appId: string, origin: Omit<Origin, "id" | "enabled">): Promise<{ id: string }> {
  return request<{ id: string }>(`/applications/${appId}/origins`, {
    method: "POST",
    body: JSON.stringify(origin),
  });
}

export function listDomains(appId: string): Promise<Domain[]> {
  return request<Domain[]>(`/applications/${appId}/domains`);
}

export function listGateways(): Promise<Gateway[]> {
  return request<Gateway[]>("/gateways");
}

export function listSecurityPolicies(): Promise<SecurityPolicy[]> {
  return request<SecurityPolicy[]>("/security-policies");
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

export type User = {
  id: string;
  email: string;
  role: string;
  status: string;
  created_at: string;
};

export function listUsers(role?: string): Promise<User[]> {
  return request<{ users: User[] }>("/users", undefined, role).then((r) => r.users ?? []);
}