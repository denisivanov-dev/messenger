import asyncio
import json
from sqlalchemy.ext.asyncio import AsyncSession

from backend_python.chat.repository.friends_repo import (
    on_request_sent,
    on_request_canceled,
    on_request_declined,
    on_request_accepted,
    on_friend_removed,
)
from backend_python.core.redis_client import redis_client
from backend_python.core.db_client import SessionFactory

async def listen_to_friend_queue(queue_name: str):
    while True:
        try:
            _, raw_data = await redis_client.blpop(queue_name)
            event_data = json.loads(raw_data)

            async with SessionFactory() as db:
                await process_friend_event(queue_name, event_data, db)
                await db.commit()

        except Exception as e:
            print(f"Ошибка в очереди [{queue_name}]:", e)

async def process_friend_event(queue_name: str, event_data: dict, db: AsyncSession):
    print("🤝 Обрабатываю событие:", queue_name, event_data)

    event_type = event_data.get("type")
    from_user_id = event_data.get("from_id")
    to_user_id = event_data.get("to_id")

    if not event_type or from_user_id is None or to_user_id is None:
        print("Некорректное событие:", event_data)
        return

    try:
        from_user_id = int(from_user_id)
        to_user_id = int(to_user_id)
    except ValueError:
        print("from_id/to_id не числа:", event_data)
        return

    try:
        match event_type:
            case "friend_request_sent":
                await on_request_sent(db, from_user_id, to_user_id)
            case "friend_request_canceled":
                await on_request_canceled(db, from_user_id, to_user_id)
            case "friend_request_declined":
                await on_request_declined(db, from_user_id, to_user_id)
            case "friend_request_accepted":
                await on_request_accepted(db, from_user_id, to_user_id)
            case "friend_removed":
                await on_friend_removed(db, from_user_id, to_user_id)
            case _:
                print("Неизвестный тип события:", event_type)

    except Exception as error:
        print("Ошибка при обработке события:", error)

async def start_friend_events_listener():
    await listen_to_friend_queue("to_save:friend_events")
