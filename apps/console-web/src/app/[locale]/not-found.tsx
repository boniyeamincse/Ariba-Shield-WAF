import Sidebar from "@/components/layout/Sidebar";
import { getLocale, getTranslations } from "next-intl/server";
import Link from "next/link";

export default async function NotFound() {
  const locale = await getLocale();
  const t = await getTranslations("notFound");

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />
      <main className="main-content">
        <div className="data-section animate-fade-in">
          <div
            className="glass-panel"
            style={{
              padding: "80px 40px",
              textAlign: "center",
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
              gap: "16px",
            }}
          >
            <div style={{ fontSize: "72px", fontWeight: 700, color: "var(--accent-primary)", fontFamily: "'Outfit', sans-serif" }}>
              404
            </div>
            <p style={{ color: "var(--text-secondary)", fontSize: "15px" }}>
              {t("message")}
            </p>
            <Link href={`/${locale}`} className="btn btn-primary">
              {t("home")}
            </Link>
          </div>
        </div>
      </main>
    </div>
  );
}
