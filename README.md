# Raven

A distributed event streaming system built in Go.  
Producers publish messages to topics. Multiple broker nodes 
replicate and store them. Consumer groups pull messages at 
their own pace, tracking offsets independently.

Think simplified Kafka — built from scratch.

## Architecture

- **Broker nodes** — replicated, each holds a copy of partitions
- **Topics** — named streams, split into partitions
- **Producers** — push messages to a topic
- **Consumers** — pull messages, track offset per partition
- **Delivery guarantee** — at-least-once

## Tech Stack

Go · gRPC · Docker · Kubernetes · Prometheus · Grafana

## Project Status

| Milestone | Status |
|-----------|--------|
| M1 — Single node, HTTP API, in-memory store | 🔄 In Progress |
| M2 — Multi-broker replication | ⏳ Pending |
| M3 — Partitions + consumer offsets | ⏳ Pending |
| M4 — At-least-once delivery guarantees | ⏳ Pending |
| M5 — Docker + Kubernetes StatefulSet | ⏳ Pending |
| M6 — Prometheus + Grafana observability | ⏳ Pending |

## Running Locally

```bash
go run cmd/broker/main.go
```

## Author

Xenios Gerolemou — [LinkedIn](https://www.linkedin.com/in/xenios-gerolemou-594086202/) · [Portfolio](https://xenios7.github.io/portfolio/)