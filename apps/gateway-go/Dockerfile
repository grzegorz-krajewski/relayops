FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o gateway ./cmd/api

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/gateway /app/gateway

EXPOSE 8080
EXPOSE 9090

CMD ["/app/gateway"]