from backend_python.chat.repository import message_repo

async def handle(env, db):
    if env["chat_type"] == "global":
        await message_repo.save_message_to_global_chat(
            db=db,
            message_id=env["message_id"],
            user_id=env["sender_id"],
            text=env["text"],
            timestamp_ms=env["timestamp"],
            attachments=env["attachments"],
            reply_to_id=env["reply_to_id"],
            raw_envelope=env["raw"],
        )
    else:
        await message_repo.save_private_message(
            db=db,
            chat_key=env["chat_id"],
            message_id=env["message_id"],
            sender_id=env["sender_id"],
            text=env["text"],
            timestamp_ms=env["timestamp"],
            attachments=env["attachments"],
            reply_to_id=env["reply_to_id"],
            raw_envelope=env["raw"],
        )
