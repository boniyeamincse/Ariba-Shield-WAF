"use client";

import { useCallback, useEffect, useState } from "react";
import { useLocale, useTranslations } from "next-intl";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import { getSettings, updateSettings } from "@/lib/api";

const TABS = [
  { key: "general", label: "General", icon: <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><circle cx="12" cy="12" r="3"></circle><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path></svg> },
  { key: "security", label: "Security", icon: <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg> },
  { key: "localization", label: "Localization", icon: <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><circle cx="12" cy="12" r="10"></circle><line x1="2" y1="12" x2="22" y2="12"></line><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path></svg> },
  { key: "retention", label: "Retention", icon: <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><polyline points="1 4 1 10 7 10"></polyline><path d="M3.51 15a9 9 0 1 0 2.13-9.36L1 10"></path></svg> },
] as const;

const SETTINGS_SCHEMA: Record<string, { key: string; label: string; type: "text" | "number" | "boolean" | "select"; options?: string[]; default: string | number | boolean; desc: string }[]> = {
  general: [
    { key: "site_name", label: "Site Name", type: "text", default: "Ariba Shield WAF", desc: "The name of the platform displayed across the dashboard." },
    { key: "admin_email", label: "Admin Email", type: "text", default: "admin@aribashield.local", desc: "System email used for core administrative alerts." },
    { key: "maintenance_mode", label: "Maintenance Mode", type: "boolean", default: false, desc: "If enabled, blocks external API changes." }
  ],
  security: [
    { key: "max_login_attempts", label: "Max Login Attempts", type: "number", default: 5, desc: "Number of failed attempts before temporary IP block." },
    { key: "session_timeout_minutes", label: "Session Timeout (Minutes)", type: "number", default: 60, desc: "Idle time before automatic logout." },
    { key: "mfa_required", label: "Require MFA for Admins", type: "boolean", default: false, desc: "Force all administrative users to configure Two-Factor Authentication." }
  ],
  localization: [
    { key: "default_language", label: "Default Language", type: "select", options: ["en", "bn"], default: "en", desc: "Primary language for newly created users." },
    { key: "timezone", label: "Timezone", type: "text", default: "UTC", desc: "Timezone used for displaying dates and times." }
  ],
  retention: [
    { key: "audit_log_days", label: "Audit Log Retention (Days)", type: "number", default: 90, desc: "How long to keep administrative audit events." },
    { key: "security_event_days", label: "Security Event Retention (Days)", type: "number", default: 30, desc: "How long to keep blocked traffic security events." }
  ]
};

type TabKey = (typeof TABS)[number]["key"];

