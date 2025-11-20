PKG = $(shell grep '^module ' go.mod | cut -f2 -d ' ')

VERSION = $(shell git describe --tags --exact-match 2> /dev/null || echo dev)
COMMIT = $(shell git rev-parse --short=10 HEAD)
DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%S")

GOLDFLAGS = \
	    -X '$(PKG)/internal/commands.Version=$(VERSION)' \
	    -X '$(PKG)/internal/commands.Commit=$(COMMIT)' \
	    -X '$(PKG)/internal/commands.Date=$(DATE)'

GOFILES = $(shell go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}}{{"\n"}}{{end}}' ./...)

GOOS   ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOARM  ?= $(shell go env GOARM)

buildSuffix ?= -$(GOOS)-$(GOARCH)

KWCTL = kwctl$(buildSuffix)

BINS = $(KWCTL)

COMPILE =	GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -o $@ -ldflags "$(GOLDFLAGS)"

all: $(BINS)

lint:
	golangci-lint run

kwctl: $(KWCTL)

$(KWCTL): $(GOFILES)
	$(COMPILE) ./cmd/kwctl

clean:
	rm -f $(BINS)
