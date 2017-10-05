all:
	#GOOS="windows" go build -v -x -o bin/emmparser.exe cmd/emmparser/emmparser.go
	#GOOS="freebsd" go build -v -x -o bin/emmparser cmd/emmparser/emmparser.go
	GOOS="darwin" go build -v -x -o bin/emmparser.darwin cmd/emmparser/emmparser.go
