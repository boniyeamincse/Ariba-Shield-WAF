#!/usr/bin/env python3
"""Generate 0027_seed_rule_catalog.up.sql — 27 categories + sub-signatures + metadata."""
import datetime

ORG = "01ARZ3NDEKTSV4RRFFQ69G5FAV"


def sql(s: str) -> str:
    return s.replace("'", "''")


# (category_key, display_name, severity, risk, [(sub_name, field, op, pattern, accuracy, confidence, attack_type)])
CATS = [
    ("sqli", "SQL Injection", "critical", "critical", [
        ("SQL Injection", "request_body", "regex", r"(?i)(union[\s]+select|select[\s]+.*[\s]+from|insert[\s]+into|delete[\s]+from)", 98, 95, "SQL Injection"),
        ("Blind SQL Injection", "request_body", "regex", r"(?i)(and|or)[\s]+[0-9]+[=][0-9]+", 92, 90, "Blind SQL Injection"),
        ("Boolean-based SQL Injection", "request_body", "regex", r"(?i)('|\))[\s]*(and|or)[\s]+('[^']*'|[0-9]+)[\s]*=[\s]*('[^']*'|[0-9]+)", 90, 88, "Boolean-based SQL Injection"),
        ("Time-based SQL Injection", "request_body", "regex", r"(?i)(sleep|benchmark|pg_sleep|waitfor[\s]+delay)[\s]*\(", 94, 92, "Time-based SQL Injection"),
        ("Union-based SQL Injection", "request_body", "regex", r"(?i)union[\s]+(all[\s]+)?select", 96, 93, "Union-based SQL Injection"),
        ("Error-based SQL Injection", "request_body", "regex", r"(?i)(extractvalue|updatexml|group[\s]+by[\s]+with[\s]+rollup)", 90, 87, "Error-based SQL Injection"),
        ("Stored SQL Injection", "request_body", "regex", r"(?i)(into[\s]+outfile|load_file|into[\s]+dumpfile)", 89, 85, "Stored SQL Injection"),
        ("Numeric SQL Injection", "request_body", "regex", r"(?i)[0-9]+[\s]*[;'\"-]+", 82, 80, "Numeric SQL Injection"),
        ("String SQL Injection", "request_body", "regex", r"(?i)('|\")[\s]*(or|and)[\s]+'", 88, 86, "String SQL Injection"),
        ("SQL Comment Injection", "request_body", "regex", r"(?i)(--|#|\/\*.*\*\/)[\s]*", 85, 82, "SQL Comment Injection"),
        ("SQL Keyword Abuse", "request_body", "regex", r"(?i)\b(having|cast|convert|declare|exec|execute|master|xp_)\b", 84, 81, "SQL Keyword Abuse"),
        ("SQL Function Abuse", "request_body", "regex", r"(?i)\b(ascii|char|chr|concat|substring|substr|length|hex|ord)\b[\s]*\(", 83, 80, "SQL Function Abuse"),
        ("Database Error Pattern", "request_body", "regex", r"(?i)(mysql|postgresql|oracle|sqlite|syntax error|sqlstate)", 80, 78, "Database Error Pattern"),
        ("Database Fingerprinting", "request_body", "regex", r"(?i)(version\(\)|@@version|dbms_random|sqlite_version)", 88, 85, "Database Fingerprinting"),
        ("SQL Meta-character Detection", "request_body", "regex", r"(?i)(;|\|\||&&)[\s]*(select|drop|insert|update|delete|alter)", 90, 87, "SQL Meta-character Detection"),
    ]),
    ("xss", "Cross-Site Scripting", "high", "high", [
        ("Reflected XSS", "request_body", "regex", r"(?i)<script[\s>]", 96, 93, "Reflected XSS"),
        ("Stored XSS", "request_body", "regex", r"(?i)<script[\s>].*</script>", 92, 90, "Stored XSS"),
        ("DOM XSS", "request_body", "regex", r"(?i)(document\.(location|cookie|write)|window\.location|eval\(|innerHTML)", 88, 85, "DOM XSS"),
        ("HTML Injection", "request_body", "regex", r"(?i)<(iframe|object|embed|applet|form|img|svg)[\s>]", 90, 87, "HTML Injection"),
        ("Script Context Injection", "request_body", "regex", r"(?i)(javascript:|vbscript:|data:text/html)", 93, 90, "Script Context Injection"),
        ("HTML Attribute Injection", "request_body", "regex", r"(?i)[\"']\s*(on\w+)\s*=", 89, 86, "HTML Attribute Injection"),
        ("JavaScript Context Injection", "request_body", "regex", r"(?i)(alert|prompt|confirm)\s*\(", 86, 83, "JavaScript Context Injection"),
        ("CSS Context Injection", "request_body", "regex", r"(?i)(expression|@import|url\()", 82, 79, "CSS Context Injection"),
        ("Event Handler Injection", "request_body", "regex", r"(?i)\s(on(load|error|click|mouseover|focus))\s*=", 91, 88, "Event Handler Injection"),
        ("Encoded XSS", "request_body", "regex", r"(?i)(%3cscript|%253c|&lt;script|&#60;script)", 90, 87, "Encoded XSS"),
        ("Polyglot XSS", "request_body", "regex", r"(?i)(\"'<>/\\|javascript|onerror)", 85, 82, "Polyglot XSS"),
    ]),
    ("cmdi", "OS Command Injection", "critical", "critical", [
        ("OS Command Injection", "request_body", "regex", r"(?i)(;|&&|\|\|)[\s]*(ls|cat|id|whoami|pwd|uname|dir|type)", 95, 92, "OS Command Injection"),
        ("Shell Command Injection", "request_body", "regex", r"(?i)(/bin/sh|/bin/bash|cmd\.exe|powershell)", 94, 91, "Shell Command Injection"),
        ("Unix Command Injection", "request_body", "regex", r"(?i)(\$\(|`[^`]*`|;[\s]*(sh|bash|perl|python))", 92, 89, "Unix Command Injection"),
        ("Windows Command Injection", "request_body", "regex", r"(?i)(cmd[/\\\\]c|type[\s]+c:\\\\|dir[\s]+c:\\\\|net[\s]+user)", 90, 87, "Windows Command Injection"),
        ("Command Separator Detection", "request_body", "regex", r"(?i)(;|\|\||&&|\n|%0a|%0d)", 84, 80, "Command Separator Detection"),
        ("Command Substitution", "request_body", "regex", r"(?i)(\$\([^)]*\)|`[^`]*`)", 89, 86, "Command Substitution"),
        ("Shell Metacharacter Detection", "request_body", "regex", r"(?i)([;&|`$()<>])", 78, 75, "Shell Metacharacter Detection"),
        ("Command Execution Pattern", "request_body", "regex", r"(?i)\b(wget|curl|nc|netcat|tftp|scp|rsync)\b", 87, 84, "Command Execution Pattern"),
    ]),
    ("rce", "Remote Code Execution", "critical", "critical", [
        ("Code Injection", "request_body", "regex", r"(?i)\b(eval|assert|system|passthru|shell_exec|exec)\s*\(", 96, 93, "Code Injection"),
        ("Expression Injection", "request_body", "regex", r"(?i)(SpEL|OGNL|EL|#\{|\\$\{)[\s\w]", 88, 85, "Expression Injection"),
        ("Template Injection", "request_body", "regex", r"(?i)(\{\{|\{%|#{|<\%)", 90, 87, "Template Injection"),
        ("Server-side Template Injection", "request_body", "regex", r"(?i)({{[^}]+}}|\${{[^}]+}})", 89, 86, "Server-side Template Injection"),
        ("Dynamic Code Execution", "request_body", "regex", r"(?i)(__import__|os\.system|Runtime\.getRuntime|ProcessBuilder)", 93, 90, "Dynamic Code Execution"),
        ("Script Injection", "request_body", "regex", r"(?i)(python|cscript|wscript|perl|php[\s-]+r)", 85, 82, "Script Injection"),
        ("Runtime Execution Pattern", "request_body", "regex", r"(?i)(<%=|<\?php|jsp:|\.cfm)", 84, 81, "Runtime Execution Pattern"),
    ]),
    ("pt", "Path Traversal", "high", "high", [
        ("Directory Traversal", "url", "regex", r"(?i)(\.\./|\.\.\\|%2e%2e)", 97, 94, "Directory Traversal"),
        ("Encoded Traversal", "url", "regex", r"(?i)(%2e%2e%2f|%252e%252e|%c0%ae%c0%ae)", 93, 90, "Encoded Traversal"),
        ("Double Encoded Traversal", "url", "regex", r"(?i)(%25252e|%252e%252e)", 88, 85, "Double Encoded Traversal"),
        ("Path Normalization Bypass", "url", "regex", r"(?i)(\.\./\.\./|/\./|//)", 80, 76, "Path Normalization Bypass"),
        ("Unix Path Traversal", "url", "regex", r"(?i)(etc/passwd|etc/shadow|proc/self|home/[a-z])", 96, 93, "Unix Path Traversal"),
        ("Windows Path Traversal", "url", "regex", r"(?i)(c:\\\\windows|boot\.ini|win\.ini)", 92, 89, "Windows Path Traversal"),
        ("Sensitive File Access", "url", "regex", r"(?i)(\.env|\.git|web\.config|php\.ini|\.htaccess)", 94, 91, "Sensitive File Access"),
    ]),
    ("fi", "File Inclusion", "high", "high", [
        ("Local File Inclusion", "url", "regex", r"(?i)(include|require)[\s=]+(file|php://|/etc|/var/www)", 95, 92, "Local File Inclusion"),
        ("Remote File Inclusion", "url", "regex", r"(?i)(include[\s=]+(http|https|ftp):)", 96, 93, "Remote File Inclusion"),
        ("PHP File Inclusion", "url", "regex", r"(?i)(php://(filter|input|data|expect))", 94, 91, "PHP File Inclusion"),
        ("Dynamic Include Abuse", "url", "regex", r"(?i)(page=|\?file=|view=|path=).*(\.\.|/etc|php://)", 89, 86, "Dynamic Include Abuse"),
        ("Remote Resource Inclusion", "url", "regex", r"(?i)(http[s]?://[^\s]+(\.txt|\.php|\.html))", 85, 82, "Remote Resource Inclusion"),
    ]),
    ("xxe", "XXE / XML", "critical", "critical", [
        ("XML External Entity", "request_body", "regex", r"(?i)<!DOCTYPE[\s]+[^>]*\[[^>]*<!ENTITY", 97, 94, "XML External Entity"),
        ("External Entity Reference", "request_body", "regex", r"(?i)(&xxe;|&ent;|SYSTEM[\s]+\"|PUBLIC[\s]+\")", 93, 90, "External Entity Reference"),
        ("Parameter Entity", "request_body", "regex", r"(?i)<!ENTITY[\s]+%[^>]*>", 91, 88, "Parameter Entity"),
        ("External DTD", "request_body", "regex", r"(?i)(<!DOCTYPE[\s]+[^>]*SYSTEM[\s]+\")", 92, 89, "External DTD"),
        ("XML Entity Expansion", "request_body", "regex", r"(?i)<!ENTITY[\s]+[a-z][^>]*>.*(<!ENTITY[\s]+)", 90, 87, "XML Entity Expansion"),
        ("XML Bomb", "request_body", "regex", r"(?i)(lol|billion)[^\"]*\"[\s]*\"[^\"]*\"[\s]*\"[^\"]*\"", 89, 86, "XML Bomb"),
        ("XML Parser Abuse", "request_body", "regex", r"(?i)(<\?(xml|soap)|<![CDATA[|&lt;)", 85, 82, "XML Parser Abuse"),
        ("Malformed XML", "request_body", "regex", r"(?i)(<xml[^>]*>.*<xml|</>|<<)", 80, 76, "Malformed XML"),
    ]),
    ("ssrf", "SSRF", "critical", "critical", [
        ("Server-Side Request Forgery", "url", "regex", r"(?i)(url=|uri=|dest=|redirect=|proxy=)", 88, 84, "Server-Side Request Forgery"),
        ("Internal Network Access", "url", "regex", r"(?i)(10\.\d+\.\d+\.\d+|192\.168\.\d+\.\d+|172\.(1[6-9]|2[0-9]|3[01])\.\d+\.\d+)", 94, 91, "Internal Network Access"),
        ("Localhost Access", "url", "regex", r"(?i)(localhost|127\.0\.0\.1|0\.0\.0\.0)", 95, 92, "Localhost Access"),
        ("Private IP Access", "url", "regex", r"(?i)(10\.|192\.168\.|172\.(1[6-9]|2[0-9]|3[01])\.)", 93, 90, "Private IP Access"),
        ("Cloud Metadata Access", "url", "regex", r"(?i)(169\.254\.169\.254|metadata\.google\.internal|100\.100\.100\.200)", 97, 95, "Cloud Metadata Access"),
        ("Loopback Access", "url", "regex", r"(?i)(::1|\[::1\]|127\.[0-9]+\.[0-9]+\.[0-9]+)", 91, 88, "Loopback Access"),
        ("Internal DNS Access", "url", "regex", r"(?i)(\.internal$|\.local$|\.corp$|\.intranet$)", 87, 84, "Internal DNS Access"),
        ("URL Scheme Abuse", "url", "regex", r"(?i)(gopher://|file://|dict://|ftp://)", 92, 89, "URL Scheme Abuse"),
    ]),
    ("http", "HTTP Protocol", "medium", "medium", [
        ("HTTP Request Smuggling", "header", "regex", r"(?i)(transfer-encoding:\s*chunked[^\n]*\ncontent-length:|content-length:[^\n]*\ntransfer-encoding:\s*chunked)", 96, 93, "HTTP Request Smuggling"),
        ("HTTP Response Splitting", "header", "regex", r"(?i)(%0d%0a|\r\n\r\n)[^\n]*(location:|set-cookie:)", 90, 87, "HTTP Response Splitting"),
        ("Invalid HTTP Method", "method", "regex", r"(?i)^(trace|track|connect)$", 88, 85, "Invalid HTTP Method"),
        ("Invalid HTTP Version", "url", "regex", r"(?i)(HTTP/0\.|HTTP/1\.0|HTTP/9)", 87, 84, "Invalid HTTP Version"),
        ("Malformed Request", "url", "regex", r"(?i)(%00|%0a|%0d|\\\\)", 84, 80, "Malformed Request"),
        ("Header Anomaly", "header", "regex", r"(?i)(^\s|\\r\\n\\r\\n.*<)", 82, 78, "Header Anomaly"),
        ("Host Header Attack", "host", "regex", r"(?i)(localhost|127\.0\.0\.1|0\.0\.0\.0|\.internal)", 89, 86, "Host Header Attack"),
        ("Transfer-Encoding Anomaly", "header", "regex", r"(?i)transfer-encoding:\s*(identity|chunked[\s,])", 91, 88, "Transfer-Encoding Anomaly"),
        ("Content-Length Anomaly", "header", "regex", r"(?i)content-length:\s*[0-9]+[\s,]*content-length:", 93, 90, "Content-Length Anomaly"),
        ("Duplicate Header", "header", "regex", r"(?i)(content-length:|transfer-encoding:|host:).*\n.*(content-length:|transfer-encoding:|host:)", 90, 87, "Duplicate Header"),
        ("Request Desynchronization", "header", "regex", r"(?i)(te:\s*chunked|cl\.te|te\.cl)", 92, 89, "Request Desynchronization"),
    ]),
    ("hpp", "HTTP Parameter Pollution", "medium", "medium", [
        ("Duplicate Parameters", "url", "regex", r"(?i)([a-z_]+=[^&]+&[a-z_]+=)", 85, 81, "Duplicate Parameters"),
        ("Conflicting Parameters", "url", "regex", r"(?i)(role|admin|user|account)=[^&]+&(role|admin|user|account)=", 88, 85, "Conflicting Parameters"),
        ("Parameter Pollution", "url", "regex", r"(?i)(&|;|\|)[^=]*=[^&]*(&|\|;)[^=]*=", 83, 80, "Parameter Pollution"),
        ("Array Parameter Abuse", "url", "regex", r"(?i)([a-z_]+\[\]=|\[\]\[)", 86, 83, "Array Parameter Abuse"),
        ("Parameter Parsing Ambiguity", "url", "regex", r"(?i)(%26|%3b|%7c|%3d)", 84, 80, "Parameter Parsing Ambiguity"),
    ]),
    ("api", "API Security", "medium", "medium", [
        ("Invalid Content-Type", "header", "regex", r"(?i)content-type:\s*(text/plain|application/octet-stream)", 82, 78, "Invalid Content-Type"),
        ("JSON Structure Violation", "request_body", "regex", r"(?i)(\{\{|\}\}|:{2}|,{2})", 80, 76, "JSON Structure Violation"),
        ("Parameter Type Violation", "request_body", "regex", r"(?i)(\"(age|count|limit|offset)\":\s*\")", 81, 78, "Parameter Type Violation"),
        ("Excessive Parameters", "request_body", "regex", r"(?i)(\"[a-z_]+\":){10,}", 79, 75, "Excessive Parameters"),
        ("Unknown Parameter", "request_body", "regex", r"(?i)(\"(__proto__|constructor|prototype)\":)", 90, 87, "Unknown Parameter"),
        ("API Schema Violation", "request_body", "regex", r"(?i)(\"type\":\s*\"(int|str|bool)\"[\s,]*\"(expected|enum)\")", 85, 82, "API Schema Violation"),
    ]),
    ("file_upload", "File Upload", "high", "high", [
        ("Executable File Upload", "url", "regex", r"(?i)(\.exe|\.dll|\.so|\.bin|\.msi)$", 92, 89, "Executable File Upload"),
        ("Script File Upload", "url", "regex", r"(?i)(\.(php|php3|php5|phtml|jsp|asp|aspx|py|pl|rb|sh|cgi))$", 95, 92, "Script File Upload"),
        ("MIME Type Mismatch", "header", "regex", r"(?i)(multipart/form-data.*filename=\"[^\"]+\.(php|exe|sh))", 88, 85, "MIME Type Mismatch"),
        ("Double Extension", "url", "regex", r"(?i)(\.(php|exe|sh)\.[a-z0-9]{2,4})$", 90, 87, "Double Extension"),
        ("Archive Abuse", "url", "regex", r"(?i)(\.(zip|tar|gz|rar|7z))$", 84, 80, "Archive Abuse"),
        ("Oversized File", "header", "regex", r"(?i)content-length:\s*[0-9]{8,}", 86, 82, "Oversized File"),
        ("Malicious File Content", "request_body", "regex", r"(?i)(GIF89a|\\x89PNG|\\xff\\xd8\\xff)(.{0,100})(<\?php|exec\(|eval\()", 87, 84, "Malicious File Content"),
        ("Upload Path Manipulation", "url", "regex", r"(?i)((\.\./|%2e%2e/).*(upload|files|media))", 89, 86, "Upload Path Manipulation"),
    ]),
    ("ldap", "LDAP Injection", "high", "high", [
        ("LDAP Filter Injection", "request_body", "regex", r"(?i)(\*\)\(|\|\(|&\(|!\(|\(\|)", 91, 88, "LDAP Filter Injection"),
        ("LDAP Query Manipulation", "request_body", "regex", r"(?i)(cn=|uid=|mail=|userid=|sAMAccountName=).*(\*|\(|\))", 92, 89, "LDAP Query Manipulation"),
        ("LDAP Wildcard Abuse", "request_body", "regex", r"(?i)(\*\)\s*\(|\)\s*\|\s*\(|&\s*\()", 89, 86, "LDAP Wildcard Abuse"),
        ("LDAP Authentication Bypass", "request_body", "regex", r"(?i)(uid=.*\)\s*\(\|\(\s*uid=|\*\)\s*\|\(objectClass=\*)", 93, 90, "LDAP Authentication Bypass Pattern"),
    ]),
    ("nosql", "NoSQL Injection", "high", "high", [
        ("MongoDB Injection", "request_body", "regex", r"(?i)(\$where|\$ne|\$gt|\$regex|\$nin)", 95, 92, "MongoDB Injection"),
        ("NoSQL Operator Injection", "request_body", "regex", r"(?i)(\"\$[a-z]+\":)", 93, 90, "NoSQL Operator Injection"),
        ("JSON Operator Abuse", "request_body", "regex", r"(?i)((\$where|\$ne|\$gt|\$lt|\$regex)\s*:\s*\")", 92, 89, "JSON Operator Abuse"),
        ("Query Manipulation", "request_body", "regex", r"(?i)(find\s*\(\s*\{|\"query\":\s*\{)", 88, 85, "Query Manipulation"),
        ("NoSQL Authentication Bypass", "request_body", "regex", r"(?i)((\"password\"|\"user\").*(\$ne|\$gt|\$regex))", 94, 91, "NoSQL Authentication Bypass"),
    ]),
    ("ssti", "Server-Side Template Injection", "high", "high", [
        ("Template Expression", "request_body", "regex", r"(?i)(\{\{|\{%|#{|\<\%)", 90, 87, "Template Expression"),
        ("Template Directive", "request_body", "regex", r"(?i)(\{%[\s\w]+%\}|\{\{[\s\w]+\}\})", 89, 86, "Template Directive"),
        ("Expression Language", "request_body", "regex", r"(?i)(\$\{[\w.]+\}|\#\{[\w.]+\})", 91, 88, "Expression Language"),
        ("Template Variable Abuse", "request_body", "regex", r"(?i)(\{\{[^}]*(config|settings|env|request|self)\b)", 87, 84, "Template Variable Abuse"),
        ("Template Execution Pattern", "request_body", "regex", r"(?i)(\{\{[^}]*(system|exec|popen|os\.|import)\b)", 93, 90, "Template Execution Pattern"),
    ]),
    ("deserialization", "Deserialization", "critical", "critical", [
        ("Unsafe Deserialization", "request_body", "regex", r"(?i)(O:[0-9]+:\"|a:[0-9]+:\{|\{\"\\x00)", 92, 89, "Unsafe Deserialization"),
        ("Serialized Object Abuse", "request_body", "regex", r"(?i)(PHP_O:|JAVA_OBJECT|ACED0005|rO0AB)", 94, 91, "Serialized Object Abuse"),
        ("Object Injection", "request_body", "regex", r"(?i)(\$\_REQUEST|\$\_GET|\$\_POST|\.\.\.O:)", 88, 85, "Object Injection"),
        ("Java Deserialization", "request_body", "regex", r"(?i)(AC ED 00 05|rO0AB|\\xac\\xed)", 96, 93, "Java Deserialization"),
        ("PHP Deserialization", "request_body", "regex", r"(?i)(O:[0-9]+:\"[A-Z][^\"]+\":[0-9]+:)", 95, 92, "PHP Deserialization"),
        (".NET Deserialization", "request_body", "regex", r"(?i)(AAEAAAD|System\.Object|/System\.Web)", 90, 87, ".NET Deserialization"),
    ]),
    ("scanner", "Scanner Detection", "low", "low", [
        ("Vulnerability Scanner", "user_agent", "regex", r"(?i)(nikto|nessus|openvas|acunetix|w3af|sqlmap)", 95, 92, "Vulnerability Scanner"),
        ("Web Scanner", "user_agent", "regex", r"(?i)(gobuster|dirb|wfuzz|burpsuite|zap)", 94, 91, "Web Scanner"),
        ("Directory Scanner", "user_agent", "regex", r"(?i)(dirbuster|dirsearch|feroxbuster|ffuf)", 93, 90, "Directory Scanner"),
        ("Automated Reconnaissance", "user_agent", "regex", r"(?i)(masscan|nmap|zmap|nuclei)", 92, 89, "Automated Reconnaissance"),
        ("Security Testing Tool", "user_agent", "regex", r"(?i)(hydra|medusa|john|hashcat)", 91, 88, "Security Testing Tool"),
        ("Fingerprinting Tool", "user_agent", "regex", r"(?i)(whatweb|wappalyzer|builtwith|recon-ng)", 89, 86, "Fingerprinting Tool"),
        ("Automated Enumeration", "user_agent", "regex", r"(?i)(crawler4j|python-requests|scrapy|apache-httpclient)", 86, 83, "Automated Enumeration"),
    ]),
    ("bot", "Bot Protection", "medium", "medium", [
        ("Malicious Bot", "user_agent", "regex", r"(?i)(sqlmap|nikto|curl|wget|python|go-http-client)", 90, 87, "Malicious Bot"),
        ("Suspicious Bot", "user_agent", "regex", r"(?i)(bot|crawler|spider|scraper)", 82, 78, "Suspicious Bot"),
        ("Automated Client", "user_agent", "regex", r"(?i)(headless|phantomjs|selenium|puppeteer|playwright)", 91, 88, "Automated Client"),
        ("Headless Browser", "user_agent", "regex", r"(?i)(headlesschrome|chrome-headless|phantomjs)", 90, 87, "Headless Browser"),
        ("Browser Anomaly", "user_agent", "regex", r"(?i)(^$|^[-_]+$|^\d+$)", 84, 80, "Browser Anomaly"),
        ("High-frequency Client", "user_agent", "regex", r"(?i)(smtp|ftp|telnet|ssh|scanner|spider)", 80, 76, "High-frequency Client"),
    ]),
    ("auth", "Authentication Attacks", "high", "high", [
        ("Credential Stuffing", "request_body", "regex", r"(?i)(\"password\":\s*\"){2,}", 90, 87, "Credential Stuffing"),
        ("Authentication Bypass", "request_body", "regex", r"(?i)(('|\")\s*(or|and)\s+['\"]?\s*['\"]?\s*=|password.*=.*'|' OR 1=1)", 95, 92, "Authentication Bypass"),
        ("Password Attack", "request_body", "regex", r"(?i)(password|passwd|pwd)=[^&]*(admin|123456|qwerty|password)", 88, 85, "Password Attack"),
        ("Session Anomaly", "header", "regex", r"(?i)(cookie:\s*.*(sessionid|jsessionid|phpsessid)=.{1,2}$)", 85, 82, "Session Anomaly"),
        ("Authentication Flood", "url", "regex", r"(?i)(/login|/auth).*(bot|spider|scanner)", 83, 79, "Authentication Flood"),
    ]),
    ("session", "Session Security", "medium", "medium", [
        ("Session Fixation", "header", "regex", r"(?i)(cookie:\s*.*(sessionid|phpsessid|jsessionid)=[^;]+)", 89, 86, "Session Fixation"),
        ("Session Token Anomaly", "header", "regex", r"(?i)(session(token|id)=[^;]{1,4}$)", 87, 84, "Session Token Anomaly"),
        ("Cookie Manipulation", "header", "regex", r"(?i)(cookie:\s*.*(isadmin|role|user)=1|admin=true)", 91, 88, "Cookie Manipulation"),
        ("Invalid Session", "header", "regex", r"(?i)(session(id|token)=\s*[\"']?\s*[\"']?)", 86, 83, "Invalid Session"),
        ("Session Replay", "header", "regex", r"(?i)(cookie:\s*.*(session|token)=[^;]{5,};\s*cookie:\s*.*)", 84, 80, "Session Replay"),
    ]),
    ("csrf", "CSRF", "medium", "medium", [
        ("Missing CSRF Token", "request_body", "regex", r"(?i)(method:\s*(post|put|delete).*(?!csrf|token))", 78, 74, "Missing CSRF Token"),
        ("Invalid CSRF Token", "request_body", "regex", r"(?i)(csrf(token|_token)?=\s*[\"']?[\"']?|_token=\s*$)", 88, 85, "Invalid CSRF Token"),
        ("CSRF Token Mismatch", "request_body", "regex", r"(?i)((csrf|_token|token)=[^&]+&.*(csrf|_token|token)=[^&]+)", 87, 84, "CSRF Token Mismatch"),
        ("Cross-origin Request Anomaly", "header", "regex", r"(?i)(origin:\s*https?://(?!([a-z0-9.-]+\.)?(example\.com|yourdomain\.com)))", 82, 78, "Cross-origin Request Anomaly"),
    ]),
    ("info_disclosure", "Information Disclosure", "medium", "medium", [
        ("Source Code Exposure", "url", "regex", r"(?i)(\.(php|js|py|rb|java|c)$|/source|/src)", 90, 87, "Source Code Exposure"),
        ("Configuration File Exposure", "url", "regex", r"(?i)(config\.(php|js|json|xml)|application\.(yml|yaml|properties))", 89, 86, "Configuration File Exposure"),
        ("Backup File Exposure", "url", "regex", r"(?i)(\.(bak|old|orig|save|swp|tmp))$", 88, 85, "Backup File Exposure"),
        ("Debug Information", "url", "regex", r"(?i)(/debug|/trace|/phpinfo|/info\.php|/status)", 91, 88, "Debug Information"),
        ("Stack Trace Exposure", "request_body", "regex", r"(?i)(Traceback \(|at [a-z]+\.[a-z]+\.|\.java:[0-9]+|\.cs:[0-9]+)", 87, 84, "Stack Trace Exposure"),
        ("Environment File Exposure", "url", "regex", r"(?i)(\.env|\.env\.local|\.env\.prod|\.env\.development)", 95, 92, "Environment File Exposure"),
        ("Version Disclosure", "header", "regex", r"(?i)(server:\s*|x-powered-by:\s*|x-aspnet-version:\s*)[a-z0-9/.-]+", 86, 82, "Version Disclosure"),
        ("Sensitive File Access", "url", "regex", r"(?i)(/etc/passwd|/etc/shadow|web\.config|php\.ini|\.git/config)", 94, 91, "Sensitive File Access"),
    ]),
    ("resource_discovery", "Resource Discovery", "low", "low", [
        ("Admin Panel Discovery", "url", "regex", r"(?i)(/admin|/administrator|/wp-admin|/cpanel|/dashboard)", 89, 86, "Admin Panel Discovery"),
        ("Login Panel Discovery", "url", "regex", r"(?i)(/login|/signin|/signup|/auth|/account)", 85, 82, "Login Panel Discovery"),
        ("Backup File Discovery", "url", "regex", r"(?i)(\.(zip|tar|gz|bak|old|sql))$", 87, 84, "Backup File Discovery"),
        ("Configuration Discovery", "url", "regex", r"(?i)(\.git|\.svn|\.env|web\.config|config\.php)", 92, 89, "Configuration Discovery"),
        ("Git Repository Discovery", "url", "regex", r"(?i)(/\.git/(config|HEAD|index)|/\.git$|/\.svn/)", 95, 92, "Git Repository Discovery"),
        ("API Discovery", "url", "regex", r"(?i)(/api/|/swagger|/openapi|/graphql|/v[0-9]+/)", 88, 85, "API Discovery"),
        ("Debug Endpoint Discovery", "url", "regex", r"(?i)(/debug|/healthz|/actuator|/metrics|/status)", 86, 83, "Debug Endpoint Discovery"),
    ]),
    ("request_anomaly", "Request Anomaly", "medium", "medium", [
        ("URL Too Long", "url", "regex", r"^.{2000,}$", 90, 87, "URL Too Long"),
        ("Header Too Large", "header", "regex", r"^.{5000,}$", 89, 86, "Header Too Large"),
        ("Body Too Large", "request_body", "regex", r"^.{1000000,}$", 88, 85, "Body Too Large"),
        ("Too Many Parameters", "url", "regex", r"([?&][a-z_]+=){30,}", 87, 84, "Too Many Parameters"),
        ("Invalid Encoding", "url", "regex", r"(?i)(%[0-9a-f]{1}$|%[g-z]|%(?![0-9a-f]{2}))", 84, 80, "Invalid Encoding"),
        ("Invalid Characters", "url", "regex", r"(?i)([\x00-\x1f\x7f])", 86, 82, "Invalid Characters"),
        ("Invalid Content-Length", "header", "regex", r"(?i)content-length:\s*0{5,}|content-length:\s*[^0-9]", 88, 85, "Invalid Content-Length"),
        ("Malformed JSON", "request_body", "regex", r"(?i)(\}[^\s,}\]]|\{:|\[:|,\s*[}\]])", 83, 80, "Malformed JSON"),
        ("Protocol Violation", "url", "regex", r"(?i)(%[0-9a-f]{1}\b|\\x[0-9a-f]{2})", 82, 78, "Protocol Violation"),
    ]),
    ("ip", "IP / Reputation", "high", "high", [
        ("IP Blocklist", "source_ip", "cidr_match", r"0.0.0.0/0", 90, 88, "IP Blocklist"),
        ("CIDR Block", "source_ip", "cidr_match", r"0.0.0.0/0", 88, 85, "CIDR Block"),
        ("Proxy Detection", "header", "regex", r"(?i)(x-forwarded-for:\s*[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+|via:\s*[a-z0-9-]+)", 85, 82, "Proxy Detection"),
        ("Threat Intelligence Match", "source_ip", "regex", r"^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$", 87, 84, "Threat Intelligence Match"),
    ]),
    ("geo", "Geo Security", "low", "low", [
        ("Country Blocklist", "source_ip", "regex", r"^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$", 90, 87, "Country Blocklist"),
        ("Country Allowlist", "source_ip", "regex", r"^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$", 90, 87, "Country Allowlist"),
        ("High-risk Location Policy", "source_ip", "regex", r"^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$", 88, 85, "High-risk Location Policy"),
    ]),
    ("rate_limit", "Rate Limiting", "medium", "medium", [
        ("Global Rate Limit", "method", "equals", r"GET", 85, 82, "Global Rate Limit"),
        ("IP Rate Limit", "source_ip", "regex", r"^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$", 86, 83, "IP Rate Limit"),
        ("Endpoint Rate Limit", "url", "regex", r"(?i)(/login|/api/|/search|/auth)", 84, 80, "Endpoint Rate Limit"),
        ("Login Rate Limit", "url", "regex", r"(?i)(/login|/signin|/auth)", 87, 84, "Login Rate Limit"),
        ("API Rate Limit", "url", "regex", r"(?i)/api/", 85, 82, "API Rate Limit"),
        ("Burst Detection", "request_body", "regex", r"(?i)(burst|flood|rapid)", 78, 74, "Burst Detection"),
        ("Request Flood Detection", "url", "regex", r"(?i)(/flood|/attack|/stress)", 80, 76, "Request Flood Detection"),
    ]),
]


