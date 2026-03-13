all:
	build test

build:
	go build -o stewfish

test:
	go test -v ./...

windows:
	GOOS=windows GOARCH=amd64 go build -o stewfish.exe

clean:
	rm -f stewfish

.PHONY: all build test clean
