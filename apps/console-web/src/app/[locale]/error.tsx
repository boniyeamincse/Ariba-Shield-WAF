"use client";

import Sidebar from "@/components/layout/Sidebar";
import { useLocale, useTranslations } from "next-intl";

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  const locale = useLocale();
  const t = useTranslations("error");

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />
      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <h1>{t("title")}</h1>
          </div>
        </div>

        <div className="data-section animate-fade-in delay-1">
          <div
            className="glass-panel"
            style={{
              padding: "60px",
              textAlign: "center",
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
              gap: "16px",
            }}
          >
            <svg
              width="48"
              height="48"
              viewBox="0 0 24 24"
              fill="none"
              stroke="var(--danger)"
              strokeWidth="1.5"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="8" x2="12" y2="12"></line>
              <line x1="12" y1="16" x2="12.01" y2="16"></line>
            </svg>
            <p style={{ color: "var(--text-secondary)", fontSize: "14px" }}>
              {t("message")}
            </p>
            {error?.message && (
              <p style={{ color: "var(--text-secondary)", fontSize: "12px", fontFamily: "monospace" }}>
                {error.message}
              </p>
            )}
            <button type="button" className="btn btn-primary" onClick={reset}>
              {t("retry")}
            </button>
          </div>
        </div>
      </main>
    </div>
  );
}
