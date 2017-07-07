.PHONY: dev build clean

all: dev

dev: build
	./pastebin

build: clean
	go get ./...
	go build .

test:
	go test ./...

clean:
	rm -rf pastebin
