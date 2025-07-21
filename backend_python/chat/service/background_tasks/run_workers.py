import asyncio
from backend_python.chat.service.background_tasks.workers import (
    message_saver_worker, message_delete_worker,
    message_editor_worker, message_pin_worker
    
)

async def main():
    await asyncio.gather(
        message_saver_worker.start_save_listener(),
        message_delete_worker.start_delete_listener(),
        message_editor_worker.start_edit_listener(),
        message_pin_worker.start_pin_listener()
    )

if __name__ == "__main__":
    print("Worker started")
    asyncio.run(main())