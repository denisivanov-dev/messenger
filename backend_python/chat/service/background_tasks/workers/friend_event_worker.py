# backend_python/service/workers/friend_events_worker.py
import asyncio
import json
import redis.asyncio as redis
from sqlalchemy.ext.asyncio import create_async_engine, async_sessionmaker, AsyncSession

from backend_python.config import DB_URL
from backend_python.chat.repository.friends_repo import (
    on_request_sent,
    on_request_canceled,
    on_request_declined,
    on_request_accepted,
    on_friend_removed,
)

engine = create_async_engine(DB_URL, echo=False, pool_pre_ping=True)
SessionFactory = async_sessionmaker(engine, expire_on_commit=False)
redis_client = redis.Redis(host="localhost", port=6379, db=0, decode_responses=True)

async def listen_to_friend_queue(queue_name: str):
    while True:
        _, raw_data = await redis_client.blpop(queue_name)
        event_data = json.loads(raw_data)

        async with SessionFactory() as db:
            await process_friend_event(queue_name, event_data, db)
            await db.commit()

async def process_friend_event(queue_name: str, event_data: dict, db: AsyncSession):
    print("Сохраняю:", queue_name, event_data)

    event_type = event_data.get("type")
    from_user_id = event_data.get("from_id")
    to_user_id = event_data.get("to_id")

    if not event_type or from_user_id is None or to_user_id is None:
        print("⚠️ Некорректное событие:", event_data)
        return

    try:
        from_user_id = int(from_user_id)
        to_user_id = int(to_user_id)
    except ValueError:
        print("⚠️ from_id/to_id не числа:", event_data)
        return

    try:
        if event_type == "friend_request_sent":
            await on_request_sent(db, from_user_id, to_user_id)
        elif event_type == "friend_request_canceled":
            await on_request_canceled(db, from_user_id, to_user_id)
        elif event_type == "friend_request_declined":
            await on_request_declined(db, from_user_id, to_user_id)
        elif event_type == "friend_request_accepted":
            await on_request_accepted(db, from_user_id, to_user_id)
        elif event_type == "friend_removed":
            await on_friend_removed(db, from_user_id, to_user_id)
        else:
            print("⚠️ Неизвестный тип события:", event_type)
    except Exception as error:
        print("❌ Ошибка при обработке события:", error)

async def start_friend_events_listener():
    await listen_to_friend_queue("to_save:friend_events")