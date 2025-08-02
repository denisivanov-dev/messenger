import uuid
import mimetypes
import io
from data.cloud_data import r2_client, R2_BUCKET_NAME

async def upload_to_r2(filename: str, content: bytes) -> tuple[str, int]:
    file_ext = filename.split('.')[-1].lower()
    
    unique_name = f"{uuid.uuid4().hex}.{file_ext}"

    mime_type = mimetypes.types_map.get(f".{file_ext}", "application/octet-stream")

    file_obj = io.BytesIO(content)
    size = len(content)

    r2_client.put_object(
        Bucket=R2_BUCKET_NAME,
        Key=unique_name,
        Body=file_obj,
        ContentType=mime_type,
        ContentLength=size
    )

    return unique_name, size

def generate_presigned_url(key: str, expires_in: int = 3600) -> str:
    return r2_client.generate_presigned_url(
        ClientMethod='get_object',
        Params={'Bucket': R2_BUCKET_NAME, 'Key': key},
        ExpiresIn=expires_in
    )
