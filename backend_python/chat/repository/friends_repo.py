from typing import Tuple, Dict
from sqlalchemy import select, insert, update, delete, or_, and_
from sqlalchemy.ext.asyncio import AsyncSession
from backend_python.chat.models import FriendLink, User


def normalize_pair(user_id_1: int, user_id_2: int) -> Tuple[int, int]:
    return (user_id_1, user_id_2) if user_id_1 < user_id_2 else (user_id_2, user_id_1)


async def on_request_sent(db: AsyncSession, from_user_id: int, to_user_id: int) -> None:
    user1_id, user2_id = normalize_pair(from_user_id, to_user_id)

    existing = await db.execute(
        select(FriendLink).where(
            FriendLink.user1_id == user1_id,
            FriendLink.user2_id == user2_id,
        )
    )
    row = existing.scalar_one_or_none()

    if row is None:
        await db.execute(
            insert(FriendLink).values(
                user1_id=user1_id,
                user2_id=user2_id,
                status="pending",
                requested_by=from_user_id,
            )
        )
    else:
        await db.execute(
            update(FriendLink)
            .where(FriendLink.id == row.id)
            .values(status="pending", requested_by=from_user_id)
        )


async def on_request_canceled(db: AsyncSession, from_user_id: int, to_user_id: int) -> None:
    user1_id, user2_id = normalize_pair(from_user_id, to_user_id)
    await db.execute(
        delete(FriendLink).where(
            FriendLink.user1_id == user1_id,
            FriendLink.user2_id == user2_id,
            FriendLink.status == "pending",
        )
    )


async def on_request_declined(db: AsyncSession, from_user_id: int, to_user_id: int) -> None:
    user1_id, user2_id = normalize_pair(from_user_id, to_user_id)
    await db.execute(
        delete(FriendLink).where(
            FriendLink.user1_id == user1_id,
            FriendLink.user2_id == user2_id,
            FriendLink.status == "pending",
        )
    )


async def on_request_accepted(db: AsyncSession, from_user_id: int, to_user_id: int) -> None:
    user1_id, user2_id = normalize_pair(from_user_id, to_user_id)

    existing = await db.execute(
        select(FriendLink).where(
            FriendLink.user1_id == user1_id,
            FriendLink.user2_id == user2_id,
        )
    )
    row = existing.scalar_one_or_none()

    if row is None:
        await db.execute(
            insert(FriendLink).values(
                user1_id=user1_id,
                user2_id=user2_id,
                status="friends",
                requested_by=None,
            )
        )
    else:
        await db.execute(
            update(FriendLink)
            .where(FriendLink.id == row.id)
            .values(status="friends", requested_by=None)
        )


async def on_friend_removed(db: AsyncSession, user_id_1: int, user_id_2: int) -> None:
    user1_id, user2_id = normalize_pair(user_id_1, user_id_2)
    await db.execute(
        delete(FriendLink).where(
            FriendLink.user1_id == user1_id,
            FriendLink.user2_id == user2_id,
        )
    )

async def get_friends_from_db(db: AsyncSession, user_id: int):
    link_to_user = or_(
        and_(FriendLink.user1_id == user_id, FriendLink.user2_id == User.id),
        and_(FriendLink.user2_id == user_id, FriendLink.user1_id == User.id),
    )

    async def fetch_list(query):
        result = await db.execute(query)
        rows = result.all()
        return [{"id": r.id, "username": r.username} for r in rows]

    q_friends = (
        select(User.id, User.username)
        .join(FriendLink, link_to_user)
        .where(FriendLink.status == "friends")
    )
    q_incoming = (
        select(User.id, User.username)
        .join(FriendLink, link_to_user)
        .where(FriendLink.status == "pending", FriendLink.requested_by != user_id)
    )
    q_outgoing = (
        select(User.id, User.username)
        .join(FriendLink, link_to_user)
        .where(FriendLink.status == "pending", FriendLink.requested_by == user_id)
    )

    friends  = await fetch_list(q_friends)
    incoming = await fetch_list(q_incoming)
    outgoing = await fetch_list(q_outgoing)

    return {
        "friends": friends,
        "incoming": incoming,
        "outgoing": outgoing,
    }