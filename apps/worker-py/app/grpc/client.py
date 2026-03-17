import grpc

from task.v1 import task_pb2, task_pb2_grpc


class PermanentTaskError(Exception):
    pass


class TransientTaskError(Exception):
    pass


class TaskProcessorClient:
    def __init__(self, target: str = "localhost:50051"):
        self.target = target

    def process_task(self, task_id: str, task_type: str, raw_text: str, trace_id: str) -> dict:
        with grpc.insecure_channel(self.target) as channel:
            stub = task_pb2_grpc.TaskProcessorStub(channel)

            try:
                response = stub.ProcessTask(
                    task_pb2.ProcessTaskRequest(
                        task_id=task_id,
                        task_type=task_type,
                        raw_text=raw_text,
                        trace_id=trace_id,
                    )
                )
            except grpc.RpcError as exc:
                code = exc.code().name if exc.code() else "UNKNOWN"
                details = exc.details() or "grpc processing error"
                message = f"grpc_error code={code} details={details}"

                if code in {"UNAVAILABLE", "DEADLINE_EXCEEDED", "RESOURCE_EXHAUSTED"}:
                    raise TransientTaskError(message) from exc

                raise PermanentTaskError(message) from exc

            result = {
                "task_id": response.task_id,
                "status": response.status,
                "normalized_text": response.normalized_text,
                "duration_ms": response.duration_ms,
            }

            if response.error_message:
                result["error_message"] = response.error_message

            return result