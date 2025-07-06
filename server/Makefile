# Simple Makefile for server

run:
	DEBUG=true LOGFILE=server.log go run main.go

run-prod:
	LOGFILE=server.log go run main.go

build:
	go build -o server main.go

clean:
	rm -f server

stop:
	@PID=$$(lsof -ti :8080); if [ "$$PID" != "" ]; then kill $$PID; echo "Stopped server on port 8080 (PID: $$PID)"; else echo "No server running on port 8080"; fi
