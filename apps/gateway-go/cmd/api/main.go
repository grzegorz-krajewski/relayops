package main

import (
	"context"
	"fmt"
	stdhttp "net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"relayops/apps/gateway-go/internal/config"
	apphttp "relayops/apps/gateway-go/internal/http"
	"relayops/apps/gateway-go/internal/metrics"
	"relayops/apps/gateway-go/internal/redisstream"
	"relayops/apps/gateway-go/internal/store"
)

func main() {
	cfg := config.Load()

	metrics.MustRegister()

	publisher := redisstream.NewPublisher(cfg.RedisAddr, cfg.RedisStreamName)
	if err := waitForRedis(publisher, 10, 2*time.Second); err != nil {
		panic(fmt.Sprintf("redis not ready: %v", err))
	}

	taskStore, err := store.NewTaskStore(cfg.PostgresDSN)
	if err != nil {
		panic(fmt.Sprintf("postgres init failed: %v", err))
	}
	if err := waitForPostgres(taskStore, 15, 2*time.Second); err != nil {
		panic(fmt.Sprintf("postgres not ready: %v", err))
	}

	go func() {
		metricsMux := stdhttp.NewServeMux()
		metricsMux.Handle("/metrics", promhttp.Handler())

		fmt.Printf("metrics listening on :%s\n", cfg.MetricsPort)
		if err := stdhttp.ListenAndServe(":"+cfg.MetricsPort, metricsMux); err != nil {
			panic(err)
		}
	}()

	mux := stdhttp.NewServeMux()

	handler := apphttp.NewHandler(publisher, taskStore, cfg.RedisStreamName)
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
			_, _ = w.Write([]byte(`{"status":"not_ready","dependency":"redis"}`))
			return
		}
		if err := taskStore.Ping(ctx); err != nil {
			w.WriteHeader(stdhttp.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"status":"not_ready","dependency":"postgres"}`))
			return
		}

		w.WriteHeader(stdhttp.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ready"}`))
	})

	mux.HandleFunc("/", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		w.WriteHeader(stdhttp.StatusOK)
		_, _ = w.Write([]byte("RelayOps gateway is running"))
	})

	handlerWithMetrics := apphttp.MetricsMiddleware(mux)

	fmt.Printf("gateway listening on :%s\n", cfg.HTTPPort)
	if err := stdhttp.ListenAndServe(":"+cfg.HTTPPort, handlerWithMetrics); err != nil {
		panic(err)
	}
}

func waitForRedis(publisher *redisstream.Publisher, attempts int, delay time.Duration) error {
	var lastErr error

	for i := 1; i <= attempts; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		err := publisher.Ping(ctx)
		cancel()

		if err == nil {
			fmt.Printf("redis ready after %d attempt(s)\n", i)
			return nil
		}

		lastErr = err
		fmt.Printf("waiting for redis (attempt %d/%d): %v\n", i, attempts, err)
		time.Sleep(delay)
	}

	return lastErr
}

func waitForPostgres(taskStore *store.TaskStore, attempts int, delay time.Duration) error {
	var lastErr error

	for i := 1; i <= attempts; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		err := taskStore.Ping(ctx)
		cancel()

		if err == nil {
			fmt.Printf("postgres ready after %d attempt(s)\n", i)
			return nil
		}

		lastErr = err
		fmt.Printf("waiting for postgres (attempt %d/%d): %v\n", i, attempts, err)
		time.Sleep(delay)
	}

	return lastErr
}
