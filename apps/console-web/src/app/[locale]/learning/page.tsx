"use client";

import { useCallback, useEffect, useState } from "react";
import type { CSSProperties } from "react";
import { useLocale } from "next-intl";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import DataTable, { type Column } from "@/components/shared/DataTable";
import { SeverityBadge, StatusBadge } from "@/components/shared/Badges";
import ConfirmDialog from "@/components/shared/ConfirmDialog";
import {
  listLearningSessions,
  createLearningSession,
  startLearningSession,
  stopLearningSession,
  listLearningSuggestions,
  acceptSuggestion,
  rejectSuggestion,
  type LearningSession,
  type LearningSuggestion,
} from "@/lib/api";

type Tab = "sessions" | "suggestions";

export default function LearningPage() {
  const locale = useLocale();
  const [tab, setTab] = useState<Tab>("sessions");

  // Sessions
  const [sessions, setSessions] = useState<LearningSession[]>([]);
  const [sLoading, setSLoading] = useState(true);
  const [sError, setSError] = useState("");
  const [creating, setCreating] = useState(false);
  const [sForm, setSForm] = useState({ name: "", source: "trusted", description: "", confidence_threshold: "0.7" });
  const [submitting, setSubmitting] = useState(false);

  // Suggestions
  const [suggestions, setSuggestions] = useState<LearningSuggestion[]>([]);
  const [gLoading, setGLoading] = useState(false);
  const [gError, setGError] = useState("");
  const [confirmSuggestion, setConfirmSuggestion] = useState<{ sug: LearningSuggestion; accept: boolean } | null>(null);

  const loadSessions = useCallback(async () => {
    setSLoading(true);
    setSError("");
    try {
      setSessions(await listLearningSessions());
    } catch {
      setSError("Failed to load learning sessions");
    } finally {
      setSLoading(false);
    }
  }, []);

  const loadSuggestions = useCallback(async () => {
    setGLoading(true);
    setGError("");
    try {
      setSuggestions(await listLearningSuggestions());
    } catch {
      setGError("Failed to load suggestions");
    } finally {
      setGLoading(false);
    }
  }, []);

  useEffect(() => {
    loadSessions();
  }, [loadSessions]);

  useEffect(() => {
    if (tab === "suggestions") loadSuggestions();
  }, [tab, loadSuggestions]);

  const submitSession = async () => {
    setSubmitting(true);
    try {
      await createLearningSession({
        name: sForm.name,
        source: sForm.source,
        description: sForm.description,
        confidence_threshold: sForm.confidence_threshold,
      });
      setCreating(false);
      setSForm({ name: "", source: "trusted", description: "", confidence_threshold: "0.7" });
      await loadSessions();
    } catch {
      setSError("Failed to create learning session");
    } finally {
      setSubmitting(false);
    }
  };

  const toggleSession = async (row: LearningSession, start: boolean) => {
    try {
      if (start) await startLearningSession(row.id);
      else await stopLearningSession(row.id);
      await loadSessions();
    } catch {
      setSError(`Failed to ${start ? "start" : "stop"} session`);
    }
  };

  const confirmSuggestionAction = async () => {
    if (!confirmSuggestion) return;
    setSubmitting(true);
    try {
      if (confirmSuggestion.accept) await acceptSuggestion(confirmSuggestion.sug.id);
      else await rejectSuggestion(confirmSuggestion.sug.id);
      setConfirmSuggestion(null);
      await loadSuggestions();
    } catch {
      setGError("Failed to update suggestion");
    } finally {
      setSubmitting(false);
    }
  };

  const sessionColumns: Column<LearningSession>[] = [
    { key: "name", label: "Name", sortable: true },
    { key: "source", label: "Source", render: (row) => <StatusBadge value={row.source} /> },
    { key: "confidence_threshold", label: "Confidence", render: (row) => row.confidence_threshold },
    { key: "status", label: "Status", render: (row) => <StatusBadge value={row.status} /> },
    { key: "created_at", label: "Created", render: (row) => new Date(row.created_at).toLocaleDateString() },
  ];

  const suggestionColumns: Column<LearningSuggestion>[] = [
    { key: "rule_id", label: "Rule", render: (row) => <span style={{ fontFamily: "monospace", fontSize: "12px" }}>{row.rule_id}</span> },
    { key: "severity", label: "Severity", render: (row) => <SeverityBadge value={row.severity} /> },
    { key: "confidence", label: "Confidence", render: (row) => row.confidence },
    { key: "rationale", label: "Rationale", render: (row) => row.rationale || "—" },
    { key: "status", label: "Status", render: (row) => <StatusBadge value={row.status} /> },
  ];

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />

      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <h1>Learning</h1>
            <p style={{ color: "var(--text-secondary)" }}>Trusted-source learning and policy suggestions.</p>
          </div>
          <div className="header-actions" style={{ display: "flex", alignItems: "center", gap: "12px" }}>
            {tab === "sessions" && (
              <button type="button" className="btn btn-primary" onClick={() => setCreating(true)}>
                + New Session
              </button>
            )}
            <UserProfileWidget />
          </div>
        </div>

        {/* Tabs */}
        <div style={{ display: "flex", gap: "8px", marginBottom: "20px" }}>
          {(["sessions", "suggestions"] as Tab[]).map((t) => (
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
              {t === "sessions" ? "Sessions" : "Suggestions"}
            </button>
          ))}
        </div>

        {tab === "sessions" ? (
          <div className="data-section animate-fade-in delay-1">
            <DataTable
              columns={sessionColumns}
              data={sessions}
              rowKey={(row) => row.id}
              loading={sLoading}
              error={sError || undefined}
              onRetry={loadSessions}
              emptyMessage="No learning sessions yet."
              actions={(row) => (
                <button
                  type="button"
                  className="btn"
                  onClick={() => toggleSession(row, row.status !== "active")}
                  style={{ padding: "6px 10px", fontSize: "12px", color: row.status === "active" ? "var(--warning)" : "var(--success)" }}
                >
                  {row.status === "active" ? "Stop" : "Start"}
                </button>
              )}
            />
          </div>
        ) : (
          <div className="data-section animate-fade-in delay-1">
            <DataTable
              columns={suggestionColumns}
              data={suggestions}
              rowKey={(row) => row.id}
              loading={gLoading}
              error={gError || undefined}
              onRetry={loadSuggestions}
              emptyMessage="No suggestions yet."
              actions={(row) =>
                row.status === "pending" ? (
                  <div style={{ display: "flex", gap: "6px" }}>
                    <button type="button" className="btn" onClick={() => setConfirmSuggestion({ sug: row, accept: true })} style={{ padding: "6px 10px", fontSize: "12px", color: "var(--success)" }}>
                      Accept
                    </button>
                    <button type="button" className="btn" onClick={() => setConfirmSuggestion({ sug: row, accept: false })} style={{ padding: "6px 10px", fontSize: "12px", color: "var(--danger)" }}>
                      Reject
                    </button>
                  </div>
                ) : null
              }
            />
          </div>
        )}
      </main>

      {/* Create session modal */}
      {creating && (
        <div
          onClick={() => setCreating(false)}
          style={{ position: "fixed", inset: 0, zIndex: 1000, background: "rgba(0,0,0,0.6)", display: "flex", alignItems: "center", justifyContent: "center", padding: "20px" }}
        >
          <div className="glass-panel animate-fade-in" onClick={(e) => e.stopPropagation()} style={{ width: "100%", maxWidth: "440px", padding: "28px", display: "flex", flexDirection: "column", gap: "16px" }}>
            <h3 style={{ fontSize: "17px", fontWeight: 600 }}>New Learning Session</h3>
            <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Name</label>
              <input type="text" value={sForm.name} onChange={(e) => setSForm({ ...sForm, name: e.target.value })} style={inputStyle} />
            </div>
            <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Source (must be trusted)</label>
              <select value={sForm.source} onChange={(e) => setSForm({ ...sForm, source: e.target.value })} style={inputStyle}>
                <option value="trusted" style={optionStyle}>trusted</option>
              </select>
            </div>
            <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Confidence Threshold</label>
              <input type="text" value={sForm.confidence_threshold} onChange={(e) => setSForm({ ...sForm, confidence_threshold: e.target.value })} style={inputStyle} />
            </div>
            <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Description</label>
              <textarea value={sForm.description} onChange={(e) => setSForm({ ...sForm, description: e.target.value })} rows={3} style={{ ...inputStyle, resize: "vertical" }} />
            </div>
            <div style={{ display: "flex", justifyContent: "flex-end", gap: "10px", marginTop: "8px" }}>
              <button type="button" className="btn" onClick={() => setCreating(false)} disabled={submitting} style={{ padding: "10px 16px" }}>Cancel</button>
              <button type="button" className="btn btn-primary" onClick={submitSession} disabled={submitting || !sForm.name} style={{ padding: "10px 16px" }}>
                {submitting ? "Creating…" : "Create Session"}
              </button>
            </div>
          </div>
        </div>
      )}

      <ConfirmDialog
        open={!!confirmSuggestion}
        title={confirmSuggestion?.accept ? "Accept suggestion" : "Reject suggestion"}
        message={`Are you sure you want to ${confirmSuggestion?.accept ? "accept" : "reject"} rule "${confirmSuggestion?.sug.rule_id}"?`}
        confirmLabel={confirmSuggestion?.accept ? "Accept" : "Reject"}
        danger={!confirmSuggestion?.accept}
        loading={submitting}
        onConfirm={confirmSuggestionAction}
        onCancel={() => setConfirmSuggestion(null)}
      />
    </div>
  );
}

const inputStyle: CSSProperties = {
  background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)",
  padding: "10px 12px", borderRadius: "8px", color: "white", outline: "none", fontSize: "14px",
};
const optionStyle: CSSProperties = { background: "#13141c", color: "#fff" };
