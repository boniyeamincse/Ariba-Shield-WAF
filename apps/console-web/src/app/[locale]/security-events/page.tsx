import { getTranslations } from "next-intl/server";
import { API_BASE } from "@/lib/api";

type SecurityEvent = {
  id: string;
  event_id: string;
  request_id: string;
  timestamp: string;
  severity: string;
  decision_action: string;
  reason: string;
  rule_ids: string[];
  client_ip: string;
  method: string;
  path: string;
  status: number;
  created_at: string;
};

async function fetchSecurityEvents(): Promise<SecurityEvent[]> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/security-events`, { cache: "no-store" });
    if (!res.ok) return [];
    const data = await res.json();
    return Array.isArray(data) ? data : (data.events ?? []);
  } catch {
    return [];
  }
}

export default async function SecurityEventsPage() {
  const t = await getTranslations("security_events");
  const events = await fetchSecurityEvents();

  return (
    <main>
      <h1>{t("title")}</h1>
      {events.length === 0 ? (
        <p>{t("no_events")}</p>
      ) : (
        <table>
          <thead>
            <tr>
              <th>{t("severity")}</th>
              <th>{t("time")}</th>
              <th>{t("method")}</th>
              <th>{t("path")}</th>
              <th>{t("reason")}</th>
              <th>{t("rules")}</th>
              <th>{t("client_ip")}</th>
            </tr>
          </thead>
          <tbody>
            {events.map((ev) => (
              <tr key={ev.id}>
                <td>{ev.severity}</td>
                <td>{new Date(ev.created_at).toLocaleString()}</td>
                <td>{ev.method}</td>
                <td>{ev.path}</td>
                <td>{ev.reason}</td>
                <td>{ev.rule_ids?.join(", ") || "—"}</td>
                <td>{ev.client_ip}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </main>
  );
}