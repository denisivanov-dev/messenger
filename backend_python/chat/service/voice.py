from backend_python.core.redis_client import redis_client

async def get_call_room_from_redis(room_id: str) -> dict:
    key = f"callroom:{room_id}"
    return await redis_client.hgetall(key)