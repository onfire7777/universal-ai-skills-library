# Launcher for LightPanda MCP bridge (called by run_hidden.vbs)
# AUDIT-006: Use environment variables instead of hardcoded user paths
$UserProfile = $env:USERPROFILE
$LogPath = "C:\ProgramData\universal-ai-mcps\lightpanda-mcp.log"
$DockerPipe = "\\.\pipe\dockerDesktopLinuxEngine"

if (-not (Test-Path $DockerPipe)) {
    New-Item -ItemType Directory -Path (Split-Path -Parent $LogPath) -Force | Out-Null
    Add-Content -LiteralPath $LogPath -Value ("[{0}] SKIP: Docker Desktop Linux engine is not running; LightPanda MCP bridge is optional." -f (Get-Date -Format o))
    exit 0
}

& "$UserProfile\.universal-ai\tools\run_mcp_bridge_forever.ps1" `
    -Name "UniversalAI-LightpandaMcp" `
    -Port 8878 `
    -Command "C:\WINDOWS\System32\cmd.exe" `
    -LogPath $LogPath `
    /c "$UserProfile\.lightpanda-ai\lightpanda-mcp.cmd"
