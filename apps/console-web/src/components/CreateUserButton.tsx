"use client";

import { useState } from "react";
import { createUser } from "@/lib/api";

export default function CreateUserButton() {
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
      setTimeout(() => {
        setShowModal(false);
        window.location.reload();
      }, 1200);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Failed to create user");
    } finally {
      setCreating(false);
    }
  };

  if (!showModal) {
    return (
      <button className="btn btn-primary" onClick={() => setShowModal(true)}>
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"></path>
          <circle cx="9" cy="7" r="4"></circle>
          <line x1="19" y1="8" x2="19" y2="14"></line>
          <line x1="22" y1="11" x2="16" y2="11"></line>
        </svg>
        Add User
      </button>
    );
  }

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

      {/* Overlay */}
      <div
        onClick={() => !creating && setShowModal(false)}
        style={{
          position: "fixed",
          inset: 0,
          background: "rgba(0,0,0,0.6)",
          backdropFilter: "blur(4px)",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          zIndex: 1000,
          animation: "fadeIn 0.15s ease",
        }}
      >
        {/* Modal card */}
        <div
          onClick={(e) => e.stopPropagation()}
          className="glass-panel"
          style={{
            width: "90%",
            maxWidth: "440px",
            padding: 0,
            overflow: "hidden",
            animation: "fadeIn 0.2s ease, slideUp 0.2s ease",
          }}
        >
          {/* Header */}
          <div
            style={{
              padding: "24px 28px 16px",
              borderBottom: "1px solid rgba(255,255,255,0.08)",
            }}
          >
            <div style={{ display: "flex", alignItems: "center", gap: "12px" }}>
              <div
                style={{
                  width: "40px",
                  height: "40px",
                  borderRadius: "10px",
                  background: "linear-gradient(135deg, rgba(59,130,246,0.2), rgba(139,92,246,0.2))",
                  border: "1px solid rgba(59,130,246,0.3)",
                  display: "grid",
                  placeItems: "center",
                  color: "#60a5fa",
                  flexShrink: 0,
                }}
              >
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"></path>
                  <circle cx="9" cy="7" r="4"></circle>
                  <line x1="19" y1="8" x2="19" y2="14"></line>
                  <line x1="22" y1="11" x2="16" y2="11"></line>
                </svg>
              </div>
              <div>
                <h2 style={{ margin: 0, fontSize: "17px", fontWeight: 600 }}>Create New User</h2>
                <p style={{ margin: "2px 0 0", fontSize: "12px", color: "var(--text-secondary)" }}>
                  Add a new team member to this workspace
                </p>
              </div>
            </div>
          </div>

          {/* Form body */}
          <form onSubmit={handleCreate} style={{ padding: "20px 28px 24px", display: "flex", flexDirection: "column", gap: "16px" }}>
            {/* Email */}
            <div>
              <label style={{ display: "block", fontSize: "12px", fontWeight: 600, color: "var(--text-secondary)", marginBottom: "6px" }}>
                Email address
              </label>
              <input
                type="email"
                className="form-control"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="user@example.com"
                autoFocus
                style={{ padding: "10px 12px", fontSize: "14px" }}
              />
            </div>

            {/* Password */}
            <div>
              <label style={{ display: "block", fontSize: "12px", fontWeight: 600, color: "var(--text-secondary)", marginBottom: "6px" }}>
                Password
              </label>
              <input
                type="password"
                className="form-control"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Secure password"
                style={{ padding: "10px 12px", fontSize: "14px" }}
              />
            </div>

            {/* Role */}
            <div>
              <label style={{ display: "block", fontSize: "12px", fontWeight: 600, color: "var(--text-secondary)", marginBottom: "6px" }}>
                Role
              </label>
              <select
                className="form-control"
                value={role}
                onChange={(e) => setRole(e.target.value)}
                style={{ padding: "10px 12px", fontSize: "14px" }}
              >
                {[
                  { value: "Super Admin", color: "#f87171", desc: "Full system access" },
                  { value: "Platform Admin", color: "#c084fc", desc: "Platform management" },
                  { value: "Security Admin", color: "#fbbf24", desc: "Security policy management" },
                  { value: "App Owner", color: "#22d3ee", desc: "Application management" },
                  { value: "SOC Analyst", color: "#60a5fa", desc: "Security event monitoring" },
                  { value: "Auditor", color: "#34d399", desc: "Read-only audit access" },
                  { value: "Read Only", color: "#9ca3af", desc: "View-only access" },
                ].map((r) => (
                  <option key={r.value} value={r.value}>
                    {r.value}
                  </option>
                ))}
              </select>
              {role && (
                <div style={{ marginTop: "6px", fontSize: "11px", color: "var(--text-secondary)" }}>
                  <span
                    style={{
                      display: "inline-block",
                      width: "8px",
                      height: "8px",
                      borderRadius: "50%",
                      background: ({
                        "Super Admin": "#f87171",
                        "Platform Admin": "#c084fc",
                        "Security Admin": "#fbbf24",
                        "App Owner": "#22d3ee",
                        "SOC Analyst": "#60a5fa",
                        "Auditor": "#34d399",
                        "Read Only": "#9ca3af",
                      }[role] || "#9ca3af"),
                      marginRight: "6px",
                    }}
                  />
                  {role} — {{
                    "Super Admin": "Full system access",
                    "Platform Admin": "Platform management",
                    "Security Admin": "Security policy management",
                    "App Owner": "Application management",
                    "SOC Analyst": "Security event monitoring",
                    "Auditor": "Read-only audit access",
                    "Read Only": "View-only access",
                  }[role] || ""}
                </div>
              )}
            </div>

            {/* Messages */}
            {error && (
              <div
                style={{
                  padding: "10px 14px",
                  borderRadius: "8px",
                  background: "rgba(239,68,68,0.1)",
                  border: "1px solid rgba(239,68,68,0.2)",
                  color: "#f87171",
                  fontSize: "13px",
                  display: "flex",
                  alignItems: "center",
                  gap: "8px",
                }}
              >
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ flexShrink: 0 }}>
                  <circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line>
                </svg>
                {error}
              </div>
            )}
            {success && (
              <div
                style={{
                  padding: "10px 14px",
                  borderRadius: "8px",
                  background: "rgba(16,185,129,0.1)",
                  border: "1px solid rgba(16,185,129,0.2)",
                  color: "#34d399",
                  fontSize: "13px",
                  display: "flex",
                  alignItems: "center",
                  gap: "8px",
                }}
              >
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ flexShrink: 0 }}>
                  <polyline points="20 6 9 17 4 12"></polyline>
                </svg>
                {success}
              </div>
            )}

            {/* Actions */}
            <div
              style={{
                display: "flex",
                gap: "10px",
                justifyContent: "flex-end",
                borderTop: "1px solid rgba(255,255,255,0.06)",
                paddingTop: "16px",
                marginTop: "4px",
              }}
            >
              <button
                type="button"
                className="btn"
                onClick={() => setShowModal(false)}
                disabled={creating}
              >
                Cancel
              </button>
              <button type="submit" className="btn btn-primary" disabled={creating}>
                {creating ? (
                  <span style={{ display: "flex", alignItems: "center", gap: "6px" }}>
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" style={{ animation: "spin 1s linear infinite" }}>
                      <circle cx="12" cy="12" r="10" strokeDasharray="31.4 31.4" strokeLinecap="round"></circle>
                    </svg>
                    Creating...
                  </span>
                ) : (
                  <span style={{ display: "flex", alignItems: "center", gap: "6px" }}>
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"></path>
                      <circle cx="9" cy="7" r="4"></circle>
                      <line x1="19" y1="8" x2="19" y2="14"></line>
                      <line x1="22" y1="11" x2="16" y2="11"></line>
                    </svg>
                    Create User
                  </span>
                )}
              </button>
            </div>
          </form>
        </div>
      </div>
    </>
  );
}