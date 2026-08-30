"use client";

import { useState } from "react";

export type FilterDef =
  | { type: "text"; key: string; label: string; placeholder?: string }
  | { type: "select"; key: string; label: string; options: { value: string; label: string }[] }
  | { type: "date"; key: string; label: string };

export type FilterValues = Record<string, string>;

type FilterBarProps = {
  filters: FilterDef[];
  values: FilterValues;
  onChange: (values: FilterValues) => void;
};

export default function FilterBar({ filters, values, onChange }: FilterBarProps) {
  const [search, setSearch] = useState("");

  const setValue = (key: string, value: string) => {
    onChange({ ...values, [key]: value });
  };

  const clearAll = () => {
    setSearch("");
    onChange({});
  };

  const hasActive = Object.values(values).some((v) => v !== "" && v !== undefined);

  return (
    <div
      className="glass-panel"
      style={{
        padding: "14px 16px",
        marginBottom: "16px",
        display: "flex",
        flexWrap: "wrap",
        alignItems: "center",
        gap: "12px",
      }}
    >
      {filters.map((filter) => {
        if (filter.type === "text") {
          return (
            <input
              key={filter.key}
              type="text"
              placeholder={filter.placeholder ?? filter.label}
              value={search}
              onChange={(e) => {
                setSearch(e.target.value);
                setValue(filter.key, e.target.value);
              }}
              style={{
                background: "rgba(255,255,255,0.05)",
                border: "1px solid rgba(255,255,255,0.1)",
                padding: "8px 12px",
                borderRadius: "8px",
                color: "white",
                outline: "none",
                fontSize: "13px",
                minWidth: "160px",
              }}
            />
          );
        }
        if (filter.type === "select") {
          return (
            <select
              key={filter.key}
              value={values[filter.key] ?? ""}
              onChange={(e) => setValue(filter.key, e.target.value)}
              style={{
                background: "rgba(255,255,255,0.05)",
                border: "1px solid rgba(255,255,255,0.1)",
                padding: "8px 12px",
                borderRadius: "8px",
                color: "white",
                outline: "none",
                fontSize: "13px",
              }}
            >
              <option value="">{filter.label}</option>
              {filter.options.map((opt) => (
                <option key={opt.value} value={opt.value} style={{ background: "#13141c", color: "#fff" }}>
                  {opt.label}
                </option>
              ))}
            </select>
          );
        }
        // date
        return (
          <input
            key={filter.key}
            type="date"
            value={values[filter.key] ?? ""}
            onChange={(e) => setValue(filter.key, e.target.value)}
            style={{
              background: "rgba(255,255,255,0.05)",
              border: "1px solid rgba(255,255,255,0.1)",
              padding: "8px 12px",
              borderRadius: "8px",
              colorScheme: "dark",
              color: "white",
              outline: "none",
              fontSize: "13px",
            }}
          />
        );
      })}

      {hasActive && (
        <button type="button" className="btn" onClick={clearAll} style={{ padding: "8px 14px", fontSize: "13px", marginTop: 0 }}>
          Clear
        </button>
      )}
    </div>
  );
}
