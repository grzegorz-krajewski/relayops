# RelayOps

RelayOps to platforma do przetwarzania zdarzeń i niezawodności systemów, zbudowana w oparciu o Go i Pythona.

Projekt powstał jako praktyczne środowisko do przetestowania modułów:

- backendu w Go
- Redis Streams
- workerów w Pythonie
- kontraktów gRPC i Protobuf
- persystencji w PostgreSQL
- obserwowalności przez Prometheus i Grafanę
- strategii retry i dead-letter flow
- podstaw projektowania systemów rozproszonych

## Obecna architektura

    Klient
      -> Nginx
      -> Go Gateway API
      -> Redis Streams
      -> Python Worker
      -> gRPC Task Processor
      -> PostgreSQL

    Prometheus <- metryki Gatewaya i Workera
    Grafana <- Prometheus

## Co jest już zaimplementowane

### Gateway (Go)
- HTTP API do przyjmowania tasków
- zapis tasków do PostgreSQL
- endpoint do odczytu statusu taska
- publikacja do Redis Streams
- endpointy health i readiness
- metryki Prometheusa

### Worker (Python)
- consumer group dla Redis Streams
- przetwarzanie oparte o gRPC
- aktualizacja statusu taska w PostgreSQL
- retry dla błędów transient
- publikacja do dead-letter stream
- metryki Prometheusa

### Platforma
- lokalne środowisko w Docker Compose
- reverse proxy w Nginx
- PostgreSQL
- Redis
- Prometheus
- Grafana
- provisionowany dashboard Grafany

## Główne przepływy

### Ścieżka poprawna
1. Klient wysyła task do API w Go
2. Gateway waliduje dane i zapisuje task do PostgreSQL
3. Gateway publikuje task do Redis Streams
4. Python worker odbiera wiadomość
5. Worker wywołuje gRPC Task Processor
6. Worker zapisuje wynik do PostgreSQL
7. Worker potwierdza wiadomość w Redis Streams przez ack

### Ścieżka błędu
1. Przetwarzanie kończy się błędem trwałym albo chwilowym
2. Błędy chwilowe są ponawiane
3. Jeżeli retry się wyczerpią albo błąd jest trwały:
   - task dostaje status `failed` w PostgreSQL
   - task trafia do dead-letter stream
   - oryginalna wiadomość jest potwierdzana

## Stos technologiczny

### Backend
- Go
- Python 3.12
- gRPC
- Protobuf

### Infrastruktura
- Docker Compose
- Redis Streams
- PostgreSQL
- Nginx

### Observability
- Prometheus
- Grafana

## Struktura repozytorium

    relayops/
    ├─ apps/
    │  ├─ gateway-go/
    │  └─ worker-py/
    ├─ proto/
    │  └─ task/v1/
    ├─ deploy/
    │  ├─ compose/
    │  ├─ grafana/
    │  ├─ nginx/
    │  ├─ postgres/
    │  └─ prometheus/
    ├─ docs/
    ├─ .github/workflows/
    ├─ Makefile
    └─ README.md

## Uruchomienie lokalne

### Wymagania
- Docker Desktop
- Go lokalnie tylko wtedy, gdy chcesz robić buildy lub codegen poza Dockerem
- Python lokalnie tylko wtedy, gdy chcesz generować protobufy poza Dockerem

### Start
    cp .env.example .env
    make up

### Zatrzymanie
    make down

### Logi
    make logs

## Główne endpointy

### Gateway
- `GET /health`
- `GET /ready`
- `POST /api/v1/tasks`
- `GET /api/v1/tasks/{id}`

### Metryki
- metryki gatewaya są dostępne na porcie hosta z `.env`
- metryki workera są wystawione w sieci Dockera i scrapowane przez Prometheusa

## Przykładowy request

    curl -X POST http://localhost:8082/api/v1/tasks       -H "Content-Type: application/json"       -d '{
        "type": "normalize_payload",
        "payload": {
          "text": "  Hello RelayOps  "
        }
      }'

## Przykładowe typy tasków

- `normalize_payload`
- `enrich_text`
- `force_permanent_error`
- `force_transient_error`

## Założenia niezawodności, które już działają

- readiness checks dla Redis i PostgreSQL
- retry dla błędów transient w workerze
- rozróżnienie błędów permanent i transient
- dead-letter stream `tasks.dlq`
- zapis failed tasków w PostgreSQL
- metryki dla requestów, przetwarzania, retry, failures i dead-letter events

## Observability

### Prometheus
Projekt wystawia metryki dotyczące:
- requestów HTTP gatewaya
- liczby utworzonych tasków
- błędów publish/persist
- liczby przetworzonych tasków w workerze
- liczby błędów w workerze
- liczby retry
- liczby dead-letter events
- czasu przetwarzania

### Grafana
Provisionowany dashboard:
- `RelayOps Overview`

## Dead-letter stream

Nieudane taski trafiają do:
- `tasks.dlq`

Wiadomość w DLQ zawiera:
- `task_id`
- `task_type`
- `trace_id`
- `raw_payload`
- `failure_kind`
- `error_message`

## Planowane kolejne kroki

- testy obciążeniowe w k6
- integracja z LocalStack i SQS
- bogatsze kontrakty gRPC
- wiele instancji gatewaya za Nginx
- bardziej rozbudowane dashboardy i alerty
- runbooki i dokumentacja incydentów

## Uwagi projektowe

To jest projekt edukacyjny i rozwijany iteracyjnie.
Część decyzji została uproszczona celowo, żeby łatwiej było zrozumieć przepływ zdarzeń, kontrakty, retry i reliability patterns, a później rozwijać to dalej.
