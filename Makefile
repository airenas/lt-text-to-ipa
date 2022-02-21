-include Makefile.options
#####################################################################################
## print usage information
help:
	@echo 'Usage:'
	@cat ${MAKEFILE_LIST} | grep -e "^## " -A 1 | grep -v '\-\-' | sed 's/^##//' | cut -f1 -d":" | \
		awk '{info=$$0; getline; print "  " $$0 ": " info;}' | column -t -s ':' | sort 
.PHONY: help
#####################################################################################
commit_count=$(shell git rev-list --count HEAD)
go_build_cmd=CGO_ENABLED=0 go build -ldflags "-X main.version=git-version.$(commit_count)" 

generate: 
	go get github.com/petergtz/pegomock/...
	go generate ./...

bin:
	mkdir -p $@

build: | bin
	cd cmd/text-to-ipa && $(go_build_cmd) -o $(CURDIR)/bin/text-to-ipa

run:
	cd cmd/text-to-ipa/ && go run --race . -c config.yml	

## call units tests
test/unit: 
	go test -v -race -count=1 ./...
.PHONY: test/unit
#####################################################################################
## code vet and lint
test/lint: 
	go vet ./...
	go install golang.org/x/lint/golint@latest
	golint -set_exit_status ./...
.PHONY: test/lint
#####################################################################################
## build docker image
docker/build:
	cd build && $(MAKE) dbuild
.PHONY: docker/build
#####################################################################################
## build and push text-to-ipa docker image
docker/push:
	cd build && $(MAKE) dpush
.PHONY: docker/push
#####################################################################################
## cleans all temporary data
clean:
	rm -rf bin

.PHONY: clean run generate
