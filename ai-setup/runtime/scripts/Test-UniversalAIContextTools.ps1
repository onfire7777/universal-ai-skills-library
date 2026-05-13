param(
  [switch]$Deep,
  [switch]$StartLightpanda,
  [switch]$LeaveLightpandaRunning,
  [switch]$JsonOnly
)

$ErrorActionPreference = 'Continue'
$HomeDir = $env:USERPROFILE
$Root = Join-Path $HomeDir '.universal-ai-stack'
$StateDir = Join-Path $Root 'state'

function Invoke-Checked {
  param(
    [string]$FilePath,
    [string[]]$Arguments = @(),
    [int]$TimeoutSeconds = 30
  )

  function Quote-Arg {
    param([string]$Value)
    if ($null -eq $Value) { return '""' }
    if ($Value -notmatch '[\s"]') { return $Value }
    return '"' + ($Value -replace '"', '\"') + '"'
  }

  $psi = New-Object System.Diagnostics.ProcessStartInfo
  if ([System.IO.Path]::GetExtension($FilePath) -in @('.cmd', '.bat')) {
    $psi.FileName = $env:ComSpec
    $psi.Arguments = '/d /s /c ' + (Quote-Arg $FilePath) + ' ' + (($Arguments | ForEach-Object { Quote-Arg $_ }) -join ' ')
  } else {
    $psi.FileName = $FilePath
    $psi.Arguments = (($Arguments | ForEach-Object { Quote-Arg $_ }) -join ' ')
  }
  $psi.RedirectStandardOutput = $true
  $psi.RedirectStandardError = $true
  $psi.UseShellExecute = $false
  $psi.CreateNoWindow = $true

  $process = New-Object System.Diagnostics.Process
  $process.StartInfo = $psi
  try {
    [void]$process.Start()
    if (!$process.WaitForExit($TimeoutSeconds * 1000)) {
      try { $process.Kill($true) } catch {}
      return [ordered]@{ ok = $false; timedOut = $true; exitCode = $null; output = ''; error = "timed out after $TimeoutSeconds seconds" }
    }
    $stdout = $process.StandardOutput.ReadToEnd()
    $stderr = $process.StandardError.ReadToEnd()
    return [ordered]@{
      ok = ($process.ExitCode -eq 0)
      timedOut = $false
      exitCode = $process.ExitCode
      output = $stdout.Trim()
      error = $stderr.Trim()
    }
  } catch {
    return [ordered]@{ ok = $false; timedOut = $false; exitCode = $null; output = ''; error = $_.Exception.Message }
  } finally {
    $process.Dispose()
  }
}

function Test-DockerEngine {
  $docker = Get-Command docker -ErrorAction SilentlyContinue
  if (!$docker) { return [ordered]@{ ok = $false; error = 'docker command missing' } }
  $info = Invoke-Checked -FilePath $docker.Source -Arguments @('info', '--format', '{{.OSType}} {{.ServerVersion}}') -TimeoutSeconds 15
  if (!$info.ok) { return [ordered]@{ ok = $false; error = ($info.error + ' ' + $info.output).Trim() } }
  return [ordered]@{ ok = $true; info = $info.output }
}

function Test-Url {
  param([string]$Url)
  try {
    $response = Invoke-WebRequest -Uri $Url -UseBasicParsing -TimeoutSec 10
    return [ordered]@{ ok = $true; status = [int]$response.StatusCode; body = $response.Content }
  } catch {
    $status = $null
    if ($_.Exception.Response) { $status = [int]$_.Exception.Response.StatusCode }
    return [ordered]@{ ok = $false; status = $status; error = $_.Exception.Message }
  }
}

$contextModeCmd = Join-Path $env:APPDATA 'npm\context-mode.cmd'
$contextMode = Get-Command $contextModeCmd -ErrorAction SilentlyContinue
if (!$contextMode) { $contextMode = Get-Command context-mode -ErrorAction SilentlyContinue }
$contextPackage = Join-Path $env:APPDATA 'npm\node_modules\context-mode\package.json'
$contextVersion = $null
if (Test-Path -LiteralPath $contextPackage) {
  try { $contextVersion = (Get-Content -LiteralPath $contextPackage -Raw | ConvertFrom-Json).version } catch {}
}

