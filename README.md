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
| M1 — Single node, HTTP API, in-memory store | ✅ Complete |
| M2 — Multi-broker replication | ✅ Complete |
| M3 — Partitions + consumer offsets | ✅ Complete |
| M4 — At-least-once delivery guarantees | ✅ Complete |
| M5 — Docker + Kubernetes StatefulSet | ✅ Complete |
| M6 — Prometheus + Grafana observability | ✅ Complete |

## Grafana Dashboard

[![Grafana Dashboard](./assets/grafana-dashboard.png)](./assets/grafana-dashboard.png)


## Running Locally

```bash
go run cmd/broker/main.go
```

## Docker Compose

```bash
docker-compose up --build
```

## Kubernetes

```bash
kind create cluster --name raven
docker build -t raven:latest -f docker/Dockerfile .
kind load docker-image raven:latest --name raven
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/statefulset.yaml
kubectl port-forward raven-0 8080:8080
```

## Author

Xenios Gerolemou — [LinkedIn](https://www.linkedin.com/in/xenios-gerolemou-594086202/) · [Portfolio](https://xenios7.github.io/portfolio/)
