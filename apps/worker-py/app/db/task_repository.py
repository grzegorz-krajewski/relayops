import json
from contextlib import contextmanager

import psycopg


class TaskRepository:
    def __init__(self, dsn: str):
        self.dsn = dsn

    @contextmanager
    def connection(self):
        with psycopg.connect(self.dsn) as conn:
            yield conn

    def mark_processed(self, task_id: str, result_payload: dict) -> None:
        with self.connection() as conn:
            with conn.cursor() as cur:
                cur.execute(
                    """
                    UPDATE tasks
                    SET status = %s,
                        result_payload = %s::jsonb,
                        processed_at = NOW(),
                        error_message = NULL
                    WHERE id = %s
                    """,
                    ("processed", json.dumps(result_payload), task_id),
                )
            conn.commit()

    def mark_failed(self, task_id: str, error_message: str) -> None:
        with self.connection() as conn:
            with conn.cursor() as cur:
                cur.execute(
                    """
                    UPDATE tasks
                    SET status = %s,
                        error_message = %s,
                        processed_at = NOW()
                    WHERE id = %s
                    """,
                    ("failed", error_message, task_id),
                )
            conn.commit()