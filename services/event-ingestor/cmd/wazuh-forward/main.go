package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// Wazuh syslog/JSON forwarder (Phase 2 integration, master plan §6.8).
// Reads JSON-lines security events from stdin and forwards to Wazuh via
// syslog over TCP/UDP. Forward-compatible with the ADR-002 event shape.

// wazuhEvent is the Wazuh-alert payload format expected by Wazuh agents/syslog
// ingestion (rule.id, rule.level, data, full_log).
type wazuhEvent struct {
	Timestamp  string         `json:"timestamp"`
	Rule       wazuhRule      `json:"rule"`
	Agent      wazuhAgent     `json:"agent"`
	Manager    wazuhManager   `json:"manager"`
	Data       map[string]any `json:"data"`
	FullLog    string         `json:"full_log"`
	Location   string         `json:"location"`
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
	managerName := flag.String("manager", "ariba-shield", "Wazuh manager name")
	agentName := flag.String("agent", "", "agent name (defaults to GATEWAY_ID)")
	output := flag.String("output", "stdout", "output mode: stdout | syslog")
	flag.Parse()

	agent := *agentName
	if agent == "" {
		agent = os.Getenv("GATEWAY_ID")
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

		ev, err := toWazuhEvent(raw, *managerName, agent)
		if err != nil {
			log.Printf("transform failed: %v", err)
			continue
		}
		out, err := json.Marshal(ev)
		if err != nil {
			continue
		}

		if *output == "syslog" {
			// syslog RFC5424-ish: <PRI>TIMESTAMP HOST TAG [id] message
			pri := 13 // LOG_USER | LOG_NOTICE
			fmt.Fprintf(os.Stdout, "<%d>%s %s waf[%s]: %s\n",
				pri, time.Now().UTC().Format("2006-01-02T15:04:05Z"), agent, ev.Rule.ID, out)
		} else {
			fmt.Fprintln(os.Stdout, string(out))
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
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