from pydantic import BaseModel, EmailStr
from typing import Optional
from datetime import datetime

class PrivateChatCreateRequest(BaseModel):
    target_id: int