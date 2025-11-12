import asyncio
import json
from sqlalchemy.ext.asyncio import AsyncSession

from backend_python.chat.repository.message_repo import fetch_messages_by_chat_id
from backend_python.chat.utils.user_utils import fetch_usernames_map
from backend_python.core.redis_client import redis_client
from backend_python.core.db_client import SessionFactory

async def preload_global_chat_history():
    chat_id = 1
    limit = 500

    async with SessionFactory() as db:
        messages = await fetch_messages_by_chat_id(db, chat_id, limit)
        usernames_map = await fetch_usernames_map(db, messages)

    envelopes = []

    for msg, reply_text, reply_user in messages:
        attachments_list = []

        if msg.attachments and isinstance(msg.attachments, list):
            attachments_list.extend(msg.attachments)

        if msg.attachments_db:
            for att in msg.attachments_db:
                attachments_list.append({
                    "key": att.key,
                    "type": att.type,
                    "size": att.size,
                    "mime": att.mime,
                    "original_name": att.original_name,
                    "width": att.width,
                    "height": att.height,
                    "duration": att.duration,
                })

        payload = {
            "message_id": msg.id,
            "text": msg.text or "",
            "reply_to": msg.reply_to_id,
            "reply_to_text": reply_text,
            "reply_to_user": reply_user,
            "pinned": msg.is_pinned,
            "is_edited": msg.is_edited,
            "edited_at": int(msg.edited_at.timestamp() * 1000) if msg.is_edited and msg.edited_at else 0,
            "attachments": attachments_list,
            "username": usernames_map.get(msg.sender_id),
        }

        envelope = {
            "kind": msg.kind or "message",
            "action": msg.action or "send",
            "chat_id": str(msg.chat_id),
            "chat_type": msg.chat_type or "global",
            "sender_id": str(msg.sender_id or ""),
            "target_id": str(msg.target_id or ""),
            "timestamp": int(msg.timestamp or int(msg.created_at.timestamp() * 1000)),
            "payload": payload,
        }

        envelopes.append(json.dumps(envelope, ensure_ascii=False))

    if envelopes:
        key = f"chat:history:{chat_id}"
        await redis_client.delete(key)
        await redis_client.rpush(key, *envelopes)
        await redis_client.ltrim(key, -1000, -1)

    print(f"[OK] Redis preload complete: {len(envelopes)} messages loaded for chat_id={chat_id}")


async def run_preload_global_chat_history():
    print("Message loader worker started")
    await preload_global_chat_history()
    print("Message loader worker finished")
