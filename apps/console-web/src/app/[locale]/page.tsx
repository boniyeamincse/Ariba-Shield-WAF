import { getTranslations } from "next-intl/server";
import { listApplications, listGateways, listSecurityPolicies } from "@/lib/api";

export default async function OverviewPage() {
  const t = await getTranslations("overview");
  const navT = await getTranslations("nav");

  let apps: { id: string; name: string; status: string }[] = [];
  let gateways: { id: string; hostname: string; status: string }[] = [];
  let policies: { id: string; name: string; enforcement_mode: string }[] = [];

  try {
    [apps, gateways, policies] = await Promise.all([
      listApplications(),
      listGateways(),
      listSecurityPolicies(),
    ]);
  } catch {
    // API not available — show empty state
  }

  return (
    <main>
      <h1>{t("title")}</h1>
      <p>{t("description")}</p>

      <section>
        <h2>{t("applications")} ({apps.length})</h2>
        {apps.length === 0 ? (
          <p>{t("no_applications")}</p>
        ) : (
          <ul>
            {apps.map((a) => (
              <li key={a.id}>{a.name} — {a.status}</li>
            ))}
          </ul>
        )}
      </section>

      <section>
        <h2>{navT("gateways")} ({gateways.length})</h2>
        {gateways.length === 0 ? (
          <p>{t("no_gateways")}</p>
        ) : (
          <ul>
            {gateways.map((g) => (
              <li key={g.id}>{g.hostname} — {g.status}</li>
            ))}
          </ul>
        )}
      </section>

      <section>
        <h2>{navT("policies")} ({policies.length})</h2>
        {policies.length === 0 ? (
          <p>{t("no_policies")}</p>
        ) : (
          <ul>
            {policies.map((p) => (
              <li key={p.id}>{p.name} — {p.enforcement_mode}</li>
            ))}
          </ul>
        )}
      </section>
    </main>
  );
}