test-azure-sim:

# Root Makefile: orchestrates build, test, and clean only
.PHONY: build test clean

build:
	$(MAKE) -C multitool build
	$(MAKE) -C cube-server build
	$(MAKE) -C generator build
	$(MAKE) -C werfty build

test:
	$(MAKE) -C multitool go-tests
	$(MAKE) -C cube-server go-tests
	$(MAKE) -C generator go-tests
	$(MAKE) -C werfty go-tests

clean:
	$(MAKE) -C multitool clean
	$(MAKE) -C cube-server clean
	$(MAKE) -C generator clean
	$(MAKE) -C werfty clean
