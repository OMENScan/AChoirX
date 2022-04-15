@Echo Off
set GOARCH=amd64
set GOOS=windows
set CGO_ENABLED=0
Echo.
Echo Copying Most Recent Source Files...
copy ..\AChoirX.go
copy ..\only_windows.go
Echo.
Echo Compiling AChoirX for Windows - YARA PoC...
go build AChoirX.go only_windows.go
del AChoirX-yara.exe
del A-AChoirX-yara.exe
ren AChoirX.exe AChoirX-yara.exe
copy AChoirX-yara.exe A-AChoirX-yara.exe
Echo.
Echo Now Embedding Manifest...
..\mt-master\X64\mt.exe /nologo /manifest "..\mt-master\x64\AChoirX.exe.manifest" /outputresource:"A-AChoirX-yara.exe;#1"
