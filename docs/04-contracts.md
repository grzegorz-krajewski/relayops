# Kontrakty gRPC i Protobuf

## Cel
Kontrakty w RelayOps definiują sposób komunikacji pomiędzy komponentami przetwarzającymi taski.

## Obecny kontrakt
Plik:
- `proto/task/v1/task.proto`

Definiuje:
- `ProcessTaskRequest`
- `ProcessTaskResponse`
- serwis `TaskProcessor`

## Główne pola requestu
- `task_id`
- `task_type`
- `raw_text`
- `trace_id`

## Główne pola response
- `task_id`
- `status`
- `normalized_text`
- `duration_ms`
- `error_message`

## Założenia projektowe
- kontrakt jest prosty i celowo ograniczony
- typy tasków są rozróżniane po `task_type`
- błędy transportowe i przetwarzania są mapowane na kody gRPC
- wersjonowanie odbywa się przez ścieżkę `v1`

## Kierunek rozwoju
W kolejnych etapach kontrakt może zostać rozszerzony o:
- bogatszy payload
- enumy dla typów tasków
- enumy dla statusów
- bardziej formalne mapowanie błędów
- dodatkowe serwisy gRPC
