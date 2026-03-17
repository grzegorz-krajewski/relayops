CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY,
    type TEXT NOT NULL,
    status TEXT NOT NULL,
    raw_payload JSONB NOT NULL,
    result_payload JSONB,
    trace_id UUID NOT NULL,
    stream_message_id TEXT,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ
);