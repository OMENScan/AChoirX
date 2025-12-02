@Echo off
Echo [+] Deleteting Old Versions
Del AChASEP.exe
Del AChDBox.exe
Del AChWinFull.exe
Del AChIISLogs.exe
Del AChWinMSF.exe
Del AChWinVolLoki.exe
Del AChWinYara.exe
Echo.
Echo [+] Building AChASEP Tool
A-AChoirX.exe /PKR:AChASEP.zip
Rename A-AChoirx_PKR.exe AChASEP.exe
Echo.
Echo [+] Building AChDBox Tool
A-AChoirX.exe /PKR:AChDBox.zip
Rename A-AChoirx_PKR.exe AChDBox.exe
Echo.
Echo [+] Building AChWinFull Tool
A-AChoirX.exe /PKR:AChWinFull.zip
Rename A-AChoirx_PKR.exe AChWinFull.exe
Echo.
Echo [+] Building AChIISLogs Tool
A-AChoirX.exe /PKR:AChIISLogs.zip
Rename A-AChoirx_PKR.exe AChIISLogs.exe
Echo.
Echo [+] Building AChWinMSF Tool
A-AChoirX.exe /PKR:AChWinMSF.zip
Rename A-AChoirx_PKR.exe AChWinMSF.exe
Echo.
Echo [+] Building AChWinVolLoki Tool
A-AChoirX.exe /PKR:AChWinVolLoki.zip
Rename A-AChoirx_PKR.exe AChWinVolLoki.exe
Echo.
Echo [+] Building AChWinYara Tool
A-AChoirX.exe /PKR:AChWinYara.zip
Rename A-AChoirx_PKR.exe AChWinYara.exe
