# Środowiska i konfiguracja

## Planowane środowiska
- local
- dev
- test
- staging
- prod

Na obecnym etapie aktywnie używane jest głównie środowisko lokalne.

## Zasady konfiguracji
- konfiguracja pochodzi ze zmiennych środowiskowych
- sekrety nie trafiają do repozytorium
- `.env.example` zawiera bezpieczne wartości przykładowe
- `.env` jest lokalny i ignorowany przez Git

## Najważniejsze zmienne

### Gateway
- `GATEWAY_HTTP_PORT`
- `GATEWAY_METRICS_PORT`
- `GATEWAY_HOST_PORT`
- `GATEWAY_METRICS_HOST_PORT`
- `REDIS_ADDR`
- `REDIS_STREAM_NAME`
- `REDIS_DLQ_STREAM_NAME`
- `POSTGRES_DSN`

### Worker
- `WORKER_NAME`
- `WORKER_GROUP`
- `GRPC_TARGET`
- `MAX_TRANSIENT_RETRIES`
- `RETRY_BACKOFF_SECONDS`

### Platforma
- `NGINX_PORT`
- `PROMETHEUS_PORT`
- `GRAFANA_PORT`

## Dobre praktyki
- rozdzielaj porty hosta od portów wewnętrznych kontenerów
- nie zakładaj, że `depends_on` oznacza gotowość usługi
- dla usług zależnych stosuj health checks albo retry przy starcie
- trzymaj konfigurację spójną między Go i Pythonem
