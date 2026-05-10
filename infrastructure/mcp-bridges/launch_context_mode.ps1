# Launcher for Context Mode MCP bridge (called by scheduled task)
# AUDIT-006: Use environment variables instead of hardcoded user paths
# CMD-FIX: Use 8.3 short path for node.exe to avoid mcp-proxy space-splitting
$UserProfile = $env:USERPROFILE
& "$UserProfile\.universal-ai\tools\run_mcp_bridge_forever.ps1" `
    -Name "UniversalAI-ContextModeMcp" `
    -Port 8877 `
    -Command "C:\PROGRA~1\nodejs\node.exe" `
    -LogPath "C:\ProgramData\universal-ai-mcps\context-mode-mcp.log" `
    "$UserProfile\AppData\Roaming\npm\node_modules\context-mode\cli.bundle.mjs"
