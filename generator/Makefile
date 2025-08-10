# Makefile for generator CLI and end-to-end testing

BINARY=werfty-generator

all: build

build:
	go build -o $(BINARY) main.go

clean:
	rm -f $(BINARY) *.tf *.json *.a

# Run all Go tests (unit and mapping logic)
go-tests:
	go test -v ./...

# End-to-end test: simulate multitool, generate Terraform, lint
end2end: clean build
	chmod +x $(BINARY)
	./$(BINARY) --simulate-import --resource-type aks --name test-aks --location eastus --resource-group test-rg --node-count 2 > test_aks.json
	./$(BINARY) --generate-terraform --input test_aks.json --output test_aks.tf
	@if command -v tflint >/dev/null 2>&1; then \
		echo "Running tflint on test_aks.tf..."; \
		tflint test_aks.tf; \
	else \
		echo "tflint not found, skipping lint."; \
	fi

.PHONY: all build clean test end2end
