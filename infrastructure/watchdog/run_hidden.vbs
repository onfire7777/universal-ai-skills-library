' Run a PowerShell script completely hidden (no window flash)
' Usage: wscript.exe run_hidden.vbs "C:\path\to\script.ps1"
Set objShell = CreateObject("WScript.Shell")
objShell.Run "powershell.exe -ExecutionPolicy Bypass -WindowStyle Hidden -File """ & WScript.Arguments(0) & """", 0, False
