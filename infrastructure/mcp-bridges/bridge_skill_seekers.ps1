# Skill Seekers MCP Bridge
# Starts the Skill Seekers HTTP/SSE server on port 8875
# Used by: Manus (cloud), OpenWebUI, any HTTP-based MCP client
# Transport: SSE (Server-Sent Events) on http://127.0.0.1:8875/sse

$ErrorActionPreference = "Stop"
$PORT = 8875
$VENV = "$env:USERPROFILE\.skill-seekers\venv\Scripts\python.exe"
$SCRIPT = "$env:USERPROFILE\.skill-seekers\venv\Lib\site-packages\skill_seekers\mcp_server.py"
$LOG = "C:\ProgramData\manus-mcps\skill-seekers-mcp.log"

# Ensure log directory exists
New-Item -ItemType Directory -Force -Path "C:\ProgramData\manus-mcps" | Out-Null

# Kill any existing process on this port
$existing = Get-NetTCPConnection -LocalPort $PORT -State Listen -ErrorAction SilentlyContinue
if ($existing) {
    Stop-Process -Id $existing.OwningProcess -Force -ErrorAction SilentlyContinue
    Start-Sleep -Seconds 2
}

# Start Skill Seekers in HTTP mode
Write-Output "[$(Get-Date)] Starting Skill Seekers on port $PORT" | Out-File -Append $LOG
& $VENV $SCRIPT --transport sse --port $PORT 2>&1 | Out-File -Append $LOG
