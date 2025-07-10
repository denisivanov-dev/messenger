import asyncio
from backend_python.chat.service.workers.savers import (
    message_saver_worker
)

async def main():
    await asyncio.gather(
        message_saver_worker.start_redis_listener()
    )

if __name__ == "__main__":
    asyncio.run(main())