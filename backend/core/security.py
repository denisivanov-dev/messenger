from passlib.context import CryptContext
from jose import jwt, JWTError
from datetime import datetime
# from config import SECRET_KEY

pwd_context = CryptContext(schemes=["bcrypt"])

def hash_password(password: str) -> str:
    return pwd_context.hash(password)

def verify_password(plain_password: str, hashed_password: str) -> bool:
    return pwd_context.verify(plain_password, hashed_password)

# def create_access_token(data: dict) -> str:
#     return jwt.encode(data, SECRET_KEY, algorithm="HS256")

# def decode_access_token(token: str) -> dict:
#     try:
#         return jwt.decode(token, SECRET_KEY, algorithms=["HS256"])
#     except JWTError:
#         raise