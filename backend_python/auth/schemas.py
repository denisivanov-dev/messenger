from pydantic import BaseModel, EmailStr
from typing import Optional
from datetime import datetime


# ==== AUTH INPUT ====

class UserLogin(BaseModel):
    email: str
    password: str

class UserCreate(BaseModel):
    email: Optional[str]
    password: Optional[str]
    username: Optional[str]
    confirmPassword: Optional[str]

class UserCreateVerify(BaseModel):
    username: str
    email: str
    code: str


# ==== AUTH OUTPUT ====

class SuccessResponse(BaseModel):
    message: str

class ConfirmRegistrationResponse(BaseModel):
    message: str
    access_token: str


# ==== USER MODELS ====

class UserOut(BaseModel):
    id: int
    username: str
    email: str
    created_at: datetime
    last_login: Optional[datetime] = None
    avatar_url: Optional[str] = None
    bio: Optional[str] = None
    is_active: bool

    model_config = {
        "from_attributes": True
    }

class unfinishedPublicUser(BaseModel):
    username: str
    email: str

    model_config = {
        "from_attributes": True
    }

class PublicUser(BaseModel):
    id: int
    username: str
    email: str

    model_config = {
        "from_attributes": True
    }

class RegisterResponse(BaseModel):
    message: str
    user: unfinishedPublicUser

class AutoLoginResponse(BaseModel):
    message: str
    user: PublicUser
    accessToken: str


# ==== TOKENS ====

class TokenCreate(BaseModel):
    user_id: int
    token: str
    expires_at: datetime
    user_agent: Optional[str] = None
    ip_address: Optional[str] = None

class TokenOut(BaseModel):
    access_token: str
    token_type: str = "bearer"

class TokenInDB(BaseModel):
    id: int
    user_agent: Optional[str]
    ip_address: Optional[str]
    created_at: datetime
    expires_at: datetime
    is_active: bool

    model_config = {
        "from_attributes": True
    }
