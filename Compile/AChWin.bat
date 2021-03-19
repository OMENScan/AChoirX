@Echo Off
Echo.
Echo Compiling AChoirX for Windows...
set GOARCH=amd64
set GOOS=windows
set CGO_ENABLED=0
go build AChoirX.go only_windows.go
REM ***
REM * Include DATA RACE Condition checking
REM * set CGO_ENABLED=1
REM * go build -race AChoirX.go only_windows.go
REM ***
copy AChoirX.exe A-AChoirX.exe
Echo.
Echo Now Embedding Manifest...
mt-master\X64\mt.exe /nologo /manifest "mt-master\x64\AChoirX.exe.manifest" /outputresource:"A-AChoirX.exe;#1"
