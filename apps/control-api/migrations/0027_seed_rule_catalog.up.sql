-- 0027_seed_rule_catalog.up.sql
-- Full 27-category WAF rule catalog with F5-style signature metadata.

-- Managed rule groups (one per category).
INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT001', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'SQL Injection', 'sqli', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT002', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'Cross-Site Scripting', 'xss', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT003', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'OS Command Injection', 'cmdi', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT004', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'Remote Code Execution', 'rce', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT005', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'Path Traversal', 'pt', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT006', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'File Inclusion', 'fi', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT007', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'XXE / XML', 'xxe', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT008', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'SSRF', 'ssrf', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT009', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'HTTP Protocol', 'http', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT010', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'HTTP Parameter Pollution', 'hpp', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT011', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'API Security', 'api', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT012', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'File Upload', 'file_upload', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT013', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'LDAP Injection', 'ldap', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT014', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'NoSQL Injection', 'nosql', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT015', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'Server-Side Template Injection', 'ssti', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT016', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'Deserialization', 'deserialization', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT017', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'Scanner Detection', 'scanner', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT018', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'Bot Protection', 'bot', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT019', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'Authentication Attacks', 'auth', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT020', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'Session Security', 'session', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT021', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'CSRF', 'csrf', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT022', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'Information Disclosure', 'info_disclosure', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT023', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'Resource Discovery', 'resource_discovery', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT024', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'Request Anomaly', 'request_anomaly', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT025', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'IP / Reputation', 'ip', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT026', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'Geo Security', 'geo', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT027', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'Rate Limiting', 'rate_limit', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

