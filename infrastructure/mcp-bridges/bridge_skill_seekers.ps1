# Launcher for Skill Seekers MCP bridge (called by run_hidden.vbs)
# AUDIT-006: Use environment variables instead of hardcoded user paths
$UserProfile = $env:USERPROFILE
& "$UserProfile\.universal-ai\tools\run_mcp_bridge_forever.ps1" `
    -Name "UniversalAI-SkillSeekersMcp" `
    -Port 8875 `
    -Command "$UserProfile\.skill-seekers\venv\Scripts\python.exe" `
    -LogPath "C:\ProgramData\universal-ai-mcps\skill-seekers-mcp.log" `
    -m skill_seekers.mcp.server_fastmcp
