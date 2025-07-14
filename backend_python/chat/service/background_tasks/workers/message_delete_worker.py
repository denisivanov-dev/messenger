import asyncio, json
import redis.asyncio as redis
from sqlalchemy.ext.asyncio import create_async_engine, async_sessionmaker, AsyncSession

from backend_python.config import DB_URL
from backend_python.chat.repository.message_repo import (
    delete_global_message, delete_private_message
)
from backend_python.chat.repository.chat_repo import get_private_chat_keys

engine = create_async_engine(DB_URL, echo=False, pool_pre_ping=True)
SessionFactory = async_sessionmaker(engine, expire_on_commit=False)
redis_client = redis.Redis(host="localhost", port=6379, db=0, decode_responses=True)

async def listen_to_delete_queue(queue: str):
    while True:
        _, data = await redis_client.blpop(queue)
        message = json.loads(data)

        async with SessionFactory() as db:
            await process_delete_message(queue, message, db)

async def process_delete_message(queue: str, message: dict, db: AsyncSession):
    print("Удаляю:", queue, message)

    if queue == "to_delete:global":
        try:
            await delete_global_message(db=db, message_id=message["message_id"])
        except Exception as e:
            print("❌ Ошибка при удалении глобального:", e)

    elif queue.startswith("to_delete:private:"):
        try:
            await delete_private_message(
                db=db,
                chat_key=message["chat_id"],
                message_id=message["message_id"]
            )
        except Exception as e:
            print("❌ Ошибка при удалении приватного:", e)

async def start_delete_listener():
    async with SessionFactory() as db:
        keys = await get_private_chat_keys(db)

    await asyncio.gather(
        listen_to_delete_queue("to_delete:global"),
        *[listen_to_delete_queue(f"to_delete:private:{k}") for k in keys]
    )