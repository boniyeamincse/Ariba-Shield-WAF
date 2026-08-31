"use client";

import { useEffect, useState } from "react";
import { useLocale, useTranslations } from "next-intl";
import Link from "next/link";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import DataTable, { type Column } from "@/components/shared/DataTable";
import { StatusBadge } from "@/components/shared/Badges";
import ConfirmDialog from "@/components/shared/ConfirmDialog";
import { usePermission } from "@/hooks/usePermission";
import ApplicationWizard from "@/components/applications/ApplicationWizard";
import {
  listApplications,
  deleteApplication,
  type Application,
} from "@/lib/api";

type AppRow = Application & { origins: number; domains: number };

export default function ApplicationsPage() {
  const locale = useLocale();
  const t = useTranslations("applications");
  const ct = useTranslations("common");
  const [rows, setRows] = useState<AppRow[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [wizardOpen, setWizardOpen] = useState(false);
  const [editing, setEditing] = useState<Application | null>(null);
  const [deleting, setDeleting] = useState<Application | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const { can: canUser } = usePermission();

  const load = async () => {
    setLoading(true);
    setError("");
    try {
      const apps = await listApplications();
      setRows(apps as AppRow[]);
    } catch {
      setError("Failed to load applications");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const openCreate = () => {
    setEditing(null);
    setWizardOpen(true);
  };

  const openEdit = (app: Application) => {
    setEditing(app);
    setWizardOpen(true);
  };

  const closeWizard = async () => {
    setWizardOpen(false);
    setEditing(null);
    await load();
  };

  const confirmDelete = async () => {
    if (!deleting) return;
    setSubmitting(true);
    try {
      await deleteApplication(deleting.id);
      setDeleting(null);
      await load();
    } catch {
      setError("Failed to delete application");
    } finally {
      setSubmitting(false);
    }
  };

  const columns: Column<AppRow>[] = [
    {
      key: "name",
      label: "Name",
      sortable: true,
      render: (row) => (
        <Link href={`/${locale}/applications/${row.id}`} style={{ color: "var(--accent-primary)", textDecoration: "none", fontWeight: 500 }}>
          {row.name}
        </Link>
      ),
    },
    {
      key: "environment",
      label: "Environment",
      render: (row) => <StatusBadge value={row.environment ?? "production"} />,
    },
    { key: "domain", label: "Domain", render: (row) => row.domain || "—" },
    {
      key: "origin_host",
      label: "Origin",
      render: (row) => (row.origin_host ? `${row.origin_host}:${row.origin_port ?? ""}` : "—"),
    },
    { key: "status", label: "Status", render: (row) => <StatusBadge value={row.status} /> },
    { key: "version", label: "Version", render: (row) => String(row.version ?? "—") },
  ];

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />

      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <h1>{t("title")}</h1>
            <p style={{ color: "var(--text-secondary)" }}>{t("description")}</p>
          </div>
          <div className="header-actions" style={{ display: "flex", alignItems: "center", gap: "12px" }}>
            {canUser("create") && (
              <button type="button" className="btn btn-primary" onClick={openCreate}>
                + {t("new_application")}
              </button>
            )}
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
            emptyMessage="No applications yet."
            actions={(row) => (
              <div style={{ display: "flex", gap: "8px" }}>
                {canUser("edit") && (
                  <button type="button" className="btn" onClick={() => openEdit(row)} style={{ padding: "6px 12px", fontSize: "12px" }}>
                    {ct("edit")}
                  </button>
                )}
                {canUser("delete") && (
                  <button type="button" className="btn" onClick={() => setDeleting(row)} style={{ padding: "6px 12px", fontSize: "12px", color: "var(--danger)" }}>
                    {ct("delete")}
                  </button>
                )}
              </div>
            )}
          />
        </div>
      </main>

      {wizardOpen && (
        <ApplicationWizard
          initial={editing}
          onClose={() => {
            setWizardOpen(false);
            setEditing(null);
          }}
          onCreated={closeWizard}
        />
      )}

      <ConfirmDialog
        open={!!deleting}
        title="Delete application"
        message={`Are you sure you want to delete "${deleting?.name}"? This action cannot be undone.`}
        confirmLabel="Delete"
        danger
        loading={submitting}
        onConfirm={confirmDelete}
        onCancel={() => setDeleting(null)}
      />
    </div>
  );
}
