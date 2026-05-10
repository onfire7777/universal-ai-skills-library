# Universal AI Skills Installer
# Installs ALL skills (core + library) from this repository to all detected AI platform roots
# Run from the repository root directory
#
# Usage:
#   .\infrastructure\scripts\install_skills.ps1           # Install with junctions
#   .\infrastructure\scripts\install_skills.ps1 -DryRun   # Preview without changes
#   .\infrastructure\scripts\install_skills.ps1 -Force    # Overwrite existing installations
#   .\infrastructure\scripts\install_skills.ps1 -Copy     # Copy files instead of junction

param(
    [switch]$DryRun,
    [switch]$Force,
    [switch]$Copy
)

$ErrorActionPreference = "Stop"
$RepoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)

Write-Host "=== Universal AI Skills Installer ===" -ForegroundColor Cyan
Write-Host "Repository: $RepoRoot"
Write-Host ""

# ─── Identify all skills ──────────────────────────────────────────────────────
# Core skills: top-level directories with SKILL.md (not skills/, infrastructure/, .git)
$coreSkillDirs = Get-ChildItem $RepoRoot -Directory | Where-Object {
    (Test-Path (Join-Path $_.FullName "SKILL.md")) -and
    $_.Name -notin @("skills", "infrastructure", ".git", "node_modules")
}

# Library skills: in the skills/ subdirectory
$librarySkillsDir = Join-Path $RepoRoot "skills"
$librarySkillDirs = @()
if (Test-Path $librarySkillsDir) {
    $librarySkillDirs = Get-ChildItem $librarySkillsDir -Directory | Where-Object {
        Test-Path (Join-Path $_.FullName "SKILL.md")
    }
}

$totalSkills = $coreSkillDirs.Count + $librarySkillDirs.Count
Write-Host "Skills found:" -ForegroundColor Green
Write-Host "  Core skills:    $($coreSkillDirs.Count) (top-level, with scripts)"
Write-Host "  Library skills: $($librarySkillDirs.Count) (in skills/ directory)"
Write-Host "  Total:          $totalSkills"
Write-Host ""

# ─── Define AI platform roots ─────────────────────────────────────────────────
$platforms = @(
    @{ Name = "Claude Desktop"; Path = "$env:APPDATA\Claude\skills" },
    @{ Name = "Claude Code";    Path = "$env:USERPROFILE\.claude\skills" },
    @{ Name = "Cursor";         Path = "$env:USERPROFILE\.cursor\skills" },
    @{ Name = "Codex";          Path = "$env:USERPROFILE\.codex\skills" },
    @{ Name = "OpenCode";       Path = "$env:USERPROFILE\.opencode\skills" },
    @{ Name = "Gemini CLI";     Path = "$env:USERPROFILE\.gemini\skills" },
    @{ Name = "Manus-compatible"; Path = "$env:USERPROFILE\.manus\skills" }
)

# ─── Install function ─────────────────────────────────────────────────────────
function Install-SkillToTarget {
    param(
        [string]$SourceDir,
        [string]$TargetDir,
        [string]$SkillName
    )
    
    $targetSkillPath = Join-Path $TargetDir $SkillName
    
    if (Test-Path $targetSkillPath) {
        if (-not $Force) {
            $item = Get-Item $targetSkillPath -Force -ErrorAction SilentlyContinue
            if ($item -and ($item.Attributes -band [System.IO.FileAttributes]::ReparsePoint)) {
                return "linked"
            }
            return "exists"
        }
        Remove-Item $targetSkillPath -Recurse -Force
    }
    
    if ($Copy) {
        Copy-Item $SourceDir -Destination $targetSkillPath -Recurse -Force
    } else {
        New-Item -ItemType Junction -Path $targetSkillPath -Target $SourceDir | Out-Null
    }
    return "installed"
}

# ─── Install to each platform ─────────────────────────────────────────────────
$installedCount = 0
$skippedCount = 0

foreach ($platform in $platforms) {
    $targetPath = $platform.Path
    $parentDir = Split-Path $targetPath -Parent
    
    # Check if the platform is installed (parent directory exists)
    if (-not (Test-Path $parentDir)) {
        Write-Host "  SKIP: $($platform.Name) (not installed)" -ForegroundColor DarkGray
        $skippedCount++
        continue
    }
    
    Write-Host "  Installing to $($platform.Name)..." -ForegroundColor Yellow
    
    # Create skills directory if it doesn't exist
    if (-not (Test-Path $targetPath)) {
        if (-not $DryRun) {
            New-Item -ItemType Directory -Path $targetPath -Force | Out-Null
        }
    }
    
    if ($DryRun) {
        Write-Host "    [DRY RUN] Would install $totalSkills skills to: $targetPath"
        continue
    }
    
    $platformInstalled = 0
    $platformLinked = 0
    $platformSkipped = 0
    
    # Install core skills (these have scripts and are the most important)
    foreach ($skill in $coreSkillDirs) {
        $result = Install-SkillToTarget -SourceDir $skill.FullName -TargetDir $targetPath -SkillName $skill.Name
        switch ($result) {
            "installed" { $platformInstalled++ }
            "linked"    { $platformLinked++ }
            "exists"    { $platformSkipped++ }
        }
    }
    
    # Install library skills
    foreach ($skill in $librarySkillDirs) {
        $result = Install-SkillToTarget -SourceDir $skill.FullName -TargetDir $targetPath -SkillName $skill.Name
        switch ($result) {
            "installed" { $platformInstalled++ }
            "linked"    { $platformLinked++ }
            "exists"    { $platformSkipped++ }
        }
    }
    
    Write-Host "    New: $platformInstalled | Already linked: $platformLinked | Skipped: $platformSkipped" -ForegroundColor Green
    $installedCount++
}

Write-Host ""
Write-Host "=== Installation Complete ===" -ForegroundColor Green
Write-Host "  Platforms configured: $installedCount"
Write-Host "  Platforms skipped:    $skippedCount (not installed)"
Write-Host ""

if (-not $Copy) {
    Write-Host "Skills are linked via directory junctions." -ForegroundColor Cyan
    Write-Host "Updates to the repo will automatically propagate to all platforms."
} else {
    Write-Host "Skills were copied. Run again after updates to re-sync." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Next steps:" -ForegroundColor White
Write-Host "  1. Set up MCP bridges:  .\infrastructure\scripts\setup_mcp_bridges.ps1"
Write-Host "  2. Verify installation: Open any AI client and check that skills are available"
Write-Host ""
