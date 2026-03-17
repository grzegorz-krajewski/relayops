# Architektura

## Komponenty

### Go Gateway
Odpowiada za:
- przyjmowanie żądań HTTP
- walidację danych wejściowych
- zapis taska do PostgreSQL
- publikację taska do Redis Streams
- udostępnianie endpointów health i ready
- wystawianie metryk Prometheusa

### Redis Streams
Odpowiada za:
- transport tasków między gatewayem a workerem
- buforowanie wiadomości
- obsługę consumer group

### Python Worker
Odpowiada za:
- odbieranie tasków z Redis Streams
- wykonywanie retry dla błędów transient
- wywołanie gRPC processora
- aktualizację stanu taska w PostgreSQL
- wysyłanie wiadomości do dead-letter stream
- wystawianie metryk Prometheusa

### gRPC Task Processor
Odpowiada za:
- kontrakt przetwarzania tasków
- typowaną komunikację opartą o Protobuf
- rozróżnienie ścieżek poprawnych i błędów kontrolowanych

### PostgreSQL
Odpowiada za:
- zapis stanu tasków
- zapis wyników przetwarzania
- zapis błędów końcowych

### Prometheus i Grafana
Odpowiadają za:
- zbieranie metryk
- wizualizację działania systemu
- obserwację requestów, retry, błędów i DLQ

## Przepływ danych

### Ścieżka poprawna

    Klient
      -> Gateway HTTP
      -> PostgreSQL: insert task (accepted)
      -> Redis Streams: publish event
      -> Worker: consume message
      -> gRPC Processor: ProcessTask
      -> PostgreSQL: update task (processed)
      -> Redis Streams: ack

### Ścieżka błędu

    Klient
      -> Gateway HTTP
      -> PostgreSQL: insert task
      -> Redis Streams: publish event
      -> Worker: consume message
      -> gRPC Processor: error
      -> retry dla transient
      -> PostgreSQL: update task (failed)
      -> Redis DLQ: publish failed event
      -> Redis Streams: ack

## Decyzje architektoniczne

### Redis Streams zamiast Kafka
Na obecnym etapie Redis Streams jest prostsze do wdrożenia, szybsze do zrozumienia i wystarczające do nauki stream processingu, consumer groups i retry patterns.

### gRPC + Protobuf zamiast JSON między usługami
Kontrakt binarny daje większą spójność typów, lepsze podstawy pod wersjonowanie i bardziej realistyczny model komunikacji między komponentami.

### Worker i gRPC processor w jednym serwisie
Na obecnym etapie to uproszczenie celowe. Dzięki temu można ćwiczyć kontrakty, retry i przetwarzanie bez rozbijania projektu na zbyt wiele usług naraz.

### Dead-letter stream w Redis
To prosty i skuteczny sposób na analizę wiadomości, które ostatecznie nie przeszły przetwarzania.
