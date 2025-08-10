# werfty Makefile: build, clean, test
.PHONY: build clean test

build:
	go build -o werfty main.go

clean:
	rm -f werfty

# Run all Go tests for werfty
# Usage: make test

go-tests:
	go test ./...
