@Echo Off
Echo.
Echo Compiling AChoirX for Windows - 32Bit ...
set CGO_ENABLED=0
set GOARCH=386
set GOOS=windows
go build AChoirX.go only_windows.go
