from datetime import datetime
from backend_python.chat.models import (
    Chat, ChatParticipant, Message, MessageEdit, User, Attachment
)
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select, update
from sqlalchemy.orm import aliased, selectinload

from backend_python.core.redis_client import redis_client

ReplyMsg = aliased(Message)
ReplyUser = aliased(User)

async def save_message_to_global_chat(
    db: AsyncSession,
    *,
    message_id: str,
    user_id: int,
    text: str,
    timestamp_ms: int,
    reply_to_id: str = None,
    attachments: list
):
    if attachments and text.strip():
        msg_type = "text/image"
    elif attachments:
        msg_type = "image"
    elif text.strip():
        msg_type = "text"
    else:
        msg_type = "empty"

    msg = Message(
        id=message_id,
        chat_id=1,
        sender_id=user_id,
        type=msg_type,
        content=text,
        reply_to_id=reply_to_id,
        created_at=datetime.fromtimestamp(timestamp_ms / 1000.0),
    )

    db.add(msg)
    await db.flush()

    for att in attachments:
        db.add(Attachment(
            message_id=message_id,
            filename=att["key"],
            filetype=att["type"],
            filesize=att["size"],
            original_name=att["original_name"]
        ))

    await db.commit()

async def save_private_message(
    db: AsyncSession,
    *,
    message_id: str,
    chat_key: str,
    sender_id: int,
    text: str,
    timestamp_ms: int,
    reply_to_id: str = None,
    attachments: list
):
    chat_id = await redis_client.get(f"chat_id:{chat_key}")
    if not chat_id:
        chat_id = await db.scalar(select(Chat.id).where(Chat.chat_key == chat_key))
        await redis_client.set(f"chat_id:{chat_key}", chat_id)

    if attachments and text.strip():
        msg_type = "text/image"
    elif attachments:
        msg_type = "image"
    elif text.strip():
        msg_type = "text"
    else:
        msg_type = "empty"

    msg = Message(
        id=message_id,
        chat_id=int(chat_id),
        sender_id=sender_id,
        type=msg_type,
        content=text,
        reply_to_id=reply_to_id,
        created_at=datetime.fromtimestamp(timestamp_ms / 1000.0),
    )

    db.add(msg)
    await db.flush()

    for att in attachments:
        db.add(Attachment(
            message_id=message_id,
            filename=att["key"],
            filetype=att["type"],
            filesize=att["size"],
            original_name=att["original_name"]
        ))

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

async def pin_global_message(
    db: AsyncSession,
    *,
    message_id: str,
    pinned: bool
):
    await db.execute(
        update(Message)
        .where(Message.id == message_id, Message.chat_id == 1)
        .values(is_pinned=pinned)
    )
    await db.commit()

async def pin_private_message(
    db: AsyncSession,
    *,
    chat_key: str,
    message_id: str,
    pinned: bool
):
    chat_id = await redis_client.get(f"chat_id:{chat_key}")
    if not chat_id:
        chat_id = await db.scalar(
            select(Chat.id).where(Chat.chat_key == chat_key)
        )
        await redis_client.set(f"chat_id:{chat_key}", chat_id)

    await db.execute(
        update(Message)
        .where(Message.id == message_id, Message.chat_id == int(chat_id))
        .values(is_pinned=pinned)
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
            ReplyUser.username.label("reply_to_user"),
        )
        .outerjoin(ReplyMsg, Message.reply_to_id == ReplyMsg.id)
        .outerjoin(ReplyUser, ReplyMsg.sender_id == ReplyUser.id)
        .where(Message.chat_id == chat_id, Message.deleted.is_(False))
        .options(selectinload(Message.attachments))
        .order_by(Message.created_at.asc())
        .limit(limit)
        .offset(offset)
    )
    result = await db.execute(query)
    return result.all()