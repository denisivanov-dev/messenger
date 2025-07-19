from sqlalchemy import (
    Column, Integer, String, Text, Boolean, DateTime,
    ForeignKey, Index, UniqueConstraint
)
from sqlalchemy.orm import relationship
from sqlalchemy.sql import func
from backend_python.database import Base
from backend_python.auth.models import User

class Chat(Base):
    __tablename__ = "chats"

    id = Column(Integer, primary_key=True)
    type  = Column(Text, nullable=False, index=True)   # 'global' | 'private' | 'group'
    title = Column(Text, nullable=True)                # null для приватных

    chat_key = Column(Text, unique=True, nullable=True)
    
    created_at = Column(DateTime(timezone=True), server_default=func.now())

    messages     = relationship("Message", back_populates="chat", cascade="all, delete-orphan")
    participants = relationship("ChatParticipant", back_populates="chat", cascade="all, delete-orphan")


class ChatParticipant(Base):
    __tablename__ = "chat_participants"

    id       = Column(Integer, primary_key=True)
    chat_id  = Column(ForeignKey("chats.id",  ondelete="CASCADE"), nullable=False, index=True)
    user_id  = Column(ForeignKey("users.id",  ondelete="CASCADE"), nullable=False, index=True)
    joined_at = Column(DateTime(timezone=True), server_default=func.now())

    __table_args__ = (
        UniqueConstraint("chat_id", "user_id", name="uc_chat_user"),
    )

    chat = relationship("Chat", back_populates="participants")
    user = relationship("User")                   


class Message(Base):
    __tablename__ = "messages"

    id        = Column(String, primary_key=True)
    chat_id   = Column(ForeignKey("chats.id",  ondelete="CASCADE"), nullable=False, index=True)
    sender_id = Column(ForeignKey("users.id",  ondelete="SET NULL"), nullable=True, index=True)

    content    = Column(Text, nullable=False)
    is_edited  = Column(Boolean, default=False)
    deleted    = Column(Boolean, default=False)
    is_pinned  = Column(Boolean, default=False)

    reply_to_id = Column(String, ForeignKey("messages.id", ondelete="SET NULL"), nullable=True)

    created_at = Column(DateTime(timezone=True), server_default=func.now())
    edited_at  = Column(DateTime(timezone=True), nullable=True)

    __table_args__ = (
        Index("idx_messages_chat_created", "chat_id", "created_at"),
    )

    chat   = relationship("Chat", back_populates="messages")
    sender = relationship("User")
    reply_to = relationship("Message", remote_side=[id], post_update=True)

class MessageEdit(Base):
    __tablename__ = "message_edits"

    id = Column(Integer, primary_key=True, autoincrement=True)
    message_id = Column(ForeignKey("messages.id", ondelete="CASCADE"), nullable=False, index=True)
    old_text = Column(Text, nullable=False)
    edited_at = Column(DateTime(timezone=True), server_default=func.now())

    message = relationship("Message", backref="edits")