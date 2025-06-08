from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession
from sqlalchemy.orm import sessionmaker, declarative_base
from backend_python.config import DB_URL  

engine = create_async_engine(DB_URL, echo=False)

async_session = sessionmaker(
    bind=engine,
    expire_on_commit=False,
    class_=AsyncSession
)

Base = declarative_base()

async def get_async_db():
    async with async_session() as session:
        try:
            yield session
        finally:
            await session.close()
