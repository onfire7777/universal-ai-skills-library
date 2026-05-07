# MCP Services Watchdog
# Monitors ports 8875, 8876, 8877 and auto-restarts any crashed bridge
# Runs every 5 minutes via Windows Scheduled Task
# Logs to C:\ProgramData\manus-mcps\watchdog.log

$ErrorActionPreference = "SilentlyContinue"
$LOG = "C:\ProgramData\manus-mcps\watchdog.log"
$MAX_LOG_SIZE = 1MB

# Rotate log if too large
if ((Test-Path $LOG) -and (Get-Item $LOG).Length -gt $MAX_LOG_SIZE) {
    Move-Item $LOG "$LOG.old" -Force
}

function Write-Log($msg) {
    "[$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')] $msg" | Out-File -Append $LOG
}

$services = @(
    @{ Name = "Skill Seekers"; Port = 8875; Task = "Manus-SkillSeekersMcp" },
    @{ Name = "MemPalace";     Port = 8876; Task = "Manus-MemPalaceMcp" },
    @{ Name = "Context Mode";  Port = 8877; Task = "Manus-ContextModeMcp" }
)

$allOk = $true

foreach ($svc in $services) {
    $listening = Get-NetTCPConnection -LocalPort $svc.Port -State Listen -ErrorAction SilentlyContinue
    
    if (-not $listening) {
        $allOk = $false
        Write-Log "WARNING: $($svc.Name) (port $($svc.Port)) is DOWN. Restarting..."
        
        # Stop and restart the scheduled task
        Stop-ScheduledTask -TaskName $svc.Task -ErrorAction SilentlyContinue
        Start-Sleep -Seconds 2
        Start-ScheduledTask -TaskName $svc.Task -ErrorAction SilentlyContinue
        
        # Wait up to 30 seconds for port to come up
        $waited = 0
        while ($waited -lt 30) {
            Start-Sleep -Seconds 5
            $waited += 5
            $check = Get-NetTCPConnection -LocalPort $svc.Port -State Listen -ErrorAction SilentlyContinue
            if ($check) {
                Write-Log "RECOVERED: $($svc.Name) is back UP on port $($svc.Port)"
                break
            }
        }
        
        if ($waited -ge 30) {
            Write-Log "CRITICAL: $($svc.Name) failed to restart after 30s"
        }
    }
}

if ($allOk) {
    # Only log "all OK" once per hour to keep log clean
    $lastOk = Get-Content $LOG -Tail 1 -ErrorAction SilentlyContinue
    if ($lastOk -notmatch "All services OK") {
        Write-Log "HEALTH: All services OK (8875, 8876, 8877 listening)"
    }
}
