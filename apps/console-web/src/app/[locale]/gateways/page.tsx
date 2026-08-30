"use client";

import { useEffect, useState } from "react";
import { useLocale, useTranslations } from "next-intl";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import DataTable, { type Column } from "@/components/shared/DataTable";
import { StatusBadge } from "@/components/shared/Badges";
import ConfirmDialog from "@/components/shared/ConfirmDialog";
import { usePermission } from "@/hooks/usePermission";
import {
  listGateways,
  updateGateway,
  deleteGateway,
  type Gateway,
} from "@/lib/api";

export default function GatewaysPage() {
  const locale = useLocale();
  const t = useTranslations("gateways");
  const [rows, setRows] = useState<Gateway[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [editing, setEditing] = useState<Gateway | null>(null);
  const [deleting, setDeleting] = useState<Gateway | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [form, setForm] = useState({ hostname: "", ip: "" });
  const { can: canUser } = usePermission();

  const load = async () => {
    setLoading(true);
    setError("");
    try {
      setRows(await listGateways());
    } catch {
      setError("Failed to load gateways");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const openEdit = (g: Gateway) => {
    setForm({ hostname: g.hostname, ip: g.ip ?? "" });
    setEditing(g);
  };

  const submit = async () => {
    if (!editing) return;
    setSubmitting(true);
    setError("");
    try {
      await updateGateway(editing.id, { hostname: form.hostname, ip: form.ip });
      setEditing(null);
      await load();
    } catch {
      setError("Failed to update gateway");
    } finally {
      setSubmitting(false);
    }
  };

  const confirmDelete = async () => {
    if (!deleting) return;
    setSubmitting(true);
    try {
      await deleteGateway(deleting.id);
      setDeleting(null);
      await load();
    } catch {
      setError("Failed to delete gateway");
    } finally {
      setSubmitting(false);
    }
  };

  const columns: Column<Gateway>[] = [
    { key: "hostname", label: "Hostname", sortable: true },
    { key: "ip", label: "IP", render: (row) => row.ip || "—" },
    { key: "status", label: "Status", render: (row) => <StatusBadge value={row.status} /> },
    { key: "version", label: "Version", render: (row) => row.version || "—" },
    {
      key: "last_seen_at",
      label: "Last Seen",
      render: (row) => (row.last_seen_at ? new Date(row.last_seen_at).toLocaleString() : "—"),
    },
    {
      key: "applied_hash",
      label: "Applied Hash",
      render: (row) => (row.applied_hash ? `${row.applied_hash.slice(0, 12)}…` : "—"),
    },
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
          <div className="header-actions" style={{ display: "flex", alignItems: "center", gap: "12px" }}>
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
            emptyMessage="No gateways registered."
            actions={(row) => (
              <div style={{ display: "flex", gap: "8px" }}>
                {canUser("edit") && (
                  <button type="button" className="btn" onClick={() => openEdit(row)} style={{ padding: "6px 12px", fontSize: "12px" }}>
                    Edit
                  </button>
                )}
                {canUser("delete") && (
                  <button type="button" className="btn" onClick={() => setDeleting(row)} style={{ padding: "6px 12px", fontSize: "12px", color: "var(--danger)" }}>
                    Delete
                  </button>
                )}
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
            <h3 style={{ fontSize: "17px", fontWeight: 600 }}>Edit Gateway</h3>
            <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Hostname</label>
              <input
                type="text"
                value={form.hostname}
                onChange={(e) => setForm({ ...form, hostname: e.target.value })}
                style={{
                  background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)",
                  padding: "10px 12px", borderRadius: "8px", color: "white", outline: "none", fontSize: "14px",
                }}
              />
            </div>
            <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>IP Address</label>
              <input
                type="text"
                value={form.ip}
                onChange={(e) => setForm({ ...form, ip: e.target.value })}
                style={{
                  background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)",
                  padding: "10px 12px", borderRadius: "8px", color: "white", outline: "none", fontSize: "14px",
                }}
              />
            </div>
            <div style={{ display: "flex", justifyContent: "flex-end", gap: "10px", marginTop: "8px" }}>
              <button type="button" className="btn" onClick={() => setEditing(null)} disabled={submitting} style={{ padding: "10px 16px" }}>
                Cancel
              </button>
              <button type="button" className="btn btn-primary" onClick={submit} disabled={submitting || !form.hostname} style={{ padding: "10px 16px" }}>
                {submitting ? "Saving…" : "Save"}
              </button>
            </div>
          </div>
        </div>
      )}

      <ConfirmDialog
        open={!!deleting}
        title="Delete gateway"
        message={`Are you sure you want to delete "${deleting?.hostname}"? This cannot be undone.`}
        confirmLabel="Delete"
        danger
        loading={submitting}
        onConfirm={confirmDelete}
        onCancel={() => setDeleting(null)}
      />
    </div>
  );
}
