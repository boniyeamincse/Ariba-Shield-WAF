"use client";

import { useCallback, useEffect, useState } from "react";
import { useLocale } from "next-intl";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import { getSettings, updateSettings } from "@/lib/api";

const TABS = [
  { key: "general", label: "General" },
  { key: "security", label: "Security" },
  { key: "localization", label: "Localization" },
  { key: "retention", label: "Retention" },
] as const;

type TabKey = (typeof TABS)[number]["key"];

export default function SettingsPage() {
  const locale = useLocale();
  const [tab, setTab] = useState<TabKey>("general");
  const [settings, setSettings] = useState<Record<string, Record<string, unknown>>>({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);
  const [form, setForm] = useState<Record<string, string>>({});

  const load = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const data = await getSettings();
      setSettings(data);
    } catch {
      setError("Failed to load settings");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  useEffect(() => {
    const current = settings[tab] ?? {};
    const str: Record<string, string> = {};
    for (const [k, v] of Object.entries(current)) {
      str[k] = String(v ?? "");
    }
    setForm(str);
    setSaved(false);
  }, [tab, settings]);

  const setValue = (key: string, value: string) => {
    setForm((prev) => ({ ...prev, [key]: value }));
    setSaved(false);
  };

  const save = async () => {
    setSaving(true);
    setError("");
    try {
      // Convert numeric strings back to numbers for the API.
      const payload: Record<string, unknown> = {};
      for (const [k, v] of Object.entries(form)) {
        const n = Number(v);
        payload[k] = v !== "" && !Number.isNaN(n) && String(n) === v ? n : v;
      }
      await updateSettings(tab, payload);
      setSaved(true);
      await load();
    } catch {
      setError("Failed to save settings");
    } finally {
      setSaving(false);
    }
  };

  const renderValue = (key: string): string => form[key] ?? "";

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />

      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <h1>System Settings</h1>
            <p style={{ color: "var(--text-secondary)" }}>Configure system-wide options.</p>
          </div>
          <div className="header-actions" style={{ display: "flex", alignItems: "center", gap: "12px" }}>
            {saved && <span style={{ color: "var(--success)", fontSize: "13px" }}>✓ Saved</span>}
            <UserProfileWidget />
          </div>
        </div>

        {error && (
          <div
            className="glass-panel"
            style={{ padding: "12px 16px", marginBottom: "16px", color: "var(--danger)", fontSize: "13px", border: "1px solid rgba(239,68,68,0.3)" }}
          >
            {error}
          </div>
        )}

        {/* Tabs */}
        <div style={{ display: "flex", gap: "8px", marginBottom: "20px", flexWrap: "wrap" }}>
          {TABS.map((t) => (
            <button
              key={t.key}
              type="button"
              onClick={() => setTab(t.key)}
              style={{
                padding: "10px 18px",
                borderRadius: "8px",
                border: `1px solid ${tab === t.key ? "var(--accent-primary)" : "rgba(255,255,255,0.1)"}`,
                background: tab === t.key ? "rgba(59,130,246,0.15)" : "rgba(255,255,255,0.03)",
                color: tab === t.key ? "#60a5fa" : "var(--text-secondary)",
                fontSize: "14px",
                fontWeight: tab === t.key ? 600 : 500,
                cursor: "pointer",
              }}
            >
              {t.label}
            </button>
          ))}
        </div>

        <div className="data-section animate-fade-in delay-1">
          <div className="glass-panel" style={{ padding: "24px" }}>
            {loading ? (
              <p style={{ color: "var(--text-secondary)" }}>Loading settings…</p>
            ) : Object.keys(form).length === 0 ? (
              <p style={{ color: "var(--text-secondary)" }}>No settings configured for this category yet.</p>
            ) : (
              <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
                {Object.keys(form).map((key) => (
                  <div key={key} style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
                    <label style={{ fontSize: "14px", color: "var(--text-secondary)", textTransform: "capitalize" }}>
                      {key.replace(/_/g, " ")}
                    </label>
                    <input
                      type="text"
                      value={renderValue(key)}
                      onChange={(e) => setValue(key, e.target.value)}
                      style={{
                        background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)",
                        padding: "10px 12px", borderRadius: "8px", color: "white", outline: "none", fontSize: "14px",
                      }}
                    />
                  </div>
                ))}
                <div style={{ display: "flex", justifyContent: "flex-end", marginTop: "8px" }}>
                  <button type="button" className="btn btn-primary" onClick={save} disabled={saving} style={{ padding: "10px 20px" }}>
                    {saving ? "Saving…" : "Save Changes"}
                  </button>
                </div>
              </div>
            )}
          </div>
        </div>
      </main>
    </div>
  );
}
