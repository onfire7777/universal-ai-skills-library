# Launcher for Skill Seekers MCP bridge (called by run_hidden.vbs)
# AUDIT-006: Use environment variables instead of hardcoded user paths
$UserProfile = $env:USERPROFILE
& "$UserProfile\.manus\tools\run_mcp_bridge_forever.ps1" `
    -Name "Manus-SkillSeekersMcp" `
    -Port 8875 `
    -Command "$UserProfile\.skill-seekers\venv\Scripts\python.exe" `
    -LogPath "C:\ProgramData\manus-mcps\skill-seekers-mcp.log" `
    -m skill_seekers.mcp.server_fastmcp
