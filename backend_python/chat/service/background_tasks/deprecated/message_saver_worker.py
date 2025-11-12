import asyncio
import json
from typing import Any, Dict, Optional

from sqlalchemy.ext.asyncio import AsyncSession

from backend_python.chat.repository.message_repo import (
    save_message_to_global_chat,
    save_private_message,
)
from backend_python.chat.repository.chat_repo import get_private_chat_keys
from backend_python.core.redis_client import redis_client
from backend_python.core.db_client import SessionFactory
from backend_python.chat.service.background_tasks.utils import messages_utils

async def process_save_envelope(queue: str, envelope_raw: Dict[str, Any], db: AsyncSession):
    env = messages_utils.parse_envelope(envelope_raw)

    if env["kind"] != "message":
        print(f"[skip] kind != 'message' (got: {env['kind']}). queue={queue}")
        return

    if env["action"] not in ("send", "edit", "delete"):
        print(f"[skip] unsupported action '{env['action']}'. queue={queue}")
        return

    if env["action"] != "send":
        print(f"[skip] action='{env['action']}' not handled yet (todo). queue={queue}")
        return

    chat_type = env["chat_type"]
    try:
        if chat_type == "global" or queue == "to_save:global":
            await save_message_to_global_chat(
                db=db,
                message_id=env["message_id"],
                user_id=int(env["sender_id"]) if env["sender_id"] is not None else None,
                text=env["text"],
                timestamp_ms=env["timestamp"],
                reply_to_id=env["reply_to_id"],
                attachments=env["attachments"],
                raw_envelope=envelope_raw,
            )
            print(f"[ok] saved global message id={env['message_id']}")

        elif chat_type == "private" or queue.startswith("to_save:private:"):
            chat_key = messages_utils.chat_key_from_queue(queue) or env["chat_id"]
            if not chat_key:
                print(f"[warn] no chat_key for private message; queue={queue}, env.chat_id={env['chat_id']}")
                return

            await save_private_message(
                db=db,
                chat_key=chat_key,
                message_id=env["message_id"],
                sender_id=int(env["sender_id"]) if env["sender_id"] is not None else None,
                text=env["text"],
                timestamp_ms=env["timestamp"],
                reply_to_id=env["reply_to_id"],
                attachments=env["attachments"],
                raw_envelope=envelope_raw,
            )
            print(f"[ok] saved private message id={env['message_id']} chat_key={chat_key}")

        else:
            print(f"[warn] unknown chat_type='{chat_type}' queue={queue}")
            return

    except Exception as e:
        print(f"[ERR] failed to save from queue={queue}: {e}")


async def listen_to_save_queue(queue: str):
    print(f"[init] listening queue: {queue}")
    while True:
        try:
            _, data = await redis_client.blpop(queue)
            envelope = json.loads(data)

            async with SessionFactory() as db:
                await process_save_envelope(queue, envelope, db)
        except Exception as e:
            print(f"[ERR] queue={queue} recv/handle error: {e}")


async def start_save_listener():
    async with SessionFactory() as db:
        keys = await get_private_chat_keys(db)

    await asyncio.gather(
        listen_to_save_queue("to_save:global"),
        *[listen_to_save_queue(f"to_save:private:{k}") for k in keys],
    )