@Echo Off
Echo.
Echo Compiling AChoirX for Windows...
set CGO_ENABLED=0
set GOARCH=amd64
set GOOS=windows
go build AChoirX.go only_windows.go
Echo.
Echo Now Embedding Manifest...
mt-master\X64\mt.exe /nologo /manifest "mt-master\x64\AChoirX.exe.manifest" /outputresource:"AChoirX.exe;#1"
