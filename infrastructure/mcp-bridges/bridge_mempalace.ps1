# Launcher for MemPalace MCP bridge (called by run_hidden.vbs)
# AUDIT-006: Use environment variables instead of hardcoded user paths
$UserProfile = $env:USERPROFILE
& "$UserProfile\.universal-ai\tools\run_mcp_bridge_forever.ps1" `
    -Name "UniversalAI-MemPalaceMcp" `
    -Port 8876 `
    -Command "$UserProfile\AppData\Roaming\Python\Python314\Scripts\mempalace-mcp.exe" `
    -LogPath "C:\ProgramData\universal-ai-mcps\mempalace-mcp.log" `
    --palace "$UserProfile\.mempalace\palace"
