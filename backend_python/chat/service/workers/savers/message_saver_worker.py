import asyncio, json
import redis.asyncio as redis
from sqlalchemy.ext.asyncio import create_async_engine, async_sessionmaker, AsyncSession

from backend_python.config import DB_URL

from backend_python.chat.repository.message_repo import (
    save_message_to_global_chat, save_private_message
)
from backend_python.chat.repository.chat_repo import get_private_chat_queues

engine = create_async_engine(DB_URL, echo=False, pool_pre_ping=True)
SessionFactory = async_sessionmaker(engine, expire_on_commit=False)

redis_client = redis.Redis(host="localhost", port=6379, db=0, decode_responses=True)

async def listen_to_queue(queue: str):
    while True:
        _, data = await redis_client.blpop(queue)
        message = json.loads(data)
        
        async with SessionFactory() as db:
            await process_message(queue, message, db)

async def process_message(queue: str, message: dict, db: AsyncSession):
    print("Обрабатываю:", queue, message)

    if queue == "to_save:global":
        try:
            await save_message_to_global_chat(
                db=db,
                message_id=message["message_id"],
                user_id=int(message["user_id"]),
                text=message["text"],
                timestamp_ms=message["timestamp"]
            )
        except Exception as e:
            print("❌ Ошибка при сохранении сообщения:", e)

    elif queue.startswith("to_save:private"):
        try:
            chat_key = message["chat_id"]
            await save_private_message(
                db=db,
                chat_key=chat_key,
                message_id=message["message_id"],
                sender_id=int(message["user_id"]),
                text=message["text"],
                timestamp_ms=message["timestamp"]
            )
        except Exception as e:
            print("❌ Ошибка при сохранении приватного сообщения:", e)

async def start_redis_listener():
    async with SessionFactory() as db:
        private_queues = await get_private_chat_queues(db)
    
    await asyncio.gather(
        listen_to_queue("to_save:global"),
        *[listen_to_queue(queue) for queue in private_queues]
    )