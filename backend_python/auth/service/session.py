from datetime import datetime, timezone, timedelta
from typing import Tuple

from fastapi import HTTPException, Request
from fastapi.responses import Response
from sqlalchemy.ext.asyncio import AsyncSession

from backend_python.auth import schemas
from backend_python.auth.models import UserSession
from backend_python.auth.repository.session_repo import (
    get_user_session_from_token,
    prolong_user_session,
)
from backend_python.auth.repository.user_repo import get_user_by_user_id
from backend_python.auth.utils.cookies import (
    get_refresh_token_from_request,
    set_refresh_cookie,
)

# ---------------------------------------------------------------------------
# Проверка refresh-токена и возврат информации о сессии и пользователе
# ---------------------------------------------------------------------------

async def validate_session(
    request: Request,
    db: AsyncSession,
) -> Tuple[str, UserSession, schemas.PublicUser]:
    refresh_token = get_refresh_token_from_request(request)
    if not refresh_token:
        raise HTTPException(status_code=401, detail="Нет сессионного токена")

    session = await get_user_session_from_token(refresh_token, db)
    if not session or not session.is_active:
        raise HTTPException(status_code=401, detail="Недействительная сессия")

    ua, ip = request.headers.get("user-agent", "unknown"), request.client.host
    if session.user_agent != ua or session.ip_address != ip:
        raise HTTPException(status_code=401, detail="session_mismatch")

    user = await get_user_by_user_id(db, session.user_id)
    return refresh_token, session, user


# ---------------------------------------------------------------------------
# Продление сессии и установка нового refresh-cookie, если «жить» ≤ threshold
# ---------------------------------------------------------------------------

async def maybe_refresh_session(
    session: UserSession,
    refresh_token: str,
    response: Response,
    db: AsyncSession,
    threshold: timedelta = timedelta(days=3),
) -> None:
    now = datetime.now(timezone.utc)
    if session.expires_at - now <= threshold:
        await prolong_user_session(refresh_token, db)
        set_refresh_cookie(response, refresh_token)
