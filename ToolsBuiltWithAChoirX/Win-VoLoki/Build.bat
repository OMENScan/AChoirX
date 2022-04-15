@Echo Off
set GOARCH=amd64
set GOOS=windows
set CGO_ENABLED=0
Echo.
Echo Copying Most Recent Source Files...
copy ..\AChoirX.go
copy ..\only_windows.go
Echo.
Echo Compiling AChoirX for Windows - VoLoki PoC...
go build AChoirX.go only_windows.go
del AChoirX-VoLoki.exe
del A-AChoirX-VoLoki.exe
ren AChoirX.exe AChoirX-VoLoki.exe
copy AChoirX-VoLoki.exe A-AChoirX-VoLoki.exe
Echo.
Echo Now Embedding Manifest...
..\mt-master\X64\mt.exe /nologo /manifest "..\mt-master\x64\AChoirX.exe.manifest" /outputresource:"A-AChoirX-VoLoki.exe;#1"
