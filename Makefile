# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test -cover

VERSION=`git log --pretty=format:'%ad-%h' --date=short | head -1`
LDFLAGS := -ldflags "-s -w -X main.buildInfo=${VERSION}"

CMD_LIST := cmd/emmparser/emmparser.go

ALL_LIST = $(CMD_LIST)

UPX = /opt/local/bin/upx -9 --brute -v -o 

all:
	GOOS="windows" go build -v -x -o bin/emmparser.exe cmd/emmparser/emmparser.go
	GOOS="freebsd" go build -v -x -o bin/emmparser cmd/emmparser/emmparser.go
	GOOS="darwin" go build -v -x -o bin/emmparser.darwin cmd/emmparser/emmparser.go

release:
	GOOS="windows" ${GOBUILD} ${LDFLAGS} -v -x -o bin/emmparser_win.exe \
		 cmd/emmparser/emmparser.go
	${UPX} bin/emmparser.exe bin/emmparser_win.exe

	GOOS="darwin" ${GOBUILD} ${LDFLAGS} -v -x -o bin/emmparser_darwin \
		 cmd/emmparser/emmparser.go
	${UPX} bin/emmparser.darwin bin/emmparser_darwin

	GOOS="linux" ${GOBUILD} ${LDFLAGS} -v -x -o bin/emmparser_linux \
		 cmd/emmparser/emmparser.go
	${UPX} bin/emmparser.linux bin/emmparser_linux

	GOOS="freebsd" ${GOBUILD} ${LDFLAGS} -v -x -o bin/emmparser.freebsd \
		 cmd/emmparser/emmparser.go

	zip bin/emmparser-${VERSION}.zip bin/emmparser.exe bin/emmparser.linux \
		bin/emmparser.freebsd bin/emmparser.darwin


test:
	${GOTEST} ./...

cover:
	go test -v -coverprofile=coverage.out github.com/certeu/emmchan/emm
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

clean:
	${GOCLEAN}
	rm -rf bin/*
