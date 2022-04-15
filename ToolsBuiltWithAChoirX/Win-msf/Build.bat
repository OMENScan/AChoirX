@Echo Off
set GOARCH=amd64
set GOOS=windows
set CGO_ENABLED=0
Echo.
Echo Copying Most Recent Source Files...
copy ..\AChoirX.go
copy ..\only_windows.go
Echo.
Echo Compiling AChoirX for Windows - MSF PoC...
go build AChoirX.go only_windows.go
del AChoirX-msf.exe
del A-AChoirX-msf.exe
ren AChoirX.exe AChoirX-msf.exe
copy AChoirX-msf.exe A-AChoirX-msf.exe
Echo.
Echo Now Embedding Manifest...
..\mt-master\X64\mt.exe /nologo /manifest "..\mt-master\x64\AChoirX.exe.manifest" /outputresource:"A-AChoirX-msf.exe;#1"
