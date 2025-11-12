from backend_python.chat.service.background_tasks.envelope.handlers import (
    message_send
)

handlers = {
    ("message", "send"): message_send.handle,
    # ("message", "edit"): message_edit.handle,
    # ("message", "delete"): message_delete.handle,
    # ("message", "pin"): message_pin.handle,
}

async def dispatch_envelope(env, db):
    key = (env["kind"], env["action"])
    handler = handlers.get(key)

    if not handler:
        print(f"[warn] no handler for {key}")
        return
    
    await handler(env, db)
