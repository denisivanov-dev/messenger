from pydantic import BaseModel, EmailStr
from typing import Optional
from datetime import datetime

class UserCreate(BaseModel):
    email: Optional[str]
    password: Optional[str]
    username: Optional[str]
    confirmPassword: Optional[str]

class UserLogin(BaseModel):
    username: str
    password: str

class UserOut(BaseModel):
    id: int
    username: str
    email: str
    created_at: datetime
    last_login: datetime | None = None
    avatar_url: str | None = None
    bio: str | None = None
    is_active: bool

    class Config:
        from_attributes = True

class TokenCreate(BaseModel):
    user_id: int
    token: str
    expires_at: datetime
    user_agent: str | None = None
    ip_address: str | None = None

class TokenOut(BaseModel):
    access_token: str
    token_type: str = "bearer"

class TokenInDB(BaseModel):
    id: int
    user_agent: str | None
    ip_address: str | None
    created_at: datetime
    expires_at: datetime
    is_active: bool

    class Config:
        from_attributes = True
