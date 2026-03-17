package redisstream

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Publisher struct {
	client *redis.Client
	stream string
}

func NewPublisher(addr, stream string) *Publisher {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &Publisher{
		client: client,
		stream: stream,
	}
}

func (p *Publisher) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return p.client.Ping(ctx).Err()
}

func (p *Publisher) PublishTask(
	ctx context.Context,
	taskID string,
	taskType string,
	rawPayload string,
	traceID string,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return p.client.XAdd(ctx, &redis.XAddArgs{
		Stream: p.stream,
		Values: map[string]any{
			"task_id":     taskID,
			"task_type":   taskType,
			"raw_payload": rawPayload,
			"trace_id":    traceID,
			"created_at":  time.Now().UTC().Format(time.RFC3339Nano),
		},
	}).Result()
}
