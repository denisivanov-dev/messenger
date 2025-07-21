from typing import List, Dict
from sqlalchemy.ext.asyncio import AsyncSession
from backend_python.chat.repository.user_repo import get_username_by_id

async def fetch_usernames_map(db: AsyncSession, messages: List) -> Dict[int, str | None]:
    unique_ids = []
    seen = set()
    for msg, *_ in messages:
        if msg.sender_id not in seen:
            unique_ids.append(msg.sender_id)
            seen.add(msg.sender_id)

    usernames_map = {}
    for user_id in unique_ids:
        username = await get_username_by_id(db, user_id)
        usernames_map[user_id] = username
        
    return usernames_map

def generate_private_chat_key(user_a: int, user_b: int) -> str:
    return f"{min(user_a, user_b)}:{max(user_a, user_b)}"