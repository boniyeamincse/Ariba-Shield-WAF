"use client";

import React, { createContext, useContext, useEffect, useState } from "react";
import { useRouter, usePathname } from "next/navigation";

interface UserProfile {
  id: string;
  name: string;
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
        const res = await fetch("http://localhost:8080/api/v1/auth/me", {
          // If using HTTP-only cookies, include credentials
          credentials: "omit", 
          headers: {
            "Authorization": `Bearer ${localStorage.getItem("token") || ""}`
          }
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
      } catch (err) {
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
      await fetch("http://localhost:8080/api/v1/auth/logout", {
        method: "POST",
      });
    } catch (e) {
      console.error("Logout failed", e);
    } finally {
      localStorage.removeItem("token");
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
