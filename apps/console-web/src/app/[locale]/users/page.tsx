"use client";

import { useEffect, useState } from "react";
import { useLocale } from "next-intl";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import CreateUserButton from "@/components/CreateUserButton";
import DataTable, { type Column } from "@/components/shared/DataTable";
import ConfirmDialog from "@/components/shared/ConfirmDialog";
import { listUsers, updateUser, deleteUser, type User } from "@/lib/api";

const ROLE_OPTIONS = ["Super Admin", "Platform Admin", "Security Admin", "App Owner", "SOC Analyst", "Auditor", "Read Only"];

export default function UsersAccessPage() {
  const locale = useLocale();
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [editing, setEditing] = useState<User | null>(null);
  const [deleting, setDeleting] = useState<User | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [formRole, setFormRole] = useState("");
  const [formStatus, setFormStatus] = useState("active");

  const load = async () => {
    setLoading(true);
    setError("");
    try {
      setUsers(await listUsers());
    } catch {
      setError("Failed to load users");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const openEdit = (user: User) => {
    setEditing(user);
    setFormRole(user.role);
    setFormStatus(user.status);
  };

  const submit = async () => {
    if (!editing) return;
    setSubmitting(true);
    setError("");
    try {
      await updateUser(editing.id, { role: formRole, status: formStatus });
      setEditing(null);
      await load();
    } catch {
      setError("Failed to update user");
    } finally {
      setSubmitting(false);
    }
  };

  const confirmDelete = async () => {
    if (!deleting) return;
    setSubmitting(true);
    try {
      await deleteUser(deleting.id);
      setDeleting(null);
      await load();
    } catch {
      setError("Failed to delete user");
    } finally {
      setSubmitting(false);
    }
  };

  const activeCount = users.filter((u) => u.status === "active").length;

  const columns: Column<User>[] = [
    {
      key: "email",
      label: "User",
      render: (row) => (
        <div>
          <div style={{ fontWeight: 500, fontSize: "14px" }}>{row.email.split("@")[0]}</div>
          <div style={{ color: "var(--text-secondary)", fontSize: "12px" }}>{row.email}</div>
        </div>
      ),
    },
    {
      key: "role",
      label: "Role",
      render: (row) => (
        <span style={{ padding: "4px 12px", borderRadius: "6px", background: "rgba(59,130,246,0.15)", color: "#60a5fa", fontSize: "12px", fontWeight: 600 }}>
          {row.role}
        </span>
      ),
    },
    { key: "status", label: "Status" },
    { key: "created_at", label: "Member Since", render: (row) => new Date(row.created_at).toLocaleDateString() },
  ];

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />

      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <h1>Users & Access</h1>
            <p style={{ color: "var(--text-secondary)" }}>Manage users, roles, and access.</p>
          </div>
          <div className="header-actions" style={{ display: "flex", alignItems: "center", gap: "12px" }}>
            {error && (
              <span style={{ fontSize: "12px", color: "var(--warning)", background: "rgba(245,158,11,0.1)", padding: "6px 12px", borderRadius: "8px", border: "1px solid rgba(245,158,11,0.2)" }}>
                ⚠ API offline
              </span>
            )}
            <UserProfileWidget />
            <CreateUserButton />
          </div>
        </div>

        {/* Stats row */}
        <div className="metrics-grid animate-fade-in delay-1" style={{ gridTemplateColumns: "repeat(3, 1fr)", maxWidth: "600px", marginBottom: "28px" }}>
          <div className="glass-panel metric-card" style={{ padding: "20px" }}>
            <div className="metric-header"><span>Total Users</span></div>
            <div className="metric-value">{users.length}</div>
          </div>
          <div className="glass-panel metric-card" style={{ padding: "20px" }}>
            <div className="metric-header"><span>Active</span></div>
            <div className="metric-value" style={{ color: "var(--success)" }}>{activeCount}</div>
          </div>
          <div className="glass-panel metric-card" style={{ padding: "20px" }}>
            <div className="metric-header"><span>Pending / Inactive</span></div>
            <div className="metric-value">{users.length - activeCount}</div>
          </div>
        </div>

        <div className="data-section animate-fade-in delay-2">
          <DataTable
            columns={columns}
            data={users}
            rowKey={(row) => row.id}
            loading={loading}
            error={error || undefined}
            onRetry={load}
            emptyMessage="No users found."
            actions={(row) => (
              <div style={{ display: "flex", gap: "8px" }}>
                <button type="button" className="btn" onClick={() => openEdit(row)} style={{ padding: "6px 12px", fontSize: "12px" }}>
                  Edit
                </button>
                <button type="button" className="btn" onClick={() => setDeleting(row)} style={{ padding: "6px 12px", fontSize: "12px", color: "var(--danger)" }}>
                  Delete
                </button>
              </div>
            )}
          />
        </div>
      </main>

      {editing && (
        <div
          onClick={() => setEditing(null)}
          style={{
            position: "fixed", inset: 0, zIndex: 1000, background: "rgba(0,0,0,0.6)",
            display: "flex", alignItems: "center", justifyContent: "center", padding: "20px",
          }}
        >
          <div
            className="glass-panel animate-fade-in"
            onClick={(e) => e.stopPropagation()}
            style={{ width: "100%", maxWidth: "440px", padding: "28px", display: "flex", flexDirection: "column", gap: "16px" }}
          >
            <h3 style={{ fontSize: "17px", fontWeight: 600 }}>Edit User — {editing.email}</h3>
            <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Role</label>
              <select
                value={formRole}
                onChange={(e) => setFormRole(e.target.value)}
                style={{
                  background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)",
                  padding: "10px 12px", borderRadius: "8px", color: "white", outline: "none", fontSize: "14px",
                }}
              >
                {ROLE_OPTIONS.map((r) => (
                  <option key={r} value={r} style={{ background: "#13141c", color: "#fff" }}>{r}</option>
                ))}
              </select>
            </div>
            <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Status</label>
              <select
                value={formStatus}
                onChange={(e) => setFormStatus(e.target.value)}
                style={{
                  background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)",
                  padding: "10px 12px", borderRadius: "8px", color: "white", outline: "none", fontSize: "14px",
                }}
              >
                <option value="active" style={{ background: "#13141c", color: "#fff" }}>Active</option>
                <option value="invited" style={{ background: "#13141c", color: "#fff" }}>Invited</option>
                <option value="disabled" style={{ background: "#13141c", color: "#fff" }}>Disabled</option>
              </select>
            </div>
            <div style={{ display: "flex", justifyContent: "flex-end", gap: "10px", marginTop: "8px" }}>
              <button type="button" className="btn" onClick={() => setEditing(null)} disabled={submitting} style={{ padding: "10px 16px" }}>
                Cancel
              </button>
              <button type="button" className="btn btn-primary" onClick={submit} disabled={submitting} style={{ padding: "10px 16px" }}>
                {submitting ? "Saving…" : "Save"}
              </button>
            </div>
          </div>
        </div>
      )}

      <ConfirmDialog
        open={!!deleting}
        title="Delete user"
        message={`Are you sure you want to delete "${deleting?.email}"? This cannot be undone.`}
        confirmLabel="Delete"
        danger
        loading={submitting}
        onConfirm={confirmDelete}
        onCancel={() => setDeleting(null)}
      />
    </div>
  );
}
