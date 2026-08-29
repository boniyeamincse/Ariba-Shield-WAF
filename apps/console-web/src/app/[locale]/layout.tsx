import type { Metadata } from "next";
import { setRequestLocale } from "next-intl/server";
import { notFound } from "next/navigation";

export const metadata: Metadata = {
  title: "Ariba Shield WAF",
  description: "Centralized enterprise Web Application Firewall management console",
};

const locales = ["en", "bn"];

// Locale segment layout: validates the locale and pins it for this subtree.
// <html>/<body> and providers live in the root layout (src/app/layout.tsx).
export default async function LocaleLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  if (!locales.includes(locale)) notFound();

  setRequestLocale(locale);
  return children;
}