# Universal Skills Installer
# Installs all skills from this repository to all detected AI platform roots
# Run from the repository root directory

param(
    [switch]$DryRun,
    [switch]$Force
)

$ErrorActionPreference = "Stop"
$RepoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$SkillsSource = Join-Path $RepoRoot "skills"

Write-Host "=== Universal Skills Installer ===" -ForegroundColor Cyan
Write-Host "Source: $SkillsSource"
Write-Host ""

# Define all AI platform skill roots
$platforms = @(
    @{ Name = "Claude Desktop"; Path = "$env:APPDATA\Claude\skills" },
    @{ Name = "Claude Code";    Path = "$env:USERPROFILE\.claude\skills" },
    @{ Name = "Cursor";         Path = "$env:USERPROFILE\.cursor\skills" },
    @{ Name = "Codex";          Path = "$env:USERPROFILE\.codex\skills" },
    @{ Name = "OpenCode";       Path = "$env:USERPROFILE\.opencode\skills" },
    @{ Name = "Gemini CLI";     Path = "$env:USERPROFILE\.gemini\skills" },
    @{ Name = "Manus Local";    Path = "$env:USERPROFILE\.manus\skills" }
)

# Count skills
$skillCount = (Get-ChildItem $SkillsSource -Directory).Count
Write-Host "Skills available: $skillCount" -ForegroundColor Green
Write-Host ""

# Install to each platform
foreach ($platform in $platforms) {
    $targetPath = $platform.Path
    $parentDir = Split-Path $targetPath -Parent
    
    # Check if the platform is installed (parent directory exists)
    if (-not (Test-Path $parentDir)) {
        Write-Host "  SKIP: $($platform.Name) (not installed)" -ForegroundColor DarkGray
        continue
    }
    
    Write-Host "  Installing to $($platform.Name)..." -ForegroundColor Yellow
    
    if ($DryRun) {
        Write-Host "    [DRY RUN] Would create junction: $targetPath -> $SkillsSource"
        continue
    }
    
    # Remove existing skills directory/junction
    if (Test-Path $targetPath) {
        if ($Force) {
            Remove-Item $targetPath -Recurse -Force
        } else {
            # Check if it's already a junction pointing to our source
            $item = Get-Item $targetPath -Force
            if ($item.Attributes -band [System.IO.FileAttributes]::ReparsePoint) {
                $target = [System.IO.Path]::GetFullPath((cmd /c "dir /al `"$parentDir`"" 2>$null | Select-String $item.Name | ForEach-Object { ($_ -split '\[')[1] -replace '\]','' }))
                if ($target -eq $SkillsSource) {
                    Write-Host "    Already linked correctly" -ForegroundColor Green
                    continue
                }
            }
            Write-Host "    WARNING: $targetPath exists. Use -Force to overwrite." -ForegroundColor Red
            continue
        }
    }
    
    # Create directory junction (symlink for directories)
    New-Item -ItemType Junction -Path $targetPath -Target $SkillsSource | Out-Null
    Write-Host "    Linked: $targetPath -> $SkillsSource" -ForegroundColor Green
}

Write-Host ""
Write-Host "=== Installation Complete ===" -ForegroundColor Green
Write-Host ""
Write-Host "All platforms now share the same skills directory via junctions."
Write-Host "Updates to the repo will automatically propagate to all platforms."
Write-Host ""
Write-Host "To also set up MCP bridges, run:"
Write-Host "  .\infrastructure\scripts\setup_mcp_bridges.ps1" -ForegroundColor Cyan
