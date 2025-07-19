from datetime import datetime
from backend_python.chat.models import Chat, ChatParticipant, Message, MessageEdit
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select, update
import redis.asyncio as redis
from sqlalchemy.orm import aliased

redis_client = redis.Redis(host="localhost", port=6379, db=0, decode_responses=True)

ReplyMsg = aliased(Message)

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
        chat_id=1,
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

async def delete_global_message(
    db: AsyncSession,
    *,
    message_id: str
):
    await db.execute(
        update(Message)
        .where(Message.id == message_id, Message.chat_id == 1)
        .values(deleted=True)
    )
    await db.commit()

async def delete_private_message(
    db: AsyncSession,
    *,
    chat_key: str,
    message_id: str
):
    chat_id = await redis_client.get(f"chat_id:{chat_key}")
    if not chat_id:
        chat_id = await db.scalar(select(Chat.id).where(Chat.chat_key == chat_key))
        await redis_client.set(f"chat_id:{chat_key}", chat_id)

    await db.execute(
        update(Message)
        .where(Message.id == message_id, Message.chat_id == int(chat_id))
        .values(deleted=True)
    )
    await db.commit()

async def edit_global_message(
    db: AsyncSession,
    message_id: str,
    new_text: str
):
    chat_id = 1

    result = await db.execute(
        select(Message.content).where(
            Message.id == message_id,
            Message.chat_id == chat_id
        )
    )
    old_content = result.scalar_one_or_none()

    if old_content is None:
        return

    edit_entry = MessageEdit(
        message_id=message_id,
        old_text=old_content,
        edited_at=datetime.utcnow()
    )
    db.add(edit_entry)

    await db.execute(
        update(Message)
        .where(Message.id == message_id, Message.chat_id == chat_id)
        .values(
            is_edited=True,
            edited_at=datetime.utcnow(),
            content=new_text
        )
    )

    await db.commit()

async def edit_private_message(
    db: AsyncSession,
    *,
    chat_key: str,
    message_id: str,
    new_text: str
):
    chat_id = await redis_client.get(f"chat_id:{chat_key}")
    if not chat_id:
        chat_id = await db.scalar(
            select(Chat.id).where(Chat.chat_key == chat_key)
        )
        await redis_client.set(f"chat_id:{chat_key}", chat_id)

    chat_id = int(chat_id)

    result = await db.execute(
        select(Message.content).where(
            Message.id == message_id,
            Message.chat_id == chat_id
        )
    )
    old_content = result.scalar_one_or_none()

    if old_content is None:
        return

    edit_entry = MessageEdit(
        message_id=message_id,
        old_text=old_content,
        edited_at=datetime.utcnow()
    )
    db.add(edit_entry)

    await db.execute(
        update(Message)
        .where(Message.id == message_id, Message.chat_id == chat_id)
        .values(
            is_edited=True,
            edited_at=datetime.utcnow(),
            content=new_text
        )
    )

    await db.commit()

async def fetch_messages_by_chat_id(
    db: AsyncSession,
    chat_id: int,
    limit: int = 50,
    offset: int = 0,
):
    query = (
        select(
            Message,
            ReplyMsg.content.label("reply_to_text"),
            ReplyMsg.username.label("reply_to_user"),
        )
        .outerjoin(ReplyMsg, Message.reply_to_id == ReplyMsg.id)
        .where(Message.chat_id == chat_id, Message.deleted.is_(False))
        .order_by(Message.created_at.asc())
        .limit(limit)
        .offset(offset)
    )
    result = await db.execute(query)
    return result.all()