package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Wazuh syslog/JSON forwarder (Phase 2 integration, master plan §6.8).
// Reads JSON-lines security events from stdin and forwards to Wazuh via
// syslog over UDP/TCP/TLS (P2.29). When no --host is configured it falls
// back to writing to stdout. Sensitive fields are redacted before sending
// (P2.30).

// wazuhEvent is the Wazuh-alert payload format expected by Wazuh agents/syslog
// ingestion (rule.id, rule.level, data, full_log).
type wazuhEvent struct {
	Timestamp string         `json:"timestamp"`
	Rule      wazuhRule      `json:"rule"`
	Agent     wazuhAgent     `json:"agent"`
	Manager   wazuhManager   `json:"manager"`
	Data      map[string]any `json:"data"`
	FullLog   string         `json:"full_log"`
	Location  string         `json:"location"`
}

type wazuhRule struct {
	ID      string   `json:"id"`
	Level   int      `json:"level"`
	Groups  []string `json:"groups"`
	Matches []string `json:"matches,omitempty"`
}

type wazuhAgent struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type wazuhManager struct {
	Name string `json:"name"`
}

func main() {
	managerName := flag.String("manager", envOr("WAZUH_MANAGER", "ariba-shield"), "Wazuh manager name")
	agentName := flag.String("agent", envOr("WAZUH_AGENT", ""), "agent name (defaults to GATEWAY_ID)")
	output := flag.String("output", "stdout", "output mode: stdout | syslog")
	transport := flag.String("transport", envOr("WAZUH_TRANSPORT", "tcp"), "syslog transport: udp | tcp | tls")
	host := flag.String("host", envOr("WAZUH_HOST", ""), "syslog host:port (empty = write to stdout)")
	rfc := flag.String("rfc", envOr("WAZUH_RFC", "3164"), "syslog format: 3164 | 5424")
	tlsInsecure := flag.Bool("tls-insecure", envBool("WAZUH_TLS_INSECURE", false), "skip TLS server verification (TLS transport only)")
	tlsCA := flag.String("tls-ca", envOr("WAZUH_TLS_CA", ""), "path to a PEM CA bundle to verify the syslog server (TLS transport only)")
	flag.Parse()

	agent := *agentName
	if agent == "" {
		agent = os.Getenv("GATEWAY_ID")
	}

	switch *transport {
	case "udp", "tcp", "tls":
	default:
		log.Fatalf("invalid transport %q: must be udp, tcp or tls", *transport)
	}
	if *rfc != "3164" && *rfc != "5424" {
		log.Fatalf("invalid rfc %q: must be 3164 or 5424", *rfc)
	}

	// P2.29: an actual syslog socket. Without --host we keep the stdout modes.
	var sender *netSender
	if *host != "" {
		s, err := newNetSender(*transport, *host, *tlsCA, *tlsInsecure)
		if err != nil {
			log.Fatalf("syslog transport: %v", err)
		}
		defer s.close()
		sender = s
		log.Printf("forwarding syslog (%s, rfc%s) to %s", *transport, *rfc, *host)
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var raw map[string]any
		if err := json.Unmarshal(line, &raw); err != nil {
			log.Printf("skip malformed event: %v", err)
			continue
		}

		// P2.30: redact sensitive fields before they leave this hop.
		redacted := redactEvent(raw).(map[string]any)

		ev, err := toWazuhEvent(redacted, *managerName, agent)
		if err != nil {
			log.Printf("transform failed: %v", err)
			continue
		}
		out, err := json.Marshal(ev)
		if err != nil {
			continue
		}

		rfcVal, _ := strconv.Atoi(*rfc)
		msg := formatSyslog(rfcVal, syslogPRI(ev.Rule.Level), agent, ev.Rule.ID, out)

		if sender != nil {
			if err := sender.send([]byte(msg)); err != nil {
				log.Printf("syslog send failed: %v", err)
			}
			continue
		}

		if *output == "syslog" {
			// Backward-compatible: formatted line on stdout.
			fmt.Fprintln(os.Stdout, msg)
		} else {
			fmt.Fprintln(os.Stdout, string(out))
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// formatSyslog renders a syslog message in RFC3164 or RFC5424 shape.
func formatSyslog(rfc, pri int, hostname, msgID string, msg []byte) string {
	if rfc == 5424 {
		ts := time.Now().UTC().Format("2006-01-02T15:04:05Z")
		return fmt.Sprintf("<%d>1 %s %s waf %d %s - %s",
			pri, ts, syslogHost(hostname), os.Getpid(), syslogMsgID(msgID), msg)
	}
	// RFC3164: no year, day is space-padded.
	ts := time.Now().Format("Jan _2 15:04:05")
	return fmt.Sprintf("<%d>%s %s waf[%d]: %s", pri, ts, syslogHost(hostname), os.Getpid(), msg)
}

// syslogPRI maps a Wazuh rule level to a syslog priority using the LOG_USER
// facility (1). Levels follow the default wazuh level mapping in toWazuhEvent.
func syslogPRI(level int) int {
	facility := 1 // LOG_USER
	var severity int
	switch {
	case level >= 10:
		severity = 2 // LOG_CRIT
	case level >= 8:
		severity = 3 // LOG_ERR
	case level >= 5:
		severity = 5 // LOG_NOTICE
	default:
		severity = 6 // LOG_INFO
	}
	return facility*8 + severity
}

func syslogHost(agent string) string {
	if agent != "" {
		return agent
	}
	host, err := os.Hostname()
	if err != nil {
		return "waf"
	}
	return host
}

func syslogMsgID(id string) string {
	if id == "" {
		return "-"
	}
	return strings.ReplaceAll(id, " ", "_")
}

func toWazuhEvent(raw map[string]any, manager, agent string) (*wazuhEvent, error) {
	sev := asString(raw["severity"])
	level := 0
	switch strings.ToLower(sev) {
	case "critical", "high":
		level = 12
	case "medium", "warn":
		level = 8
	case "low", "info":
		level = 3
	default:
		level = 6
	}

	reason := asString(raw["reason"])
	ruleID := firstString(raw["rule_ids"])

	groups := []string{"ariba_shield", "waf"}
	if reason != "" {
		groups = append(groups, strings.ToLower(reason))
	}

	// Build a human-readable full_log for Wazuh dashboards.
	fullLog := fmt.Sprintf("[%s] method=%s path=%s client_ip=%s reason=%s rule=%s",
		reason, asString(raw["method"]), asString(raw["path"]), asString(raw["client_ip"]), reason, ruleID)

	return &wazuhEvent{
		Timestamp: asString(raw["timestamp"]),
		Rule: wazuhRule{
			ID:      ruleID,
			Level:   level,
			Groups:  groups,
			Matches: stringSlice(raw["rule_ids"]),
		},
		Agent:    wazuhAgent{Name: agent},
		Manager:  wazuhManager{Name: manager},
		Data:     raw,
		FullLog:  fullLog,
		Location: "waf-engine",
	}, nil
}

func firstString(v any) string {
	if arr, ok := v.([]any); ok && len(arr) > 0 {
		return fmt.Sprintf("%v", arr[0])
	}
	if arr, ok := v.([]string); ok && len(arr) > 0 {
		return arr[0]
	}
	return ""
}

func stringSlice(v any) []string {
	var out []string
	switch arr := v.(type) {
	case []any:
		for _, item := range arr {
			out = append(out, fmt.Sprintf("%v", item))
		}
	case []string:
		out = arr
	}
	return out
}

func asString(v any) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envBool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	switch strings.ToLower(v) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}
