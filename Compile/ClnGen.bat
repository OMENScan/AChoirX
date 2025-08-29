@Echo Off
Echo.
Echo ClnGen 1.0 - Regenerates all AChoir Cleaner Versions
Echo Ctl-C to Exit or
Pause
REM *********************************************************************
REM * AChoir 64Bit Compile both regular and elevated (manifest)         *
REM *********************************************************************
Del  AChCleanr.exe
:Ach64
Echo.
Echo Compiling AChoirX Cleaner for Windows - 64Bit...
set CGO_ENABLED=0
set GOARCH=amd64
set GOOS=windows
go build AChCleanr.go AchCln_win.go
:AchOSXAMD
Echo.
Echo Compiling AChoirX Cleaner for OSX AMD64...
set CGO_ENABLED=0
set GOARCH=amd64
set GOOS=darwin
go build AChCleanr.go AchCln_osx.go
Copy AChCleanr AChCleanr-AMD /y
REM *********************************************************************
REM * AChoir OSX Compile - ARM 64 -Apple Silicon                        *
REM *********************************************************************
:AchOSXARM
Echo.
Echo Compiling AChoirX Cleaner for OSX ARM64...
set CGO_ENABLED=0
set GOARCH=arm64
set GOOS=darwin
go build AChCleanr.go AchCln_osx.go
Copy AChCleanr AChCleanr-ARM /y
REM *********************************************************************
REM * AChoir Linux Compile                                              *
REM *********************************************************************
:AchLin
Echo.
Echo Compiling AChoirX for Linux...
set CGO_ENABLED=0
set GOARCH=amd64
set GOOS=linux
go build AChCleanr.go AchCln_lin.go
