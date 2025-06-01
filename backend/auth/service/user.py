from sqlalchemy.orm import Session
from backend.auth.models import User

def get_user_by_username(db: Session, username: str) -> User | None:
    query = db.query(User)
    query = query.filter(User.username == username)
    user = query.first()

    return user

def get_user_by_email(db: Session, email: str) -> User | None:
    query = db.query(User)
    query = query.filter(User.email == email)
    user = query.first()

    return user