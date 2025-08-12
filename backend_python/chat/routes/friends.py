from fastapi import APIRouter, Depends
from backend_python.chat.events import friends_events
from backend_python.chat.schemas import FriendRequest
from sqlalchemy.ext.asyncio import AsyncSession
from backend_python.chat.service import friend
from backend_python.database import get_async_db

router = APIRouter(prefix="/friends", tags=["friends"])

@router.get("/list/{user_id}")
async def get_friends_overview(user_id: int, db: AsyncSession = Depends(get_async_db)):
    return await friend.get_friend_status_map(db, user_id)

@router.post("/request")
async def send_friend_request(data: FriendRequest):
    return await friends_events.send_friend_request(data.from_id, data.to_id)

@router.post("/cancel")
async def cancel_friend_request(data: FriendRequest):
    return await friends_events.cancel_friend_request(data.from_id, data.to_id)

@router.post("/accept")
async def accept_friend_request(data: FriendRequest):
    return await friends_events.accept_friend_request(data.from_id, data.to_id)

@router.post("/decline")
async def decline_friend_request(data: FriendRequest):
    return await friends_events.decline_friend_request(data.from_id, data.to_id)

@router.post("/delete")
async def decline_friend_request(data: FriendRequest):
    return await friends_events.delete_friend(data.from_id, data.to_id)
    