$codexConfig = Join-Path $HomeDir '.codex\config.toml'
$codexConfigText = ''
if (Test-Path -LiteralPath $codexConfig) {
  $codexConfigText = [System.IO.File]::ReadAllText($codexConfig)
}
$codexHooks = Join-Path $HomeDir '.codex\hooks.json'
$codexHooksText = ''
if (Test-Path -LiteralPath $codexHooks) {
  $codexHooksText = [System.IO.File]::ReadAllText($codexHooks)
}
$codexHookMatchers = @()
$expectedCodexHookEvents = @('PreToolUse', 'PostToolUse', 'SessionStart', 'UserPromptSubmit', 'Stop', 'PreCompact')
$codexHookEventsPresent = @{}
$codexContextHookTimeoutsOk = $true
if ($codexHooksText) {
  try {
    $codexHooksJson = $codexHooksText | ConvertFrom-Json
    foreach ($expectedEvent in $expectedCodexHookEvents) {
      $codexHookEventsPresent[$expectedEvent] = $false
    }
    foreach ($event in $codexHooksJson.hooks.PSObject.Properties.Name) {
      foreach ($entry in @($codexHooksJson.hooks.$event)) {
        if ($entry.PSObject.Properties.Name -contains 'matcher' -and $entry.matcher) {
          $codexHookMatchers += [string]$entry.matcher
        }
        foreach ($hook in @($entry.hooks)) {
          if ($hook.command -match 'context-mode hook codex') {
            if ($codexHookEventsPresent.ContainsKey($event)) {
              $codexHookEventsPresent[$event] = $true
            }
            if (!($hook.PSObject.Properties.Name -contains 'timeout') -or [int]$hook.timeout -ne 30) {
              $codexContextHookTimeoutsOk = $false
            }
          }
        }
      }
    }
  } catch {}
}
$codexAllLifecycleHooksConfigured = $true
foreach ($expectedEvent in $expectedCodexHookEvents) {
  if (!$codexHookEventsPresent.ContainsKey($expectedEvent) -or !$codexHookEventsPresent[$expectedEvent]) {
    $codexAllLifecycleHooksConfigured = $false
  }
}

$contextDoctor = $null
if ($Deep -and $contextMode) {
  $contextDoctor = Invoke-Checked -FilePath $contextMode.Source -Arguments @('doctor') -TimeoutSeconds 60
}

$lightpandaRoot = Join-Path $HomeDir '.lightpanda-ai'
$lightpandaFetch = Join-Path $lightpandaRoot 'lightpanda-fetch.cmd'
$lightpandaMcp = Join-Path $lightpandaRoot 'lightpanda-mcp.cmd'
$lightpandaServe = Join-Path $lightpandaRoot 'lightpanda-serve.ps1'
$lightpandaStop = Join-Path $lightpandaRoot 'lightpanda-stop.ps1'
$dockerEngine = Test-DockerEngine
$lightpandaFetchSmoke = $null
if ($Deep -and $dockerEngine.ok -and (Test-Path -LiteralPath $lightpandaFetch)) {
  $fetch = Invoke-Checked -FilePath $lightpandaFetch -Arguments @('--dump', 'markdown', 'https://example.com') -TimeoutSeconds 90
  $lightpandaFetchSmoke = [ordered]@{
    ok = ($fetch.ok -and $fetch.output -match 'Example Domain')
    exitCode = $fetch.exitCode
    hasExpectedContent = [bool]($fetch.output -match 'Example Domain')
    error = $fetch.error
  }
}

