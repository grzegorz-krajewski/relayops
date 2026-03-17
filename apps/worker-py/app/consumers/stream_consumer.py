import time
from redis import Redis
from redis.exceptions import ResponseError

from app.processors.task_processor import process_task


class StreamConsumer:
    def __init__(self, redis_addr: str, stream_name: str, group_name: str, consumer_name: str):
        host, port = redis_addr.split(":")
        self.redis = Redis(host=host, port=int(port), decode_responses=True)
        self.stream_name = stream_name
        self.group_name = group_name
        self.consumer_name = consumer_name

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
            result = process_task(task_type=task_type, raw_payload=raw_payload)
            duration_ms = int((time.perf_counter() - started) * 1000)

            print(
                f"processed task_id={task_id} message_id={message_id} "
                f"duration_ms={duration_ms} result={result}"
            )

            self.redis.xack(stream_name, self.group_name, message_id)
            print(f"acked message_id={message_id}")

        except Exception as exc:
            print(
                f"processing failed task_id={task_id} message_id={message_id} "
                f"error={exc}"
            )