-- Signature rules with metadata + conditions.
INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0001', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SQLI-001', 'SQL Injection', 'SQL Injection - SQL Injection signature', 'managed', 'sqli', 'critical', 101, 'block', 'active', 'AND', 1, 'SQL Injection', 'regex', 98, 'critical', 95, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0001C', '01MRRUL0001', 0, 'request_body', 'regex', '(?i)(union[\s]+select|select[\s]+.*[\s]+from|insert[\s]+into|delete[\s]+from)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0002', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SQLI-002', 'Blind SQL Injection', 'Blind SQL Injection - SQL Injection signature', 'managed', 'sqli', 'critical', 102, 'block', 'active', 'AND', 1, 'Blind SQL Injection', 'regex', 92, 'critical', 90, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0002C', '01MRRUL0002', 0, 'request_body', 'regex', '(?i)(and|or)[\s]+[0-9]+[=][0-9]+', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0003', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SQLI-003', 'Boolean-based SQL Injection', 'Boolean-based SQL Injection - SQL Injection signature', 'managed', 'sqli', 'critical', 103, 'block', 'active', 'AND', 1, 'Boolean-based SQL Injection', 'regex', 90, 'critical', 88, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0003C', '01MRRUL0003', 0, 'request_body', 'regex', '(?i)(''|\))[\s]*(and|or)[\s]+(''[^'']*''|[0-9]+)[\s]*=[\s]*(''[^'']*''|[0-9]+)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0004', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SQLI-004', 'Time-based SQL Injection', 'Time-based SQL Injection - SQL Injection signature', 'managed', 'sqli', 'critical', 104, 'block', 'active', 'AND', 1, 'Time-based SQL Injection', 'regex', 94, 'critical', 92, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0004C', '01MRRUL0004', 0, 'request_body', 'regex', '(?i)(sleep|benchmark|pg_sleep|waitfor[\s]+delay)[\s]*\(', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0005', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SQLI-005', 'Union-based SQL Injection', 'Union-based SQL Injection - SQL Injection signature', 'managed', 'sqli', 'critical', 105, 'block', 'active', 'AND', 1, 'Union-based SQL Injection', 'regex', 96, 'critical', 93, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0005C', '01MRRUL0005', 0, 'request_body', 'regex', '(?i)union[\s]+(all[\s]+)?select', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0006', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SQLI-006', 'Error-based SQL Injection', 'Error-based SQL Injection - SQL Injection signature', 'managed', 'sqli', 'critical', 106, 'block', 'active', 'AND', 1, 'Error-based SQL Injection', 'regex', 90, 'critical', 87, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0006C', '01MRRUL0006', 0, 'request_body', 'regex', '(?i)(extractvalue|updatexml|group[\s]+by[\s]+with[\s]+rollup)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0007', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SQLI-007', 'Stored SQL Injection', 'Stored SQL Injection - SQL Injection signature', 'managed', 'sqli', 'critical', 107, 'block', 'active', 'AND', 1, 'Stored SQL Injection', 'regex', 89, 'critical', 85, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0007C', '01MRRUL0007', 0, 'request_body', 'regex', '(?i)(into[\s]+outfile|load_file|into[\s]+dumpfile)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0008', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SQLI-008', 'Numeric SQL Injection', 'Numeric SQL Injection - SQL Injection signature', 'managed', 'sqli', 'critical', 108, 'block', 'active', 'AND', 1, 'Numeric SQL Injection', 'regex', 82, 'critical', 80, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0008C', '01MRRUL0008', 0, 'request_body', 'regex', '(?i)[0-9]+[\s]*[;''\"-]+', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0009', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SQLI-009', 'String SQL Injection', 'String SQL Injection - SQL Injection signature', 'managed', 'sqli', 'critical', 109, 'block', 'active', 'AND', 1, 'String SQL Injection', 'regex', 88, 'critical', 86, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0009C', '01MRRUL0009', 0, 'request_body', 'regex', '(?i)(''|\")[\s]*(or|and)[\s]+''', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0010', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SQLI-010', 'SQL Comment Injection', 'SQL Comment Injection - SQL Injection signature', 'managed', 'sqli', 'critical', 110, 'block', 'active', 'AND', 1, 'SQL Comment Injection', 'regex', 85, 'critical', 82, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0010C', '01MRRUL0010', 0, 'request_body', 'regex', '(?i)(--|#|\/\*.*\*\/)[\s]*', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0011', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SQLI-011', 'SQL Keyword Abuse', 'SQL Keyword Abuse - SQL Injection signature', 'managed', 'sqli', 'critical', 111, 'block', 'active', 'AND', 1, 'SQL Keyword Abuse', 'regex', 84, 'critical', 81, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0011C', '01MRRUL0011', 0, 'request_body', 'regex', '(?i)\b(having|cast|convert|declare|exec|execute|master|xp_)\b', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0012', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SQLI-012', 'SQL Function Abuse', 'SQL Function Abuse - SQL Injection signature', 'managed', 'sqli', 'critical', 112, 'block', 'active', 'AND', 1, 'SQL Function Abuse', 'regex', 83, 'critical', 80, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0012C', '01MRRUL0012', 0, 'request_body', 'regex', '(?i)\b(ascii|char|chr|concat|substring|substr|length|hex|ord)\b[\s]*\(', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0013', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SQLI-013', 'Database Error Pattern', 'Database Error Pattern - SQL Injection signature', 'managed', 'sqli', 'critical', 113, 'block', 'active', 'AND', 1, 'Database Error Pattern', 'regex', 80, 'critical', 78, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0013C', '01MRRUL0013', 0, 'request_body', 'regex', '(?i)(mysql|postgresql|oracle|sqlite|syntax error|sqlstate)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0014', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SQLI-014', 'Database Fingerprinting', 'Database Fingerprinting - SQL Injection signature', 'managed', 'sqli', 'critical', 114, 'block', 'active', 'AND', 1, 'Database Fingerprinting', 'regex', 88, 'critical', 85, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0014C', '01MRRUL0014', 0, 'request_body', 'regex', '(?i)(version\(\)|@@version|dbms_random|sqlite_version)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0015', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SQLI-015', 'SQL Meta-character Detection', 'SQL Meta-character Detection - SQL Injection signature', 'managed', 'sqli', 'critical', 115, 'block', 'active', 'AND', 1, 'SQL Meta-character Detection', 'regex', 90, 'critical', 87, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0015C', '01MRRUL0015', 0, 'request_body', 'regex', '(?i)(;|\|\||&&)[\s]*(select|drop|insert|update|delete|alter)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0016', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-XSS-001', 'Reflected XSS', 'Reflected XSS - Cross-Site Scripting signature', 'managed', 'xss', 'high', 201, 'block', 'active', 'AND', 1, 'Reflected XSS', 'regex', 96, 'high', 93, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0016C', '01MRRUL0016', 0, 'request_body', 'regex', '(?i)<script[\s>]', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0017', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-XSS-002', 'Stored XSS', 'Stored XSS - Cross-Site Scripting signature', 'managed', 'xss', 'high', 202, 'block', 'active', 'AND', 1, 'Stored XSS', 'regex', 92, 'high', 90, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0017C', '01MRRUL0017', 0, 'request_body', 'regex', '(?i)<script[\s>].*</script>', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0018', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-XSS-003', 'DOM XSS', 'DOM XSS - Cross-Site Scripting signature', 'managed', 'xss', 'high', 203, 'block', 'active', 'AND', 1, 'DOM XSS', 'regex', 88, 'high', 85, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0018C', '01MRRUL0018', 0, 'request_body', 'regex', '(?i)(document\.(location|cookie|write)|window\.location|eval\(|innerHTML)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0019', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-XSS-004', 'HTML Injection', 'HTML Injection - Cross-Site Scripting signature', 'managed', 'xss', 'high', 204, 'block', 'active', 'AND', 1, 'HTML Injection', 'regex', 90, 'high', 87, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0019C', '01MRRUL0019', 0, 'request_body', 'regex', '(?i)<(iframe|object|embed|applet|form|img|svg)[\s>]', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0020', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-XSS-005', 'Script Context Injection', 'Script Context Injection - Cross-Site Scripting signature', 'managed', 'xss', 'high', 205, 'block', 'active', 'AND', 1, 'Script Context Injection', 'regex', 93, 'high', 90, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0020C', '01MRRUL0020', 0, 'request_body', 'regex', '(?i)(javascript:|vbscript:|data:text/html)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0021', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-XSS-006', 'HTML Attribute Injection', 'HTML Attribute Injection - Cross-Site Scripting signature', 'managed', 'xss', 'high', 206, 'block', 'active', 'AND', 1, 'HTML Attribute Injection', 'regex', 89, 'high', 86, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0021C', '01MRRUL0021', 0, 'request_body', 'regex', '(?i)[\"'']\s*(on\w+)\s*=', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0022', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-XSS-007', 'JavaScript Context Injection', 'JavaScript Context Injection - Cross-Site Scripting signature', 'managed', 'xss', 'high', 207, 'block', 'active', 'AND', 1, 'JavaScript Context Injection', 'regex', 86, 'high', 83, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0022C', '01MRRUL0022', 0, 'request_body', 'regex', '(?i)(alert|prompt|confirm)\s*\(', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0023', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-XSS-008', 'CSS Context Injection', 'CSS Context Injection - Cross-Site Scripting signature', 'managed', 'xss', 'high', 208, 'block', 'active', 'AND', 1, 'CSS Context Injection', 'regex', 82, 'high', 79, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0023C', '01MRRUL0023', 0, 'request_body', 'regex', '(?i)(expression|@import|url\()', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0024', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-XSS-009', 'Event Handler Injection', 'Event Handler Injection - Cross-Site Scripting signature', 'managed', 'xss', 'high', 209, 'block', 'active', 'AND', 1, 'Event Handler Injection', 'regex', 91, 'high', 88, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0024C', '01MRRUL0024', 0, 'request_body', 'regex', '(?i)\s(on(load|error|click|mouseover|focus))\s*=', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0025', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-XSS-010', 'Encoded XSS', 'Encoded XSS - Cross-Site Scripting signature', 'managed', 'xss', 'high', 210, 'block', 'active', 'AND', 1, 'Encoded XSS', 'regex', 90, 'high', 87, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0025C', '01MRRUL0025', 0, 'request_body', 'regex', '(?i)(%3cscript|%253c|&lt;script|&#60;script)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0026', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-XSS-011', 'Polyglot XSS', 'Polyglot XSS - Cross-Site Scripting signature', 'managed', 'xss', 'high', 211, 'block', 'active', 'AND', 1, 'Polyglot XSS', 'regex', 85, 'high', 82, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0026C', '01MRRUL0026', 0, 'request_body', 'regex', '(?i)(\"''<>/\\|javascript|onerror)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0027', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-CMDI-001', 'OS Command Injection', 'OS Command Injection - OS Command Injection signature', 'managed', 'cmdi', 'critical', 301, 'block', 'active', 'AND', 1, 'OS Command Injection', 'regex', 95, 'critical', 92, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0027C', '01MRRUL0027', 0, 'request_body', 'regex', '(?i)(;|&&|\|\|)[\s]*(ls|cat|id|whoami|pwd|uname|dir|type)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0028', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-CMDI-002', 'Shell Command Injection', 'Shell Command Injection - OS Command Injection signature', 'managed', 'cmdi', 'critical', 302, 'block', 'active', 'AND', 1, 'Shell Command Injection', 'regex', 94, 'critical', 91, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0028C', '01MRRUL0028', 0, 'request_body', 'regex', '(?i)(/bin/sh|/bin/bash|cmd\.exe|powershell)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0029', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-CMDI-003', 'Unix Command Injection', 'Unix Command Injection - OS Command Injection signature', 'managed', 'cmdi', 'critical', 303, 'block', 'active', 'AND', 1, 'Unix Command Injection', 'regex', 92, 'critical', 89, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0029C', '01MRRUL0029', 0, 'request_body', 'regex', '(?i)(\$\(|`[^`]*`|;[\s]*(sh|bash|perl|python))', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0030', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-CMDI-004', 'Windows Command Injection', 'Windows Command Injection - OS Command Injection signature', 'managed', 'cmdi', 'critical', 304, 'block', 'active', 'AND', 1, 'Windows Command Injection', 'regex', 90, 'critical', 87, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0030C', '01MRRUL0030', 0, 'request_body', 'regex', '(?i)(cmd[/\\\\]c|type[\s]+c:\\\\|dir[\s]+c:\\\\|net[\s]+user)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0031', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-CMDI-005', 'Command Separator Detection', 'Command Separator Detection - OS Command Injection signature', 'managed', 'cmdi', 'critical', 305, 'block', 'active', 'AND', 1, 'Command Separator Detection', 'regex', 84, 'critical', 80, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0031C', '01MRRUL0031', 0, 'request_body', 'regex', '(?i)(;|\|\||&&|\n|%0a|%0d)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0032', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-CMDI-006', 'Command Substitution', 'Command Substitution - OS Command Injection signature', 'managed', 'cmdi', 'critical', 306, 'block', 'active', 'AND', 1, 'Command Substitution', 'regex', 89, 'critical', 86, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0032C', '01MRRUL0032', 0, 'request_body', 'regex', '(?i)(\$\([^)]*\)|`[^`]*`)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0033', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-CMDI-007', 'Shell Metacharacter Detection', 'Shell Metacharacter Detection - OS Command Injection signature', 'managed', 'cmdi', 'critical', 307, 'block', 'active', 'AND', 1, 'Shell Metacharacter Detection', 'regex', 78, 'critical', 75, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0033C', '01MRRUL0033', 0, 'request_body', 'regex', '(?i)([;&|`$()<>])', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0034', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-CMDI-008', 'Command Execution Pattern', 'Command Execution Pattern - OS Command Injection signature', 'managed', 'cmdi', 'critical', 308, 'block', 'active', 'AND', 1, 'Command Execution Pattern', 'regex', 87, 'critical', 84, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0034C', '01MRRUL0034', 0, 'request_body', 'regex', '(?i)\b(wget|curl|nc|netcat|tftp|scp|rsync)\b', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0035', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-RCE-001', 'Code Injection', 'Code Injection - Remote Code Execution signature', 'managed', 'rce', 'critical', 401, 'block', 'active', 'AND', 1, 'Code Injection', 'regex', 96, 'critical', 93, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0035C', '01MRRUL0035', 0, 'request_body', 'regex', '(?i)\b(eval|assert|system|passthru|shell_exec|exec)\s*\(', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0036', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-RCE-002', 'Expression Injection', 'Expression Injection - Remote Code Execution signature', 'managed', 'rce', 'critical', 402, 'block', 'active', 'AND', 1, 'Expression Injection', 'regex', 88, 'critical', 85, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0036C', '01MRRUL0036', 0, 'request_body', 'regex', '(?i)(SpEL|OGNL|EL|#\{|\\$\{)[\s\w]', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0037', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-RCE-003', 'Template Injection', 'Template Injection - Remote Code Execution signature', 'managed', 'rce', 'critical', 403, 'block', 'active', 'AND', 1, 'Template Injection', 'regex', 90, 'critical', 87, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0037C', '01MRRUL0037', 0, 'request_body', 'regex', '(?i)(\{\{|\{%|#{|<\%)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0038', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-RCE-004', 'Server-side Template Injection', 'Server-side Template Injection - Remote Code Execution signature', 'managed', 'rce', 'critical', 404, 'block', 'active', 'AND', 1, 'Server-side Template Injection', 'regex', 89, 'critical', 86, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0038C', '01MRRUL0038', 0, 'request_body', 'regex', '(?i)({{[^}]+}}|\${{[^}]+}})', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0039', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-RCE-005', 'Dynamic Code Execution', 'Dynamic Code Execution - Remote Code Execution signature', 'managed', 'rce', 'critical', 405, 'block', 'active', 'AND', 1, 'Dynamic Code Execution', 'regex', 93, 'critical', 90, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0039C', '01MRRUL0039', 0, 'request_body', 'regex', '(?i)(__import__|os\.system|Runtime\.getRuntime|ProcessBuilder)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0040', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-RCE-006', 'Script Injection', 'Script Injection - Remote Code Execution signature', 'managed', 'rce', 'critical', 406, 'block', 'active', 'AND', 1, 'Script Injection', 'regex', 85, 'critical', 82, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0040C', '01MRRUL0040', 0, 'request_body', 'regex', '(?i)(python|cscript|wscript|perl|php[\s-]+r)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0041', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-RCE-007', 'Runtime Execution Pattern', 'Runtime Execution Pattern - Remote Code Execution signature', 'managed', 'rce', 'critical', 407, 'block', 'active', 'AND', 1, 'Runtime Execution Pattern', 'regex', 84, 'critical', 81, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0041C', '01MRRUL0041', 0, 'request_body', 'regex', '(?i)(<%=|<\?php|jsp:|\.cfm)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0042', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-PT-001', 'Directory Traversal', 'Directory Traversal - Path Traversal signature', 'managed', 'pt', 'high', 501, 'block', 'active', 'AND', 1, 'Directory Traversal', 'regex', 97, 'high', 94, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0042C', '01MRRUL0042', 0, 'url', 'regex', '(?i)(\.\./|\.\.\\|%2e%2e)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0043', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-PT-002', 'Encoded Traversal', 'Encoded Traversal - Path Traversal signature', 'managed', 'pt', 'high', 502, 'block', 'active', 'AND', 1, 'Encoded Traversal', 'regex', 93, 'high', 90, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0043C', '01MRRUL0043', 0, 'url', 'regex', '(?i)(%2e%2e%2f|%252e%252e|%c0%ae%c0%ae)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0044', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-PT-003', 'Double Encoded Traversal', 'Double Encoded Traversal - Path Traversal signature', 'managed', 'pt', 'high', 503, 'block', 'active', 'AND', 1, 'Double Encoded Traversal', 'regex', 88, 'high', 85, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0044C', '01MRRUL0044', 0, 'url', 'regex', '(?i)(%25252e|%252e%252e)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0045', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-PT-004', 'Path Normalization Bypass', 'Path Normalization Bypass - Path Traversal signature', 'managed', 'pt', 'high', 504, 'block', 'active', 'AND', 1, 'Path Normalization Bypass', 'regex', 80, 'high', 76, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0045C', '01MRRUL0045', 0, 'url', 'regex', '(?i)(\.\./\.\./|/\./|//)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0046', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-PT-005', 'Unix Path Traversal', 'Unix Path Traversal - Path Traversal signature', 'managed', 'pt', 'high', 505, 'block', 'active', 'AND', 1, 'Unix Path Traversal', 'regex', 96, 'high', 93, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0046C', '01MRRUL0046', 0, 'url', 'regex', '(?i)(etc/passwd|etc/shadow|proc/self|home/[a-z])', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0047', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-PT-006', 'Windows Path Traversal', 'Windows Path Traversal - Path Traversal signature', 'managed', 'pt', 'high', 506, 'block', 'active', 'AND', 1, 'Windows Path Traversal', 'regex', 92, 'high', 89, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0047C', '01MRRUL0047', 0, 'url', 'regex', '(?i)(c:\\\\windows|boot\.ini|win\.ini)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0048', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-PT-007', 'Sensitive File Access', 'Sensitive File Access - Path Traversal signature', 'managed', 'pt', 'high', 507, 'block', 'active', 'AND', 1, 'Sensitive File Access', 'regex', 94, 'high', 91, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0048C', '01MRRUL0048', 0, 'url', 'regex', '(?i)(\.env|\.git|web\.config|php\.ini|\.htaccess)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0049', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-FI-001', 'Local File Inclusion', 'Local File Inclusion - File Inclusion signature', 'managed', 'fi', 'high', 601, 'block', 'active', 'AND', 1, 'Local File Inclusion', 'regex', 95, 'high', 92, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0049C', '01MRRUL0049', 0, 'url', 'regex', '(?i)(include|require)[\s=]+(file|php://|/etc|/var/www)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0050', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-FI-002', 'Remote File Inclusion', 'Remote File Inclusion - File Inclusion signature', 'managed', 'fi', 'high', 602, 'block', 'active', 'AND', 1, 'Remote File Inclusion', 'regex', 96, 'high', 93, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0050C', '01MRRUL0050', 0, 'url', 'regex', '(?i)(include[\s=]+(http|https|ftp):)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0051', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-FI-003', 'PHP File Inclusion', 'PHP File Inclusion - File Inclusion signature', 'managed', 'fi', 'high', 603, 'block', 'active', 'AND', 1, 'PHP File Inclusion', 'regex', 94, 'high', 91, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0051C', '01MRRUL0051', 0, 'url', 'regex', '(?i)(php://(filter|input|data|expect))', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0052', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-FI-004', 'Dynamic Include Abuse', 'Dynamic Include Abuse - File Inclusion signature', 'managed', 'fi', 'high', 604, 'block', 'active', 'AND', 1, 'Dynamic Include Abuse', 'regex', 89, 'high', 86, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0052C', '01MRRUL0052', 0, 'url', 'regex', '(?i)(page=|\?file=|view=|path=).*(\.\.|/etc|php://)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0053', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-FI-005', 'Remote Resource Inclusion', 'Remote Resource Inclusion - File Inclusion signature', 'managed', 'fi', 'high', 605, 'block', 'active', 'AND', 1, 'Remote Resource Inclusion', 'regex', 85, 'high', 82, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0053C', '01MRRUL0053', 0, 'url', 'regex', '(?i)(http[s]?://[^\s]+(\.txt|\.php|\.html))', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0054', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-XXE-001', 'XML External Entity', 'XML External Entity - XXE / XML signature', 'managed', 'xxe', 'critical', 701, 'block', 'active', 'AND', 1, 'XML External Entity', 'regex', 97, 'critical', 94, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0054C', '01MRRUL0054', 0, 'request_body', 'regex', '(?i)<!DOCTYPE[\s]+[^>]*\[[^>]*<!ENTITY', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0055', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-XXE-002', 'External Entity Reference', 'External Entity Reference - XXE / XML signature', 'managed', 'xxe', 'critical', 702, 'block', 'active', 'AND', 1, 'External Entity Reference', 'regex', 93, 'critical', 90, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0055C', '01MRRUL0055', 0, 'request_body', 'regex', '(?i)(&xxe;|&ent;|SYSTEM[\s]+\"|PUBLIC[\s]+\")', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0056', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-XXE-003', 'Parameter Entity', 'Parameter Entity - XXE / XML signature', 'managed', 'xxe', 'critical', 703, 'block', 'active', 'AND', 1, 'Parameter Entity', 'regex', 91, 'critical', 88, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0056C', '01MRRUL0056', 0, 'request_body', 'regex', '(?i)<!ENTITY[\s]+%[^>]*>', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0057', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-XXE-004', 'External DTD', 'External DTD - XXE / XML signature', 'managed', 'xxe', 'critical', 704, 'block', 'active', 'AND', 1, 'External DTD', 'regex', 92, 'critical', 89, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0057C', '01MRRUL0057', 0, 'request_body', 'regex', '(?i)(<!DOCTYPE[\s]+[^>]*SYSTEM[\s]+\")', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0058', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-XXE-005', 'XML Entity Expansion', 'XML Entity Expansion - XXE / XML signature', 'managed', 'xxe', 'critical', 705, 'block', 'active', 'AND', 1, 'XML Entity Expansion', 'regex', 90, 'critical', 87, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0058C', '01MRRUL0058', 0, 'request_body', 'regex', '(?i)<!ENTITY[\s]+[a-z][^>]*>.*(<!ENTITY[\s]+)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0059', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-XXE-006', 'XML Bomb', 'XML Bomb - XXE / XML signature', 'managed', 'xxe', 'critical', 706, 'block', 'active', 'AND', 1, 'XML Bomb', 'regex', 89, 'critical', 86, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0059C', '01MRRUL0059', 0, 'request_body', 'regex', '(?i)(lol|billion)[^\"]*\"[\s]*\"[^\"]*\"[\s]*\"[^\"]*\"', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0060', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-XXE-007', 'XML Parser Abuse', 'XML Parser Abuse - XXE / XML signature', 'managed', 'xxe', 'critical', 707, 'block', 'active', 'AND', 1, 'XML Parser Abuse', 'regex', 85, 'critical', 82, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0060C', '01MRRUL0060', 0, 'request_body', 'regex', '(?i)(<\?(xml|soap)|<![CDATA[|&lt;)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0061', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-XXE-008', 'Malformed XML', 'Malformed XML - XXE / XML signature', 'managed', 'xxe', 'critical', 708, 'block', 'active', 'AND', 1, 'Malformed XML', 'regex', 80, 'critical', 76, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0061C', '01MRRUL0061', 0, 'request_body', 'regex', '(?i)(<xml[^>]*>.*<xml|</>|<<)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0062', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SSRF-001', 'Server-Side Request Forgery', 'Server-Side Request Forgery - SSRF signature', 'managed', 'ssrf', 'critical', 801, 'block', 'active', 'AND', 1, 'Server-Side Request Forgery', 'regex', 88, 'critical', 84, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0062C', '01MRRUL0062', 0, 'url', 'regex', '(?i)(url=|uri=|dest=|redirect=|proxy=)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0063', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SSRF-002', 'Internal Network Access', 'Internal Network Access - SSRF signature', 'managed', 'ssrf', 'critical', 802, 'block', 'active', 'AND', 1, 'Internal Network Access', 'regex', 94, 'critical', 91, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0063C', '01MRRUL0063', 0, 'url', 'regex', '(?i)(10\.\d+\.\d+\.\d+|192\.168\.\d+\.\d+|172\.(1[6-9]|2[0-9]|3[01])\.\d+\.\d+)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0064', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SSRF-003', 'Localhost Access', 'Localhost Access - SSRF signature', 'managed', 'ssrf', 'critical', 803, 'block', 'active', 'AND', 1, 'Localhost Access', 'regex', 95, 'critical', 92, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0064C', '01MRRUL0064', 0, 'url', 'regex', '(?i)(localhost|127\.0\.0\.1|0\.0\.0\.0)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0065', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SSRF-004', 'Private IP Access', 'Private IP Access - SSRF signature', 'managed', 'ssrf', 'critical', 804, 'block', 'active', 'AND', 1, 'Private IP Access', 'regex', 93, 'critical', 90, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0065C', '01MRRUL0065', 0, 'url', 'regex', '(?i)(10\.|192\.168\.|172\.(1[6-9]|2[0-9]|3[01])\.)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0066', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SSRF-005', 'Cloud Metadata Access', 'Cloud Metadata Access - SSRF signature', 'managed', 'ssrf', 'critical', 805, 'block', 'active', 'AND', 1, 'Cloud Metadata Access', 'regex', 97, 'critical', 95, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0066C', '01MRRUL0066', 0, 'url', 'regex', '(?i)(169\.254\.169\.254|metadata\.google\.internal|100\.100\.100\.200)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0067', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SSRF-006', 'Loopback Access', 'Loopback Access - SSRF signature', 'managed', 'ssrf', 'critical', 806, 'block', 'active', 'AND', 1, 'Loopback Access', 'regex', 91, 'critical', 88, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0067C', '01MRRUL0067', 0, 'url', 'regex', '(?i)(::1|\[::1\]|127\.[0-9]+\.[0-9]+\.[0-9]+)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0068', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SSRF-007', 'Internal DNS Access', 'Internal DNS Access - SSRF signature', 'managed', 'ssrf', 'critical', 807, 'block', 'active', 'AND', 1, 'Internal DNS Access', 'regex', 87, 'critical', 84, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0068C', '01MRRUL0068', 0, 'url', 'regex', '(?i)(\.internal$|\.local$|\.corp$|\.intranet$)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0069', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SSRF-008', 'URL Scheme Abuse', 'URL Scheme Abuse - SSRF signature', 'managed', 'ssrf', 'critical', 808, 'block', 'active', 'AND', 1, 'URL Scheme Abuse', 'regex', 92, 'critical', 89, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0069C', '01MRRUL0069', 0, 'url', 'regex', '(?i)(gopher://|file://|dict://|ftp://)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0070', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-HTTP-001', 'HTTP Request Smuggling', 'HTTP Request Smuggling - HTTP Protocol signature', 'managed', 'http', 'medium', 901, 'block', 'active', 'AND', 1, 'HTTP Request Smuggling', 'regex', 96, 'medium', 93, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0070C', '01MRRUL0070', 0, 'header', 'regex', '(?i)(transfer-encoding:\s*chunked[^\n]*\ncontent-length:|content-length:[^\n]*\ntransfer-encoding:\s*chunked)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0071', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-HTTP-002', 'HTTP Response Splitting', 'HTTP Response Splitting - HTTP Protocol signature', 'managed', 'http', 'medium', 902, 'block', 'active', 'AND', 1, 'HTTP Response Splitting', 'regex', 90, 'medium', 87, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0071C', '01MRRUL0071', 0, 'header', 'regex', '(?i)(%0d%0a|\r\n\r\n)[^\n]*(location:|set-cookie:)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0072', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-HTTP-003', 'Invalid HTTP Method', 'Invalid HTTP Method - HTTP Protocol signature', 'managed', 'http', 'medium', 903, 'block', 'active', 'AND', 1, 'Invalid HTTP Method', 'regex', 88, 'medium', 85, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0072C', '01MRRUL0072', 0, 'method', 'regex', '(?i)^(trace|track|connect)$', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0073', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-HTTP-004', 'Invalid HTTP Version', 'Invalid HTTP Version - HTTP Protocol signature', 'managed', 'http', 'medium', 904, 'block', 'active', 'AND', 1, 'Invalid HTTP Version', 'regex', 87, 'medium', 84, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0073C', '01MRRUL0073', 0, 'url', 'regex', '(?i)(HTTP/0\.|HTTP/1\.0|HTTP/9)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0074', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-HTTP-005', 'Malformed Request', 'Malformed Request - HTTP Protocol signature', 'managed', 'http', 'medium', 905, 'block', 'active', 'AND', 1, 'Malformed Request', 'regex', 84, 'medium', 80, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0074C', '01MRRUL0074', 0, 'url', 'regex', '(?i)(%00|%0a|%0d|\\\\)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0075', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-HTTP-006', 'Header Anomaly', 'Header Anomaly - HTTP Protocol signature', 'managed', 'http', 'medium', 906, 'block', 'active', 'AND', 1, 'Header Anomaly', 'regex', 82, 'medium', 78, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0075C', '01MRRUL0075', 0, 'header', 'regex', '(?i)(^\s|\\r\\n\\r\\n.*<)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0076', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-HTTP-007', 'Host Header Attack', 'Host Header Attack - HTTP Protocol signature', 'managed', 'http', 'medium', 907, 'block', 'active', 'AND', 1, 'Host Header Attack', 'regex', 89, 'medium', 86, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0076C', '01MRRUL0076', 0, 'host', 'regex', '(?i)(localhost|127\.0\.0\.1|0\.0\.0\.0|\.internal)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0077', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-HTTP-008', 'Transfer-Encoding Anomaly', 'Transfer-Encoding Anomaly - HTTP Protocol signature', 'managed', 'http', 'medium', 908, 'block', 'active', 'AND', 1, 'Transfer-Encoding Anomaly', 'regex', 91, 'medium', 88, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0077C', '01MRRUL0077', 0, 'header', 'regex', '(?i)transfer-encoding:\s*(identity|chunked[\s,])', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0078', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-HTTP-009', 'Content-Length Anomaly', 'Content-Length Anomaly - HTTP Protocol signature', 'managed', 'http', 'medium', 909, 'block', 'active', 'AND', 1, 'Content-Length Anomaly', 'regex', 93, 'medium', 90, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0078C', '01MRRUL0078', 0, 'header', 'regex', '(?i)content-length:\s*[0-9]+[\s,]*content-length:', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0079', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-HTTP-010', 'Duplicate Header', 'Duplicate Header - HTTP Protocol signature', 'managed', 'http', 'medium', 910, 'block', 'active', 'AND', 1, 'Duplicate Header', 'regex', 90, 'medium', 87, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0079C', '01MRRUL0079', 0, 'header', 'regex', '(?i)(content-length:|transfer-encoding:|host:).*\n.*(content-length:|transfer-encoding:|host:)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0080', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-HTTP-011', 'Request Desynchronization', 'Request Desynchronization - HTTP Protocol signature', 'managed', 'http', 'medium', 911, 'block', 'active', 'AND', 1, 'Request Desynchronization', 'regex', 92, 'medium', 89, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0080C', '01MRRUL0080', 0, 'header', 'regex', '(?i)(te:\s*chunked|cl\.te|te\.cl)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0081', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-HPP-001', 'Duplicate Parameters', 'Duplicate Parameters - HTTP Parameter Pollution signature', 'managed', 'hpp', 'medium', 1001, 'block', 'active', 'AND', 1, 'Duplicate Parameters', 'regex', 85, 'medium', 81, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0081C', '01MRRUL0081', 0, 'url', 'regex', '(?i)([a-z_]+=[^&]+&[a-z_]+=)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0082', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-HPP-002', 'Conflicting Parameters', 'Conflicting Parameters - HTTP Parameter Pollution signature', 'managed', 'hpp', 'medium', 1002, 'block', 'active', 'AND', 1, 'Conflicting Parameters', 'regex', 88, 'medium', 85, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0082C', '01MRRUL0082', 0, 'url', 'regex', '(?i)(role|admin|user|account)=[^&]+&(role|admin|user|account)=', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0083', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-HPP-003', 'Parameter Pollution', 'Parameter Pollution - HTTP Parameter Pollution signature', 'managed', 'hpp', 'medium', 1003, 'block', 'active', 'AND', 1, 'Parameter Pollution', 'regex', 83, 'medium', 80, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0083C', '01MRRUL0083', 0, 'url', 'regex', '(?i)(&|;|\|)[^=]*=[^&]*(&|\|;)[^=]*=', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0084', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-HPP-004', 'Array Parameter Abuse', 'Array Parameter Abuse - HTTP Parameter Pollution signature', 'managed', 'hpp', 'medium', 1004, 'block', 'active', 'AND', 1, 'Array Parameter Abuse', 'regex', 86, 'medium', 83, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0084C', '01MRRUL0084', 0, 'url', 'regex', '(?i)([a-z_]+\[\]=|\[\]\[)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0085', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-HPP-005', 'Parameter Parsing Ambiguity', 'Parameter Parsing Ambiguity - HTTP Parameter Pollution signature', 'managed', 'hpp', 'medium', 1005, 'block', 'active', 'AND', 1, 'Parameter Parsing Ambiguity', 'regex', 84, 'medium', 80, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0085C', '01MRRUL0085', 0, 'url', 'regex', '(?i)(%26|%3b|%7c|%3d)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0086', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-API-001', 'Invalid Content-Type', 'Invalid Content-Type - API Security signature', 'managed', 'api', 'medium', 1101, 'block', 'active', 'AND', 1, 'Invalid Content-Type', 'regex', 82, 'medium', 78, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0086C', '01MRRUL0086', 0, 'header', 'regex', '(?i)content-type:\s*(text/plain|application/octet-stream)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0087', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-API-002', 'JSON Structure Violation', 'JSON Structure Violation - API Security signature', 'managed', 'api', 'medium', 1102, 'block', 'active', 'AND', 1, 'JSON Structure Violation', 'regex', 80, 'medium', 76, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0087C', '01MRRUL0087', 0, 'request_body', 'regex', '(?i)(\{\{|\}\}|:{2}|,{2})', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0088', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-API-003', 'Parameter Type Violation', 'Parameter Type Violation - API Security signature', 'managed', 'api', 'medium', 1103, 'block', 'active', 'AND', 1, 'Parameter Type Violation', 'regex', 81, 'medium', 78, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0088C', '01MRRUL0088', 0, 'request_body', 'regex', '(?i)(\"(age|count|limit|offset)\":\s*\")', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0089', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-API-004', 'Excessive Parameters', 'Excessive Parameters - API Security signature', 'managed', 'api', 'medium', 1104, 'block', 'active', 'AND', 1, 'Excessive Parameters', 'regex', 79, 'medium', 75, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0089C', '01MRRUL0089', 0, 'request_body', 'regex', '(?i)(\"[a-z_]+\":){10,}', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0090', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-API-005', 'Unknown Parameter', 'Unknown Parameter - API Security signature', 'managed', 'api', 'medium', 1105, 'block', 'active', 'AND', 1, 'Unknown Parameter', 'regex', 90, 'medium', 87, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0090C', '01MRRUL0090', 0, 'request_body', 'regex', '(?i)(\"(__proto__|constructor|prototype)\":)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0091', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-API-006', 'API Schema Violation', 'API Schema Violation - API Security signature', 'managed', 'api', 'medium', 1106, 'block', 'active', 'AND', 1, 'API Schema Violation', 'regex', 85, 'medium', 82, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0091C', '01MRRUL0091', 0, 'request_body', 'regex', '(?i)(\"type\":\s*\"(int|str|bool)\"[\s,]*\"(expected|enum)\")', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0092', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-FILE_UPLOAD-001', 'Executable File Upload', 'Executable File Upload - File Upload signature', 'managed', 'file_upload', 'high', 1201, 'block', 'active', 'AND', 1, 'Executable File Upload', 'regex', 92, 'high', 89, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0092C', '01MRRUL0092', 0, 'url', 'regex', '(?i)(\.exe|\.dll|\.so|\.bin|\.msi)$', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0093', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-FILE_UPLOAD-002', 'Script File Upload', 'Script File Upload - File Upload signature', 'managed', 'file_upload', 'high', 1202, 'block', 'active', 'AND', 1, 'Script File Upload', 'regex', 95, 'high', 92, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0093C', '01MRRUL0093', 0, 'url', 'regex', '(?i)(\.(php|php3|php5|phtml|jsp|asp|aspx|py|pl|rb|sh|cgi))$', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0094', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-FILE_UPLOAD-003', 'MIME Type Mismatch', 'MIME Type Mismatch - File Upload signature', 'managed', 'file_upload', 'high', 1203, 'block', 'active', 'AND', 1, 'MIME Type Mismatch', 'regex', 88, 'high', 85, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0094C', '01MRRUL0094', 0, 'header', 'regex', '(?i)(multipart/form-data.*filename=\"[^\"]+\.(php|exe|sh))', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0095', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-FILE_UPLOAD-004', 'Double Extension', 'Double Extension - File Upload signature', 'managed', 'file_upload', 'high', 1204, 'block', 'active', 'AND', 1, 'Double Extension', 'regex', 90, 'high', 87, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0095C', '01MRRUL0095', 0, 'url', 'regex', '(?i)(\.(php|exe|sh)\.[a-z0-9]{2,4})$', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0096', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-FILE_UPLOAD-005', 'Archive Abuse', 'Archive Abuse - File Upload signature', 'managed', 'file_upload', 'high', 1205, 'block', 'active', 'AND', 1, 'Archive Abuse', 'regex', 84, 'high', 80, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0096C', '01MRRUL0096', 0, 'url', 'regex', '(?i)(\.(zip|tar|gz|rar|7z))$', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0097', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-FILE_UPLOAD-006', 'Oversized File', 'Oversized File - File Upload signature', 'managed', 'file_upload', 'high', 1206, 'block', 'active', 'AND', 1, 'Oversized File', 'regex', 86, 'high', 82, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0097C', '01MRRUL0097', 0, 'header', 'regex', '(?i)content-length:\s*[0-9]{8,}', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0098', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-FILE_UPLOAD-007', 'Malicious File Content', 'Malicious File Content - File Upload signature', 'managed', 'file_upload', 'high', 1207, 'block', 'active', 'AND', 1, 'Malicious File Content', 'regex', 87, 'high', 84, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0098C', '01MRRUL0098', 0, 'request_body', 'regex', '(?i)(GIF89a|\\x89PNG|\\xff\\xd8\\xff)(.{0,100})(<\?php|exec\(|eval\()', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0099', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-FILE_UPLOAD-008', 'Upload Path Manipulation', 'Upload Path Manipulation - File Upload signature', 'managed', 'file_upload', 'high', 1208, 'block', 'active', 'AND', 1, 'Upload Path Manipulation', 'regex', 89, 'high', 86, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0099C', '01MRRUL0099', 0, 'url', 'regex', '(?i)((\.\./|%2e%2e/).*(upload|files|media))', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0100', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-LDAP-001', 'LDAP Filter Injection', 'LDAP Filter Injection - LDAP Injection signature', 'managed', 'ldap', 'high', 1301, 'block', 'active', 'AND', 1, 'LDAP Filter Injection', 'regex', 91, 'high', 88, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0100C', '01MRRUL0100', 0, 'request_body', 'regex', '(?i)(\*\)\(|\|\(|&\(|!\(|\(\|)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0101', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-LDAP-002', 'LDAP Query Manipulation', 'LDAP Query Manipulation - LDAP Injection signature', 'managed', 'ldap', 'high', 1302, 'block', 'active', 'AND', 1, 'LDAP Query Manipulation', 'regex', 92, 'high', 89, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0101C', '01MRRUL0101', 0, 'request_body', 'regex', '(?i)(cn=|uid=|mail=|userid=|sAMAccountName=).*(\*|\(|\))', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0102', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-LDAP-003', 'LDAP Wildcard Abuse', 'LDAP Wildcard Abuse - LDAP Injection signature', 'managed', 'ldap', 'high', 1303, 'block', 'active', 'AND', 1, 'LDAP Wildcard Abuse', 'regex', 89, 'high', 86, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0102C', '01MRRUL0102', 0, 'request_body', 'regex', '(?i)(\*\)\s*\(|\)\s*\|\s*\(|&\s*\()', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0103', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-LDAP-004', 'LDAP Authentication Bypass', 'LDAP Authentication Bypass - LDAP Injection signature', 'managed', 'ldap', 'high', 1304, 'block', 'active', 'AND', 1, 'LDAP Authentication Bypass Pattern', 'regex', 93, 'high', 90, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0103C', '01MRRUL0103', 0, 'request_body', 'regex', '(?i)(uid=.*\)\s*\(\|\(\s*uid=|\*\)\s*\|\(objectClass=\*)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0104', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-NOSQL-001', 'MongoDB Injection', 'MongoDB Injection - NoSQL Injection signature', 'managed', 'nosql', 'high', 1401, 'block', 'active', 'AND', 1, 'MongoDB Injection', 'regex', 95, 'high', 92, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0104C', '01MRRUL0104', 0, 'request_body', 'regex', '(?i)(\$where|\$ne|\$gt|\$regex|\$nin)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0105', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-NOSQL-002', 'NoSQL Operator Injection', 'NoSQL Operator Injection - NoSQL Injection signature', 'managed', 'nosql', 'high', 1402, 'block', 'active', 'AND', 1, 'NoSQL Operator Injection', 'regex', 93, 'high', 90, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0105C', '01MRRUL0105', 0, 'request_body', 'regex', '(?i)(\"\$[a-z]+\":)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0106', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-NOSQL-003', 'JSON Operator Abuse', 'JSON Operator Abuse - NoSQL Injection signature', 'managed', 'nosql', 'high', 1403, 'block', 'active', 'AND', 1, 'JSON Operator Abuse', 'regex', 92, 'high', 89, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0106C', '01MRRUL0106', 0, 'request_body', 'regex', '(?i)((\$where|\$ne|\$gt|\$lt|\$regex)\s*:\s*\")', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0107', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-NOSQL-004', 'Query Manipulation', 'Query Manipulation - NoSQL Injection signature', 'managed', 'nosql', 'high', 1404, 'block', 'active', 'AND', 1, 'Query Manipulation', 'regex', 88, 'high', 85, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0107C', '01MRRUL0107', 0, 'request_body', 'regex', '(?i)(find\s*\(\s*\{|\"query\":\s*\{)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0108', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-NOSQL-005', 'NoSQL Authentication Bypass', 'NoSQL Authentication Bypass - NoSQL Injection signature', 'managed', 'nosql', 'high', 1405, 'block', 'active', 'AND', 1, 'NoSQL Authentication Bypass', 'regex', 94, 'high', 91, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0108C', '01MRRUL0108', 0, 'request_body', 'regex', '(?i)((\"password\"|\"user\").*(\$ne|\$gt|\$regex))', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0109', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SSTI-001', 'Template Expression', 'Template Expression - Server-Side Template Injection signature', 'managed', 'ssti', 'high', 1501, 'block', 'active', 'AND', 1, 'Template Expression', 'regex', 90, 'high', 87, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0109C', '01MRRUL0109', 0, 'request_body', 'regex', '(?i)(\{\{|\{%|#{|\<\%)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0110', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SSTI-002', 'Template Directive', 'Template Directive - Server-Side Template Injection signature', 'managed', 'ssti', 'high', 1502, 'block', 'active', 'AND', 1, 'Template Directive', 'regex', 89, 'high', 86, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0110C', '01MRRUL0110', 0, 'request_body', 'regex', '(?i)(\{%[\s\w]+%\}|\{\{[\s\w]+\}\})', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0111', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SSTI-003', 'Expression Language', 'Expression Language - Server-Side Template Injection signature', 'managed', 'ssti', 'high', 1503, 'block', 'active', 'AND', 1, 'Expression Language', 'regex', 91, 'high', 88, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0111C', '01MRRUL0111', 0, 'request_body', 'regex', '(?i)(\$\{[\w.]+\}|\#\{[\w.]+\})', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0112', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SSTI-004', 'Template Variable Abuse', 'Template Variable Abuse - Server-Side Template Injection signature', 'managed', 'ssti', 'high', 1504, 'block', 'active', 'AND', 1, 'Template Variable Abuse', 'regex', 87, 'high', 84, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0112C', '01MRRUL0112', 0, 'request_body', 'regex', '(?i)(\{\{[^}]*(config|settings|env|request|self)\b)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0113', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SSTI-005', 'Template Execution Pattern', 'Template Execution Pattern - Server-Side Template Injection signature', 'managed', 'ssti', 'high', 1505, 'block', 'active', 'AND', 1, 'Template Execution Pattern', 'regex', 93, 'high', 90, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0113C', '01MRRUL0113', 0, 'request_body', 'regex', '(?i)(\{\{[^}]*(system|exec|popen|os\.|import)\b)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0114', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-DESERIALIZATION-001', 'Unsafe Deserialization', 'Unsafe Deserialization - Deserialization signature', 'managed', 'deserialization', 'critical', 1601, 'block', 'active', 'AND', 1, 'Unsafe Deserialization', 'regex', 92, 'critical', 89, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0114C', '01MRRUL0114', 0, 'request_body', 'regex', '(?i)(O:[0-9]+:\"|a:[0-9]+:\{|\{\"\\x00)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0115', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-DESERIALIZATION-002', 'Serialized Object Abuse', 'Serialized Object Abuse - Deserialization signature', 'managed', 'deserialization', 'critical', 1602, 'block', 'active', 'AND', 1, 'Serialized Object Abuse', 'regex', 94, 'critical', 91, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0115C', '01MRRUL0115', 0, 'request_body', 'regex', '(?i)(PHP_O:|JAVA_OBJECT|ACED0005|rO0AB)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0116', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-DESERIALIZATION-003', 'Object Injection', 'Object Injection - Deserialization signature', 'managed', 'deserialization', 'critical', 1603, 'block', 'active', 'AND', 1, 'Object Injection', 'regex', 88, 'critical', 85, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0116C', '01MRRUL0116', 0, 'request_body', 'regex', '(?i)(\$\_REQUEST|\$\_GET|\$\_POST|\.\.\.O:)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0117', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-DESERIALIZATION-004', 'Java Deserialization', 'Java Deserialization - Deserialization signature', 'managed', 'deserialization', 'critical', 1604, 'block', 'active', 'AND', 1, 'Java Deserialization', 'regex', 96, 'critical', 93, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0117C', '01MRRUL0117', 0, 'request_body', 'regex', '(?i)(AC ED 00 05|rO0AB|\\xac\\xed)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0118', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-DESERIALIZATION-005', 'PHP Deserialization', 'PHP Deserialization - Deserialization signature', 'managed', 'deserialization', 'critical', 1605, 'block', 'active', 'AND', 1, 'PHP Deserialization', 'regex', 95, 'critical', 92, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0118C', '01MRRUL0118', 0, 'request_body', 'regex', '(?i)(O:[0-9]+:\"[A-Z][^\"]+\":[0-9]+:)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0119', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-DESERIALIZATION-006', '.NET Deserialization', '.NET Deserialization - Deserialization signature', 'managed', 'deserialization', 'critical', 1606, 'block', 'active', 'AND', 1, '.NET Deserialization', 'regex', 90, 'critical', 87, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0119C', '01MRRUL0119', 0, 'request_body', 'regex', '(?i)(AAEAAAD|System\.Object|/System\.Web)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0120', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SCANNER-001', 'Vulnerability Scanner', 'Vulnerability Scanner - Scanner Detection signature', 'managed', 'scanner', 'low', 1701, 'block', 'active', 'AND', 1, 'Vulnerability Scanner', 'regex', 95, 'low', 92, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0120C', '01MRRUL0120', 0, 'user_agent', 'regex', '(?i)(nikto|nessus|openvas|acunetix|w3af|sqlmap)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0121', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SCANNER-002', 'Web Scanner', 'Web Scanner - Scanner Detection signature', 'managed', 'scanner', 'low', 1702, 'block', 'active', 'AND', 1, 'Web Scanner', 'regex', 94, 'low', 91, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0121C', '01MRRUL0121', 0, 'user_agent', 'regex', '(?i)(gobuster|dirb|wfuzz|burpsuite|zap)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0122', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SCANNER-003', 'Directory Scanner', 'Directory Scanner - Scanner Detection signature', 'managed', 'scanner', 'low', 1703, 'block', 'active', 'AND', 1, 'Directory Scanner', 'regex', 93, 'low', 90, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0122C', '01MRRUL0122', 0, 'user_agent', 'regex', '(?i)(dirbuster|dirsearch|feroxbuster|ffuf)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0123', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SCANNER-004', 'Automated Reconnaissance', 'Automated Reconnaissance - Scanner Detection signature', 'managed', 'scanner', 'low', 1704, 'block', 'active', 'AND', 1, 'Automated Reconnaissance', 'regex', 92, 'low', 89, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0123C', '01MRRUL0123', 0, 'user_agent', 'regex', '(?i)(masscan|nmap|zmap|nuclei)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0124', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SCANNER-005', 'Security Testing Tool', 'Security Testing Tool - Scanner Detection signature', 'managed', 'scanner', 'low', 1705, 'block', 'active', 'AND', 1, 'Security Testing Tool', 'regex', 91, 'low', 88, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0124C', '01MRRUL0124', 0, 'user_agent', 'regex', '(?i)(hydra|medusa|john|hashcat)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0125', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SCANNER-006', 'Fingerprinting Tool', 'Fingerprinting Tool - Scanner Detection signature', 'managed', 'scanner', 'low', 1706, 'block', 'active', 'AND', 1, 'Fingerprinting Tool', 'regex', 89, 'low', 86, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0125C', '01MRRUL0125', 0, 'user_agent', 'regex', '(?i)(whatweb|wappalyzer|builtwith|recon-ng)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0126', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SCANNER-007', 'Automated Enumeration', 'Automated Enumeration - Scanner Detection signature', 'managed', 'scanner', 'low', 1707, 'block', 'active', 'AND', 1, 'Automated Enumeration', 'regex', 86, 'low', 83, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0126C', '01MRRUL0126', 0, 'user_agent', 'regex', '(?i)(crawler4j|python-requests|scrapy|apache-httpclient)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0127', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-BOT-001', 'Malicious Bot', 'Malicious Bot - Bot Protection signature', 'managed', 'bot', 'medium', 1801, 'block', 'active', 'AND', 1, 'Malicious Bot', 'regex', 90, 'medium', 87, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0127C', '01MRRUL0127', 0, 'user_agent', 'regex', '(?i)(sqlmap|nikto|curl|wget|python|go-http-client)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0128', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-BOT-002', 'Suspicious Bot', 'Suspicious Bot - Bot Protection signature', 'managed', 'bot', 'medium', 1802, 'block', 'active', 'AND', 1, 'Suspicious Bot', 'regex', 82, 'medium', 78, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0128C', '01MRRUL0128', 0, 'user_agent', 'regex', '(?i)(bot|crawler|spider|scraper)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0129', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-BOT-003', 'Automated Client', 'Automated Client - Bot Protection signature', 'managed', 'bot', 'medium', 1803, 'block', 'active', 'AND', 1, 'Automated Client', 'regex', 91, 'medium', 88, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0129C', '01MRRUL0129', 0, 'user_agent', 'regex', '(?i)(headless|phantomjs|selenium|puppeteer|playwright)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0130', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-BOT-004', 'Headless Browser', 'Headless Browser - Bot Protection signature', 'managed', 'bot', 'medium', 1804, 'block', 'active', 'AND', 1, 'Headless Browser', 'regex', 90, 'medium', 87, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0130C', '01MRRUL0130', 0, 'user_agent', 'regex', '(?i)(headlesschrome|chrome-headless|phantomjs)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0131', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-BOT-005', 'Browser Anomaly', 'Browser Anomaly - Bot Protection signature', 'managed', 'bot', 'medium', 1805, 'block', 'active', 'AND', 1, 'Browser Anomaly', 'regex', 84, 'medium', 80, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0131C', '01MRRUL0131', 0, 'user_agent', 'regex', '(?i)(^$|^[-_]+$|^\d+$)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0132', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-BOT-006', 'High-frequency Client', 'High-frequency Client - Bot Protection signature', 'managed', 'bot', 'medium', 1806, 'block', 'active', 'AND', 1, 'High-frequency Client', 'regex', 80, 'medium', 76, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0132C', '01MRRUL0132', 0, 'user_agent', 'regex', '(?i)(smtp|ftp|telnet|ssh|scanner|spider)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0133', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-AUTH-001', 'Credential Stuffing', 'Credential Stuffing - Authentication Attacks signature', 'managed', 'auth', 'high', 1901, 'block', 'active', 'AND', 1, 'Credential Stuffing', 'regex', 90, 'high', 87, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0133C', '01MRRUL0133', 0, 'request_body', 'regex', '(?i)(\"password\":\s*\"){2,}', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0134', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-AUTH-002', 'Authentication Bypass', 'Authentication Bypass - Authentication Attacks signature', 'managed', 'auth', 'high', 1902, 'block', 'active', 'AND', 1, 'Authentication Bypass', 'regex', 95, 'high', 92, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0134C', '01MRRUL0134', 0, 'request_body', 'regex', '(?i)((''|\")\s*(or|and)\s+[''\"]?\s*[''\"]?\s*=|password.*=.*''|'' OR 1=1)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0135', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-AUTH-003', 'Password Attack', 'Password Attack - Authentication Attacks signature', 'managed', 'auth', 'high', 1903, 'block', 'active', 'AND', 1, 'Password Attack', 'regex', 88, 'high', 85, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0135C', '01MRRUL0135', 0, 'request_body', 'regex', '(?i)(password|passwd|pwd)=[^&]*(admin|123456|qwerty|password)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0136', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-AUTH-004', 'Session Anomaly', 'Session Anomaly - Authentication Attacks signature', 'managed', 'auth', 'high', 1904, 'block', 'active', 'AND', 1, 'Session Anomaly', 'regex', 85, 'high', 82, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0136C', '01MRRUL0136', 0, 'header', 'regex', '(?i)(cookie:\s*.*(sessionid|jsessionid|phpsessid)=.{1,2}$)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0137', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-AUTH-005', 'Authentication Flood', 'Authentication Flood - Authentication Attacks signature', 'managed', 'auth', 'high', 1905, 'block', 'active', 'AND', 1, 'Authentication Flood', 'regex', 83, 'high', 79, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0137C', '01MRRUL0137', 0, 'url', 'regex', '(?i)(/login|/auth).*(bot|spider|scanner)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0138', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SESSION-001', 'Session Fixation', 'Session Fixation - Session Security signature', 'managed', 'session', 'medium', 2001, 'block', 'active', 'AND', 1, 'Session Fixation', 'regex', 89, 'medium', 86, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0138C', '01MRRUL0138', 0, 'header', 'regex', '(?i)(cookie:\s*.*(sessionid|phpsessid|jsessionid)=[^;]+)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0139', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SESSION-002', 'Session Token Anomaly', 'Session Token Anomaly - Session Security signature', 'managed', 'session', 'medium', 2002, 'block', 'active', 'AND', 1, 'Session Token Anomaly', 'regex', 87, 'medium', 84, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0139C', '01MRRUL0139', 0, 'header', 'regex', '(?i)(session(token|id)=[^;]{1,4}$)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0140', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SESSION-003', 'Cookie Manipulation', 'Cookie Manipulation - Session Security signature', 'managed', 'session', 'medium', 2003, 'block', 'active', 'AND', 1, 'Cookie Manipulation', 'regex', 91, 'medium', 88, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0140C', '01MRRUL0140', 0, 'header', 'regex', '(?i)(cookie:\s*.*(isadmin|role|user)=1|admin=true)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0141', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SESSION-004', 'Invalid Session', 'Invalid Session - Session Security signature', 'managed', 'session', 'medium', 2004, 'block', 'active', 'AND', 1, 'Invalid Session', 'regex', 86, 'medium', 83, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0141C', '01MRRUL0141', 0, 'header', 'regex', '(?i)(session(id|token)=\s*[\"'']?\s*[\"'']?)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0142', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SESSION-005', 'Session Replay', 'Session Replay - Session Security signature', 'managed', 'session', 'medium', 2005, 'block', 'active', 'AND', 1, 'Session Replay', 'regex', 84, 'medium', 80, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0142C', '01MRRUL0142', 0, 'header', 'regex', '(?i)(cookie:\s*.*(session|token)=[^;]{5,};\s*cookie:\s*.*)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0143', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-CSRF-001', 'Missing CSRF Token', 'Missing CSRF Token - CSRF signature', 'managed', 'csrf', 'medium', 2101, 'block', 'active', 'AND', 1, 'Missing CSRF Token', 'regex', 78, 'medium', 74, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0143C', '01MRRUL0143', 0, 'request_body', 'regex', '(?i)(method:\s*(post|put|delete).*(?!csrf|token))', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0144', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-CSRF-002', 'Invalid CSRF Token', 'Invalid CSRF Token - CSRF signature', 'managed', 'csrf', 'medium', 2102, 'block', 'active', 'AND', 1, 'Invalid CSRF Token', 'regex', 88, 'medium', 85, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0144C', '01MRRUL0144', 0, 'request_body', 'regex', '(?i)(csrf(token|_token)?=\s*[\"'']?[\"'']?|_token=\s*$)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0145', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-CSRF-003', 'CSRF Token Mismatch', 'CSRF Token Mismatch - CSRF signature', 'managed', 'csrf', 'medium', 2103, 'block', 'active', 'AND', 1, 'CSRF Token Mismatch', 'regex', 87, 'medium', 84, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0145C', '01MRRUL0145', 0, 'request_body', 'regex', '(?i)((csrf|_token|token)=[^&]+&.*(csrf|_token|token)=[^&]+)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0146', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-CSRF-004', 'Cross-origin Request Anomaly', 'Cross-origin Request Anomaly - CSRF signature', 'managed', 'csrf', 'medium', 2104, 'block', 'active', 'AND', 1, 'Cross-origin Request Anomaly', 'regex', 82, 'medium', 78, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0146C', '01MRRUL0146', 0, 'header', 'regex', '(?i)(origin:\s*https?://(?!([a-z0-9.-]+\.)?(example\.com|yourdomain\.com)))', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0147', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-INFO_DISCLOSURE-001', 'Source Code Exposure', 'Source Code Exposure - Information Disclosure signature', 'managed', 'info_disclosure', 'medium', 2201, 'block', 'active', 'AND', 1, 'Source Code Exposure', 'regex', 90, 'medium', 87, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0147C', '01MRRUL0147', 0, 'url', 'regex', '(?i)(\.(php|js|py|rb|java|c)$|/source|/src)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0148', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-INFO_DISCLOSURE-002', 'Configuration File Exposure', 'Configuration File Exposure - Information Disclosure signature', 'managed', 'info_disclosure', 'medium', 2202, 'block', 'active', 'AND', 1, 'Configuration File Exposure', 'regex', 89, 'medium', 86, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0148C', '01MRRUL0148', 0, 'url', 'regex', '(?i)(config\.(php|js|json|xml)|application\.(yml|yaml|properties))', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0149', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-INFO_DISCLOSURE-003', 'Backup File Exposure', 'Backup File Exposure - Information Disclosure signature', 'managed', 'info_disclosure', 'medium', 2203, 'block', 'active', 'AND', 1, 'Backup File Exposure', 'regex', 88, 'medium', 85, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0149C', '01MRRUL0149', 0, 'url', 'regex', '(?i)(\.(bak|old|orig|save|swp|tmp))$', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0150', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-INFO_DISCLOSURE-004', 'Debug Information', 'Debug Information - Information Disclosure signature', 'managed', 'info_disclosure', 'medium', 2204, 'block', 'active', 'AND', 1, 'Debug Information', 'regex', 91, 'medium', 88, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0150C', '01MRRUL0150', 0, 'url', 'regex', '(?i)(/debug|/trace|/phpinfo|/info\.php|/status)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0151', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-INFO_DISCLOSURE-005', 'Stack Trace Exposure', 'Stack Trace Exposure - Information Disclosure signature', 'managed', 'info_disclosure', 'medium', 2205, 'block', 'active', 'AND', 1, 'Stack Trace Exposure', 'regex', 87, 'medium', 84, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0151C', '01MRRUL0151', 0, 'request_body', 'regex', '(?i)(Traceback \(|at [a-z]+\.[a-z]+\.|\.java:[0-9]+|\.cs:[0-9]+)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0152', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-INFO_DISCLOSURE-006', 'Environment File Exposure', 'Environment File Exposure - Information Disclosure signature', 'managed', 'info_disclosure', 'medium', 2206, 'block', 'active', 'AND', 1, 'Environment File Exposure', 'regex', 95, 'medium', 92, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0152C', '01MRRUL0152', 0, 'url', 'regex', '(?i)(\.env|\.env\.local|\.env\.prod|\.env\.development)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0153', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-INFO_DISCLOSURE-007', 'Version Disclosure', 'Version Disclosure - Information Disclosure signature', 'managed', 'info_disclosure', 'medium', 2207, 'block', 'active', 'AND', 1, 'Version Disclosure', 'regex', 86, 'medium', 82, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0153C', '01MRRUL0153', 0, 'header', 'regex', '(?i)(server:\s*|x-powered-by:\s*|x-aspnet-version:\s*)[a-z0-9/.-]+', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0154', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-INFO_DISCLOSURE-008', 'Sensitive File Access', 'Sensitive File Access - Information Disclosure signature', 'managed', 'info_disclosure', 'medium', 2208, 'block', 'active', 'AND', 1, 'Sensitive File Access', 'regex', 94, 'medium', 91, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0154C', '01MRRUL0154', 0, 'url', 'regex', '(?i)(/etc/passwd|/etc/shadow|web\.config|php\.ini|\.git/config)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0155', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-RESOURCE_DISCOVERY-001', 'Admin Panel Discovery', 'Admin Panel Discovery - Resource Discovery signature', 'managed', 'resource_discovery', 'low', 2301, 'block', 'active', 'AND', 1, 'Admin Panel Discovery', 'regex', 89, 'low', 86, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0155C', '01MRRUL0155', 0, 'url', 'regex', '(?i)(/admin|/administrator|/wp-admin|/cpanel|/dashboard)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0156', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-RESOURCE_DISCOVERY-002', 'Login Panel Discovery', 'Login Panel Discovery - Resource Discovery signature', 'managed', 'resource_discovery', 'low', 2302, 'block', 'active', 'AND', 1, 'Login Panel Discovery', 'regex', 85, 'low', 82, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0156C', '01MRRUL0156', 0, 'url', 'regex', '(?i)(/login|/signin|/signup|/auth|/account)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0157', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-RESOURCE_DISCOVERY-003', 'Backup File Discovery', 'Backup File Discovery - Resource Discovery signature', 'managed', 'resource_discovery', 'low', 2303, 'block', 'active', 'AND', 1, 'Backup File Discovery', 'regex', 87, 'low', 84, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0157C', '01MRRUL0157', 0, 'url', 'regex', '(?i)(\.(zip|tar|gz|bak|old|sql))$', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0158', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-RESOURCE_DISCOVERY-004', 'Configuration Discovery', 'Configuration Discovery - Resource Discovery signature', 'managed', 'resource_discovery', 'low', 2304, 'block', 'active', 'AND', 1, 'Configuration Discovery', 'regex', 92, 'low', 89, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0158C', '01MRRUL0158', 0, 'url', 'regex', '(?i)(\.git|\.svn|\.env|web\.config|config\.php)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0159', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-RESOURCE_DISCOVERY-005', 'Git Repository Discovery', 'Git Repository Discovery - Resource Discovery signature', 'managed', 'resource_discovery', 'low', 2305, 'block', 'active', 'AND', 1, 'Git Repository Discovery', 'regex', 95, 'low', 92, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0159C', '01MRRUL0159', 0, 'url', 'regex', '(?i)(/\.git/(config|HEAD|index)|/\.git$|/\.svn/)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0160', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-RESOURCE_DISCOVERY-006', 'API Discovery', 'API Discovery - Resource Discovery signature', 'managed', 'resource_discovery', 'low', 2306, 'block', 'active', 'AND', 1, 'API Discovery', 'regex', 88, 'low', 85, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0160C', '01MRRUL0160', 0, 'url', 'regex', '(?i)(/api/|/swagger|/openapi|/graphql|/v[0-9]+/)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0161', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-RESOURCE_DISCOVERY-007', 'Debug Endpoint Discovery', 'Debug Endpoint Discovery - Resource Discovery signature', 'managed', 'resource_discovery', 'low', 2307, 'block', 'active', 'AND', 1, 'Debug Endpoint Discovery', 'regex', 86, 'low', 83, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0161C', '01MRRUL0161', 0, 'url', 'regex', '(?i)(/debug|/healthz|/actuator|/metrics|/status)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0162', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-REQUEST_ANOMALY-001', 'URL Too Long', 'URL Too Long - Request Anomaly signature', 'managed', 'request_anomaly', 'medium', 2401, 'block', 'active', 'AND', 1, 'URL Too Long', 'regex', 90, 'medium', 87, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0162C', '01MRRUL0162', 0, 'url', 'regex', '^.{2000,}$', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0163', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-REQUEST_ANOMALY-002', 'Header Too Large', 'Header Too Large - Request Anomaly signature', 'managed', 'request_anomaly', 'medium', 2402, 'block', 'active', 'AND', 1, 'Header Too Large', 'regex', 89, 'medium', 86, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0163C', '01MRRUL0163', 0, 'header', 'regex', '^.{5000,}$', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0164', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-REQUEST_ANOMALY-003', 'Body Too Large', 'Body Too Large - Request Anomaly signature', 'managed', 'request_anomaly', 'medium', 2403, 'block', 'active', 'AND', 1, 'Body Too Large', 'regex', 88, 'medium', 85, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0164C', '01MRRUL0164', 0, 'request_body', 'regex', '^.{1000000,}$', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0165', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-REQUEST_ANOMALY-004', 'Too Many Parameters', 'Too Many Parameters - Request Anomaly signature', 'managed', 'request_anomaly', 'medium', 2404, 'block', 'active', 'AND', 1, 'Too Many Parameters', 'regex', 87, 'medium', 84, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0165C', '01MRRUL0165', 0, 'url', 'regex', '([?&][a-z_]+=){30,}', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0166', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-REQUEST_ANOMALY-005', 'Invalid Encoding', 'Invalid Encoding - Request Anomaly signature', 'managed', 'request_anomaly', 'medium', 2405, 'block', 'active', 'AND', 1, 'Invalid Encoding', 'regex', 84, 'medium', 80, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0166C', '01MRRUL0166', 0, 'url', 'regex', '(?i)(%[0-9a-f]{1}$|%[g-z]|%(?![0-9a-f]{2}))', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0167', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-REQUEST_ANOMALY-006', 'Invalid Characters', 'Invalid Characters - Request Anomaly signature', 'managed', 'request_anomaly', 'medium', 2406, 'block', 'active', 'AND', 1, 'Invalid Characters', 'regex', 86, 'medium', 82, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0167C', '01MRRUL0167', 0, 'url', 'regex', '(?i)([\x00-\x1f\x7f])', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0168', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-REQUEST_ANOMALY-007', 'Invalid Content-Length', 'Invalid Content-Length - Request Anomaly signature', 'managed', 'request_anomaly', 'medium', 2407, 'block', 'active', 'AND', 1, 'Invalid Content-Length', 'regex', 88, 'medium', 85, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0168C', '01MRRUL0168', 0, 'header', 'regex', '(?i)content-length:\s*0{5,}|content-length:\s*[^0-9]', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0169', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-REQUEST_ANOMALY-008', 'Malformed JSON', 'Malformed JSON - Request Anomaly signature', 'managed', 'request_anomaly', 'medium', 2408, 'block', 'active', 'AND', 1, 'Malformed JSON', 'regex', 83, 'medium', 80, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0169C', '01MRRUL0169', 0, 'request_body', 'regex', '(?i)(\}[^\s,}\]]|\{:|\[:|,\s*[}\]])', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0170', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-REQUEST_ANOMALY-009', 'Protocol Violation', 'Protocol Violation - Request Anomaly signature', 'managed', 'request_anomaly', 'medium', 2409, 'block', 'active', 'AND', 1, 'Protocol Violation', 'regex', 82, 'medium', 78, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0170C', '01MRRUL0170', 0, 'url', 'regex', '(?i)(%[0-9a-f]{1}\b|\\x[0-9a-f]{2})', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0171', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-IP-001', 'IP Blocklist', 'IP Blocklist - IP / Reputation signature', 'managed', 'ip', 'high', 2501, 'block', 'active', 'AND', 1, 'IP Blocklist', 'cidr_match', 90, 'high', 88, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0171C', '01MRRUL0171', 0, 'source_ip', 'cidr_match', '0.0.0.0/0', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0172', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-IP-002', 'CIDR Block', 'CIDR Block - IP / Reputation signature', 'managed', 'ip', 'high', 2502, 'block', 'active', 'AND', 1, 'CIDR Block', 'cidr_match', 88, 'high', 85, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0172C', '01MRRUL0172', 0, 'source_ip', 'cidr_match', '0.0.0.0/0', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0173', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-IP-003', 'Proxy Detection', 'Proxy Detection - IP / Reputation signature', 'managed', 'ip', 'high', 2503, 'block', 'active', 'AND', 1, 'Proxy Detection', 'regex', 85, 'high', 82, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0173C', '01MRRUL0173', 0, 'header', 'regex', '(?i)(x-forwarded-for:\s*[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+|via:\s*[a-z0-9-]+)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0174', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-IP-004', 'Threat Intelligence Match', 'Threat Intelligence Match - IP / Reputation signature', 'managed', 'ip', 'high', 2504, 'block', 'active', 'AND', 1, 'Threat Intelligence Match', 'regex', 87, 'high', 84, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0174C', '01MRRUL0174', 0, 'source_ip', 'regex', '^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0175', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-GEO-001', 'Country Blocklist', 'Country Blocklist - Geo Security signature', 'managed', 'geo', 'low', 2601, 'block', 'active', 'AND', 1, 'Country Blocklist', 'regex', 90, 'low', 87, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0175C', '01MRRUL0175', 0, 'source_ip', 'regex', '^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0176', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-GEO-002', 'Country Allowlist', 'Country Allowlist - Geo Security signature', 'managed', 'geo', 'low', 2602, 'block', 'active', 'AND', 1, 'Country Allowlist', 'regex', 90, 'low', 87, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0176C', '01MRRUL0176', 0, 'source_ip', 'regex', '^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0177', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-GEO-003', 'High-risk Location Policy', 'High-risk Location Policy - Geo Security signature', 'managed', 'geo', 'low', 2603, 'block', 'active', 'AND', 1, 'High-risk Location Policy', 'regex', 88, 'low', 85, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0177C', '01MRRUL0177', 0, 'source_ip', 'regex', '^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0178', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-RATE_LIMIT-001', 'Global Rate Limit', 'Global Rate Limit - Rate Limiting signature', 'managed', 'rate_limit', 'medium', 2701, 'block', 'active', 'AND', 1, 'Global Rate Limit', 'equals', 85, 'medium', 82, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0178C', '01MRRUL0178', 0, 'method', 'equals', 'GET', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0179', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-RATE_LIMIT-002', 'IP Rate Limit', 'IP Rate Limit - Rate Limiting signature', 'managed', 'rate_limit', 'medium', 2702, 'block', 'active', 'AND', 1, 'IP Rate Limit', 'regex', 86, 'medium', 83, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0179C', '01MRRUL0179', 0, 'source_ip', 'regex', '^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0180', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-RATE_LIMIT-003', 'Endpoint Rate Limit', 'Endpoint Rate Limit - Rate Limiting signature', 'managed', 'rate_limit', 'medium', 2703, 'block', 'active', 'AND', 1, 'Endpoint Rate Limit', 'regex', 84, 'medium', 80, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0180C', '01MRRUL0180', 0, 'url', 'regex', '(?i)(/login|/api/|/search|/auth)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0181', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-RATE_LIMIT-004', 'Login Rate Limit', 'Login Rate Limit - Rate Limiting signature', 'managed', 'rate_limit', 'medium', 2704, 'block', 'active', 'AND', 1, 'Login Rate Limit', 'regex', 87, 'medium', 84, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0181C', '01MRRUL0181', 0, 'url', 'regex', '(?i)(/login|/signin|/auth)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0182', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-RATE_LIMIT-005', 'API Rate Limit', 'API Rate Limit - Rate Limiting signature', 'managed', 'rate_limit', 'medium', 2705, 'block', 'active', 'AND', 1, 'API Rate Limit', 'regex', 85, 'medium', 82, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0182C', '01MRRUL0182', 0, 'url', 'regex', '(?i)/api/', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0183', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-RATE_LIMIT-006', 'Burst Detection', 'Burst Detection - Rate Limiting signature', 'managed', 'rate_limit', 'medium', 2706, 'block', 'active', 'AND', 1, 'Burst Detection', 'regex', 78, 'medium', 74, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0183C', '01MRRUL0183', 0, 'request_body', 'regex', '(?i)(burst|flood|rapid)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0184', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-RATE_LIMIT-007', 'Request Flood Detection', 'Request Flood Detection - Rate Limiting signature', 'managed', 'rate_limit', 'medium', 2707, 'block', 'active', 'AND', 1, 'Request Flood Detection', 'regex', 80, 'medium', 76, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0184C', '01MRRUL0184', 0, 'url', 'regex', '(?i)(/flood|/attack|/stress)', 'lowercase', false);

-- 27 categories, 184 signature rules seeded.