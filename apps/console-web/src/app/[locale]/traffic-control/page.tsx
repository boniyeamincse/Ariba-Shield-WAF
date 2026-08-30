"use client";

import { useCallback, useEffect, useState } from "react";
import { useLocale } from "next-intl";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import DataTable, { type Column } from "@/components/shared/DataTable";
import { StatusBadge } from "@/components/shared/Badges";
import ConfirmDialog from "@/components/shared/ConfirmDialog";
import {
  listIPLists,
  createIPList,
  deleteIPList,
  listRateLimits,
  createRateLimit,
  deleteRateLimit,
  type IPList,
  type RateLimit,
} from "@/lib/api";

type Tab = "ip-lists" | "rate-limits";

export default function TrafficControlPage() {
  const locale = useLocale();
  const [tab, setTab] = useState<Tab>("ip-lists");

  // IP lists
  const [ipLists, setIpLists] = useState<IPList[]>([]);
  const [ipLoading, setIpLoading] = useState(true);
  const [ipError, setIpError] = useState("");
  const [creatingList, setCreatingList] = useState(false);
  const [listForm, setListForm] = useState({ name: "", list_type: "allow", description: "" });
  const [deletingList, setDeletingList] = useState<IPList | null>(null);

  // Rate limits
  const [rateLimits, setRateLimits] = useState<RateLimit[]>([]);
  const [rlLoading, setRlLoading] = useState(false);
  const [rlError, setRlError] = useState("");
  const [creatingRL, setCreatingRL] = useState(false);
  const [rlForm, setRlForm] = useState({ name: "", route_prefix: "/", limit_count: "100", window_seconds: "60", action: "throttle" });
  const [deletingRL, setDeletingRL] = useState<RateLimit | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const loadIPs = useCallback(async () => {
    setIpLoading(true);
    setIpError("");
    try {
      setIpLists(await listIPLists());
    } catch {
      setIpError("Failed to load IP lists");
    } finally {
      setIpLoading(false);
    }
  }, []);

  const loadRates = useCallback(async () => {
    setRlLoading(true);
    setRlError("");
    try {
      setRateLimits(await listRateLimits());
    } catch {
      setRlError("Failed to load rate limits");
    } finally {
      setRlLoading(false);
    }
  }, []);

  useEffect(() => {
    loadIPs();
  }, [loadIPs]);

  useEffect(() => {
    if (tab === "rate-limits") loadRates();
  }, [tab, loadRates]);

  const submitList = async () => {
    setSubmitting(true);
    try {
      await createIPList({ name: listForm.name, list_type: listForm.list_type, description: listForm.description, entries: [] });
      setCreatingList(false);
      setListForm({ name: "", list_type: "allow", description: "" });
      await loadIPs();
    } catch {
      setIpError("Failed to create IP list");
    } finally {
      setSubmitting(false);
    }
  };

  const submitRL = async () => {
    setSubmitting(true);
    try {
      await createRateLimit({
        name: rlForm.name,
        route_prefix: rlForm.route_prefix,
        limit_count: Number(rlForm.limit_count),
        window_seconds: Number(rlForm.window_seconds),
        action: rlForm.action,
      });
      setCreatingRL(false);
      setRlForm({ name: "", route_prefix: "/", limit_count: "100", window_seconds: "60", action: "throttle" });
      await loadRates();
    } catch {
      setRlError("Failed to create rate limit");
    } finally {
      setSubmitting(false);
    }
  };

  const confirmDeleteList = async () => {
    if (!deletingList) return;
    setSubmitting(true);
    try {
      await deleteIPList(deletingList.id);
      setDeletingList(null);
      await loadIPs();
    } catch {
      setIpError("Failed to delete IP list");
    } finally {
      setSubmitting(false);
    }
  };

  const confirmDeleteRL = async () => {
    if (!deletingRL) return;
    setSubmitting(true);
    try {
      await deleteRateLimit(deletingRL.id);
      setDeletingRL(null);
      await loadRates();
    } catch {
      setRlError("Failed to delete rate limit");
    } finally {
      setSubmitting(false);
    }
  };

  const ipColumns: Column<IPList>[] = [
    { key: "name", label: "Name", sortable: true },
    { key: "list_type", label: "Type", render: (row) => <StatusBadge value={row.list_type} /> },
    { key: "entries", label: "Entries", render: (row) => String(row.entries?.length ?? 0) },
    { key: "description", label: "Description", render: (row) => row.description || "—" },
  ];

  const rlColumns: Column<RateLimit>[] = [
    { key: "name", label: "Name", sortable: true },
    { key: "route_prefix", label: "Path" },
    { key: "limit_count", label: "Limit", render: (row) => `${row.limit_count} req` },
    { key: "window_seconds", label: "Window", render: (row) => `${row.window_seconds}s` },
    { key: "action", label: "Action", render: (row) => <StatusBadge value={row.action} /> },
  ];

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />

      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <h1>Traffic Control</h1>
            <p style={{ color: "var(--text-secondary)" }}>IP lists and rate limiting.</p>
          </div>
          <div className="header-actions" style={{ display: "flex", alignItems: "center", gap: "12px" }}>
            <button
              type="button"
              className="btn btn-primary"
              onClick={() => (tab === "ip-lists" ? setCreatingList(true) : setCreatingRL(true))}
            >
              {tab === "ip-lists" ? "+ New IP List" : "+ New Rate Limit"}
            </button>
            <UserProfileWidget />
          </div>
        </div>

        {/* Tabs */}
        <div style={{ display: "flex", gap: "8px", marginBottom: "20px" }}>
          {(["ip-lists", "rate-limits"] as Tab[]).map((t) => (
            <button
              key={t}
              type="button"
              onClick={() => setTab(t)}
              style={{
                padding: "10px 18px",
                borderRadius: "8px",
                border: `1px solid ${tab === t ? "var(--accent-primary)" : "rgba(255,255,255,0.1)"}`,
                background: tab === t ? "rgba(59,130,246,0.15)" : "rgba(255,255,255,0.03)",
                color: tab === t ? "#60a5fa" : "var(--text-secondary)",
                fontSize: "14px",
                fontWeight: tab === t ? 600 : 500,
                cursor: "pointer",
              }}
            >
              {t === "ip-lists" ? "IP Lists" : "Rate Limits"}
            </button>
          ))}
        </div>

        {tab === "ip-lists" ? (
          <div className="data-section animate-fade-in delay-1">
            <DataTable
              columns={ipColumns}
              data={ipLists}
              rowKey={(row) => row.id}
              loading={ipLoading}
              error={ipError || undefined}
              onRetry={loadIPs}
              emptyMessage="No IP lists yet."
              actions={(row) => (
                <button type="button" className="btn" onClick={() => setDeletingList(row)} style={{ padding: "6px 12px", fontSize: "12px", color: "var(--danger)" }}>
                  Delete
                </button>
              )}
            />
          </div>
        ) : (
          <div className="data-section animate-fade-in delay-1">
            <DataTable
              columns={rlColumns}
              data={rateLimits}
              rowKey={(row) => row.id}
              loading={rlLoading}
              error={rlError || undefined}
              onRetry={loadRates}
              emptyMessage="No rate limits yet."
              actions={(row) => (
                <button type="button" className="btn" onClick={() => setDeletingRL(row)} style={{ padding: "6px 12px", fontSize: "12px", color: "var(--danger)" }}>
                  Delete
                </button>
              )}
            />
          </div>
        )}
      </main>

      {/* Create IP List modal */}
      {creatingList && (
        <ModalShell title="New IP List" onClose={() => setCreatingList(false)}>
          <Field label="Name">
            <input type="text" value={listForm.name} onChange={(e) => setListForm({ ...listForm, name: e.target.value })} style={inputStyle} />
          </Field>
          <Field label="Type">
            <select value={listForm.list_type} onChange={(e) => setListForm({ ...listForm, list_type: e.target.value })} style={inputStyle}>
              <option value="allow" style={optionStyle}>Allow</option>
              <option value="block" style={optionStyle}>Block</option>
            </select>
          </Field>
          <Field label="Description">
            <textarea value={listForm.description} onChange={(e) => setListForm({ ...listForm, description: e.target.value })} rows={2} style={inputStyle} />
          </Field>
          <ModalFooter onCancel={() => setCreatingList(false)} onSubmit={submitList} submitting={submitting} submitLabel="Create" />
        </ModalShell>
      )}

      {/* Create Rate Limit modal */}
      {creatingRL && (
        <ModalShell title="New Rate Limit" onClose={() => setCreatingRL(false)}>
          <Field label="Name">
            <input type="text" value={rlForm.name} onChange={(e) => setRlForm({ ...rlForm, name: e.target.value })} style={inputStyle} />
          </Field>
          <Field label="Path Prefix">
            <input type="text" value={rlForm.route_prefix} onChange={(e) => setRlForm({ ...rlForm, route_prefix: e.target.value })} style={inputStyle} />
          </Field>
          <Field label="Limit (requests)">
            <input type="number" value={rlForm.limit_count} onChange={(e) => setRlForm({ ...rlForm, limit_count: e.target.value })} style={inputStyle} />
          </Field>
          <Field label="Window (seconds)">
            <input type="number" value={rlForm.window_seconds} onChange={(e) => setRlForm({ ...rlForm, window_seconds: e.target.value })} style={inputStyle} />
          </Field>
          <Field label="Action">
            <select value={rlForm.action} onChange={(e) => setRlForm({ ...rlForm, action: e.target.value })} style={inputStyle}>
              <option value="throttle" style={optionStyle}>Throttle</option>
              <option value="block" style={optionStyle}>Block</option>
              <option value="log" style={optionStyle}>Log</option>
            </select>
          </Field>
          <ModalFooter onCancel={() => setCreatingRL(false)} onSubmit={submitRL} submitting={submitting} submitLabel="Create" />
        </ModalShell>
      )}

      <ConfirmDialog
        open={!!deletingList}
        title="Delete IP list"
        message={`Are you sure you want to delete "${deletingList?.name}"?`}
        confirmLabel="Delete"
        danger
        loading={submitting}
        onConfirm={confirmDeleteList}
        onCancel={() => setDeletingList(null)}
      />
      <ConfirmDialog
        open={!!deletingRL}
        title="Delete rate limit"
        message={`Are you sure you want to delete "${deletingRL?.name}"?`}
        confirmLabel="Delete"
        danger
        loading={submitting}
        onConfirm={confirmDeleteRL}
        onCancel={() => setDeletingRL(null)}
      />
    </div>
  );
}

