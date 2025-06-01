from sqlalchemy.orm import Session
from fastapi import HTTPException, status
from backend.core.security import hash_password, verify_password
from backend.auth.models import User
from backend.auth.schemas import UserCreate, UserLogin
from datetime import datetime, timedelta

# def save_tvaliden(db: Session, tvaliden_data: TvalidenCreate) -> TvalidenSession:
#     tvaliden = TvalidenSession(**tvaliden_data.dict())
#     db.add(tvaliden)
#     db.commit()
#     db.refresh(tvaliden)
#     return tvaliden


# def generate_access_tvaliden(user: User, expire_minutes: int = 60 * 24) -> str:
#     expire_time = datetime.utcnow() + timedelta(minutes=expire_minutes)
#     tvaliden = create_access_tvaliden(data={"sub": user.username, "exp": expire_time})
#     return tvaliden, expire_time