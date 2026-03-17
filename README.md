# RelayOps

RelayOps is a distributed event processing and reliability learning platform built with Go and Python.

## Stack
- Go API gateway
- Python worker
- Redis Streams
- gRPC + Protobuf
- PostgreSQL
- Nginx
- Prometheus
- Grafana
- Docker Compose
- GitHub Actions

## MVP goal
Accept tasks through an HTTP API, publish them to Redis Streams, process them in a Python worker, persist results to PostgreSQL, and observe the system through metrics and dashboards.

## Local run
cp .env.example .env
make up

## Main services
- gateway-go: HTTP ingestion
- worker-py: async processing
- redis: stream backbone
- postgres: persistence
- nginx: reverse proxy
- prometheus: metrics scraping
- grafana: dashboards

## First milestone
- infra boots locally
- gateway health endpoint works
- worker boots
- Prometheus and Grafana are reachable