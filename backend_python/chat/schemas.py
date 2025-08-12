from pydantic import BaseModel, EmailStr
from typing import Optional
from datetime import datetime

class PrivateChatCreateRequest(BaseModel):
    target_id: int

class FriendRequest(BaseModel):
    from_id: int
    to_id: int