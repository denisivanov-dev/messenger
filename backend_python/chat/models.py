from sqlalchemy import (
    Column, Integer, String, Text, Boolean, DateTime,
    ForeignKey, Index, UniqueConstraint, CheckConstraint, 
    JSON, BigInteger, Float
)
import sqlalchemy as sa 
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

    id = Column(String, primary_key=True)
    chat_id = Column(Integer, ForeignKey("chats.id", ondelete="CASCADE"), nullable=False, index=True)
    sender_id = Column(ForeignKey("users.id", ondelete="SET NULL"), nullable=True, index=True)
    target_id = Column(ForeignKey("users.id", ondelete="SET NULL"), nullable=True, index=True)

    kind = Column(String, nullable=False, default="message", server_default="message")
    action = Column(String, nullable=False, default="send", server_default="send")
    chat_type = Column(String, nullable=False, default="private", server_default="private")
    timestamp = Column(BigInteger, nullable=True)

    text = Column(Text, nullable=True)
    attachments = Column(JSON, nullable=True)
    attachments_db = relationship(
        "Attachment",
        back_populates="message",
        cascade="all, delete-orphan"
    )
    reply_to_id = Column(String, ForeignKey("messages.id", ondelete="SET NULL"), nullable=True)

    is_edited = Column(Boolean, default=False, server_default=sa.text("false"))
    deleted = Column(Boolean, default=False, server_default=sa.text("false"))
    is_pinned = Column(Boolean, default=False, server_default=sa.text("false"))

    created_at = Column(DateTime(timezone=True), server_default=func.now())
    edited_at = Column(DateTime(timezone=True), nullable=True)

    raw_envelope = Column(JSON, nullable=True)

    __table_args__ = (
        Index("idx_messages_chat_created", "chat_id", "created_at"),
    )

    chat = relationship("Chat", back_populates="messages")
    sender = relationship("User", foreign_keys=[sender_id])
    target = relationship("User", foreign_keys=[target_id])
    reply_to = relationship("Message", remote_side=[id], post_update=True)

class MessageEdit(Base):
    __tablename__ = "message_edits"

    id = Column(Integer, primary_key=True, autoincrement=True)
    message_id = Column(ForeignKey("messages.id", ondelete="CASCADE"), nullable=False, index=True)
    old_text = Column(Text, nullable=False)
    edited_at = Column(DateTime(timezone=True), server_default=func.now())

    message = relationship("Message", backref="edits")

class Attachment(Base):
    __tablename__ = "attachments"

    id = Column(Integer, primary_key=True, autoincrement=True)
    message_id = Column(ForeignKey("messages.id", ondelete="CASCADE"), nullable=False, index=True)

    key = Column(String, nullable=False)
    type = Column(String, nullable=False)
    size = Column(Integer, nullable=False)
    mime = Column(String, nullable=True)
    original_name = Column(String, nullable=True)

    # 🔹 Метаданные для мультимедиа
    width = Column(Integer, nullable=True)
    height = Column(Integer, nullable=True)
    duration = Column(Float, nullable=True)

    uploaded_at = Column(DateTime(timezone=True), server_default=func.now())

    message = relationship("Message", back_populates="attachments_db")

class FriendLink(Base):
    __tablename__ = "friend_links"

    id = Column(Integer, primary_key=True)
    user1_id = Column(ForeignKey("users.id", ondelete="CASCADE"), nullable=False, index=True)
    user2_id = Column(ForeignKey("users.id", ondelete="CASCADE"), nullable=False, index=True)

    status = Column(Text, nullable=False, index=True)           # 'pending' | 'friends'
    requested_by = Column(Integer, nullable=True, index=True)   # кто отправил (только при pending)

    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), server_default=func.now(), onupdate=func.now())

    __table_args__ = (
        UniqueConstraint("user1_id", "user2_id", name="uc_friend_pair_single"),
        CheckConstraint("user1_id < user2_id", name="ck_friend_pair_order"),
        Index("ix_friend_links_users", "user1_id", "user2_id"),
    )

    user1 = relationship(User, foreign_keys=[user1_id])
    user2 = relationship(User, foreign_keys=[user2_id])