"use client";

import { useEffect, useState } from "react";
import { useLocale, useTranslations } from "next-intl";
import Link from "next/link";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import DataTable, { type Column } from "@/components/shared/DataTable";
import { StatusBadge } from "@/components/shared/Badges";
import ConfirmDialog from "@/components/shared/ConfirmDialog";
import {
  listSecurityPolicies,
  createSecurityPolicy,
  deleteSecurityPolicy,
  activateSecurityPolicy,
  disableSecurityPolicy,
  type SecurityPolicy,
} from "@/lib/api";

export default function PoliciesPage() {
  const locale = useLocale();
  const t = useTranslations("policies");
  const ct = useTranslations("common");
  const [rows, setRows] = useState<SecurityPolicy[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [deleting, setDeleting] = useState<SecurityPolicy | null>(null);

  const load = async () => {
    setLoading(true);
    setError("");
    try {
      setRows(await listSecurityPolicies());
    } catch {
      setError("Failed to load policies");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const doActivate = async (id: string) => {
    try {
      await activateSecurityPolicy(id);
      await load();
    } catch {
      setError("Failed to activate policy");
    }
  };

  const doDisable = async (id: string) => {
    try {
      await disableSecurityPolicy(id);
      await load();
    } catch {
      setError("Failed to disable policy");
    }
  };

  const confirmDelete = async () => {
    if (!deleting) return;
    setSubmitting(true);
    try {
      await deleteSecurityPolicy(deleting.id);
      setDeleting(null);
      await load();
    } catch {
      setError("Failed to delete policy");
    } finally {
      setSubmitting(false);
    }
  };

  const columns: Column<SecurityPolicy>[] = [
    {
      key: "name",
      label: "Name",
      sortable: true,
      render: (row) => (
        <Link href={`/${locale}/policies/${row.id}`} style={{ color: "var(--accent-primary)", textDecoration: "none", fontWeight: 500 }}>
          {row.name}
        </Link>
      ),
    },
    { key: "enforcement_mode", label: "Mode", render: (row) => <StatusBadge value={row.enforcement_mode} /> },
    { key: "application_id", label: "Application", render: (row) => (row.application_id ? row.application_id.slice(0, 8) + "…" : "—") },
    { key: "version", label: "Version", render: (row) => `v${row.version}` },
    { key: "description", label: "Description", render: (row) => row.description || "—" },
  ];

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />

      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <h1>{t("title")}</h1>
            <p style={{ color: "var(--text-secondary)" }}>Manage WAF security policies.</p>
          </div>
          <div className="header-actions" style={{ display: "flex", alignItems: "center", gap: "12px" }}>
            <Link href={`/${locale}/policies/create`} className="btn btn-primary" style={{ textDecoration: 'none' }}>
              + New Policy
            </Link>
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
            emptyMessage="No security policies yet."
            actions={(row) => (
              <div style={{ display: "flex", gap: "8px" }}>
                <button type="button" className="btn" onClick={() => doActivate(row.id)} style={{ padding: "6px 10px", fontSize: "12px", color: "var(--success)" }}>
                  {ct("activate")}
                </button>
                <button type="button" className="btn" onClick={() => doDisable(row.id)} style={{ padding: "6px 10px", fontSize: "12px", color: "var(--warning)" }}>
                  {ct("disable")}
                </button>
                <button type="button" className="btn" onClick={() => setDeleting(row)} style={{ padding: "6px 10px", fontSize: "12px", color: "var(--danger)" }}>
                  {ct("delete")}
                </button>
              </div>
            )}
          />
        </div>
      </main>

      <ConfirmDialog
        open={!!deleting}
        title="Delete policy"
        message={`Are you sure you want to delete "${deleting?.name}"? This cannot be undone.`}
        confirmLabel="Delete"
        danger
        loading={submitting}
        onConfirm={confirmDelete}
        onCancel={() => setDeleting(null)}
      />
    </div>
  );
}
