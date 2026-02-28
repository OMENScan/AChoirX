@echo off
REM save-chat.bat - append chat transcript from clipboard to CHAT_HISTORY.md










echo (open CHAT_HISTORY.md to review and commit changes)echo Transcript appended to CHAT_HISTORY.mdpowershell -NoProfile -Command "Add-Content -Path '..\CHAT_HISTORY.md' -Value ("`n`n## " + (Get-Date) + "`n" + (Get-Clipboard))"
:: use PowerShell to grab clipboard and append with timestamp)    echo # Chat History>"..\CHAT_HISTORY.md"    echo Creating CHAT_HISTORY.md...if not exist "..\CHAT_HISTORY.md" (:: ensure file exists