package compiler

import (
	"bytes"
	"fmt"
)

// RenderWAFConfig produces the Coraza configuration from the policy document.
func RenderWAFConfig(doc *PolicyDocument) (string, error) {
	if doc.WAF == nil || !doc.WAF.Enabled {
		// If WAF is not enabled or config missing, generate a detect-only minimal baseline.
		return renderBaseline(), nil
	}

	var buf bytes.Buffer
	buf.WriteString("# Ariba Shield WAF — Dynamically Generated Configuration\n")
	buf.WriteString(fmt.Sprintf("# config_id=%s\n\n", doc.ConfigID))

	// Core engine directives
	buf.WriteString("Include @coraza.conf-recommended\n\n")
	buf.WriteString("SecRuleEngine On\n")
	buf.WriteString("SecRequestBodyAccess On\n")
	buf.WriteString("SecRequestBodyLimit 13107200\n")
	buf.WriteString("SecRequestBodyNoFilesLimit 131072\n")
	buf.WriteString("SecResponseBodyAccess Off\n")
	buf.WriteString("SecPcreMatchLimit 100000\n")
	buf.WriteString("SecPcreMatchLimitRecursion 100000\n\n")

	// CRS Setup
	buf.WriteString("# --- CRS Setup ---\n")
	buf.WriteString("Include @crs-setup.conf.example\n\n")

	// Paranoia Level
	pl := doc.WAF.ParanoiaLevel
	if pl < 1 || pl > 4 {
		pl = 1 // Default PL
	}
	buf.WriteString(fmt.Sprintf("SecAction \"id:900000,phase:1,nolog,pass,t:none,setvar:tx.paranoia_level=%d\"\n\n", pl))

	// Include all CRS rules if any managed rule is enabled
	hasManaged := false
	for _, mr := range doc.WAF.ManagedRules {
		if mr.Enabled {
			hasManaged = true
			break
		}
	}
	
	if hasManaged {
		buf.WriteString("# --- OWASP CRS Managed Rules ---\n")
		// For now we include all rules. Later we can filter by specific files
		// if we only want certain categories.
		buf.WriteString("Include @owasp_crs/*.conf\n\n")
	}

	// Anomaly Score Threshold
	thresh := doc.WAF.AnomalyThreshold
	if thresh <= 0 {
		thresh = 5 // Default threshold
	}
	buf.WriteString("# --- Anomaly Score Threshold ---\n")
	buf.WriteString(fmt.Sprintf("SecAction \"id:900001,phase:1,nolog,pass,t:none,setvar:tx.blocking_anomaly_score=%d\"\n", thresh))

	return buf.String(), nil
}

func renderBaseline() string {
	return `# Ariba Shield WAF — Minimal Baseline Configuration (WAF Disabled or Missing)
Include @coraza.conf-recommended
SecRuleEngine DetectionOnly
SecRequestBodyAccess On
SecRequestBodyLimit 13107200
SecRequestBodyNoFilesLimit 131072
SecResponseBodyAccess Off
SecPcreMatchLimit 100000
SecPcreMatchLimitRecursion 100000

Include @crs-setup.conf.example
SecAction "id:900000,phase:1,nolog,pass,t:none,setvar:tx.paranoia_level=1"
Include @owasp_crs/*.conf
SecAction "id:900001,phase:1,nolog,pass,t:none,setvar:tx.blocking_anomaly_score=5"
`
}
