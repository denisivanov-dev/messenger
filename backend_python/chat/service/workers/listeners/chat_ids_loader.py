import asyncio, json
import redis.asyncio as redis
from sqlalchemy.ext.asyncio import create_async_engine, async_sessionmaker, AsyncSession

from backend_python.config import DB_URL
from backend_python.chat.repository.user_repo import get_all_users

engine = create_async_engine(DB_URL, echo=False, pool_pre_ping=True)
SessionFactory = async_sessionmaker(engine, expire_on_commit=False)

redis_client = redis.Redis(host="localhost", port=6379, db=0, decode_responses=True)

GLOBAL_CHAT_ID = 1  

async def load_all_chat_ids():
    async with SessionFactory() as db:
        user_list = await get_all_users(db) 

    pipe = redis_client.pipeline()

    for user_id, _ in user_list:
        uid = str(user_id)
        pipe.sadd(f"user:{uid}:chats", GLOBAL_CHAT_ID)

    await pipe.execute()

async def run_load_all_chat_ids():
    print("Chat IDs loader worker started")
    await load_all_chat_ids()
    print("Chat IDs loader worker finished")
