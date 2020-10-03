@Echo Off
Echo.
Echo Compiling AChoirX for Windows...
set CGO_ENABLED=0
set GOARCH=amd64
set GOOS=windows
go build AChoirX.go only_windows.go
