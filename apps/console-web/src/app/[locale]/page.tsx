import { getTranslations } from "next-intl/server";
import {
  listApplications,
  listGateways,
  listSecurityPolicies,
  getDashboardOverview,
  getDashboardTraffic,
  getDashboardSecurity,
} from "@/lib/api";
import UserProfileWidget from "@/components/UserProfileWidget";

import Sidebar from "@/components/layout/Sidebar";

export default async function OverviewPage({ params }: { params: Promise<{ locale: string }> }) {
  const { locale } = await params;
  const t = await getTranslations("overview");

  let apps: { id: string; name: string; status: string }[] = [];
  let gateways: { id: string; hostname: string; status: string }[] = [];
  let policies: { id: string; name: string; enforcement_mode: string }[] = [];
  let overview: {
    period_days: number;
    total_events: number;
    blocked_events: number;
    total_requests: number;
    applications: number;
    gateways: number;
    active_incidents: number;
  } = {
    period_days: 7,
    total_events: 0,
    blocked_events: 0,
    total_requests: 0,
    applications: 0,
    gateways: 0,
    active_incidents: 0,
  };
  let fetchError = false;

  let traffic: {
    period_days: number;
    total_requests: number;
    avg_latency_ms: number;
    p99_latency_ms: number;
    by_status: { status: string; count: number }[];
  } = { period_days: 7, total_requests: 0, avg_latency_ms: 0, p99_latency_ms: 0, by_status: [] };

  let security = {
    period_days: 7, total_events: 0, blocked_events: 0, unique_ips: 0, by_severity: [] as { severity: string; count: number }[],
  };

  try {
    [apps, gateways, policies, overview, traffic, security] = await Promise.all([
      listApplications(),
      listGateways(),
      listSecurityPolicies(),
      getDashboardOverview(),
      getDashboardTraffic(),
      getDashboardSecurity(),
    ]);
  } catch {
    fetchError = true;
  }

  const activeGateways = gateways.filter((g) => g.status === "active" || g.status === "starting").length;

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />

      {/* Main Content */}
      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <h1>{t("title")}</h1>
            <p>{t("description")}</p>
          </div>
          <div className="header-actions" style={{ display: 'flex', alignItems: 'center', gap: '24px' }}>
            {fetchError && (
              <span style={{ fontSize: "12px", color: "var(--warning)", background: "rgba(245,158,11,0.1)", padding: "6px 12px", borderRadius: "8px", border: "1px solid rgba(245,158,11,0.2)" }}>
                ⚠ API offline
              </span>
            )}
            <UserProfileWidget />
            <div style={{ height: '24px', width: '1px', backgroundColor: 'rgba(255,255,255,0.1)' }}></div>
            <button className="btn btn-primary">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
              Deploy Config
            </button>
          </div>
        </div>

        {/* Metrics Grid — real data from /dashboard/overview */}
        <div className="metrics-grid">
          <div className="metric-card glass-panel animate-fade-in delay-1">
            <div className="metric-header">
              <span>Total Applications</span>
              <div className="icon-wrapper icon-blue">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polygon points="12 2 2 7 12 12 22 7 12 2"></polygon><polyline points="2 17 12 22 22 17"></polyline><polyline points="2 12 12 17 22 12"></polyline></svg>
              </div>
            </div>
            <div>
              <div className="metric-value">{overview.applications}</div>
              <div className="metric-trend trend-up">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polyline points="23 6 13.5 15.5 8.5 10.5 1 18"></polyline><polyline points="17 6 23 6 23 12"></polyline></svg>
                Active in last {overview.period_days}d
              </div>
            </div>
          </div>

          <div className="metric-card glass-panel animate-fade-in delay-2">
            <div className="metric-header">
              <span>Active Gateways</span>
              <div className="icon-wrapper icon-purple">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="2" y="2" width="20" height="8" rx="2" ry="2"></rect><rect x="2" y="14" width="20" height="8" rx="2" ry="2"></rect><line x1="6" y1="6" x2="6.01" y2="6"></line><line x1="6" y1="18" x2="6.01" y2="18"></line></svg>
              </div>
            </div>
            <div>
              <div className="metric-value">{overview.gateways} / {gateways.length}</div>
              <div className={`metric-trend ${activeGateways < gateways.length && gateways.length > 0 ? "trend-down" : "trend-up"}`}>
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polyline points="23 18 13.5 8.5 8.5 13.5 1 6"></polyline><polyline points="17 18 23 18 23 12"></polyline></svg>
                {activeGateways} online
              </div>
            </div>
          </div>

          <div className="metric-card glass-panel animate-fade-in delay-3">
            <div className="metric-header">
              <span>Enforced Policies</span>
              <div className="icon-wrapper icon-green">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path></svg>
              </div>
            </div>
            <div>
              <div className="metric-value">{policies.length}</div>
              <div className="metric-trend trend-up">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polyline points="23 6 13.5 15.5 8.5 10.5 1 18"></polyline><polyline points="17 6 23 6 23 12"></polyline></svg>
                Active policies
              </div>
            </div>
          </div>

          <div className="metric-card glass-panel animate-fade-in delay-3">
            <div className="metric-header">
              <span>Security Events</span>
              <div className="icon-wrapper" style={{ background: 'rgba(239,68,68,0.12)' }}>
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#f87171" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path></svg>
              </div>
            </div>
            <div>
              <div className="metric-value">{overview.total_events}</div>
              <div className="metric-trend trend-down">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polyline points="23 18 13.5 8.5 8.5 13.5 1 6"></polyline><polyline points="17 18 23 18 23 12"></polyline></svg>
                {overview.blocked_events} blocked
              </div>
            </div>
          </div>
        </div>

        {/* Traffic Widget */}
        <div className="data-section animate-fade-in delay-3">
          <div className="section-header">
            <h2>Traffic Summary</h2>
          </div>
          <div className="metrics-grid" style={{ marginBottom: 0 }}>
            <div className="metric-card glass-panel">
              <div className="metric-header"><span>Total Requests</span></div>
              <div className="metric-value">{traffic.total_requests.toLocaleString()}</div>
            </div>
            <div className="metric-card glass-panel">
              <div className="metric-header"><span>Avg Latency</span></div>
              <div className="metric-value">{traffic.avg_latency_ms.toFixed(1)}ms</div>
            </div>
            <div className="metric-card glass-panel">
              <div className="metric-header"><span>P99 Latency</span></div>
              <div className="metric-value">{traffic.p99_latency_ms.toFixed(1)}ms</div>
            </div>
            <div className="metric-card glass-panel">
              <div className="metric-header"><span>Status Breakdown</span></div>
              <div style={{ fontSize: "13px", color: "var(--text-secondary)", lineHeight: "1.8" }}>
                {traffic.by_status.slice(0, 5).map((s) => (
                  <div key={s.status} style={{ display: "flex", justifyContent: "space-between", gap: "12px" }}>
                    <span>{s.status}</span>
                    <span style={{ color: "var(--text-primary)" }}>{s.count.toLocaleString()}</span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>

        {/* Security Widget */}
        <div className="data-section animate-fade-in delay-3">
          <div className="section-header">
            <h2>Security Summary</h2>
          </div>
          <div className="metrics-grid" style={{ marginBottom: 0 }}>
            <div className="metric-card glass-panel">
              <div className="metric-header"><span>Total Events</span></div>
              <div className="metric-value">{security.total_events.toLocaleString()}</div>
            </div>
            <div className="metric-card glass-panel">
              <div className="metric-header"><span>Blocked</span></div>
              <div className="metric-value" style={{ color: "var(--danger)" }}>{security.blocked_events.toLocaleString()}</div>
            </div>
            <div className="metric-card glass-panel">
              <div className="metric-header"><span>Unique IPs</span></div>
              <div className="metric-value">{security.unique_ips.toLocaleString()}</div>
            </div>
            <div className="metric-card glass-panel">
              <div className="metric-header"><span>By Severity</span></div>
              <div style={{ fontSize: "13px", color: "var(--text-secondary)", lineHeight: "1.8" }}>
                {security.by_severity.slice(0, 5).map((s) => (
                  <div key={s.severity} style={{ display: "flex", justifyContent: "space-between", gap: "12px" }}>
                    <span style={{ textTransform: "capitalize" }}>{s.severity}</span>
                    <span style={{ color: "var(--text-primary)" }}>{s.count.toLocaleString()}</span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>

        {/* Lists Section */}
        <div className="data-section animate-fade-in delay-3">
          <div className="section-header">
            <h2>Protected Applications</h2>
          </div>
          <div className="data-list glass-panel">
            {apps.length === 0 && !fetchError && (
              <div style={{ padding: "40px", textAlign: "center", color: "var(--text-secondary)" }}>
                {t("no_applications")}
              </div>
            )}
            {apps.map((a) => (
              <div key={a.id} className="list-item">
                <div className="item-info">
                  <div className="icon-wrapper" style={{ background: 'rgba(255,255,255,0.05)', padding: '10px' }}>
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="2" y1="12" x2="22" y2="12"></line><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path></svg>
                  </div>
                  <div>
                    <div className="item-title">{a.name}</div>
                    <div className="item-subtitle">ID: {a.id} • WAF Enabled</div>
                  </div>
                </div>
                <span className={`badge ${a.status === 'active' ? 'badge-active' : a.status === 'warning' ? 'badge-warning' : 'badge-danger'}`}>
                  {a.status}
                </span>
              </div>
            ))}
          </div>
        </div>

      </main>
    </div>
  );
}
