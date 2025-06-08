from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession
from sqlalchemy.orm import sessionmaker
from fastapi import HTTPException
from backend_python.core.security import hash_password, verify_password
from backend_python.auth.models import User
from backend_python.auth.schemas import UserCreate, UserLogin
import asyncio
from fastapi import Request, Depends

from backend_python.auth.utils.validators import (
    validate_empty_data, validate_username, validate_email, validate_password,
)

from backend_python.auth.utils.helpers import validate_or_raise

from backend_python.auth.service.user import (
    get_user_by_username, get_user_by_email,
)

async def authenticate_user(db: AsyncSession, login_data: UserLogin) -> User | None:
    user = await get_user_by_email(db, login_data.email)

    if not user:
        raise HTTPException(status_code=400, detail={"email": "- Пользователя не существует"})
    
    if not verify_password(login_data.password, user.hashed_password):
        raise HTTPException(status_code=400, detail={"password": "- Неверный пароль"})
    
    return user

async def validate_user(user_data: UserCreate) -> User:
    # valid, field, error_msg = validate_empty_data([
    # ("email", user_data.email),
    # ("username", user_data.username),
    # ("password", user_data.password),
    # ("confirmPassword", user_data.confirmPassword)
    # ])
    # validate_or_raise(valid, field, error_msg)

    valid, error_msg, checked_email = validate_email(user_data.email)
    # validate_or_raise(valid, "email", error_msg)

    valid, error_msg, checked_username = validate_username(user_data.username)
    # validate_or_raise(valid, "username", error_msg)

    # valid, error_msg = validate_password(user_data.password, user_data.confirmPassword)
    # validate_or_raise(valid, "password", error_msg)
    
    # existing_email = await get_user_by_email(db, user_data.email)
    # if existing_email:
    #     raise HTTPException(status_code=400, detail={ "email": "Почта занята" })

    # existing_user = await get_user_by_username(db, user_data.username)
    # if existing_user:
    #     raise HTTPException(status_code=400, detail={ "username": "Юзернейм занят" })

    user = User(
        username = checked_username,
        email = checked_email
    )

    return user

async def create_user(username: str, email: str, 
                      hashed_password: str, db: AsyncSession) -> User:
    user = User(
        username=username,
        email=email,
        hashed_password=hashed_password,
        is_active=True 
    )

    db.add(user)
    await db.commit()
    await db.refresh(user)

    return user