$lightpandaCdp = $null
if ($Deep -and $StartLightpanda -and $dockerEngine.ok -and (Test-Path -LiteralPath $lightpandaServe)) {
  $serve = Invoke-Checked -FilePath 'powershell.exe' -Arguments @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', $lightpandaServe) -TimeoutSeconds 120
  $version = Test-Url -Url 'http://127.0.0.1:9222/json/version'
  $lightpandaCdp = [ordered]@{
    ok = ($serve.ok -and $version.ok -and $version.body -match 'webSocketDebuggerUrl')
    serveExitCode = $serve.exitCode
    endpointOk = $version.ok
    hasWebSocketDebuggerUrl = [bool]($version.body -match 'webSocketDebuggerUrl')
    error = ($serve.error + ' ' + $version.error).Trim()
  }
  if (!$LeaveLightpandaRunning -and (Test-Path -LiteralPath $lightpandaStop)) {
    [void](Invoke-Checked -FilePath 'powershell.exe' -Arguments @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', $lightpandaStop) -TimeoutSeconds 30)
  }
}

$result = [ordered]@{
  time = (Get-Date).ToString('o')
  contextMode = [ordered]@{
    commandPresent = [bool]$contextMode
    command = if ($contextMode) { $contextMode.Source } else { $null }
    version = $contextVersion
    codexMcpRegistered = [bool]($codexConfigText -match '(?m)^\[mcp_servers\."context-mode"\]')
    codexHooksConfigured = [bool]($codexHooksText -match 'context-mode hook codex pretooluse') -and [bool]($codexHooksText -match 'context-mode hook codex posttooluse')
    codexAllLifecycleHooksConfigured = $codexAllLifecycleHooksConfigured
    codexContextHookTimeouts30 = $codexContextHookTimeoutsOk
    codexHookMatchersNoLookaround = -not [bool]($codexHookMatchers -match '\(\?<?[!=]')
    doctor = $contextDoctor
  }
  lightpanda = [ordered]@{
    rootPresent = Test-Path -LiteralPath $lightpandaRoot
    fetchWrapperPresent = Test-Path -LiteralPath $lightpandaFetch
    mcpWrapperPresent = Test-Path -LiteralPath $lightpandaMcp
    serveScriptPresent = Test-Path -LiteralPath $lightpandaServe
    stopScriptPresent = Test-Path -LiteralPath $lightpandaStop
    dockerEngine = $dockerEngine
    fetchSmoke = $lightpandaFetchSmoke
    cdpSmoke = $lightpandaCdp
  }
}

$failures = New-Object System.Collections.Generic.List[string]
if (!$result.contextMode.commandPresent) { $failures.Add('contextMode.commandPresent') | Out-Null }
if (!$result.contextMode.codexMcpRegistered) { $failures.Add('contextMode.codexMcpRegistered') | Out-Null }
if (!$result.contextMode.codexHooksConfigured) { $failures.Add('contextMode.codexHooksConfigured') | Out-Null }
if (!$result.contextMode.codexAllLifecycleHooksConfigured) { $failures.Add('contextMode.codexAllLifecycleHooksConfigured') | Out-Null }
if (!$result.contextMode.codexContextHookTimeouts30) { $failures.Add('contextMode.codexContextHookTimeouts30') | Out-Null }
if (!$result.contextMode.codexHookMatchersNoLookaround) { $failures.Add('contextMode.codexHookMatchersNoLookaround') | Out-Null }
if ($Deep -and (!$result.contextMode.doctor -or !$result.contextMode.doctor.ok)) { $failures.Add('contextMode.doctor') | Out-Null }
if (!$result.lightpanda.rootPresent) { $failures.Add('lightpanda.rootPresent') | Out-Null }
if (!$result.lightpanda.fetchWrapperPresent) { $failures.Add('lightpanda.fetchWrapperPresent') | Out-Null }
if (!$result.lightpanda.mcpWrapperPresent) { $failures.Add('lightpanda.mcpWrapperPresent') | Out-Null }
if (!$result.lightpanda.serveScriptPresent) { $failures.Add('lightpanda.serveScriptPresent') | Out-Null }
if (!$result.lightpanda.stopScriptPresent) { $failures.Add('lightpanda.stopScriptPresent') | Out-Null }
if ($Deep -and !$result.lightpanda.dockerEngine.ok) { $failures.Add('lightpanda.dockerEngine') | Out-Null }
if ($Deep -and (!$result.lightpanda.fetchSmoke -or !$result.lightpanda.fetchSmoke.ok)) { $failures.Add('lightpanda.fetchSmoke') | Out-Null }
if ($Deep -and $StartLightpanda -and (!$result.lightpanda.cdpSmoke -or !$result.lightpanda.cdpSmoke.ok)) { $failures.Add('lightpanda.cdpSmoke') | Out-Null }

$result['ok'] = ($failures.Count -eq 0)
$result['failures'] = @($failures)

New-Item -ItemType Directory -Force -Path $StateDir | Out-Null
$utf8NoBom = New-Object System.Text.UTF8Encoding($false)
[System.IO.File]::WriteAllText((Join-Path $StateDir 'last-context-tools-test.json'), ($result | ConvertTo-Json -Depth 10), $utf8NoBom)
$result | ConvertTo-Json -Depth 10
