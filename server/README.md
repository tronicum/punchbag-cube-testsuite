# Server

The `server` component provides a REST API for managing Kubernetes clusters and running tests. It supports multi-cloud environments and includes simulation endpoints for cloud providers.

## Usage

1. Navigate to the `server` directory:
   ```bash
   cd server
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the server:
   ```bash
   go run main.go
   ```

The API will be available at `http://localhost:8080`.

## Features

- Multi-cloud support.
- Simulation endpoints for cloud providers.
- REST API for cluster management and testing.

## Overview

The server component provides a REST API for managing AKS clusters and running tests against them. It includes endpoints for cluster management, test execution, and monitoring.

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
