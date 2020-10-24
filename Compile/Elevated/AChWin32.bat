@Echo Off
Echo.
Echo Compiling AChoirX for Windows - 32Bit ...
set CGO_ENABLED=0
set GOARCH=386
set GOOS=windows
go build AChoirX.go only_windows.go
Echo.
Echo Now Embedding Manifest...
mt-master\x86\mt.exe /nologo /manifest "mt-master\x86\AChoirX.exe.manifest" /outputresource:"AChoirX.exe;#1"
Del AChoirX-32.exe
Rename AChoirX.exe AChoirX-32.exe
