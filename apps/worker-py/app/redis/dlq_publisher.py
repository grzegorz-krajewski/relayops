import json
from redis import Redis


class DLQPublisher:
    def __init__(self, redis_addr: str, dlq_stream_name: str):
        host, port = redis_addr.split(":")
        self.redis = Redis(host=host, port=int(port), decode_responses=True)
        self.dlq_stream_name = dlq_stream_name

    def publish(
        self,
        task_id: str,
        task_type: str,
        trace_id: str,
        raw_payload: str,
        failure_kind: str,
        error_message: str,
    ) -> str:
        return self.redis.xadd(
            self.dlq_stream_name,
            {
                "task_id": task_id,
                "task_type": task_type,
                "trace_id": trace_id,
                "raw_payload": raw_payload,
                "failure_kind": failure_kind,
                "error_message": error_message,
            },
        )