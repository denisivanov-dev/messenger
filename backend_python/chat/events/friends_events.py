import redis.asyncio as redis
import json

redis_client = redis.Redis(host="localhost", port=6379, decode_responses=True)

async def queue_friend_event(event: dict):
    await redis_client.rpush("to_save:friend_events", json.dumps(event))

async def send_friend_request(from_id: int, to_id: int):
    if from_id == to_id:
        return {"error": "Нельзя добавить себя"}

    # Для получателя — входящая заявка
    await redis_client.sadd(f"user:{to_id}:incoming", from_id)
    # Для отправителя — исходящая заявка
    await redis_client.sadd(f"user:{from_id}:outgoing", to_id)

    payload = {
        "type": "friend_request_sent",
        "from_id": str(from_id),
        "to_id": str(to_id)
    }
    await redis_client.publish("friend_requests", json.dumps(payload))
    await queue_friend_event(payload)

    return {"status": "ok"}

async def cancel_friend_request(from_id: int, to_id: int):
    await redis_client.srem(f"user:{to_id}:incoming", from_id)
    await redis_client.srem(f"user:{from_id}:outgoing", to_id)

    payload = {
        "type": "friend_request_canceled",
        "from_id": str(from_id),
        "to_id": str(to_id)
    }
    await redis_client.publish("friend_requests", json.dumps(payload))
    await queue_friend_event(payload)

    return {"status": "ok"}

async def accept_friend_request(from_id: int, to_id: int):
    if from_id == to_id:
        return {"error": "Неверный запрос"}

    # Чистим все статусы заявок
    await redis_client.srem(f"user:{from_id}:incoming", to_id)
    await redis_client.srem(f"user:{from_id}:outgoing", to_id)
    await redis_client.srem(f"user:{to_id}:incoming", from_id)
    await redis_client.srem(f"user:{to_id}:outgoing", from_id)

    # Добавляем в друзья
    await redis_client.sadd(f"user:{from_id}:friends", to_id)
    await redis_client.sadd(f"user:{to_id}:friends", from_id)

    payload = {
        "type": "friend_request_accepted",
        "from_id": str(from_id),
        "to_id": str(to_id)
    }
    await redis_client.publish("friend_requests", json.dumps(payload))
    await queue_friend_event(payload)

    return {"status": "ok"}

async def decline_friend_request(from_id: int, to_id: int):
    await redis_client.srem(f"user:{from_id}:incoming", to_id)
    await redis_client.srem(f"user:{to_id}:outgoing", from_id)

    payload = {
        "type": "friend_request_declined",
        "from_id": str(from_id),
        "to_id": str(to_id)
    }
    await redis_client.publish("friend_requests", json.dumps(payload))
    await queue_friend_event(payload)

    return {"status": "ok"}

async def delete_friend(user1_id: int, user2_id: int):
    await redis_client.srem(f"user:{user1_id}:friends", user2_id)
    await redis_client.srem(f"user:{user2_id}:friends", user1_id)

    payload = {
        "type": "friend_removed",
        "from_id": str(user1_id),
        "to_id": str(user2_id)
    }
    await redis_client.publish("friend_requests", json.dumps(payload))
    await queue_friend_event(payload)

    return {"status": "ok"}