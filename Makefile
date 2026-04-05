all:
	build test

build:
	go build -ldflags "-X github.com/alvissraghnall/stewfish/internal.Version=0.1.0" -o stewfish

test:
	go test -v ./...

windows:
	GOOS=windows GOARCH=amd64 go build -o stewfish.exe

clean:
	rm -f stewfish

.PHONY: all build test clean
