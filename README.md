# Punchbag Cube Test Suite

A comprehensive test suite for testing punchbag cube functionality with server, client, and Terraform provider components.

## Overview

This project provides a complete ecosystem for testing various aspects of the punchbag cube system including:
- REST API server for cluster management and test execution
- Command-line client for interacting with the API
- Terraform provider for Infrastructure as Code (IaC) support
- AKS (Azure Kubernetes Service) integration testing
- Performance and load testing capabilities

## Project Structure

```
├── README.md                    # This file
├── LICENSE                      # License file
├── server/                      # REST API Server
│   ├── main.go                  # Server entry point
│   ├── go.mod                   # Server dependencies
│   ├── Dockerfile               # Server container config
│   ├── README.md                # Server documentation
│   ├── api/                     # API layer
│   │   ├── handlers.go          # HTTP request handlers
│   │   ├── routes.go            # Route definitions
│   │   └── openapi.yaml         # API specification
│   ├── models/                  # Data models
│   │   └── aks.go               # AKS-related models
│   └── store/                   # Data storage layer
│       └── store.go             # Storage interface
├── client/                      # CLI Client
│   ├── main.go                  # Client entry point
│   ├── go.mod                   # Client dependencies
│   ├── README.md                # Client documentation
│   ├── cmd/                     # CLI commands
│   │   ├── root.go              # Root command
│   │   ├── cluster.go           # Cluster commands
│   │   └── test.go              # Test commands
│   └── pkg/                     # Client packages
│       ├── api/                 # API client
│       │   └── client.go        # HTTP client
│       └── output/              # Output formatting
│           └── formatter.go     # Output formatters
└── terraform-provider/          # Terraform Provider
    ├── main.go                  # Provider entry point
    ├── go.mod                   # Provider dependencies
    ├── README.md                # Provider documentation
    └── internal/                # Provider implementation
        └── provider/            # Provider logic
            ├── provider.go      # Main provider
            ├── cluster_resource.go
            ├── test_resource.go
            └── clusters_data_source.go
```

## Components

### 🖥️ Server (`/server`)

REST API server that provides endpoints for:
- AKS cluster management (CRUD operations)
- Test execution and monitoring
- Health checks and metrics
- Complete OpenAPI specification

**Quick Start:**
```bash
cd server
go run main.go
```

Server will be available at `http://localhost:8080`

### 📱 Client (`/client`)

Command-line interface for interacting with the server:
- Manage clusters from the command line
- Run tests and monitor progress
- Multiple output formats (table, JSON, YAML)
- Configuration file support

**Quick Start:**
```bash
cd client
go build -o punchbag-client .
./punchbag-client cluster list
```

### 🏗️ Terraform Provider (`/terraform-provider`)

Terraform provider for Infrastructure as Code:
- Manage clusters with Terraform
- Run tests as part of infrastructure deployment
- Import existing resources
- Complete resource and data source support

**Quick Start:**
```bash
cd terraform-provider
go build -o terraform-provider-punchbag
```

## Getting Started

### Prerequisites
- Go 1.21 or later
- Docker (optional, for containerized deployment)
- Terraform 1.0+ (for using the provider)

### Quick Setup

1. **Start the Server:**
```bash
cd server
go mod tidy
go run main.go
```

2. **Use the Client:**
```bash
cd client
go mod tidy
go build -o punchbag-client .
./punchbag-client cluster list
```

3. **Try the Terraform Provider:**
```bash
cd terraform-provider
go mod tidy
go build -o terraform-provider-punchbag
# See terraform-provider/README.md for setup instructions
```

## Docker Support

Each component can be containerized:

```bash
# Server
cd server
docker build -t punchbag-server .
docker run -p 8080:8080 punchbag-server

# Client (for CI/CD pipelines)
cd client
docker build -t punchbag-client .
```

## API Documentation

The server provides comprehensive API documentation:
- Interactive docs: `http://localhost:8080/docs`
- OpenAPI spec: `server/api/openapi.yaml`

## Contributing

Each component has its own documentation and development guidelines. See the README.md files in each directory for specific instructions.

## License

This project is licensed under the terms specified in the LICENSE file.
