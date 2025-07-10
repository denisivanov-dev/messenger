import asyncio
from backend_python.chat.service.workers.listeners import (
    chat_ids_loader, message_loader_worker, user_init_worker
)

async def main():
    await asyncio.gather(
        chat_ids_loader.run_load_all_chat_ids(),
        message_loader_worker.run_preload_global_chat_history(),
        user_init_worker.run_preload_all_users()
    )

if __name__ == "__main__":
    asyncio.run(main())