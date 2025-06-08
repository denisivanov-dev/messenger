import redis.asyncio as redis
from typing import Optional

redis_client = redis.Redis(
    host="localhost",
    port=6379,     
    db=0,          
    decode_responses=True
)

DEFAULT_EXPIRE_SECONDS = 300

async def set_code(email: str, code: str, category: str, ttl: int = DEFAULT_EXPIRE_SECONDS):
    key = f"{category}:{email}"

    await redis_client.set(key, code, ex=ttl)

async def get_code(email: str, category: str) -> Optional[str]:
    key = f"{category}:{email}"

    return await redis_client.get(key)

async def delete_code(email: str, category: str):
    key = f"{category}:{email}"

    await redis_client.delete(key)

async def set_registration_password(email: str, password: str, expire: int = DEFAULT_EXPIRE_SECONDS) -> None:
    key = f"registration_password:{email}"
    
    await redis_client.set(key, password, ex=expire)

async def get_registration_password(email: str) -> str | None:
    key = f"registration_password:{email}"
    
    return await redis_client.get(key)

async def delete_registration_password(email: str) -> None:
    key = f"registration_password:{email}"
    
    await redis_client.delete(key)
