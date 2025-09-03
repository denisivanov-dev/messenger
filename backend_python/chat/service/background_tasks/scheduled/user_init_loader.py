import asyncio
from backend_python.chat.repository.user_repo import get_all_users
from backend_python.core.db_client import SessionFactory
from backend_python.core.redis_client import redis_client
from data.avatar_data import DEFAULT_AVATAR_URL
from data.cloud_data import R2_BASE_URL

async def preload_all_users() -> None:
    async with SessionFactory() as db:
        user_list = await get_all_users(db)

    pipe = redis_client.pipeline()
    for user_id, username, avatar_key, bio in user_list:
        avatar_url = (
            DEFAULT_AVATAR_URL
            if avatar_key is None
            else f"{R2_BASE_URL}/{avatar_key}"
        )

        pipe.hset(
            f"user:{user_id}",
            mapping={
                "id": str(user_id),
                "username": username,
                "status": "offline",
                "avatar_url": avatar_url,
                "bio": bio or ""
            },
        )
        pipe.sadd("users_all", user_id)

    await pipe.execute()

async def run_preload_all_users():
    print("Users loader worker started")
    await preload_all_users()
    print("Users loader worker finished")