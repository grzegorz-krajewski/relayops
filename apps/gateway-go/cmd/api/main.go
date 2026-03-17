package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	httpPort := getenv("GATEWAY_HTTP_PORT", "8080")
	metricsPort := getenv("GATEWAY_METRICS_PORT", "9090")

	go func() {
		metricsMux := http.NewServeMux()
		metricsMux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("# metrics placeholder\n"))
		})

		fmt.Printf("metrics listening on :%s\n", metricsPort)
		if err := http.ListenAndServe(":"+metricsPort, metricsMux); err != nil {
			panic(err)
		}
	}()

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	r.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ready"}`))
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("RelayOps gateway is running"))
	})

	fmt.Printf("gateway listening on :%s\n", httpPort)
	if err := http.ListenAndServe(":"+httpPort, r); err != nil {
		panic(err)
	}
}

func getenv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
