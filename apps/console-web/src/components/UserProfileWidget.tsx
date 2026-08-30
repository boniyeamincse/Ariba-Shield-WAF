"use client";

import { useAuth } from "../context/AuthContext";

export default function UserProfileWidget() {
  const { user, logout } = useAuth();

  if (!user) return null;

  return (
    <div style={{ 
      display: 'flex', 
      alignItems: 'center', 
      gap: '12px',
      background: 'var(--glass-bg)',
      padding: '6px 16px 6px 6px',
      borderRadius: '999px',
      border: '1px solid var(--glass-border)',
      transition: 'all 0.2s ease',
      boxShadow: '0 4px 12px rgba(0,0,0,0.1)'
    }}>
      <div 
        style={{
          width: '34px',
          height: '34px',
          borderRadius: '50%',
          backgroundColor: 'rgba(59, 130, 246, 0.15)',
          border: '1px solid rgba(59, 130, 246, 0.3)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          color: '#60a5fa',
          fontWeight: '600',
          fontSize: '15px'
        }}
      >
        {(user.email || user.role || "?").charAt(0).toUpperCase()}
      </div>
      <div style={{ display: 'flex', flexDirection: 'column', marginRight: '4px' }}>
        <div style={{ fontSize: '13px', fontWeight: 600, color: '#e5e7eb', lineHeight: '1' }}>
          {user.email.split('@')[0]}
        </div>
        <div style={{ fontSize: '11px', color: '#9ca3af', marginTop: '4px', textTransform: 'capitalize', letterSpacing: '0.02em' }}>
          {user.role}
        </div>
      </div>
      
      <div style={{ width: '1px', height: '24px', background: 'rgba(255,255,255,0.1)' }} />
      
      <button 
        onClick={logout}
        title="Logout"
        style={{
          background: 'transparent',
          border: 'none',
          color: '#9ca3af',
          cursor: 'pointer',
          padding: '6px',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          borderRadius: '50%',
          transition: 'all 0.2s ease'
        }}
        onMouseOver={(e) => { e.currentTarget.style.color = '#ef4444'; e.currentTarget.style.background = 'rgba(239, 68, 68, 0.1)'; }}
        onMouseOut={(e) => { e.currentTarget.style.color = '#9ca3af'; e.currentTarget.style.background = 'transparent'; }}
      >
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
          <polyline points="16 17 21 12 16 7"></polyline>
          <line x1="21" y1="12" x2="9" y2="12"></line>
        </svg>
      </button>
    </div>
  );
}
