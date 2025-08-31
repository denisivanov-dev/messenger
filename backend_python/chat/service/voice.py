import redis.asyncio as redis

rdb = redis.Redis(host="localhost", port=6379, decode_responses=True)

async def get_call_room_from_redis(room_id: str) -> dict:
    key = f"callroom:{room_id}"
    return await rdb.hgetall(key)