def gen():
    out = []
    out.append("-- 0027_seed_rule_catalog.up.sql")
    out.append("-- Full 27-category WAF rule catalog with F5-style signature metadata.")
    out.append("")
    out.append("-- Managed rule groups (one per category).")
    for i, (cat, name, sev, risk, sigs) in enumerate(CATS, start=1):
        mid = f"01MRCAT{i:03d}"
        out.append(f"INSERT INTO managed_rules (id, organization_id, name, category, enabled, sensitivity, paranoia_level, action, anomaly_threshold, status) VALUES")
        out.append(f"  ('{mid}', '{ORG}', '{sql(name)}', '{cat}', true, 'low', 1, 'block', 5, 'active')")
        out.append("ON CONFLICT DO NOTHING;")
        out.append("")

    out.append("-- Signature rules with metadata + conditions.")
    n = 0
    for ci, (cat, name, sev, risk, sigs) in enumerate(CATS, start=1):
        for si, (sname, field, op, pattern, acc, conf, atk) in enumerate(sigs, start=1):
            n += 1
            rid = f"01MRRUL{n:04d}"
            rule_id = f"ARB-{cat.upper()}-{si:03d}"
            priority = ci * 100 + si
            desc = f"{sname} - {name} signature"
            out.append(f"INSERT INTO rules (id, organization_id, rule_id, name, description, type, category, severity, priority, action, status, logic, version, attack_type, pattern_type, accuracy, risk, confidence, source, staging, remediation) VALUES")
            out.append(f"  ('{rid}', '{ORG}', '{rule_id}', '{sql(sname)}', '{sql(desc)}', 'managed', '{cat}', '{sev}', {priority}, 'block', 'active', 'AND', 1, '{sql(atk)}', '{op}', {acc}, '{risk}', {conf}, 'ariba-core', false, '')")
            out.append("ON CONFLICT DO NOTHING;")
            p = sql(pattern)
            out.append(f"INSERT INTO waf_rule_conditions (id, rule_id, group_id, field, operator, value, transformation, case_sensitive) VALUES")
            out.append(f"  ('{rid}C', '{rid}', 0, '{field}', '{op}', '{p}', 'lowercase', false);")
            out.append("")
    out.append(f"-- {len(CATS)} categories, {n} signature rules seeded.")
    return "\n".join(out)


if __name__ == "__main__":
    path = "apps/control-api/migrations/0027_seed_rule_catalog.up.sql"
    with open(path, "w") as f:
        f.write(gen())
    total = sum(len(s) for _, _, _, _, s in CATS)
    print(f"Generated: {len(CATS)} categories, {total} signature rules -> {path}")
