@echo off
set "UNIVERSAL_AI_STACK_HOME=%USERPROFILE%\.universal-ai-stack"
set "UNIVERSAL_AI_STACK_BASE_URL=http://127.0.0.1:18100/v1"
set "UNIVERSAL_AI_STACK_MODEL=auto-coding"
powershell -NoProfile -ExecutionPolicy Bypass -Command ". '%USERPROFILE%\.universal-ai-stack\env\universal-ai-stack.ps1'; Write-Host 'Loaded for this PowerShell child process.'"
