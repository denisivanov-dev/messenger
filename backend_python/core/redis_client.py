import os
import redis.asyncio as redis

def get_redis_url() -> str:
    if url := os.getenv("REDIS_URL"):
        return url

    host = os.getenv("REDIS_HOST", "redis")
    port = os.getenv("REDIS_PORT", "6379")
    db   = os.getenv("REDIS_DB", "0")
    pwd  = os.getenv("REDIS_PASSWORD")

    auth = f":{pwd}@" if pwd else ""
    
    return f"redis://{auth}{host}:{port}/{db}"

print("🔧 REDIS_URL used:", get_redis_url())
redis_client = redis.from_url(get_redis_url(), decode_responses=True)