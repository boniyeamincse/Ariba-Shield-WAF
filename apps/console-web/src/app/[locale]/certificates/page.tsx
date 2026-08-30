"use client";

import { useCallback, useEffect, useState } from "react";
import { useLocale, useTranslations } from "next-intl";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import DataTable, { type Column } from "@/components/shared/DataTable";
import { StatusBadge } from "@/components/shared/Badges";
import ConfirmDialog from "@/components/shared/ConfirmDialog";
import {
  listCertificates,
  deleteCertificate,
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
    </div>
  );
}
