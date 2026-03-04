all:
	build test

build:
	go build -o stewfish

test:
	go test -v ./...

clean:
	rm -f stewfish

.PHONY: all build test clean