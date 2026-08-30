"use client";

import { useAuth } from "@/context/AuthContext";
import { can, canAccessPage } from "@/lib/permissions";

/** Hook for permission-aware UI: check if the current user can perform an action. */
export function usePermission() {
  const { user } = useAuth();
  const role = user?.role;

  return {
    role,
    /** e.g. can("create"), can("delete"), can("activate") */
    can: (action: string) => can(role, action),
    /** Check page access by path. */
    canAccessPage: (path: string) => canAccessPage(role, path),
    /** True if the user has at least write capability. */
    canWrite: () => can(role, "create") || can(role, "edit"),
    /** True if the user is an admin. */
    isAdmin: () => can(role, "admin"),
  };
}
