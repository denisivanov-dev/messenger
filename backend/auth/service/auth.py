from sqlalchemy.orm import Session
from fastapi import HTTPException
from backend.core.security import hash_password, verify_password
from backend.auth.models import User
from backend.auth.schemas import UserCreate, UserLogin

from backend.auth.utils import (
    validate_empty_data,
    validate_username, validate_email, validate_password
)

from backend.auth.service.user import get_user_by_username, get_user_by_email

def create_user(db: Session, user_data: UserCreate) -> User:
    valid, field, error_msg = validate_empty_data([
    ("email", user_data.email),
    ("username", user_data.username),
    ("password", user_data.password),
    ("confirmPassword", user_data.confirmPassword)
    ])
    if not valid:
        raise HTTPException(status_code=400, detail={ field: error_msg })

    valid, error_msg = validate_email(user_data.email)
    if not valid:
        raise HTTPException(status_code=400, detail={ "email": error_msg })

    valid, error_msg = validate_username(user_data.username)
    if not valid:    
        raise HTTPException(status_code=400, detail={ "username": error_msg })

    valid, error_msg = validate_password(user_data.password, user_data.confirmPassword)
    if not valid: 
        raise HTTPException(status_code=400, detail={ "password": error_msg })
    
    existing_email = get_user_by_email(db, user_data.email)
    if existing_email:
        raise HTTPException(status_code=400, detail={ "email": "Почта занята" })

    existing_user = get_user_by_username(db, user_data.username)
    if existing_user:
        raise HTTPException(status_code=400, detail={ "username": "Юзернейм занят" })

    user = User(
        username=user_data.username,
        email=user_data.email,
        hashed_password=hash_password(user_data.password)
    )

    db.add(user)
    db.commit()
    db.refresh(user)

    return user

def authenticate_user(db: Session, login_data: UserLogin) -> User | None:
    user = get_user_by_username(db, login_data.username)

    if not user:
        raise HTTPException(status_code=400, detail="Пользователя не существует")
    
    if not verify_password(login_data.password, user.hashed_password):
        raise HTTPException(status_code=400, detail="Неверный пароль")
    
    return user