from datetime import timezone
from typing import List, Dict
from backend_python.chat.models import Message

def orm_to_dto(msg: Message, receiver_id: int) -> dict:
    return {
        "message_id": msg.id,
        "text":       msg.content,
        "timestamp":  int(msg.created_at.replace(tzinfo=timezone.utc).timestamp() * 1000),
        "username":   getattr(msg.sender, "username", "неизвестный"),
        "user_id":    str(msg.sender_id),
    }