from typing import Tuple, Optional, Dict
import re
from data.validation_data import FORBIDDEN_USERNAMES, BLOCKED_DOMAINS, COMMON_WEAK_PASSWORDS

def validate_empty_data(fields: list[tuple[str, str]]) -> tuple[bool, dict[str, str]]:
    field_error_messages = {
        "email": "- Введите почту",
        "username": "- Введите имя пользователя",
        "password": "- Введите пароль",
        "confirmPassword": "- Подтвердите пароль"
    }
    
    for field_name, value in fields:
        if not value:
            return False, field_name, field_error_messages[field_name]

    return True, None, "OK"

def validate_email(email: str) -> Tuple[bool, str]:
    email = email.lower()

    if not re.match(r"^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$", email):
        return False, "- Неверный формат email"

    domain = email.split("@")[1]

    if domain in BLOCKED_DOMAINS:
        return False, "- Недопустимый почтовый домен"

    return True, "OK"

def validate_username(username: str) -> Tuple[bool, str]:
    username = username.lower()

    if len(username) < 3 or len(username) > 20:
        return False, "- Имя пользователя должно быть от 3 до 20 символов"
    
    if username in FORBIDDEN_USERNAMES:
        return False, "- Это имя запрещено"
    
    if not re.match(r"^[a-z][a-z0-9_]*$", username):
        return False, "- Имя должно содержать только буквы и цифры"
    
    if username == username[0] * len(username):
        return False, "- Имя не может состоять из одного символа"
    
    if "__" in username:
        return False, "- Имя не должно содержать двойные подчёркивания"
    
    if username[-1] == "_":
        return False, "- Имя не должно заканчиваться подчёркиванием"

    return True, "OK"

def validate_password(password: str, confirm_password: str) -> Tuple[bool, str]:
    if not password:
        return "- Введите пароль"
    elif not confirm_password:
        return "- Подтвердите пароль"

    if password != confirm_password:
        return False, "- Пароли не совпадают"
    
    if len(password) < 8 or len(password) > 128:
        return False, "- Пароль должен быть от 8 до 128 символов"
    
    if password == password[0] * len(password):
        return False, "- Пароль не может состоять из одного символа"

    if not re.search(r"[A-Z]", password):
        return False, "- Пароль должен содержать хотя бы одну заглавную букву"

    if not re.search(r"[a-z]", password):
        return False, "- Пароль должен содержать хотя бы одну строчную букву"

    if not re.search(r"[0-9]", password):
        return False, "- Пароль должен содержать хотя бы одну цифру"

    if not re.search(r"[!@#$%^&*()_+\-=\[\]{};':\"\\|,.<>\/?]", password):
        return False, "- Пароль должен содержать хотя бы один специальный символ"
    
    if password.lower() in COMMON_WEAK_PASSWORDS:
        return False, "- Пароль слишком простой"

    return True, "OK"