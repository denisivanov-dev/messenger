from fastapi import Response, Request

def set_refresh_cookie(response: Response, token: str):
    response.set_cookie(
        key="refresh_token",
        value=token,
        httponly=True,
        secure=False,  # не забудь поставить True на проде
        samesite="Lax",
        # samesite="None",
        max_age=7 * 24 * 60 * 60
    )

def clear_refresh_cookie(response: Response):
    response.delete_cookie("refresh_token")

def get_refresh_token_from_request(request: Request) -> str | None:
    return request.cookies.get("refresh_token")