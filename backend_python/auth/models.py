from sqlalchemy import Column, Integer, String, Boolean, DateTime, ForeignKey
from sqlalchemy.sql import func
from sqlalchemy.orm import relationship
from backend_python.database import Base

class User(Base):
    __tablename__ = "users"

    id = Column(Integer, primary_key=True, index=True)
    username = Column(String, unique=True, index=True, nullable=False)
    email = Column(String, unique=True, index=True, nullable=False)
    hashed_password = Column(String, nullable=False)

    created_at = Column(DateTime(timezone=True), server_default=func.now())
    last_login = Column(DateTime(timezone=True), nullable=True)
    is_active = Column(Boolean, default=False)

    avatar_url = Column(String, nullable=True)
    bio = Column(String, nullable=True)

    sessions = relationship("UserSession", back_populates="user", cascade="all, delete-orphan")

class UserSession(Base):
    __tablename__ = "user_sessions"

    id = Column(Integer, primary_key=True)
    user_id = Column(ForeignKey("users.id", ondelete="CASCADE"), nullable=False)
    refresh_token_hash = Column(String, unique=True, nullable=False, index=True)
    
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    expires_at = Column(DateTime(timezone=True), nullable=False)
    is_active = Column(Boolean, default=True, nullable=False)

    user_agent = Column(String, nullable=True)
    ip_address = Column(String, nullable=True)

    user = relationship("User", back_populates="sessions")