@Echo Off
Echo.
Echo Compiling AChoirX for Linux...
set CGO_ENABLED=0
set GOARCH=amd64
set GOOS=linux
go build AChoirX.go only_linux.go LinEmbed.go
