import threading
import time
from http.server import BaseHTTPRequestHandler, HTTPServer

from app.config.settings import Settings
from app.consumers.stream_consumer import StreamConsumer
from app.db.task_repository import TaskRepository


class MetricsHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == "/metrics":
            self.send_response(200)
            self.send_header("Content-Type", "text/plain; version=0.0.4")
            self.end_headers()
            self.wfile.write(b"# worker metrics placeholder\n")
        else:
            self.send_response(404)
            self.end_headers()


def serve_metrics():
    server = HTTPServer(("0.0.0.0", 9100), MetricsHandler)
    server.serve_forever()


def main():
    settings = Settings()

    thread = threading.Thread(target=serve_metrics, daemon=True)
    thread.start()

    print(
        f"worker started env={settings.app_env} "
        f"name={settings.worker_name} stream={settings.redis_stream_name}"
    )

    repository = TaskRepository(settings.postgres_dsn)

    consumer = StreamConsumer(
        redis_addr=settings.redis_addr,
        stream_name=settings.redis_stream_name,
        group_name=settings.worker_group,
        consumer_name=settings.worker_name,
        task_repository=repository,
    )

    while True:
        try:
            consumer.run()
        except Exception as exc:
            print(f"worker loop error={exc}")
            time.sleep(3)


if __name__ == "__main__":
    main()