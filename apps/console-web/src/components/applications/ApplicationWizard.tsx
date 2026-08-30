"use client";

import { useState } from "react";
import { createApplication, type ApplicationCreatePayload } from "@/lib/api";

type FormState = ApplicationCreatePayload & {
  tagsInput: string;
};

const emptyForm = (): FormState => ({
  name: "",
  description: "",
  environment: "production",
  status: "active",
  tagsInput: "",
  tags: [],
  domain: "",
  origin_type: "ip",
  origin_host: "",
  origin_port: 443,
  origin_protocol: "https",
  origin_path: "/",
  origin_load_balancing: "single",
  waf_policy_id: "",
  waf_mode: "block",
  tls_enabled: false,
  certificate_id: "",
  min_tls_version: "1.2",
  http_redirect: false,
  rate_limit_enabled: false,
  rate_limit: 1000,
  health_check_enabled: false,
  health_check_method: "GET",
  health_check_path: "/health",
  health_check_interval: 30,
  health_check_timeout: 5,
  health_check_retries: 3,
  health_check_expected_status: 200,
  request_body_limit_mb: 10,
  connection_timeout_s: 30,
  keep_alive: true,
  real_client_ip_header: "X-Forwarded-For",
  log_request_headers: true,
  log_response_status: true,
});

const STEPS = ["Basic Information", "Domain & Origin", "Security & WAF", "Health Check & Advanced", "Review"];

type Props = {
  onClose: () => void;
  onCreated: () => void;
};

