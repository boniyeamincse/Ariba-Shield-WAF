import { getTranslations } from "next-intl/server";
import { listGateways } from "@/lib/api";

export default async function GatewaysPage() {
  const t = await getTranslations("gateways");
  const gateways = await listGateways().catch(() => []);

  return (
    <main>
      <h1>{t("title")}</h1>
      {gateways.length === 0 ? (
        <p>{t("no_gateways")}</p>
      ) : (
        <table>
          <thead>
            <tr>
              <th>{t("hostname")}</th>
              <th>{t("status")}</th>
              <th>{t("version")}</th>
              <th>{t("last_seen")}</th>
              <th>{t("applied_hash")}</th>
            </tr>
          </thead>
          <tbody>
            {gateways.map((g) => (
              <tr key={g.id}>
                <td>{g.hostname}</td>
                <td>{g.status}</td>
                <td>{g.version}</td>
                <td>{g.last_seen_at ? new Date(g.last_seen_at).toLocaleString() : "—"}</td>
                <td>{g.applied_hash ? g.applied_hash.slice(0, 12) + "…" : "—"}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </main>
  );
}