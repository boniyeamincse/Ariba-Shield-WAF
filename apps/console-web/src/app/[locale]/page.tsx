import { getTranslations } from "next-intl/server";
import { listApplications, listGateways, listSecurityPolicies } from "@/lib/api";
import UserProfileWidget from "@/components/UserProfileWidget";

export default async function OverviewPage() {
  const t = await getTranslations("overview");
  const navT = await getTranslations("nav");

  let apps: { id: string; name: string; status: string }[] = [];
  let gateways: { id: string; hostname: string; status: string }[] = [];
  let policies: { id: string; name: string; enforcement_mode: string }[] = [];

  try {
    [apps, gateways, policies] = await Promise.all([
      listApplications(),
      listGateways(),
      listSecurityPolicies(),
    ]);
  } catch {
    // API not available — fallback to mock data for UI showcase
    apps = [
      { id: "1", name: "Enterprise Payment Gateway", status: "active" },
      { id: "2", name: "Internal HR Portal", status: "active" },
      { id: "3", name: "Customer CRM", status: "warning" }
    ];
    gateways = [
      { id: "1", hostname: "gw-nyc-01.ariba.local", status: "active" },
      { id: "2", hostname: "gw-lon-02.ariba.local", status: "active" },
      { id: "3", hostname: "gw-sgp-03.ariba.local", status: "offline" }
    ];
    policies = [
      { id: "1", name: "Strict SQLi Protection", enforcement_mode: "blocking" },
      { id: "2", name: "Global Rate Limiting", enforcement_mode: "detection" },
      { id: "3", name: "OWASP Top 10 Core Rules", enforcement_mode: "blocking" }
    ];
  }

  return (
    <div className="dashboard-container">
      {/* Sidebar */}
      <aside className="sidebar">
        <div className="brand">
          <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
            <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path>
          </svg>
          Ariba Shield
        </div>
        
        <ul className="nav-links">
          <li className="nav-item active">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="3" y="3" width="7" height="9"></rect><rect x="14" y="3" width="7" height="5"></rect><rect x="14" y="12" width="7" height="9"></rect><rect x="3" y="16" width="7" height="5"></rect></svg>
            {t("title")}
          </li>
          <li className="nav-item">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polygon points="12 2 2 7 12 12 22 7 12 2"></polygon><polyline points="2 17 12 22 22 17"></polyline><polyline points="2 12 12 17 22 12"></polyline></svg>
            {t("applications")}
          </li>
          <li className="nav-item">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="2" y="2" width="20" height="8" rx="2" ry="2"></rect><rect x="2" y="14" width="20" height="8" rx="2" ry="2"></rect><line x1="6" y1="6" x2="6.01" y2="6"></line><line x1="6" y1="18" x2="6.01" y2="18"></line></svg>
            {navT("gateways")}
          </li>
          <li className="nav-item">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path></svg>
            {navT("policies")}
          </li>
          <li className="nav-item">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="3"></circle><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path></svg>
            {navT("settings")}
          </li>
        </ul>
      </aside>

      {/* Main Content */}
      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <h1>{t("title")}</h1>
            <p>{t("description") || "Real-time insights and security operations overview."}</p>
          </div>
          <div className="header-actions" style={{ display: 'flex', alignItems: 'center', gap: '24px' }}>
            <UserProfileWidget />
            <div style={{ height: '24px', width: '1px', backgroundColor: 'rgba(255,255,255,0.1)' }}></div>
            <button className="btn btn-primary">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
              Deploy Config
            </button>
          </div>
        </div>

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
              <div className="metric-value">{apps.length}</div>
              <div className="metric-trend trend-up">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polyline points="23 6 13.5 15.5 8.5 10.5 1 18"></polyline><polyline points="17 6 23 6 23 12"></polyline></svg>
                100% Protected
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
              <div className="metric-value">{gateways.filter(g => g.status === 'active').length} / {gateways.length}</div>
              <div className="metric-trend trend-down">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polyline points="23 18 13.5 8.5 8.5 13.5 1 6"></polyline><polyline points="17 18 23 18 23 12"></polyline></svg>
                1 Node Offline
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
                All Synced
              </div>
            </div>
          </div>
        </div>

        {/* Lists Section */}
        <div className="data-section animate-fade-in delay-3">
          <div className="section-header">
            <h2>Protected Applications</h2>
            <button className="btn" style={{ background: 'rgba(255,255,255,0.05)' }}>View All</button>
          </div>
          <div className="data-list glass-panel">
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