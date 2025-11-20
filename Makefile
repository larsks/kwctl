PKG = $(shell grep '^module ' go.mod | cut -f2 -d ' ')

VERSION = $(shell git describe --tags --exact-match 2> /dev/null || echo dev)
COMMIT = $(shell git rev-parse --short=10 HEAD)
DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%S")

GOLDFLAGS = \
	    -X '$(PKG)/internal/commands.Version=$(VERSION)' \
	    -X '$(PKG)/internal/commands.Commit=$(COMMIT)' \
	    -X '$(PKG)/internal/commands.Date=$(DATE)'

GOFILES = $(shell go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}}{{"\n"}}{{end}}' ./...)

all: kwctl

kwctl: $(GOFILES)
	go build -o $@ -ldflags "$(GOLDFLAGS)" ./cmd/kwctl

clean:
	rm -f kwctl
