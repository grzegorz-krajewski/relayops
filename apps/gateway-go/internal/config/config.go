package config

import "os"

type Config struct {
	AppEnv             string
	HTTPPort           string
	MetricsPort        string
	RedisAddr          string
	RedisStreamName    string
	RedisDLQStreamName string
	PostgresDSN        string
}

func Load() Config {
	return Config{
		AppEnv:             getenv("APP_ENV", "local"),
		HTTPPort:           getenv("GATEWAY_HTTP_PORT", "8080"),
		MetricsPort:        getenv("GATEWAY_METRICS_PORT", "9090"),
		RedisAddr:          getenv("REDIS_ADDR", "redis:6379"),
		RedisStreamName:    getenv("REDIS_STREAM_NAME", "tasks.stream"),
		RedisDLQStreamName: getenv("REDIS_DLQ_STREAM_NAME", "tasks.dlq"),
		PostgresDSN:        getenv("POSTGRES_DSN", "postgresql://relayops:relayops@postgres:5432/relayops"),
	}
}

func getenv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
