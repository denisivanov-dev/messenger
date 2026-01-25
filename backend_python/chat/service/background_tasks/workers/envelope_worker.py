import asyncio
import json
from backend_python.core.redis_client import redis_client
from backend_python.core.db_client import SessionFactory
from backend_python.chat.repository.chat_repo import get_private_chat_keys
from backend_python.chat.service.background_tasks.envelope.parser import parse_envelope
from backend_python.chat.service.background_tasks.envelope.validator import validate_envelope
from backend_python.chat.service.background_tasks.envelope.dispatcher import dispatch_envelope

async def listen_queue(queue: str):
    print(f"[init] listening {queue}")
    while True:
        try:
            _, data = await redis_client.blpop(queue)
            raw = json.loads(data)
            env = parse_envelope(raw)

            if not validate_envelope(env, queue):
                continue

            async with SessionFactory() as db:
                await dispatch_envelope(env, db)

        except Exception as e:
            print(f"[ERR] queue={queue}: {e}")

async def start_envelope_worker():
    async with SessionFactory() as db:
        chat_keys = await get_private_chat_keys(db)

    queues = ["to_save:global", *[f"to_save:private:{k}" for k in chat_keys]]

    print(f"[init] envelope worker listening {len(queues)} queues")

    await asyncio.gather(*[listen_queue(q) for q in queues])
