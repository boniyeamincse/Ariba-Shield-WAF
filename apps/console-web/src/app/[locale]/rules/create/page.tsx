"use client";

import { useLocale } from "next-intl";
import { useRouter } from "next/navigation";
import Sidebar from "@/components/layout/Sidebar";
import UserProfileWidget from "@/components/UserProfileWidget";
import RuleWizard from "@/components/rules/RuleWizard";

export default function CreateRulePage() {
  const locale = useLocale();
  const router = useRouter();

  return (
    <div className="dashboard-container">
      <Sidebar locale={locale} />
      <main className="main-content">
        <div className="top-header animate-fade-in">
          <div className="header-title">
            <button type="button" onClick={() => router.push(`/${locale}/rules`)} style={{ background: "none", border: "none", color: "var(--text-secondary)", cursor: "pointer", padding: 0, marginBottom: "4px", fontSize: "13px" }}>
              ← Back to rules
            </button>
            <h1>Create Rule</h1>
            <p style={{ color: "var(--text-secondary)" }}>Build a WAF detection rule step by step.</p>
          </div>
          <div className="header-actions" style={{ display: "flex", alignItems: "center", gap: "12px" }}>
            <UserProfileWidget />
          </div>
        </div>

        <div className="data-section animate-fade-in delay-1" style={{ display: "flex", justifyContent: "center" }}>
          <RuleWizard onCancel={() => router.push(`/${locale}/rules`)} />
        </div>
      </main>
    </div>
  );
}
