# Design Artifacts

Standalone UI prototypes and design references for the Ariba Shield WAF console.

## Contents

| File | Purpose |
|---|---|
| `ariba-shield-waf-dashboard.html` | Standalone, self-contained HTML dashboard prototype. Single-file artifact with embedded CSS/JS — open directly in a browser. Not a production screen; used for visual design reference and layout iteration. |

## Notes

- These artifacts are **not** part of the production build and should not be confused with the Next.js console in `apps/console-web/`.
- The dashboard HTML references external CDN resources (react, chart.js, etc.) via its CSP; it requires internet access to render.