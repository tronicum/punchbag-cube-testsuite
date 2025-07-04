# Client

Command-line interface for the Punchbag Cube Test Suite.

## Overview

The client provides a CLI tool for interacting with the Punchbag Cube Test Suite API. It allows you to manage AKS clusters, run tests, and view results from the command line.

## Features

- **Cluster Management**: Create, list, get, update, and delete AKS clusters
- **Test Execution**: Run tests on clusters and monitor their progress
- **Multiple Output Formats**: Support for table, JSON, and YAML output
- **Configuration**: Support for configuration files and environment variables

## Installation

```bash
cd client

# Install dependencies
go mod tidy

# Build the client
go build -o punchbag-client .

# Install globally (optional)
go install .
```

## Usage

### Configuration

The client can be configured using:
- Command-line flags
- Configuration file (`$HOME/.punchbag-client.yaml`)
- Environment variables

Example configuration file:
```yaml
server: http://localhost:8080
format: table
verbose: false
```

### Basic Commands

```bash
# List all clusters
punchbag-client cluster list

# Get cluster details
punchbag-client cluster get <cluster-id>

# Create a new cluster
punchbag-client cluster create \
  --name my-cluster \
  --resource-group my-rg \
  --location eastus \
  --kubernetes-version 1.28.0 \
  --node-count 3

# Delete a cluster
punchbag-client cluster delete <cluster-id>

# Run a test on a cluster
punchbag-client cluster test <cluster-id> --type load_test

# Get test result
punchbag-client test get <test-id>

# List test results for a cluster
punchbag-client test list <cluster-id>

# Watch test progress
punchbag-client test watch <test-id>
```

### Output Formats

The client supports multiple output formats:

```bash
# Table format (default)
punchbag-client cluster list

# JSON format
punchbag-client cluster list --format json

# YAML format
punchbag-client cluster list --format yaml
```

### Global Flags

- `--server`: Server URL (default: http://localhost:8080)
- `--format`: Output format (table, json, yaml)
- `--verbose`: Verbose output
- `--config`: Configuration file path

## Examples

### Create and Test a Cluster

```bash
# Create a cluster
punchbag-client cluster create \
  --name test-cluster \
  --resource-group test-rg \
  --location eastus

# Run a load test
punchbag-client cluster test <cluster-id> \
  --type load_test \
  --config test-config.json

# Watch the test progress
punchbag-client test watch <test-id>
```

### Test Configuration

Create a `test-config.json` file for test parameters:

```json
{
  "duration": "5m",
  "concurrency": 10,
  "request_rate": 100,
  "target_url": "http://example.com",
  "method": "GET",
  "expected_code": 200
}
```
