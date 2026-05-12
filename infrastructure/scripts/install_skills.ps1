# Universal AI Skills Installer
# Default mode installs only compact router wrapper skills. This keeps agent
# roots connected to skill-router without copying the full corpus everywhere.
#
# Usage:
#   .\infrastructure\scripts\install_skills.ps1
#   .\infrastructure\scripts\install_skills.ps1 -DryRun
#   .\infrastructure\scripts\install_skills.ps1 -Force
#   .\infrastructure\scripts\install_skills.ps1 -Copy
#   .\infrastructure\scripts\install_skills.ps1 -SkillNames universal-ai-skills,printable-cards
#   .\infrastructure\scripts\install_skills.ps1 -FullCopy   # explicit, intentionally redundant
#   .\infrastructure\scripts\install_skills.ps1 -Target "$env:USERPROFILE\.codex\skills"

param(
    [switch]$DryRun,
    [switch]$Force,
    [switch]$Copy,
    [switch]$FullCopy,
    [string[]]$SkillNames = @(),
    [string]$Target
)

$ErrorActionPreference = "Stop"
$RepoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$LibrarySkillsDir = Join-Path $RepoRoot "skills"
$DefaultWrapperSkillNames = @("universal-ai-skills", "printable-cards")

Write-Host "=== Universal AI Skills Installer ===" -ForegroundColor Cyan
Write-Host "Repository: $RepoRoot"
Write-Host "Mode: $(if ($FullCopy) { 'full-copy' } elseif ($SkillNames.Count -gt 0) { 'selected' } else { 'wrapper-only' })"
Write-Host ""

if (-not (Test-Path $LibrarySkillsDir)) {
    throw "Missing skills directory: $LibrarySkillsDir"
}

$allSkillDirs = Get-ChildItem $LibrarySkillsDir -Directory | Where-Object {
    Test-Path (Join-Path $_.FullName "SKILL.md")
}

if ($FullCopy) {
    $selectedSkillDirs = $allSkillDirs
} else {
    $wanted = if ($SkillNames.Count -gt 0) { $SkillNames } else { $DefaultWrapperSkillNames }
    $wantedSet = @{}
    foreach ($name in $wanted) {
        $wantedSet[$name] = $true
    }
    $selectedSkillDirs = $allSkillDirs | Where-Object { $wantedSet.ContainsKey($_.Name) }
    $selectedNames = @($selectedSkillDirs | ForEach-Object { $_.Name })
    $missing = @($wanted | Where-Object { $_ -notin $selectedNames })
    if ($missing.Count -gt 0) {
        throw "Missing selected skill(s): $($missing -join ', ')"
    }
}

Write-Host "Skills selected: $($selectedSkillDirs.Count)" -ForegroundColor Green
foreach ($skill in $selectedSkillDirs) {
    Write-Host "  - $($skill.Name)"
}
Write-Host ""

if ($Target) {
    $platforms = @(@{ Name = "Explicit target"; Path = $Target })
} else {
    $platforms = @(
        @{ Name = "OpenSkills / .agent"; Path = "$env:USERPROFILE\.agent\skills" },
        @{ Name = "Claude Code"; Path = "$env:USERPROFILE\.claude\skills" },
        @{ Name = "Codex"; Path = "$env:USERPROFILE\.codex\skills" },
        @{ Name = "Manus-compatible"; Path = "$env:USERPROFILE\.manus\skills" },
        @{ Name = "Gemini CLI"; Path = "$env:USERPROFILE\.gemini\skills" },
        @{ Name = "Cursor"; Path = "$env:USERPROFILE\.cursor\skills" },
        @{ Name = "OpenCode"; Path = "$env:USERPROFILE\.config\opencode\skills" },
        @{ Name = "Kiro"; Path = "$env:USERPROFILE\.kiro\skills" }
    )
}

function Install-SkillToTarget {
    param(
        [string]$SourceDir,
        [string]$TargetDir,
        [string]$SkillName
    )

    $targetSkillPath = Join-Path $TargetDir $SkillName

    if (Test-Path $targetSkillPath) {
        if (-not $Force) {
            return "exists"
        }
        Remove-Item -LiteralPath $targetSkillPath -Recurse -Force
    }

    if ($Copy) {
        Copy-Item -LiteralPath $SourceDir -Destination $targetSkillPath -Recurse -Force
    } else {
        New-Item -ItemType Junction -Path $targetSkillPath -Target $SourceDir | Out-Null
    }
    return "installed"
}

$configuredCount = 0
$skippedCount = 0

foreach ($platform in $platforms) {
    $targetPath = $platform.Path
    $parentDir = Split-Path $targetPath -Parent

    if (-not (Test-Path $parentDir)) {
        Write-Host "  SKIP: $($platform.Name) (parent missing: $parentDir)" -ForegroundColor DarkGray
        $skippedCount++
        continue
    }

    Write-Host "  Installing to $($platform.Name)..." -ForegroundColor Yellow

    if ($DryRun) {
        Write-Host "    [DRY RUN] Would ensure $targetPath"
        Write-Host "    [DRY RUN] Would install $($selectedSkillDirs.Count) selected skill(s)"
        $configuredCount++
        continue
    }

    if (-not (Test-Path $targetPath)) {
        New-Item -ItemType Directory -Path $targetPath -Force | Out-Null
    }

    $installed = 0
    $existing = 0
    foreach ($skill in $selectedSkillDirs) {
        $result = Install-SkillToTarget -SourceDir $skill.FullName -TargetDir $targetPath -SkillName $skill.Name
        switch ($result) {
            "installed" { $installed++ }
            "exists" { $existing++ }
        }
    }

    Write-Host "    Installed: $installed | Existing: $existing" -ForegroundColor Green
    $configuredCount++
}

Write-Host ""
Write-Host "=== Installation Complete ===" -ForegroundColor Green
Write-Host "  Platforms configured: $configuredCount"
Write-Host "  Platforms skipped:    $skippedCount"
Write-Host ""

if ($FullCopy) {
    Write-Host "Full-copy mode was explicitly requested. This is not the recommended default." -ForegroundColor Yellow
} elseif (-not $Copy) {
    Write-Host "Wrapper skills are linked via directory junctions." -ForegroundColor Cyan
} else {
    Write-Host "Wrapper skills were copied. Run again after updates to re-sync." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Next steps:" -ForegroundColor White
Write-Host "  1. Verify: skill-router doctor"
Write-Host "  2. Inspect matrix: skill-router sync matrix"
