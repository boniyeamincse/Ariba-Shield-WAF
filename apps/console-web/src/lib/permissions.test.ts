import { describe, it, expect } from "vitest";
import { can, canAccessPage, navItemsForRole } from "@/lib/permissions";

describe("permissions", () => {
  it("can() returns true for allowed actions", () => {
    expect(can("Super Admin", "delete")).toBe(true);
    expect(can("Security Admin", "create")).toBe(true);
    expect(can("App Owner", "edit")).toBe(true);
  });

  it("can() returns false for disallowed actions", () => {
    expect(can("Read Only", "delete")).toBe(false);
    expect(can("SOC Analyst", "create")).toBe(false);
    expect(can("Auditor", "edit")).toBe(false);
  });

  it("can() handles missing role", () => {
    expect(can(undefined, "read")).toBe(false);
    expect(can("", "read")).toBe(false);
  });

  it("canAccessPage gates admin-only pages", () => {
    expect(canAccessPage("Super Admin", "/users")).toBe(true);
    expect(canAccessPage("Read Only", "/users")).toBe(false);
    expect(canAccessPage("Auditor", "/audit")).toBe(true);
    expect(canAccessPage("Read Only", "/audit")).toBe(false);
  });

  it("canAccessPage allows pages with no role restriction", () => {
    expect(canAccessPage("Read Only", "/applications")).toBe(true);
    expect(canAccessPage("App Owner", "/overview")).toBe(true);
  });

  it("navItemsForRole filters nav by role", () => {
    const admin = navItemsForRole("Super Admin");
    expect(admin.platform.some((i) => i.href === "/users")).toBe(true);
    expect(admin.platform.some((i) => i.href === "/license")).toBe(true);

    const auditor = navItemsForRole("Auditor");
    expect(auditor.platform.some((i) => i.href === "/users")).toBe(false);
    expect(auditor.platform.some((i) => i.href === "/audit")).toBe(true);
  });
});
