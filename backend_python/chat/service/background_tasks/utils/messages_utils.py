from typing import Any, Dict, Optional

ACTIONS = ["send", "edit", "delete", "pin"]

def parse_envelope(raw: Dict[str, Any]) -> Dict[str, Any]:
    payload = raw.get("payload") or {}
    env = {
        "kind": raw.get("kind"),
        "action": raw.get("action"),
        "chat_id": raw.get("chat_id"),
        "chat_type": raw.get("chat_type"),
        "sender_id": raw.get("sender_id"),
        "target_id": raw.get("target_id"),
        "timestamp": raw.get("timestamp"),
        "message_id": payload.get("id"),
        "text": payload.get("text") or "",
        "attachments": payload.get("attachments") or [],
        "reply_to_id": payload.get("reply_to"),
        "raw": raw,
    }
    return env


def chat_key_from_queue(queue: str) -> Optional[str]:
    # type of queue: to_save:private:<chat_key>
    parts = queue.split(":", 2)
    if len(parts) == 3 and parts[1] == "private":
        return parts[2]
    return None

def validate_envelope(env: Dict[str, Any], queue: str) -> bool:
    if env["kind"] != "message":
        print(f"[skip] kind != 'message' (got: {env['kind']}). queue={queue}")
        return False

    if env["action"] not in ("send", "edit", "delete"):
        print(f"[skip] unsupported action '{env['action']}'. queue={queue}")
        return False

    if env["action"] not in ACTIONS:
        print(f"[skip] action='{env['action']}' not handled yet (todo). queue={queue}")
        return False

    return True