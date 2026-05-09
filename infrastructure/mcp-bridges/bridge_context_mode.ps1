# Launcher for Context Mode MCP bridge (called by run_hidden.vbs)
# AUDIT-006: Use environment variables instead of hardcoded user paths
$UserProfile = $env:USERPROFILE
& "$UserProfile\.manus\tools\run_mcp_bridge_forever.ps1" `
    -Name "Manus-ContextModeMcp" `
    -Port 8877 `
    -Command "C:\Program Files\nodejs\node.exe" `
    -LogPath "C:\ProgramData\manus-mcps\context-mode-mcp.log" `
    "$UserProfile\AppData\Roaming\npm\node_modules\context-mode\cli.bundle.mjs"
