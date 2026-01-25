from typing import Dict

ALLOWED_ACTIONS = {"send", "edit", "delete", "pin"}

def validate_envelope(env: Dict[str, any], queue: str) -> bool:
    print(env)
    
    if env["kind"] != "message":
        print(f"[skip] not message kind (got={env['kind']}) queue={queue}")
        return False

    if env["action"] not in ALLOWED_ACTIONS:
        print(f"[skip] unsupported action={env['action']} queue={queue}")
        return False

    if not env["chat_id"] or not env["sender_id"]:
        print(f"[skip] missing chat_id or sender_id queue={queue}")
        return False

    return True
