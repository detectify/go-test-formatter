PREFIX?=$(shell pwd)

VERSION=$(shell git describe --match 'v[0-9]*' --dirty='.m' --always)
GO_LDFLAGS=-ldflags "-X `go list ./version`.Version $(VERSION)"

.PHONY: clean all fmt vet lint build test binaries
.DEFAULT: default
all: AUTHORS clean fmt vet fmt lint build test binaries

AUTHORS: .mailmap .git/HEAD
    git log --format='%aN <%aE>' | sort -fu > $@

version/version.go:
    ./version/version.sh > $@

${PREFIX}/bin/go-test-formatter: version/version.go $(shell find . -type f -name '*.go')
    @echo "+ $@"
    @go build -o $@ ${GO_LDFLAGS} ./cmd/go-test-formatter

vet:
    @echo "+ $@"
    @go vet ./...

fmt:
    @echo "+ $@"
    @test -z "$$(gofmt -s -l . | grep -v Godeps/_workspace/src/ | tee /dev/stderr)" || \
        echo "+ please format Go code with 'gofmt -s'"

lint:
    @echo "+ $@"
    @test -z "$$(golint ./... | grep -v Godeps/_workspace/src/ | tee /dev/stderr)"

build:
    @echo "+ $@"
    @go build -v ${GO_LDFLAGS} ./...

test:
    @echo "+ $@"
    @go test -test.short ./...

test-full:
    @echo "+ $@"
    @go test -v ./...

binaries: ${PREFIX}/bin/go-test-formatter
    @echo "+ $@"

clean:
    @echo "+ $@"
    @rm -rf "${PREFIX}/bin/go-test-formatter"
