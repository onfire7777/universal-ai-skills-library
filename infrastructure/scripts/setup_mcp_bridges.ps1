# Universal MCP Bridge Setup Script
# Run this script as Administrator to set up all MCP bridges and watchdog
# This configures Windows Scheduled Tasks for persistent, auto-healing MCP services

param(
    [string]$InstallDir = "C:\ProgramData\manus-mcps",
    [string]$ToolsDir = "$env:USERPROFILE\.manus\tools"
)

$ErrorActionPreference = "Stop"
Write-Host "=== Universal MCP Bridge Setup ===" -ForegroundColor Cyan
Write-Host ""

# 1. Create directories
Write-Host "[1/6] Creating directories..." -ForegroundColor Yellow
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
New-Item -ItemType Directory -Force -Path $ToolsDir | Out-Null

# 2. Copy bridge scripts
Write-Host "[2/6] Installing bridge scripts..." -ForegroundColor Yellow
$scriptRoot = Split-Path -Parent $PSScriptRoot
Copy-Item "$scriptRoot\mcp-bridges\bridge_skill_seekers.ps1" "$ToolsDir\" -Force
Copy-Item "$scriptRoot\mcp-bridges\bridge_mempalace.ps1" "$ToolsDir\" -Force
Copy-Item "$scriptRoot\mcp-bridges\bridge_context_mode.ps1" "$ToolsDir\" -Force
Write-Host "  Installed 3 bridge scripts to $ToolsDir"

# 3. Install watchdog
Write-Host "[3/6] Installing watchdog..." -ForegroundColor Yellow
Copy-Item "$scriptRoot\watchdog\mcp_watchdog.ps1" "$InstallDir\" -Force
Copy-Item "$scriptRoot\watchdog\run_hidden.vbs" "$InstallDir\" -Force
Write-Host "  Installed watchdog to $InstallDir"

# 4. Create scheduled tasks for bridges
Write-Host "[4/6] Creating bridge scheduled tasks..." -ForegroundColor Yellow

$bridges = @(
    @{ Name = "Manus-SkillSeekersMcp"; Script = "$ToolsDir\bridge_skill_seekers.ps1" },
    @{ Name = "Manus-MemPalaceMcp"; Script = "$ToolsDir\bridge_mempalace.ps1" },
    @{ Name = "Manus-ContextModeMcp"; Script = "$ToolsDir\bridge_context_mode.ps1" }
)

foreach ($bridge in $bridges) {
    # Remove existing task
    Unregister-ScheduledTask -TaskName $bridge.Name -Confirm:$false -ErrorAction SilentlyContinue
    
    # Create new task
    $action = New-ScheduledTaskAction -Execute "powershell.exe" -Argument "-ExecutionPolicy Bypass -WindowStyle Hidden -File `"$($bridge.Script)`""
    $trigger = New-ScheduledTaskTrigger -AtLogOn
    $settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries -StartWhenAvailable -RestartInterval (New-TimeSpan -Minutes 1) -RestartCount 3
    
    Register-ScheduledTask -TaskName $bridge.Name -Action $action -Trigger $trigger -Settings $settings -RunLevel Highest -Description "MCP Bridge: $($bridge.Name)" | Out-Null
    Write-Host "  Created task: $($bridge.Name)"
}

# 5. Create watchdog scheduled task (every 5 minutes, hidden)
Write-Host "[5/6] Creating watchdog scheduled task..." -ForegroundColor Yellow
Unregister-ScheduledTask -TaskName "Manus-McpWatchdog" -Confirm:$false -ErrorAction SilentlyContinue

schtasks /Create /TN "Manus-McpWatchdog" /TR "wscript.exe `"$InstallDir\run_hidden.vbs`" `"$InstallDir\mcp_watchdog.ps1`"" /SC MINUTE /MO 5 /F | Out-Null
Write-Host "  Created watchdog task (every 5 minutes, hidden)"

# 6. Start everything
Write-Host "[6/6] Starting all services..." -ForegroundColor Yellow
foreach ($bridge in $bridges) {
    Start-ScheduledTask -TaskName $bridge.Name
    Write-Host "  Started: $($bridge.Name)"
}
Start-ScheduledTask -TaskName "Manus-McpWatchdog"
Write-Host "  Started: Manus-McpWatchdog"

Write-Host ""
Write-Host "=== Setup Complete ===" -ForegroundColor Green
Write-Host "Services will be available in ~15 seconds:"
Write-Host "  - Skill Seekers: http://127.0.0.1:8875/sse"
Write-Host "  - MemPalace:     http://127.0.0.1:8876/mcp"
Write-Host "  - Context Mode:  http://127.0.0.1:8877/mcp"
Write-Host ""
Write-Host "Watchdog monitors every 5 minutes and auto-restarts crashed services."
Write-Host "Logs: $InstallDir\watchdog.log"
