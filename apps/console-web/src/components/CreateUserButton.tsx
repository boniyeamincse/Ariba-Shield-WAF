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