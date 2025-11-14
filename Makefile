PKG = $(shell grep '^module ' go.mod | cut -f2 -d ' ')

VERSION = $(shell git describe --tags --exact-match 2> /dev/null || echo dev)
COMMIT = $(shell git rev-parse --short=10 HEAD)
DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%S")

GOLDFLAGS = \
	    -X '$(PKG)/internal/commands.Version=$(VERSION)' \
	    -X '$(PKG)/internal/commands.Commit=$(COMMIT)' \
	    -X '$(PKG)/internal/commands.Date=$(DATE)'

all: build

build:
	go build -ldflags "$(GOLDFLAGS)" ./cmd/kwctl
