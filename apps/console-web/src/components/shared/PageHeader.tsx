"use client";

import type { ReactNode } from "react";

export type HeaderAction = {
  label: string;
  onClick?: () => void;
  variant?: "primary" | "secondary" | "danger";
  icon?: ReactNode;
};

type PageHeaderProps = {
  title: string;
  description?: string;
  actions?: HeaderAction[];
  children?: ReactNode;
};

export default function PageHeader({ title, description, actions, children }: PageHeaderProps) {
  return (
    <div className="top-header animate-fade-in">
      <div className="header-title">
        <h1>{title}</h1>
        {description && <p style={{ color: "var(--text-secondary)" }}>{description}</p>}
      </div>
      <div className="header-actions" style={{ display: "flex", alignItems: "center", gap: "12px" }}>
        {actions?.map((action) => (
          <button
            key={action.label}
            type="button"
            className={`btn ${action.variant === "danger" ? "" : "btn-primary"}`}
            onClick={action.onClick}
            style={{
              ...(action.variant === "danger"
                ? { background: "var(--danger)", borderColor: "var(--danger)", color: "#fff" }
                : {}),
              ...(action.variant === "secondary" ? { background: "transparent" } : {}),
            }}
          >
            {action.icon}
            {action.label}
          </button>
        ))}
        {children}
      </div>
    </div>
  );
}
