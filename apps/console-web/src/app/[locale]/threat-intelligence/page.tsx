"use client";

import { useCallback, useEffect, useState } from "react";
import type { CSSProperties } from "react";
import { useLocale } from "next-intl";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import DataTable, { type Column } from "@/components/shared/DataTable";
import { StatusBadge } from "@/components/shared/Badges";
import ConfirmDialog from "@/components/shared/ConfirmDialog";
import {
  listThreatFeeds,
  createThreatFeed,
  deleteThreatFeed,
  syncThreatFeed,
  testThreatFeed,
  type ThreatFeed,
} from "@/lib/api";

export default function ThreatIntelligencePage() {
  const locale = useLocale();
  const [rows, setRows] = useState<ThreatFeed[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [creating, setCreating] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [form, setForm] = useState({ name: "", source: "", indicator_type: "ip", confidence: "low" });
  const [deleting, setDeleting] = useState<ThreatFeed | null>(null);
  const [notice, setNotice] = useState<{ kind: "ok" | "err"; msg: string } | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      setRows(await listThreatFeeds());
    } catch {
      setError("Failed to load threat feeds");
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
      await createThreatFeed({
        name: form.name,
        source: form.source,
        indicator_type: form.indicator_type,
        confidence: form.confidence,
        indicators: [],
      });
      setCreating(false);
      setForm({ name: "", source: "", indicator_type: "ip", confidence: "low" });
      await load();
    } catch {
      setError("Failed to create threat feed");
    } finally {
      setSubmitting(false);
    }
  };

  const doSync = async (row: ThreatFeed) => {
    try {
      await syncThreatFeed(row.id);
      setNotice({ kind: "ok", msg: `${row.name}: sync triggered` });
      await load();
    } catch {
      setNotice({ kind: "err", msg: `${row.name}: sync failed` });
    }
  };

  const doTest = async (row: ThreatFeed) => {
    try {
      await testThreatFeed(row.id);
      setNotice({ kind: "ok", msg: `${row.name}: test passed` });
    } catch {
      setNotice({ kind: "err", msg: `${row.name}: test failed` });
    }
  };

  const confirmDelete = async () => {
    if (!deleting) return;
    setSubmitting(true);
    try {
      await deleteThreatFeed(deleting.id);
      setDeleting(null);
      await load();
    } catch {
      setError("Failed to delete threat feed");
    } finally {
      setSubmitting(false);
    }
  };

  const columns: Column<ThreatFeed>[] = [
    { key: "name", label: "Name", sortable: true },
    { key: "source", label: "Source" },
    { key: "indicator_type", label: "Type", render: (row) => <StatusBadge value={row.indicator_type} /> },
    { key: "confidence", label: "Confidence", render: (row) => <StatusBadge value={row.confidence} /> },
    { key: "indicators", label: "Indicators", render: (row) => String(Array.isArray(row.indicators) ? row.indicators.length : 0) },
    { key: "status", label: "Status", render: (row) => <StatusBadge value={row.status} /> },
  ];

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />

      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <h1>Threat Intelligence</h1>
            <p style={{ color: "var(--text-secondary)" }}>Manage threat intelligence feeds.</p>
          </div>
          <div className="header-actions" style={{ display: "flex", alignItems: "center", gap: "12px" }}>
            <button type="button" className="btn btn-primary" onClick={() => setCreating(true)}>
              + New Feed
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
            emptyMessage="No threat feeds yet."
            actions={(row) => (
              <div style={{ display: "flex", gap: "6px" }}>
                <button type="button" className="btn" onClick={() => doTest(row)} style={{ padding: "6px 10px", fontSize: "12px" }}>
                  Test
                </button>
                <button type="button" className="btn" onClick={() => doSync(row)} style={{ padding: "6px 10px", fontSize: "12px" }}>
                  Sync
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
            <h3 style={{ fontSize: "17px", fontWeight: 600 }}>New Threat Feed</h3>
            <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Name</label>
              <input type="text" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} style={inputStyle} />
            </div>
            <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Source</label>
              <input type="text" value={form.source} onChange={(e) => setForm({ ...form, source: e.target.value })} style={inputStyle} />
            </div>
            <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Indicator Type</label>
              <select value={form.indicator_type} onChange={(e) => setForm({ ...form, indicator_type: e.target.value })} style={inputStyle}>
                {["ip", "domain", "url", "asn"].map((t) => (
                  <option key={t} value={t} style={optionStyle}>{t}</option>
                ))}
              </select>
            </div>
            <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Confidence</label>
              <select value={form.confidence} onChange={(e) => setForm({ ...form, confidence: e.target.value })} style={inputStyle}>
                {["low", "medium", "high"].map((c) => (
                  <option key={c} value={c} style={optionStyle}>{c}</option>
                ))}
              </select>
            </div>
            <div style={{ display: "flex", justifyContent: "flex-end", gap: "10px", marginTop: "8px" }}>
              <button type="button" className="btn" onClick={() => setCreating(false)} disabled={submitting} style={{ padding: "10px 16px" }}>
                Cancel
              </button>
              <button type="button" className="btn btn-primary" onClick={submit} disabled={submitting || !form.name || !form.source} style={{ padding: "10px 16px" }}>
                {submitting ? "Creating…" : "Create Feed"}
              </button>
            </div>
          </div>
        </div>
      )}

      <ConfirmDialog
        open={!!deleting}
        title="Delete threat feed"
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

const inputStyle: CSSProperties = {
  background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)",
  padding: "10px 12px", borderRadius: "8px", color: "white", outline: "none", fontSize: "14px",
};
const optionStyle: CSSProperties = { background: "#13141c", color: "#fff" };
