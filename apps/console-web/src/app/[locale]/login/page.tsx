"use client";

import { useTranslations } from "next-intl";
import Link from "next/link";
import { useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { API_BASE } from "@/lib/api";

export default function LoginPage() {
  const t = useTranslations("login");
  const router = useRouter();
  const params = useParams<{ locale: string }>();
  const locale = params?.locale ?? "en";
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const res = await fetch(`${API_BASE}/api/v1/auth/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ email, password }),
      });

      if (!res.ok) {
        throw new Error("Invalid credentials");
      }

      router.push(`/${locale}`);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Failed to login");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{
      minHeight: '100vh',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      backgroundImage: 'radial-gradient(circle at 15% 50%, rgba(59, 130, 246, 0.15), transparent 25%), radial-gradient(circle at 85% 30%, rgba(16, 185, 129, 0.15), transparent 25%)',
      backgroundColor: '#0a0a0f',
      fontFamily: "'Inter', sans-serif"
    }}>
      <div className="glass-panel animate-fade-in" style={{
        width: '100%',
        maxWidth: '420px',
        padding: '40px',
        display: 'flex',
        flexDirection: 'column',
        gap: '24px'
      }}>
        
        <div style={{ textAlign: 'center' }}>
          <div className="brand" style={{ justifyContent: 'center', marginBottom: '8px' }}>
            <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path>
            </svg>
            Ariba Shield
          </div>
          <p style={{ color: '#9ca3af', fontSize: '14px' }}>Enterprise WAF Management Console</p>
        </div>

        <form onSubmit={handleLogin} style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
          {error && (
            <div style={{ padding: '12px', background: 'rgba(239, 68, 68, 0.1)', color: '#f87171', borderRadius: '8px', fontSize: '14px', border: '1px solid rgba(239, 68, 68, 0.2)' }}>
              {error}
            </div>
          )}

          <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
            <label style={{ fontSize: '14px', fontWeight: 500, color: '#e5e7eb' }}>{t("email")}</label>
            <input 
              type="email" 
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required 
              style={{
                background: 'rgba(255, 255, 255, 0.05)',
                border: '1px solid rgba(255, 255, 255, 0.1)',
                padding: '12px 16px',
                borderRadius: '8px',
                color: 'white',
                outline: 'none',
                transition: 'border-color 0.3s'
              }}
              onFocus={(e) => e.target.style.borderColor = '#3b82f6'}
              onBlur={(e) => e.target.style.borderColor = 'rgba(255, 255, 255, 0.1)'}
            />
          </div>

          <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
            <label style={{ fontSize: '14px', fontWeight: 500, color: '#e5e7eb' }}>{t("password")}</label>
            <input 
              type="password" 
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required 
              style={{
                background: 'rgba(255, 255, 255, 0.05)',
                border: '1px solid rgba(255, 255, 255, 0.1)',
                padding: '12px 16px',
                borderRadius: '8px',
                color: 'white',
                outline: 'none',
                transition: 'border-color 0.3s'
              }}
              onFocus={(e) => e.target.style.borderColor = '#3b82f6'}
              onBlur={(e) => e.target.style.borderColor = 'rgba(255, 255, 255, 0.1)'}
            />
          </div>

          <button 
            type="submit" 
            disabled={loading}
            className="btn btn-primary"
            style={{ width: '100%', justifyContent: 'center', marginTop: '8px', padding: '12px' }}
          >
            {loading ? "Authenticating..." : t("submit")}
          </button>
        </form>

        <div style={{ textAlign: 'center', marginTop: '16px', fontSize: '13px', color: '#9ca3af', display: 'flex', justifyContent: 'center', gap: '12px' }}>
          <Link href="/en/login" locale="en" style={{ color: '#60a5fa', textDecoration: 'none' }}>English</Link> 
          <span>|</span> 
          <Link href="/bn/login" locale="bn" style={{ color: '#60a5fa', textDecoration: 'none' }}>বাংলা</Link>
        </div>
      </div>
    </div>
  );
}