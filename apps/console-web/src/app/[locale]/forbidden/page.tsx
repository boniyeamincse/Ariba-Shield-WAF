import Sidebar from "@/components/layout/Sidebar";
import { getLocale } from "next-intl/server";
import Link from "next/link";

export default async function ForbiddenPage() {
  const locale = await getLocale();

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
            <div style={{ fontSize: "72px", fontWeight: 700, color: "var(--warning)", fontFamily: "'Outfit', sans-serif" }}>
              403
            </div>
            <p style={{ color: "var(--text-secondary)", fontSize: "15px" }}>
              You do not have permission to access this page.
            </p>
            <Link href={`/${locale}`} className="btn btn-primary">
              Go to Dashboard
            </Link>
          </div>
        </div>
      </main>
    </div>
  );
}
