from datetime import datetime, timedelta
from sqlalchemy.ext.asyncio import AsyncSession
from backend_python.auth.models import UserSession
from backend_python.core.security import generate_refresh_token, hash_token
from sqlalchemy.future import select 
from typing import Optional
from sqlalchemy import update

REFRESH_TOKEN_EXPIRES_DAYS = 7


async def create_and_store_refresh_token(user_id: int, db: AsyncSession, 
                                         user_agent: str, ip_address: str) -> str:
    refresh_token = generate_refresh_token()
    refresh_token_hash = hash_token(refresh_token)

    print("Token →", refresh_token)
    print("Hash  →", refresh_token_hash)

    session = UserSession(
        user_id=user_id,
        refresh_token_hash=refresh_token_hash,
        expires_at=datetime.utcnow() + timedelta(days=REFRESH_TOKEN_EXPIRES_DAYS),
        user_agent=user_agent,
        ip_address=ip_address
    )

    db.add(session)
    await db.commit()

    return refresh_token 


async def prolong_user_session(token: str, db):
    new_exp = datetime.utcnow() + timedelta(days=7)
    token_hash = hash_token(token)
    
    await db.execute(
        update(UserSession)
        .where(UserSession.refresh_token_hash == token_hash)
        .values(expires_at=new_exp)
    )
    await db.commit()


async def get_user_session_from_token(token: str, db: AsyncSession) -> Optional[UserSession]:
    token_hash = hash_token(token)

    result = await db.execute(
        select(UserSession).where(UserSession.refresh_token_hash == token_hash)
    )

    return result.scalar_one_or_none()


async def get_user_session_from_user_id(user_id: str, db: AsyncSession) -> Optional[UserSession]:
    result = await db.execute(
        select(UserSession).where(UserSession.user_id == user_id)
    )

    return result.scalar_one_or_none()