from prometheus_client import Counter, Histogram

tasks_processed_total = Counter(
    "relayops_worker_tasks_processed_total",
    "Total number of processed tasks.",
)

tasks_failed_total = Counter(
    "relayops_worker_tasks_failed_total",
    "Total number of failed tasks.",
)

tasks_acked_total = Counter(
    "relayops_worker_tasks_acked_total",
    "Total number of acknowledged stream messages.",
)

grpc_calls_total = Counter(
    "relayops_worker_grpc_calls_total",
    "Total number of gRPC task processor calls.",
)

transient_retries_total = Counter(
    "relayops_worker_transient_retries_total",
    "Total number of transient retry attempts.",
)

task_processing_duration_seconds = Histogram(
    "relayops_worker_task_processing_duration_seconds",
    "Task processing duration in seconds.",
)