"use client";

import { useEffect, useState } from "react";
import { useLocale } from "next-intl";
import Link from "next/link";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import DataTable, { type Column } from "@/components/shared/DataTable";
import { StatusBadge } from "@/components/shared/Badges";
import ConfirmDialog from "@/components/shared/ConfirmDialog";
import {
  listSecurityPolicies,
  createSecurityPolicy,
  deleteSecurityPolicy,
  activateSecurityPolicy,
  disableSecurityPolicy,
  type SecurityPolicy,
} from "@/lib/api";

export default function PoliciesPage() {
  const locale = useLocale();
  const [rows, setRows] = useState<SecurityPolicy[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [creating, setCreating] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [form, setForm] = useState({ name: "", description: "", enforcement_mode: "transparent" });
  const [deleting, setDeleting] = useState<SecurityPolicy | null>(null);

  const load = async () => {
    setLoading(true);
    setError("");
    try {
      setRows(await listSecurityPolicies());
    } catch {
      setError("Failed to load policies");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const submit = async () => {
    setSubmitting(true);
    setError("");
    try {
      await createSecurityPolicy(form.name, form.enforcement_mode);
      setCreating(false);
      setForm({ name: "", description: "", enforcement_mode: "transparent" });
      await load();
    } catch {
      setError("Failed to create policy");
    } finally {
      setSubmitting(false);
    }
  };

  const doActivate = async (id: string) => {
    try {
      await activateSecurityPolicy(id);
      await load();
    } catch {
      setError("Failed to activate policy");
    }
  };

  const doDisable = async (id: string) => {
    try {
      await disableSecurityPolicy(id);
      await load();
    } catch {
      setError("Failed to disable policy");
    }
  };

  const confirmDelete = async () => {
    if (!deleting) return;
    setSubmitting(true);
    try {
      await deleteSecurityPolicy(deleting.id);
      setDeleting(null);
      await load();
    } catch {
      setError("Failed to delete policy");
    } finally {
      setSubmitting(false);
    }
  };

  const columns: Column<SecurityPolicy>[] = [
    {
      key: "name",
      label: "Name",
      sortable: true,
      render: (row) => (
        <Link href={`/${locale}/policies/${row.id}`} style={{ color: "var(--accent-primary)", textDecoration: "none", fontWeight: 500 }}>
          {row.name}
        </Link>
      ),
    },
    { key: "enforcement_mode", label: "Mode", render: (row) => <StatusBadge value={row.enforcement_mode} /> },
    { key: "application_id", label: "Application", render: (row) => (row.application_id ? row.application_id.slice(0, 8) + "…" : "—") },
    { key: "version", label: "Version", render: (row) => `v${row.version}` },
    { key: "description", label: "Description", render: (row) => row.description || "—" },
  ];

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />

      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <h1>Security Policies</h1>
            <p style={{ color: "var(--text-secondary)" }}>Manage WAF security policies.</p>
          </div>
          <div className="header-actions" style={{ display: "flex", alignItems: "center", gap: "12px" }}>
            <button type="button" className="btn btn-primary" onClick={() => setCreating(true)}>
              + New Policy
            </button>
            <UserProfileWidget />
          </div>
        </div>

        <div className="data-section animate-fade-in delay-1">
          <DataTable
            columns={columns}
            data={rows}
            rowKey={(row) => row.id}
            loading={loading}
            error={error || undefined}
            onRetry={load}
            emptyMessage="No security policies yet."
            actions={(row) => (
              <div style={{ display: "flex", gap: "8px" }}>
                <button type="button" className="btn" onClick={() => doActivate(row.id)} style={{ padding: "6px 10px", fontSize: "12px", color: "var(--success)" }}>
                  Activate
                </button>
                <button type="button" className="btn" onClick={() => doDisable(row.id)} style={{ padding: "6px 10px", fontSize: "12px", color: "var(--warning)" }}>
                  Disable
                </button>
                <button type="button" className="btn" onClick={() => setDeleting(row)} style={{ padding: "6px 10px", fontSize: "12px", color: "var(--danger)" }}>
                  Delete
                </button>
              </div>
            )}
          />
        </div>
      </main>

      {creating && (
        <div
          onClick={() => setCreating(false)}
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
            <h3 style={{ fontSize: "17px", fontWeight: 600 }}>New Security Policy</h3>
            <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Name</label>
              <input
                type="text"
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
                style={{
                  background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)",
                  padding: "10px 12px", borderRadius: "8px", color: "white", outline: "none", fontSize: "14px",
                }}
              />
            </div>
            <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Description</label>
              <textarea
                value={form.description}
                onChange={(e) => setForm({ ...form, description: e.target.value })}
                rows={3}
                style={{
                  background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)",
                  padding: "10px 12px", borderRadius: "8px", color: "white", outline: "none", fontSize: "14px",
                  resize: "vertical",
                }}
              />
            </div>
            <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Enforcement Mode</label>
              <select
                value={form.enforcement_mode}
                onChange={(e) => setForm({ ...form, enforcement_mode: e.target.value })}
                style={{
                  background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)",
                  padding: "10px 12px", borderRadius: "8px", color: "white", outline: "none", fontSize: "14px",
                }}
              >
                <option value="transparent" style={{ background: "#13141c", color: "#fff" }}>Transparent</option>
                <option value="alarm" style={{ background: "#13141c", color: "#fff" }}>Alarm</option>
                <option value="blocking" style={{ background: "#13141c", color: "#fff" }}>Blocking</option>
              </select>
            </div>
            <div style={{ display: "flex", justifyContent: "flex-end", gap: "10px", marginTop: "8px" }}>
              <button type="button" className="btn" onClick={() => setCreating(false)} disabled={submitting} style={{ padding: "10px 16px" }}>
                Cancel
              </button>
              <button type="button" className="btn btn-primary" onClick={submit} disabled={submitting || !form.name} style={{ padding: "10px 16px" }}>
                {submitting ? "Creating…" : "Create Policy"}
              </button>
            </div>
          </div>
        </div>
      )}

      <ConfirmDialog
        open={!!deleting}
        title="Delete policy"
        message={`Are you sure you want to delete "${deleting?.name}"? This cannot be undone.`}
        confirmLabel="Delete"
        danger
        loading={submitting}
        onConfirm={confirmDelete}
        onCancel={() => setDeleting(null)}
      />
    </div>
  );
}
