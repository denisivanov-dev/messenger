import asyncio, json
import redis.asyncio as redis
from sqlalchemy.ext.asyncio import create_async_engine, async_sessionmaker, AsyncSession

from backend_python.config import DB_URL
from backend_python.chat.repository.message_repo import fetch_messages_by_chat_id
from backend_python.chat.utils.user_utils import fetch_usernames_map

engine = create_async_engine(DB_URL, echo=False, pool_pre_ping=True)
SessionFactory = async_sessionmaker(engine, expire_on_commit=False)

redis_client = redis.Redis(host="localhost", port=6379, db=0, decode_responses=True)

async def preload_global_chat_history():
    chat_id = 1
    limit = 500

    async with SessionFactory() as db:
        messages = await fetch_messages_by_chat_id(db, chat_id, limit)
        usernames_map = await fetch_usernames_map(db, messages)

    values = []
    for msg, reply_text, reply_user in messages:
        attachments_list = []
        for attachment in msg.attachments:
            attachments_list.append({
                "key": attachment.filename,
                "type": attachment.filetype,
                "size": attachment.filesize,
                "original_name": attachment.original_name,
            })

        msg_dict = {
            "message_id": msg.id,
            "user_id": str(msg.sender_id),
            "text": msg.content,
            "timestamp": int(msg.created_at.timestamp() * 1000),
            "username": usernames_map.get(msg.sender_id),
            "type": "global",
            "receiver_id": "",
            "edited_at": int(msg.edited_at.timestamp() * 1000) if msg.is_edited and msg.edited_at else 0,
            "reply_to": msg.reply_to_id,
            "reply_to_text": reply_text,
            "reply_to_user": reply_user,
            "pinned": msg.is_pinned,
            "attachments": attachments_list,
        }

        values.append(json.dumps(msg_dict))

    if values:
        await redis_client.delete(f"chat:history:{chat_id}")
        await redis_client.rpush(f"chat:history:{chat_id}", *values)
        await redis_client.ltrim(f"chat:history:{chat_id}", -1000, -1)

async def run_preload_global_chat_history():
    print("Message loader worker started")
    await preload_global_chat_history()
    print("Message loader worker finished")