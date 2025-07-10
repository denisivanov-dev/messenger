from datetime import datetime
from backend_python.auth.models import User
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select
from typing import Tuple, List
import redis.asyncio as redis

redis_client = redis.Redis(host="localhost", port=6379, decode_responses=True)

async def get_username_by_id(db: AsyncSession, user_id: int) -> str | None:
    query = select(User.username).where(User.id == user_id)
    result = await db.execute(query)
    username = result.scalar_one_or_none()
    
    return username

async def get_all_users(db: AsyncSession) -> List[Tuple[int, str]]:
    query = select(User.id, User.username).where(User.is_active.is_(True))
    result = await db.execute(query)
    users = result.all()

    return users

async def get_all_users_from_redis():
        user_ids = await redis_client.smembers("users_all")
        pipe = redis_client.pipeline()

        for user_id in user_ids:
            pipe.hgetall(f"user:{user_id}")

        users = await pipe.execute()

        return [user for user in users if user]