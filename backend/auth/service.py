from sqlalchemy.orm import Session
from fastapi import HTTPException, status
from backend.core.security import hash_password, verify_password
from backend.auth.models import User, TokenSession
from backend.auth.schemas import UserCreate, UserLogin, TokenCreate
from datetime import datetime, timedelta

from backend.auth.utils import (
    validate_username, validate_email, validate_password
)

def get_user_by_username(db: Session, username: str) -> User | None:
    query = db.query(User)
    query = query.filter(User.username == username)
    user = query.first()

    return user

def create_user(db: Session, user_data: UserCreate) -> User:
    existing = get_user_by_username(db, user_data.username)
    if existing:
        raise HTTPException(status_code=400, detail="Юзернейм занят")

    ok, error_msg = validate_username(user_data.username)
    if not ok:    
        raise HTTPException(status_code=400, detail=error_msg)
    
    ok, error_msg = validate_email(user_data.email)
    if not ok:
        raise HTTPException(status_code=400, detail=error_msg)
    
    ok, error_msg = validate_password(user_data.password)
    if not ok: 
        raise HTTPException(status_code=400, detail=error_msg)

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

# ------------------------------------------------------------

# def save_token(db: Session, token_data: TokenCreate) -> TokenSession:
#     token = TokenSession(**token_data.dict())
#     db.add(token)
#     db.commit()
#     db.refresh(token)
#     return token


# def generate_access_token(user: User, expire_minutes: int = 60 * 24) -> str:
#     expire_time = datetime.utcnow() + timedelta(minutes=expire_minutes)
#     token = create_access_token(data={"sub": user.username, "exp": expire_time})
#     return token, expire_time
