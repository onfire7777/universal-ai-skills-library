' Run a PowerShell script completely hidden (no window flash at all)
' Usage: wscript.exe run_hidden.vbs "path\to\script.ps1" [args...]
' The key: WScript.Shell.Run with window style 0 prevents ANY window from appearing.
'
' Security fixes applied:
'   AUDIT-001: Escape embedded quotes to prevent command injection (CWE-78)
'   Newline/CR characters rejected to prevent argument corruption

Option Explicit

Dim shell, cmd, i, arg
Set shell = CreateObject("WScript.Shell")

cmd = "powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Bypass -WindowStyle Hidden -File"

For i = 0 To WScript.Arguments.Count - 1
    arg = WScript.Arguments(i)

    ' Reject control characters that can corrupt command-line parsing
    If InStr(arg, vbCr) > 0 Or InStr(arg, vbLf) > 0 Then
        WScript.Quit 1
    End If

    ' Escape embedded double quotes for safe command-line passing
    arg = Replace(arg, Chr(34), Chr(34) & Chr(34))
    cmd = cmd & " " & Chr(34) & arg & Chr(34)
Next

shell.Run cmd, 0, False
