import asyncio
import json
from sqlalchemy.ext.asyncio import AsyncSession

from backend_python.chat.repository.message_repo import (
    delete_global_message,
    delete_private_message,
)
from backend_python.chat.repository.chat_repo import get_private_chat_keys
from backend_python.core.redis_client import redis_client
from backend_python.core.db_client import SessionFactory

async def listen_to_delete_queue(queue: str):
    while True:
        try:
            _, data = await redis_client.blpop(queue)
            message = json.loads(data)

            async with SessionFactory() as db:
                await process_delete_message(queue, message, db)

        except Exception as e:
            print(f"Ошибка в очереди [{queue}]:", e)

async def process_delete_message(queue: str, message: dict, db: AsyncSession):
    print("Удаляю:", queue, message)

    try:
        if queue == "to_delete:global":
            await delete_global_message(db=db, message_id=message["message_id"])

        elif queue.startswith("to_delete:private:"):
            await delete_private_message(
                db=db,
                chat_key=message["chat_id"],
                message_id=message["message_id"]
            )

    except Exception as e:
        print(f"Ошибка при удалении сообщения из {queue}:", e)

async def start_delete_listener():
    async with SessionFactory() as db:
        keys = await get_private_chat_keys(db)

    await asyncio.gather(
        listen_to_delete_queue("to_delete:global"),
        *[listen_to_delete_queue(f"to_delete:private:{k}") for k in keys]
    )