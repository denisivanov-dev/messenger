import os
from sqlalchemy.ext.asyncio import create_async_engine, async_sessionmaker, AsyncSession

def get_db_url() -> str:
    if url := os.getenv("DATABASE_URL"):
        return url

    user = os.getenv("DB_USER", "app")
    password = os.getenv("DB_PASSWORD", "app")
    host = os.getenv("DB_HOST", "db")
    port = os.getenv("DB_PORT", "5432")
    name = os.getenv("DB_NAME", "app")

    return f"postgresql+asyncpg://{user}:{password}@{host}:{port}/{name}"

engine = create_async_engine(
    get_db_url(),
    echo=False,
    pool_pre_ping=True
)

SessionFactory = async_sessionmaker(
    bind=engine,
    expire_on_commit=False,
    class_=AsyncSession,
)