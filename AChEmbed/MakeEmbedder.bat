@Echo Off
Echo.
Echo Compiling Embedder for Windows...
set CGO_ENABLED=0
set GOARCH=amd64
set GOOS=windows
go build Embedder.go
