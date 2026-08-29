-- 0016_policy_approvals.up.sql
-- Four-eyes approval workflow (master plan §7.1: approval workflows,
-- §6.2: policy diff, review, approval, rollback).

CREATE TABLE IF NOT EXISTS policy_approvals (
    id              TEXT PRIMARY KEY,        -- ULID
    organization_id TEXT NOT NULL REFERENCES organizations(id),
    policy_id       TEXT NOT NULL REFERENCES security_policies(id) ON DELETE CASCADE,
    policy_version_id TEXT REFERENCES policy_versions(id),
    status          TEXT NOT NULL DEFAULT 'pending', -- pending | approved | rejected
    requested_by    TEXT REFERENCES users(id),
    approved_by     TEXT REFERENCES users(id),
    reviewer_notes  TEXT NOT NULL DEFAULT '',
    reviewed_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
