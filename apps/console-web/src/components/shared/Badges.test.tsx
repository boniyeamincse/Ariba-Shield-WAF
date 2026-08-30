import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { StatusBadge, SeverityBadge } from "@/components/shared/Badges";

describe("StatusBadge", () => {
  it("renders the value", () => {
    render(<StatusBadge value="active" />);
    expect(screen.getByText("active")).toBeInTheDocument();
  });

  it("falls back to em dash for missing value", () => {
    render(<StatusBadge />);
    expect(screen.getByText("—")).toBeInTheDocument();
  });
});

describe("SeverityBadge", () => {
  it("renders severity value", () => {
    render(<SeverityBadge value="critical" />);
    expect(screen.getByText("critical")).toBeInTheDocument();
  });

  it("handles missing value", () => {
    render(<SeverityBadge />);
    expect(screen.getByText("—")).toBeInTheDocument();
  });
});
