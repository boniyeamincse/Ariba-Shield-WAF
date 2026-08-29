# ADR-007 — High-Volume Event Transport

- **Status:** Proposed
- **Date:** 2026-08-29
- **Owner:** Platform Architect
- **References:** `../enterprise_waf_development_master_plan.md` (§4, §6.8, §8), `adr-006-logging-event-pipeline.md`, `adr-005-coraza-integration.md`

---

## Context

The current event pipeline (ADR-006) uses stdout JSON-lines → best-effort sidecar `event-ingestor` → PostgreSQL. This is sufficient for Release 0.1–0.5 volumes but will not scale to enterprise-grade throughput:

1. PostgreSQL is not a high-volume event sink; `security_events` concurrent insert/query pressure grows with traffic.
2. The sidecar is single-process, single-threaded; a restarted ingestor loses in-flight events.
3. No backpressure or persistent queue — a DB outage causes unbounded memory growth (mitigated by P2.21 cap, but events are still dropped).
4. No replay or redelivery for downstream consumers (Wazuh, SIEM, SOAR, analytics).

The master plan (§4) prescribes: "Redis Streams initially → NATS JetStream or Kafka according to measured scale."

---

## Decision

Transition the event pipeline through two phases:

### Phase A — Redis Streams (current, Release 0.5–1.0)

- **When:** Immediately after the DB-backed pipeline proves insufficient under load.
- **What:** The gateway/engine writes events to a **Redis Stream** (`XADD`) instead of stdout. The `event-ingestor` reads from the stream (`XREADGROUP`), batches, and writes to PostgreSQL.
- **Why Redis:** Already deployed as a dependency (sessions, cache). Streams provide at-least-once delivery, consumer groups, and bounded memory (configurable `MAXLEN`).
- **Limitations:** Redis is not a durable long-term store; stream length is bounded. Events older than `MAXLEN` are lost. No replay for backfill. No partitioning.

### Phase B — NATS JetStream or Kafka (Release 1.5+)

- **When:** Fleet scales beyond 10 gateways or event volume exceeds Redis stream capacity.
- **What:** Replace Redis Streams with NATS JetStream (preferred for Kubernetes-native) or Kafka (preferred for existing SIEM/Kafka investment).
- **Why:** Durable, replayable, partitioned, multi-consumer event streams. Enables independent scaling of ingestion, analytics, and SIEM forwarding.
- **ClickHouse:** Replace PostgreSQL as the security-event analytics store. All events land in ClickHouse for sub-second analytical queries; PostgreSQL retains config, audit, and metadata (master plan §8).

### D1 — Event schema stability

The ADR-002 event schema v0 **must not change** across transport phases. The transport layer is opaque to event producers:
- Gateway/engine: emit JSON-lines to stdout (Phase A) or Redis Stream `XADD` (Phase B).
- Event-ingestor: reads from stdin (Phase A) or Redis Stream consumer group (Phase B). Same batch/flush/retry logic.
- Consumers (Wazuh forwarder, console, ClickHouse ETL): read from PostgreSQL or directly from the stream.

### D2 — Backpressure and buffer

Ingestor must:
- Cap in-memory buffer (10k events, P2.21).
- Apply exponential backoff on consumer-side failure (100ms–30s, P2.22).
- When the buffer is full, **drop the oldest events** (never block the producer).

### D3 — Migration path

1. Deploy Redis Streams alongside the stdout path (dual-write or gateway sends to both).
2. Validate stream-based ingestion parity with the stdout path on the replay corpus.
3. Switch gateway output to Redis Streams.
4. Remove stdout path.
5. When volume outgrows Redis, deploy NATS JetStream and migrate consumers.

---

## Alternatives considered

### A1 — Kafka from day one
**Rejected.** Operational overhead (ZooKeeper/KRaft, partitioning, rebalancing) is unjustified for Release 0.1–0.5 volumes. Redis Streams provides equivalent semantics with zero additional infrastructure.

### A2 — PostgreSQL COPY for bulk ingestion
**Rejected for events.** Batch inserts are already used; the bottleneck is insert throughput on the single primary. ClickHouse columnar storage is the correct long-term solution.

### A3 — gRPC streaming from gateway to ingestor
**Rejected for this phase.** stdout is simpler, debuggable, and works with container orchestration (Docker/K8s log collectors). gRPC is reserved for the Phase 4 gateway-control-plane mTLS contract.

---

## Consequences

### Positive
- Redis Streams provides at-least-once delivery with consumer groups, bounding memory.
- Migration path is incremental: same schema, same ingestor logic, swap transport.
- ClickHouse provides sub-second analytical queries on high-volume event data.

### Negative
- Redis adds complexity (memory sizing, `MAXLEN` tuning, replica for durability).
- Dual-write during migration increases gateway stdout output.
- Events dropped during ingestor buffer overflow are not recoverable (best-effort guarantee).

### Risks and mitigations
| Risk | Mitigation |
|---|---|
| Redis OOM from unbounded stream | `MAXLEN` (configurable, default 100k) |
| Stream consumer lag | Monitor `XINFO` lag; scale consumer group |
| ClickHouse migration cost | Same schema; write via `clickhouse-go` native insert |
| NATS/Kafka migration cost | Adapter pattern: same `Event` type, swap transport backend |

---

## Acceptance for this ADR

- Redis Streams path validated with event-ingestor parity test against the stdout path.
- `MAXLEN` eviction confirmed: oldest events are dropped without blocking the producer.
- Consumer group rebalancing: a new ingestor instance picks up from the last committed offset.
- ClickHouse schema auto-generated from the ADR-002 event schema (Phase B gate).

---

## Follow-up

- Implement Redis Streams adapter in `services/event-ingestor`.
- Add ClickHouse migration + ETL when volume exceeds PostgreSQL retention.