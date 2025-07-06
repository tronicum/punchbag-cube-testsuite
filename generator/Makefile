# Makefile for generator CLI and end-to-end testing

# Use the new tflint CLI style for v0.47+
BINARY=werfty-generator

all: build

build:
	go build -o $(BINARY) main.go

clean:
	rm -f $(BINARY) *.tf *.json *.a

# Run all Go tests (unit and mapping logic)
test:
	go test -v ./...

lint:
	golangci-lint run

fmt:
	gofmt -s -w .

vet:
	go vet ./...

deploy:
	./deploy.sh

# End-to-end test: simulate multitool, generate Terraform, lint
end2end: clean build
	chmod +x $(BINARY)
	./$(BINARY) --simulate-import --resource-type aks --name test-aks --location eastus --resource-group test-rg --node-count 2 > test_aks.json
	./$(BINARY) --generate-terraform --input test_aks.json --output test_aks.tf
	@if command -v tflint >/dev/null 2>&1; then \
		echo "Running tflint on test_aks.tf..."; \
		mkdir -p .tflint-tmp && cp test_aks.tf .tflint-tmp/main.tf && cd .tflint-tmp && tflint --chdir=. --filter=main.tf && cd .. && rm -rf .tflint-tmp; \
	else \
		echo "tflint not found, skipping lint."; \
	fi

.PHONY: all build clean test lint fmt vet deploy end2end
