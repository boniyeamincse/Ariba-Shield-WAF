"use client";

import { useCallback, useEffect, useState } from "react";
import { useLocale, useTranslations } from "next-intl";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import DataTable, { type Column } from "@/components/shared/DataTable";
import { SeverityBadge, StatusBadge } from "@/components/shared/Badges";
import ConfirmDialog from "@/components/shared/ConfirmDialog";
import FilterBar from "@/components/shared/FilterBar";
import {
  listIncidents,
  createIncident,
  deleteIncident,
  closeIncident,
  reopenIncident,
  escalateIncident,
  type Incident,
} from "@/lib/api";

export default function IncidentsPage() {
  const locale = useLocale();
  const t = useTranslations("incidents");
  const ct = useTranslations("common");
  const [rows, setRows] = useState<Incident[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [creating, setCreating] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [form, setForm] = useState({ title: "", severity: "medium", notes: "" });
  const [filters, setFilters] = useState<Record<string, string>>({});
  const [deleting, setDeleting] = useState<Incident | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      let data = await listIncidents();
      if (filters.status) data = data.filter((i) => i.status === filters.status);
      if (filters.severity) data = data.filter((i) => i.severity === filters.severity);
      setRows(data);
    } catch {
      setError("Failed to load incidents");
    } finally {
      setLoading(false);
    }
  }, [filters]);

  useEffect(() => {
    load();
  }, [load]);

  const submit = async () => {
    setSubmitting(true);
    setError("");
    try {
      await createIncident({ title: form.title, severity: form.severity, notes: form.notes });
      setCreating(false);
      setForm({ title: "", severity: "medium", notes: "" });
      await load();
    } catch {
      setError("Failed to create incident");
    } finally {
      setSubmitting(false);
    }
  };

  const doAction = async (action: "close" | "reopen" | "escalate", row: Incident) => {
    try {
      if (action === "close") await closeIncident(row.id);
      if (action === "reopen") await reopenIncident(row.id);
      if (action === "escalate") await escalateIncident(row.id);
      await load();
    } catch {
      setError(`Failed to ${action} incident`);
    }
  };

  const confirmDelete = async () => {
    if (!deleting) return;
    setSubmitting(true);
    try {
      await deleteIncident(deleting.id);
      setDeleting(null);
      await load();
    } catch {
      setError("Failed to delete incident");
    } finally {
      setSubmitting(false);
    }
  };

  const columns: Column<Incident>[] = [
    { key: "title", label: "Title", sortable: true },
    { key: "severity", label: "Severity", render: (row) => <SeverityBadge value={row.severity} /> },
    { key: "status", label: "Status", render: (row) => <StatusBadge value={row.status} /> },
    { key: "owner_user_id", label: "Owner", render: (row) => (row.owner_user_id ? row.owner_user_id.slice(0, 8) + "…" : "—") },
    { key: "related_events", label: "Events", render: (row) => String(row.related_events?.length ?? 0) },
    { key: "created_at", label: "Created", render: (row) => new Date(row.created_at).toLocaleDateString() },
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
<button type="button" className="btn btn-primary" onClick={() => setCreating(true)}>
                + {t("new_incident")}
            </button>
            <UserProfileWidget />
          </div>
        </div>

        <FilterBar
          filters={[
            { type: "select", key: "status", label: "Status", options: [
              { value: "open", label: "Open" },
              { value: "investigating", label: "Investigating" },
              { value: "resolved", label: "Resolved" },
            ]},
            { type: "select", key: "severity", label: "Severity", options: [
              { value: "critical", label: "Critical" },
              { value: "high", label: "High" },
              { value: "medium", label: "Medium" },
              { value: "low", label: "Low" },
            ]},
          ]}
          values={filters}
          onChange={setFilters}
        />

        <div className="data-section animate-fade-in delay-1">
          <DataTable
            columns={columns}
            data={rows}
            rowKey={(row) => row.id}
            loading={loading}
            error={error || undefined}
            onRetry={load}
            emptyMessage="No incidents yet."
            actions={(row) => (
              <div style={{ display: "flex", gap: "6px" }}>
                {row.status === "resolved" ? (
                  <button type="button" className="btn" onClick={() => doAction("reopen", row)} style={{ padding: "6px 10px", fontSize: "12px", color: "var(--warning)" }}>
                    {ct("reopen")}
                  </button>
                ) : (
                  <>
                    <button type="button" className="btn" onClick={() => doAction("escalate", row)} style={{ padding: "6px 10px", fontSize: "12px", color: "var(--danger)" }}>
                      {ct("escalate")}
                    </button>
                    <button type="button" className="btn" onClick={() => doAction("close", row)} style={{ padding: "6px 10px", fontSize: "12px", color: "var(--success)" }}>
                      {ct("close")}
                    </button>
                  </>
                )}
                <button type="button" className="btn" onClick={() => setDeleting(row)} style={{ padding: "6px 10px", fontSize: "12px", color: "var(--danger)" }}>
                  {ct("delete")}
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
            <h3 style={{ fontSize: "17px", fontWeight: 600 }}>New Incident</h3>
            <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Title</label>
              <input
                type="text"
                value={form.title}
                onChange={(e) => setForm({ ...form, title: e.target.value })}
                style={{
                  background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)",
                  padding: "10px 12px", borderRadius: "8px", color: "white", outline: "none", fontSize: "14px",
                }}
              />
            </div>
            <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Severity</label>
              <select
                value={form.severity}
                onChange={(e) => setForm({ ...form, severity: e.target.value })}
                style={{
                  background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)",
                  padding: "10px 12px", borderRadius: "8px", color: "white", outline: "none", fontSize: "14px",
                }}
              >
                <option value="critical" style={{ background: "#13141c", color: "#fff" }}>Critical</option>
                <option value="high" style={{ background: "#13141c", color: "#fff" }}>High</option>
                <option value="medium" style={{ background: "#13141c", color: "#fff" }}>Medium</option>
                <option value="low" style={{ background: "#13141c", color: "#fff" }}>Low</option>
              </select>
            </div>
            <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Notes</label>
              <textarea
                value={form.notes}
                onChange={(e) => setForm({ ...form, notes: e.target.value })}
                rows={3}
                style={{
                  background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)",
                  padding: "10px 12px", borderRadius: "8px", color: "white", outline: "none", fontSize: "14px",
                  resize: "vertical",
                }}
              />
            </div>
            <div style={{ display: "flex", justifyContent: "flex-end", gap: "10px", marginTop: "8px" }}>
              <button type="button" className="btn" onClick={() => setCreating(false)} disabled={submitting} style={{ padding: "10px 16px" }}>
                {ct("cancel")}
              </button>
              <button type="button" className="btn btn-primary" onClick={submit} disabled={submitting || !form.title} style={{ padding: "10px 16px" }}>
                {submitting ? "Creating…" : "Create Incident"}
              </button>
            </div>
          </div>
        </div>
      )}

      <ConfirmDialog
        open={!!deleting}
        title="Delete incident"
        message={`Are you sure you want to delete "${deleting?.title}"? This cannot be undone.`}
        confirmLabel="Delete"
        danger
        loading={submitting}
        onConfirm={confirmDelete}
        onCancel={() => setDeleting(null)}
      />
    </div>
  );
}
