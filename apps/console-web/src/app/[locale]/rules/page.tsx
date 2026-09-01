"use client";

import { useCallback, useEffect, useState } from "react";
import { useLocale } from "next-intl";
import Link from "next/link";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import DataTable, { type Column } from "@/components/shared/DataTable";
import { SeverityBadge, StatusBadge } from "@/components/shared/Badges";
import ConfirmDialog from "@/components/shared/ConfirmDialog";
import FilterBar from "@/components/shared/FilterBar";
import { usePermission } from "@/hooks/usePermission";
import {
  listRules,
  deleteRule,
  enableRule,
  disableRule,
  type RuleFull,
} from "@/lib/api";

function ExportRulesButton() {
  const [busy, setBusy] = useState(false);
  const doExport = async () => {
    setBusy(true);
    try {
      const r = await fetch(`${process.env.NEXT_PUBLIC_API_BASE ?? "http://127.0.0.1:8443"}/api/v1/rules/export`, { credentials: "include" });
      const blob = await r.blob();
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = "rules-export.json";
      a.click();
    } finally { setBusy(false); }
  };
  return <button className="btn" onClick={doExport} disabled={busy} style={{ padding: "8px 14px", fontSize: "13px" }}>{busy ? "Exporting…" : "Export"}</button>;
}

