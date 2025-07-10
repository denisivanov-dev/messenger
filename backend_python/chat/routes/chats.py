from fastapi import APIRouter, Depends, Request
from sqlalchemy.ext.asyncio import AsyncSession

from backend_python.database import get_async_db
from backend_python.chat.schemas import PrivateChatCreateRequest
from backend_python.chat.service import chat

router = APIRouter(prefix="/chats", tags=["chats"])

@router.post("/private-chat")
async def open_or_create_private_chat_route(
    data: PrivateChatCreateRequest,
    request: Request,
    db: AsyncSession = Depends(get_async_db),
):
    return await chat.open_or_create_private_chat(
        data=data,
        request=request,
        db=db
    )
