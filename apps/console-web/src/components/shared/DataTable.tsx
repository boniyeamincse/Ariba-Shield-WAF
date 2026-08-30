"use client";

import { useState } from "react";
import type { ReactNode } from "react";

export type Column<T> = {
  key: string;
  label: string;
  sortable?: boolean;
  render?: (row: T) => ReactNode;
};

type DataTableProps<T> = {
  columns: Column<T>[];
  data: T[];
  rowKey: (row: T) => string;
  loading?: boolean;
  emptyMessage?: string;
  error?: string;
  onRetry?: () => void;
  onSort?: (key: string, direction: "asc" | "desc") => void;
  defaultSort?: { key: string; direction: "asc" | "desc" };
  pageSize?: number;
  page?: number;
  total?: number;
  onPageChange?: (page: number) => void;
  onRowClick?: (row: T) => void;
  actions?: (row: T) => ReactNode;
};

export default function DataTable<T>({
  columns,
  data,
  rowKey,
  loading = false,
  emptyMessage = "No data",
  error,
  onRetry,
  onSort,
  defaultSort,
  pageSize,
  page = 1,
  total,
  onPageChange,
  onRowClick,
  actions,
}: DataTableProps<T>) {
  const [localSort, setLocalSort] = useState<{ key: string; direction: "asc" | "desc" } | undefined>(defaultSort);
  const sort = onSort ? undefined : localSort;

  const sorted = (() => {
    if (sort && sort.key) {
      const dir = sort.direction === "asc" ? 1 : -1;
      return [...data].sort((a, b) => {
        const av = String((a as Record<string, unknown>)[sort.key] ?? "");
        const bv = String((b as Record<string, unknown>)[sort.key] ?? "");
        return av.localeCompare(bv) * dir;
      });
    }
    return data;
  })();

  const totalCount = total ?? data.length;
  const totalPages = pageSize ? Math.max(1, Math.ceil(totalCount / pageSize)) : 1;

  const handleSort = (col: Column<T>) => {
    if (!col.sortable) return;
    const nextDir = sort?.key === col.key && sort.direction === "asc" ? "desc" : "asc";
    if (onSort) {
      onSort(col.key, nextDir);
    } else {
      setLocalSort({ key: col.key, direction: nextDir });
    }
  };

  if (loading) {
    return (
      <div className="glass-panel" style={{ padding: "40px", textAlign: "center", color: "var(--text-secondary)", fontSize: "14px" }}>
        Loading…
      </div>
    );
  }

  if (error) {
    return (
      <div className="glass-panel" style={{ padding: "40px", textAlign: "center", display: "flex", flexDirection: "column", alignItems: "center", gap: "12px" }}>
        <p style={{ color: "var(--danger)", fontSize: "14px" }}>{error}</p>
        {onRetry && (
          <button type="button" className="btn btn-primary" onClick={onRetry}>
            Retry
          </button>
        )}
      </div>
    );
  }

  return (
    <div className="glass-panel" style={{ padding: 0, overflow: "hidden" }}>
      <div style={{ overflowX: "auto" }}>
        <table style={{ width: "100%", borderCollapse: "collapse", textAlign: "left" }}>
          <thead>
            <tr style={{ borderBottom: "1px solid rgba(255,255,255,0.08)", background: "rgba(255,255,255,0.02)" }}>
              {columns.map((col) => (
                <th
                  key={col.key}
                  onClick={() => handleSort(col)}
                  style={{
                    padding: "14px 20px",
                    color: "var(--text-secondary)",
                    fontWeight: 500,
                    fontSize: "13px",
                    whiteSpace: "nowrap",
                    cursor: col.sortable ? "pointer" : "default",
                    userSelect: "none",
                  }}
                >
                  {col.label}
                  {col.sortable && sort?.key === col.key && (
                    <span style={{ marginLeft: 6, fontSize: 10 }}>{sort.direction === "asc" ? "▲" : "▼"}</span>
                  )}
                </th>
              ))}
              {actions && <th style={{ padding: "14px 20px", color: "var(--text-secondary)", fontWeight: 500, fontSize: "13px", whiteSpace: "nowrap" }}>Actions</th>}
            </tr>
          </thead>
          <tbody>
            {sorted.map((row, idx) => (
              <tr
                key={rowKey(row)}
                onClick={() => onRowClick?.(row)}
                style={{
                  borderBottom: idx !== sorted.length - 1 ? "1px solid rgba(255,255,255,0.05)" : "none",
                  cursor: onRowClick ? "pointer" : "default",
                }}
              >
                {columns.map((col) => (
                  <td key={col.key} style={{ padding: "12px 20px", fontSize: "13px" }}>
                    {col.render ? col.render(row) : String((row as Record<string, unknown>)[col.key] ?? "—")}
                  </td>
                ))}
                {actions && <td style={{ padding: "12px 20px", whiteSpace: "nowrap" }}>{actions(row)}</td>}
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {sorted.length === 0 && (
        <div style={{ padding: "60px", textAlign: "center", color: "var(--text-secondary)" }}>
          <p>{emptyMessage}</p>
        </div>
      )}

      {pageSize && totalPages > 1 && (
        <div
          style={{
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            padding: "14px 20px",
            borderTop: "1px solid rgba(255,255,255,0.08)",
            fontSize: "13px",
            color: "var(--text-secondary)",
          }}
        >
          <span>
            Page {page} of {totalPages}
          </span>
          <div style={{ display: "flex", gap: "8px" }}>
            <button
              type="button"
              className="btn"
              disabled={page <= 1}
              onClick={() => onPageChange?.(page - 1)}
              style={{ padding: "6px 12px", fontSize: "12px" }}
            >
              Prev
            </button>
            <button
              type="button"
              className="btn"
              disabled={page >= totalPages}
              onClick={() => onPageChange?.(page + 1)}
              style={{ padding: "6px 12px", fontSize: "12px" }}
            >
              Next
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
