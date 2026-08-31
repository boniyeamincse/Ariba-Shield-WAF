"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { createRule, type RuleCondition, type RuleScope } from "@/lib/api";

const STEPS = ["Basic", "Match", "Conditions", "Action", "Scope", "Review"];

const FIELDS = ["url", "query_param", "header", "cookie", "request_body", "method", "source_ip", "user_agent", "host"];
const OPERATORS = ["equals", "not_equals", "contains", "not_contains", "starts_with", "ends_with", "regex", "ip_match", "cidr_match", "gt", "lt"];
const TRANSFORMS = ["", "lowercase", "url_decode", "base64_decode"];
const CATEGORIES = [
  "sqli", "xss", "pt", "cmdi", "fi", "ssrf", "xxe", "rce", "http",
  "scanner", "bot", "rate_limit", "ip", "geo", "header", "custom",
];
const METHODS = ["GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"];

type FormState = {
  name: string;
  rule_id: string;
  description: string;
  type: string;
  category: string;
  severity: string;
  priority: number;
  status: string;
  conditions: RuleCondition[];
  logic: string;
  action: string;
  status_code: number;
  response_type: string;
  message: string;
  scopes: RuleScope[];
};

const emptyForm = (): FormState => ({
  name: "", rule_id: "", description: "", type: "custom", category: "sqli",
  severity: "medium", priority: 10, status: "active",
  conditions: [{ group_id: 0, field: "request_body", operator: "contains", value: "", transformation: "lowercase", case_sensitive: false }],
  logic: "AND", action: "block", status_code: 403, response_type: "json",
  message: "Request blocked by security policy",
  scopes: [{ path_pattern: "/*", methods: ["GET", "POST", "PUT", "PATCH", "DELETE"] }],
});

