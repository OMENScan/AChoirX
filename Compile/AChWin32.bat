@Echo Off
Echo.
Echo Compiling AChoirX for Windows - 32Bit ...
set CGO_ENABLED=0
set GOARCH=386
set GOOS=windows
go build AChoirX.go only_windows.go
copy AChoirX.exe AChoirX32.exe
copy AChoirX.exe A-AChoirX32.exe
Del  AChoirX.exe
Echo.
Echo Now Embedding Manifest...
mt-master\x86\mt.exe /nologo /manifest "mt-master\x86\AChoirX.exe.manifest" /outputresource:"A-AChoirX32.exe;#1"
