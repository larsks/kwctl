export GOARCH=arm64
export GOARM=v8

all: build

build:
	go build
