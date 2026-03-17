package main

import (
	"context"
	"fmt"
	stdhttp "net/http"
	"time"

	"relayops/apps/gateway-go/internal/config"
	apphttp "relayops/apps/gateway-go/internal/http"
	"relayops/apps/gateway-go/internal/redisstream"
)

func main() {
	cfg := config.Load()

	publisher := redisstream.NewPublisher(cfg.RedisAddr, cfg.RedisStreamName)

	if err := publisher.Ping(context.Background()); err != nil {
		panic(fmt.Sprintf("redis ping failed: %v", err))
	}

	go func() {
		metricsMux := stdhttp.NewServeMux()
		metricsMux.HandleFunc("/metrics", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			w.WriteHeader(stdhttp.StatusOK)
			_, _ = w.Write([]byte("# metrics placeholder\n"))
		})

		fmt.Printf("metrics listening on :%s\n", cfg.MetricsPort)
		if err := stdhttp.ListenAndServe(":"+cfg.MetricsPort, metricsMux); err != nil {
			panic(err)
		}
	}()

	mux := stdhttp.NewServeMux()

	handler := apphttp.NewHandler(publisher, cfg.RedisStreamName)
	handler.RegisterRoutes(mux)

	mux.HandleFunc("/health", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		w.WriteHeader(stdhttp.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("/ready", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()

		if err := publisher.Ping(ctx); err != nil {
			w.WriteHeader(stdhttp.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"status":"not_ready"}`))
			return
		}

		w.WriteHeader(stdhttp.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ready"}`))
	})

	mux.HandleFunc("/", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		w.WriteHeader(stdhttp.StatusOK)
		_, _ = w.Write([]byte("RelayOps gateway is running"))
	})

	fmt.Printf("gateway listening on :%s\n", cfg.HTTPPort)
	if err := stdhttp.ListenAndServe(":"+cfg.HTTPPort, mux); err != nil {
		panic(err)
	}
}
