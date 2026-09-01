-- 0035_phase2_expand_categories.up.sql
-- Phase 2: expand existing categories with sub-type rules (unique rule_ids).

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0301', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SQLI-P2-001', 'Out-of-Band SQL Injection', 'Out-of-Band SQL Injection - SQL Injection signature', 'managed', 'sqli', 'critical', 2001, 'block', 'active', 'AND', 1, 'Out-of-Band SQL Injection', 'regex', 88, 'critical', 85, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0301C', '01MRRUL0301', 0, 'request_body', 'regex', '(?i)(load_file\(|into\s+outfile|into\s+dumpfile|dnslog|burpcollaborator\.oast\.)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0302', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-XSS-P2-001', 'Blind XSS', 'Blind XSS - Cross-Site Scripting signature', 'managed', 'xss', 'high', 2101, 'block', 'active', 'AND', 1, 'Blind XSS', 'regex', 89, 'high', 86, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0302C', '01MRRUL0302', 0, 'request_body', 'regex', '(?i)(<script[^>]*>.*(fetch|XMLHttpRequest|new\s+Image|document\.cookie|location\.href)|(http|https):\/\/[a-z0-9.-]+\.(burpcollaborator|xss\.ht|oast\.))', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0303', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-XSS-P2-002', 'Self XSS', 'Self XSS - Cross-Site Scripting signature', 'managed', 'xss', 'high', 2102, 'block', 'active', 'AND', 1, 'Self XSS', 'regex', 84, 'high', 80, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0303C', '01MRRUL0303', 0, 'request_body', 'regex', '(?i)(javascript:document\.(location|cookie|title)|<script>alert\()', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0304', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-XSS-P2-003', 'SVG XSS', 'SVG XSS - Cross-Site Scripting signature', 'managed', 'xss', 'high', 2103, 'block', 'active', 'AND', 1, 'SVG XSS', 'regex', 93, 'high', 90, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0304C', '01MRRUL0304', 0, 'request_body', 'regex', '(?i)(<svg[\s>][^>]*on(load|error|click)=|</svg>.*<script)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0305', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-XSS-P2-004', 'JSONP XSS', 'JSONP XSS - Cross-Site Scripting signature', 'managed', 'xss', 'high', 2104, 'block', 'active', 'AND', 1, 'JSONP XSS', 'regex', 91, 'high', 88, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0305C', '01MRRUL0305', 0, 'request_body', 'regex', '(?i)(callback=|jsonp=)[^&]*[\(\{\"<]', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0306', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-XSS-P2-005', 'Cookie-Based XSS', 'Cookie-Based XSS - Cross-Site Scripting signature', 'managed', 'xss', 'high', 2105, 'block', 'active', 'AND', 1, 'Cookie-Based XSS', 'regex', 87, 'high', 84, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0306C', '01MRRUL0306', 0, 'request_body', 'regex', '(?i)(document\.cookie\s*=|;cookie\s*=|setcookie\(.*script)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0307', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-XSS-P2-006', 'DOM-Based XSS Sink', 'DOM-Based XSS Sink - Cross-Site Scripting signature', 'managed', 'xss', 'high', 2106, 'block', 'active', 'AND', 1, 'DOM-Based XSS', 'regex', 90, 'high', 87, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0307C', '01MRRUL0307', 0, 'request_body', 'regex', '(?i)(document\.(write|writeln)|\.innerHTML\s*=\s*.*(location|hash|search|referrer)|eval\(.*(location|hash|search))', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0308', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-AUTH-P2-001', '2FA Bypass', '2FA Bypass - Authentication Attacks signature', 'managed', 'auth', 'high', 2201, 'block', 'active', 'AND', 1, '2FA Bypass', 'regex', 85, 'high', 82, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0308C', '01MRRUL0308', 0, 'request_body', 'regex', '(?i)((otp|code|token|pin)=.{0,4}$|skip.*(2fa|otp|mfa)|debug.*(2fa|otp|mfa))', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0309', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-AUTH-P2-002', 'OAuth Misconfiguration', 'OAuth Misconfiguration - Authentication Attacks signature', 'managed', 'auth', 'high', 2202, 'block', 'active', 'AND', 1, 'OAuth Misconfiguration', 'regex', 89, 'high', 86, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0309C', '01MRRUL0309', 0, 'url', 'regex', '(?i)(redirect_uri=[^&\s]*(@|localhost|127\.0\.0\.1|%0d%0a)|response_type=token.*scope=admin)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0310', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-AUTH-P2-003', 'Password Reset Poisoning', 'Password Reset Poisoning - Authentication Attacks signature', 'managed', 'auth', 'high', 2203, 'block', 'active', 'AND', 1, 'Password Reset Poisoning', 'regex', 93, 'high', 90, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0310C', '01MRRUL0310', 0, 'header', 'regex', '(?i)(Host:\s*(evil|attacker|test)\.(com|net|io)|X-Forwarded-Host:\s*evil)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0311', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-AUTH-P2-004', 'Privilege Escalation Probe', 'Privilege Escalation Probe - Authentication Attacks signature', 'managed', 'auth', 'high', 2204, 'block', 'active', 'AND', 1, 'Privilege Escalation', 'regex', 92, 'high', 89, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0311C', '01MRRUL0311', 0, 'request_body', 'regex', '(?i)(\"role\"\s*:\s*\"(admin|superuser|root)\"|\"is_admin\"\s*:\s*(true|1)|role=admin&|admin=true)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0312', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-AUTH-P2-005', 'Account Takeover Token', 'Account Takeover Token - Authentication Attacks signature', 'managed', 'auth', 'high', 2205, 'block', 'active', 'AND', 1, 'Account Takeover Token', 'regex', 90, 'high', 87, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0312C', '01MRRUL0312', 0, 'request_body', 'regex', '(?i)(reset_token=[a-z0-9]{1,6}$|password_reset.*token.*=\s*[''\"]{0,1}[a-z0-9]{1,6})', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0313', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SESSION-P2-001', 'Session Hijacking Probe', 'Session Hijacking Probe - Session Security signature', 'managed', 'session', 'medium', 2301, 'block', 'active', 'AND', 1, 'Session Hijacking Probe', 'regex', 83, 'medium', 79, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0313C', '01MRRUL0313', 0, 'header', 'regex', '(?i)(cookie:\s*.*(session|token|auth)=.{1,8}$|X-Forwarded-For:\s*[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+.*cookie)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0314', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-SESSION-P2-002', 'Token Leakage in Referer', 'Token Leakage in Referer - Session Security signature', 'managed', 'session', 'medium', 2302, 'block', 'active', 'AND', 1, 'Token Leakage in Referer', 'regex', 88, 'medium', 85, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0314C', '01MRRUL0314', 0, 'header', 'regex', '(?i)referer:\s*[^\s]*\?(.*(token|session|auth|jwt|key)=[a-z0-9_-]{8,})', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0315', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-RATE_LIMIT-P2-001', 'Rate Limit Bypass - Header', 'Rate Limit Bypass - Header - Rate Limiting signature', 'managed', 'rate_limit', 'medium', 2401, 'block', 'active', 'AND', 1, 'Rate Limit Bypass via Header', 'regex', 80, 'medium', 76, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0315C', '01MRRUL0315', 0, 'header', 'regex', '(?i)(X-Forwarded-For:\s*[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+.*X-Forwarded-For|X-Real-IP:\s*[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0316', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-RATE_LIMIT-P2-002', 'Rate Limit Bypass - Param', 'Rate Limit Bypass - Param - Rate Limiting signature', 'managed', 'rate_limit', 'medium', 2402, 'block', 'active', 'AND', 1, 'Rate Limit Bypass via Param', 'regex', 82, 'medium', 78, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0316C', '01MRRUL0316', 0, 'url', 'regex', '(?i)(\?|&)(client_ip|forwarded|real_ip|cf-ip)=[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+', 'lowercase', false);

INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES
  ('01MRRUL0317', '01ARZ3NDEKTSV4RRFFQ69G5FAV', 'ARB-RATE_LIMIT-P2-003', 'Brute Force Burst', 'Brute Force Burst - Rate Limiting signature', 'managed', 'rate_limit', 'medium', 2403, 'block', 'active', 'AND', 1, 'Brute Force Burst Pattern', 'regex', 75, 'medium', 70, 'ariba-core', false, '')
ON CONFLICT DO NOTHING;
INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES
  ('01MRRUL0317C', '01MRRUL0317', 0, 'request_body', 'regex', '(?i)((password|passwd|pwd)=[^&\s]*)&', 'lowercase', false);
