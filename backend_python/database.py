from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession
from sqlalchemy.orm import sessionmaker, declarative_base

try:
    from backend_python.config import DB_URL
except ImportError:
    DB_URL = "postgresql+asyncpg://app:app@db:5432/app"

engine = create_async_engine(DB_URL, echo=False)

async_session = sessionmaker(
    bind=engine,
    expire_on_commit=False,
    class_=AsyncSession
)

Base = declarative_base()

async def get_async_db():
    async with async_session() as session:
        yield session
