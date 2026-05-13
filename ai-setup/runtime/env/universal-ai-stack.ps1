# Dot-source this file to load Universal AI Stack settings into the current PowerShell session:
# . $env:USERPROFILE\.universal-ai-stack\env\universal-ai-stack.ps1

$UniversalAIStackHome = Join-Path $env:USERPROFILE '.universal-ai-stack'
$UniversalAIStackEnv = Join-Path $UniversalAIStackHome 'secrets\.env'

if (!(Test-Path -LiteralPath $UniversalAIStackEnv)) {
  throw "Universal AI Stack secrets file not found: $UniversalAIStackEnv"
}

foreach ($raw in Get-Content -LiteralPath $UniversalAIStackEnv) {
  $line = $raw.Trim()
  if (!$line -or $line.StartsWith('#') -or !$line.Contains('=')) { continue }
  $parts = $line.Split('=', 2)
  $name = $parts[0].Trim()
  $value = $parts[1].Trim().Trim('"')
  if (!$name) { continue }
  Set-Item -Path "Env:$name" -Value $value
}

$env:UNIVERSAL_AI_STACK_HOME = $UniversalAIStackHome
$env:UNIVERSAL_AI_STACK_BASE_URL = 'http://127.0.0.1:18100/v1'
$env:UNIVERSAL_AI_STACK_MODEL = 'auto-coding'
if ($env:UNIVERSAL_AI_STACK_API_KEY) {
  $env:OPENAI_BASE_URL = $env:UNIVERSAL_AI_STACK_BASE_URL
  $env:OPENAI_API_KEY = $env:UNIVERSAL_AI_STACK_API_KEY
}

Write-Host "Universal AI Stack loaded: $env:UNIVERSAL_AI_STACK_BASE_URL model=$env:UNIVERSAL_AI_STACK_MODEL"
