# multitool Makefile: build, clean, test
.PHONY: build clean go-tests

build:
	go build -o mt main.go

clean:
	rm -f mt

# Run all Go tests for multitool
# Usage: make test

go-tests:
	go test ./...
