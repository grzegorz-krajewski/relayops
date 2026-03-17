import os
import time
import threading
from http.server import BaseHTTPRequestHandler, HTTPServer


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
    worker_name = os.getenv("WORKER_NAME", "worker-1")
    env = os.getenv("APP_ENV", "local")

    thread = threading.Thread(target=serve_metrics, daemon=True)
    thread.start()

    print(f"worker started: {worker_name} env={env}")

    while True:
        print("worker heartbeat")
        time.sleep(10)


if __name__ == "__main__":
    main()