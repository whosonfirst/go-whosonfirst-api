CWD=$(shell pwd)
GOPATH := $(CWD)

build:	rmdeps deps fmt bin

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src/github.com/whosonfirst/go-whosonfirst-api; then rm -rf src/github.com/whosonfirst/go-whosonfirst-api; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-api/client
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-api/endpoint
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-api/flags
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-api/mapzen
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-api/response
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-api/result
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-api/throttle
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-api/util
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-api/writer
	cp client/*.go src/github.com/whosonfirst/go-whosonfirst-api/client/
	cp endpoint/*.go src/github.com/whosonfirst/go-whosonfirst-api/endpoint/
	cp flags/*.go src/github.com/whosonfirst/go-whosonfirst-api/flags/
	cp mapzen/*.go src/github.com/whosonfirst/go-whosonfirst-api/mapzen/
	cp response/*.go src/github.com/whosonfirst/go-whosonfirst-api/response/
	cp result/*.go src/github.com/whosonfirst/go-whosonfirst-api/result/
	cp throttle/*.go src/github.com/whosonfirst/go-whosonfirst-api/throttle/
	cp util/*.go src/github.com/whosonfirst/go-whosonfirst-api/util/
	cp writer/*.go src/github.com/whosonfirst/go-whosonfirst-api/writer/
	cp api.go src/github.com/whosonfirst/go-whosonfirst-api/
	if test ! -d src; then mkdir src; fi
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

deps:   
	@GOPATH=$(GOPATH) go get -u "github.com/tidwall/gjson"
	@GOPATH=$(GOPATH) go get -u "github.com/tidwall/pretty"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-csv"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-placetypes"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-uri"

vendor-deps: deps
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt api.go
	go fmt client/*.go
	go fmt cmd/*.go
	go fmt endpoint/*.go
	go fmt flags/*.go
	go fmt mapzen/*.go
	go fmt response/*.go
	go fmt result/*.go
	go fmt throttle/*.go
	go fmt util/*.go
	go fmt writer/*.go

bin:	self
	@GOPATH=$(shell pwd) go build -o bin/wof-api cmd/wof-api.go
