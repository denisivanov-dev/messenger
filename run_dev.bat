@echo off

echo === backend ===
start cmd /k "python -m uvicorn backend_python.main:app --reload"

echo === frontend ===
start cmd /k "cd frontend && npm run dev"

pause