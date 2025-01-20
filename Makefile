MODULE = $(shell go list -m)
all: build
generate:
	go generate ./...

build: # build a server
	CGO_ENABLED=0 go build -a -o go_news_server $(MODULE)/cmd/go_news_server

release:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -o go_news_server $(MODULE)/cmd/go_news_server/
	zip -9 -r ./go_news_server.zip go_news_server

lint:
	gofmt -l .