// Login page
import { useTranslations } from "next-intl";
import Link from "next/link";

export default function LoginPage() {
  const t = useTranslations("login");

  return (
    <main>
      <h1>{t("title")}</h1>
      <form>
        <label>
          {t("email")}
          <input type="email" name="email" required />
        </label>
        <label>
          {t("password")}
          <input type="password" name="password" required />
        </label>
        <button type="submit">{t("submit")}</button>
      </form>
      <nav>
        <Link href="/en/login" locale="en">English</Link> | <Link href="/bn/login" locale="bn">বাংলা</Link>
      </nav>
    </main>
  );
}