import asyncio
import json
from sqlalchemy.ext.asyncio import AsyncSession

from backend_python.chat.repository.message_repo import (
    pin_global_message,
    pin_private_message
)
from backend_python.chat.repository.chat_repo import get_private_chat_keys
from backend_python.core.redis_client import redis_client
from backend_python.core.db_client import SessionFactory

async def listen_to_pin_queue(queue: str):
    while True:
        try:
            _, data = await redis_client.blpop(queue)
            message = json.loads(data)

            async with SessionFactory() as db:
                await process_pin_message(queue, message, db)

        except Exception as e:
            print(f"Ошибка в очереди [{queue}]:", e)

async def process_pin_message(queue: str, message: dict, db: AsyncSession):
    print("Пиную:", queue, message)

    pin = message.get("pinned", False)

    try:
        if queue == "to_pin:global":
            await pin_global_message(
                db=db,
                message_id=message["message_id"],
                pinned=pin
            )

        elif queue.startswith("to_pin:private:"):
            await pin_private_message(
                db=db,
                chat_key=message["chat_id"],
                message_id=message["message_id"],
                pinned=pin
            )

    except Exception as e:
        print(f"Ошибка при пине сообщения из {queue}:", e)

async def start_pin_listener():
    async with SessionFactory() as db:
        keys = await get_private_chat_keys(db)

    await asyncio.gather(
        listen_to_pin_queue("to_pin:global"),
        *[listen_to_pin_queue(f"to_pin:private:{k}") for k in keys]
    )