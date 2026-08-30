"use client";

import { useCallback, useEffect, useState } from "react";
import { useLocale, useTranslations } from "next-intl";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import DataTable, { type Column } from "@/components/shared/DataTable";
import { SeverityBadge } from "@/components/shared/Badges";
import FilterBar from "@/components/shared/FilterBar";
import { listSecurityEvents, type SecurityEvent } from "@/lib/api";

const PAGE_SIZE = 25;

export default function SecurityEventsPage() {
  const locale = useLocale();
  const t = useTranslations("security_events");
  const [rows, setRows] = useState<SecurityEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [filters, setFilters] = useState<Record<string, string>>({});

  const load = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const params: Record<string, string> = {
        limit: String(PAGE_SIZE),
        offset: String((page - 1) * PAGE_SIZE),
        ...filters,
      };
      const data = await listSecurityEvents(params);
      setRows(data.events ?? []);
      setTotal(data.pagination?.count ?? 0);
    } catch {
      setError("Failed to load security events");
    } finally {
      setLoading(false);
    }
  }, [page, filters]);

  useEffect(() => {
    load();
  }, [load]);

  const columns: Column<SecurityEvent>[] = [
    { key: "severity", label: "Severity", render: (row) => <SeverityBadge value={row.severity} /> },
    { key: "created_at", label: "Time", render: (row) => new Date(row.created_at).toLocaleString() },
    { key: "method", label: "Method", render: (row) => row.method || "—" },
    { key: "path", label: "Path", render: (row) => <span style={{ fontFamily: "monospace", fontSize: "12px" }}>{row.path || "—"}</span> },
    { key: "reason", label: "Reason", render: (row) => row.reason || "—" },
    { key: "rule_ids", label: "Rules", render: (row) => (row.rule_ids?.length ? row.rule_ids.slice(0, 3).join(", ") : "—") },
    { key: "client_ip", label: "Client IP", render: (row) => <span style={{ fontFamily: "monospace", fontSize: "12px" }}>{row.client_ip || "—"}</span> },
  ];

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />

      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <h1>{t("title")}</h1>
            <p style={{ color: "var(--text-secondary)" }}>Inspected and blocked requests.</p>
          </div>
          <div className="header-actions" style={{ display: "flex", alignItems: "center", gap: "12px" }}>
            <UserProfileWidget />
          </div>
        </div>

        <FilterBar
          filters={[
            { type: "select", key: "severity", label: "Severity", options: [
              { value: "critical", label: "Critical" },
              { value: "high", label: "High" },
              { value: "medium", label: "Medium" },
              { value: "low", label: "Low" },
            ]},
            { type: "text", key: "client_ip", label: "Client IP", placeholder: "Filter by IP" },
            { type: "select", key: "action", label: "Action", options: [
              { value: "block", label: "Block" },
              { value: "pass", label: "Pass" },
              { value: "detect", label: "Detect" },
            ]},
          ]}
          values={filters}
          onChange={(v) => {
            setFilters(v);
            setPage(1);
          }}
        />

        <div className="data-section animate-fade-in delay-1">
          <DataTable
            columns={columns}
            data={rows}
            rowKey={(row) => row.id}
            loading={loading}
            error={error || undefined}
            onRetry={load}
            emptyMessage="No security events yet."
            pageSize={PAGE_SIZE}
            page={page}
            total={total}
            onPageChange={setPage}
          />
        </div>
      </main>
    </div>
  );
}
