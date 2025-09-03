import asyncio
import json
from sqlalchemy.ext.asyncio import AsyncSession

from backend_python.chat.repository.message_repo import (
    save_message_to_global_chat,
    save_private_message,
)
from backend_python.chat.repository.chat_repo import get_private_chat_keys
from backend_python.core.redis_client import redis_client
from backend_python.core.db_client import SessionFactory

async def listen_to_save_queue(queue: str):
    while True:
        try:
            _, data = await redis_client.blpop(queue)
            message = json.loads(data)

            async with SessionFactory() as db:
                await process_save_message(queue, message, db)

        except Exception as e:
            print(f"Ошибка при получении или обработке сообщения из очереди [{queue}]:", e)

async def process_save_message(queue: str, message: dict, db: AsyncSession):
    print("Сохраняю:", queue, message)

    attachments = message.get("attachments", [])
    reply_to_id = message.get("reply_to")

    try:
        if queue == "to_save:global":
            await save_message_to_global_chat(
                db=db,
                message_id=message["message_id"],
                user_id=int(message["user_id"]),
                text=message["text"],
                timestamp_ms=message["timestamp"],
                reply_to_id=reply_to_id,
                attachments=attachments
            )

        elif queue.startswith("to_save:private:"):
            await save_private_message(
                db=db,
                chat_key=message["chat_id"],
                message_id=message["message_id"],
                sender_id=int(message["user_id"]),
                text=message["text"],
                timestamp_ms=message["timestamp"],
                reply_to_id=reply_to_id,
                attachments=attachments
            )

    except Exception as e:
        print(f"Ошибка при сохранении сообщения из {queue}:", e)

async def start_save_listener():
    async with SessionFactory() as db:
        keys = await get_private_chat_keys(db)

    await asyncio.gather(
        listen_to_save_queue("to_save:global"),
        *[listen_to_save_queue(f"to_save:private:{k}") for k in keys]
    )