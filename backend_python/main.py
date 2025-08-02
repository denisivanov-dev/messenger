from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from contextlib import asynccontextmanager

from backend_python.auth.routes import router as auth_router
from backend_python.chat.routes.users import router as users_router
from backend_python.chat.routes.chats import router as chat_router
from backend_python.chat.routes.messages import router as messages_router
from backend_python.auth.utils.redis_client import redis_client


@asynccontextmanager
async def lifespan(app: FastAPI):
    try:
        await redis_client.ping()
        print("Redis connected")
    except Exception as e:
        print(f"Redis connection failed: {e}")
    yield 
    await redis_client.close()
    print("Redis disconnected")

app = FastAPI(title="Messenger API", lifespan=lifespan)

app.add_middleware(
    CORSMiddleware,
    allow_origins=[
        "http://localhost:5173",
        "http://127.0.0.1:5173"
        # "https://k3ch20c2-5173.euw.devtunnels.ms"
    ],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(auth_router, prefix="/api")
app.include_router(users_router, prefix="/api")
app.include_router(chat_router, prefix="/api")
app.include_router(messages_router, prefix="/api")
