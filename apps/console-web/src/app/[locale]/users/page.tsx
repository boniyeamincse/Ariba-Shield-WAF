import { getTranslations } from "next-intl/server";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";

export default async function UsersAccessPage({ params }: { params: Promise<{ locale: string }> }) {
  const { locale } = await params;
  const t = await getTranslations("users");

  // Mock data for the Users table
  const usersList = [
    { id: "1", name: "Boni Yeamin", email: "boni@aribashield.local", role: "Super Admin", status: "Active", lastLogin: "2 mins ago" },
    { id: "2", name: "Raihan Ali", email: "raihan@aribashield.local", role: "Security Analyst", status: "Active", lastLogin: "1 hour ago" },
    { id: "3", name: "Sara Rahman", email: "sara@aribashield.local", role: "Viewer", status: "Invited", lastLogin: "Never" },
  ];

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
            <UserProfileWidget />
            <div style={{ height: '24px', width: '1px', backgroundColor: 'rgba(255,255,255,0.1)' }}></div>
            <button className="btn btn-primary">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><line x1="19" y1="8" x2="19" y2="14"></line><line x1="22" y1="11" x2="16" y2="11"></line></svg>
              {t("add_user")}
            </button>
          </div>
        </div>

        <div className="data-section animate-fade-in delay-1">
          <div className="data-list glass-panel" style={{ padding: '0', overflow: 'hidden' }}>
            <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left' }}>
              <thead>
                <tr style={{ borderBottom: '1px solid rgba(255,255,255,0.08)', backgroundColor: 'rgba(255,255,255,0.02)' }}>
                  <th style={{ padding: '16px 24px', color: 'var(--text-secondary)', fontWeight: 500, whiteSpace: 'nowrap' }}>{t("name")}</th>
                  <th style={{ padding: '16px 24px', color: 'var(--text-secondary)', fontWeight: 500, whiteSpace: 'nowrap' }}>{t("email")}</th>
                  <th style={{ padding: '16px 24px', color: 'var(--text-secondary)', fontWeight: 500, whiteSpace: 'nowrap' }}>{t("role")}</th>
                  <th style={{ padding: '16px 24px', color: 'var(--text-secondary)', fontWeight: 500, whiteSpace: 'nowrap' }}>{t("status")}</th>
                  <th style={{ padding: '16px 24px', color: 'var(--text-secondary)', fontWeight: 500, whiteSpace: 'nowrap' }}>{t("last_login")}</th>
                  <th style={{ padding: '16px 24px', color: 'var(--text-secondary)', fontWeight: 500, textAlign: 'right', whiteSpace: 'nowrap' }}>{t("actions")}</th>
                </tr>
              </thead>
              <tbody>
                {usersList.map((user, idx) => (
                  <tr key={user.id} style={{ borderBottom: idx !== usersList.length - 1 ? '1px solid rgba(255,255,255,0.05)' : 'none', transition: 'background 0.2s' }}>
                    <td style={{ padding: '16px 24px', whiteSpace: 'nowrap' }}>
                      <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                        <div style={{ width: '32px', height: '32px', borderRadius: '50%', background: 'linear-gradient(135deg, #3b82f6, #8b5cf6)', display: 'grid', placeItems: 'center', fontSize: '14px', fontWeight: 600 }}>
                          {user.name.charAt(0)}
                        </div>
                        <span style={{ fontWeight: 500 }}>{user.name}</span>
                      </div>
                    </td>
                    <td style={{ padding: '16px 24px', color: 'var(--text-secondary)', whiteSpace: 'nowrap' }}>{user.email}</td>
                    <td style={{ padding: '16px 24px', whiteSpace: 'nowrap' }}>
                      <span style={{ padding: '4px 10px', borderRadius: '6px', background: 'rgba(255,255,255,0.08)', fontSize: '12px', fontWeight: 500, display: 'inline-block' }}>
                        {user.role}
                      </span>
                    </td>
                    <td style={{ padding: '16px 24px', whiteSpace: 'nowrap' }}>
                      <span className={`badge ${user.status === 'Active' ? 'badge-active' : 'badge-warning'}`}>
                        {user.status === 'Active' ? t("active") : t("invited")}
                      </span>
                    </td>
                    <td style={{ padding: '16px 24px', color: 'var(--text-secondary)', fontSize: '14px', whiteSpace: 'nowrap' }}>{user.lastLogin}</td>
                    <td style={{ padding: '16px 24px', textAlign: 'right', whiteSpace: 'nowrap' }}>
                      <button className="btn btn-ghost" style={{ padding: '6px', background: 'transparent', border: 'none', color: 'var(--text-secondary)', cursor: 'pointer', display: 'inline-flex', alignItems: 'center', justifyContent: 'center' }}>
                        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="1"></circle><circle cx="19" cy="12" r="1"></circle><circle cx="5" cy="12" r="1"></circle></svg>
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

      </main>
    </div>
  );
}
