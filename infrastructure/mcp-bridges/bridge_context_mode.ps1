# Launcher for Context Mode MCP bridge (called by run_hidden.vbs)
# AUDIT-006: Use environment variables instead of hardcoded user paths
$UserProfile = $env:USERPROFILE
& "$UserProfile\.universal-ai\tools\run_mcp_bridge_forever.ps1" `
    -Name "UniversalAI-ContextModeMcp" `
    -Port 8877 `
    -Command "C:\Program Files\nodejs\node.exe" `
    -LogPath "C:\ProgramData\universal-ai-mcps\context-mode-mcp.log" `
    "$UserProfile\AppData\Roaming\npm\node_modules\context-mode\cli.bundle.mjs"
