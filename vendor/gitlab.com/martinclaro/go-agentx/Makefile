APPNAME = go-agentx
GOPATH := ${PWD}:${PWD}/vendor
GOBIN := ${PWD}/bin

.PHONY: all clean fmt test vendor_clean vendor_update vet

default: test

vendor_clean:
	rm -dRfv ./vendor/src

vendor_update:
	GOPATH=${PWD}/vendor \
		go get -v \
			'gopkg.in/errgo.v1' \
			'gitlab.com/martinclaro/go-oidsort' \
			'gitlab.com/martinclaro/go-agentx'

fmt:
	GOPATH=$(GOPATH) go fmt ./...
	@echo "Fmt done."

vet: vendor_update
	GOPATH=$(GOPATH) go vet ./...
	@echo "Vet done."

test: vendor_update
	GOPATH=$(GOPATH) go test -v

all: fmt vet test

clean:
	rm -rfv vendor
	@echo "Done."