export default function SettingsPage() {
  const locale = useLocale();
  const t = useTranslations("settings");
  const [tab, setTab] = useState<TabKey>("general");
  const [dbSettings, setDbSettings] = useState<Record<string, Record<string, unknown>>>({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);
  const [form, setForm] = useState<Record<string, string | number | boolean>>({});

  const load = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const data = await getSettings();
      setDbSettings(data || {});
    } catch {
      setError("Failed to load settings. The database might be unreachable.");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  useEffect(() => {
    const currentCategory = dbSettings[tab] || {};
    const schema = SETTINGS_SCHEMA[tab] || [];
    const initialForm: Record<string, string | number | boolean> = {};
    
    schema.forEach(field => {
      const v = currentCategory[field.key];
      if (v !== undefined && v !== null) {
        initialForm[field.key] = typeof v === "object" ? field.default : (v as string | number | boolean);
      } else {
        initialForm[field.key] = field.default;
      }
    });
    setForm(initialForm);
    setSaved(false);
  }, [tab, dbSettings]);

  const setValue = (key: string, value: string | number | boolean) => {
    setForm((prev) => ({ ...prev, [key]: value }));
    setSaved(false);
  };

  const save = async () => {
    setSaving(true);
    setError("");
    setSaved(false);
    try {
      const payload: Record<string, unknown> = {};
      const schema = SETTINGS_SCHEMA[tab] || [];
      for (const field of schema) {
        let val = form[field.key];
        if (field.type === 'number') {
          const numVal = Number(val);
          val = isNaN(numVal) ? field.default : numVal;
        }
        payload[field.key] = val;
      }
      await updateSettings(tab, payload);
      setSaved(true);
      setTimeout(() => setSaved(false), 3000);
      
      // Update local dbSettings cache immediately to prevent jitter
      setDbSettings(prev => ({
        ...prev,
        [tab]: { ...prev[tab], ...payload }
      }));
    } catch {
      setError("Failed to save settings. Please try again.");
    } finally {
      setSaving(false);
    }
  };

  const schema = SETTINGS_SCHEMA[tab] || [];

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />

      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <h1>{t("title") || "Settings"}</h1>
            <p style={{ color: "var(--text-secondary)" }}>Configure system-wide platform options and preferences.</p>
          </div>
          <div className="header-actions" style={{ display: "flex", alignItems: "center", gap: "16px" }}>
            <UserProfileWidget />
          </div>
        </div>

        {error && (
          <div
            className="glass-panel animate-fade-in"
            style={{ display: 'flex', alignItems: 'center', gap: '12px', padding: "16px 20px", marginBottom: "24px", color: "var(--danger)", border: "1px solid rgba(239,68,68,0.3)" }}
          >
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>
            {error}
          </div>
        )}

        {/* Tabs Layout */}
        <div style={{ display: 'grid', gridTemplateColumns: '240px 1fr', gap: '32px', alignItems: 'start' }} className="animate-fade-in delay-1">
          
          <div className="glass-panel" style={{ display: 'flex', flexDirection: 'column', gap: '8px', padding: '16px' }}>
            <h3 style={{ fontSize: '12px', textTransform: 'uppercase', letterSpacing: '0.1em', color: 'var(--text-secondary)', marginBottom: '8px', marginLeft: '8px' }}>Categories</h3>
            {TABS.map((t) => {
              const active = tab === t.key;
              return (
                <button
                  key={t.key}
                  type="button"
                  onClick={() => setTab(t.key)}
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: '12px',
                    padding: "12px 16px",
                    borderRadius: "12px",
                    border: 'none',
                    background: active ? "rgba(59,130,246,0.15)" : "transparent",
                    color: active ? "#60a5fa" : "var(--text-secondary)",
                    fontSize: "14px",
                    fontWeight: active ? 600 : 500,
                    cursor: "pointer",
                    transition: 'all 0.2s ease',
                    textAlign: 'left'
                  }}
                  onMouseOver={(e) => {
                    if (!active) e.currentTarget.style.background = 'rgba(255,255,255,0.05)';
                  }}
                  onMouseOut={(e) => {
                    if (!active) e.currentTarget.style.background = 'transparent';
                  }}
                >
                  <span style={{ color: active ? '#60a5fa' : 'var(--text-secondary)' }}>{t.icon}</span>
                  {t.label}
                </button>
              );
            })}
          </div>

          <div className="glass-panel" style={{ padding: "32px", position: 'relative' }}>
            
            {saved && (
              <div style={{ position: 'absolute', top: '32px', right: '32px', display: 'flex', alignItems: 'center', gap: '8px', color: 'var(--success)', background: 'rgba(16,185,129,0.1)', padding: '6px 12px', borderRadius: '999px', border: '1px solid rgba(16,185,129,0.2)' }} className="animate-fade-in">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><polyline points="20 6 9 17 4 12"></polyline></svg>
                <span style={{ fontSize: '13px', fontWeight: 600 }}>Saved successfully</span>
              </div>
            )}

            <h2 style={{ fontSize: '20px', marginBottom: '8px', display: 'flex', alignItems: 'center', gap: '12px' }}>
              {TABS.find(t => t.key === tab)?.label} Settings
            </h2>
            <p style={{ color: 'var(--text-secondary)', fontSize: '14px', marginBottom: '32px' }}>
              Update your {tab} configurations. Changes are applied globally.
            </p>

            {loading ? (
              <div style={{ display: 'flex', justifyContent: 'center', padding: '40px 0' }}>
                 <div style={{ width: '32px', height: '32px', borderRadius: '50%', border: '3px solid rgba(59, 130, 246, 0.1)', borderTopColor: '#3b82f6', animation: 'spin 1s linear infinite' }} />
              </div>
            ) : (
              <div style={{ display: "flex", flexDirection: "column", gap: "24px" }}>
                {schema.map((field) => (
                  <div key={field.key} style={{ display: "flex", flexDirection: "column", gap: "8px" }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                      <label style={{ fontSize: "14px", fontWeight: 500, color: "var(--text-primary)" }}>
                        {field.label}
                      </label>
                    </div>
                    
                    {field.type === 'boolean' ? (
                      <div style={{ display: 'flex', alignItems: 'center', gap: '12px', marginTop: '4px' }}>
                         <button 
                            onClick={() => setValue(field.key, !form[field.key])}
                            style={{
                              width: '44px', height: '24px', borderRadius: '12px', 
                              background: form[field.key] ? 'var(--accent-primary)' : 'rgba(255,255,255,0.1)',
                              border: 'none', position: 'relative', cursor: 'pointer', transition: 'all 0.3s ease'
                            }}
                         >
                           <div style={{
                             width: '18px', height: '18px', borderRadius: '50%', background: 'white',
                             position: 'absolute', top: '3px', left: form[field.key] ? '23px' : '3px',
                             transition: 'all 0.3s ease'
                           }} />
                         </button>
                         <span style={{ fontSize: '13px', color: 'var(--text-secondary)' }}>{form[field.key] ? 'Enabled' : 'Disabled'}</span>
                      </div>
                    ) : field.type === 'select' ? (
                      <select
                        value={(form[field.key] as string | number) || ''}
                        onChange={(e) => setValue(field.key, e.target.value)}
                        style={{
                          background: "rgba(255,255,255,0.03)", border: "1px solid rgba(255,255,255,0.1)",
                          padding: "12px 16px", borderRadius: "10px", color: "white", outline: "none", fontSize: "14px",
                          width: '100%', maxWidth: '400px', appearance: 'none', cursor: 'pointer'
                        }}
                      >
                        {field.options?.map(opt => (
                          <option key={opt} value={opt} style={{ background: 'var(--bg-secondary)', color: 'white' }}>
                            {opt === 'en' ? 'English' : opt === 'bn' ? 'Bengali' : opt}
                          </option>
                        ))}
                      </select>
                    ) : (
                      <input
                        type={field.type === 'number' ? 'number' : 'text'}
                        value={form[field.key] !== undefined ? (form[field.key] as string | number) : ''}
                        onChange={(e) => setValue(field.key, e.target.value)}
                        style={{
                          background: "rgba(255,255,255,0.03)", border: "1px solid rgba(255,255,255,0.1)",
                          padding: "12px 16px", borderRadius: "10px", color: "white", outline: "none", fontSize: "14px",
                          width: '100%', maxWidth: '400px', transition: 'all 0.2s ease'
                        }}
                        onFocus={(e) => e.target.style.borderColor = 'var(--accent-primary)'}
                        onBlur={(e) => e.target.style.borderColor = 'rgba(255,255,255,0.1)'}
                      />
                    )}
                    <div style={{ fontSize: '12px', color: 'var(--text-secondary)' }}>{field.desc}</div>
                    
                    {/* Separator for all but last item */}
                    {schema[schema.length - 1].key !== field.key && (
                      <div style={{ height: '1px', background: 'rgba(255,255,255,0.05)', marginTop: '16px' }} />
                    )}
                  </div>
                ))}

                <div style={{ display: "flex", justifyContent: "flex-start", marginTop: "16px" }}>
                  <button type="button" className="btn btn-primary" onClick={save} disabled={saving} style={{ padding: "12px 24px" }}>
                    {saving ? (
                      <>
                        <div style={{ width: '16px', height: '16px', borderRadius: '50%', border: '2px solid rgba(255,255,255,0.3)', borderTopColor: 'white', animation: 'spin 1s linear infinite' }} />
                        Saving...
                      </>
                    ) : "Save Changes"}
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
