# Architecture

## Components
- Go Gateway: accepts HTTP requests and publishes events
- Redis Streams: event backbone
- Python Worker: consumes and processes tasks
- PostgreSQL: stores task results
- Nginx: reverse proxy
- Prometheus + Grafana: observability

## Flow
1. client sends task
2. gateway validates request
3. gateway publishes event
4. worker processes task
5. worker stores result
6. system exposes metrics