from sqlalchemy.ext.asyncio import AsyncSession
import redis.asyncio as redis
import backend_python.chat.repository.friends_repo as friends_repo
from fastapi import status
from fastapi.responses import JSONResponse

redis_client = redis.Redis(host="localhost", port=6379, decode_responses=True)

async def get_friend_status_map(db: AsyncSession, user_id: int):
    friends_key = f"user:{user_id}:friends"
    incoming_key = f"user:{user_id}:pending"
    outgoing_key = f"user:{user_id}:outgoing"

    keys_exist = await redis_client.exists(friends_key, incoming_key, outgoing_key)

    if keys_exist < 3:
        result = await friends_repo.get_friends_from_db(db, user_id)

        pipe = redis_client.pipeline(transaction=True)
        pipe.delete(friends_key, incoming_key, outgoing_key)

        if result["friends"]:
            pipe.sadd(friends_key, *(str(u["id"]) for u in result["friends"]))
        if result["incoming"]:
            pipe.sadd(incoming_key, *(str(u["id"]) for u in result["incoming"]))
        if result["outgoing"]:
            pipe.sadd(outgoing_key, *(str(u["id"]) for u in result["outgoing"]))

        await pipe.execute()

    friends_ids = await redis_client.smembers(friends_key)
    incoming_ids = await redis_client.smembers(incoming_key)
    outgoing_ids = await redis_client.smembers(outgoing_key)

    status_map = {}
    for fid in friends_ids:
        status_map[fid] = "friends"
    for fid in incoming_ids:
        status_map[fid] = "incoming"
    for fid in outgoing_ids:
        status_map[fid] = "outgoing"

    return JSONResponse(content=status_map, status_code=status.HTTP_200_OK)