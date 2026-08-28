# Reference Hardware and Performance Method — Phase 0

- **Status:** Draft v0.1
- **Date:** 2026-08-28
- **Owner:** Platform Architect / SRE
- **References:** `../enterprise_waf_development_master_plan.md` (§13, §20), `../phase0_srs.md` (NFR-0.1-001..005), `adr-001-initial-architecture.md`

---

## 1. Purpose

Performance targets in the SRS and master plan are **only meaningful against a defined, published reference environment** (master plan §13: "never publish unqualified throughput"; rule 7). This document fixes the reference hardware, workload definitions, and measurement methodology so that:

- Release 0.1 performance claims are reproducible.
- Later releases (especially the Phase 9 Rust gateway) are compared **on the same bench**, not on vibes.
- Regression is measured per commit on the traffic path.

---

## 2. Reference hardware (benchmark server)

| Component | Specification (target) | Notes |
|---|---|---|
| CPU | 2× 16-core (32 threads/core group), ≥ 3.0 GHz boost, Intel Xeon or AMD EPYC | Hyperthreading **disabled** for consistent results |
| RAM | 128 GB DDR4/DDR5 | Enough to avoid page cache as a variable |
| NIC | 2× 25 GbE (or 2× 10 GbE minimum), RSS enabled | One port for client-facing traffic |
| Storage | NVMe ≥ 1 TB | Local write-ahead config store + logs |
| OS | Ubuntu Server LTS (same major version as reference deployment) | Fresh install, no background load |
| Kernel | LTS default, tuned: `net.core.somaxconn`, `net.ipv4.tcp_max_syn_backlog`, TCP keepalive defaults | Documented `sysctl` snapshot captured with results |
| BIOS/firmware | Performance profile, C-states allowed, no power caps | Documented with results |

**Publish with every result:** exact CPU model + core count + base/boost, RAM size/type, NIC model + negotiated speed, OS + kernel version, OpenResty/Nginx + OpenSSL version, TLS version + cipher suite, and the full benchmark command + flags. A result without this context is not a result.

---

## 3. Workload definitions

Every measurement is one cell in a matrix (master plan §13 lists the axes):

### 3.1 Fixed base workload
- HTTP/1.1, keep-alive enabled, TLS 1.3 (TLS_AES_128_GCM_SHA256), 50 concurrent connections, steady-state, warmup 60 s, measurement 180 s.
- Small body: request 1 KB / response 4 KB. Medium body: request 4 KB / response 64 KB. Large body: request 64 KB / response 1 MB.
- Load generator on a **separate machine** on the same L2/L3 segment (never co-located with the gateway; the generator must not be the bottleneck — verify generator headroom).

### 3.2 Measurement axes (each a full matrix cell)

| Axis | Values |
|---|---|
| Transport | plain HTTP vs TLS |
| HTTP version | HTTP/1.1 vs HTTP/2 |
| Connection reuse | keep-alive vs new connection per request |
| Body size | small / medium / large |
| Enforcement | detection-only vs blocking (Phase 2+) |
| CRS sensitivity | paranoia levels 1–4 (Phase 2+) |
| Payload type | JSON / XML / multipart (Phase 2+) |
| Response inspection | off vs on (Phase 2+) |
| File upload scan | off vs on (Phase 2+) |
| Traffic mix | normal vs attack-heavy (Phase 2+, documented ratio) |

### 3.3 Reported metrics
- Requests/sec (steady-state mean, p50/p95/p99)
- Throughput (MB/s both directions)
- **Added latency** p50/p95/p99: difference between gateway-in-path and direct-to-backend (same workload) — this is the honest metric for "WAF cost"
- TLS handshakes/sec
- CPU % per core (per-process), memory/request
- Event drops (must be 0 at rated load; NFR-0.1-005)

### 3.4 Rated load definition
Rated load = the sustained requests/sec at which added p99 ≤ 10 ms **and** event drops = 0 **and** CPU headroom ≥ 30% on the busiest core. Publish this as the official "rated load" per configuration, with all conditions.

---

## 4. Methodology

1. **Baseline:** measure direct client→backend (no gateway) to get the reference latency/throughput curve.
2. **Calibrate generator:** confirm the generator sustains ≥ 2× the gateway's expected capacity before trusting its numbers.
3. **Warmup:** 60 s at target load; discard the first 30 s of samples.
4. **Measurement:** 180 s steady state; record p50/p95/p99, req/s, errors, event drops.
5. **Repeat 3×**; report median of runs; if run-to-run spread > 5%, re-check environment (noise, thermal, co-tenancy) and re-run.
6. **Regression rule:** a change is a performance regression if p99 added latency or req/s degrades > 5% at the same rated load on the same bench; blocks release per §14.3.
7. **Tooling:** `wrk2`/`h2load`/`ghz`-class generators; `perf`/`flamegraph` and allocation profiling for hotspots; deterministic replay for bug reproduction; fuzz coverage tracked for parsers.

---

## 5. Reporting template

A benchmark result must be published in this shape:

```text
Environment: [CPU model xN, RAM, NIC, OS+kernel, OpenResty/OpenSSL versions, sysctl snapshot]
Workload:    [HTTP/1.1|h2, TLS|plain, keep-alive|new-conn, body small|med|large, ...]
Enforcement: [detection-only|blocking, CRS paranoia N]
Result:
  req/s          = ...
  added p50/p99  = ... ms
  handshakes/s   = ...
  CPU/mem        = ...
  event drops    = ...
  rated load     = ...  (load where p99<=10ms, drops=0, CPU headroom>=30%)
Command:        [full generator invocation]
```

---

## 6. Release 0.1 baseline scope

Before Release 0.1 ships, publish the base matrix for: plain HTTP + TLS, HTTP/1.1 keep-alive, small/medium/large bodies, detection-off (proxy only). This becomes the benchmark the Phase 2 WAF path and Phase 9 Rust gateway must match or beat.

---

## 7. Open items

- OQ-2 (from SRS): exact reference-server procurement/spec sign-off — owner: architect.
- OQ-4 (from SRS): confirm load-generator hardware availability for the bench.
- Publish the `sysctl` + tuning snapshot as a tracked file in the repo (`deployments/`).
