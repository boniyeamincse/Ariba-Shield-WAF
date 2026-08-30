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
import {
  listApplications,
  createApplication,
  updateApplication,
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
  const [modal, setModal] = useState<"create" | "edit" | null>(null);
  const [editing, setEditing] = useState<Application | null>(null);
  const [deleting, setDeleting] = useState<Application | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [form, setForm] = useState({ name: "", description: "" });
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
    setForm({ name: "", description: "" });
    setModal("create");
  };

  const openEdit = (app: Application) => {
    setForm({ name: app.name, description: app.description ?? "" });
    setEditing(app);
    setModal("edit");
  };

  const submit = async () => {
    setSubmitting(true);
    setError("");
    try {
      if (modal === "create") {
        await createApplication(form.name, form.description);
      } else if (editing) {
        await updateApplication(editing.id, { name: form.name, description: form.description });
      }
      setModal(null);
      setEditing(null);
      await load();
    } catch {
      setError("Failed to save application");
    } finally {
      setSubmitting(false);
    }
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
    { key: "status", label: "Status", render: (row) => <StatusBadge value={row.status} /> },
    { key: "origins", label: "Origins", render: (row) => String(row.origins ?? "—") },
    { key: "domains", label: "Domains", render: (row) => String(row.domains ?? "—") },
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

      {modal && (
        <div
          onClick={() => setModal(null)}
          style={{
            position: "fixed", inset: 0, zIndex: 1000, background: "rgba(0,0,0,0.6)",
            display: "flex", alignItems: "center", justifyContent: "center", padding: "20px",
          }}
        >
          <div
            className="glass-panel animate-fade-in"
            onClick={(e) => e.stopPropagation()}
            style={{ width: "100%", maxWidth: "440px", padding: "28px", display: "flex", flexDirection: "column", gap: "16px" }}
          >
            <h3 style={{ fontSize: "17px", fontWeight: 600 }}>
              {modal === "create" ? "New Application" : "Edit Application"}
            </h3>
            <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Name</label>
              <input
                type="text"
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
                style={{
                  background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)",
                  padding: "10px 12px", borderRadius: "8px", color: "white", outline: "none", fontSize: "14px",
                }}
              />
            </div>
            <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
              <label style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Description</label>
              <textarea
                value={form.description}
                onChange={(e) => setForm({ ...form, description: e.target.value })}
                rows={3}
                style={{
                  background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)",
                  padding: "10px 12px", borderRadius: "8px", color: "white", outline: "none", fontSize: "14px",
                  resize: "vertical",
                }}
              />
            </div>
            <div style={{ display: "flex", justifyContent: "flex-end", gap: "10px", marginTop: "8px" }}>
              <button type="button" className="btn" onClick={() => setModal(null)} disabled={submitting} style={{ padding: "10px 16px" }}>
                {ct("cancel")}
              </button>
              <button type="button" className="btn btn-primary" onClick={submit} disabled={submitting || !form.name} style={{ padding: "10px 16px" }}>
                {submitting ? "Saving…" : "Save"}
              </button>
            </div>
          </div>
        </div>
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