export default function RulesPage() {
  const locale = useLocale();
  const { can: canUser } = usePermission();
  const [rows, setRows] = useState<RuleFull[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [filters, setFilters] = useState<Record<string, string>>({});
  const [deleting, setDeleting] = useState<RuleFull | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      setRows(await listRules(filters));
    } catch {
      setError("Failed to load rules");
    } finally {
      setLoading(false);
    }
  }, [filters]);

  useEffect(() => {
    load();
  }, [load]);

  const toggle = async (rule: RuleFull, enable: boolean) => {
    try {
      if (enable) await enableRule(rule.id);
      else await disableRule(rule.id);
      await load();
    } catch {
      setError("Failed to update rule status");
    }
  };

  const confirmDelete = async () => {
    if (!deleting) return;
    setSubmitting(true);
    try {
      await deleteRule(deleting.id);
      setDeleting(null);
      await load();
    } catch {
      setError("Failed to delete rule");
    } finally {
      setSubmitting(false);
    }
  };

  const columns: Column<RuleFull>[] = [
    {
      key: "rule_id",
      label: "Rule ID",
      render: (row) => <span style={{ fontFamily: "monospace", fontSize: "12px", color: "var(--accent-primary)" }}>{row.rule_id}</span>,
    },
    {
      key: "name",
      label: "Rule Name",
      sortable: true,
      render: (row) => (
        <Link href={`/${locale}/rules/${row.id}`} style={{ color: "var(--text-primary)", textDecoration: "none", fontWeight: 500 }}>
          {row.name}
        </Link>
      ),
    },
    { key: "category", label: "Category", render: (row) => (row.category ? <StatusBadge value={row.category} /> : "—") },
    { key: "severity", label: "Severity", render: (row) => <SeverityBadge value={row.severity} /> },
    { key: "action", label: "Action", render: (row) => <StatusBadge value={row.action} /> },
    { key: "status", label: "Status", render: (row) => <StatusBadge value={row.status ?? "active"} /> },
    { key: "priority", label: "Priority", render: (row) => String(row.priority ?? "—") },
  ];

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />

      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <h1>Rules</h1>
            <p style={{ color: "var(--text-secondary)" }}>Manage WAF detection and protection rules.</p>
          </div>
          <div className="header-actions" style={{ display: "flex", alignItems: "center", gap: "12px" }}>
            {canUser("edit") && (
              <Link href={`/${locale}/rules/managed`} className="btn" style={{ textDecoration: "none", background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)", padding: "8px 16px" }}>
                Managed Rules (CRS)
              </Link>
            )}
            {canUser("create") && (
              <>
                <Link href={`/${locale}/rules/create`} className="btn btn-primary" style={{ textDecoration: "none", padding: "8px 16px" }}>
                  + New Rule
                </Link>
                <ExportRulesButton />
              </>
            )}
            <UserProfileWidget />
          </div>
        </div>

        <FilterBar
          filters={[
            { type: "text", key: "q", label: "Search", placeholder: "Search rules…" },
            { type: "select", key: "category", label: "Category", options: [
              { value: "sqli", label: "SQL Injection" },
              { value: "xss", label: "Cross-Site Scripting" },
              { value: "cmdi", label: "OS Command Injection" },
              { value: "rce", label: "Remote Code Execution" },
              { value: "pt", label: "Path Traversal" },
              { value: "fi", label: "File Inclusion" },
              { value: "xxe", label: "XXE / XML" },
              { value: "ssrf", label: "SSRF" },
              { value: "http", label: "HTTP Protocol" },
              { value: "hpp", label: "HTTP Parameter Pollution" },
              { value: "api", label: "API Security" },
              { value: "file_upload", label: "File Upload" },
              { value: "ldap", label: "LDAP Injection" },
              { value: "nosql", label: "NoSQL Injection" },
              { value: "ssti", label: "SSTI" },
              { value: "deserialization", label: "Deserialization" },
              { value: "scanner", label: "Scanner Detection" },
              { value: "bot", label: "Bot Protection" },
              { value: "auth", label: "Authentication Attacks" },
              { value: "session", label: "Session Security" },
              { value: "csrf", label: "CSRF" },
              { value: "info_disclosure", label: "Information Disclosure" },
              { value: "resource_discovery", label: "Resource Discovery" },
              { value: "request_anomaly", label: "Request Anomaly" },
              { value: "ip", label: "IP / Reputation" },
              { value: "geo", label: "Geo Security" },
              { value: "rate_limit", label: "Rate Limiting" },
              { value: "idor", label: "Insecure Direct Object Reference" },
              { value: "open_redirect", label: "Open Redirect" },
              { value: "clickjacking", label: "Clickjacking" },
              { value: "jwt", label: "JWT Misconfiguration" },
              { value: "prototype_pollution", label: "Prototype Pollution" },
              { value: "cors", label: "CORS Misconfiguration" },
            ]},
            { type: "select", key: "type", label: "Type", options: [
              { value: "managed", label: "Managed" },
              { value: "custom", label: "Custom" },
            ]},
            { type: "select", key: "severity", label: "Severity", options: [
              { value: "critical", label: "Critical" },
              { value: "high", label: "High" },
              { value: "medium", label: "Medium" },
              { value: "low", label: "Low" },
            ]},
            { type: "select", key: "action", label: "Action", options: [
              { value: "allow", label: "Allow" },
              { value: "log", label: "Log" },
              { value: "block", label: "Block" },
              { value: "challenge", label: "Challenge" },
              { value: "rate_limit", label: "Rate Limit" },
            ]},
            { type: "select", key: "status", label: "Status", options: [
              { value: "active", label: "Active" },
              { value: "disabled", label: "Disabled" },
            ]},
          ]}
          values={filters}
          onChange={setFilters}
        />

        <div className="data-section animate-fade-in delay-1">
          <DataTable
            columns={columns}
            data={rows}
            rowKey={(row) => row.id}
            loading={loading}
            error={error || undefined}
            onRetry={load}
            emptyMessage="No rules found."
            actions={(row) => (
              <div style={{ display: "flex", gap: "6px", flexWrap: "wrap" }}>
                <Link href={`/${locale}/rules/${row.id}`} className="btn" style={{ padding: "6px 10px", fontSize: "12px", textDecoration: "none" }}>
                  View
                </Link>
                {canUser("edit") && (
                  <>
                    {row.status === "active" ? (
                      <button type="button" className="btn" onClick={() => toggle(row, false)} style={{ padding: "6px 10px", fontSize: "12px", color: "var(--warning)" }}>
                        Disable
                      </button>
                    ) : (
                      <button type="button" className="btn" onClick={() => toggle(row, true)} style={{ padding: "6px 10px", fontSize: "12px", color: "var(--success)" }}>
                        Enable
                      </button>
                    )}
                    <Link href={`/${locale}/rules/${row.id}?test=1`} className="btn" style={{ padding: "6px 10px", fontSize: "12px", textDecoration: "none" }}>
                      Test Rule
                    </Link>
                    <button type="button" className="btn" onClick={() => setDeleting(row)} style={{ padding: "6px 10px", fontSize: "12px", color: "var(--danger)" }}>
                      Delete
                    </button>
                  </>
                )}
              </div>
            )}
          />
        </div>
      </main>

      <ConfirmDialog
        open={!!deleting}
        title="Delete rule"
        message={`Are you sure you want to delete "${deleting?.name}" (${deleting?.rule_id})?`}
        confirmLabel="Delete"
        danger
        loading={submitting}
        onConfirm={confirmDelete}
        onCancel={() => setDeleting(null)}
      />
    </div>
  );
}
