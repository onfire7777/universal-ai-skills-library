# MCP Bridge Watchdog - Ensures all 4 MCP services are always running
# Runs every 5 minutes via scheduled task. Checks port availability and restarts if needed.
# Location: C:\ProgramData\manus-mcps\mcp_watchdog.ps1
#
# Audit fixes applied:
#   AUDIT-003: Process verification instead of raw TCP connect (CWE-400)
#   AUDIT-009: LightPanda (port 8878) added to monitoring (CWE-840)
#   AUDIT-010: Task stop errors logged instead of silently discarded (CWE-391)
#   AUDIT-011: Regex pattern match variable usage fixed (CWE-754)

$LogFile = "C:\ProgramData\manus-mcps\watchdog.log"
$MaxLogSize = 1MB

# Rotate log if too large
if ((Test-Path $LogFile) -and (Get-Item $LogFile).Length -gt $MaxLogSize) {
    Move-Item $LogFile "$LogFile.old" -Force
}

function Write-Log {
    param([string]$Message)
    $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    Add-Content -Path $LogFile -Value "[$timestamp] $Message"
}

# AUDIT-003: Use Get-NetTCPConnection + process verification instead of raw TCP connect
# This prevents port hijacking attacks and resource leaks
function Test-Port {
    param([int]$Port)

    $listener = Get-NetTCPConnection -LocalPort $Port -State Listen -ErrorAction SilentlyContinue |
        Where-Object { $_.LocalAddress -in @('127.0.0.1', '::1', '0.0.0.0') } |
        Select-Object -First 1

    if (-not $listener) { return $false }

    # Verify the owning process is actually mcp-proxy for this port
    try {
        $proc = Get-CimInstance Win32_Process -Filter "ProcessId = $($listener.OwningProcess)" -ErrorAction Stop
        return ($proc.CommandLine -like '*mcp-proxy*') -and ($proc.CommandLine -like "*--port $Port*")
    } catch {
        return $false
    }
}

# AUDIT-010: Proper error handling for task stop/start operations
function Restart-McpService {
    param(
        [string]$TaskName,
        [int]$Port,
        [string]$ServiceName
    )

    Write-Log "WARN: $ServiceName (port $Port) is DOWN. Restarting task '$TaskName'..."

    # Stop the task if it's in a bad state
    try {
        Stop-ScheduledTask -TaskName $TaskName -ErrorAction Stop
        Start-Sleep -Seconds 2
    } catch {
        Write-Log "WARN: Failed to stop task '$TaskName' cleanly: $($_.Exception.Message). Continuing with start attempt."
    }

    # Start the task
    try {
        Start-ScheduledTask -TaskName $TaskName -ErrorAction Stop
        Write-Log "INFO: Started task '$TaskName'. Waiting for port $Port..."

        # Wait up to 30 seconds for the port to come up
        $attempts = 0
        while ($attempts -lt 6) {
            Start-Sleep -Seconds 5
            if (Test-Port -Port $Port) {
                Write-Log "OK: $ServiceName (port $Port) is back UP."
                return $true
            }
            $attempts++
        }

        Write-Log "ERROR: $ServiceName (port $Port) failed to start after 30 seconds."
        return $false
    } catch {
        Write-Log "ERROR: Failed to start task '$TaskName': $($_.Exception.Message)"
        return $false
    }
}

# --- Main Watchdog Logic ---

# AUDIT-009: Include LightPanda (port 8878) in monitoring
$services = @(
    @{ Name = 'Skill Seekers'; Port = 8875; Task = 'Manus-SkillSeekersMcp' },
    @{ Name = 'MemPalace';     Port = 8876; Task = 'Manus-MemPalaceMcp' },
    @{ Name = 'Context Mode';  Port = 8877; Task = 'Manus-ContextModeMcp' },
    @{ Name = 'LightPanda';    Port = 8878; Task = 'Manus-LightpandaMcp' }
)

$allHealthy = $true

foreach ($svc in $services) {
    if (-not (Test-Port -Port $svc.Port)) {
        $allHealthy = $false
        Restart-McpService -TaskName $svc.Task -Port $svc.Port -ServiceName $svc.Name
    }
}

if ($allHealthy) {
    # Only log periodic health confirmations every hour (check if last OK was >55 min ago)
    $lastOk = Select-String -Path $LogFile -Pattern "HEALTH: All services OK" -ErrorAction SilentlyContinue | Select-Object -Last 1
    $shouldLog = $true
    if ($lastOk) {
        # AUDIT-011: Fix regex variable usage - use if() directly with -match
        if ($lastOk.Line -match "^\[(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\]") {
            $lastTime = [datetime]::ParseExact($matches[1], "yyyy-MM-dd HH:mm:ss", $null)
            if (((Get-Date) - $lastTime).TotalMinutes -lt 55) {
                $shouldLog = $false
            }
        }
    }
    if ($shouldLog) {
        Write-Log "HEALTH: All services OK (8875, 8876, 8877, 8878 listening)"
    }
}
