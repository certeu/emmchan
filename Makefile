all:
	#GOOS="windows" go build -v -x -o bin/emmparser.exe cmd/emmparser/emmparser.go
	#GOOS="freebsd" go build -v -x -o bin/emmparser cmd/emmparser/emmparser.go
	GOOS="darwin" go build -v -x -o bin/emmparser.darwin cmd/emmparser/emmparser.go

test:
	go test -v -cover github.com/ics/emm/pkg/rss
	go test -v -cover github.com/ics/emm/pkg/emm

coverage:
	go test -v -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
