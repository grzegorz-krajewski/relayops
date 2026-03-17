import json
from typing import Any


def process_task(task_type: str, raw_payload: str) -> dict[str, Any]:
    payload = json.loads(raw_payload)

    if task_type == "normalize_payload":
        text = str(payload.get("text", ""))
        normalized = " ".join(text.split())
        return {
            "status": "processed",
            "task_type": task_type,
            "normalized_text": normalized,
            "original_payload": payload,
        }

    if task_type == "enrich_text":
        text = str(payload.get("text", ""))
        return {
            "status": "processed",
            "task_type": task_type,
            "enriched_text": f"{text} :: enriched",
            "original_payload": payload,
        }

    if task_type == "classify_priority":
        text = str(payload.get("text", "")).lower()
        priority = "high" if "urgent" in text or "asap" in text else "normal"
        return {
            "status": "processed",
            "task_type": task_type,
            "priority": priority,
            "original_payload": payload,
        }

    return {
        "status": "processed",
        "task_type": task_type,
        "note": "unknown task type, no-op",
        "original_payload": payload,
    }