# Launcher for LightPanda MCP bridge (called by run_hidden.vbs)
# AUDIT-006: Use environment variables instead of hardcoded user paths
$UserProfile = $env:USERPROFILE
& "$UserProfile\.manus\tools\run_mcp_bridge_forever.ps1" `
    -Name "Manus-LightpandaMcp" `
    -Port 8878 `
    -Command "C:\WINDOWS\System32\cmd.exe" `
    -LogPath "C:\ProgramData\manus-mcps\lightpanda-mcp.log" `
    /c "$UserProfile\.lightpanda-ai\lightpanda-mcp.cmd"
