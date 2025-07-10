from fastapi import APIRouter, Depends, Request, HTTPException
from sqlalchemy.ext.asyncio import AsyncSession

from backend_python.database import get_async_db
from backend_python.chat import schemas
from backend_python.core.security import decode_access_token
from backend_python.chat.repository.chat_repo import get_or_create_private_chat

async def open_or_create_private_chat(
    data: schemas.PrivateChatCreateRequest,
    request: Request,
    db: AsyncSession = Depends(get_async_db),
):
    auth_header = request.headers.get("Authorization")
    if not auth_header or not auth_header.startswith("Bearer "):
        raise HTTPException(status_code=401, detail="Missing or invalid Authorization header")

    token = auth_header[len("Bearer "):]
    payload = decode_access_token(token)
    target_id = data.target_id

    chat, messages = await get_or_create_private_chat(
        db=db,
        user1_id=int(payload["sub"]),
        user2_id=target_id
    )

    return {
        "success": True,
        "chat": chat,
        "messages": messages
    }