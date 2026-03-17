import json
import time
from redis import Redis
from redis.exceptions import ResponseError

from app.db.task_repository import TaskRepository
from app.grpc.client import (
    PermanentTaskError,
    TaskProcessorClient,
    TransientTaskError,
)
from app.metrics.metrics import (
    grpc_calls_total,
    task_processing_duration_seconds,
    tasks_acked_total,
    tasks_failed_total,
    tasks_processed_total,
    transient_retries_total,
)


class StreamConsumer:
    def __init__(
        self,
        redis_addr: str,
        stream_name: str,
        group_name: str,
        consumer_name: str,
        task_repository: TaskRepository,
        grpc_client: TaskProcessorClient,
        max_transient_retries: int,
        retry_backoff_seconds: float,
    ):
        host, port = redis_addr.split(":")
        self.redis = Redis(host=host, port=int(port), decode_responses=True)
        self.stream_name = stream_name
        self.group_name = group_name
        self.consumer_name = consumer_name
        self.task_repository = task_repository
        self.grpc_client = grpc_client
        self.max_transient_retries = max_transient_retries
        self.retry_backoff_seconds = retry_backoff_seconds

    def ensure_group(self) -> None:
        try:
            self.redis.xgroup_create(
                name=self.stream_name,
                groupname=self.group_name,
                id="0",
                mkstream=True,
            )
            print(f"created consumer group={self.group_name} stream={self.stream_name}")
        except ResponseError as exc:
            if "BUSYGROUP" in str(exc):
                print(f"consumer group already exists group={self.group_name}")
            else:
                raise

    def run(self) -> None:
        self.ensure_group()
        print(
            f"worker consuming stream={self.stream_name} "
            f"group={self.group_name} consumer={self.consumer_name}"
        )

        while True:
            response = self.redis.xreadgroup(
                groupname=self.group_name,
                consumername=self.consumer_name,
                streams={self.stream_name: ">"},
                count=10,
                block=5000,
            )

            if not response:
                continue

            for stream_name, messages in response:
                for message_id, fields in messages:
                    self.handle_message(stream_name, message_id, fields)

    def handle_message(self, stream_name: str, message_id: str, fields: dict) -> None:
        task_id = fields.get("task_id", "")
        task_type = fields.get("task_type", "")
        raw_payload = fields.get("raw_payload", "{}")
        trace_id = fields.get("trace_id", "")

        print(
            f"received message_id={message_id} task_id={task_id} "
            f"task_type={task_type} trace_id={trace_id}"
        )

        started = time.perf_counter()

        try:
            payload = json.loads(raw_payload)
            raw_text = str(payload.get("text", ""))

            result = self._process_with_retry(
                task_id=task_id,
                task_type=task_type,
                raw_text=raw_text,
                trace_id=trace_id,
            )

            duration_seconds = time.perf_counter() - started
            duration_ms = int(duration_seconds * 1000)

            result["worker_duration_ms"] = duration_ms
            result["processor"] = "grpc"

            self.task_repository.mark_processed(task_id=task_id, result_payload=result)

            tasks_processed_total.inc()
            task_processing_duration_seconds.observe(duration_seconds)

            print(
                f"processed via grpc task_id={task_id} message_id={message_id} "
                f"duration_ms={duration_ms} result={result}"
            )

            self.redis.xack(stream_name, self.group_name, message_id)
            tasks_acked_total.inc()
            print(f"acked message_id={message_id}")

        except PermanentTaskError as exc:
            tasks_failed_total.inc()
            self.task_repository.mark_failed(task_id=task_id, error_message=str(exc))
            self.redis.xack(stream_name, self.group_name, message_id)
            tasks_acked_total.inc()
            print(
                f"permanent failure task_id={task_id} message_id={message_id} "
                f"task_type={task_type} error={exc}"
            )

        except TransientTaskError as exc:
            tasks_failed_total.inc()
            self.task_repository.mark_failed(task_id=task_id, error_message=str(exc))
            self.redis.xack(stream_name, self.group_name, message_id)
            tasks_acked_total.inc()
            print(
                f"transient failure after retries task_id={task_id} "
                f"message_id={message_id} task_type={task_type} error={exc}"
            )

        except Exception as exc:
            tasks_failed_total.inc()
            self.task_repository.mark_failed(task_id=task_id, error_message=str(exc))
            self.redis.xack(stream_name, self.group_name, message_id)
            tasks_acked_total.inc()
            print(
                f"unexpected failure task_id={task_id} message_id={message_id} "
                f"task_type={task_type} error={exc}"
            )

    def _process_with_retry(self, task_id: str, task_type: str, raw_text: str, trace_id: str) -> dict:
        attempt = 0

        while True:
            attempt += 1
            grpc_calls_total.inc()

            try:
                return self.grpc_client.process_task(
                    task_id=task_id,
                    task_type=task_type,
                    raw_text=raw_text,
                    trace_id=trace_id,
                )

            except PermanentTaskError:
                raise

            except TransientTaskError as exc:
                if attempt >= self.max_transient_retries:
                    raise

                transient_retries_total.inc()
                print(
                    f"transient retry task_id={task_id} task_type={task_type} "
                    f"attempt={attempt}/{self.max_transient_retries} error={exc}"
                )
                time.sleep(self.retry_backoff_seconds)