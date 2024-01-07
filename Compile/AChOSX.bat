@Echo Off
REM *********************************************************************
REM * AChoir OSX Compile - AMD 64 - Intel                               *
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
REM * AChoir OSX Compile - ARM 64 - M1                                  *
REM *********************************************************************
:AchOSXARM
Echo.
Echo Compiling AChoirX for OSX ARM64...
set CGO_ENABLED=0
set GOARCH=arm64
set GOOS=darwin
go build AChoirX.go only_OSX.go
Copy AChoirX AChoirX-ARM /y
