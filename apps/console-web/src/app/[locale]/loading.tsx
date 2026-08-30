import Sidebar from "@/components/layout/Sidebar";
import { getLocale, getTranslations } from "next-intl/server";

export default async function Loading() {
  const locale = await getLocale();
  const t = await getTranslations("loading");

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />
      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <h1>{t("title")}</h1>
            <p style={{ color: "var(--text-secondary)" }}>{t("description")}</p>
          </div>
        </div>

        <div className="data-section animate-fade-in delay-1">
          <div className="glass-panel" style={{ padding: "24px" }}>
            <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
              <div className="skeleton-row" style={{ width: "40%" }} />
              <div className="skeleton-row" style={{ width: "70%" }} />
              <div className="skeleton-row" style={{ width: "55%" }} />
              <div className="skeleton-row" style={{ width: "85%" }} />
              <div className="skeleton-row" style={{ width: "30%" }} />
            </div>
          </div>
        </div>
      </main>

      <style dangerouslySetInnerHTML={{ __html: `
        .skeleton-row {
          height: 18px;
          border-radius: 6px;
          background: linear-gradient(90deg, rgba(255,255,255,0.05) 25%, rgba(255,255,255,0.12) 50%, rgba(255,255,255,0.05) 75%);
          background-size: 200% 100%;
          animation: shimmer 1.5s infinite;
        }
        @keyframes shimmer {
          0% { background-position: 200% 0; }
          100% { background-position: -200% 0; }
        }
      `}} />
    </div>
  );
}
