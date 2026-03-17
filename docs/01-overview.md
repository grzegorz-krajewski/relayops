# Przegląd projektu

RelayOps to edukacyjna platforma do przetwarzania zdarzeń w architekturze rozproszonej.

## Główne cele
- nauka Go w kontekście backendowym
- nauka Redis Streams
- nauka workerów w Pythonie
- użycie gRPC i Protobuf jako kontraktu
- wdrożenie multi-environment config
- dodanie observability i podstaw reliability
- ćwiczenie retry i dead-letter strategy

## Główna idea
System przyjmuje taski przez HTTP API, publikuje je do Redis Streams, przetwarza asynchronicznie w Python workerze, zapisuje wyniki do PostgreSQL i wystawia metryki przez Prometheus oraz Grafanę.

## Główne komponenty
- Go Gateway
- Redis Streams
- Python Worker
- gRPC Task Processor
- PostgreSQL
- Prometheus
- Grafana
- Nginx

## Aktualne możliwości
RelayOps obsługuje już:
- przyjmowanie tasków
- zapis tasków do bazy
- odczyt statusu taska
- asynchroniczne przetwarzanie
- retry dla błędów transient
- dead-letter stream dla błędów końcowych
- dashboard obserwowalności
