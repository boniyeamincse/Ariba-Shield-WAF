"use client";

import { useCallback, useEffect, useState } from "react";
import { useLocale } from "next-intl";
import Link from "next/link";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import { listManagedRules, updateManagedRule, updateManagedRulesGlobal, type ManagedRule } from "@/lib/api";

export default function ManagedRulesPage() {
  const locale = useLocale();
  const [rules, setRules] = useState<ManagedRule[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [updating, setUpdating] = useState<string | null>(null);
  const [globalParanoia, setGlobalParanoia] = useState(1);
  const [globalThreshold, setGlobalThreshold] = useState(5);
  const [savingGlobal, setSavingGlobal] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const data = await listManagedRules();
      setRules(data);
      if (data.length > 0) {
        setGlobalParanoia(data[0].paranoia_level ?? 1);
        setGlobalThreshold(data[0].anomaly_threshold ?? 5);
      }
    } catch {
      setError("Failed to load managed rules.");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const allEnabled = rules.length > 0 && rules.every((r) => r.enabled);

  const handleGlobal = async () => {
    setSavingGlobal(true);
    setError("");
    try {
      await updateManagedRulesGlobal({
        paranoia_level: globalParanoia,
        anomaly_threshold: globalThreshold,
      });
      await load();
    } catch {
      setError("Failed to update global CRS settings.");
    } finally {
      setSavingGlobal(false);
    }
  };

  const toggleAll = async () => {
    setSavingGlobal(true);
    setError("");
    try {
      await updateManagedRulesGlobal({ enabled: !allEnabled });
      await load();
    } catch {
      setError("Failed to toggle OWASP CRS.");
    } finally {
      setSavingGlobal(false);
    }
  };

  const handleToggle = async (rule: ManagedRule) => {
    setUpdating(rule.id);
    try {
      await updateManagedRule(rule.id, { enabled: !rule.enabled });
      await load();
    } catch {
      setError("Failed to update rule.");
    } finally {
      setUpdating(null);
    }
  };

  const handleSensitivity = async (rule: ManagedRule, sensitivity: string) => {
    setUpdating(rule.id);
    try {
      await updateManagedRule(rule.id, { sensitivity });
      await load();
    } catch {
      setError("Failed to update sensitivity.");
    } finally {
      setUpdating(null);
    }
  };

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />
      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <div style={{ display: "flex", alignItems: "center", gap: "10px", marginBottom: "8px" }}>
              <Link href={`/${locale}/rules`} style={{ color: "var(--text-secondary)", textDecoration: "none", fontSize: "14px" }}>Custom Rules</Link>
              <span style={{ color: "var(--text-secondary)" }}>/</span>
              <span style={{ color: "var(--text-primary)", fontWeight: 500, fontSize: "14px" }}>Managed Rules (OWASP CRS)</span>
            </div>
            <h1>Managed Rules (OWASP CRS)</h1>
            <p style={{ color: "var(--text-secondary)" }}>Enable and configure industry-standard protection rule groups.</p>
          </div>
          <div className="header-actions">
            <UserProfileWidget />
          </div>
        </div>

        {error && (
          <div style={{ padding: "12px 16px", background: "rgba(239,68,68,0.1)", border: "1px solid rgba(239,68,68,0.2)", color: "var(--danger)", borderRadius: "8px", margin: "20px 0" }}>
            {error}
          </div>
        )}

        {/* OWASP CRS global settings */}
        <div className="glass-panel animate-fade-in delay-1" style={{ padding: "24px", marginTop: "20px", display: "flex", flexDirection: "column", gap: "16px" }}>
          <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", flexWrap: "wrap", gap: "12px" }}>
            <div>
              <h3 style={{ fontSize: "16px", fontWeight: 600, color: "var(--text-primary)" }}>OWASP Core Rule Set (CRS)</h3>
              <p style={{ fontSize: "13px", color: "var(--text-secondary)", marginTop: "4px" }}>
                Master toggle + global sensitivity and anomaly threshold applied to all CRS groups.
              </p>
            </div>
            <button
              onClick={toggleAll}
              disabled={savingGlobal || rules.length === 0}
              className="btn"
              style={{
                padding: "10px 20px",
                background: allEnabled ? "var(--success)" : "rgba(255,255,255,0.1)",
                color: "white", border: "none", fontSize: "14px", fontWeight: 600,
              }}
            >
              {savingGlobal ? "…" : allEnabled ? "CRS Enabled" : "Enable CRS"}
            </button>
          </div>

          <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "16px", flexWrap: "wrap" }}>
            <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
              <label style={{ fontSize: "13px", color: "var(--text-secondary)" }}>Global Paranoia Level (Sensitivity)</label>
              <select
                value={globalParanoia}
                onChange={(e) => setGlobalParanoia(Number(e.target.value))}
                style={{
                  background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)",
                  color: "white", padding: "10px 12px", borderRadius: "8px", fontSize: "14px", outline: "none",
                }}
              >
                <option value={1}>Paranoia 1 — Low (fewer false positives)</option>
                <option value={2}>Paranoia 2 — Medium</option>
                <option value={3}>Paranoia 3 — High</option>
                <option value={4}>Paranoia 4 — Strict (more false positives)</option>
              </select>
            </div>
            <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
              <label style={{ fontSize: "13px", color: "var(--text-secondary)" }}>Anomaly Threshold (score-based blocking)</label>
              <input
                type="number"
                value={globalThreshold}
                onChange={(e) => setGlobalThreshold(Number(e.target.value))}
                min={1}
                style={{
                  background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)",
                  color: "white", padding: "10px 12px", borderRadius: "8px", fontSize: "14px", outline: "none",
                }}
              />
            </div>
          </div>

          <div style={{ display: "flex", justifyContent: "flex-end" }}>
            <button className="btn btn-primary" onClick={handleGlobal} disabled={savingGlobal} style={{ padding: "10px 20px" }}>
              {savingGlobal ? "Saving…" : "Apply Global Settings"}
            </button>
          </div>
        </div>

        <div className="data-section animate-fade-in delay-1" style={{ display: "flex", flexDirection: "column", gap: "16px", marginTop: "20px" }}>
          {loading && rules.length === 0 ? (
            <p style={{ color: "var(--text-secondary)" }}>Loading rules...</p>
          ) : rules.length === 0 ? (
             <div className="glass-panel" style={{ padding: "30px", textAlign: "center", color: "var(--text-secondary)" }}>
               No managed rules available.
             </div>
          ) : (
            rules.map((rule) => (
              <div key={rule.id} className="glass-panel" style={{ padding: "20px", display: "flex", justifyContent: "space-between", alignItems: "center", gap: "20px", opacity: updating === rule.id ? 0.6 : 1 }}>
                <div>
                  <h4 style={{ fontSize: "16px", fontWeight: 600, color: "var(--text-primary)", marginBottom: "4px" }}>{rule.name}</h4>
                  <p style={{ fontSize: "13px", color: "var(--text-secondary)", marginBottom: "8px" }}>Category: {rule.category}</p>
                  
                  <div style={{ display: "flex", alignItems: "center", gap: "12px" }}>
                    <span style={{ fontSize: "13px", color: "var(--text-secondary)" }}>Sensitivity (Paranoia Level):</span>
                    <select
                      value={rule.sensitivity || "low"}
                      onChange={(e) => handleSensitivity(rule, e.target.value)}
                      disabled={!rule.enabled || updating === rule.id}
                      style={{
                        background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)",
                        color: "white", padding: "4px 8px", borderRadius: "6px", fontSize: "13px", outline: "none",
                        opacity: !rule.enabled ? 0.5 : 1
                      }}
                    >
                      <option value="low">Low (Fewer false positives)</option>
                      <option value="medium">Medium</option>
                      <option value="high">High</option>
                      <option value="strict">Strict (More false positives)</option>
                    </select>
                  </div>
                </div>

                <div>
                  <button
                    onClick={() => handleToggle(rule)}
                    disabled={updating === rule.id}
                    className="btn"
                    style={{
                      padding: "8px 16px",
                      background: rule.enabled ? "var(--success)" : "rgba(255,255,255,0.1)",
                      color: "white",
                      border: "none"
                    }}
                  >
                    {rule.enabled ? "Enabled" : "Enable"}
                  </button>
                </div>
              </div>
            ))
          )}
        </div>
      </main>
    </div>
  );
}
