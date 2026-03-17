import grpc

from task.v1 import task_pb2, task_pb2_grpc


class TaskProcessorClient:
    def __init__(self, target: str = "localhost:50051"):
        self.target = target

    def process_task(self, task_id: str, task_type: str, raw_text: str, trace_id: str) -> dict:
        with grpc.insecure_channel(self.target) as channel:
            stub = task_pb2_grpc.TaskProcessorStub(channel)

            response = stub.ProcessTask(
                task_pb2.ProcessTaskRequest(
                    task_id=task_id,
                    task_type=task_type,
                    raw_text=raw_text,
                    trace_id=trace_id,
                )
            )

            result = {
                "task_id": response.task_id,
                "status": response.status,
                "normalized_text": response.normalized_text,
                "duration_ms": response.duration_ms,
            }

            if response.error_message:
                result["error_message"] = response.error_message

            return result