# Root Makefile to control the server from the project root

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
