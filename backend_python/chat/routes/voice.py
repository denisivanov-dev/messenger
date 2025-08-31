from fastapi import APIRouter, HTTPException
import redis.asyncio as redis
from backend_python.chat.service import voice

router = APIRouter(prefix="/voice", tags=["voice"])

@router.get("/{room_id}")
async def get_call_room_status(room_id: str):
    return await voice.get_call_room_from_redis(room_id)