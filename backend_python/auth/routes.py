from fastapi import APIRouter, Depends, Request, HTTPException
from sqlalchemy.orm import Session
from fastapi.security import OAuth2PasswordRequestForm
from fastapi import BackgroundTasks, Request
from fastapi.responses import JSONResponse
from backend_python.database import get_async_db
from backend_python.auth import schemas, deps
from backend_python.auth.service import auth
from backend_python.auth.service.user import get_user_by_user_id
from backend_python.auth.models import UserSession
from datetime import datetime, timedelta
from sqlalchemy.ext.asyncio import AsyncSession
import asyncio
import time

from backend_python.auth.utils.cookies import set_refresh_cookie, get_refresh_token_from_request
from backend_python.auth.utils.email import send_verification_email, generate_code
from backend_python.auth.utils.redis_client import (
    set_code, set_registration_password,
    get_code, get_registration_password,
    delete_code, delete_registration_password
)
from backend_python.auth.service.token import (
    create_and_store_refresh_token, get_user_session_from_token,
    get_user_session_from_user_id
)
from backend_python.core.security import hash_password, create_access_token

router = APIRouter(prefix="/auth", tags=["auth"])

@router.post("/login")
async def login(user_data: schemas.UserLogin, 
                request: Request,
                db: AsyncSession = Depends(get_async_db)):
    user = await auth.authenticate_user(db, user_data)

    session = await get_user_session_from_user_id(user.id, db)

    if session and session.is_active:
        return JSONResponse(
            status_code=202,
            content={"message": "need_code_verification"}
        )

    refresh_token = await create_and_store_refresh_token(
        user_id=user.id,
        db=db,
        user_agent=request.headers.get("user-agent", "unknown"),
        ip_address=request.client.host
    )
    access_token = create_access_token({"sub": str(user.id)})

    response = JSONResponse(
        status_code=200,
        content={
            "user": schemas.PublicUser.model_validate(user).model_dump(),
            "access_token": access_token
        }
    )

    set_refresh_cookie(response, refresh_token)

    return response


@router.post("/register", response_model=schemas.RegisterResponse, status_code=202)
async def register(user_data: schemas.UserCreate, background_tasks: BackgroundTasks):  
    user = await auth.validate_user(user_data)

    code = generate_code()
    hashed_password = hash_password(user_data.password)

    await set_code(user_data.email, code, "email_verification")
    await set_registration_password(user_data.email, hashed_password)
    
    background_tasks.add_task(send_verification_email, user_data.email, code)

    return schemas.RegisterResponse(
        message="success",
        user=schemas.PublicUser.model_validate(user)
    )


@router.post("/confirm-registration")
async def confirm_code(data: schemas.UserCreateVerify,  
                       request: Request,
                       db: AsyncSession = Depends(get_async_db)
):
    redis_code = await get_code(data.email, "email_verification")
    if not redis_code:
        raise HTTPException(status_code=400, detail="- Просроченный код")

    if data.code != redis_code:
        raise HTTPException(status_code=400, detail="- Недействительный код")
    
    hashed_password = await get_registration_password(data.email)

    user = await auth.create_user(
        username=data.username,
        email=data.email, 
        hashed_password=hashed_password,
        db=db
    )

    refresh_token = await create_and_store_refresh_token(
        user_id=user.id,
        db=db,
        user_agent=request.headers.get("user-agent", "unknown"),
        ip_address=request.client.host
    )
    access_token = create_access_token({"sub": str(user.id)})
    
    await delete_code(data.email, "email_verification")
    await delete_registration_password(data.email)

    response = JSONResponse(
        status_code=200,
        content={"access_token": access_token}
    )

    set_refresh_cookie(response, refresh_token) #вернуть юзер айди еще
    
    return response


@router.get("/auto-login", response_model=schemas.AutoLoginResponse)
async def verify_session(request: Request, db: AsyncSession = Depends(get_async_db)):
    refresh_token = get_refresh_token_from_request()
    if not refresh_token:
        raise HTTPException(status_code=401, detail="Нет сессионого токена")

    session = await get_user_session_from_token(refresh_token, db)
    if not session or not session.is_active:
        raise HTTPException(status_code=401, detail="Недействительная сессия")
    
    user_agent=request.headers.get("user-agent", "unknown")
    ip_address=request.client.host

    if session.user_agent != user_agent or session.ip_address != ip_address:
        raise HTTPException(status_code=401, detail="session_mismatch")
    
    access_token = create_access_token({"sub": str(session.user_id)})
    user = await get_user_by_user_id(db, user_id=session.user_id)
    user = schemas.PublicUser.model_validate(user)
    
    return {
        "message": "success",
        "user": user,
        "accessToken": access_token
    }