"use client";

import { useCallback, useEffect, useState } from "react";
import { useLocale } from "next-intl";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import DataTable, { type Column } from "@/components/shared/DataTable";
import { StatusBadge } from "@/components/shared/Badges";
import ConfirmDialog from "@/components/shared/ConfirmDialog";
import { listIntegrations, createIntegration, deleteIntegration, testIntegration, setIntegrationEnabled, type Integration } from "@/lib/api";

const TYPE_OPTIONS = ["splunk_hec", "wazuh", "syslog", "webhook", "teams", "slack"];

export default function IntegrationsPage() {
  const locale = useLocale();
  const [rows, setRows] = useState<Integration[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [creating, setCreating] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [form, setForm] = useState({ name: "", type: "webhook", endpoint: "" });
  const [deleting, setDeleting] = useState<Integration | null>(null);
  const [testResult, setTestResult] = useState<{ id: string; message: string } | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      setRows(await listIntegrations());
    } catch {
      setError("Failed to load integrations");
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
    try {
      await createIntegration({ name: form.name, type: form.type, endpoint: form.endpoint, enabled: true });
      setCreating(false);
      setForm({ name: "", type: "webhook", endpoint: "" });
      await load();
    } catch {
      setError("Failed to create integration");
    } finally {
      setSubmitting(false);
    }
  };

  const doTest = async (row: Integration) => {
    try {
      const r = await testIntegration(row.id);
      setTestResult({ id: row.id, message: r.success === "true" ? "Test passed" : "Test failed" });
    } catch {
      setTestResult({ id: row.id, message: "Test failed" });
    }
  };

  const doToggle = async (row: Integration) => {
    try {
      await setIntegrationEnabled(row.id, !row.enabled);
      await load();
    } catch {
      setError("Failed to update integration");
    }
  };

  const columns: Column<Integration>[] = [
    { key: "name", label: "Name", sortable: true },
    { key: "type", label: "Type", render: (row) => <StatusBadge value={row.type} /> },
    { key: "endpoint", label: "Endpoint", render: (row) => row.endpoint || "—" },
    {
      key: "log_types",
      label: "Log Types",
      render: (row) => (row.log_types?.length ? row.log_types.slice(0, 3).join(", ") : "—"),
    },
    { key: "enabled", label: "Enabled", render: (row) => (row.enabled ? <StatusBadge value="enabled" /> : <StatusBadge value="disabled" />) },
  ];

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />

      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <h1>Integrations</h1>
            <p style={{ color: "var(--text-secondary)" }}>SIEM, logging, and notification integrations.</p>
          </div>
          <div className="header-actions" style={{ display: "flex", alignItems: "center", gap: "12px" }}>
            <button type="button" className="btn btn-primary" onClick={() => setCreating(true)}>
              + New Integration
            </button>
            <UserProfileWidget />
          </div>
        </div>

        {testResult && (
          <div
            className="glass-panel animate-fade-in"
            style={{
              padding: "12px 16px", marginBottom: "16px", fontSize: "13px",
              color: testResult.message === "Test passed" ? "var(--success)" : "var(--danger)",
              border: `1px solid ${testResult.message === "Test passed" ? "rgba(16,185,129,0.3)" : "rgba(239,68,68,0.3)"}`,
            }}
          >
            Integration {testResult.id.slice(0, 8)}: {testResult.message}
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
            emptyMessage="No integrations yet."
            actions={(row) => (
              <div style={{ display: "flex", gap: "6px" }}>
                <button type="button" className="btn" onClick={() => doTest(row)} style={{ padding: "6px 10px", fontSize: "12px" }}>
                  Test
                </button>
                <button type="button" className="btn" onClick={() => doToggle(row)} style={{ padding: "6px 10px", fontSize: "12px", color: row.enabled ? "var(--warning)" : "var(--success)" }}>
                  {row.enabled ? "Disable" : "Enable"}
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
            <h3 style={{ fontSize: "17px", fontWeight: 600 }}>New Integration</h3>
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
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Type</label>
              <select
                value={form.type}
                onChange={(e) => setForm({ ...form, type: e.target.value })}
                style={{
                  background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)",
                  padding: "10px 12px", borderRadius: "8px", color: "white", outline: "none", fontSize: "14px",
                }}
              >
                {TYPE_OPTIONS.map((t) => (
                  <option key={t} value={t} style={{ background: "#13141c", color: "#fff" }}>{t}</option>
                ))}
              </select>
            </div>
            <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Endpoint</label>
              <input
                type="text"
                value={form.endpoint}
                onChange={(e) => setForm({ ...form, endpoint: e.target.value })}
                style={{
                  background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)",
                  padding: "10px 12px", borderRadius: "8px", color: "white", outline: "none", fontSize: "14px",
                }}
              />
            </div>
            <div style={{ display: "flex", justifyContent: "flex-end", gap: "10px", marginTop: "8px" }}>
              <button type="button" className="btn" onClick={() => setCreating(false)} disabled={submitting} style={{ padding: "10px 16px" }}>
                Cancel
              </button>
              <button type="button" className="btn btn-primary" onClick={submit} disabled={submitting || !form.name} style={{ padding: "10px 16px" }}>
                {submitting ? "Creating…" : "Create Integration"}
              </button>
            </div>
          </div>
        </div>
      )}

      <ConfirmDialog
        open={!!deleting}
        title="Delete integration"
        message={`Are you sure you want to delete "${deleting?.name}"?`}
        confirmLabel="Delete"
        danger
        loading={submitting}
        onConfirm={async () => {
          if (!deleting) return;
          setSubmitting(true);
          try {
            await deleteIntegration(deleting.id);
            setDeleting(null);
            await load();
          } catch {
            setError("Failed to delete integration");
          } finally {
            setSubmitting(false);
          }
        }}
        onCancel={() => setDeleting(null)}
      />
    </div>
  );
}
