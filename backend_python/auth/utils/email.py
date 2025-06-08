import smtplib
import random
from email.mime.text import MIMEText

from data.gmail_data import email as sender_email, password as sender_password

import asyncio

def generate_code() -> str:
    return str(random.randint(100000, 999999))

def send_verification_email(to_email: str, code: str):
    print(code)
    
    subject = "Код подтверждения"
    body = f"Ваш код подтверждения: {code}"

    msg = MIMEText(body)
    msg['Subject'] = subject
    msg['From'] = sender_email
    msg['To'] = to_email

    with smtplib.SMTP("smtp.gmail.com", 587) as server:
        server.starttls()
        server.login(sender_email, sender_password)
        server.send_message(msg)