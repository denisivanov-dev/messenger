from sqlalchemy import select, func
from backend_python.chat.models import Chat, ChatParticipant, Message
from sqlalchemy.ext.asyncio import AsyncSession
from datetime import datetime
from sqlalchemy.orm import selectinload
import json

from backend_python.core.redis_client import redis_client
from backend_python.chat.utils import user_utils, message_utils

async def get_or_create_private_chat(db: AsyncSession, user1_id: int, user2_id: int):
    chat_key = user_utils.generate_private_chat_key(user1_id, user2_id)
    redis_queue_key = f"chat:history:{chat_key}"

    raw_messages = await redis_client.lrange(redis_queue_key, 0, -1)
    if raw_messages:
        await redis_client.expire(redis_queue_key, 900)
        messages_dto = [json.loads(m) for m in raw_messages]

        chat_result = await db.execute(select(Chat).where(Chat.chat_key == chat_key))
        chat = chat_result.scalars().first()

        return chat, messages_dto

    query = (
        select(Chat)
        .join(ChatParticipant)
        .where(
            Chat.type == 'private',
            ChatParticipant.user_id.in_([user1_id, user2_id])
        )
        .group_by(Chat.id)
        .having(func.count(Chat.id) == 2)
    )
    result = await db.execute(query)
    chat = result.scalars().first()

    if chat:
        messages_query = (
            select(Message)
            .options(
                selectinload(Message.sender),
                selectinload(Message.attachments),
                selectinload(Message.reply_to).selectinload(Message.sender),
            )
            .where(
                Message.chat_id == chat.id,
                Message.deleted == False
            )
            .order_by(Message.created_at)
        )
        messages_result = await db.execute(messages_query)
        messages = messages_result.scalars().all()
        messages_dto = [
            message_utils.orm_to_dto(msg=m, receiver_id=user2_id) for m in messages
        ]

        if messages_dto:
            redis_pipe = redis_client.pipeline()
            redis_pipe.rpush(redis_queue_key, *[json.dumps(m) for m in messages_dto])
            redis_pipe.expire(redis_queue_key, 900)
            await redis_pipe.execute()

        await redis_client.set(f"chat_id:{chat_key}", chat.id)
        await redis_client.sadd(f"chat:{chat.id}:participants", user1_id, user2_id)

        return chat, messages_dto

    new_chat = Chat(
        type='private',
        title=None,
        created_at=datetime.utcnow(),
        chat_key=chat_key
    )
    db.add(new_chat)
    await db.flush()

    participants = [
        ChatParticipant(chat_id=new_chat.id, user_id=user1_id),
        ChatParticipant(chat_id=new_chat.id, user_id=user2_id)
    ]
    db.add_all(participants)
    await db.commit()

    await redis_client.set(f"chat_id:{chat_key}", new_chat.id)
    await redis_client.sadd(f"chat:{new_chat.id}:participants", user1_id, user2_id)

    return new_chat, []

async def get_private_chat_keys(db: AsyncSession) -> list[str]:
    query = select(Chat.chat_key).where(Chat.type == "private")
    result = await db.execute(query)
    chat_keys = result.scalars().all()
    return [key for key in chat_keys if key]