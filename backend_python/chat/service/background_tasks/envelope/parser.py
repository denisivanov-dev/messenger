from typing import Any, Dict

def parse_envelope(raw: Dict[str, Any]) -> Dict[str, Any]:
    payload = raw.get("payload") or {}

    return {
        "kind": raw.get("kind"),
        "action": raw.get("action"),
        "chat_id": int(raw.get("chat_id")) if raw.get("chat_id") else None,
        "chat_type": raw.get("chat_type"),
        "sender_id": int(raw.get("sender_id")) if raw.get("sender_id") else None,
        "target_id": int(raw.get("target_id")) if raw.get("target_id") else None,
        "timestamp": raw.get("timestamp"),
        "message_id": payload.get("message_id"),
        "text": payload.get("text") or "",
        "attachments": payload.get("attachments") or [],
        "reply_to_id": payload.get("reply_to"),
        "raw": raw,
    }
