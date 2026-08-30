"use client";

import { useCallback, useEffect, useState } from "react";
import { useLocale, useTranslations } from "next-intl";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import DataTable, { type Column } from "@/components/shared/DataTable";
import { StatusBadge } from "@/components/shared/Badges";
import ConfirmDialog from "@/components/shared/ConfirmDialog";
import {
  listNotificationChannels,
  createNotificationChannel,
  deleteNotificationChannel,
  testNotificationChannel,
  type NotificationChannel,
} from "@/lib/api";

const KINDS = ["wazuh", "syslog", "cef", "leef", "webhook", "email", "teams", "slack", "soar"];

export default function NotificationChannelsPage() {
  const locale = useLocale();
  const t = useTranslations("notification_channels");
  const [rows, setRows] = useState<NotificationChannel[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [creating, setCreating] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [form, setForm] = useState({ name: "", kind: "webhook", is_default: false });
  const [deleting, setDeleting] = useState<NotificationChannel | null>(null);
  const [notice, setNotice] = useState<{ kind: "ok" | "err"; msg: string } | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      setRows(await listNotificationChannels());
    } catch {
      setError("Failed to load notification channels");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const submit = async () => {
    setSubmitting(true);
    setError("");
    setNotice(null);
    try {
      await createNotificationChannel({ name: form.name, kind: form.kind, is_default: form.is_default });
      setCreating(false);
      setForm({ name: "", kind: "webhook", is_default: false });
      await load();
    } catch {
      setError("Failed to create notification channel");
    } finally {
      setSubmitting(false);
    }
  };

  const doTest = async (row: NotificationChannel) => {
    try {
      const r = await testNotificationChannel(row.id);
      setNotice({ kind: "ok", msg: `${row.name}: test ${r.success === "true" ? "passed" : "failed"}` });
    } catch {
      setNotice({ kind: "err", msg: `${row.name}: test failed` });
    }
  };

  const confirmDelete = async () => {
    if (!deleting) return;
    setSubmitting(true);
    try {
      await deleteNotificationChannel(deleting.id);
      setDeleting(null);
      await load();
    } catch {
      setError("Failed to delete notification channel");
    } finally {
      setSubmitting(false);
    }
  };

  const columns: Column<NotificationChannel>[] = [
    { key: "name", label: "Name", sortable: true },
    { key: "kind", label: "Kind", render: (row) => <StatusBadge value={row.kind} /> },
    { key: "is_default", label: "Default", render: (row) => (row.is_default ? "✓" : "—") },
    { key: "enabled", label: "Enabled", render: (row) => (row.enabled ? <StatusBadge value="enabled" /> : <StatusBadge value="disabled" />) },
    { key: "created_at", label: "Created", render: (row) => new Date(row.created_at).toLocaleDateString() },
  ];

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />

      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <h1>{t("title")}</h1>
            <p style={{ color: "var(--text-secondary)" }}>Configure alert delivery channels.</p>
          </div>
          <div className="header-actions" style={{ display: "flex", alignItems: "center", gap: "12px" }}>
            <button type="button" className="btn btn-primary" onClick={() => setCreating(true)}>
              + New Channel
            </button>
            <UserProfileWidget />
          </div>
        </div>

        {notice && (
          <div
            className="glass-panel animate-fade-in"
            style={{
              padding: "12px 16px", marginBottom: "16px", fontSize: "13px",
              color: notice.kind === "ok" ? "var(--success)" : "var(--danger)",
              border: `1px solid ${notice.kind === "ok" ? "rgba(16,185,129,0.3)" : "rgba(239,68,68,0.3)"}`,
            }}
          >
            {notice.msg}
          </div>
        )}

        <div className="data-section animate-fade-in delay-1">
          <DataTable
            columns={columns}
            data={rows}
            rowKey={(row) => row.id}
            loading={loading}
            error={error || undefined}
            onRetry={load}
            emptyMessage="No notification channels yet."
            actions={(row) => (
              <div style={{ display: "flex", gap: "6px" }}>
                <button type="button" className="btn" onClick={() => doTest(row)} style={{ padding: "6px 10px", fontSize: "12px" }}>
                  Test
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
            <h3 style={{ fontSize: "17px", fontWeight: 600 }}>New Notification Channel</h3>
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
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Kind</label>
              <select
                value={form.kind}
                onChange={(e) => setForm({ ...form, kind: e.target.value })}
                style={{
                  background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)",
                  padding: "10px 12px", borderRadius: "8px", color: "white", outline: "none", fontSize: "14px",
                }}
              >
                {KINDS.map((k) => (
                  <option key={k} value={k} style={{ background: "#13141c", color: "#fff" }}>{k}</option>
                ))}
              </select>
            </div>
            <div style={{ display: "flex", alignItems: "center", gap: "8px" }}>
              <input
                type="checkbox"
                checked={form.is_default}
                onChange={(e) => setForm({ ...form, is_default: e.target.checked })}
                style={{ width: "16px", height: "16px" }}
              />
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Set as default channel</label>
            </div>
            <div style={{ display: "flex", justifyContent: "flex-end", gap: "10px", marginTop: "8px" }}>
              <button type="button" className="btn" onClick={() => setCreating(false)} disabled={submitting} style={{ padding: "10px 16px" }}>
                Cancel
              </button>
              <button type="button" className="btn btn-primary" onClick={submit} disabled={submitting || !form.name} style={{ padding: "10px 16px" }}>
                {submitting ? "Creating…" : "Create Channel"}
              </button>
            </div>
          </div>
        </div>
      )}

      <ConfirmDialog
        open={!!deleting}
        title="Delete notification channel"
        message={`Are you sure you want to delete "${deleting?.name}"?`}
        confirmLabel="Delete"
        danger
        loading={submitting}
        onConfirm={confirmDelete}
        onCancel={() => setDeleting(null)}
      />
    </div>
  );
}
