import { useTranslations } from "next-intl";

export default function OverviewPage() {
  const t = useTranslations("overview");

  return (
    <main>
      <h1>{t("title")}</h1>
      <p>{t("description")}</p>
      <section>
        <h2>{t("applications")}</h2>
        <p>{t("no_applications")}</p>
      </section>
    </main>
  );
}