"use client";

import { useState } from "react";
import { useLocale } from "next-intl";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import { enableMFA, verifyMFA, disableMFA } from "@/lib/api";

export default function MFAPage() {
  const locale = useLocale();
  const [enrolling, setEnrolling] = useState(false);
  const [secret, setSecret] = useState("");
  const [otpauthUrl, setOtpauthUrl] = useState("");
  const [code, setCode] = useState("");
  const [status, setStatus] = useState<"idle" | "enabled" | "disabled">("idle");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const startEnroll = async () => {
    setLoading(true);
    setError("");
    try {
      const { secret, otpauth_url } = await enableMFA();
      setSecret(secret);
      setOtpauthUrl(otpauth_url);
      setEnrolling(true);
    } catch {
      setError("Failed to start MFA enrollment");
    } finally {
      setLoading(false);
    }
  };

  const confirmEnroll = async () => {
    setLoading(true);
    setError("");
    try {
      await verifyMFA(code);
      setEnrolling(false);
      setCode("");
      setSecret("");
      setOtpauthUrl("");
      setStatus("enabled");
    } catch {
      setError("Invalid verification code");
    } finally {
      setLoading(false);
    }
  };

  const doDisable = async () => {
    setLoading(true);
    setError("");
    try {
      await disableMFA();
      setStatus("disabled");
    } catch {
      setError("Failed to disable MFA");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />

      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <h1>MFA / Two-Factor Authentication</h1>
            <p style={{ color: "var(--text-secondary)" }}>Protect your account with TOTP.</p>
          </div>
          <div className="header-actions" style={{ display: "flex", alignItems: "center", gap: "12px" }}>
            <UserProfileWidget />
          </div>
        </div>

        {error && (
          <div className="glass-panel" style={{ padding: "12px 16px", marginBottom: "16px", color: "var(--danger)", fontSize: "13px", border: "1px solid rgba(239,68,68,0.3)" }}>
            {error}
          </div>
        )}

        <div className="data-section animate-fade-in delay-1" style={{ maxWidth: "560px" }}>
          <div className="glass-panel" style={{ padding: "28px", display: "flex", flexDirection: "column", gap: "20px" }}>
            {status === "enabled" ? (
              <>
                <div style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                  <span style={{ color: "var(--success)", fontSize: "18px" }}>✓</span>
                  <span style={{ fontSize: "15px", color: "var(--text-primary)" }}>MFA is enabled for your account.</span>
                </div>
                <button type="button" className="btn" onClick={doDisable} disabled={loading} style={{ padding: "10px 16px", color: "var(--danger)", borderColor: "var(--danger)", alignSelf: "flex-start" }}>
                  {loading ? "Disabling…" : "Disable MFA"}
                </button>
              </>
            ) : enrolling ? (
              <>
                <h3 style={{ fontSize: "16px", fontWeight: 600 }}>Scan with your authenticator app</h3>
                <p style={{ fontSize: "13px", color: "var(--text-secondary)" }}>
                  Scan the QR code (or manually enter the secret) in your authenticator app, then enter the 6-digit code below.
                </p>
                {otpauthUrl && (
                  <div style={{ display: "flex", justifyContent: "center" }}>
                    {/* Render a QR via Google Charts API (public, no key). */}
                    {/* eslint-disable-next-line @next/next/no-img-element */}
                    <img
                      src={`https://api.qrserver.com/v1/create-qr-code/?size=180x180&data=${encodeURIComponent(otpauthUrl)}`}
                      alt="TOTP QR code"
                      width={180}
                      height={180}
                      style={{ borderRadius: "8px", background: "#fff", padding: "8px" }}
                    />
                  </div>
                )}
                <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
                  <label style={{ fontSize: "13px", color: "var(--text-secondary)" }}>Manual secret</label>
                  <code style={{ fontSize: "13px", color: "var(--text-primary)", background: "rgba(255,255,255,0.06)", padding: "10px 12px", borderRadius: "8px", wordBreak: "break-all" }}>
                    {secret}
                  </code>
                </div>
                <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
                  <label style={{ fontSize: "13px", color: "var(--text-secondary)" }}>Verification code</label>
                  <input
                    type="text"
                    value={code}
                    onChange={(e) => setCode(e.target.value)}
                    placeholder="6-digit code"
                    maxLength={6}
                    style={{
                      background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.1)",
                      padding: "10px 12px", borderRadius: "8px", color: "white", outline: "none", fontSize: "16px",
                      letterSpacing: "6px",
                    }}
                  />
                </div>
                <div style={{ display: "flex", justifyContent: "flex-end", gap: "10px" }}>
                  <button type="button" className="btn" onClick={() => setEnrolling(false)} disabled={loading} style={{ padding: "10px 16px" }}>
                    Cancel
                  </button>
                  <button type="button" className="btn btn-primary" onClick={confirmEnroll} disabled={loading || code.length !== 6} style={{ padding: "10px 16px" }}>
                    {loading ? "Verifying…" : "Verify & Enable"}
                  </button>
                </div>
              </>
            ) : (
              <>
                <div style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                  <span style={{ color: "var(--text-secondary)", fontSize: "18px" }}>—</span>
                  <span style={{ fontSize: "15px", color: "var(--text-primary)" }}>MFA is not enabled for your account.</span>
                </div>
                <button type="button" className="btn btn-primary" onClick={startEnroll} disabled={loading} style={{ padding: "10px 16px", alignSelf: "flex-start" }}>
                  {loading ? "Setting up…" : "Enable MFA"}
                </button>
              </>
            )}
          </div>
        </div>
      </main>
    </div>
  );
}
