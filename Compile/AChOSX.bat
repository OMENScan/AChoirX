@Echo Off
Echo.
Echo Compiling AChoirX for OSX...
set CGO_ENABLED=0
set GOARCH=amd64
set GOOS=darwin
go build AChoirX.go only_OSX.go OSXEmbed.go
