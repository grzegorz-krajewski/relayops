# Reliability

## Cel
Warstwa reliability w RelayOps ma uczyć podstaw zachowania systemu w warunkach błędów i przeciążeń.

## Obecnie zaimplementowane mechanizmy

### Readiness checks
Gateway sprawdza gotowość Redis i PostgreSQL.

### Retry dla błędów transient
Worker ponawia wywołania gRPC dla błędów chwilowych, takich jak:
- `UNAVAILABLE`
- `DEADLINE_EXCEEDED`
- `RESOURCE_EXHAUSTED`

### Rozróżnienie błędów
System rozróżnia:
- błędy permanent
- błędy transient

### Dead-letter stream
Taski, które kończą się ostatecznym niepowodzeniem, trafiają do:
- `tasks.dlq`

### Failed state persistence
Nieudane taski są zapisywane w PostgreSQL z informacją o błędzie.

## Obecna polityka
- permanent error -> `failed` + DLQ + ack
- transient error po wyczerpaniu retry -> `failed` + DLQ + ack
- success -> `processed` + ack

## Ograniczenia obecnej wersji
Na tym etapie system nie implementuje jeszcze:
- reclaim pending messages
- pełnego mechanizmu requeue
- osobnego dead-letter processora
- zaawansowanej polityki backoff
- idempotency guarantees dla ponownych dostarczeń

## Naturalne następne kroki
- stress testy
- poison message analysis
- LocalStack + SQS
- dead-letter reprocessing flow
- bardziej formalne runbooki incydentów
