PKG = $(shell grep '^module ' go.mod | cut -f2 -d ' ')

VERSION = $(shell git describe --tags --exact-match 2> /dev/null || echo dev)

GOFILES = $(shell go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}}{{"\n"}}{{end}}' ./...)

GOOS   ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOARM  ?= $(shell go env GOARM)

buildSuffix ?= -$(GOOS)-$(GOARCH)

KWCTL = kwctl$(buildSuffix)
KWUI = kwui$(buildSuffix)

BINS = $(KWCTL) $(KWUI)

COMPILE =	go build -o $@ -ldflags '-X $(PKG)/internal/version.Version=$(VERSION)'

all: $(BINS)

lint:
	golangci-lint run

.PHONY: kwctl
kwctl: $(KWCTL)

.PHONY: kwui
kwui: $(KWUI)

$(KWCTL): $(GOFILES)
	$(COMPILE) ./cmd/kwctl

$(KWUI): $(GOFILES)
	$(COMPILE) ./cmd/kwui

clean:
	rm -f $(BINS)

realclean:
	rm -f kwctl-* kwui-* kwctl kwui
