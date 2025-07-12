# Root Makefile to control the server from the project root

.PHONY: help build test test-azure test-azure-sim test-azure-direct clean

# Default target
help:
	@echo "Available targets:"
	@echo "  build           Build all binaries"
	@echo "  test           Run all tests"
	@echo "  test-azure     Test Azure functions (simulation mode)"
	@echo "  test-azure-sim Test Azure functions (simulation mode)"
	@echo "  test-azure-direct Test Azure functions (direct mode)"
	@echo "  clean          Clean build artifacts"

# Build werfty and other binaries
build:
	@echo "Building werfty..."
	cd werfty && go build -o werfty ./cmd
	@echo "Building cube-server..."
	cd cube-server && go build -o cube-server ./
	@echo "Building multitool..."
	cd multitool && go build -o mt ./

# Run all tests
test:
	go test ./...

# Test Azure functions in simulation mode (default)
test-azure: test-azure-sim

test-azure-sim:
	@echo "Testing Azure functions in simulation mode..."
	./scripts/test-azure-functions.sh

# Test Azure functions in direct mode
test-azure-direct:
	@echo "Testing Azure functions in direct mode..."
	./scripts/test-azure-functions.sh --direct

# Clean build artifacts
clean:
	rm -f werfty/werfty
	rm -f cube-server/cube-server
	rm -f multitool/mt
	rm -f generator/werfty-generator

# Start cube-server for testing
start-server:
	@echo "Starting cube-server..."
	cd cube-server && ./cube-server &

# Stop cube-server
stop-server:
	@echo "Stopping cube-server..."
	pkill -f cube-server || true

server-run:
	$(MAKE) -C server run

server-stop:
	$(MAKE) -C server stop

server-build:
	$(MAKE) -C server build

server-clean:
	$(MAKE) -C server clean

server-restart:
	$(MAKE) -C server stop
	$(MAKE) -C server run
