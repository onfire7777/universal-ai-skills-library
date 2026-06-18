# Universal MCP Bridge Setup Script
# Run this script to set up optional local MCP bridges and the watchdog
# This configures Windows Scheduled Tasks for persistent, auto-healing MCP services

param(
    [string]$InstallDir = "C:\ProgramData\universal-ai-mcps",
    [string]$ToolsDir = "$env:USERPROFILE\.universal-ai\tools"
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
Copy-Item "$scriptRoot\mcp-bridges\bridge_lightpanda.ps1" "$ToolsDir\" -Force
Copy-Item "$scriptRoot\mcp-bridges\run_mcp_bridge_forever.ps1" "$ToolsDir\" -Force
Write-Host "  Installed 4 bridge scripts and runner to $ToolsDir"

# 3. Install watchdog
Write-Host "[3/6] Installing watchdog..." -ForegroundColor Yellow
Copy-Item "$scriptRoot\watchdog\mcp_watchdog.ps1" "$InstallDir\" -Force
Copy-Item "$scriptRoot\watchdog\run_hidden.vbs" "$InstallDir\" -Force
Write-Host "  Installed watchdog to $InstallDir"

# 4. Create scheduled tasks for bridges
Write-Host "[4/6] Creating bridge scheduled tasks..." -ForegroundColor Yellow

$bridges = @(
    @{ Name = "UniversalAI-SkillSeekersMcp"; Script = "$ToolsDir\bridge_skill_seekers.ps1" },
    @{ Name = "UniversalAI-MemPalaceMcp"; Script = "$ToolsDir\bridge_mempalace.ps1" },
    @{ Name = "UniversalAI-ContextModeMcp"; Script = "$ToolsDir\bridge_context_mode.ps1" },
    @{ Name = "UniversalAI-LightpandaMcp"; Script = "$ToolsDir\bridge_lightpanda.ps1" }
)

$legacyTaskPrefix = "Man" + "us-"
$legacyTasks = @(
    "SkillSeekersMcp",
    "MemPalaceMcp",
    "ContextModeMcp",
    "LightPandaMcp",
    "LightpandaMcp",
    "McpWatchdog"
) | ForEach-Object { "$legacyTaskPrefix$_" }

foreach ($legacyTask in $legacyTasks) {
    Unregister-ScheduledTask -TaskName $legacyTask -Confirm:$false -ErrorAction SilentlyContinue
}

foreach ($bridge in $bridges) {
    # Remove existing task
    Unregister-ScheduledTask -TaskName $bridge.Name -Confirm:$false -ErrorAction SilentlyContinue

    # Create new task with schtasks.exe for stable quoting across Windows hosts.
    $action = "powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Bypass -WindowStyle Hidden -File `"$($bridge.Script)`""
    & schtasks.exe /Create /TN $bridge.Name /TR $action /SC ONCE /ST 23:59 /F | Out-Null
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to create scheduled task: $($bridge.Name)"
    }
    Write-Host "  Created task: $($bridge.Name)"
}

# 5. Create watchdog scheduled task (every 5 minutes, hidden)
Write-Host "[5/6] Creating watchdog scheduled task..." -ForegroundColor Yellow
Unregister-ScheduledTask -TaskName "UniversalAI-McpWatchdog" -Confirm:$false -ErrorAction SilentlyContinue

$watchdogAction = "powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Bypass -WindowStyle Hidden -File `"$InstallDir\mcp_watchdog.ps1`""
& schtasks.exe /Create /TN "UniversalAI-McpWatchdog" /TR $watchdogAction /SC MINUTE /MO 5 /F | Out-Null
if ($LASTEXITCODE -ne 0) {
    throw "Failed to create scheduled task: UniversalAI-McpWatchdog"
}
Write-Host "  Created watchdog task (every 5 minutes, hidden)"

# 6. Start everything
Write-Host "[6/6] Starting all services..." -ForegroundColor Yellow
foreach ($bridge in $bridges) {
    if (($bridge.Name -eq "UniversalAI-LightpandaMcp") -and (-not (Test-Path "\\.\pipe\dockerDesktopLinuxEngine"))) {
        Write-Host "  Skipped: $($bridge.Name) (Docker Desktop Linux engine not running)"
        continue
    }
    Start-ScheduledTask -TaskName $bridge.Name
    Write-Host "  Started: $($bridge.Name)"
}
try {
    Start-ScheduledTask -TaskName "UniversalAI-McpWatchdog" -ErrorAction Stop
    Write-Host "  Started: UniversalAI-McpWatchdog"
} catch {
    Write-Warning "  Watchdog task was created but Windows refused a manual start for this trigger. It will run on its schedule."
}

Write-Host ""
Write-Host "=== Setup Complete ===" -ForegroundColor Green
Write-Host "Services will be available in ~15 seconds:"
Write-Host "  - Skill Seekers: http://127.0.0.1:8875/sse"
Write-Host "  - MemPalace:     http://127.0.0.1:8876/mcp"
Write-Host "  - Context Mode:  http://127.0.0.1:8877/mcp"
Write-Host ""
Write-Host "Watchdog monitors every 5 minutes and auto-restarts crashed services."
Write-Host "Logs: $InstallDir\watchdog.log"
