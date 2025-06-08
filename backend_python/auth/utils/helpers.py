from fastapi import HTTPException
import asyncio
from backend_python.auth.utils.email import send_verification_email

def validate_or_raise(valid: bool, field: str, error_msg: str):
    if not valid:
        raise HTTPException(status_code=400, detail={field: error_msg})