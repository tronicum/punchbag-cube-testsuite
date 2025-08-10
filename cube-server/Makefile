# cube-server Makefile: server orchestration (start/stop/build/clean)
.PHONY: build start stop clean run go-tests

build:
	go build -o cube-server ./

start:
	@echo "Starting cube-server..."
	./cube-server &

stop:
	@echo "Stopping cube-server..."
	pkill -f cube-server || true

clean:
	rm -f cube-server
	@echo "Cleaned cube-server binary."

run:
	./cube-server

go-tests:
	go test ./...