export default function ApplicationWizard({ onClose, onCreated }: Props) {
  const [step, setStep] = useState(0);
  const [form, setForm] = useState<FormState>(emptyForm);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");

  const set = <K extends keyof FormState>(key: K, value: FormState[K]) =>
    setForm((prev) => ({ ...prev, [key]: value }));

  const addTag = (v: string) => {
    const tag = v.trim().replace(/,+$/, "");
    if (tag && !(form.tags ?? []).includes(tag)) {
      set("tags", [...(form.tags ?? []), tag]);
    }
    set("tagsInput", "");
  };

  const validateStep = (): boolean => {
    if (step === 0) {
      if (!form.name) { setError("Application name is required."); return false; }
      if (!form.environment) { setError("Environment is required."); return false; }
    }
    if (step === 1) {
      if (!form.domain) { setError("Domain / hostname is required."); return false; }
      if (!form.origin_host) { setError("Origin address is required."); return false; }
    }
    setError("");
    return true;
  };

  const next = () => {
    if (!validateStep()) return;
    setStep((s) => Math.min(s + 1, STEPS.length - 1));
  };
  const back = () => setStep((s) => Math.max(s - 1, 0));

  const submit = async () => {
    setSubmitting(true);
    setError("");
    try {
      const payload: ApplicationCreatePayload = { ...form };
      delete (payload as Partial<FormState>).tagsInput;
      await createApplication(payload);
      onCreated();
    } catch {
      setError("Failed to create application.");
    } finally {
      setSubmitting(false);
    }
  };

  const field = (label: string, children: React.ReactNode, required = false) => (
    <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
      <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>
        {label} {required && <span style={{ color: "var(--danger)" }}>*</span>}
      </label>
      {children}
    </div>
  );

  const input: React.CSSProperties = {
    background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)",
    padding: "10px 12px", borderRadius: "8px", color: "white", outline: "none", fontSize: "14px", width: "100%",
  };
  const select: React.CSSProperties = { ...input, appearance: "auto" };
  const toggle = { display: "flex", alignItems: "center", gap: "10px" } as React.CSSProperties;

  return (
    <div onClick={onClose} style={{ position: "fixed", inset: 0, zIndex: 1200, background: "rgba(0,0,0,0.7)", display: "flex", alignItems: "center", justifyContent: "center", padding: "20px", overflowY: "auto" }}>
      <div className="glass-panel animate-fade-in" onClick={(e) => e.stopPropagation()} style={{ width: "100%", maxWidth: "640px", padding: "28px", display: "flex", flexDirection: "column", gap: "18px", maxHeight: "92vh", overflowY: "auto" }}>
        <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
          <h3 style={{ fontSize: "18px", fontWeight: 700 }}>New Application</h3>
          <button onClick={onClose} style={{ background: "none", border: "none", color: "var(--text-secondary)", fontSize: "20px", cursor: "pointer" }}>×</button>
        </div>

        {/* Stepper */}
        <div style={{ display: "flex", gap: "6px", flexWrap: "wrap" }}>
          {STEPS.map((s, i) => (
            <div key={s} style={{ display: "flex", alignItems: "center", gap: "6px" }}>
              <span
                style={{
                  width: "26px", height: "26px", borderRadius: "50%", display: "grid", placeItems: "center",
                  fontSize: "12px", fontWeight: 700,
                  background: i < step ? "var(--success)" : i === step ? "var(--accent-primary)" : "rgba(255,255,255,0.08)",
                  color: i <= step ? "#fff" : "var(--text-secondary)",
                }}
              >
                {i < step ? "✓" : i + 1}
              </span>
              <span style={{ fontSize: "12px", color: i === step ? "var(--text-primary)" : "var(--text-secondary)", display: i === step ? "inline" : "none" }}>
                {s}
              </span>
            </div>
          ))}
        </div>

        {error && <div style={{ padding: "10px 14px", background: "rgba(239,68,68,0.1)", border: "1px solid rgba(239,68,68,0.3)", color: "var(--danger)", borderRadius: "8px", fontSize: "13px" }}>{error}</div>}

        {/* STEP 0 — Basic */}
        {step === 0 && (
          <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
            {field("Application Name", <input style={input} value={form.name} onChange={(e) => set("name", e.target.value)} placeholder="ERP Portal" />, true)}
            {field("Description", <textarea style={{ ...input, resize: "vertical" }} rows={3} value={form.description} onChange={(e) => set("description", e.target.value)} placeholder="Describe this application" />)}
            {field("Environment", (
              <select style={select} value={form.environment} onChange={(e) => set("environment", e.target.value)}>
                <option value="production">Production</option>
                <option value="staging">Staging</option>
                <option value="development">Development</option>
              </select>
            ), true)}
            <div style={toggle}>
              <input type="checkbox" checked={form.status === "active"} onChange={(e) => set("status", e.target.checked ? "active" : "disabled")} />
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Active</label>
            </div>
            {field("Tags", (
              <div>
                <div style={{ display: "flex", gap: "8px" }}>
                  <input style={input} value={form.tagsInput} onChange={(e) => set("tagsInput", e.target.value)} onKeyDown={(e) => { if (e.key === "Enter" || e.key === ",") { e.preventDefault(); addTag(form.tagsInput); } }} placeholder="Type a tag, press Enter" />
                  <button className="btn" onClick={() => addTag(form.tagsInput)} style={{ padding: "8px 14px" }}>Add</button>
                </div>
                <div style={{ display: "flex", gap: "6px", marginTop: "8px", flexWrap: "wrap" }}>
                  {(form.tags ?? []).map((tg) => (
                    <span key={tg} style={{ padding: "4px 10px", borderRadius: "6px", background: "rgba(59,130,246,0.15)", color: "#60a5fa", fontSize: "12px" }}>
                      {tg} <button style={{ background: "none", border: "none", color: "inherit", cursor: "pointer", marginLeft: "4px" }} onClick={() => set("tags", (form.tags ?? []).filter((x) => x !== tg))}>×</button>
                    </span>
                  ))}
                </div>
              </div>
            ))}
          </div>
        )}

        {/* STEP 1 — Domain & Origin */}
        {step === 1 && (
          <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
            {field("Domain / Hostname", <input style={input} value={form.domain} onChange={(e) => set("domain", e.target.value)} placeholder="erp.example.com" />, true)}
            {field("Origin Type", (
              <div style={{ display: "flex", gap: "10px" }}>
                {(["ip", "hostname", "load_balancer"] as const).map((ot) => (
                  <label key={ot} style={{ display: "flex", alignItems: "center", gap: "6px", fontSize: "14px", color: "var(--text-secondary)", cursor: "pointer" }}>
                    <input type="radio" name="origin_type" checked={form.origin_type === ot} onChange={() => set("origin_type", ot)} />
                    {ot === "ip" ? "IP Address" : ot === "hostname" ? "Hostname" : "Load Balancer"}
                  </label>
                ))}
              </div>
            ))}
            {field("Origin Address", <input style={input} value={form.origin_host} onChange={(e) => set("origin_host", e.target.value)} placeholder="10.24.20.50" />, true)}
            <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "12px" }}>
              {field("Origin Port", <input type="number" style={input} value={form.origin_port} onChange={(e) => set("origin_port", Number(e.target.value))} />)}
              {field("Protocol", (
                <select style={select} value={form.origin_protocol} onChange={(e) => set("origin_protocol", e.target.value)}>
                  <option value="http">HTTP</option>
                  <option value="https">HTTPS</option>
                </select>
              ))}
            </div>
            <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "12px" }}>
              {field("Origin Path", <input style={input} value={form.origin_path} onChange={(e) => set("origin_path", e.target.value)} placeholder="/" />)}
              {field("Load Balancing", (
                <select style={select} value={form.origin_load_balancing} onChange={(e) => set("origin_load_balancing", e.target.value)}>
                  <option value="single">Single Origin</option>
                  <option value="round_robin">Round Robin</option>
                  <option value="ip_hash">IP Hash</option>
                </select>
              ))}
            </div>
          </div>
        )}

        {/* STEP 2 — Security & WAF */}
        {step === 2 && (
          <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
            {field("WAF Mode", (
              <div style={{ display: "flex", gap: "10px" }}>
                {(["block", "detection", "disabled"] as const).map((m) => (
                  <label key={m} style={{ display: "flex", alignItems: "center", gap: "6px", fontSize: "14px", color: "var(--text-secondary)", cursor: "pointer" }}>
                    <input type="radio" name="waf_mode" checked={form.waf_mode === m} onChange={() => set("waf_mode", m)} />
                    {m === "block" ? "Block" : m === "detection" ? "Detection Only" : "Disabled"}
                  </label>
                ))}
              </div>
            ))}
            {field("Security Policy", <input style={input} value={form.waf_policy_id} onChange={(e) => set("waf_policy_id", e.target.value)} placeholder="Default Web Protection (policy ID)" />)}
            {field("Enable HTTPS", (
              <div style={toggle}>
                <input type="checkbox" checked={form.tls_enabled} onChange={(e) => set("tls_enabled", e.target.checked)} />
                <span style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Terminate TLS on the WAF</span>
              </div>
            ))}
            {form.tls_enabled && (
              <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "12px" }}>
                {field("Certificate ID", <input style={input} value={form.certificate_id} onChange={(e) => set("certificate_id", e.target.value)} placeholder="Upload under Certificates first" />)}
                {field("Min TLS Version", (
                  <select style={select} value={form.min_tls_version} onChange={(e) => set("min_tls_version", e.target.value)}>
                    <option value="1.0">TLS 1.0</option>
                    <option value="1.1">TLS 1.1</option>
                    <option value="1.2">TLS 1.2</option>
                    <option value="1.3">TLS 1.3</option>
                  </select>
                ))}
              </div>
            )}
            <div style={toggle}>
              <input type="checkbox" checked={form.http_redirect} onChange={(e) => set("http_redirect", e.target.checked)} />
              <span style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Redirect HTTP → HTTPS</span>
            </div>
            <div style={toggle}>
              <input type="checkbox" checked={form.rate_limit_enabled} onChange={(e) => set("rate_limit_enabled", e.target.checked)} />
              <span style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Enable Rate Limiting</span>
            </div>
            {form.rate_limit_enabled && (
              field("Rate Limit (requests / minute)", <input type="number" style={input} value={form.rate_limit} onChange={(e) => set("rate_limit", Number(e.target.value))} />)
            )}
          </div>
        )}

        {/* STEP 3 — Health Check & Advanced */}
        {step === 3 && (
          <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
            <div style={toggle}>
              <input type="checkbox" checked={form.health_check_enabled} onChange={(e) => set("health_check_enabled", e.target.checked)} />
              <span style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Enable Health Check</span>
            </div>
            {form.health_check_enabled && (
              <>
                <div style={{ display: "grid", gridTemplateColumns: "1fr 2fr", gap: "12px" }}>
                  {field("Method", (
                    <select style={select} value={form.health_check_method} onChange={(e) => set("health_check_method", e.target.value)}>
                      <option value="GET">GET</option>
                      <option value="HEAD">HEAD</option>
                      <option value="TCP">TCP</option>
                    </select>
                  ))}
                  {field("Path", <input style={input} value={form.health_check_path} onChange={(e) => set("health_check_path", e.target.value)} placeholder="/health" />)}
                </div>
                <div style={{ display: "grid", gridTemplateColumns: "repeat(2, 1fr)", gap: "12px" }}>
                  {field("Expected Status", <input type="number" style={input} value={form.health_check_expected_status} onChange={(e) => set("health_check_expected_status", Number(e.target.value))} />)}
                  {field("Interval (s)", <input type="number" style={input} value={form.health_check_interval} onChange={(e) => set("health_check_interval", Number(e.target.value))} />)}
                  {field("Timeout (s)", <input type="number" style={input} value={form.health_check_timeout} onChange={(e) => set("health_check_timeout", Number(e.target.value))} />)}
                  {field("Retries", <input type="number" style={input} value={form.health_check_retries} onChange={(e) => set("health_check_retries", Number(e.target.value))} />)}
                </div>
              </>
            )}

            <h4 style={{ fontSize: "15px", fontWeight: 600, marginTop: "8px" }}>Advanced Settings</h4>
            <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "12px" }}>
              {field("Request Body Limit (MB)", <input type="number" style={input} value={form.request_body_limit_mb} onChange={(e) => set("request_body_limit_mb", Number(e.target.value))} />)}
              {field("Connection Timeout (s)", <input type="number" style={input} value={form.connection_timeout_s} onChange={(e) => set("connection_timeout_s", Number(e.target.value))} />)}
              {field("Real Client IP Header", (
                <select style={select} value={form.real_client_ip_header} onChange={(e) => set("real_client_ip_header", e.target.value)}>
                  <option value="X-Forwarded-For">X-Forwarded-For</option>
                  <option value="X-Real-IP">X-Real-IP</option>
                  <option value="CF-Connecting-IP">CF-Connecting-IP</option>
                </select>
              ))}
              <div />
            </div>
            <div style={{ display: "flex", flexDirection: "column", gap: "8px" }}>
              <div style={toggle}><input type="checkbox" checked={form.keep_alive} onChange={(e) => set("keep_alive", e.target.checked)} /><span style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Keep-Alive</span></div>
              <div style={toggle}><input type="checkbox" checked={form.log_request_headers} onChange={(e) => set("log_request_headers", e.target.checked)} /><span style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Log Request Headers</span></div>
              <div style={toggle}><input type="checkbox" checked={form.log_response_status} onChange={(e) => set("log_response_status", e.target.checked)} /><span style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Log Response Status</span></div>
            </div>
          </div>
        )}

        {/* STEP 4 — Review */}
        {step === 4 && (
          <div className="glass-panel" style={{ padding: "20px", border: "1px solid var(--glass-border)" }}>
            <h4 style={{ fontSize: "15px", fontWeight: 600, marginBottom: "12px" }}>Review Application</h4>
            <table style={{ width: "100%", borderCollapse: "collapse", fontSize: "14px" }}>
              <tbody>
                {[
                  ["Name", form.name], ["Environment", form.environment], ["Domain", form.domain],
                  ["Origin", `${form.origin_host}:${form.origin_port}`], ["Protocol", form.origin_protocol],
                  ["WAF Mode", form.waf_mode], ["TLS", form.tls_enabled ? `Enabled (${form.min_tls_version})` : "Disabled"],
                  ["Rate Limit", form.rate_limit_enabled ? `${form.rate_limit} req/min` : "Disabled"],
                  ["Health Check", form.health_check_enabled ? `${form.health_check_method} ${form.health_check_path}` : "Disabled"],
                  ["Status", form.status],
                ].map(([k, v]) => (
                  <tr key={k} style={{ borderBottom: "1px solid rgba(255,255,255,0.06)" }}>
                    <td style={{ padding: "8px 0", color: "var(--text-secondary)" }}>{k}</td>
                    <td style={{ padding: "8px 0", textAlign: "right", color: "var(--text-primary)" }}>{String(v || "—")}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {/* Nav buttons */}
        <div style={{ display: "flex", justifyContent: "space-between", marginTop: "8px" }}>
          <div>
            {step > 0 && <button type="button" className="btn" onClick={back} style={{ padding: "10px 18px" }}>← Back</button>}
          </div>
          <div style={{ display: "flex", gap: "10px" }}>
            <button type="button" className="btn" onClick={onClose} disabled={submitting} style={{ padding: "10px 18px" }}>Cancel</button>
            {step < STEPS.length - 1 ? (
              <button type="button" className="btn btn-primary" onClick={next} style={{ padding: "10px 18px" }}>Next →</button>
            ) : (
              <button type="button" className="btn btn-primary" onClick={submit} disabled={submitting} style={{ padding: "10px 18px" }}>
                {submitting ? "Creating…" : "Create Application"}
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