const inputStyle: React.CSSProperties = {
  background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)",
  padding: "10px 12px", borderRadius: "8px", color: "white", outline: "none", fontSize: "14px", width: "100%",
};
const optionStyle: React.CSSProperties = { background: "#13141c", color: "#fff" };

function ModalShell({ title, onClose, children }: { title: string; onClose: () => void; children: React.ReactNode }) {
  return (
    <div onClick={onClose} style={{ position: "fixed", inset: 0, zIndex: 1000, background: "rgba(0,0,0,0.6)", display: "flex", alignItems: "center", justifyContent: "center", padding: "20px" }}>
      <div className="glass-panel animate-fade-in" onClick={(e) => e.stopPropagation()} style={{ width: "100%", maxWidth: "440px", padding: "28px", display: "flex", flexDirection: "column", gap: "14px" }}>
        <h3 style={{ fontSize: "17px", fontWeight: 600 }}>{title}</h3>
        {children}
      </div>
    </div>
  );
}

function Field({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
      <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>{label}</label>
      {children}
    </div>
  );
}

function ModalFooter({ onCancel, onSubmit, submitting, submitLabel }: { onCancel: () => void; onSubmit: () => void; submitting: boolean; submitLabel: string }) {
  return (
    <div style={{ display: "flex", justifyContent: "flex-end", gap: "10px", marginTop: "8px" }}>
      <button type="button" className="btn" onClick={onCancel} disabled={submitting} style={{ padding: "10px 16px" }}>Cancel</button>
      <button type="button" className="btn btn-primary" onClick={onSubmit} disabled={submitting} style={{ padding: "10px 16px" }}>
        {submitting ? "Saving…" : submitLabel}
      </button>
    </div>
  );
}
