import threading
import time
from http.server import BaseHTTPRequestHandler, HTTPServer

from prometheus_client import CONTENT_TYPE_LATEST, generate_latest

from app.config.settings import Settings
from app.consumers.stream_consumer import StreamConsumer
from app.db.task_repository import TaskRepository
from app.grpc.client import TaskProcessorClient
from app.grpc.server import serve as serve_grpc


class MetricsHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == "/metrics":
            data = generate_latest()
            self.send_response(200)
            self.send_header("Content-Type", CONTENT_TYPE_LATEST)
            self.send_header("Content-Length", str(len(data)))
            self.end_headers()
            self.wfile.write(data)
        else:
            self.send_response(404)
            self.end_headers()


def serve_metrics():
    server = HTTPServer(("0.0.0.0", 9100), MetricsHandler)
    server.serve_forever()


def main():
    settings = Settings()

    metrics_thread = threading.Thread(target=serve_metrics, daemon=True)
    metrics_thread.start()

    grpc_thread = threading.Thread(target=serve_grpc, daemon=True)
    grpc_thread.start()

    print(
        f"worker started env={settings.app_env} "
        f"name={settings.worker_name} stream={settings.redis_stream_name}"
    )

    repository = TaskRepository(settings.postgres_dsn)
    grpc_client = TaskProcessorClient(target="localhost:50051")

    consumer = StreamConsumer(
        redis_addr=settings.redis_addr,
        stream_name=settings.redis_stream_name,
        group_name=settings.worker_group,
        consumer_name=settings.worker_name,
        task_repository=repository,
        grpc_client=grpc_client,
    )

    while True:
        try:
            consumer.run()
        except Exception as exc:
            print(f"worker loop error={exc}")
            time.sleep(3)


if __name__ == "__main__":
    main()