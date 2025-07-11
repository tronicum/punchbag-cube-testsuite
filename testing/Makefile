# Testing Makefile to control builds and tests from the testing directory

.PHONY: help build test test-azure test-azure-sim test-azure-direct clean fix-deps

# Default target
help:
	@echo "Available targets:"
	@echo "  fix-deps        Fix Go module dependencies"
	@echo "  build           Build all binaries"
	@echo "  test           Run all tests"
	@echo "  test-azure     Test Azure functions (simulation mode)"
	@echo "  test-azure-sim Test Azure functions (simulation mode)"
	@echo "  test-azure-direct Test Azure functions (direct mode)"
	@echo "  clean          Clean build artifacts"

# Fix dependencies first
fix-deps:
	@echo "Fixing Go module dependencies..."
	cd ../multitool && go get github.com/aws/aws-sdk-go-v2/aws
	cd ../multitool && go get github.com/aws/aws-sdk-go-v2/config
	cd ../multitool && go get github.com/aws/aws-sdk-go-v2/service/cloudformation
	cd ../multitool && go get github.com/aws/aws-sdk-go-v2/credentials
	cd ../multitool && go get github.com/aws/aws-sdk-go-v2/service/s3
	cd ../multitool && go mod tidy
	cd ../terraform-multicloud-provider && go mod tidy || true
	cd .. && rm -rf vendor/
	cd .. && go work vendor

# Build werfty and other binaries
build: fix-deps
	@echo "Building multitool..."
	cd ../multitool && go build -mod=readonly -o mt ./
	@echo "Building cube-server..."
	cd ../cube-server && go build -mod=readonly -o cube-server ./
	@echo "Building werfty..."
	cd ../werfty && go build -mod=readonly -o werfty ./cmd

# Run all tests
test:
	cd .. && go test -mod=readonly ./...

# Test Azure functions in simulation mode (default)
test-azure: test-azure-sim

test-azure-sim:
	@echo "Testing Azure functions in simulation mode..."
	../scripts/test-azure-functions.sh

# Test Azure functions in direct mode
test-azure-direct:
	@echo "Testing Azure functions in direct mode..."
	../scripts/test-azure-functions.sh --direct

# Clean build artifacts
clean:
	rm -f ../werfty/werfty
	rm -f ../cube-server/cube-server
	rm -f ../multitool/mt
	rm -f ../generator/werfty-generator

# Start cube-server for testing
start-server:
	@echo "Starting cube-server..."
	cd ../cube-server && ./cube-server &

# Stop cube-server
stop-server:
	@echo "Stopping cube-server..."
	pkill -f cube-server || true

server-run:
	$(MAKE) -C ../server run

server-stop:
	$(MAKE) -C ../server stop

server-build:
	$(MAKE) -C ../server build

server-clean:
	$(MAKE) -C ../server clean

server-restart:
	$(MAKE) -C ../server stop
	$(MAKE) -C ../server run
