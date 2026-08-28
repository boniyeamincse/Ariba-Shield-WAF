# event-ingestor

Go sidecar that reads the gateway's event stream (stdout, JSON-lines) and forwards to
sinks (PostgreSQL, OpenTelemetry, Prometheus) asynchronously. Never backpressures the data plane.

**Sprint target:** Sprint 5 (initial ingestor + PostgreSQL persistence), Sprint 6 (release).