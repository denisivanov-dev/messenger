@echo off

echo === backend_python ===
start cmd /k "venv\Scripts\python -m uvicorn backend_python.main:app --reload"

echo === backend_golang ===
start cmd /k "cd backend_golang && air"

echo === frontend ===
start cmd /k "cd frontend && npm run dev"

echo === python workers ===
start cmd /k "python -m backend_python.chat.service.background_tasks.run_scheduled"
start cmd /k "python -m backend_python.chat.service.background_tasks.run_workers"

pause