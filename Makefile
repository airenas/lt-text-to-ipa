-include Makefile.options

commit_count=$(shell git rev-list --count HEAD)
go_build_cmd=CGO_ENABLED=0 go build -ldflags "-X main.version=git-version.$(commit_count)" 

test: 
	go test -v ./...

generate: 
	go get github.com/petergtz/pegomock/...
	go generate ./...

bin:
	mkdir -p $@

build: | bin
	cd cmd/text-to-ipa && $(go_build_cmd) -o $(CURDIR)/bin/text-to-ipa

run:
	cd cmd/text-to-ipa/ && go run --race . -c config.yml	

clean:
	rm -rf bin

.PHONY: clean build run generate
