"use client";

import { useCallback, useEffect, useState } from "react";
import { useLocale } from "next-intl";
import { useParams, useRouter, useSearchParams } from "next/navigation";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import { SeverityBadge, StatusBadge } from "@/components/shared/Badges";
import ConfirmDialog from "@/components/shared/ConfirmDialog";
import {
  getRule,
  deleteRule,
  duplicateRule,
  testRule,
  enableRule,
  disableRule,
  type RuleFull,
} from "@/lib/api";

type Tab = "overview" | "conditions" | "applications" | "events" | "history";

export default function RuleDetailPage() {
  const locale = useLocale();
  const router = useRouter();
  const params = useParams<{ id: string }>();
  const search = useSearchParams();
  const [tab, setTab] = useState<Tab>("overview");
  const [rule, setRule] = useState<RuleFull | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [deleting, setDeleting] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  // Test modal
  const [testOpen, setTestOpen] = useState(search.get("test") === "1");
  const [testForm, setTestForm] = useState({ method: "GET", url: "/", body: "", headers: "" });
  const [testResult, setTestResult] = useState<{ matched: boolean; action: string; matched_fields: string[] } | null>(null);
  const [testErr, setTestErr] = useState("");

  const load = useCallback(async () => {
    if (!params?.id) return;
    setLoading(true);
    setError("");
    try {
      setRule(await getRule(params.id));
    } catch {
      setError("Failed to load rule");
    } finally {
      setLoading(false);
    }
  }, [params?.id]);

  useEffect(() => {
    load();
  }, [load]);

  const doToggle = async (enable: boolean) => {
    if (!rule) return;
    try {
      if (enable) await enableRule(rule.id);
      else await disableRule(rule.id);
      await load();
    } catch {
      setError("Failed to update status");
    }
  };

  const doDuplicate = async () => {
    if (!rule) return;
    try {
      await duplicateRule(rule.id);
      await load();
    } catch {
      setError("Failed to duplicate rule");
    }
  };

  const confirmDelete = async () => {
    if (!rule) return;
    setSubmitting(true);
    try {
      await deleteRule(rule.id);
      router.push(`/${locale}/rules`);
    } catch {
      setError("Failed to delete rule");
    } finally {
      setSubmitting(false);
    }
  };

  const runTest = async () => {
    setTestErr(""); setTestResult(null);
    const headers: Record<string, string> = {};
    testForm.headers.split("\n").forEach((line) => {
      const i = line.indexOf(":");
      if (i > 0) headers[line.slice(0, i).trim()] = line.slice(i + 1).trim();
    });
    try {
      setTestResult(await testRule(rule!.id, { method: testForm.method, url: testForm.url, headers, body: testForm.body }));
    } catch {
      setTestErr("Test failed — could not evaluate rule.");
    }
  };

  if (loading) {
    return (
      <div className="dashboard-container">
        <Sidebar locale={locale} />
        <main className="main-content"><p style={{ color: "var(--text-secondary)" }}>Loading…</p></main>
      </div>
    );
  }

  if (error || !rule) {
    return (
      <div className="dashboard-container">
        <Sidebar locale={locale} />
        <main className="main-content">
          <div className="glass-panel" style={{ padding: "40px", textAlign: "center" }}>
            <p style={{ color: "var(--danger)" }}>{error || "Rule not found"}</p>
            <button className="btn btn-primary" onClick={load} style={{ marginTop: "12px" }}>Retry</button>
          </div>
        </main>
      </div>
    );
  }

  const tabs: { key: Tab; label: string }[] = [
    { key: "overview", label: "Overview" },
    { key: "conditions", label: "Conditions" },
    { key: "applications", label: "Applications" },
    { key: "events", label: "Events" },
    { key: "history", label: "History" },
  ];

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />

      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <button type="button" onClick={() => router.push(`/${locale}/rules`)} style={{ background: "none", border: "none", color: "var(--text-secondary)", cursor: "pointer", padding: 0, marginBottom: "4px", fontSize: "13px" }}>
              ← Back to rules
            </button>
            <h1>{rule.name}</h1>
            <p style={{ color: "var(--text-secondary)", fontFamily: "monospace", fontSize: "13px" }}>{rule.rule_id}</p>
          </div>
          <div className="header-actions" style={{ display: "flex", alignItems: "center", gap: "10px", flexWrap: "wrap" }}>
            <StatusBadge value={rule.status ?? "active"} />
            <button className="btn" onClick={() => setTestOpen(true)} style={{ padding: "8px 14px", fontSize: "13px" }}>Test Rule</button>
            <button className="btn" onClick={doDuplicate} style={{ padding: "8px 14px", fontSize: "13px" }}>Duplicate</button>
            <button className="btn" onClick={() => doToggle(rule.status !== "active")} style={{ padding: "8px 14px", fontSize: "13px", color: rule.status === "active" ? "var(--warning)" : "var(--success)" }}>
              {rule.status === "active" ? "Disable" : "Enable"}
            </button>
            <button className="btn" onClick={() => setDeleting(true)} style={{ padding: "8px 14px", fontSize: "13px", color: "var(--danger)" }}>Delete</button>
            <UserProfileWidget />
          </div>
        </div>

        {/* Tabs */}
        <div style={{ display: "flex", gap: "8px", marginBottom: "20px", flexWrap: "wrap" }}>
          {tabs.map((t) => (
            <button key={t.key} onClick={() => setTab(t.key)} style={{
              padding: "10px 16px", borderRadius: "8px",
              border: `1px solid ${tab === t.key ? "var(--accent-primary)" : "rgba(255,255,255,0.1)"}`,
              background: tab === t.key ? "rgba(59,130,246,0.15)" : "rgba(255,255,255,0.03)",
              color: tab === t.key ? "#60a5fa" : "var(--text-secondary)",
              fontSize: "13px", fontWeight: tab === t.key ? 600 : 500, cursor: "pointer",
            }}>{t.label}</button>
          ))}
        </div>

        {tab === "overview" && (
          <div className="metrics-grid animate-fade-in delay-1" style={{ maxWidth: "900px" }}>
            <div className="glass-panel metric-card" style={{ padding: "20px" }}>
              <div className="metric-header"><span>Category</span></div>
              <div style={{ marginTop: "8px" }}>{rule.category ? <StatusBadge value={rule.category} /> : "—"}</div>
            </div>
            <div className="glass-panel metric-card" style={{ padding: "20px" }}>
              <div className="metric-header"><span>Severity</span></div>
              <div style={{ marginTop: "8px" }}><SeverityBadge value={rule.severity} /></div>
            </div>
            <div className="glass-panel metric-card" style={{ padding: "20px" }}>
              <div className="metric-header"><span>Action</span></div>
              <div style={{ marginTop: "8px" }}><StatusBadge value={rule.action} /></div>
            </div>
            <div className="glass-panel metric-card" style={{ padding: "20px" }}>
              <div className="metric-header"><span>Priority</span></div>
              <div className="metric-value" style={{ fontSize: "24px" }}>{rule.priority ?? "—"}</div>
            </div>
          </div>
        )}

        {tab === "conditions" && (
          <div className="data-section animate-fade-in delay-1">
            <div className="glass-panel" style={{ padding: "24px" }}>
              <h4 style={{ fontSize: "15px", fontWeight: 600, marginBottom: "12px" }}>Rule Logic</h4>
              <p style={{ fontSize: "13px", color: "var(--text-secondary)", marginBottom: "12px" }}>IF</p>
              {(rule.conditions ?? []).map((c, i) => (
                <div key={i} style={{ display: "flex", gap: "10px", alignItems: "center", marginBottom: "8px", fontSize: "14px" }}>
                  <span style={{ fontFamily: "monospace", fontSize: "12px", color: "var(--accent-primary)", minWidth: "120px" }}>{c.field}</span>
                  <span style={{ color: "var(--text-secondary)" }}>{c.operator}</span>
                  <span style={{ fontFamily: "monospace" }}>&ldquo;{c.value}&rdquo;</span>
                  {i < (rule.conditions ?? []).length - 1 && <span style={{ color: "var(--warning)", marginLeft: "6px" }}>{rule.logic ?? "AND"}</span>}
                </div>
              ))}
              <p style={{ fontSize: "13px", color: "var(--text-secondary)", marginTop: "12px" }}>THEN <span style={{ color: "var(--danger)", fontWeight: 700 }}>{rule.action?.toUpperCase()}</span></p>
            </div>
          </div>
        )}

        {tab === "applications" && (
          <div className="data-section animate-fade-in delay-1">
            <div className="glass-panel" style={{ padding: "24px" }}>
              <h4 style={{ fontSize: "15px", fontWeight: 600, marginBottom: "12px" }}>Scope</h4>
              {(rule.scopes ?? []).length === 0 ? (
                <p style={{ color: "var(--text-secondary)" }}>All applications</p>
              ) : (
                (rule.scopes ?? []).map((s, i) => (
                  <div key={i} style={{ display: "flex", gap: "12px", marginBottom: "8px", fontSize: "14px" }}>
                    <span style={{ fontFamily: "monospace", color: "var(--accent-primary)" }}>{s.path_pattern || "/*"}</span>
                    <span style={{ color: "var(--text-secondary)" }}>{(s.methods ?? []).join(", ") || "All methods"}</span>
                  </div>
                ))
              )}
            </div>
          </div>
        )}

        {tab === "events" && (
          <div className="data-section animate-fade-in delay-1">
            <div className="glass-panel" style={{ padding: "24px" }}>
              <p style={{ color: "var(--text-secondary)" }}>Security events for this rule appear here once traffic flows through the gateway.</p>
            </div>
          </div>
        )}

        {tab === "history" && (
          <div className="data-section animate-fade-in delay-1">
            <div className="glass-panel" style={{ padding: "24px" }}>
              <h4 style={{ fontSize: "15px", fontWeight: 600, marginBottom: "12px" }}>Version History</h4>
              <p style={{ color: "var(--text-secondary)", fontSize: "14px" }}>v{rule.version ?? 1} — current</p>
            </div>
          </div>
        )}
      </main>

      {/* Test Rule modal */}
      {testOpen && (
        <div onClick={() => setTestOpen(false)} style={{ position: "fixed", inset: 0, zIndex: 1200, background: "rgba(0,0,0,0.7)", display: "flex", alignItems: "center", justifyContent: "center", padding: "20px" }}>
          <div className="glass-panel animate-fade-in" onClick={(e) => e.stopPropagation()} style={{ width: "100%", maxWidth: "520px", padding: "28px", display: "flex", flexDirection: "column", gap: "14px" }}>
            <h3 style={{ fontSize: "17px", fontWeight: 700 }}>Test Rule — {rule.rule_id}</h3>
            <div style={{ display: "grid", gridTemplateColumns: "1fr 2fr", gap: "10px" }}>
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Method</label>
              <select style={inputStyle} value={testForm.method} onChange={(e) => setTestForm({ ...testForm, method: e.target.value })}>
                {["GET", "POST", "PUT", "PATCH", "DELETE"].map((m) => <option key={m} value={m}>{m}</option>)}
              </select>
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>URL</label>
              <input style={inputStyle} value={testForm.url} onChange={(e) => setTestForm({ ...testForm, url: e.target.value })} />
            </div>
            <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Headers (one per line, Key: Value)</label>
            <textarea style={{ ...inputStyle, resize: "vertical" }} rows={2} value={testForm.headers} onChange={(e) => setTestForm({ ...testForm, headers: e.target.value })} placeholder="Content-Type: application/x-www-form-urlencoded" />
            <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Body</label>
            <textarea style={{ ...inputStyle, resize: "vertical" }} rows={3} value={testForm.body} onChange={(e) => setTestForm({ ...testForm, body: e.target.value })} placeholder="username=admin' OR 1=1--" />

            {testErr && <p style={{ color: "var(--danger)", fontSize: "13px" }}>{testErr}</p>}

            {testResult && (
              <div style={{ padding: "14px", borderRadius: "8px", background: testResult.matched ? "rgba(239,68,68,0.1)" : "rgba(16,185,129,0.1)", border: `1px solid ${testResult.matched ? "rgba(239,68,68,0.3)" : "rgba(16,185,129,0.3)"}` }}>
                <div style={{ fontSize: "15px", fontWeight: 700, color: testResult.matched ? "var(--danger)" : "var(--success)" }}>
                  {testResult.matched ? "✓ RULE MATCHED" : "✓ RULE NOT MATCHED"}
                </div>
                <div style={{ fontSize: "13px", color: "var(--text-secondary)", marginTop: "6px" }}>
                  Action: <span style={{ textTransform: "uppercase", fontWeight: 600 }}>{testResult.action}</span>
                </div>
                {testResult.matched_fields.length > 0 && (
                  <div style={{ fontSize: "13px", color: "var(--text-secondary)", marginTop: "4px" }}>
                    Matched: {testResult.matched_fields.join(", ")}
                  </div>
                )}
              </div>
            )}

            <div style={{ display: "flex", justifyContent: "flex-end", gap: "10px" }}>
              <button className="btn" onClick={() => setTestOpen(false)} style={{ padding: "10px 16px" }}>Close</button>
              <button className="btn btn-primary" onClick={runTest} style={{ padding: "10px 16px" }}>Test Rule</button>
            </div>
          </div>
        </div>
      )}

      <ConfirmDialog
        open={deleting}
        title="Delete rule"
        message={`Are you sure you want to delete "${rule.name}" (${rule.rule_id})?`}
        confirmLabel="Delete"
        danger
        loading={submitting}
        onConfirm={confirmDelete}
        onCancel={() => setDeleting(false)}
      />
    </div>
  );
}

const inputStyle: React.CSSProperties = {
  background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)",
  padding: "10px 12px", borderRadius: "8px", color: "white", outline: "none", fontSize: "14px", width: "100%",
};
