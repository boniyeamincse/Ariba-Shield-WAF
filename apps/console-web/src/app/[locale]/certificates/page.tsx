"use client";

import { useCallback, useEffect, useState } from "react";
import type { CSSProperties } from "react";
import { useLocale, useTranslations } from "next-intl";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import DataTable, { type Column } from "@/components/shared/DataTable";
import { StatusBadge } from "@/components/shared/Badges";
import ConfirmDialog from "@/components/shared/ConfirmDialog";
import {
  listCertificates,
  deleteCertificate,
  provisionACME,
  type Certificate,
} from "@/lib/api";

export default function CertificatesPage() {
  const locale = useLocale();
  const t = useTranslations("certificates");
  const [rows, setRows] = useState<Certificate[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [deleting, setDeleting] = useState<Certificate | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [acmeOpen, setAcmeOpen] = useState(false);
  const [acmeForm, setAcmeForm] = useState({ domain: "", email: "", staging: false });

  const load = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      setRows(await listCertificates());
    } catch {
      setError("Failed to load certificates");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const confirmDelete = async () => {
    if (!deleting) return;
    setSubmitting(true);
    try {
      await deleteCertificate(deleting.id);
      setDeleting(null);
      await load();
    } catch {
      setError("Failed to delete certificate");
    } finally {
      setSubmitting(false);
    }
  };

  const handleProvision = async () => {
    setSubmitting(true);
    setError("");
    try {
      await provisionACME(acmeForm.domain, acmeForm.email, acmeForm.staging);
      setAcmeOpen(false);
      setAcmeForm({ domain: "", email: "", staging: false });
      await load();
    } catch {
      setError("ACME provisioning failed. Ensure the domain resolves and port 80 is reachable.");
    } finally {
      setSubmitting(false);
    }
  };

  const columns: Column<Certificate>[] = [
    { key: "name", label: "Name", sortable: true },
    { key: "domain", label: "Domain" },
    { key: "issuer", label: "Issuer", render: (row) => row.issuer || "—" },
    {
      key: "not_after",
      label: "Expiry",
      render: (row) => {
        if (!row.not_after) return "—";
        const d = new Date(row.not_after);
        const expired = d < new Date();
        const soon = !expired && d < new Date(Date.now() + 30 * 86400000);
        return (
          <span style={{ color: expired ? "var(--danger)" : soon ? "var(--warning)" : "var(--text-secondary)" }}>
            {d.toLocaleDateString()}
          </span>
        );
      },
    },
    { key: "status", label: "Status", render: (row) => <StatusBadge value={row.status ?? "active"} /> },
  ];

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />

      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <h1>{t("title")}</h1>
            <p style={{ color: "var(--text-secondary)" }}>TLS certificates and profiles.</p>
          </div>
          <div className="header-actions" style={{ display: "flex", alignItems: "center", gap: "12px" }}>
            <UserProfileWidget />
          </div>
        </div>

        <div className="glass-panel animate-fade-in delay-1" style={{ padding: "16px 20px", marginBottom: "20px", display: "flex", justifyContent: "space-between", alignItems: "center", gap: "12px" }}>
          <div>
            <h3 style={{ fontSize: "15px", fontWeight: 600 }}>Let&apos;s Encrypt Auto-Provisioning</h3>
            <p style={{ fontSize: "13px", color: "var(--text-secondary)" }}>Automatically provision a free TLS certificate via ACME. Requires the domain to be publicly reachable on port 80.</p>
          </div>
          <button type="button" className="btn btn-primary" onClick={() => setAcmeOpen(true)} style={{ padding: "10px 18px", whiteSpace: "nowrap" }}>
            + Provision Certificate
          </button>
        </div>

        <div className="data-section animate-fade-in delay-1">
          <DataTable
            columns={columns}
            data={rows}
            rowKey={(row) => row.id}
            loading={loading}
            error={error || undefined}
            onRetry={load}
            emptyMessage="No certificates yet."
            actions={(row) => (
              <button type="button" className="btn" onClick={() => setDeleting(row)} style={{ padding: "6px 12px", fontSize: "12px", color: "var(--danger)" }}>
                Delete
              </button>
            )}
          />
        </div>
      </main>

      <ConfirmDialog
        open={!!deleting}
        title="Delete certificate"
        message={`Are you sure you want to delete "${deleting?.name}" (${deleting?.domain})?`}
        confirmLabel="Delete"
        danger
        loading={submitting}
        onConfirm={confirmDelete}
        onCancel={() => setDeleting(null)}
      />

      {acmeOpen && (
        <div onClick={() => setAcmeOpen(false)} style={{ position: "fixed", inset: 0, zIndex: 1200, background: "rgba(0,0,0,0.7)", display: "flex", alignItems: "center", justifyContent: "center", padding: "20px" }}>
          <div className="glass-panel animate-fade-in" onClick={(e) => e.stopPropagation()} style={{ width: "100%", maxWidth: "440px", padding: "28px", display: "flex", flexDirection: "column", gap: "14px" }}>
            <h3 style={{ fontSize: "17px", fontWeight: 700 }}>Provision Let&apos;s Encrypt Certificate</h3>
            <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Domain</label>
              <input style={inputStyle} value={acmeForm.domain} onChange={(e) => setAcmeForm({ ...acmeForm, domain: e.target.value })} placeholder="example.com" />
            </div>
            <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Account Email</label>
              <input style={inputStyle} type="email" value={acmeForm.email} onChange={(e) => setAcmeForm({ ...acmeForm, email: e.target.value })} placeholder="admin@example.com" />
            </div>
            <div style={{ display: "flex", alignItems: "center", gap: "8px" }}>
              <input type="checkbox" checked={acmeForm.staging} onChange={(e) => setAcmeForm({ ...acmeForm, staging: e.target.checked })} />
              <span style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Use staging (test) environment</span>
            </div>
            <div style={{ display: "flex", justifyContent: "flex-end", gap: "10px" }}>
              <button className="btn" onClick={() => setAcmeOpen(false)} disabled={submitting} style={{ padding: "10px 16px" }}>Cancel</button>
              <button className="btn btn-primary" onClick={handleProvision} disabled={submitting || !acmeForm.domain || !acmeForm.email} style={{ padding: "10px 16px" }}>
                {submitting ? "Provisioning…" : "Provision"}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

const inputStyle: CSSProperties = {
  background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)",
  padding: "10px 12px", borderRadius: "8px", color: "white", outline: "none", fontSize: "14px", width: "100%",
};
