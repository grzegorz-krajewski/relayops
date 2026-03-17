# Observability

## Cel
Observability w RelayOps ma umożliwiać:
- zrozumienie przepływu requestów i tasków
- szybkie wykrywanie błędów
- analizę retry i DLQ
- przygotowanie systemu pod testy obciążeniowe

## Metryki Gatewaya
Gateway wystawia metryki dotyczące:
- liczby requestów HTTP
- liczby zaakceptowanych tasków
- błędów publikacji do streamu
- błędów zapisu do bazy
- czasu trwania requestów

## Metryki Workera
Worker wystawia metryki dotyczące:
- liczby przetworzonych tasków
- liczby failed tasków
- liczby acków
- liczby wywołań gRPC
- liczby retry dla błędów transient
- liczby tasków wysłanych do DLQ
- czasu przetwarzania

## Prometheus
Prometheus scrapuje:
- gateway
- worker

## Grafana
Provisionowany dashboard:
- `RelayOps Overview`

Dashboard pokazuje m.in.:
- tempo requestów gatewaya
- tempo tworzenia tasków
- liczbę tasków przetworzonych przez workera
- liczbę failed tasków
- czas przetwarzania
- tempo wywołań gRPC

## Dalszy rozwój observability
Kolejne kroki mogą obejmować:
- alerty
- metryki per task type
- lepsze dashboardy p95/p99
- structured logging z correlation ID
- tracing
