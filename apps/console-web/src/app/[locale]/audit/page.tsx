import { getTranslations } from "next-intl/server";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import { API_BASE } from "@/lib/api";

type AuditEvent = {
  id: string;
  action: string;
  resource: string;
  resource_id: string;
  actor_user_id: string;
  ip: string;
  request_id: string;
  created_at: string;
};

async function fetchAuditEvents(): Promise<AuditEvent[]> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/audit-events`, { cache: "no-store", credentials: "include" });
    if (!res.ok) return [];
    return res.json();
  } catch {
    return [];
  }
}

export default async function AuditLogPage({ params }: { params: Promise<{ locale: string }> }) {
  const { locale } = await params;
  const t = await getTranslations("audit");

  let events: AuditEvent[] = [];
  let fetchError = false;
  try {
    events = await fetchAuditEvents();
  } catch {
    fetchError = true;
  }

  const actionColors: Record<string, { bg: string; text: string }> = {
    POST: { bg: "rgba(59,130,246,0.15)", text: "#60a5fa" },
    PUT: { bg: "rgba(168,85,247,0.15)", text: "#c084fc" },
    PATCH: { bg: "rgba(245,158,11,0.15)", text: "#fbbf24" },
    DELETE: { bg: "rgba(239,68,68,0.15)", text: "#f87171" },
    GET: { bg: "rgba(16,185,129,0.15)", text: "#34d399" },
  };

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />

      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <h1>{t("title")}</h1>
            <p style={{ color: "var(--text-secondary)" }}>{t("description")}</p>
          </div>
          <div className="header-actions" style={{ display: "flex", alignItems: "center", gap: "16px" }}>
            {fetchError && (
              <span style={{ fontSize: "12px", color: "var(--warning)", background: "rgba(245,158,11,0.1)", padding: "6px 12px", borderRadius: "8px", border: "1px solid rgba(245,158,11,0.2)" }}>
                ⚠ API offline
              </span>
            )}
            <UserProfileWidget />
          </div>
        </div>

        {/* Immutable audit log banner */}
        <div
          className="glass-panel animate-fade-in delay-1"
          style={{ padding: "16px 20px", marginBottom: "20px", display: "flex", alignItems: "center", gap: "12px" }}
        >
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="var(--success)" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ flexShrink: 0 }}>
            <rect x="3" y="11" width="18" height="11" rx="2"></rect>
            <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
          </svg>
          <div style={{ fontSize: "13px", color: "var(--text-secondary)" }}>
            {t("immutable_note")}
          </div>
        </div>

        {/* Audit table */}
        <div className="data-section animate-fade-in delay-2">
          <div className="glass-panel" style={{ padding: "0", overflow: "hidden" }}>
            <table style={{ width: "100%", borderCollapse: "collapse", textAlign: "left" }}>
              <thead>
                <tr style={{ borderBottom: "1px solid rgba(255,255,255,0.08)", background: "rgba(255,255,255,0.02)" }}>
                  <th style={{ padding: "14px 20px", color: "var(--text-secondary)", fontWeight: 500, fontSize: "13px", whiteSpace: "nowrap" }}>{t("time")}</th>
                  <th style={{ padding: "14px 20px", color: "var(--text-secondary)", fontWeight: 500, fontSize: "13px", whiteSpace: "nowrap" }}>{t("action")}</th>
                  <th style={{ padding: "14px 20px", color: "var(--text-secondary)", fontWeight: 500, fontSize: "13px", whiteSpace: "nowrap" }}>{t("resource")}</th>
                  <th style={{ padding: "14px 20px", color: "var(--text-secondary)", fontWeight: 500, fontSize: "13px", whiteSpace: "nowrap" }}>{t("actor")}</th>
                  <th style={{ padding: "14px 20px", color: "var(--text-secondary)", fontWeight: 500, fontSize: "13px", whiteSpace: "nowrap" }}>{t("ip")}</th>
                  <th style={{ padding: "14px 20px", color: "var(--text-secondary)", fontWeight: 500, fontSize: "13px", whiteSpace: "nowrap" }}>{t("request_id")}</th>
                </tr>
              </thead>
              <tbody>
                {events.map((ev, idx) => {
                  const colors = actionColors[ev.action] ?? { bg: "rgba(255,255,255,0.06)", text: "#9ca3af" };
                  const time = new Date(ev.created_at).toLocaleString();
                  return (
                    <tr key={ev.id} style={{ borderBottom: idx !== events.length - 1 ? "1px solid rgba(255,255,255,0.05)" : "none" }}>
                      <td style={{ padding: "12px 20px", color: "var(--text-secondary)", fontSize: "13px", whiteSpace: "nowrap" }}>{time}</td>
                      <td style={{ padding: "12px 20px", whiteSpace: "nowrap" }}>
                        <span style={{ padding: "4px 10px", borderRadius: "6px", background: colors.bg, color: colors.text, fontSize: "12px", fontWeight: 600 }}>
                          {ev.action}
                        </span>
                      </td>
                      <td style={{ padding: "12px 20px", fontSize: "13px" }}>
                        <div style={{ fontWeight: 500 }}>{ev.resource}</div>
                        <div style={{ color: "var(--text-secondary)", fontSize: "12px" }}>{ev.resource_id}</div>
                      </td>
                      <td style={{ padding: "12px 20px", fontSize: "13px", whiteSpace: "nowrap" }}>
                        {ev.actor_user_id ? (
                          <span style={{ fontFamily: "monospace", fontSize: "12px" }}>{ev.actor_user_id.slice(0, 12)}…</span>
                        ) : (
                          <span style={{ color: "var(--text-secondary)" }}>system</span>
                        )}
                      </td>
                      <td style={{ padding: "12px 20px", color: "var(--text-secondary)", fontSize: "13px", whiteSpace: "nowrap" }}>{ev.ip || "—"}</td>
                      <td style={{ padding: "12px 20px", fontFamily: "monospace", fontSize: "12px", color: "var(--text-secondary)", whiteSpace: "nowrap" }}>
                        {ev.request_id ? ev.request_id.slice(0, 12) + "…" : "—"}
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>

            {events.length === 0 && !fetchError && (
              <div style={{ padding: "60px", textAlign: "center", color: "var(--text-secondary)" }}>
                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" style={{ marginBottom: "16px", opacity: 0.4 }} strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
                <p>{t("no_events")}</p>
              </div>
            )}
          </div>
        </div>
      </main>
    </div>
  );
}