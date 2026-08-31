-- 0025_seed_owasp_crs.up.sql
-- Seed default OWASP CRS groups into the managed_rules table.

-- First, ensure the default organization exists (usually seeded in setup, but we use the known dev Org ID)
-- 01ARZ3NDEKTSV4RRFFQ69G5FAV

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, status)
VALUES
  ('01M1BVLQ3J5CRS_SQLI', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'SQL Injection (SQLi) Protection', 'owasp-crs', true, 'low', 'active'),
  ('01M1BVLQ3J5CRS_XSS', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'Cross-Site Scripting (XSS)', 'owasp-crs', true, 'low', 'active'),
  ('01M1BVLQ3J5CRS_LFI', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'Local File Inclusion (LFI)', 'owasp-crs', true, 'low', 'active'),
  ('01M1BVLQ3J5CRS_RFI', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'Remote File Inclusion (RFI)', 'owasp-crs', true, 'low', 'active'),
  ('01M1BVLQ3J5CRS_RCE', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'Remote Code Execution (RCE)', 'owasp-crs', true, 'low', 'active'),
  ('01M1BVLQ3J5CRS_PHP', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'PHP Injection Attacks', 'owasp-crs', true, 'low', 'active'),
  ('01M1BVLQ3J5CRS_JAVA', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'Java Injection Attacks', 'owasp-crs', true, 'low', 'active'),
  ('01M1BVLQ3J5CRS_PROTO', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'HTTP Protocol Validation', 'owasp-crs', true, 'low', 'active'),
  ('01M1BVLQ3J5CRS_SCAN', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'Scanner Detection', 'owasp-crs', true, 'low', 'active')
ON CONFLICT DO NOTHING;
