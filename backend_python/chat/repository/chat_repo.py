from sqlalchemy import select, func
from backend_python.chat.models import Chat, ChatParticipant, Message
from sqlalchemy.ext.asyncio import AsyncSession
from datetime import datetime
from sqlalchemy.orm import selectinload, selectinload

from backend_python.chat.utils import user_utils, message_utils

async def get_or_create_private_chat(db: AsyncSession, user1_id: int, user2_id: int):
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
            .options(selectinload(Message.sender))
            .where(Message.chat_id == chat.id)
            .order_by(Message.created_at)
        )
        messages_result = await db.execute(messages_query)
        messages = messages_result.scalars().all()
        messages_dto = [message_utils.orm_to_dto(msg=m, receiver_id=user2_id) for m in messages]
        
        return chat, messages_dto

    chat_key = user_utils.generate_private_chat_key(user1_id, user2_id)

    new_chat = Chat(type='private', title=None, created_at=datetime.utcnow(), chat_key=chat_key)
    db.add(new_chat)
    await db.flush()

    participants = [
        ChatParticipant(chat_id=new_chat.id, user_id=user1_id),
        ChatParticipant(chat_id=new_chat.id, user_id=user2_id)
    ]
    db.add_all(participants)
    await db.commit()

    return new_chat, []

async def get_private_chat_queues(db: AsyncSession) -> list[str]:
    query = select(Chat.chat_key).where(Chat.type == "private")
    result = await db.execute(query)
    chat_keys = result.scalars().all()

    return [f"to_save:private:{key}" for key in chat_keys if key]