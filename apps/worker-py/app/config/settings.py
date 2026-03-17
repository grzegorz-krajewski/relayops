from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    app_env: str = "local"

    redis_addr: str = "redis:6379"
    redis_stream_name: str = "tasks.stream"
    worker_name: str = "worker-1"
    worker_group: str = "worker-group"

    model_config = SettingsConfigDict(
        env_file=".env",
        extra="ignore",
        case_sensitive=False,
    )