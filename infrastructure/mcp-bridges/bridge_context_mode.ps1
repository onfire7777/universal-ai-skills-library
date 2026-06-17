# Launcher for Context Mode MCP bridge (called by run_hidden.vbs)
# AUDIT-006: Use environment variables instead of hardcoded user paths
# BUGFIX: node.exe is invoked via the 8.3 short path (C:\PROGRA~1\nodejs) so the
# "Program Files" space does not split the -Command argument under mcp-proxy.
$UserProfile = $env:USERPROFILE
& "$UserProfile\.universal-ai\tools\run_mcp_bridge_forever.ps1" `
    -Name "UniversalAI-ContextModeMcp" `
    -Port 8877 `
    -Command "C:\PROGRA~1\nodejs\node.exe" `
    -LogPath "C:\ProgramData\universal-ai-mcps\context-mode-mcp.log" `
    "$UserProfile\AppData\Roaming\npm\node_modules\context-mode\cli.bundle.mjs"
