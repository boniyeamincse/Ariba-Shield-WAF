"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useLocale } from "next-intl";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import { createSecurityPolicy } from "@/lib/api";

const STEPS = [
  "Basic Information",
  "Security Rules",
  "Traffic Controls",
  "Policy Settings",
  "Review & Deploy"
];

export default function CreatePolicyWizard() {
  const router = useRouter();
  const locale = useLocale();
  const [step, setStep] = useState(0);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");

  const [form, setForm] = useState({
    name: "",
    description: "",
    environment: "production",
    rules: {
      crs: true,
      sqli: true,
      xss: true,
      rce: true,
      lfi: true,
      ssrf: true,
      custom: false
    },
    traffic: {
      ipControl: false,
      geoControl: false,
      rateLimit: false,
      botProtection: true
    },
    settings: {
      mode: "blocking", // transparent, alarm, blocking
      logging: "all",
      severityThreshold: "low"
    }
  });

  const nextStep = () => {
    if (step === 0 && !form.name) return; // Basic validation
    setStep((s) => Math.min(s + 1, STEPS.length - 1));
  };

  const prevStep = () => setStep((s) => Math.max(s - 1, 0));

  const submit = async (activate: boolean) => {
    setSubmitting(true);
    setError("");
    try {
      // In a real app, we would pass all the rules/traffic data. 
      // For now, the API expects name and enforcement_mode.
      const mode = activate ? form.settings.mode : "transparent"; // Draft mode is transparent
      await createSecurityPolicy(form.name, mode);
      router.push(`/${locale}/policies`);
    } catch (err) {
      setError("Failed to create security policy.");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />
      
      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <h1>Create Security Policy</h1>
            <p style={{ color: "var(--text-secondary)" }}>Design a robust defense profile using modular security rules.</p>
          </div>
          <div className="header-actions">
            <UserProfileWidget />
          </div>
        </div>

        {error && (
          <div style={{ background: "rgba(239,68,68,0.1)", border: "1px solid rgba(239,68,68,0.2)", padding: "12px 16px", borderRadius: "8px", color: "var(--danger)", marginBottom: "20px" }}>
            {error}
          </div>
        )}

        <div className="glass-panel animate-fade-in delay-1" style={{ padding: "0" }}>
          
          {/* Progress Tracker */}
          <div style={{ display: "flex", borderBottom: "1px solid rgba(255,255,255,0.05)", padding: "24px 32px" }}>
            {STEPS.map((s, idx) => (
              <div key={s} style={{ flex: 1, display: "flex", flexDirection: "column", gap: "8px", position: "relative" }}>
                <div style={{ display: "flex", alignItems: "center" }}>
                  <div style={{ 
                    width: "24px", height: "24px", borderRadius: "50%", 
                    background: step >= idx ? "var(--accent-primary)" : "rgba(255,255,255,0.1)",
                    color: step >= idx ? "white" : "rgba(255,255,255,0.5)",
                    display: "flex", alignItems: "center", justifyContent: "center", fontSize: "12px", fontWeight: 600, zIndex: 2
                  }}>
                    {idx + 1}
                  </div>
                  {idx < STEPS.length - 1 && (
                    <div style={{ 
                      flex: 1, height: "2px", 
                      background: step > idx ? "var(--accent-primary)" : "rgba(255,255,255,0.1)",
                      margin: "0 8px"
                    }} />
                  )}
                </div>
                <span style={{ fontSize: "13px", fontWeight: step === idx ? 600 : 400, color: step >= idx ? "white" : "var(--text-secondary)" }}>
                  {s}
                </span>
              </div>
            ))}
          </div>

          <div style={{ padding: "32px", minHeight: "400px" }}>
            
            {/* STEP 1: Basic Information */}
            {step === 0 && (
              <div className="animate-fade-in" style={{ display: "flex", flexDirection: "column", gap: "24px", maxWidth: "600px" }}>
                <h2 style={{ fontSize: "18px", fontWeight: 600 }}>1. Basic Information</h2>
                <div style={{ display: "flex", flexDirection: "column", gap: "8px" }}>
                  <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Policy Name *</label>
                  <input
                    type="text"
                    value={form.name}
                    onChange={(e) => setForm({ ...form, name: e.target.value })}
                    placeholder="e.g. Core Banking API Protection"
                    style={{ background: "rgba(255,255,255,0.03)", border: "1px solid rgba(255,255,255,0.1)", padding: "12px 16px", borderRadius: "8px", color: "white", outline: "none" }}
                  />
                </div>
                <div style={{ display: "flex", flexDirection: "column", gap: "8px" }}>
                  <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Description</label>
                  <textarea
                    value={form.description}
                    onChange={(e) => setForm({ ...form, description: e.target.value })}
                    rows={3}
                    placeholder="Describe what this policy protects..."
                    style={{ background: "rgba(255,255,255,0.03)", border: "1px solid rgba(255,255,255,0.1)", padding: "12px 16px", borderRadius: "8px", color: "white", outline: "none", resize: "vertical" }}
                  />
                </div>
                <div style={{ display: "flex", flexDirection: "column", gap: "8px" }}>
                  <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Target Environment</label>
                  <select
                    value={form.environment}
                    onChange={(e) => setForm({ ...form, environment: e.target.value })}
                    style={{ background: "rgba(255,255,255,0.03)", border: "1px solid rgba(255,255,255,0.1)", padding: "12px 16px", borderRadius: "8px", color: "white", outline: "none" }}
                  >
                    <option value="production" style={{ background: "#13141c" }}>Production</option>
                    <option value="staging" style={{ background: "#13141c" }}>Staging</option>
                    <option value="development" style={{ background: "#13141c" }}>Development</option>
                  </select>
                </div>
              </div>
            )}

            {/* STEP 2: Security Rules */}
            {step === 1 && (
              <div className="animate-fade-in" style={{ display: "flex", flexDirection: "column", gap: "24px" }}>
                <h2 style={{ fontSize: "18px", fontWeight: 600 }}>2. Security Rules Engine</h2>
                <p style={{ color: "var(--text-secondary)", fontSize: "14px" }}>Select the attack vectors this policy should inspect and mitigate.</p>
                
                <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "16px" }}>
                  {[
                    { id: "crs", label: "OWASP Core Rule Set (CRS)", desc: "Baseline protection against top web threats." },
                    { id: "sqli", label: "SQL Injection (SQLi)", desc: "Block database manipulation attempts." },
                    { id: "xss", label: "Cross-Site Scripting (XSS)", desc: "Prevent malicious script injection." },
                    { id: "rce", label: "Remote Code Execution (RCE)", desc: "Block OS command injection." },
                    { id: "lfi", label: "Local File Inclusion (LFI)", desc: "Prevent unauthorized file access." },
                    { id: "ssrf", label: "Server-Side Request Forgery", desc: "Block backend infrastructure probes." }
                  ].map((rule) => (
                    <div 
                      key={rule.id} 
                      onClick={() => setForm({ ...form, rules: { ...form.rules, [rule.id]: !form.rules[rule.id as keyof typeof form.rules] } })}
                      style={{ 
                        border: `1px solid ${form.rules[rule.id as keyof typeof form.rules] ? 'var(--accent-primary)' : 'rgba(255,255,255,0.1)'}`, 
                        background: form.rules[rule.id as keyof typeof form.rules] ? 'rgba(59,130,246,0.05)' : 'transparent',
                        padding: "16px", borderRadius: "12px", cursor: "pointer", display: "flex", gap: "16px", alignItems: "flex-start", transition: "all 0.2s" 
                      }}
                    >
                      <div style={{ 
                        width: "20px", height: "20px", borderRadius: "4px", border: "1px solid", 
                        borderColor: form.rules[rule.id as keyof typeof form.rules] ? "var(--accent-primary)" : "rgba(255,255,255,0.3)",
                        background: form.rules[rule.id as keyof typeof form.rules] ? "var(--accent-primary)" : "transparent",
                        display: "flex", alignItems: "center", justifyContent: "center", marginTop: "2px"
                      }}>
                        {form.rules[rule.id as keyof typeof form.rules] && <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="white" strokeWidth="3"><polyline points="20 6 9 17 4 12"></polyline></svg>}
                      </div>
                      <div>
                        <div style={{ fontSize: "14px", fontWeight: 600, color: form.rules[rule.id as keyof typeof form.rules] ? "white" : "var(--text-secondary)" }}>{rule.label}</div>
                        <div style={{ fontSize: "12px", color: "var(--text-secondary)", marginTop: "4px" }}>{rule.desc}</div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* STEP 3: Traffic Controls */}
            {step === 2 && (
              <div className="animate-fade-in" style={{ display: "flex", flexDirection: "column", gap: "24px" }}>
                <h2 style={{ fontSize: "18px", fontWeight: 600 }}>3. Traffic Controls</h2>
                <p style={{ color: "var(--text-secondary)", fontSize: "14px" }}>Enforce access rules at the network and application layer.</p>
                
                <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "16px" }}>
                  {[
                    { id: "ipControl", label: "IP Allow/Block Lists", desc: "Restrict access based on static IP lists." },
                    { id: "geoControl", label: "Geo-blocking", desc: "Block traffic from sanctioned countries." },
                    { id: "rateLimit", label: "Rate Limiting", desc: "Prevent brute-force and volumetric DDoS." },
                    { id: "botProtection", label: "Bot Protection", desc: "Challenge suspicious automated traffic." }
                  ].map((tc) => (
                    <div 
                      key={tc.id} 
                      style={{ 
                        border: "1px solid rgba(255,255,255,0.1)", background: "rgba(255,255,255,0.02)",
                        padding: "16px", borderRadius: "12px", display: "flex", justifyContent: "space-between", alignItems: "center" 
                      }}
                    >
                      <div>
                        <div style={{ fontSize: "14px", fontWeight: 500 }}>{tc.label}</div>
                        <div style={{ fontSize: "12px", color: "var(--text-secondary)", marginTop: "2px" }}>{tc.desc}</div>
                      </div>
                      
                      <button 
                         onClick={() => setForm({ ...form, traffic: { ...form.traffic, [tc.id]: !form.traffic[tc.id as keyof typeof form.traffic] } })}
                         style={{
                           width: '44px', height: '24px', borderRadius: '12px', 
                           background: form.traffic[tc.id as keyof typeof form.traffic] ? 'var(--accent-primary)' : 'rgba(255,255,255,0.1)',
                           border: 'none', position: 'relative', cursor: 'pointer', transition: 'all 0.3s ease'
                         }}
                      >
                        <div style={{
                          width: '18px', height: '18px', borderRadius: '50%', background: 'white',
                          position: 'absolute', top: '3px', left: form.traffic[tc.id as keyof typeof form.traffic] ? '23px' : '3px',
                          transition: 'all 0.3s ease'
                        }} />
                      </button>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* STEP 4: Policy Settings */}
            {step === 3 && (
              <div className="animate-fade-in" style={{ display: "flex", flexDirection: "column", gap: "24px", maxWidth: "600px" }}>
                <h2 style={{ fontSize: "18px", fontWeight: 600 }}>4. Policy Settings & Actions</h2>
                <div style={{ display: "flex", flexDirection: "column", gap: "8px" }}>
                  <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Enforcement Mode</label>
                  <select
                    value={form.settings.mode}
                    onChange={(e) => setForm({ ...form, settings: { ...form.settings, mode: e.target.value } })}
                    style={{ background: "rgba(255,255,255,0.03)", border: "1px solid rgba(255,255,255,0.1)", padding: "12px 16px", borderRadius: "8px", color: "white", outline: "none" }}
                  >
                    <option value="transparent" style={{ background: "#13141c" }}>Transparent (Log only, do not block)</option>
                    <option value="alarm" style={{ background: "#13141c" }}>Alarm (Log and trigger alerts)</option>
                    <option value="blocking" style={{ background: "#13141c" }}>Blocking (Drop malicious traffic)</option>
                  </select>
                </div>
                
                <div style={{ display: "flex", flexDirection: "column", gap: "8px" }}>
                  <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Anomaly Scoring Threshold</label>
                  <select
                    value={form.settings.severityThreshold}
                    onChange={(e) => setForm({ ...form, settings: { ...form.settings, severityThreshold: e.target.value } })}
                    style={{ background: "rgba(255,255,255,0.03)", border: "1px solid rgba(255,255,255,0.1)", padding: "12px 16px", borderRadius: "8px", color: "white", outline: "none" }}
                  >
                    <option value="low" style={{ background: "#13141c" }}>Low (Paranoia Level 1 - Recommended)</option>
                    <option value="medium" style={{ background: "#13141c" }}>Medium (Paranoia Level 2)</option>
                    <option value="high" style={{ background: "#13141c" }}>High (Paranoia Level 3 - Strict)</option>
                    <option value="critical" style={{ background: "#13141c" }}>Critical (Paranoia Level 4 - Very Strict)</option>
                  </select>
                </div>
              </div>
            )}

            {/* STEP 5: Review */}
            {step === 4 && (
              <div className="animate-fade-in" style={{ display: "flex", flexDirection: "column", gap: "24px" }}>
                <h2 style={{ fontSize: "18px", fontWeight: 600 }}>5. Review & Deploy</h2>
                
                <div style={{ background: "rgba(255,255,255,0.02)", border: "1px solid rgba(255,255,255,0.05)", borderRadius: "12px", padding: "24px" }}>
                  <div style={{ display: "flex", justifyContent: "space-between", marginBottom: "20px" }}>
                    <div>
                      <div style={{ fontSize: "20px", fontWeight: 600 }}>{form.name}</div>
                      <div style={{ color: "var(--text-secondary)", fontSize: "14px", marginTop: "4px" }}>Environment: <span style={{ textTransform: "capitalize", color: "white" }}>{form.environment}</span></div>
                    </div>
                    <div style={{ 
                      padding: "6px 12px", borderRadius: "999px", fontSize: "12px", fontWeight: 600, height: "fit-content",
                      background: form.settings.mode === 'blocking' ? 'rgba(239,68,68,0.1)' : 'rgba(59,130,246,0.1)',
                      color: form.settings.mode === 'blocking' ? '#f87171' : '#60a5fa',
                      border: `1px solid ${form.settings.mode === 'blocking' ? 'rgba(239,68,68,0.2)' : 'rgba(59,130,246,0.2)'}`
                    }}>
                      {form.settings.mode.toUpperCase()} MODE
                    </div>
                  </div>
                  
                  <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "24px" }}>
                    <div>
                      <div style={{ fontSize: "12px", color: "var(--text-secondary)", textTransform: "uppercase", letterSpacing: "0.05em", marginBottom: "12px" }}>Active Rules</div>
                      <ul style={{ margin: 0, padding: 0, listStyle: "none", display: "flex", flexDirection: "column", gap: "8px" }}>
                        {Object.entries(form.rules).filter(([, v]) => v).map(([k]) => (
                          <li key={k} style={{ fontSize: "14px", display: "flex", alignItems: "center", gap: "8px" }}>
                            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="var(--success)" strokeWidth="3"><polyline points="20 6 9 17 4 12"></polyline></svg>
                            <span style={{ textTransform: "uppercase" }}>{k}</span> Protection
                          </li>
                        ))}
                      </ul>
                    </div>
                    <div>
                      <div style={{ fontSize: "12px", color: "var(--text-secondary)", textTransform: "uppercase", letterSpacing: "0.05em", marginBottom: "12px" }}>Traffic Controls</div>
                      <ul style={{ margin: 0, padding: 0, listStyle: "none", display: "flex", flexDirection: "column", gap: "8px" }}>
                        {Object.entries(form.traffic).filter(([, v]) => v).map(([k]) => (
                          <li key={k} style={{ fontSize: "14px", display: "flex", alignItems: "center", gap: "8px" }}>
                            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="var(--accent-primary)" strokeWidth="3"><polyline points="20 6 9 17 4 12"></polyline></svg>
                            {k.replace(/([A-Z])/g, ' $1').replace(/^./, str => str.toUpperCase())}
                          </li>
                        ))}
                      </ul>
                    </div>
                  </div>
                </div>
                
                <div style={{ display: "flex", gap: "16px", marginTop: "16px" }}>
                  <button 
                    className="btn" 
                    onClick={() => submit(false)} 
                    disabled={submitting}
                    style={{ flex: 1, padding: "16px", border: "1px solid rgba(255,255,255,0.1)", background: "transparent" }}
                  >
                    Save as Draft (Transparent)
                  </button>
                  <button 
                    className="btn btn-primary" 
                    onClick={() => submit(true)} 
                    disabled={submitting}
                    style={{ flex: 1, padding: "16px" }}
                  >
                    {submitting ? "Deploying..." : "Activate Policy"}
                  </button>
                </div>
              </div>
            )}
            
          </div>

          {/* Footer Controls */}
          {step < STEPS.length - 1 && (
            <div style={{ display: "flex", justifyContent: "space-between", borderTop: "1px solid rgba(255,255,255,0.05)", padding: "24px 32px", background: "rgba(0,0,0,0.2)" }}>
              <button 
                className="btn" 
                onClick={prevStep} 
                disabled={step === 0}
                style={{ opacity: step === 0 ? 0.3 : 1 }}
              >
                Previous Step
              </button>
              
              <button 
                className="btn btn-primary" 
                onClick={nextStep} 
                disabled={step === 0 && !form.name}
              >
                Next Step
              </button>
            </div>
          )}

        </div>
      </main>
    </div>
  );
}
