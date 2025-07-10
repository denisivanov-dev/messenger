from fastapi import BackgroundTasks, HTTPException, Request
from fastapi.responses import JSONResponse
from sqlalchemy.ext.asyncio import AsyncSession

from backend_python.auth import schemas
from backend_python.auth.models import User
from backend_python.auth.repository import user_repo
from backend_python.auth.repository.session_repo import (
    create_and_store_refresh_token,
    get_user_session_from_user_id
)
from backend_python.auth.repository.user_repo import (
    get_user_by_email,
    get_user_by_username,
)
from backend_python.auth.utils.cookies import (
    set_refresh_cookie,
)
from backend_python.auth.utils.email import generate_code, send_verification_email
from backend_python.auth.utils.redis_client import (
    delete_code,
    delete_registration_password,
    get_code,
    get_registration_password,
    set_code,
    set_registration_password,
)
from backend_python.auth.utils.validators import (
    validate_email,
    validate_empty_data,
    validate_password,
    validate_username,
)
from backend_python.auth.service.session import validate_session, maybe_refresh_session
from backend_python.auth.utils.helpers import validate_or_raise
from backend_python.core.security import create_access_token, hash_password, verify_password


# ---------------------------------------------------------------------------
# Базовые проверки
# ---------------------------------------------------------------------------

async def authenticate_user(db: AsyncSession, login_data: schemas.UserLogin) -> User:
    user = await get_user_by_email(db, login_data.email)
    if not user:
        raise HTTPException(status_code=400, detail={"email": "- Пользователя не существует"})

    if not verify_password(login_data.password, user.hashed_password):
        raise HTTPException(status_code=400, detail={"password": "- Неверный пароль"})

    return user


async def validate_user(user_data: schemas.UserCreate, db: AsyncSession) -> User:
    valid, field, err = validate_empty_data(
        [
            ("email", user_data.email),
            ("username", user_data.username),
            ("password", user_data.password),
            ("confirmPassword", user_data.confirmPassword),
        ]
    )
    validate_or_raise(valid, field, err)

    valid, err, checked_email = validate_email(user_data.email)
    validate_or_raise(valid, "email", err)

    valid, err, checked_username = validate_username(user_data.username)
    validate_or_raise(valid, "username", err)

    valid, err = validate_password(user_data.password, user_data.confirmPassword)
    validate_or_raise(valid, "password", err)

    if await get_user_by_email(db, checked_email):
        raise HTTPException(status_code=400, detail={"email": "- Почта занята"})

    if await get_user_by_username(db, checked_username):
        raise HTTPException(status_code=400, detail={"username": "- Юзернейм занят"})

    return User(username=checked_username, email=checked_email)

# ---------------------------------------------------------------------------
# Сервис-функции
# ---------------------------------------------------------------------------

async def login_user(request: Request, db: AsyncSession, user_data: schemas.UserLogin):
    user = await authenticate_user(db, user_data)

    session = await get_user_session_from_user_id(user.id, db)
    # if session and session.is_active:
    #     return JSONResponse(      !!! допилить подтверждение через устройства как в тг
    #         status_code=202,
    #         content={"message": "need_code_verification"}
    #     )

    refresh_token = await create_and_store_refresh_token(
        user_id=user.id,
        db=db,
        user_agent=request.headers.get("user-agent", "unknown"),
        ip_address=request.client.host,
    )
    access_token = create_access_token({"sub": str(user.id), "username": user.username})

    response = JSONResponse(
        status_code=200,
        content={
            "message": "success",
            "user": schemas.PublicUser.model_validate(user).model_dump(),
            "accessToken": access_token,
        },
    )
    set_refresh_cookie(response, refresh_token)
    return response


async def register_user(
    user_data: schemas.UserCreate, background_tasks: BackgroundTasks, db: AsyncSession
):
    user = await validate_user(user_data, db)

    code = generate_code()
    hashed_password = hash_password(user_data.password)

    await set_code(user_data.email, code, "email_verification")
    await set_registration_password(user_data.email, hashed_password)

    background_tasks.add_task(send_verification_email, user_data.email, code)

    print(1, user.email, user.username)
    return schemas.RegisterResponse(
        message="success", 
        user=schemas.unfinishedPublicUser.model_validate(user)
    )


async def confirm_registration(
    data: schemas.UserCreateVerify, request: Request, db: AsyncSession
):
    redis_code = await get_code(data.email, "email_verification")
    if not redis_code:
        raise HTTPException(status_code=400, detail="- Просроченный код")
    if data.code != redis_code:
        raise HTTPException(status_code=400, detail="- Недействительный код")

    hashed_password = await get_registration_password(data.email)
    user = await user_repo.create_user(
        username=data.username, email=data.email, hashed_password=hashed_password, db=db
    )

    refresh_token = await create_and_store_refresh_token(
        user_id=user.id,
        db=db,
        user_agent=request.headers.get("user-agent", "unknown"),
        ip_address=request.client.host,
    )
    access_token = create_access_token({"sub": str(user.id), "username": user.username})

    await delete_code(data.email, "email_verification")
    await delete_registration_password(data.email)

    response = JSONResponse(
        status_code=200,
        content={
            "message": "success", 
            "accessToken": access_token,
            "user": schemas.PublicUser.model_validate(user).model_dump(),
        },
    )
    set_refresh_cookie(response, refresh_token)
    return response


async def auto_login(request: Request, db: AsyncSession):
    refresh_token, session, user = await validate_session(request, db)

    access_token = create_access_token({"sub": str(user.id), "username": user.username})

    response = JSONResponse(content={
        "message": "success",
        "user": schemas.PublicUser.model_validate(user).model_dump(),
        "accessToken": access_token,
    })

    await maybe_refresh_session(session, refresh_token, response, db)

    return response


async def refresh_token(request: Request, db: AsyncSession):
    refresh_token, session, user = await validate_session(request, db)

    access_token = create_access_token({"sub": str(user.id), "username": user.username})

    response = JSONResponse(content={"access_token": access_token})

    await maybe_refresh_session(session, refresh_token, response, db)

    return response