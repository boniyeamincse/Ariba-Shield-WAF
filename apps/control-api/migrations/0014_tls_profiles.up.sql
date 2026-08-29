-- 0014_tls_profiles.up.sql
-- TLS profiles (master plan §8: tls_profiles).

CREATE TABLE IF NOT EXISTS tls_profiles (
    id                 TEXT PRIMARY KEY,     -- ULID
    organization_id    TEXT NOT NULL REFERENCES organizations(id),
    name               TEXT NOT NULL,
    min_version        TEXT NOT NULL DEFAULT '1.2',
    max_version        TEXT NOT NULL DEFAULT '1.3',
    cipher_suites      TEXT[] NOT NULL DEFAULT '{}',
    certificate_ref    TEXT NOT NULL DEFAULT '',
    renegotiation      TEXT NOT NULL DEFAULT 'deny',
    status             TEXT NOT NULL DEFAULT 'active',
    version            BIGINT NOT NULL DEFAULT 1,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE certificates ADD COLUMN IF NOT EXISTS chain_pem TEXT;
