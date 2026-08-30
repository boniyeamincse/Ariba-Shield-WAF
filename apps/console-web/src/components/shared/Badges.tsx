const STATUS_COLORS: Record<string, { bg: string; text: string }> = {
  active: { bg: "rgba(16,185,129,0.15)", text: "#34d399" },
  enabled: { bg: "rgba(16,185,129,0.15)", text: "#34d399" },
  healthy: { bg: "rgba(16,185,129,0.15)", text: "#34d399" },
  online: { bg: "rgba(16,185,129,0.15)", text: "#34d399" },
  starting: { bg: "rgba(245,158,11,0.15)", text: "#fbbf24" },
  pending: { bg: "rgba(245,158,11,0.15)", text: "#fbbf24" },
  investigating: { bg: "rgba(245,158,11,0.15)", text: "#fbbf24" },
  degraded: { bg: "rgba(245,158,11,0.15)", text: "#fbbf24" },
  inactive: { bg: "rgba(156,163,175,0.15)", text: "#9ca3af" },
  disabled: { bg: "rgba(156,163,175,0.15)", text: "#9ca3af" },
  offline: { bg: "rgba(239,68,68,0.15)", text: "#f87171" },
  unregistered: { bg: "rgba(156,163,175,0.15)", text: "#9ca3af" },
  failed: { bg: "rgba(239,68,68,0.15)", text: "#f87171" },
  blocked: { bg: "rgba(239,68,68,0.15)", text: "#f87171" },
  open: { bg: "rgba(59,130,246,0.15)", text: "#60a5fa" },
  resolved: { bg: "rgba(16,185,129,0.15)", text: "#34d399" },
  accepted: { bg: "rgba(16,185,129,0.15)", text: "#34d399" },
  rejected: { bg: "rgba(239,68,68,0.15)", text: "#f87171" },
  ready: { bg: "rgba(16,185,129,0.15)", text: "#34d399" },
  critical: { bg: "rgba(239,68,68,0.15)", text: "#f87171" },
  high: { bg: "rgba(245,158,11,0.15)", text: "#fbbf24" },
  medium: { bg: "rgba(59,130,246,0.15)", text: "#60a5fa" },
  low: { bg: "rgba(156,163,175,0.15)", text: "#9ca3af" },
  info: { bg: "rgba(156,163,175,0.15)", text: "#9ca3af" },
};

type BadgeProps = {
  value?: string;
  label?: string;
};

export function StatusBadge({ value, label }: BadgeProps) {
  const normalized = String(value ?? "").toLowerCase();
  const colors = STATUS_COLORS[normalized] ?? { bg: "rgba(255,255,255,0.06)", text: "#9ca3af" };
  return (
    <span
      style={{
        display: "inline-block",
        padding: "4px 10px",
        borderRadius: "6px",
        background: colors.bg,
        color: colors.text,
        fontSize: "12px",
        fontWeight: 600,
        whiteSpace: "nowrap",
      }}
    >
      {label ?? String(value ?? "—")}
    </span>
  );
}

export function SeverityBadge({ value }: { value?: string }) {
  const normalized = String(value ?? "").toLowerCase();
  const colors = STATUS_COLORS[normalized] ?? { bg: "rgba(255,255,255,0.06)", text: "#9ca3af" };
  return (
    <span
      style={{
        display: "inline-block",
        padding: "4px 10px",
        borderRadius: "6px",
        background: colors.bg,
        color: colors.text,
        fontSize: "12px",
        fontWeight: 600,
        textTransform: "capitalize",
        whiteSpace: "nowrap",
      }}
    >
      {String(value ?? "—")}
    </span>
  );
}
