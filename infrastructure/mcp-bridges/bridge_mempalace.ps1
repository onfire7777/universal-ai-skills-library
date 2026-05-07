# MemPalace MCP Bridge
# Starts mcp-proxy to bridge MemPalace stdio MCP server to HTTP on port 8876
# Used by: Manus (cloud), OpenWebUI, any HTTP-based MCP client
# Transport: Streamable HTTP on http://127.0.0.1:8876/mcp

$ErrorActionPreference = "Stop"
$PORT = 8876
$LOG = "C:\ProgramData\manus-mcps\mempalace-mcp.log"

# Ensure log directory exists
New-Item -ItemType Directory -Force -Path "C:\ProgramData\manus-mcps" | Out-Null

# Kill any existing process on this port
$existing = Get-NetTCPConnection -LocalPort $PORT -State Listen -ErrorAction SilentlyContinue
if ($existing) {
    Stop-Process -Id $existing.OwningProcess -Force -ErrorAction SilentlyContinue
    Start-Sleep -Seconds 2
}

# Start mcp-proxy bridging MemPalace stdio to HTTP
# Uses Start-Process with single ArgumentList string to avoid PowerShell splatting issues
Write-Output "[$(Get-Date)] Starting MemPalace bridge on port $PORT" | Out-File -Append $LOG

$npxPath = (Get-Command npx.cmd -ErrorAction SilentlyContinue).Source
if (-not $npxPath) { $npxPath = "npx" }

$argString = "-y mcp-proxy --port $PORT -- npx -y @anthropic/mempalace-mcp"
Start-Process -FilePath $npxPath -ArgumentList $argString -NoNewWindow -Wait -RedirectStandardOutput $LOG -RedirectStandardError "$LOG.err"
