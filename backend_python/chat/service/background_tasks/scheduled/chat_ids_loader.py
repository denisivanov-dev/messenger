import asyncio
import json
from backend_python.chat.repository.user_repo import get_all_users
from backend_python.core.redis_client import redis_client
from backend_python.core.db_client import SessionFactory  # ✅ заменено

GLOBAL_CHAT_ID = 1  

async def load_all_chat_ids():
    async with SessionFactory() as db:
        user_list = await get_all_users(db)

    pipe = redis_client.pipeline()

    for user_id, *_ in user_list:
        uid = str(user_id)
        pipe.sadd(f"user:{uid}:chats", GLOBAL_CHAT_ID)

    await pipe.execute()

async def run_load_all_chat_ids():
    print("Chat IDs loader worker started")
    await load_all_chat_ids()
    print("Chat IDs loader worker finished")