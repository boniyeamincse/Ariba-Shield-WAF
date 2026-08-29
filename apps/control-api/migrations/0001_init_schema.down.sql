-- 0001_init_schema.down.sql
-- Reverse of 0001_init_schema.up.sql (order respects FKs).

DROP TABLE IF EXISTS access_events;
DROP TABLE IF EXISTS config_deployments;
DROP TABLE IF EXISTS config_versions;
DROP TABLE IF EXISTS health_monitors;
DROP TABLE IF EXISTS backend_nodes;
DROP TABLE IF EXISTS backend_pools;
DROP TABLE IF EXISTS routes;
DROP TABLE IF EXISTS virtual_servers;
DROP TABLE IF EXISTS applications;
DROP TABLE IF EXISTS audit_events;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS user_group_memberships;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS organizations;
