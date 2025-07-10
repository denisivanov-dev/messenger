from fastapi import APIRouter, Depends, Request, BackgroundTasks
from sqlalchemy.ext.asyncio import AsyncSession

from backend_python.database import get_async_db
from backend_python.auth import schemas
from backend_python.auth.service import auth

router = APIRouter(prefix="/auth", tags=["auth"])


@router.post("/login")
async def login(
    user_data: schemas.UserLogin,
    request: Request,
    db: AsyncSession = Depends(get_async_db)
):
    return await auth.login_user(request, db, user_data)


@router.post("/register", response_model=schemas.RegisterResponse, status_code=202)
async def register(
    user_data: schemas.UserCreate,
    background_tasks: BackgroundTasks,
    db: AsyncSession = Depends(get_async_db)
):
    return await auth.register_user(user_data, background_tasks, db)


@router.post("/confirm-registration")
async def confirm_registration(
    data: schemas.UserCreateVerify,
    request: Request,
    db: AsyncSession = Depends(get_async_db)
):
    return await auth.confirm_registration(data, request, db)


@router.get("/auto-login")
async def auto_login(
    request: Request,
    db: AsyncSession = Depends(get_async_db)
):
    return await auth.auto_login(request, db)


@router.post("/refresh")
async def refresh_token(
    request: Request,
    db: AsyncSession = Depends(get_async_db)
):
    return await auth.refresh_token(request, db)