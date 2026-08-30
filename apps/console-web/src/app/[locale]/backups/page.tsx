"use client";

import { useCallback, useEffect, useState } from "react";
import { useLocale } from "next-intl";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import DataTable, { type Column } from "@/components/shared/DataTable";
import { StatusBadge } from "@/components/shared/Badges";
import ConfirmDialog from "@/components/shared/ConfirmDialog";
import {
  listBackups,
  createBackup,
  restoreBackup,
  type Backup,
} from "@/lib/api";

export default function BackupsPage() {
  const locale = useLocale();
  const [rows, setRows] = useState<Backup[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [creating, setCreating] = useState(false);
  const [restoring, setRestoring] = useState<Backup | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      setRows(await listBackups());
    } catch {
      setError("Failed to load backups");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const doCreate = async () => {
    setCreating(true);
    setError("");
    try {
      await createBackup();
      await load();
    } catch {
      setError("Failed to create backup");
    } finally {
      setCreating(false);
    }
  };

  const confirmRestore = async () => {
    if (!restoring) return;
    setSubmitting(true);
    try {
      await restoreBackup(restoring.id);
      setRestoring(null);
      await load();
    } catch {
      setError("Failed to restore backup");
    } finally {
      setSubmitting(false);
    }
  };

  const columns: Column<Backup>[] = [
    { key: "id", label: "ID", render: (row) => <span style={{ fontFamily: "monospace", fontSize: "12px" }}>{row.id.slice(0, 12)}…</span> },
    { key: "status", label: "Status", render: (row) => <StatusBadge value={row.status} /> },
    { key: "artifact_ref", label: "Artifact", render: (row) => row.artifact_ref || "—" },
    {
      key: "size_bytes",
      label: "Size",
      render: (row) => (row.size_bytes ? `${(row.size_bytes / 1024 / 1024).toFixed(1)} MB` : "—"),
    },
    { key: "created_at", label: "Created", render: (row) => new Date(row.created_at).toLocaleString() },
  ];

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />

      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <h1>Backups</h1>
            <p style={{ color: "var(--text-secondary)" }}>Backup and disaster recovery.</p>
          </div>
          <div className="header-actions" style={{ display: "flex", alignItems: "center", gap: "12px" }}>
            <button type="button" className="btn btn-primary" onClick={doCreate} disabled={creating} style={{ padding: "8px 14px", fontSize: "13px" }}>
              {creating ? "Creating…" : "+ New Backup"}
            </button>
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
            emptyMessage="No backups yet."
            actions={(row) => (
              <button type="button" className="btn" onClick={() => setRestoring(row)} disabled={row.status === "in_progress"} style={{ padding: "6px 12px", fontSize: "12px", color: "var(--warning)" }}>
                Restore
              </button>
            )}
          />
        </div>
      </main>

      <ConfirmDialog
        open={!!restoring}
        title="Restore backup"
        message={`Are you sure you want to restore backup "${restoring?.id.slice(0, 12)}…"? This will overwrite current configuration.`}
        confirmLabel="Restore"
        danger
        loading={submitting}
        onConfirm={confirmRestore}
        onCancel={() => setRestoring(null)}
      />
    </div>
  );
}
