from fastapi import APIRouter
import redis.asyncio as redis

from backend_python.chat.repository import user_repo

router = APIRouter(prefix="/users", tags=["users"])


@router.get("/all")
async def get_all_users():
    return await user_repo.get_all_users_from_redis()