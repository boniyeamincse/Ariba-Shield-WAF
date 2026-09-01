-- 0034_phase1_new_categories.up.sql
-- Phase 1: 6 new rule categories (IDOR, Open Redirect, Clickjacking, JWT, Prototype Pollution, CORS) + managed groups.

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT101', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'Insecure Direct Object Reference', 'idor', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT102', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'Open Redirect', 'open_redirect', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT103', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'Clickjacking', 'clickjacking', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT104', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'JWT Misconfiguration', 'jwt', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT105', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'Prototype Pollution', 'prototype_pollution', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES
  ('01MRCAT106', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'CORS Misconfiguration', 'cors', true, 'low', 1, 'block', 5, 'active')
ON CONFLICT DO NOTHING;

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0201', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-IDOR-001', 'Predictable Object ID', 'Predictable Object ID - Insecure Direct Object Reference signature', 'managed', 'idor', 'high', 1001, 'block', 'active', 'AND', 1, 'IDOR - Predictable Object ID', 'regex', 85, 'high', 82, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0201C', '01MRRUL0201', 0, 'url', 'regex', '(?i)(id|user_id|account_id|doc_id|file_id|order_id|profile_id)=\d{1,6}$', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0202', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-IDOR-002', 'Numeric ID Enumeration', 'Numeric ID Enumeration - Insecure Direct Object Reference signature', 'managed', 'idor', 'high', 1002, 'block', 'active', 'AND', 1, 'IDOR - Numeric ID Enumeration', 'regex', 88, 'high', 85, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0202C', '01MRRUL0202', 0, 'url', 'regex', '(?i)(/(user|account|order|profile|invoice|transaction|document)/\d{1,6}$)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0203', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-IDOR-003', 'UUID Reference Abuse', 'UUID Reference Abuse - Insecure Direct Object Reference signature', 'managed', 'idor', 'high', 1003, 'block', 'active', 'AND', 1, 'IDOR - UUID Reference Abuse', 'regex', 80, 'high', 76, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0203C', '01MRRUL0203', 0, 'url', 'regex', '(?i)(/api/[\w-]+/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0204', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-IDOR-004', 'Object Reference in Params', 'Object Reference in Params - Insecure Direct Object Reference signature', 'managed', 'idor', 'high', 1004, 'block', 'active', 'AND', 1, 'IDOR - Object Reference in Params', 'regex', 87, 'high', 84, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0204C', '01MRRUL0204', 0, 'url', 'regex', '(?i)(\?|&)(id|uid|object|ref|resource|target)=\d{1,6}(&|$)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0205', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-IDOR-005', 'User ID Parameter Tampering', 'User ID Parameter Tampering - Insecure Direct Object Reference signature', 'managed', 'idor', 'high', 1005, 'block', 'active', 'AND', 1, 'IDOR - User ID Parameter Tampering', 'regex', 84, 'high', 81, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0205C', '01MRRUL0205', 0, 'url', 'regex', '(?i)(user[id]?|account|profile)[-_]?[id]?=(\d{1,6}|[a-z0-9]{4,16})$', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0206', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-IDOR-006', 'Account ID Enumeration', 'Account ID Enumeration - Insecure Direct Object Reference signature', 'managed', 'idor', 'high', 1006, 'block', 'active', 'AND', 1, 'IDOR - Account ID Enumeration', 'regex', 86, 'high', 83, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0206C', '01MRRUL0206', 0, 'url', 'regex', '(?i)(/account/|/users/|/customers/|/members/)\d{1,8}', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0207', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-IDOR-007', 'Document ID Access', 'Document ID Access - Insecure Direct Object Reference signature', 'managed', 'idor', 'high', 1007, 'block', 'active', 'AND', 1, 'IDOR - Document ID Access', 'regex', 85, 'high', 82, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0207C', '01MRRUL0207', 0, 'url', 'regex', '(?i)(/document|/download|/attachment|/file|/invoice)[/-]\d{1,8}', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0208', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-IDOR-008', 'Hidden Endpoint Reference', 'Hidden Endpoint Reference - Insecure Direct Object Reference signature', 'managed', 'idor', 'high', 1008, 'block', 'active', 'AND', 1, 'IDOR - Hidden Endpoint Reference', 'regex', 82, 'high', 79, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0208C', '01MRRUL0208', 0, 'url', 'regex', '(?i)(/internal|/admin|/debug|/backup)[-_]?[\w]*/\d{1,8}', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0209', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-OPEN_REDIRECT-001', 'Open Redirect (url param)', 'Open Redirect (url param) - Open Redirect signature', 'managed', 'open_redirect', 'high', 1101, 'block', 'active', 'AND', 1, 'Open Redirect', 'regex', 92, 'high', 89, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0209C', '01MRRUL0209', 0, 'url', 'regex', '(?i)(\?|&)(url|redirect|return|next|dest|destination|target|goto|r|u)=https?://', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0210', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-OPEN_REDIRECT-002', 'Redirect Parameter Abuse', 'Redirect Parameter Abuse - Open Redirect signature', 'managed', 'open_redirect', 'high', 1102, 'block', 'active', 'AND', 1, 'Open Redirect - Redirect Parameter Abuse', 'regex', 90, 'high', 87, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0210C', '01MRRUL0210', 0, 'url', 'regex', '(?i)(\?|&)(redirect|return|next|dest|target|url|link|out)=//', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0211', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-OPEN_REDIRECT-003', 'Next URL Redirect', 'Next URL Redirect - Open Redirect signature', 'managed', 'open_redirect', 'high', 1103, 'block', 'active', 'AND', 1, 'Open Redirect - Next URL Redirect', 'regex', 91, 'high', 88, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0211C', '01MRRUL0211', 0, 'url', 'regex', '(?i)(\?|&)(next|return_url|returnUrl|redirect_uri|redirect_url)=[^\s&]*(http|//)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0212', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-OPEN_REDIRECT-004', 'External Redirect Detection', 'External Redirect Detection - Open Redirect signature', 'managed', 'open_redirect', 'high', 1104, 'block', 'active', 'AND', 1, 'Open Redirect - External Redirect', 'regex', 88, 'high', 85, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0212C', '01MRRUL0212', 0, 'url', 'regex', '(?i)(\?|&)(url|redirect|return)=[a-z0-9.-]+\.(com|net|org|io|co|ru|cn|info|xyz)\b', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0213', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-OPEN_REDIRECT-005', 'Whitelist Bypass Redirect', 'Whitelist Bypass Redirect - Open Redirect signature', 'managed', 'open_redirect', 'high', 1105, 'block', 'active', 'AND', 1, 'Open Redirect - Whitelist Bypass', 'regex', 89, 'high', 86, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0213C', '01MRRUL0213', 0, 'url', 'regex', '(?i)(\?|&)(url|redirect|next)=[^&\s]*@[^&\s]*(http|//)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0214', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-OPEN_REDIRECT-006', 'Protocol-Relative Redirect', 'Protocol-Relative Redirect - Open Redirect signature', 'managed', 'open_redirect', 'high', 1106, 'block', 'active', 'AND', 1, 'Open Redirect - Protocol-Relative', 'regex', 87, 'high', 84, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0214C', '01MRRUL0214', 0, 'url', 'regex', '(?i)(\?|&)(redirect|next|url|return)=//[^&\s]+', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0215', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-OPEN_REDIRECT-007', 'OAuth Redirect_uri Abuse', 'OAuth Redirect_uri Abuse - Open Redirect signature', 'managed', 'open_redirect', 'high', 1107, 'block', 'active', 'AND', 1, 'OAuth Redirect_uri Abuse', 'regex', 93, 'high', 90, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0215C', '01MRRUL0215', 0, 'url', 'regex', '(?i)redirect_uri=[^&\s]*(http://localhost|http://127\.0\.0\.1|https?://evil)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0216', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-CLICKJACKING-001', 'Frameable Response (no XFO)', 'Frameable Response (no XFO) - Clickjacking signature', 'managed', 'clickjacking', 'medium', 1201, 'block', 'active', 'AND', 1, 'Clickjacking - Missing X-Frame-Options', 'not_contains', 80, 'medium', 76, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0216C', '01MRRUL0216', 0, 'header', 'not_contains', 'X-Frame-Options', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0217', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-CLICKJACKING-002', 'CSP frame-ancestors Missing', 'CSP frame-ancestors Missing - Clickjacking signature', 'managed', 'clickjacking', 'medium', 1202, 'block', 'active', 'AND', 1, 'Clickjacking - Missing CSP frame-ancestors', 'not_contains', 78, 'medium', 74, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0217C', '01MRRUL0217', 0, 'header', 'not_contains', 'frame-ancestors', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0218', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-CLICKJACKING-003', 'Iframe Embedding Attempt', 'Iframe Embedding Attempt - Clickjacking signature', 'managed', 'clickjacking', 'medium', 1203, 'block', 'active', 'AND', 1, 'Clickjacking - Iframe Embedding', 'regex', 85, 'medium', 82, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0218C', '01MRRUL0218', 0, 'request_body', 'regex', '(?i)(<iframe[\s>]|framebust|top\.location|window\.top)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0219', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-CLICKJACKING-004', 'UI Redressing Probe', 'UI Redressing Probe - Clickjacking signature', 'managed', 'clickjacking', 'medium', 1204, 'block', 'active', 'AND', 1, 'Clickjacking - UI Redressing', 'regex', 75, 'medium', 71, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0219C', '01MRRUL0219', 0, 'request_body', 'regex', '(?i)(opacity:\s*0|z-index|transparent.*iframe)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0220', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-JWT-001', 'JWT Token Pattern', 'JWT Token Pattern - JWT Misconfiguration signature', 'managed', 'jwt', 'critical', 1301, 'block', 'active', 'AND', 1, 'JWT Token Present', 'regex', 70, 'critical', 65, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0220C', '01MRRUL0220', 0, 'header', 'regex', 'Authorization:\s*Bearer\s+[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0221', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-JWT-002', 'JWT alg:none Attack', 'JWT alg:none Attack - JWT Misconfiguration signature', 'managed', 'jwt', 'critical', 1302, 'block', 'active', 'AND', 1, 'JWT - alg:none (missing signature)', 'regex', 94, 'critical', 92, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0221C', '01MRRUL0221', 0, 'header', 'regex', 'eyJ[A-Za-z0-9_-]*\.[A-Za-z0-9_-]*\.$', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0222', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-JWT-003', 'JWT Algorithm Confusion', 'JWT Algorithm Confusion - JWT Misconfiguration signature', 'managed', 'jwt', 'critical', 1303, 'block', 'active', 'AND', 1, 'JWT - HS256/RS256 Confusion', 'regex', 88, 'critical', 85, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0222C', '01MRRUL0222', 0, 'header', 'regex', '(\"alg\"\s*:\s*\"HS(256|384|512)\"|alg.{0,4}HS256)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0223', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-JWT-004', 'JWT Header Injection (kid)', 'JWT Header Injection (kid) - JWT Misconfiguration signature', 'managed', 'jwt', 'critical', 1304, 'block', 'active', 'AND', 1, 'JWT - kid Path Traversal', 'regex', 92, 'critical', 89, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0223C', '01MRRUL0223', 0, 'request_body', 'regex', '(?i)\"kid\"\s*:\s*\"[^\"]*\.\./|/etc/passwd|/dev/', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0224', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-JWT-005', 'JWT Weak Algorithm', 'JWT Weak Algorithm - JWT Misconfiguration signature', 'managed', 'jwt', 'critical', 1305, 'block', 'active', 'AND', 1, 'JWT - Weak Algorithm', 'regex', 82, 'critical', 78, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0224C', '01MRRUL0224', 0, 'header', 'regex', '(\"alg\"\s*:\s*\"(none|NONE|HS256|RS256|ES256)\")', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0225', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-JWT-006', 'JWT Token in URL Leak', 'JWT Token in URL Leak - JWT Misconfiguration signature', 'managed', 'jwt', 'critical', 1306, 'block', 'active', 'AND', 1, 'JWT - Token Leakage in URL', 'regex', 93, 'critical', 90, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0225C', '01MRRUL0225', 0, 'url', 'regex', '(?i)(\?|&)(token|jwt|access_token)=[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0226', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-PROTOTYPE_POLLUTION-001', '__proto__ Injection', '__proto__ Injection - Prototype Pollution signature', 'managed', 'prototype_pollution', 'critical', 1401, 'block', 'active', 'AND', 1, 'Prototype Pollution - __proto__', 'regex', 97, 'critical', 95, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0226C', '01MRRUL0226', 0, 'request_body', 'regex', '\"__proto__\"\s*:', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0227', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-PROTOTYPE_POLLUTION-002', 'constructor.prototype', 'constructor.prototype - Prototype Pollution signature', 'managed', 'prototype_pollution', 'critical', 1402, 'block', 'active', 'AND', 1, 'Prototype Pollution - constructor.prototype', 'regex', 96, 'critical', 94, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0227C', '01MRRUL0227', 0, 'request_body', 'regex', '\"constructor\"\s*:\s*\{\s*\"prototype\"', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0228', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-PROTOTYPE_POLLUTION-003', 'Object Merge Pollution', 'Object Merge Pollution - Prototype Pollution signature', 'managed', 'prototype_pollution', 'critical', 1403, 'block', 'active', 'AND', 1, 'Prototype Pollution - Merge Sink', 'regex', 90, 'critical', 87, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0228C', '01MRRUL0228', 0, 'request_body', 'regex', '(?i)(_.merge|Object\.assign|\.extend\()', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0229', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-PROTOTYPE_POLLUTION-004', 'JSON __proto__ Key', 'JSON __proto__ Key - Prototype Pollution signature', 'managed', 'prototype_pollution', 'critical', 1404, 'block', 'active', 'AND', 1, 'Prototype Pollution - JSON Key', 'regex', 93, 'critical', 90, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0229C', '01MRRUL0229', 0, 'request_body', 'regex', '(\[\"__proto__\"\]|\[\"constructor\"\])', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0230', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-PROTOTYPE_POLLUTION-005', 'Nested Prototype Chain', 'Nested Prototype Chain - Prototype Pollution signature', 'managed', 'prototype_pollution', 'critical', 1405, 'block', 'active', 'AND', 1, 'Prototype Pollution - Nested Chain', 'regex', 91, 'critical', 88, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0230C', '01MRRUL0230', 0, 'request_body', 'regex', '(\"__proto__\"[^\}]{0,200}\"__proto__\"|\"constructor\"[^\}]{0,200}\"prototype\")', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0231', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-CORS-001', 'Reflected Origin', 'Reflected Origin - CORS Misconfiguration signature', 'managed', 'cors', 'medium', 1501, 'block', 'active', 'AND', 1, 'CORS - Origin Reflected', 'regex', 72, 'medium', 68, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0231C', '01MRRUL0231', 0, 'header', 'regex', 'Origin:\s*https?://[^\s]+', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0232', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-CORS-002', 'Null Origin', 'Null Origin - CORS Misconfiguration signature', 'managed', 'cors', 'medium', 1502, 'block', 'active', 'AND', 1, 'CORS - Null Origin', 'regex', 88, 'medium', 85, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0232C', '01MRRUL0232', 0, 'header', 'regex', 'Origin:\s*null', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0233', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-CORS-003', 'Credentials with Origin', 'Credentials with Origin - CORS Misconfiguration signature', 'managed', 'cors', 'medium', 1503, 'block', 'active', 'AND', 1, 'CORS - Credentialed Request', 'regex', 76, 'medium', 72, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0233C', '01MRRUL0233', 0, 'header', 'regex', '(Origin:\s*https?://[^\s]+(access-control|Access-Control))', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0234', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-CORS-004', 'CORS Origin Spoof Probe', 'CORS Origin Spoof Probe - CORS Misconfiguration signature', 'managed', 'cors', 'medium', 1504, 'block', 'active', 'AND', 1, 'CORS - Origin Spoof', 'regex', 90, 'medium', 87, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0234C', '01MRRUL0234', 0, 'header', 'regex', '(Origin:\s*https?://[a-z0-9-]+\.(evil|attacker|test)\.(com|net|io)|Origin:\s*https?://localhost)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0235', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-CORS-005', 'Preflight OPTIONS Abuse', 'Preflight OPTIONS Abuse - CORS Misconfiguration signature', 'managed', 'cors', 'medium', 1505, 'block', 'active', 'AND', 1, 'CORS - Preflight Probe', 'equals', 70, 'medium', 65, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0235C', '01MRRUL0235', 0, 'method', 'equals', 'OPTIONS', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0236', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-CORS-006', 'Subdomain Regex Bypass', 'Subdomain Regex Bypass - CORS Misconfiguration signature', 'managed', 'cors', 'medium', 1506, 'block', 'active', 'AND', 1, 'CORS - Subdomain Regex Bypass', 'regex', 80, 'medium', 76, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0236C', '01MRRUL0236', 0, 'header', 'regex', '(Origin:\s*https?://(?!([a-z0-9-]+\.)*[a-z0-9-]+\.(com|net|org)$))', 'lowercase', false);

-- 6 categories, 36 new signature rules added.