.PHONY: dev build clean

all: dev

dev: build
	./pastebin -bind 127.0.0.1:8000

build: clean
	go get ./...
	go build -o ./pastebin .

clean:
	rm -rf bin pastebin
