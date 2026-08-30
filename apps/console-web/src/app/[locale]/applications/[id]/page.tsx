"use client";

import { useCallback, useEffect, useState } from "react";
import { useLocale } from "next-intl";
import { useParams, useRouter } from "next/navigation";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import { StatusBadge } from "@/components/shared/Badges";
import DataTable, { type Column } from "@/components/shared/DataTable";
import { SeverityBadge } from "@/components/shared/Badges";
import {
  getApplication,
  listApplicationDomains,
  listApplicationOrigins,
  listApplicationEvents,
  listApplicationIncidents,
  listApplicationPolicies,
  getApplicationHealth,
  type Application,
  type Domain,
  type Origin,
} from "@/lib/api";

type Tab = "overview" | "domains" | "origins" | "events" | "incidents" | "policies" | "health";

type Incident = { id: string; title: string; severity: string; status: string };
type SecurityPolicy = { id: string; name: string; enforcement_mode: string };

export default function ApplicationDetailPage() {
  const locale = useLocale();
  const router = useRouter();
  const params = useParams<{ id: string }>();
  const [tab, setTab] = useState<Tab>("overview");
  const [app, setApp] = useState<Application | null>(null);
  const [domains, setDomains] = useState<Domain[]>([]);
  const [origins, setOrigins] = useState<Origin[]>([]);
  const [events, setEvents] = useState<unknown[]>([]);
  const [incidents, setIncidents] = useState<Incident[]>([]);
  const [policies, setPolicies] = useState<SecurityPolicy[]>([]);
  const [health, setHealth] = useState<Record<string, unknown>>({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const load = useCallback(async () => {
    if (!params?.id) return;
    setLoading(true);
    setError("");
    try {
      const [a, d, o, e, i, p, h] = await Promise.all([
        getApplication(params.id),
        listApplicationDomains(params.id),
        listApplicationOrigins(params.id),
        listApplicationEvents(params.id),
        listApplicationIncidents(params.id),
        listApplicationPolicies(params.id),
        getApplicationHealth(params.id),
      ]);
      setApp(a);
      setDomains(d);
      setOrigins(o);
      setEvents(e ?? []);
      setIncidents(i ?? []);
      setPolicies(p ?? []);
      setHealth(h ?? {});
    } catch {
      setError("Failed to load application");
    } finally {
      setLoading(false);
    }
  }, [params?.id]);

  useEffect(() => {
    load();
  }, [load]);

  if (loading) {
    return (
      <div className="dashboard-container">
        <Sidebar locale={locale} />
        <main className="main-content"><p style={{ color: "var(--text-secondary)" }}>Loading…</p></main>
      </div>
    );
  }

  if (error || !app) {
    return (
      <div className="dashboard-container">
        <Sidebar locale={locale} />
        <main className="main-content">
          <div className="glass-panel" style={{ padding: "40px", textAlign: "center" }}>
            <p style={{ color: "var(--danger)" }}>{error || "Application not found"}</p>
            <button type="button" className="btn btn-primary" onClick={load} style={{ marginTop: "12px" }}>Retry</button>
          </div>
        </main>
      </div>
    );
  }

  const tabs: { key: Tab; label: string }[] = [
    { key: "overview", label: "Overview" },
    { key: "domains", label: `Domains (${domains.length})` },
    { key: "origins", label: `Origins (${origins.length})` },
    { key: "events", label: `Events (${(events as unknown[]).length})` },
    { key: "incidents", label: `Incidents (${incidents.length})` },
    { key: "policies", label: `Policies (${policies.length})` },
    { key: "health", label: "Health" },
  ];

  const eventColumns: Column<Record<string, unknown>>[] = [
    { key: "severity", label: "Severity", render: (r) => <SeverityBadge value={String(r.severity ?? "")} /> },
    { key: "created_at", label: "Time", render: (r) => (r.created_at ? new Date(String(r.created_at)).toLocaleString() : "—") },
    { key: "reason", label: "Reason", render: (r) => String(r.reason ?? "—") },
    { key: "client_ip", label: "Client IP", render: (r) => String(r.client_ip ?? "—") },
  ];

  const incidentColumns: Column<Incident>[] = [
    { key: "title", label: "Title" },
    { key: "severity", label: "Severity", render: (r) => <SeverityBadge value={r.severity} /> },
    { key: "status", label: "Status", render: (r) => <StatusBadge value={r.status} /> },
  ];

  const policyColumns: Column<SecurityPolicy>[] = [
    { key: "name", label: "Name" },
    { key: "enforcement_mode", label: "Mode", render: (r) => <StatusBadge value={r.enforcement_mode} /> },
  ];

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />

      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <button type="button" onClick={() => router.push(`/${locale}/applications`)} style={{ background: "none", border: "none", color: "var(--text-secondary)", cursor: "pointer", padding: 0, marginBottom: "4px", fontSize: "13px" }}>
              ← Back to applications
            </button>
            <h1>{app.name}</h1>
            <p style={{ color: "var(--text-secondary)" }}>{app.description || app.id}</p>
          </div>
          <div className="header-actions" style={{ display: "flex", alignItems: "center", gap: "12px" }}>
            <StatusBadge value={app.status} />
            <UserProfileWidget />
          </div>
        </div>

        {/* Tabs */}
        <div style={{ display: "flex", gap: "8px", marginBottom: "20px", flexWrap: "wrap" }}>
          {tabs.map((t) => (
            <button
              key={t.key}
              type="button"
              onClick={() => setTab(t.key)}
              style={{
                padding: "10px 16px",
                borderRadius: "8px",
                border: `1px solid ${tab === t.key ? "var(--accent-primary)" : "rgba(255,255,255,0.1)"}`,
                background: tab === t.key ? "rgba(59,130,246,0.15)" : "rgba(255,255,255,0.03)",
                color: tab === t.key ? "#60a5fa" : "var(--text-secondary)",
                fontSize: "13px",
                fontWeight: tab === t.key ? 600 : 500,
                cursor: "pointer",
              }}
            >
              {t.label}
            </button>
          ))}
        </div>

        {tab === "overview" && (
          <div className="metrics-grid animate-fade-in delay-1" style={{ maxWidth: "900px" }}>
            <div className="glass-panel metric-card" style={{ padding: "20px" }}>
              <div className="metric-header"><span>Status</span></div>
              <div style={{ marginTop: "8px" }}><StatusBadge value={app.status} /></div>
            </div>
            <div className="glass-panel metric-card" style={{ padding: "20px" }}>
              <div className="metric-header"><span>Domains</span></div>
              <div className="metric-value" style={{ fontSize: "24px" }}>{domains.length}</div>
            </div>
            <div className="glass-panel metric-card" style={{ padding: "20px" }}>
              <div className="metric-header"><span>Origins</span></div>
              <div className="metric-value" style={{ fontSize: "24px" }}>{origins.length}</div>
            </div>
            <div className="glass-panel metric-card" style={{ padding: "20px" }}>
              <div className="metric-header"><span>Version</span></div>
              <div className="metric-value" style={{ fontSize: "24px" }}>{app.version}</div>
            </div>
          </div>
        )}

        {tab === "domains" && (
          <div className="data-section animate-fade-in delay-1">
            <div className="glass-panel" style={{ padding: "0", overflow: "hidden" }}>
              <table style={{ width: "100%", borderCollapse: "collapse", textAlign: "left" }}>
                <thead>
                  <tr style={{ borderBottom: "1px solid rgba(255,255,255,0.08)", background: "rgba(255,255,255,0.02)" }}>
                    <th style={{ padding: "14px 20px", color: "var(--text-secondary)", fontWeight: 500, fontSize: "13px" }}>Hostname</th>
                    <th style={{ padding: "14px 20px", color: "var(--text-secondary)", fontWeight: 500, fontSize: "13px" }}>Enabled</th>
                  </tr>
                </thead>
                <tbody>
                  {domains.map((d, idx) => (
                    <tr key={d.id} style={{ borderBottom: idx !== domains.length - 1 ? "1px solid rgba(255,255,255,0.05)" : "none" }}>
                      <td style={{ padding: "12px 20px", fontSize: "13px", fontFamily: "monospace" }}>{d.hostname}</td>
                      <td style={{ padding: "12px 20px", fontSize: "13px" }}>{d.enabled ? "✓" : "—"}</td>
                    </tr>
                  ))}
                  {domains.length === 0 && (
                    <tr><td colSpan={2} style={{ padding: "40px", textAlign: "center", color: "var(--text-secondary)" }}>No domains.</td></tr>
                  )}
                </tbody>
              </table>
            </div>
          </div>
        )}

        {tab === "origins" && (
          <div className="data-section animate-fade-in delay-1">
            <div className="glass-panel" style={{ padding: "0", overflow: "hidden" }}>
              <table style={{ width: "100%", borderCollapse: "collapse", textAlign: "left" }}>
                <thead>
                  <tr style={{ borderBottom: "1px solid rgba(255,255,255,0.08)", background: "rgba(255,255,255,0.02)" }}>
                    <th style={{ padding: "14px 20px", color: "var(--text-secondary)", fontWeight: 500, fontSize: "13px" }}>Name</th>
                    <th style={{ padding: "14px 20px", color: "var(--text-secondary)", fontWeight: 500, fontSize: "13px" }}>Protocol</th>
                    <th style={{ padding: "14px 20px", color: "var(--text-secondary)", fontWeight: 500, fontSize: "13px" }}>Host</th>
                    <th style={{ padding: "14px 20px", color: "var(--text-secondary)", fontWeight: 500, fontSize: "13px" }}>Port</th>
                    <th style={{ padding: "14px 20px", color: "var(--text-secondary)", fontWeight: 500, fontSize: "13px" }}>Enabled</th>
                  </tr>
                </thead>
                <tbody>
                  {origins.map((o, idx) => (
                    <tr key={o.id} style={{ borderBottom: idx !== origins.length - 1 ? "1px solid rgba(255,255,255,0.05)" : "none" }}>
                      <td style={{ padding: "12px 20px", fontSize: "13px" }}>{o.name}</td>
                      <td style={{ padding: "12px 20px", fontSize: "13px" }}>{o.protocol}</td>
                      <td style={{ padding: "12px 20px", fontSize: "13px", fontFamily: "monospace" }}>{o.host}</td>
                      <td style={{ padding: "12px 20px", fontSize: "13px" }}>{o.port}</td>
                      <td style={{ padding: "12px 20px", fontSize: "13px" }}>{o.enabled ? "✓" : "—"}</td>
                    </tr>
                  ))}
                  {origins.length === 0 && (
                    <tr><td colSpan={5} style={{ padding: "40px", textAlign: "center", color: "var(--text-secondary)" }}>No origins.</td></tr>
                  )}
                </tbody>
              </table>
            </div>
          </div>
        )}

        {tab === "events" && (
          <div className="data-section animate-fade-in delay-1">
            <DataTable columns={eventColumns} data={events as Record<string, unknown>[]} rowKey={(r) => String(r.id ?? r.event_id ?? "")} emptyMessage="No security events." />
          </div>
        )}

        {tab === "incidents" && (
          <div className="data-section animate-fade-in delay-1">
            <DataTable columns={incidentColumns} data={incidents} rowKey={(r) => r.id} emptyMessage="No incidents." />
          </div>
        )}

        {tab === "policies" && (
          <div className="data-section animate-fade-in delay-1">
            <DataTable columns={policyColumns} data={policies} rowKey={(r) => r.id} emptyMessage="No policies bound." />
          </div>
        )}

        {tab === "health" && (
          <div className="data-section animate-fade-in delay-1">
            <div className="glass-panel" style={{ padding: "24px" }}>
              <pre style={{ fontSize: "13px", color: "var(--text-secondary)", whiteSpace: "pre-wrap", wordBreak: "break-word" }}>
                {JSON.stringify(health, null, 2)}
              </pre>
            </div>
          </div>
        )}
      </main>
    </div>
  );
}
