# Embodied Robotics Mission Control

This repository is a self-built production-style Go backend for multi-tenant
embodied robot fleet operations. It coordinates robot registration, capability
versions, mission approval, scheduling, telemetry ingestion, policy review,
remote handoff, alerts, audit trails and durable outbox delivery.

## Run

```text
GOTOOLCHAIN=local go run ./cmd/server
```

The service uses SQLite with migrations in `migrations/001_init.sql`. Set
`DATABASE_URL` and `PORT` through the environment. `GET /healthz` is the
liveness endpoint; mission and robot APIs require a Bearer session token.

## Verification

```text
GOTOOLCHAIN=local go build ./...
GOTOOLCHAIN=local go test ./... -count=1
GOTOOLCHAIN=local go test -race ./... -count=1
GOTOOLCHAIN=local go vet ./...
```

The code is intentionally layered: HTTP handlers depend on services, services
own transactions and invariants, repositories own SQL, and workers propagate
context cancellation while recording durable retries and audit events.
