---
name: CLI, Windows & PowerShell Mastery
description: Comprehensive command-line expertise across Bash, PowerShell, Windows CMD, and cross-platform terminal operations for maximum productivity and system control.
version: 1.0.0
author: Custom Meta-Skill
tags: [cli, bash, powershell, windows, terminal, command-line, shell-scripting, automation]
---

# CLI, Windows & PowerShell Mastery

## Purpose
Master command-line operations across all major platforms — Linux/macOS (Bash/Zsh), Windows (PowerShell/CMD), and cross-platform tools — for system administration, automation, development, and troubleshooting.

## Bash/Linux Shell Mastery

### Essential Patterns
- **Piping**: `command1 | command2 | command3` — chain operations
- **Redirection**: `>` (overwrite), `>>` (append), `2>&1` (stderr to stdout), `2>/dev/null` (suppress errors)
- **Process substitution**: `diff <(command1) <(command2)`
- **Command substitution**: `result=$(command)` or `` result=`command` ``
- **Here documents**: `cat <<EOF ... EOF`
- **Background jobs**: `command &`, `nohup command &`, `jobs`, `fg`, `bg`

### Power Tools
- **find**: `find /path -name "*.py" -mtime -7 -exec grep -l "pattern" {} \;`
- **grep/ripgrep**: `rg -i --type py "pattern" /path` (ripgrep is 10x faster)
- **sed**: `sed -i 's/old/new/g' file` (in-place substitution)
- **awk**: `awk -F',' '{print $2, $4}' data.csv` (column extraction)
- **xargs**: `find . -name "*.log" | xargs rm` (parallel execution with `-P`)
- **jq**: `cat data.json | jq '.results[] | {name, value}'` (JSON processing)
- **curl**: `curl -s -H "Authorization: Bearer $TOKEN" https://api.example.com/data | jq .`
- **watch**: `watch -n 5 'command'` (repeat every 5 seconds)
- **tmux/screen**: Session persistence, split panes, detach/reattach

### Shell Scripting Best Practices
- Always use `set -euo pipefail` at the top of scripts
- Quote all variables: `"$variable"` not `$variable`
- Use `[[ ]]` instead of `[ ]` for conditionals
- Use `$(command)` instead of backticks
- Use functions for reusable logic
- Add `trap` for cleanup on exit
- Use `shellcheck` for linting

### File System Operations
- **Disk usage**: `du -sh */ | sort -rh | head -20`
- **File counting**: `find . -type f | wc -l`
- **Large files**: `find / -type f -size +100M -exec ls -lh {} \;`
- **Permissions**: `chmod`, `chown`, `umask`
- **Links**: `ln -s target link` (symbolic), `ln target link` (hard)
- **Archives**: `tar -czf archive.tar.gz dir/`, `tar -xzf archive.tar.gz`

## Windows PowerShell Mastery

### Core Concepts
- **Everything is an object** — not text streams like Bash
- **Verb-Noun naming**: `Get-Process`, `Set-Location`, `New-Item`
- **Pipeline passes objects**: `Get-Process | Where-Object {$_.CPU -gt 100} | Sort-Object CPU -Descending`
- **Aliases exist**: `ls` → `Get-ChildItem`, `cd` → `Set-Location`, `cat` → `Get-Content`

### Essential Cmdlets
- **Get-Help**: `Get-Help Get-Process -Full` (documentation)
- **Get-Command**: `Get-Command *service*` (discover commands)
- **Get-Member**: `Get-Process | Get-Member` (inspect object properties)
- **Select-Object**: `Get-Process | Select-Object Name, CPU, WorkingSet`
- **Where-Object**: `Get-Service | Where-Object {$_.Status -eq 'Running'}`
- **ForEach-Object**: `1..10 | ForEach-Object { $_ * 2 }`
- **Measure-Object**: `Get-ChildItem | Measure-Object -Property Length -Sum`

### System Administration
- **Services**: `Get-Service`, `Start-Service`, `Stop-Service`, `Restart-Service`
- **Processes**: `Get-Process`, `Stop-Process`, `Start-Process`
- **Registry**: `Get-ItemProperty HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion`
- **Event Logs**: `Get-EventLog -LogName System -Newest 50`
- **Network**: `Test-NetConnection`, `Get-NetAdapter`, `Get-NetIPAddress`
- **Firewall**: `Get-NetFirewallRule`, `New-NetFirewallRule`
- **Users/Groups**: `Get-LocalUser`, `Get-LocalGroup`, `Add-LocalGroupMember`
- **Scheduled Tasks**: `Get-ScheduledTask`, `Register-ScheduledTask`

### PowerShell Scripting
- Use `$ErrorActionPreference = 'Stop'` for strict error handling
- Use `try/catch/finally` for error handling
- Use `[CmdletBinding()]` for advanced functions
- Use `[Parameter(Mandatory)]` for required parameters
- Use `-WhatIf` and `-Confirm` for destructive operations
- Use `Write-Verbose` for debug output

### Remote Operations
- **PSRemoting**: `Enter-PSSession -ComputerName server01`
- **Invoke-Command**: `Invoke-Command -ComputerName server01 -ScriptBlock { Get-Process }`
- **WinRM**: Configure with `Enable-PSRemoting -Force`

## Windows CMD (Legacy but Essential)

### Key Commands
- **dir**: `dir /s /b *.txt` (recursive, bare format)
- **findstr**: `findstr /s /i "pattern" *.log` (grep equivalent)
- **robocopy**: `robocopy source dest /mir /mt:8` (robust file copy)
- **netstat**: `netstat -ano | findstr :80` (network connections)
- **tasklist/taskkill**: `tasklist /fi "imagename eq chrome.exe"`
- **sfc**: `sfc /scannow` (system file checker)
- **DISM**: `DISM /Online /Cleanup-Image /RestoreHealth`
- **wmic**: `wmic os get caption,version` (system info)

## Cross-Platform Patterns

### Equivalent Commands
| Task | Linux/macOS | PowerShell | CMD |
|------|------------|------------|-----|
| List files | `ls -la` | `Get-ChildItem` | `dir` |
| Find text | `grep -r "text" .` | `Select-String -Path *.txt -Pattern "text"` | `findstr /s "text" *.*` |
| Process list | `ps aux` | `Get-Process` | `tasklist` |
| Network info | `ifconfig`/`ip addr` | `Get-NetIPAddress` | `ipconfig` |
| Kill process | `kill -9 PID` | `Stop-Process -Id PID -Force` | `taskkill /PID PID /F` |
| Disk space | `df -h` | `Get-PSDrive` | `wmic logicaldisk get size,freespace` |
| Environment | `echo $VAR` | `$env:VAR` | `echo %VAR%` |

### MCP Integration for Desktop Operations
When operating on a user's Windows computer via MCP:
- Use `desktop:{sessionId}` prefix for shell sessions
- Files are mounted at `/mnt/desktop/C:\`
- Prefer PowerShell over CMD for modern Windows
- Always test commands non-destructively first
- Use `-WhatIf` flag before destructive PowerShell operations
- Check execution policy: `Get-ExecutionPolicy`
- For GUI automation, use PowerShell with COM objects or UIAutomation

## Automation Patterns

### Task Automation Checklist
1. Can this be done with a one-liner? (prefer simplicity)
2. Will this be repeated? (write a script)
3. Does it need error handling? (add try/catch or set -e)
4. Does it need logging? (add timestamps and output capture)
5. Does it need scheduling? (cron on Linux, Task Scheduler on Windows)
6. Does it need to be cross-platform? (use Python or Node.js instead)
