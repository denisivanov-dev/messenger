from datetime import datetime
from backend_python.chat.models import Chat, ChatParticipant, Message
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select
import redis.asyncio as redis

redis_client = redis.Redis(host="localhost", port=6379, db=0, decode_responses=True)

async def save_message_to_global_chat(
    db: AsyncSession,
    *,
    message_id: str,
    user_id: int,
    text: str,
    timestamp_ms: int
):
    msg = Message(
        id=message_id,
        chat_id=1,  # глобальный чат — id всегда 1
        sender_id=user_id,
        content=text,
        created_at=datetime.fromtimestamp(timestamp_ms / 1000.0)
    )
    
    db.add(msg)
    await db.commit()

async def save_private_message(
    db: AsyncSession,
    *,
    message_id: str,
    chat_key: int,
    sender_id: int,
    text: str,
    timestamp_ms: int
):
    chat_id = await redis_client.get(f"chat_id:{chat_key}")
    if not chat_id:
        chat_id = await db.scalar(select(Chat.id).where(Chat.chat_key == chat_key))
        await redis_client.set(f"chat_id:{chat_key}", chat_id)
    
    msg = Message(
        id=message_id,
        chat_id=int(chat_id),
        sender_id=sender_id,
        content=text,
        created_at=datetime.fromtimestamp(timestamp_ms / 1000.0)
    )

    db.add(msg)
    await db.commit()

async def fetch_messages_by_chat_id(
    db: AsyncSession,
    chat_id: int,
    limit: int = 50,
    offset: int = 0,
):
    query = (
        select(Message)
        .where(Message.chat_id == chat_id, Message.deleted.is_(False))
        .order_by(Message.created_at.asc())
        .limit(limit)
        .offset(offset)
    )
    result = await db.execute(query )

    return result.scalars().all()