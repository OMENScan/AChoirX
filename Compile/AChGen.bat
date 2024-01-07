@Echo Off
Echo.
Echo ACHGen 1.0 - Regenerates all AChoir Versions
Echo Ctl-C to Exit or
Pause
REM *********************************************************************
REM * AChoir 32Bit Compile both regular and elevated (manifest)         *
REM *********************************************************************
Echo.
Echo Compiling AChoirX for Windows - 32Bit ...
set CGO_ENABLED=0
set GOARCH=386
set GOOS=windows
go build AChoirX.go only_windows.go
copy AChoirX.exe AChoirX32.exe /y
REM *********************************************************************
REM * Embedding a Manifest requires the Microsoft manifest utility      *
REM *********************************************************************
Echo.
Echo Checking for Manifest Utility...
if exist mt-master\x86\mt.exe GOTO Manif32
GOTO NoManif32
:Manif32
Echo.
Echo Now Embedding Manifest...
copy AChoirX.exe A-AChoirX32.exe /y
Del  AChoirX.exe
mt-master\x86\mt.exe /nologo /manifest "mt-master\x86\AChoirX.exe.manifest" /outputresource:"A-AChoirX32.exe;#1"
GOTO Ach64
:NoManif32
Echo.
Echo Manifest Utility not found, Manifest bypassed... 
REM *********************************************************************
REM * AChoir 64Bit Compile both regular and elevated (manifest)         *
REM *********************************************************************
Del  AChoirX.exe
:Ach64
Echo.
Echo Compiling AChoirX for Windows - 64Bit...
set CGO_ENABLED=0
set GOARCH=amd64
set GOOS=windows
go build AChoirX.go only_windows.go
copy AChoirX.exe A-AChoirX.exe /y
Echo.
Echo Checking for Manifest Utility...
if exist mt-master\x86\mt.exe GOTO Manif64
GOTO NoManif64
:Manif64
Echo.
Echo Now Embedding Manifest...
mt-master\X64\mt.exe /nologo /manifest "mt-master\x64\AChoirX.exe.manifest" /outputresource:"A-AChoirX.exe;#1"
GOTO AchOSXAMD
:NoManif64
Echo.
Echo Manifest Utility not found, Manifest bypassed... 
REM *********************************************************************
REM * AChoir OSX Compile AMD 64                                         *
REM *********************************************************************
:AchOSXAMD
Echo.
Echo Compiling AChoirX for OSX AMD64...
set CGO_ENABLED=0
set GOARCH=amd64
set GOOS=darwin
go build AChoirX.go only_OSX.go
Copy AChoirX AChoirX-AMD /y
REM *********************************************************************
REM * AChoir OSX Compile - ARM 64 -Apple Silicon                        *
REM *********************************************************************
:AchOSXARM
Echo.
Echo Compiling AChoirX for OSX ARM64...
set CGO_ENABLED=0
set GOARCH=arm64
set GOOS=darwin
go build AChoirX.go only_OSX.go
Copy AChoirX AChoirX-ARM /y
REM *********************************************************************
REM * AChoir Linux Compile                                              *
REM *********************************************************************
:AchLin
Echo.
Echo Compiling AChoirX for Linux...
set CGO_ENABLED=0
set GOARCH=amd64
set GOOS=linux
go build AChoirX.go only_linux.go
