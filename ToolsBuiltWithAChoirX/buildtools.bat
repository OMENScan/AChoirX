@Echo off
Echo [+] Deleteting Old Versions
Del A-AChoirX_PKR.exe
Del AChASEP\AChASEP.exe
Del AChDBox\AChDBox.exe
Del AChWinFull\AChWinFull.exe
Del AChIISLogs\AChIISLogs.exe
Del Win-MSF\AChWinMSF.exe
Del Win-VolLoki\AChWinVolLoki.exe
Del Win-Yara\AChWinYara.exe
Echo.
Echo [+] Building AChASEP Tool
A-AChoirX.exe /PKR:AChASEP\AChASEP.zip
Move A-AChoirx_PKR.exe AChASEP\AChASEP.exe
Echo.
Echo [+] Building AChDBox Tool
A-AChoirX.exe /PKR:AChDBox\AChDBox.zip
Move A-AChoirx_PKR.exe AChDBox\AChDBox.exe
Echo.
Echo [+] Building AChWinFull Tool
A-AChoirX.exe /PKR:AChWinFull\AChWinFull.zip
Move A-AChoirx_PKR.exe AChWinFull\AChWinFull.exe
Echo.
Echo [+] Building AChIISLogs Tool
A-AChoirX.exe /PKR:AChIISLogs\AChIISLogs.zip
Move A-AChoirx_PKR.exe AChIISLogs\AChIISLogs.exe
Echo.
Echo [+] Building AChWinMSF Tool
A-AChoirX.exe /PKR:Win-MSF\AChWinMSF.zip
Move A-AChoirx_PKR.exe Win-MSF\AChWinMSF.exe
Echo.
Echo [+] Building AChWinVolLoki Tool
A-AChoirX.exe /PKR:Win-Voloki\AChWinVolLoki.zip
Move A-AChoirx_PKR.exe Win-VoLoki\AChWinVolLoki.exe
Echo.
Echo [+] Building AChWinYara Tool
A-AChoirX.exe /PKR:Win-Yara\AChWinYara.zip
Move A-AChoirx_PKR.exe Win-Yara\AChWinYara.exe
