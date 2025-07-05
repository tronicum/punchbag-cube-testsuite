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

## New Endpoints

### Validate Provider

Validate the configuration of a specific cloud provider.

**GET** `/api/v1/validate/:provider`

Example:
```bash
curl http://localhost:8080/api/v1/validate/aws
```

### Simulate Provider Operation

Simulate operations for a specific cloud provider.

**POST** `/api/v1/providers/:provider/operations/:operation`

Example:
```bash
curl -X POST http://localhost:8080/api/v1/providers/aws/operations/create-cluster
```

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

# Cube Server

The Cube Server provides a unified REST API for simulating and (optionally) executing cloud resource operations across multiple providers. It is the central backend for all simulation, validation, and resource orchestration in the Punchbag Cube ecosystem.

## Features
- **Simulation endpoints** for AKS clusters, Azure Budgets, and more (see `/api/v1/simulator/...`).
- **Provider validation** and info endpoints.
- **Pluggable execution endpoints** (future: `/executor`) for real cloud actions.
- **Backed by a shared Go simulation library** for consistent logic across all tools.
- **Stateless and easy to run locally or in CI.**

## Quick Start

1. **Build and run the server:**
   ```sh
   go build -o cube-server main.go
   ./cube-server
   # or simply: go run main.go
   ```

2. **Check health:**
   ```sh
   curl http://localhost:8080/health
   ```

3. **Simulate an AKS cluster:**
   ```sh
   curl -X POST http://localhost:8080/api/v1/simulator/azure/aks \
     -H 'Content-Type: application/json' \
     -d '{"name":"demo-aks","resource_group":"demo-rg","location":"eastus","node_count":3}'
   ```

4. **Simulate an Azure Budget:**
   ```sh
   curl -X POST http://localhost:8080/api/v1/simulator/azure/budget \
     -H 'Content-Type: application/json' \
     -d '{"name":"demo-budget","amount":1000,"period":"Monthly"}'
   ```

5. **Validate a provider:**
   ```sh
   curl -X POST http://localhost:8080/api/v1/providers/validate \
     -H 'Content-Type: application/json' \
     -d '{"provider":"azure"}'
   ```

## API Overview
- `/api/v1/simulator/azure/aks` - Simulate AKS cluster creation
- `/api/v1/simulator/azure/budget` - Simulate Azure Budget
- `/api/v1/providers/validate` - Validate provider
- `/api/v1/clusters` - CRUD for clusters (simulated)
- `/api/v1/metrics/health` - Health check

## Configuration
- No config required for simulation. For real execution, see `/executor` (future).

## Development
- All simulation logic is in `shared/simulation/service.go`.
- Add new providers or operations by extending the shared package.
- Run tests with:
  ```sh
  go test ./...
  ```

## License
MIT
