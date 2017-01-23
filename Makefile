CWD=$(shell pwd)
GOPATH := $(CWD)

build:	rmdeps deps fmt bin

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src/github.com/whosonfirst/go-whosonfirst-api; then rm -rf src/github.com/whosonfirst/go-whosonfirst-api; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-api/api
	cp api.go src/github.com/whosonfirst/go-whosonfirst-api/
	cp api/*.go src/github.com/whosonfirst/go-whosonfirst-api/api
	if test ! -d src; then mkdir src; fi
	cp -r vendor/src/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

deps:   
	@GOPATH=$(GOPATH) go get -u "github.com/tidwall/gjson"

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor/src; then rm -rf vendor/src; fi
	cp -r src vendor/src
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt api.go

bin:	self
	@GOPATH=$(shell pwd) go build -o bin/test cmd/test.go
