# Server

REST API server for the Punchbag Cube Test Suite.

## Overview

The server component provides a REST API for managing AKS clusters and running tests against them. It includes endpoints for cluster management, test execution, and monitoring.

## Features

- **Cluster Management**: Create, read, update, delete AKS clusters
- **Test Execution**: Run load tests, performance tests, and stress tests
- **Results Tracking**: Store and retrieve test results with detailed metrics
- **Health Monitoring**: Health checks and status endpoints
- **API Documentation**: Complete OpenAPI specification

## Getting Started

### Prerequisites
- Go 1.21 or later
- Docker (optional, for containerized deployment)

### Running the Server

```bash
cd server

# Install dependencies
go mod tidy

# Run the server
go run main.go
```

The API will be available at `http://localhost:8080`

### Environment Variables

- `PORT`: Server port (default: 8080)
- `GIN_MODE`: Gin mode (development/release)

## API Endpoints

### Health & Status
- `GET /health` - Health check
- `GET /docs` - API documentation

### Clusters
- `POST /api/v1/clusters` - Create cluster
- `GET /api/v1/clusters` - List clusters
- `GET /api/v1/clusters/:id` - Get cluster details
- `PUT /api/v1/clusters/:id` - Update cluster
- `DELETE /api/v1/clusters/:id` - Delete cluster

### Tests
- `POST /api/v1/clusters/:id/tests` - Run test on cluster
- `GET /api/v1/clusters/:id/tests` - List tests for cluster
- `GET /api/v1/tests/:id` - Get test result details

### Metrics
- `GET /api/v1/metrics/health` - Health check
- `GET /api/v1/metrics/status` - Service status

## Docker

```bash
# Build the container
docker build -t punchbag-cube-server .

# Run the container
docker run -p 8080:8080 punchbag-cube-server
```

## Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```