export default function RuleWizard({ onCancel }: { onCancel?: () => void }) {
  const router = useRouter();
  const [step, setStep] = useState(0);
  const [form, setForm] = useState<FormState>(emptyForm);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");

  const set = <K extends keyof FormState>(k: K, v: FormState[K]) => setForm((p) => ({ ...p, [k]: v }));

  const validate = (): boolean => {
    if (step === 0) {
      if (!form.name || !form.rule_id) { setError("Rule name and Rule ID are required."); return false; }
    }
    if (step === 1 || step === 2) {
      if (!form.conditions.length || form.conditions.some((c) => !c.value)) {
        setError("At least one condition with a value is required."); return false;
      }
    }
    setError("");
    return true;
  };

  const next = () => { if (validate()) setStep((s) => Math.min(s + 1, STEPS.length - 1)); };
  const back = () => setStep((s) => Math.max(s - 1, 0));

  const submit = async () => {
    setSubmitting(true); setError("");
    try {
      await createRule({
        rule_id: form.rule_id, name: form.name, description: form.description,
        type: form.type, category: form.category, severity: form.severity,
        priority: form.priority, action: form.action, status: form.status,
        logic: form.logic, conditions: form.conditions, scopes: form.scopes,
      });
      router.push("/en/rules");
      router.refresh();
    } catch {
      setError("Failed to create rule.");
    } finally { setSubmitting(false); }
  };

  const field = (label: string, children: React.ReactNode, required = false) => (
    <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
      <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>{label} {required && <span style={{ color: "var(--danger)" }}>*</span>}</label>
      {children}
    </div>
  );
  const input: React.CSSProperties = {
    background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)",
    padding: "10px 12px", borderRadius: "8px", color: "white", outline: "none", fontSize: "14px", width: "100%",
  };
  const select: React.CSSProperties = { ...input };

  const updateCond = (i: number, patch: Partial<RuleCondition>) => {
    set("conditions", form.conditions.map((c, idx) => (idx === i ? { ...c, ...patch } : c)));
  };

  const categoryLabel = (c: string) => c.replace(/_/g, " ").replace(/^\w/, (x) => x.toUpperCase());

  return (
    <div className="glass-panel animate-fade-in" style={{ width: "100%", maxWidth: "680px", padding: "28px", display: "flex", flexDirection: "column", gap: "18px", margin: "0 auto" }}>
      <h3 style={{ fontSize: "18px", fontWeight: 700 }}>New Rule</h3>

      {/* Stepper */}
      <div style={{ display: "flex", gap: "6px", flexWrap: "wrap" }}>
        {STEPS.map((s, i) => (
          <div key={s} style={{ display: "flex", alignItems: "center", gap: "4px" }}>
            <span style={{
              width: "24px", height: "24px", borderRadius: "50%", display: "grid", placeItems: "center",
              fontSize: "11px", fontWeight: 700,
              background: i < step ? "var(--success)" : i === step ? "var(--accent-primary)" : "rgba(255,255,255,0.08)",
              color: i <= step ? "#fff" : "var(--text-secondary)",
            }}>{i < step ? "✓" : i + 1}</span>
            <span style={{ fontSize: "11px", color: i === step ? "var(--text-primary)" : "var(--text-secondary)", display: i === step ? "inline" : "none" }}>{s}</span>
          </div>
        ))}
      </div>

      {error && <div style={{ padding: "10px 14px", background: "rgba(239,68,68,0.1)", border: "1px solid rgba(239,68,68,0.3)", color: "var(--danger)", borderRadius: "8px", fontSize: "13px" }}>{error}</div>}

      {/* Step 0 — Basic */}
      {step === 0 && (
        <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
          {field("Rule Name", <input style={input} value={form.name} onChange={(e) => set("name", e.target.value)} placeholder="SQL Injection Detection" />, true)}
          {field("Rule ID", <input style={input} value={form.rule_id} onChange={(e) => set("rule_id", e.target.value)} placeholder="ARB-SQL-001" />, true)}
          {field("Description", <textarea style={{ ...input, resize: "vertical" }} rows={2} value={form.description} onChange={(e) => set("description", e.target.value)} />)}
          <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "12px" }}>
            {field("Rule Type", (
              <select style={select} value={form.type} onChange={(e) => set("type", e.target.value)}>
                <option value="custom">Custom Rule</option>
                <option value="managed">Managed Rule</option>
              </select>
            ), true)}
            {field("Category", (
              <select style={select} value={form.category} onChange={(e) => set("category", e.target.value)}>
                {CATEGORIES.map((c) => <option key={c} value={c}>{categoryLabel(c)}</option>)}
              </select>
            ), true)}
          </div>
          <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "12px" }}>
            {field("Severity", (
              <select style={select} value={form.severity} onChange={(e) => set("severity", e.target.value)}>
                {["critical", "high", "medium", "low"].map((s) => <option key={s} value={s}>{s[0].toUpperCase() + s.slice(1)}</option>)}
              </select>
            ), true)}
            {field("Priority", <input type="number" style={input} value={form.priority} onChange={(e) => set("priority", Number(e.target.value))} />, true)}
          </div>
          {field("Status", (
            <div style={{ display: "flex", alignItems: "center", gap: "8px" }}>
              <input type="checkbox" checked={form.status === "active"} onChange={(e) => set("status", e.target.checked ? "active" : "disabled")} />
              <span style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Enabled</span>
            </div>
          ))}
        </div>
      )}

      {/* Step 1 — Match */}
      {step === 1 && (
        <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
          {field("Match Target", (
            <select style={select} value={form.conditions[0]?.field} onChange={(e) => updateCond(0, { field: e.target.value })}>
              {FIELDS.map((f) => <option key={f} value={f}>{categoryLabel(f)}</option>)}
            </select>
          ), true)}
          {field("Operator", (
            <select style={select} value={form.conditions[0]?.operator} onChange={(e) => updateCond(0, { operator: e.target.value })}>
              {OPERATORS.map((o) => <option key={o} value={o}>{categoryLabel(o)}</option>)}
            </select>
          ), true)}
          {field("Value", <input style={input} value={form.conditions[0]?.value ?? ""} onChange={(e) => updateCond(0, { value: e.target.value })} placeholder="union select" />, true)}
          <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "12px" }}>
            {field("Transformation", (
              <select style={select} value={form.conditions[0]?.transformation ?? ""} onChange={(e) => updateCond(0, { transformation: e.target.value })}>
                {TRANSFORMS.map((t) => <option key={t} value={t}>{t ? categoryLabel(t) : "None"}</option>)}
              </select>
            ))}
            {field("Case Sensitive", (
              <select style={select} value={form.conditions[0]?.case_sensitive ? "yes" : "no"} onChange={(e) => updateCond(0, { case_sensitive: e.target.value === "yes" })}>
                <option value="no">No</option>
                <option value="yes">Yes</option>
              </select>
            ))}
          </div>
        </div>
      )}

      {/* Step 2 — Conditions */}
      {step === 2 && (
        <div style={{ display: "flex", flexDirection: "column", gap: "12px" }}>
          <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
            <span style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Conditions ({form.conditions.length})</span>
            <button className="btn" onClick={() => set("conditions", [...form.conditions, { group_id: 0, field: "method", operator: "equals", value: "", transformation: "", case_sensitive: false }])} style={{ padding: "6px 12px", fontSize: "12px" }}>+ Add Condition</button>
          </div>
          {form.conditions.map((c, i) => (
            <div key={i} className="glass-panel" style={{ padding: "14px", display: "grid", gridTemplateColumns: "1fr 1fr 1.4fr auto", gap: "8px", alignItems: "center" }}>
              <select style={select} value={c.field} onChange={(e) => updateCond(i, { field: e.target.value })}>
                {FIELDS.map((f) => <option key={f} value={f}>{categoryLabel(f)}</option>)}
              </select>
              <select style={select} value={c.operator} onChange={(e) => updateCond(i, { operator: e.target.value })}>
                {OPERATORS.map((o) => <option key={o} value={o}>{categoryLabel(o)}</option>)}
              </select>
              <input style={input} value={c.value} onChange={(e) => updateCond(i, { value: e.target.value })} placeholder="value" />
              <button className="btn" onClick={() => set("conditions", form.conditions.filter((_, idx) => idx !== i))} style={{ padding: "6px 10px", fontSize: "12px", color: "var(--danger)" }}>✕</button>
            </div>
          ))}
          <div style={{ display: "flex", gap: "10px", alignItems: "center" }}>
            <span style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Logical operator between conditions:</span>
            {["AND", "OR"].map((l) => (
              <label key={l} style={{ fontSize: "14px", color: "var(--text-secondary)", cursor: "pointer" }}>
                <input type="radio" name="logic" checked={form.logic === l} onChange={() => set("logic", l)} /> {l}
              </label>
            ))}
          </div>
        </div>
      )}

      {/* Step 3 — Action */}
      {step === 3 && (
        <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
          {field("Action", (
            <div style={{ display: "flex", gap: "12px", flexWrap: "wrap" }}>
              {["allow", "log", "block", "challenge", "rate_limit"].map((a) => (
                <label key={a} style={{ fontSize: "14px", color: "var(--text-secondary)", cursor: "pointer" }}>
                  <input type="radio" name="action" checked={form.action === a} onChange={() => set("action", a)} /> {categoryLabel(a)}
                </label>
              ))}
            </div>
          ), true)}
          {form.action === "block" && (
            <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "12px" }}>
              {field("HTTP Status", <input type="number" style={input} value={form.status_code} onChange={(e) => set("status_code", Number(e.target.value))} />)}
              {field("Response Type", (
                <select style={select} value={form.response_type} onChange={(e) => set("response_type", e.target.value)}>
                  <option value="json">JSON</option>
                  <option value="html">HTML</option>
                  <option value="text">Text</option>
                </select>
              ))}
            </div>
          )}
          {field("Response Message", <input style={input} value={form.message} onChange={(e) => set("message", e.target.value)} />)}
        </div>
      )}

      {/* Step 4 — Scope */}
      {step === 4 && (
        <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
          {field("Apply Path", <input style={input} value={form.scopes[0]?.path_pattern ?? "/*"} onChange={(e) => set("scopes", [{ ...form.scopes[0], path_pattern: e.target.value }])} placeholder="/*" />)}
          {field("Methods", (
            <div style={{ display: "flex", gap: "10px", flexWrap: "wrap" }}>
              {METHODS.map((m) => {
                const selected = form.scopes[0]?.methods?.includes(m) ?? false;
                return (
                  <label key={m} style={{ fontSize: "14px", color: "var(--text-secondary)", cursor: "pointer" }}>
                    <input
                      type="checkbox"
                      checked={selected}
                      onChange={(e) => {
                        const cur = form.scopes[0]?.methods ?? [];
                        const next = e.target.checked ? [...cur, m] : cur.filter((x) => x !== m);
                        set("scopes", [{ ...form.scopes[0], methods: next }]);
                      }}
                    /> {m}
                  </label>
                );
              })}
            </div>
          ))}
        </div>
      )}

      {/* Step 5 — Review */}
      {step === 5 && (
        <div className="glass-panel" style={{ padding: "20px", border: "1px solid var(--glass-border)" }}>
          <h4 style={{ fontSize: "15px", fontWeight: 600, marginBottom: "12px" }}>Review Rule</h4>
          <table style={{ width: "100%", borderCollapse: "collapse", fontSize: "14px" }}>
            <tbody>
              {[
                ["Name", form.name], ["ID", form.rule_id], ["Category", categoryLabel(form.category)],
                ["Severity", form.severity], ["Priority", String(form.priority)],
                ["Logic", form.logic], ["Action", form.action.toUpperCase()],
                ["Status", form.status],
                ["Match", form.conditions.map((c) => `${c.field} ${c.operator} "${c.value}"`).join(` ${form.logic} `)],
                ["Path", form.scopes[0]?.path_pattern ?? "/*"],
                ["Methods", (form.scopes[0]?.methods ?? []).join(", ") || "All"],
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

      {/* Nav */}
      <div style={{ display: "flex", justifyContent: "space-between", marginTop: "8px" }}>
        <div>{step > 0 && <button type="button" className="btn" onClick={back} style={{ padding: "10px 18px" }}>← Back</button>}</div>
        <div style={{ display: "flex", gap: "10px" }}>
          {onCancel && <button type="button" className="btn" onClick={onCancel} disabled={submitting} style={{ padding: "10px 18px" }}>Cancel</button>}
          {step < STEPS.length - 1 ? (
            <button type="button" className="btn btn-primary" onClick={next} style={{ padding: "10px 18px" }}>Next →</button>
          ) : (
            <button type="button" className="btn btn-primary" onClick={submit} disabled={submitting} style={{ padding: "10px 18px" }}>
              {submitting ? "Saving…" : "Save Rule"}
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
