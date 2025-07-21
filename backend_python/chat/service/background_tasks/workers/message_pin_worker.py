import asyncio
import json
import redis.asyncio as redis
from sqlalchemy.ext.asyncio import create_async_engine, async_sessionmaker, AsyncSession

from backend_python.config import DB_URL
from backend_python.chat.repository.message_repo import (
    pin_global_message, pin_private_message
)
from backend_python.chat.repository.chat_repo import get_private_chat_keys

engine = create_async_engine(DB_URL, echo=False, pool_pre_ping=True)
SessionFactory = async_sessionmaker(engine, expire_on_commit=False)
redis_client = redis.Redis(host="localhost", port=6379, db=0, decode_responses=True)

async def listen_to_pin_queue(queue: str):
    while True:
        _, data = await redis_client.blpop(queue)
        message = json.loads(data)

        async with SessionFactory() as db:
            await process_pin_message(queue, message, db)

async def process_pin_message(queue: str, message: dict, db: AsyncSession):
    print("📌 Пиную:", queue, message)

    pin = message.get("pinned", False)

    if queue == "to_pin:global":
        try:
            await pin_global_message(
                db=db,
                message_id=message["message_id"],
                pinned=pin
            )
        except Exception as e:
            print("❌ Ошибка при пине глобального:", e)

    elif queue.startswith("to_pin:private:"):
        try:
            await pin_private_message(
                db=db,
                chat_key=message["chat_id"],
                message_id=message["message_id"],
                pinned=pin
            )
        except Exception as e:
            print("❌ Ошибка при пине приватного:", e)

async def start_pin_listener():
    async with SessionFactory() as db:
        keys = await get_private_chat_keys(db)

    await asyncio.gather(
        listen_to_pin_queue("to_pin:global"),
        *[listen_to_pin_queue(f"to_pin:private:{k}") for k in keys]
    )

if __name__ == "__main__":
    asyncio.run(start_pin_listener())