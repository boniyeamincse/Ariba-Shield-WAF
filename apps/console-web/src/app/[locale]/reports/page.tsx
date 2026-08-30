"use client";

import { useCallback, useEffect, useState } from "react";
import { useLocale, useTranslations } from "next-intl";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import DataTable, { type Column } from "@/components/shared/DataTable";
import { StatusBadge } from "@/components/shared/Badges";
import ConfirmDialog from "@/components/shared/ConfirmDialog";
import {
  listReports,
  getReport,
  deleteReport,
  generateReport,
  type Report,
} from "@/lib/api";

const REPORT_TYPES = [
  { kind: "security" as const, label: "Security" },
  { kind: "traffic" as const, label: "Traffic" },
  { kind: "incidents" as const, label: "Incidents" },
  { kind: "compliance" as const, label: "Compliance" },
];

export default function ReportsPage() {
  const locale = useLocale();
  const t = useTranslations("reports");
  const [rows, setRows] = useState<Report[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [generating, setGenerating] = useState<string | null>(null);
  const [deleting, setDeleting] = useState<Report | null>(null);
  const [detail, setDetail] = useState<{ id: string; summary: unknown } | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      setRows(await listReports());
    } catch {
      setError("Failed to load reports");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const doGenerate = async (kind: "security" | "traffic" | "incidents" | "compliance") => {
    setGenerating(kind);
    setError("");
    try {
      await generateReport(kind);
      await load();
    } catch {
      setError(`Failed to generate ${kind} report`);
    } finally {
      setGenerating(null);
    }
  };

  const doView = async (row: Report) => {
    try {
      const r = await getReport(row.id);
      setDetail({ id: r.id, summary: r.summary });
    } catch {
      setError("Failed to load report detail");
    }
  };

  const confirmDelete = async () => {
    if (!deleting) return;
    setSubmitting(true);
    try {
      await deleteReport(deleting.id);
      setDeleting(null);
      await load();
    } catch {
      setError("Failed to delete report");
    } finally {
      setSubmitting(false);
    }
  };

  const downloadReport = (row: Report) => {
    const base = process.env.NEXT_PUBLIC_API_BASE ?? (window.location.protocol + "//" + window.location.hostname + ":8443");
    window.open(`${base}/api/v1/reports/${row.id}/download`, "_blank");
  };

  const columns: Column<Report>[] = [
    { key: "name", label: "Name", sortable: true },
    { key: "kind", label: "Kind", render: (row) => <StatusBadge value={row.kind} /> },
    { key: "status", label: "Status", render: (row) => <StatusBadge value={row.status} /> },
    { key: "created_by", label: "Created By", render: (row) => row.created_by || "—" },
    { key: "created_at", label: "Created", render: (row) => new Date(row.created_at).toLocaleString() },
  ];

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />

      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <h1>{t("title")}</h1>
            <p style={{ color: "var(--text-secondary)" }}>Generate security, traffic, incident, and compliance reports.</p>
          </div>
          <div className="header-actions" style={{ display: "flex", alignItems: "center", gap: "12px" }}>
            <UserProfileWidget />
          </div>
        </div>

        {/* Generate buttons */}
        <div className="metrics-grid animate-fade-in delay-1" style={{ marginBottom: "24px" }}>
          {REPORT_TYPES.map((rt) => (
            <div key={rt.kind} className="metric-card glass-panel" style={{ padding: "20px", gap: "12px" }}>
              <div className="metric-header"><span style={{ textTransform: "capitalize" }}>{rt.label} Report</span></div>
              <button
                type="button"
                className="btn btn-primary"
                onClick={() => doGenerate(rt.kind)}
                disabled={generating !== null}
                style={{ padding: "10px 16px", fontSize: "13px" }}
              >
                {generating === rt.kind ? "Generating…" : `Generate ${rt.label}`}
              </button>
            </div>
          ))}
        </div>

        <div className="data-section animate-fade-in delay-2">
          <DataTable
            columns={columns}
            data={rows}
            rowKey={(row) => row.id}
            loading={loading}
            error={error || undefined}
            onRetry={load}
            emptyMessage="No reports generated yet."
            actions={(row) => (
              <div style={{ display: "flex", gap: "6px" }}>
                <button type="button" className="btn" onClick={() => doView(row)} style={{ padding: "6px 10px", fontSize: "12px" }}>
                  View
                </button>
                <button type="button" className="btn" onClick={() => downloadReport(row)} style={{ padding: "6px 10px", fontSize: "12px" }}>
                  Download
                </button>
                <button type="button" className="btn" onClick={() => setDeleting(row)} style={{ padding: "6px 10px", fontSize: "12px", color: "var(--danger)" }}>
                  Delete
                </button>
              </div>
            )}
          />
        </div>
      </main>

      {detail && (
        <div
          onClick={() => setDetail(null)}
          style={{
            position: "fixed", inset: 0, zIndex: 1000, background: "rgba(0,0,0,0.6)",
            display: "flex", alignItems: "center", justifyContent: "center", padding: "20px",
          }}
        >
          <div
            className="glass-panel animate-fade-in"
            onClick={(e) => e.stopPropagation()}
            style={{ width: "100%", maxWidth: "640px", maxHeight: "70vh", overflow: "auto", padding: "28px", display: "flex", flexDirection: "column", gap: "16px" }}
          >
            <h3 style={{ fontSize: "17px", fontWeight: 600 }}>Report {detail.id.slice(0, 8)}</h3>
            <pre
              style={{
                background: "rgba(255,255,255,0.04)", border: "1px solid rgba(255,255,255,0.08)",
                padding: "16px", borderRadius: "8px", fontSize: "12px", color: "var(--text-secondary)",
                overflow: "auto", whiteSpace: "pre-wrap", wordBreak: "break-word",
              }}
            >
              {JSON.stringify(detail.summary, null, 2)}
            </pre>
            <div style={{ display: "flex", justifyContent: "flex-end" }}>
              <button type="button" className="btn" onClick={() => setDetail(null)} style={{ padding: "10px 16px" }}>
                Close
              </button>
            </div>
          </div>
        </div>
      )}

      <ConfirmDialog
        open={!!deleting}
        title="Delete report"
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
