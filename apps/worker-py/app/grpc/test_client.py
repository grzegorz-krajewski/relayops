import grpc

from task.v1 import task_pb2, task_pb2_grpc


def main():
    channel = grpc.insecure_channel("localhost:50051")
    stub = task_pb2_grpc.TaskProcessorStub(channel)

    response = stub.ProcessTask(
        task_pb2.ProcessTaskRequest(
            task_id="test-1",
            task_type="normalize_payload",
            raw_text="  Hello    gRPC   RelayOps  ",
            trace_id="trace-1",
        )
    )

    print(response)


if __name__ == "__main__":
    main()