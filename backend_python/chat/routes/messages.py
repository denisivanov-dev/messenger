from fastapi import APIRouter, UploadFile, File, Depends, Request
from fastapi.responses import JSONResponse
from sqlalchemy.ext.asyncio import AsyncSession

from backend_python.database import get_async_db
from backend_python.chat.utils import cloud_utils

from fastapi import HTTPException
from botocore.exceptions import ClientError

router = APIRouter(prefix="/messages", tags=["messages"])

@router.post("/upload-image")
async def upload_file_route(
    request: Request,
    file: UploadFile = File(...),
    db: AsyncSession = Depends(get_async_db),
):
    contents = await file.read()

    key, filesize = await cloud_utils.upload_to_r2(file.filename, contents)

    return JSONResponse({
        "key": key,
        "filesize": filesize,
        "original_name": file.filename
    })

@router.get("/{key}")
async def get_signed_image_url(key: str):
    try:
        url = cloud_utils.generate_presigned_url(key)
        return {"url": url}
    except ClientError as e:
        raise HTTPException(status_code=404, detail="File not found")

