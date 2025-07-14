import asyncio
from backend_python.chat.service.background_tasks.scheduled import (
    chat_ids_loader, message_loader, user_init_loader
)

async def main():
    await asyncio.gather(
        chat_ids_loader.run_load_all_chat_ids(),
        message_loader.run_preload_global_chat_history(),
        user_init_loader.run_preload_all_users()
    )

if __name__ == "__main__":
    asyncio.run(main())