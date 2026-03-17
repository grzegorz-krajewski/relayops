package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type TaskStore struct {
	db *sql.DB
}

type TaskRecord struct {
	ID              string
	Type            string
	Status          string
	RawPayload      []byte
	ResultPayload   []byte
	TraceID         string
	StreamMessageID sql.NullString
	ErrorMessage    sql.NullString
	CreatedAt       time.Time
	ProcessedAt     sql.NullTime
}

func NewTaskStore(dsn string) (*TaskStore, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	return &TaskStore{db: db}, nil
}

func (s *TaskStore) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return s.db.PingContext(ctx)
}

func (s *TaskStore) InsertTask(
	ctx context.Context,
	taskID string,
	taskType string,
	status string,
	rawPayload map[string]any,
	traceID string,
	streamMessageID string,
) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rawPayloadBytes, err := json.Marshal(rawPayload)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(
		ctx,
		`
		INSERT INTO tasks (id, type, status, raw_payload, trace_id, stream_message_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		`,
		taskID,
		taskType,
		status,
		rawPayloadBytes,
		traceID,
		streamMessageID,
	)

	return err
}

func (s *TaskStore) GetTaskByID(ctx context.Context, taskID string) (*TaskRecord, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(
		ctx,
		`
		SELECT id, type, status, raw_payload, result_payload, trace_id,
		       stream_message_id, error_message, created_at, processed_at
		FROM tasks
		WHERE id = $1
		`,
		taskID,
	)

	var rec TaskRecord
	err := row.Scan(
		&rec.ID,
		&rec.Type,
		&rec.Status,
		&rec.RawPayload,
		&rec.ResultPayload,
		&rec.TraceID,
		&rec.StreamMessageID,
		&rec.ErrorMessage,
		&rec.CreatedAt,
		&rec.ProcessedAt,
	)
	if err != nil {
		return nil, err
	}

	return &rec, nil
}
