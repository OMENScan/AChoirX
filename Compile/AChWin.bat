@Echo Off
Echo.
Echo Compiling AChoirX for Windows...
set CGO_ENABLED=0
set GOARCH=amd64
set GOOS=windows
go build AChoirX.go only_windows.go WinEmbed.go
copy AChoirX.exe A-AChoirX.exe
Echo.
Echo Now Embedding Manifest...
mt-master\X64\mt.exe /nologo /manifest "mt-master\x64\AChoirX.exe.manifest" /outputresource:"A-AChoirX.exe;#1"
