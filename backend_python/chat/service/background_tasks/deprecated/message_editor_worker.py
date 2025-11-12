import asyncio
import json
from sqlalchemy.ext.asyncio import AsyncSession

from backend_python.chat.repository.message_repo import (
    edit_global_message,
    edit_private_message,
)
from backend_python.chat.repository.chat_repo import get_private_chat_keys
from backend_python.core.redis_client import redis_client
from backend_python.core.db_client import SessionFactory

async def listen_to_edit_queue(queue: str):
    while True:
        try:
            _, data = await redis_client.blpop(queue)
            message = json.loads(data)

            async with SessionFactory() as db:
                await process_edit_message(queue, message, db)

        except Exception as e:
            print(f"Ошибка в очереди [{queue}]:", e)

async def process_edit_message(queue: str, message: dict, db: AsyncSession):
    print("Редактирую:", queue, message)

    try:
        if queue == "to_edit:global":
            await edit_global_message(
                db=db,
                message_id=message["message_id"],
                new_text=message["text"]
            )

        elif queue.startswith("to_edit:private:"):
            await edit_private_message(
                db=db,
                chat_key=message["chat_id"],
                message_id=message["message_id"],
                new_text=message["text"]
            )

    except Exception as e:
        print(f"Ошибка при редактировании сообщения из {queue}:", e)

async def start_edit_listener():
    async with SessionFactory() as db:
        keys = await get_private_chat_keys(db)

    await asyncio.gather(
        listen_to_edit_queue("to_edit:global"),
        *[listen_to_edit_queue(f"to_edit:private:{k}") for k in keys]
    )