import { getTranslations } from "next-intl/server";
import { listApplications, listOrigins, listDomains } from "@/lib/api";

export default async function ApplicationsPage() {
  const t = await getTranslations("applications");
  const apps = await listApplications().catch(() => []);

  const rows = await Promise.all(
    apps.map(async (app) => {
      const [origins, domains] = await Promise.all([
        listOrigins(app.id).catch(() => []),
        listDomains(app.id).catch(() => []),
      ]);
      return { app, origins, domains };
    })
  );

  return (
    <main>
      <h1>{t("title")}</h1>
      {rows.length === 0 ? (
        <p>{t("no_applications")}</p>
      ) : (
        <ul>
          {rows.map(({ app, origins, domains }) => (
            <li key={app.id}>
              <strong>{app.name}</strong> — {app.status}
              {domains.length > 0 && (
                <p>{t("domains")}: {domains.map((d) => d.hostname).join(", ")}</p>
              )}
              {origins.length > 0 && (
                <p>{t("origins")}: {origins.map((o) => `${o.host}:${o.port}`).join(", ")}</p>
              )}
            </li>
          ))}
        </ul>
      )}
    </main>
  );
}