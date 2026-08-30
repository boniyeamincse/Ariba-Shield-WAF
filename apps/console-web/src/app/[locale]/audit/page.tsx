"use client";

import { useCallback, useEffect, useState } from "react";
import { useLocale, useTranslations } from "next-intl";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import DataTable, { type Column } from "@/components/shared/DataTable";
import FilterBar from "@/components/shared/FilterBar";
import { listAuditEvents, exportAuditEvents, type AuditEvent } from "@/lib/api";

const PAGE_SIZE = 50;

export default function AuditLogPage() {
  const locale = useLocale();
  const t = useTranslations("audit");
  const [rows, setRows] = useState<AuditEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [filters, setFilters] = useState<Record<string, string>>({});
  const [exporting, setExporting] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const events = await listAuditEvents();
      // Backend returns all events; filter client-side for the current view.
      const filtered = events.filter((ev) => {
        if (filters.action && ev.action !== filters.action) return false;
        return true;
      });
      setRows(filtered.slice(0, PAGE_SIZE));
    } catch {
      setError("Failed to load audit events");
    } finally {
      setLoading(false);
    }
  }, [filters]);

  useEffect(() => {
    load();
  }, [load]);

  const doExport = async () => {
    setExporting(true);
    try {
      const blob = await exportAuditEvents();
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = `audit-events-${new Date().toISOString().slice(0, 10)}.csv`;
      a.click();
      URL.revokeObjectURL(url);
    } catch {
      setError("Failed to export audit events");
    } finally {
      setExporting(false);
    }
  };

  const columns: Column<AuditEvent>[] = [
    { key: "created_at", label: "Time", render: (row) => new Date(row.created_at).toLocaleString() },
    {
      key: "action",
      label: "Action",
      render: (row) => {
        const colors: Record<string, { bg: string; text: string }> = {
          POST: { bg: "rgba(59,130,246,0.15)", text: "#60a5fa" },
          PUT: { bg: "rgba(168,85,247,0.15)", text: "#c084fc" },
          PATCH: { bg: "rgba(245,158,11,0.15)", text: "#fbbf24" },
          DELETE: { bg: "rgba(239,68,68,0.15)", text: "#f87171" },
          GET: { bg: "rgba(16,185,129,0.15)", text: "#34d399" },
        };
        const c = colors[row.action] ?? { bg: "rgba(255,255,255,0.06)", text: "#9ca3af" };
        return (
          <span style={{ padding: "4px 10px", borderRadius: "6px", background: c.bg, color: c.text, fontSize: "12px", fontWeight: 600 }}>
            {row.action}
          </span>
        );
      },
    },
    { key: "resource", label: "Resource" },
    { key: "resource_id", label: "Resource ID", render: (row) => <span style={{ fontFamily: "monospace", fontSize: "12px" }}>{row.resource_id || "—"}</span> },
    { key: "actor_user_id", label: "Actor", render: (row) => (row.actor_user_id ? `${row.actor_user_id.slice(0, 12)}…` : "system") },
    { key: "ip", label: "IP", render: (row) => row.ip || "—" },
    { key: "request_id", label: "Request ID", render: (row) => (row.request_id ? `${row.request_id.slice(0, 12)}…` : "—") },
  ];

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />

      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <h1>{t("title")}</h1>
            <p style={{ color: "var(--text-secondary)" }}>Immutable record of all management actions.</p>
          </div>
          <div className="header-actions" style={{ display: "flex", alignItems: "center", gap: "12px" }}>
            <button type="button" className="btn" onClick={doExport} disabled={exporting} style={{ padding: "8px 14px", fontSize: "13px" }}>
              {exporting ? "Exporting…" : "Export CSV"}
            </button>
            <UserProfileWidget />
          </div>
        </div>

        <div
          className="glass-panel animate-fade-in delay-1"
          style={{ padding: "16px 20px", marginBottom: "20px", display: "flex", alignItems: "center", gap: "12px" }}
        >
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="var(--success)" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ flexShrink: 0 }}>
            <rect x="3" y="11" width="18" height="11" rx="2"></rect>
            <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
          </svg>
          <div style={{ fontSize: "13px", color: "var(--text-secondary)" }}>This log is append-only and cannot be modified.</div>
        </div>

        <FilterBar
          filters={[
            { type: "select", key: "action", label: "Action", options: [
              { value: "POST", label: "POST" },
              { value: "PUT", label: "PUT" },
              { value: "PATCH", label: "PATCH" },
              { value: "DELETE", label: "DELETE" },
            ]},
          ]}
          values={filters}
          onChange={setFilters}
        />

        <div className="data-section animate-fade-in delay-2">
          <DataTable
            columns={columns}
            data={rows}
            rowKey={(row) => row.id}
            loading={loading}
            error={error || undefined}
            onRetry={load}
            emptyMessage="No audit events yet."
          />
        </div>
      </main>
    </div>
  );
}
