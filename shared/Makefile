# shared Makefile: build, clean, test
.PHONY: build clean go-tests

build:
	go build ./...

clean:
	rm -f *.out

# Run all Go tests for shared module
# Usage: make test

go-tests:
	go test ./...
