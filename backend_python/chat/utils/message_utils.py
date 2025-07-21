from datetime import timezone
from typing import List, Dict
from backend_python.chat.models import Message

def orm_to_dto(msg: Message, receiver_id: int) -> dict:
    return {
        "message_id": msg.id,
        "user_id": str(msg.sender_id),
        "text": msg.content,
        "timestamp": int(msg.created_at.timestamp() * 1000),
        "username": msg.sender.username if msg.sender else None,
        "type": "private",
        "receiver_id": str(receiver_id),  # ← вот тут кастуем в строку
        "edited_at": int(msg.edited_at.timestamp() * 1000) if msg.is_edited and msg.edited_at else 0,
        "reply_to": msg.reply_to_id,
        "reply_to_text": msg.reply_to.content if msg.reply_to else None,
        "reply_to_user": msg.reply_to.sender.username if msg.reply_to and msg.reply_to.sender else None,
        "pinned": msg.is_pinned,
    }