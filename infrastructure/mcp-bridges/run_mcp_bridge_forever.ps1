<#
.SYNOPSIS
  Keep one local MCP HTTP bridge alive for Manus.
  Security: binds to 127.0.0.1 only. Includes log rotation (10MB max).

  Audit fixes applied:
    AUDIT-002: Logging errors caught with fallback to %TEMP% (CWE-703)
    AUDIT-004: Pinned mcp-proxy version via MCP_PROXY_VERSION env var (CWE-494)
    AUDIT-005: Process verification before skipping startup (CWE-20)
    AUDIT-008: Sensitive data redacted from log output (CWE-532)
#>
[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)][string]$Name,
    [Parameter(Mandatory = $true)][int]$Port,
    [Parameter(Mandatory = $true)][string]$Command,
    [Parameter(Mandatory = $true)][string]$LogPath,
    [Parameter(ValueFromRemainingArguments = $true)][string[]]$CommandArgs = @()
)

$ErrorActionPreference = 'Continue'
Set-StrictMode -Version Latest

[int64]$MaxLogBytes = 10MB

# AUDIT-004: Pin mcp-proxy version to prevent supply-chain attacks
$McpProxyVersion = $env:MCP_PROXY_VERSION
if ([string]::IsNullOrWhiteSpace($McpProxyVersion)) {
    # Default to a known-good version if env var not set
    $McpProxyVersion = '6.4.6'
}

$npx = Get-Command npx.cmd -ErrorAction SilentlyContinue
if (-not $npx) {
    $npx = Get-Command npx -ErrorAction SilentlyContinue
}
if (-not $npx) {
    throw "Missing prerequisite: npx is not on PATH."
}

New-Item -ItemType Directory -Path (Split-Path -Parent $LogPath) -Force | Out-Null
$env:PYTHONIOENCODING = 'utf-8'
$env:PYTHONUTF8 = '1'

# Remove the --% stop-parsing artifact that PowerShell passes as a literal string
# when scheduled tasks use --% in -File mode arguments
$CleanArgs = @($CommandArgs | Where-Object { $_ -ne '--%' })

# AUDIT-002: Log rotation helper with error fallback
function Write-BridgeLog {
    param([string]$Value)
    try {
        if ((Test-Path $LogPath) -and ((Get-Item $LogPath).Length -gt $MaxLogBytes)) {
            $archive = "$LogPath.1"
            if (Test-Path $archive) { Remove-Item $archive -Force }
            Rename-Item $LogPath $archive -Force
        }
        Add-Content -LiteralPath $LogPath -Value $Value -ErrorAction Stop
    } catch {
        # Fallback: write to %TEMP% if primary log path fails (symlink attack mitigation)
        $fallback = Join-Path $env:TEMP 'manus-mcp-bridge.log'
        $errorLine = "[{0}] [LOGGING-ERROR] {1} :: {2}" -f (Get-Date -Format o), $_.Exception.Message, $Value
        Add-Content -LiteralPath $fallback -Value $errorLine -ErrorAction SilentlyContinue
    }
}

while ($true) {
    # AUDIT-005: Verify the listener is actually our mcp-proxy process, not a port hijacker
    $listener = Get-NetTCPConnection -LocalPort $Port -State Listen -ErrorAction SilentlyContinue |
        Where-Object { $_.LocalAddress -in @('127.0.0.1', '::1', '0.0.0.0') } |
        Select-Object -First 1

    $expectedListener = $false
    if ($listener) {
        try {
            $proc = Get-CimInstance Win32_Process -Filter "ProcessId = $($listener.OwningProcess)" -ErrorAction Stop
            $expectedListener = ($proc.CommandLine -like '*mcp-proxy*') -and ($proc.CommandLine -like "*--port $Port*")
        } catch {
            $expectedListener = $false
        }
    }

    if ($expectedListener) {
        $timestamp = Get-Date -Format o
        Write-BridgeLog "[$timestamp] $Name already listening on port $Port (verified); checking again in 30s"
        Start-Sleep -Seconds 30
        continue
    }

    $timestamp = Get-Date -Format o
    Write-BridgeLog "[$timestamp] starting $Name on 127.0.0.1:$Port (mcp-proxy@$McpProxyVersion)"

    # SECURITY: --host 127.0.0.1 ensures no external network access
    # AUDIT-004: Pinned package version
    $proxyArgs = @(
        '--yes',
        '--package', "mcp-proxy@$McpProxyVersion",
        'mcp-proxy',
        '--host', '127.0.0.1',
        '--port', "$Port",
        '--',
        $Command
    ) + $CleanArgs

    & $npx.Source @proxyArgs 2>&1 | ForEach-Object {
        # AUDIT-008: Redact sensitive data (tokens, API keys) from log output
        $line = "$_" -replace '(Bearer\s+|token[=:]\s*|api[_-]?key[=:]\s*|password[=:]\s*)[^\s"'']+', '$1[REDACTED]'
        Write-BridgeLog $line
    }

    $timestamp = Get-Date -Format o
    Write-BridgeLog "[$timestamp] $Name exited with code $LASTEXITCODE; restarting in 2s"
    Start-Sleep -Seconds 2
}
