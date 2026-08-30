"use client";

import { useCallback, useEffect, useState } from "react";
import { useLocale } from "next-intl";
import { useParams, useRouter } from "next/navigation";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import { StatusBadge } from "@/components/shared/Badges";
import {
  getSecurityPolicy,
  listPolicyVersions,
  rollbackSecurityPolicy,
  cloneSecurityPolicy,
  type SecurityPolicy,
} from "@/lib/api";

type PolicyVersion = {
  id: string;
  version: number;
  status: string;
  created_at?: string;
};

export default function PolicyDetailPage() {
  const locale = useLocale();
  const router = useRouter();
  const params = useParams<{ id: string }>();
  const [policy, setPolicy] = useState<SecurityPolicy | null>(null);
  const [versions, setVersions] = useState<PolicyVersion[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const load = useCallback(async () => {
    if (!params?.id) return;
    setLoading(true);
    setError("");
    try {
      const [p, v] = await Promise.all([
        getSecurityPolicy(params.id),
        listPolicyVersions(params.id).catch(() => []),
      ]);
      setPolicy(p);
      setVersions((v as PolicyVersion[]) ?? []);
    } catch {
      setError("Failed to load policy");
    } finally {
      setLoading(false);
    }
  }, [params?.id]);

  useEffect(() => {
    load();
  }, [load]);

  const doRollback = async () => {
    if (!policy) return;
    try {
      await rollbackSecurityPolicy(policy.id);
      await load();
    } catch {
      setError("Failed to rollback policy");
    }
  };

  const doClone = async () => {
    if (!policy) return;
    try {
      const { id } = await cloneSecurityPolicy(policy.id);
      router.push(`/${locale}/policies/${id}`);
    } catch {
      setError("Failed to clone policy");
    }
  };

  if (loading) {
    return (
      <div className="dashboard-container">
        <Sidebar locale={locale} />
        <main className="main-content"><p style={{ color: "var(--text-secondary)" }}>Loading…</p></main>
      </div>
    );
  }

  if (error || !policy) {
    return (
      <div className="dashboard-container">
        <Sidebar locale={locale} />
        <main className="main-content">
          <div className="glass-panel" style={{ padding: "40px", textAlign: "center" }}>
            <p style={{ color: "var(--danger)" }}>{error || "Policy not found"}</p>
            <button type="button" className="btn btn-primary" onClick={load} style={{ marginTop: "12px" }}>Retry</button>
          </div>
        </main>
      </div>
    );
  }

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />

      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <button type="button" onClick={() => router.push(`/${locale}/policies`)} style={{ background: "none", border: "none", color: "var(--text-secondary)", cursor: "pointer", padding: 0, marginBottom: "4px", fontSize: "13px" }}>
              ← Back to policies
            </button>
            <h1>{policy.name}</h1>
            <p style={{ color: "var(--text-secondary)" }}>{policy.description}</p>
          </div>
          <div className="header-actions" style={{ display: "flex", alignItems: "center", gap: "12px" }}>
            <button type="button" className="btn" onClick={doClone} style={{ padding: "8px 14px", fontSize: "13px" }}>Clone</button>
            <button type="button" className="btn" onClick={doRollback} style={{ padding: "8px 14px", fontSize: "13px", color: "var(--warning)" }}>Rollback</button>
            <UserProfileWidget />
          </div>
        </div>

        <div className="metrics-grid animate-fade-in delay-1" style={{ maxWidth: "800px", marginBottom: "28px" }}>
          <div className="glass-panel metric-card" style={{ padding: "20px" }}>
            <div className="metric-header"><span>Mode</span></div>
            <div style={{ marginTop: "8px" }}><StatusBadge value={policy.enforcement_mode} /></div>
          </div>
          <div className="glass-panel metric-card" style={{ padding: "20px" }}>
            <div className="metric-header"><span>Version</span></div>
            <div className="metric-value" style={{ fontSize: "24px" }}>v{policy.version}</div>
          </div>
          <div className="glass-panel metric-card" style={{ padding: "20px" }}>
            <div className="metric-header"><span>Application</span></div>
            <div style={{ marginTop: "8px", fontSize: "13px", color: "var(--text-secondary)" }}>
              {policy.application_id ? policy.application_id.slice(0, 12) + "…" : "Not bound"}
            </div>
          </div>
        </div>

        <div className="data-section animate-fade-in delay-2">
          <div className="section-header">
            <h2>Version History</h2>
          </div>
          <div className="glass-panel" style={{ padding: "0", overflow: "hidden" }}>
            <table style={{ width: "100%", borderCollapse: "collapse", textAlign: "left" }}>
              <thead>
                <tr style={{ borderBottom: "1px solid rgba(255,255,255,0.08)", background: "rgba(255,255,255,0.02)" }}>
                  <th style={{ padding: "14px 20px", color: "var(--text-secondary)", fontWeight: 500, fontSize: "13px" }}>Version</th>
                  <th style={{ padding: "14px 20px", color: "var(--text-secondary)", fontWeight: 500, fontSize: "13px" }}>Status</th>
                  <th style={{ padding: "14px 20px", color: "var(--text-secondary)", fontWeight: 500, fontSize: "13px" }}>Created</th>
                </tr>
              </thead>
              <tbody>
                {versions.map((v, idx) => (
                  <tr key={v.id} style={{ borderBottom: idx !== versions.length - 1 ? "1px solid rgba(255,255,255,0.05)" : "none" }}>
                    <td style={{ padding: "12px 20px", fontSize: "13px" }}>v{v.version}</td>
                    <td style={{ padding: "12px 20px", fontSize: "13px" }}><StatusBadge value={v.status} /></td>
                    <td style={{ padding: "12px 20px", fontSize: "13px", color: "var(--text-secondary)" }}>
                      {v.created_at ? new Date(v.created_at).toLocaleString() : "—"}
                    </td>
                  </tr>
                ))}
                {versions.length === 0 && (
                  <tr><td colSpan={3} style={{ padding: "40px", textAlign: "center", color: "var(--text-secondary)" }}>No version history available.</td></tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
      </main>
    </div>
  );
}
