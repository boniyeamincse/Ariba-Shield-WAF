"use client";

import React, { createContext, useContext, useEffect, useState } from "react";
import { useRouter, usePathname } from "next/navigation";
import { API_BASE } from "@/lib/api";

interface UserProfile {
  id: string;
  email: string;
  role: string;
}

interface AuthContextType {
  user: UserProfile | null;
  loading: boolean;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType>({
  user: null,
  loading: true,
  logout: async () => {},
});

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<UserProfile | null>(null);
  const [loading, setLoading] = useState(true);
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    const fetchUser = async () => {
      try {
        // Cookie-based session (HTTP-only shield_session cookie).
        const res = await fetch(`${API_BASE}/api/v1/auth/me`, {
          credentials: "include",
        });

        if (res.ok) {
          const data = await res.json();
          setUser(data.user);
        } else {
          setUser(null);
          if (!pathname.includes("/login")) {
            router.push("/en/login");
          }
        }
      } catch {
        setUser(null);
        if (!pathname.includes("/login")) {
          router.push("/en/login");
        }
      } finally {
        setLoading(false);
      }
    };

    fetchUser();
  }, [pathname, router]);

  const logout = async () => {
    try {
      await fetch(`${API_BASE}/api/v1/auth/logout`, {
        method: "POST",
        credentials: "include",
      });
    } catch {
      // Ignore; still clear client state.
    } finally {
      setUser(null);
      router.push("/en/login");
    }
  };

  return (
    <AuthContext.Provider value={{ user, loading, logout }}>
      {loading && !pathname.includes("/login") ? (
        <div style={{ height: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center', backgroundColor: '#0a0a0f', color: '#60a5fa' }}>
          Loading Ariba Shield...
        </div>
      ) : (
        children
      )}
    </AuthContext.Provider>
  );
}

export const useAuth = () => useContext(AuthContext);