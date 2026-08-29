"use client";

import { useAuth } from "../context/AuthContext";

export default function UserProfileWidget() {
  const { user, logout } = useAuth();

  if (!user) return null;

  return (
    <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
      <div style={{ textAlign: 'right' }}>
        <div style={{ fontSize: '14px', fontWeight: 600, color: '#e5e7eb' }}>{user.email}</div>
        <div style={{ fontSize: '12px', color: '#9ca3af' }}>{user.role}</div>
      </div>
      <div 
        style={{
          width: '36px',
          height: '36px',
          borderRadius: '50%',
          backgroundColor: 'rgba(59, 130, 246, 0.2)',
          border: '1px solid rgba(59, 130, 246, 0.4)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          color: '#60a5fa',
          fontWeight: 'bold',
          fontSize: '14px'
        }}
      >
        {(user.email || user.role || "?").charAt(0)}
      </div>
      <button 
        onClick={logout}
        style={{
          background: 'rgba(239, 68, 68, 0.1)',
          border: '1px solid rgba(239, 68, 68, 0.3)',
          color: '#f87171',
          padding: '6px 12px',
          borderRadius: '6px',
          cursor: 'pointer',
          fontSize: '12px',
          fontWeight: 600,
          transition: 'all 0.2s ease'
        }}
        onMouseOver={(e) => e.currentTarget.style.background = 'rgba(239, 68, 68, 0.2)'}
        onMouseOut={(e) => e.currentTarget.style.background = 'rgba(239, 68, 68, 0.1)'}
      >
        Logout
      </button>
    </div>
  );
}
