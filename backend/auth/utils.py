from typing import Tuple
import re

FORBIDDEN_USERNAMES = {
    "admin", "root", "superuser", "me", "self", "you", "null", "none",
    "undefined", "true", "false", "owner", "moderator", "mod", "sys",
    "system", "support", "help", "contact", "about", "info", "user", "username",
    "login", "logout", "register", "signup", "signin", "auth", "authentication",
    "api", "v1", "v2", "adminpanel", "backend", "frontend", "server", "client",
    "guest", "test", "testing", "bot", "staff", "official", "account", "accounts",
    "email", "password", "session", "config", "security", "settings", "profile",
    "team", "dev", "developer", "superadmin", "database", "data", "json", "public",
    "private", "search", "explore", "feed", "home", "dashboard", "docs", "static",
    "media", "upload", "downloads", "terms", "privacy", "legal", "status", "api_key"
}

BLOCKED_DOMAINS = {
    "mailinator.com", "tempmail.com", "10minutemail.com", "example.com",
    "guerrillamail.com", "throwawaymail.com", "maildrop.cc", "getairmail.com",
    "fakeinbox.com", "trashmail.com", "dispostable.com", "mintemail.com",
    "moakt.com", "mailnesia.com", "sharklasers.com", "spamgourmet.com",
    "mailtemp.net", "yopmail.com", "tempinbox.com", "emailondeck.com",
    "spambog.com", "mytemp.email", "tmpmail.net", "meltmail.com",
    "trashmail.de", "spambox.us", "fakemail.net", "mailcatch.com",
    "instantemailaddress.com", "emailsensei.com", "spam4.me", "pokemail.net",
    "throwam.com", "dropmail.me", "tmail.ws", "nowmymail.com",
    "tempail.com", "eyepaste.com", "temp-mail.org", "kultmail.com",
    "yopmail.fr", "nospam.today", "binkmail.com", "fakebox.org"
}

COMMON_WEAK_PASSWORDS = {
    "123456", "123456789", "12345678", "12345", "1234567", "1234567890",
    "111111", "000000", "qwerty", "password", "123123", "abc123", "qwerty123",
    "1q2w3e4r", "admin", "letmein", "welcome", "monkey", "dragon", "passw0rd",
    "qweasd", "zaq12wsx", "user", "test", "love", "sunshine", "iloveyou",
    "login", "starwars", "hello", "freedom", "whatever", "qazwsx", "trustno1",
    "superman", "ninja", "q1w2e3r4", "asdfgh", "qwertyuiop", "baseball",
    "football", "shadow", "master", "michael", "jordan", "fuckyou",
    "killer", "hottie", "george", "tigger", "pokemon", "batman", "cheese",
    "696969", "pepper", "princess", "flower", "jessica", "ashley",
    "bailey", "cookie", "maggie", "buster", "soccer", "harley", "ranger",
    "123qwe", "1qaz2wsx", "welcome1", "pass123", "11111111", "1234qwer"
}

def validate_username(username: str) -> Tuple[bool, str]:
    username_to_check = username.lower()
    
    if len(username_to_check) < 3 or len(username_to_check) > 20:
        return False, "Имя пользователя должно быть от 3 до 20 символов"
    
    if username_to_check in FORBIDDEN_USERNAMES:
        return False, "Это имя запрещено"
    
    if not re.match(r"^[a-z][a-z0-9_]*$", username_to_check):
        return False, "Имя должно начинаться с буквы и содержать только буквы, цифры и подчёркивания"

    if username == username[0] * len(username):
        return False, "Имя не может состоять из одного символа"
    
    if "__" in username_to_check:
        return False, "Имя не должно содержать двойные подчёркивания"

    if username_to_check[-1] == "_":
        return False, "Имя не должно заканчиваться подчёркиванием"

    return True, "OK"

def validate_email(email: str) -> Tuple[bool, str]:
    email = email.lower()

    if not re.match(r"(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$)", email):
        return False, "Неверный формат email"

    domain = email.split("@")[1]

    if domain in BLOCKED_DOMAINS:
        return False, "Недопустимый почтовый домен"

    return True, "OK"

def validate_password(password: str) -> Tuple[bool, str]:
    if len(password) < 8 or len(password) > 128:
        return False, "Пароль должен быть от 8 до 128 символов"
    
    if password == password[0] * len(password):
        return False, "Пароль не может состоять из одного символа"

    if not re.search(r"[A-Z]", password):
        return False, "Пароль должен содержать хотя бы одну заглавную букву"

    if not re.search(r"[a-z]", password):
        return False, "Пароль должен содержать хотя бы одну строчную букву"

    if not re.search(r"[0-9]", password):
        return False, "Пароль должен содержать хотя бы одну цифру"

    if not re.search(r"[!@#$%^&*()_+\-=\[\]{};':\"\\|,.<>\/?]", password):
        return False, "Пароль должен содержать хотя бы один специальный символ"
    
    if password.lower() in COMMON_WEAK_PASSWORDS:
        return False, "Пароль слишком простой"

    return True, "OK"