"use client";

import { useCallback, useEffect, useState } from "react";
import { useLocale, useTranslations } from "next-intl";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import { StatusBadge } from "@/components/shared/Badges";
import ConfirmDialog from "@/components/shared/ConfirmDialog";
import {
  getLicense,
  activateLicense,
  deactivateLicense,
  getLicenseUsage,
  getLicenseEntitlements,
  type License,
} from "@/lib/api";

type Usage = {
  gateways?: { current: number; limit: number; percent: number };
  applications?: { current: number; limit: number; percent: number };
};

export default function LicensePage() {
  const locale = useLocale();
  const t = useTranslations("license_page");
  const [license, setLicense] = useState<License | null>(null);
  const [usage, setUsage] = useState<Usage>({});
  const [entitlements, setEntitlements] = useState<Record<string, boolean>>({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [key, setKey] = useState("");
  const [activating, setActivating] = useState(false);
  const [confirmDeactivate, setConfirmDeactivate] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const [lic, usg, ent] = await Promise.all([
        getLicense(),
        getLicenseUsage(),
        getLicenseEntitlements(),
      ]);
      setLicense(lic as License);
      setUsage(usg as Usage);
      setEntitlements((ent as { entitlements: Record<string, boolean> }).entitlements ?? {});
    } catch {
      setError("Failed to load license");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const doActivate = async () => {
    setActivating(true);
    setError("");
    try {
      await activateLicense(key);
      setKey("");
      await load();
    } catch {
      setError("Failed to activate license");
    } finally {
      setActivating(false);
    }
  };

  const doDeactivate = async () => {
    setSubmitting(true);
    try {
      await deactivateLicense();
      setConfirmDeactivate(false);
      await load();
    } catch {
      setError("Failed to deactivate license");
    } finally {
      setSubmitting(false);
    }
  };

  const edition = license?.edition ?? "community";
  const status = license?.status ?? "active";

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />

      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <h1>{t("title")}</h1>
            <p style={{ color: "var(--text-secondary)" }}>License status, usage, and entitlements.</p>
          </div>
          <div className="header-actions" style={{ display: "flex", alignItems: "center", gap: "12px" }}>
            {license && <StatusBadge value={status} />}
            <UserProfileWidget />
          </div>
        </div>

        {error && (
          <div
            className="glass-panel"
            style={{ padding: "12px 16px", marginBottom: "16px", color: "var(--danger)", fontSize: "13px", border: "1px solid rgba(239,68,68,0.3)" }}
          >
            {error}
          </div>
        )}

        {loading ? (
          <p style={{ color: "var(--text-secondary)" }}>Loading license…</p>
        ) : (
          <>
            {/* License info */}
            <div className="metrics-grid animate-fade-in delay-1" style={{ marginBottom: "24px" }}>
              <div className="glass-panel metric-card" style={{ padding: "20px" }}>
                <div className="metric-header"><span>Edition</span></div>
                <div className="metric-value" style={{ fontSize: "24px", textTransform: "capitalize" }}>{edition}</div>
              </div>
              <div className="glass-panel metric-card" style={{ padding: "20px" }}>
                <div className="metric-header"><span>Seats</span></div>
                <div className="metric-value" style={{ fontSize: "24px" }}>{license?.seats ?? 1}</div>
              </div>
              <div className="glass-panel metric-card" style={{ padding: "20px" }}>
                <div className="metric-header"><span>Expires</span></div>
                <div style={{ marginTop: "8px", fontSize: "14px", color: "var(--text-secondary)" }}>
                  {license?.expires_at ? new Date(license.expires_at).toLocaleDateString() : "—"}
                </div>
              </div>
            </div>

            {/* Usage bars */}
            <div className="data-section animate-fade-in delay-2">
              <div className="glass-panel" style={{ padding: "24px" }}>
                <h3 style={{ fontSize: "16px", fontWeight: 600, marginBottom: "16px" }}>Usage</h3>
                {(["gateways", "applications"] as const).map((resource) => {
                  const u = usage[resource];
                  const label = resource.charAt(0).toUpperCase() + resource.slice(1);
                  const pct = Math.min(100, u?.percent ?? 0);
                  const color = pct >= 90 ? "var(--danger)" : pct >= 70 ? "var(--warning)" : "var(--success)";
                  return (
                    <div key={resource} style={{ marginBottom: "16px" }}>
                      <div style={{ display: "flex", justifyContent: "space-between", fontSize: "13px", marginBottom: "6px" }}>
                        <span style={{ color: "var(--text-secondary)" }}>{label}</span>
                        <span>
                          {u?.current ?? 0} / {u?.limit ?? (resource === "gateways" ? 1 : 10)}
                        </span>
                      </div>
                      <div style={{ height: "8px", borderRadius: "4px", background: "rgba(255,255,255,0.06)", overflow: "hidden" }}>
                        <div style={{ height: "100%", width: `${pct}%`, borderRadius: "4px", background: color, transition: "width 0.3s" }} />
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>

            {/* Entitlements */}
            <div className="data-section animate-fade-in delay-2">
              <div className="glass-panel" style={{ padding: "24px" }}>
                <h3 style={{ fontSize: "16px", fontWeight: 600, marginBottom: "16px" }}>Entitlements</h3>
                <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(200px, 1fr))", gap: "10px" }}>
                  {Object.entries(entitlements).map(([feature, enabled]) => (
                    <div key={feature} style={{ display: "flex", alignItems: "center", gap: "8px", fontSize: "13px" }}>
                      <span style={{ color: enabled ? "var(--success)" : "var(--text-secondary)", fontSize: "15px" }}>
                        {enabled ? "✓" : "✕"}
                      </span>
                      <span style={{ color: enabled ? "var(--text-primary)" : "var(--text-secondary)", textTransform: "capitalize" }}>
                        {feature.replace(/_/g, " ")}
                      </span>
                    </div>
                  ))}
                </div>
              </div>
            </div>

            {/* Activate / Deactivate */}
            <div className="data-section animate-fade-in delay-2">
              <div className="glass-panel" style={{ padding: "24px" }}>
                <h3 style={{ fontSize: "16px", fontWeight: 600, marginBottom: "16px" }}>Manage License</h3>
                {license && license.status === "active" ? (
                  <button type="button" className="btn" onClick={() => setConfirmDeactivate(true)} style={{ padding: "10px 16px", color: "var(--danger)", borderColor: "var(--danger)" }}>
                    Deactivate License
                  </button>
                ) : (
                  <div style={{ display: "flex", gap: "10px", flexWrap: "wrap" }}>
                    <input
                      type="text"
                      placeholder="Enter license key"
                      value={key}
                      onChange={(e) => setKey(e.target.value)}
                      style={{
                        flex: "1", minWidth: "200px", background: "rgba(255,255,255,0.05)",
                        border: "1px solid rgba(255,255,255,0.1)", padding: "10px 12px",
                        borderRadius: "8px", color: "white", outline: "none", fontSize: "14px",
                      }}
                    />
                    <button type="button" className="btn btn-primary" onClick={doActivate} disabled={activating || !key} style={{ padding: "10px 16px" }}>
                      {activating ? "Activating…" : "Activate"}
                    </button>
                  </div>
                )}
              </div>
            </div>
          </>
        )}
      </main>

      <ConfirmDialog
        open={confirmDeactivate}
        title="Deactivate license"
        message="Are you sure you want to deactivate this license? The product will revert to the community edition."
        confirmLabel="Deactivate"
        danger
        loading={submitting}
        onConfirm={doDeactivate}
        onCancel={() => setConfirmDeactivate(false)}
      />
    </div>
  );
}
