import grpc
from concurrent import futures

from task.v1 import task_pb2, task_pb2_grpc


class TaskProcessorService(task_pb2_grpc.TaskProcessorServicer):
    def ProcessTask(self, request, context):
        raw_text = request.raw_text or ""
        normalized = " ".join(raw_text.split())

        if request.task_type == "normalize_payload":
            return task_pb2.ProcessTaskResponse(
                task_id=request.task_id,
                status="processed",
                normalized_text=normalized,
                duration_ms=1,
                error_message="",
            )

        if request.task_type == "enrich_text":
            return task_pb2.ProcessTaskResponse(
                task_id=request.task_id,
                status="processed",
                normalized_text=f"{normalized} :: enriched",
                duration_ms=1,
                error_message="",
            )

        return task_pb2.ProcessTaskResponse(
            task_id=request.task_id,
            status="processed",
            normalized_text=normalized,
            duration_ms=1,
            error_message="",
        )


def serve() -> None:
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    task_pb2_grpc.add_TaskProcessorServicer_to_server(TaskProcessorService(), server)
    server.add_insecure_port("[::]:50051")
    server.start()

    print("grpc server listening on :50051")
    server.wait_for_termination()


if __name__ == "__main__":
    serve()