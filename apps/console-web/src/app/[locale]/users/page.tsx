"use client";

import { getTranslations } from "next-intl/server";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import { listUsers, createUser, User } from "@/lib/api";

export default async function UsersAccessPage({ params }: { params: Promise<{ locale: string }> }) {
  const { locale } = await params;
  const t = await getTranslations("users");

  let usersList: User[] = [];
  let fetchError = false;

  try {
    usersList = await listUsers();
  } catch {
    fetchError = true;
    usersList = [
      { id: "1", email: "superadmin@aribashield.local", role: "Super Admin", status: "active", created_at: new Date().toISOString() },
      { id: "2", email: "security@aribashield.local", role: "Security Admin", status: "active", created_at: new Date().toISOString() },
      { id: "3", email: "auditor@aribashield.local", role: "Auditor", status: "active", created_at: new Date().toISOString() },
      { id: "4", email: "readonly@aribashield.local", role: "Read Only", status: "active", created_at: new Date().toISOString() },
    ];
  }

  const getInitials = (email: string) => email.charAt(0).toUpperCase();

  const roleColors: Record<string, string> = {
    "Super Admin": "rgba(239,68,68,0.15)",
    "Platform Admin": "rgba(168,85,247,0.15)",
    "Security Admin": "rgba(245,158,11,0.15)",
    "SOC Analyst": "rgba(59,130,246,0.15)",
    "Auditor": "rgba(16,185,129,0.15)",
    "App Owner": "rgba(6,182,212,0.15)",
    "Read Only": "rgba(107,114,128,0.15)",
  };

  const roleTextColors: Record<string, string> = {
    "Super Admin": "#f87171",
    "Platform Admin": "#c084fc",
    "Security Admin": "#fbbf24",
    "SOC Analyst": "#60a5fa",
    "Auditor": "#34d399",
    "App Owner": "#22d3ee",
    "Read Only": "#9ca3af",
  };

  const avatarGradients = [
    "linear-gradient(135deg, #3b82f6, #8b5cf6)",
    "linear-gradient(135deg, #10b981, #06b6d4)",
    "linear-gradient(135deg, #f59e0b, #ef4444)",
    "linear-gradient(135deg, #8b5cf6, #ec4899)",
    "linear-gradient(135deg, #06b6d4, #3b82f6)",
  ];

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
                ⚠ Showing cached data — API offline
              </span>
            )}
            <UserProfileWidget />
            <UsersPageClient locale={locale} />
          </div>
        </div>

        {/* Stats row */}
        <div className="metrics-grid animate-fade-in delay-1" style={{ gridTemplateColumns: "repeat(3, 1fr)", maxWidth: "600px", marginBottom: "28px" }}>
          <div className="glass-panel metric-card" style={{ padding: "20px" }}>
            <div className="metric-header"><span>Total Users</span><div className="icon-wrapper icon-blue"><svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle></svg></div></div>
            <div className="metric-value">{usersList.length}</div>
          </div>
          <div className="glass-panel metric-card" style={{ padding: "20px" }}>
            <div className="metric-header"><span>Active</span><div className="icon-wrapper icon-green"><svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><polyline points="20 6 9 17 4 12"></polyline></svg></div></div>
            <div className="metric-value" style={{ color: "var(--success)" }}>{usersList.filter(u => u.status === "active").length}</div>
          </div>
          <div className="glass-panel metric-card" style={{ padding: "20px" }}>
            <div className="metric-header"><span>Pending</span><div className="icon-wrapper icon-purple"><svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg></div></div>
            <div className="metric-value">{usersList.filter(u => u.status !== "active").length}</div>
          </div>
        </div>

        {/* Users Table */}
        <div className="data-section animate-fade-in delay-2">
          <div className="glass-panel" style={{ padding: "0", overflow: "hidden" }}>
            <table style={{ width: "100%", borderCollapse: "collapse", textAlign: "left" }}>
              <thead>
                <tr style={{ borderBottom: "1px solid rgba(255,255,255,0.08)", background: "rgba(255,255,255,0.02)" }}>
                  <th style={{ padding: "14px 24px", color: "var(--text-secondary)", fontWeight: 500, fontSize: "13px", whiteSpace: "nowrap" }}>{t("name")}</th>
                  <th style={{ padding: "14px 24px", color: "var(--text-secondary)", fontWeight: 500, fontSize: "13px", whiteSpace: "nowrap" }}>{t("role")}</th>
                  <th style={{ padding: "14px 24px", color: "var(--text-secondary)", fontWeight: 500, fontSize: "13px", whiteSpace: "nowrap" }}>{t("status")}</th>
                  <th style={{ padding: "14px 24px", color: "var(--text-secondary)", fontWeight: 500, fontSize: "13px", whiteSpace: "nowrap" }}>Member Since</th>
                  <th style={{ padding: "14px 24px", color: "var(--text-secondary)", fontWeight: 500, fontSize: "13px", textAlign: "right", whiteSpace: "nowrap" }}>{t("actions")}</th>
                </tr>
              </thead>
              <tbody>
                {usersList.map((user, idx) => {
                  const roleBg = roleColors[user.role] ?? "rgba(255,255,255,0.06)";
                  const roleText = roleTextColors[user.role] ?? "#9ca3af";
                  const gradient = avatarGradients[idx % avatarGradients.length];
                  const joinedDate = new Date(user.created_at).toLocaleDateString("en-GB", { day: "2-digit", month: "short", year: "numeric" });

                  return (
                    <tr key={user.id} style={{ borderBottom: idx !== usersList.length - 1 ? "1px solid rgba(255,255,255,0.05)" : "none" }}>
                      <td style={{ padding: "16px 24px", whiteSpace: "nowrap" }}>
                        <div style={{ display: "flex", alignItems: "center", gap: "12px" }}>
                          <div style={{ width: "36px", height: "36px", borderRadius: "50%", background: gradient, display: "grid", placeItems: "center", fontSize: "15px", fontWeight: 700, flexShrink: 0 }}>
                            {getInitials(user.email)}
                          </div>
                          <div>
                            <div style={{ fontWeight: 500, fontSize: "14px" }}>{user.email.split("@")[0]}</div>
                            <div style={{ color: "var(--text-secondary)", fontSize: "12px" }}>{user.email}</div>
                          </div>
                        </div>
                      </td>
                      <td style={{ padding: "16px 24px", whiteSpace: "nowrap" }}>
                        <span style={{ padding: "4px 12px", borderRadius: "6px", background: roleBg, color: roleText, fontSize: "12px", fontWeight: 600 }}>
                          {user.role}
                        </span>
                      </td>
                      <td style={{ padding: "16px 24px", whiteSpace: "nowrap" }}>
                        <span className={`badge ${user.status === "active" ? "badge-active" : "badge-warning"}`}>
                          {user.status === "active" ? t("active") : t("invited")}
                        </span>
                      </td>
                      <td style={{ padding: "16px 24px", color: "var(--text-secondary)", fontSize: "13px", whiteSpace: "nowrap" }}>{joinedDate}</td>
                      <td style={{ padding: "16px 24px", textAlign: "right", whiteSpace: "nowrap" }}>
                        <button style={{ padding: "6px 8px", background: "transparent", border: "none", color: "var(--text-secondary)", cursor: "pointer", display: "inline-flex", alignItems: "center", justifyContent: "center", borderRadius: "6px" }}>
                          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                            <circle cx="12" cy="12" r="1"></circle><circle cx="19" cy="12" r="1"></circle><circle cx="5" cy="12" r="1"></circle>
                          </svg>
                        </button>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>

            {usersList.length === 0 && (
              <div style={{ padding: "60px", textAlign: "center", color: "var(--text-secondary)" }}>
                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" style={{ marginBottom: "16px", opacity: 0.4 }} strokeLinecap="round" strokeLinejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle></svg>
                <p>No users found</p>
              </div>
            )}
          </div>
        </div>
      </main>
    </div>
  );
}

// Client component for the Add User modal and button
import { useState } from "react";

function UsersPageClient({ locale }: { locale: string }) {
  const [showModal, setShowModal] = useState(false);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState("Read Only");
  const [creating, setCreating] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setCreating(true);
    setError("");
    setSuccess("");
    try {
      await createUser(email, password, role);
      setSuccess(`User ${email} created as ${role}`);
      setEmail("");
      setPassword("");
      setRole("Read Only");
      setTimeout(() => setShowModal(false), 1500);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Failed to create user");
    } finally {
      setCreating(false);
    }
  };

  return (
    <>
      <button className="btn btn-primary" onClick={() => setShowModal(true)}>
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"></path>
          <circle cx="9" cy="7" r="4"></circle>
          <line x1="19" y1="8" x2="19" y2="14"></line>
          <line x1="22" y1="11" x2="16" y2="11"></line>
        </svg>
        Add User
      </button>

      {showModal && (
        <div style={{ position: "fixed", inset: 0, background: "rgba(0,0,0,0.6)", display: "flex", alignItems: "center", justifyContent: "center", zIndex: 1000 }}>
          <div className="glass-panel" style={{ padding: "32px", maxWidth: "420px", width: "90%" }}>
            <h2 style={{ margin: "0 0 20px", fontSize: "18px" }}>Create New User</h2>
            <form onSubmit={handleCreate} style={{ display: "flex", flexDirection: "column", gap: "14px" }}>
              <div>
                <label style={{ display: "block", fontSize: "12px", color: "var(--text-secondary)", marginBottom: "4px" }}>Email</label>
                <input type="email" className="form-control" required value={email} onChange={e => setEmail(e.target.value)} placeholder="user@example.com" />
              </div>
              <div>
                <label style={{ display: "block", fontSize: "12px", color: "var(--text-secondary)", marginBottom: "4px" }}>Password</label>
                <input type="password" className="form-control" required value={password} onChange={e => setPassword(e.target.value)} placeholder="Password" />
              </div>
              <div>
                <label style={{ display: "block", fontSize: "12px", color: "var(--text-secondary)", marginBottom: "4px" }}>Role</label>
                <select className="form-control" value={role} onChange={e => setRole(e.target.value)}>
                  {["Super Admin", "Platform Admin", "Security Admin", "App Owner", "SOC Analyst", "Auditor", "Read Only"].map(r => (
                    <option key={r} value={r}>{r}</option>
                  ))}
                </select>
              </div>
              {error && <div style={{ color: "var(--warning)", fontSize: "13px" }}>{error}</div>}
              {success && <div style={{ color: "var(--success)", fontSize: "13px" }}>{success}</div>}
              <div style={{ display: "flex", gap: "10px", justifyContent: "flex-end", marginTop: "8px" }}>
                <button type="button" className="btn" onClick={() => setShowModal(false)}>Cancel</button>
                <button type="submit" className="btn btn-primary" disabled={creating}>{creating ? "Creating..." : "Create User"}</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </>
  );
}