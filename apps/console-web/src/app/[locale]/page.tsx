"use client";

import { useEffect, useState } from "react";
import { useLocale, useTranslations } from "next-intl";
import {
  listApplications,
  listGateways,
  listSecurityPolicies,
  getDashboardOverview,
  getDashboardTraffic,
  getDashboardSecurity,
  getDashboardTopIPs,
  getDashboardTopRules,
} from "@/lib/api";
import UserProfileWidget from "@/components/UserProfileWidget";
import Sidebar from "@/components/layout/Sidebar";

export default function OverviewPage() {
  const locale = useLocale();
  const t = useTranslations("overview");

  const [apps, setApps] = useState<{ id: string; name: string; status: string }[]>([]);
  const [gateways, setGateways] = useState<{ id: string; hostname: string; status: string }[]>([]);
  const [policies, setPolicies] = useState<{ id: string; name: string; enforcement_mode: string }[]>([]);
  const [overview, setOverview] = useState({
    period_days: 7,
    total_events: 0,
    blocked_events: 0,
    total_requests: 0,
    applications: 0,
    gateways: 0,
    active_incidents: 0,
  });
  const [traffic, setTraffic] = useState({ period_days: 7, total_requests: 0, avg_latency_ms: 0, p99_latency_ms: 0, by_status: [] as { status: string; count: number }[] });
  const [security, setSecurity] = useState({ period_days: 7, total_events: 0, blocked_events: 0, unique_ips: 0, by_severity: [] as { severity: string; count: number }[] });
  const [topIPs, setTopIPs] = useState<{ client_ip: string; hits: number; blocked: number }[]>([]);
  const [topRules, setTopRules] = useState<{ rule_id: string; hits: number }[]>([]);
  
  const [loading, setLoading] = useState(true);
  const [fetchError, setFetchError] = useState(false);

  useEffect(() => {
    async function loadData() {
      try {
        const [a, g, p, o, tr, sec, tip, trul] = await Promise.all([
          listApplications(),
          listGateways(),
          listSecurityPolicies(),
          getDashboardOverview(),
          getDashboardTraffic(),
          getDashboardSecurity(),
          getDashboardTopIPs().then((r) => r.top_ips ?? []),
          getDashboardTopRules().then((r) => r.top_rules ?? []),
        ]);
        setApps(a);
        setGateways(g);
        setPolicies(p);
        setOverview(o);
        setTraffic(tr);
        setSecurity(sec);
        setTopIPs(tip);
        setTopRules(trul);
        setFetchError(false);
      } catch {
        // Ignored to prevent eslint error, handled via state
        setFetchError(true);
      } finally {
        setLoading(false);
      }
    }
    loadData();
  }, []);

  const activeGateways = gateways.filter((g) => g.status === "active" || g.status === "starting").length;

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />

      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <h1>{t("title")}</h1>
            <p>{t("description")}</p>
          </div>
          <div className="header-actions" style={{ display: 'flex', alignItems: 'center', gap: '24px' }}>
            {fetchError && !loading && (
              <div style={{ display: 'flex', alignItems: 'center', gap: '8px', background: 'rgba(239, 68, 68, 0.1)', padding: '6px 12px', borderRadius: '999px', border: '1px solid rgba(239, 68, 68, 0.2)' }}>
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#ef4444" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>
                <span style={{ fontSize: "12px", color: "#f87171", fontWeight: 500 }}>API Offline</span>
              </div>
            )}
            
            <UserProfileWidget />
            
            <button className="btn btn-primary" style={{ borderRadius: '999px', padding: '10px 20px', fontSize: '14px' }}>
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
              Deploy Config
            </button>
          </div>
        </div>

        {loading ? (
          <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '60vh', flexDirection: 'column', gap: '16px' }}>
             <div style={{ width: '40px', height: '40px', borderRadius: '50%', border: '3px solid rgba(59, 130, 246, 0.1)', borderTopColor: '#3b82f6', animation: 'spin 1s linear infinite' }} />
             <div style={{ color: 'var(--text-secondary)', fontSize: '14px', letterSpacing: '0.05em' }}>LOADING DASHBOARD...</div>
          </div>
        ) : (
          <>
            {/* Metrics Grid */}
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
                  <div className="metric-value">{overview.gateways} <span style={{fontSize: '20px', color: 'var(--text-secondary)'}}>/ {gateways.length}</span></div>
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

              <div className="metric-card glass-panel animate-fade-in delay-3">
                <div className="metric-header">
                  <span>Active Incidents</span>
                  <div className="icon-wrapper" style={{ background: 'rgba(245,158,11,0.12)' }}>
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#fbbf24" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path><line x1="12" y1="9" x2="12" y2="13"></line><line x1="12" y1="17" x2="12.01" y2="17"></line></svg>
                  </div>
                </div>
                <div>
                  <div className="metric-value">{overview.active_incidents}</div>
                  <div className={`metric-trend ${overview.active_incidents > 0 ? "trend-down" : "trend-up"}`}>
                    {overview.active_incidents > 0 ? "Open" : "No open incidents"}
                  </div>
                </div>
              </div>
            </div>

            {/* Traffic & Security Summaries */}
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '24px', marginBottom: '32px' }}>
              
              <div className="data-section animate-fade-in delay-3" style={{ marginBottom: 0 }}>
                <div className="section-header">
                  <h2>Traffic Summary</h2>
                </div>
                <div className="metrics-grid" style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                  <div style={{ display: 'flex', gap: '16px' }}>
                    <div className="metric-card glass-panel" style={{ flex: 1, padding: '16px 20px' }}>
                      <div className="metric-header" style={{ fontSize: '13px' }}><span>Total Requests</span></div>
                      <div className="metric-value" style={{ fontSize: '28px' }}>{traffic.total_requests.toLocaleString()}</div>
                    </div>
                    <div className="metric-card glass-panel" style={{ flex: 1, padding: '16px 20px' }}>
                      <div className="metric-header" style={{ fontSize: '13px' }}><span>Avg Latency</span></div>
                      <div className="metric-value" style={{ fontSize: '28px' }}>{traffic.avg_latency_ms.toFixed(1)}<span style={{fontSize: '14px', color: 'var(--text-secondary)'}}>ms</span></div>
                    </div>
                  </div>
                  <div className="metric-card glass-panel" style={{ padding: '20px' }}>
                    <div className="metric-header" style={{ marginBottom: '12px', fontSize: '13px' }}><span>Status Breakdown</span></div>
                    <div style={{ fontSize: "13px", color: "var(--text-secondary)", display: 'flex', flexDirection: 'column', gap: '10px' }}>
                      {traffic.by_status.slice(0, 4).map((s) => (
                        <div key={s.status} style={{ display: "flex", justifyContent: "space-between", alignItems: 'center' }}>
                          <span style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                            <span style={{ width: '8px', height: '8px', borderRadius: '50%', background: s.status.startsWith('2') ? 'var(--success)' : s.status.startsWith('4') ? 'var(--warning)' : s.status.startsWith('5') ? 'var(--danger)' : 'var(--accent-primary)' }} />
                            HTTP {s.status}
                          </span>
                          <span style={{ color: "var(--text-primary)", fontWeight: 500 }}>{s.count.toLocaleString()}</span>
                        </div>
                      ))}
                      {traffic.by_status.length === 0 && <span style={{ opacity: 0.5 }}>No traffic data available.</span>}
                    </div>
                  </div>
                </div>
              </div>

              <div className="data-section animate-fade-in delay-3" style={{ marginBottom: 0 }}>
                <div className="section-header">
                  <h2>Security Summary</h2>
                </div>
                <div className="metrics-grid" style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                  <div style={{ display: 'flex', gap: '16px' }}>
                    <div className="metric-card glass-panel" style={{ flex: 1, padding: '16px 20px' }}>
                      <div className="metric-header" style={{ fontSize: '13px' }}><span>Total Events</span></div>
                      <div className="metric-value" style={{ fontSize: '28px' }}>{security.total_events.toLocaleString()}</div>
                    </div>
                    <div className="metric-card glass-panel" style={{ flex: 1, padding: '16px 20px' }}>
                      <div className="metric-header" style={{ fontSize: '13px' }}><span>Blocked</span></div>
                      <div className="metric-value" style={{ fontSize: '28px', color: 'var(--danger)' }}>{security.blocked_events.toLocaleString()}</div>
                    </div>
                  </div>
                  <div className="metric-card glass-panel" style={{ padding: '20px' }}>
                    <div className="metric-header" style={{ marginBottom: '12px', fontSize: '13px' }}><span>By Severity</span></div>
                    <div style={{ fontSize: "13px", color: "var(--text-secondary)", display: 'flex', flexDirection: 'column', gap: '10px' }}>
                      {security.by_severity.slice(0, 4).map((s) => (
                        <div key={s.severity} style={{ display: "flex", justifyContent: "space-between", alignItems: 'center' }}>
                          <span style={{ display: 'flex', alignItems: 'center', gap: '8px', textTransform: "capitalize" }}>
                            <span style={{ width: '8px', height: '8px', borderRadius: '50%', background: s.severity === 'critical' ? 'var(--danger)' : s.severity === 'high' ? 'var(--warning)' : s.severity === 'medium' ? 'var(--accent-primary)' : 'var(--success)' }} />
                            {s.severity}
                          </span>
                          <span style={{ color: "var(--text-primary)", fontWeight: 500 }}>{s.count.toLocaleString()}</span>
                        </div>
                      ))}
                      {security.by_severity.length === 0 && <span style={{ opacity: 0.5 }}>No security events recorded.</span>}
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Top IPs + Top Rules */}
            <div className="metrics-grid" style={{ marginBottom: 32 }}>
              <div className="metric-card glass-panel">
                <div className="metric-header" style={{ marginBottom: '12px' }}><span>Top Client IPs</span></div>
                <div style={{ fontSize: "13px", color: "var(--text-secondary)", display: 'flex', flexDirection: 'column', gap: '10px' }}>
                  {topIPs.slice(0, 5).map((ip) => (
                    <div key={ip.client_ip} style={{ display: "flex", justifyContent: "space-between", alignItems: 'center' }}>
                      <span style={{ fontFamily: "monospace", fontSize: "13px", padding: '2px 6px', background: 'rgba(255,255,255,0.05)', borderRadius: '4px' }}>{ip.client_ip}</span>
                      <span>
                        {ip.hits}
                        {ip.blocked > 0 && <span style={{ color: "var(--danger)", marginLeft: "6px" }}>({ip.blocked} blk)</span>}
                      </span>
                    </div>
                  ))}
                  {topIPs.length === 0 && <span style={{ opacity: 0.5 }}>No IPs recorded.</span>}
                </div>
              </div>
              <div className="metric-card glass-panel">
                <div className="metric-header" style={{ marginBottom: '12px' }}><span>Top Triggered Rules</span></div>
                <div style={{ fontSize: "13px", color: "var(--text-secondary)", display: 'flex', flexDirection: 'column', gap: '10px' }}>
                  {topRules.slice(0, 5).map((rule) => (
                    <div key={rule.rule_id} style={{ display: "flex", justifyContent: "space-between", alignItems: 'center' }}>
                      <span style={{ fontFamily: "monospace", fontSize: "13px", padding: '2px 6px', background: 'rgba(255,255,255,0.05)', borderRadius: '4px' }}>{rule.rule_id}</span>
                      <span style={{ color: "var(--text-primary)", fontWeight: 500 }}>{rule.hits}</span>
                    </div>
                  ))}
                  {topRules.length === 0 && <span style={{ opacity: 0.5 }}>No rules triggered.</span>}
                </div>
              </div>
            </div>  

            {/* Applications List */}
            <div className="data-section animate-fade-in delay-3">
              <div className="section-header">
                <h2>Protected Applications</h2>
              </div>
              <div className="data-list glass-panel" style={{ overflow: 'hidden' }}>
                {apps.length === 0 && !fetchError && (
                  <div style={{ padding: "60px", textAlign: "center", color: "var(--text-secondary)", display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '16px' }}>
                    <div style={{ background: 'rgba(255,255,255,0.03)', padding: '16px', borderRadius: '50%' }}>
                      <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="2" y1="12" x2="22" y2="12"></line><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path></svg>
                    </div>
                    <div>{t("no_applications")}</div>
                  </div>
                )}
                {apps.map((a, index) => (
                  <div key={a.id} className="list-item" style={{ borderBottom: index < apps.length - 1 ? '1px solid var(--glass-border)' : 'none', borderRadius: 0 }}>
                    <div className="item-info">
                      <div className="icon-wrapper" style={{ background: 'rgba(59, 130, 246, 0.1)', padding: '12px', color: '#60a5fa' }}>
                        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="2" y1="12" x2="22" y2="12"></line><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path></svg>
                      </div>
                      <div>
                        <div className="item-title" style={{ fontSize: '15px' }}>{a.name}</div>
                        <div className="item-subtitle" style={{ fontSize: '12px', fontFamily: 'monospace' }}>ID: {a.id} • WAF Enabled</div>
                      </div>
                    </div>
                    <span className={`badge ${a.status === 'active' ? 'badge-active' : a.status === 'warning' ? 'badge-warning' : 'badge-danger'}`}>
                      {a.status}
                    </span>
                  </div>
                ))}
              </div>
            </div>

          </>
        )}
      </main>
    </div>
  );
}
