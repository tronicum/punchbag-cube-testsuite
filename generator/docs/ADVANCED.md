# Advanced Usage & Plugin System

## Plugin/Extensibility Framework (Planned)

- All resource generators and providers are modular.
- To add a new resource: implement a generator in `internal/generator/` and register it in the provider.
- To add a new provider: implement the provider interface and register in orchestration.
- Future: Support dynamic plugins via a `plugins/` directory and Go plugin system.

## CI/CD Integration

Example GitHub Actions workflow:

```yaml
name: CI
on: [push, pull_request]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - name: Install golangci-lint
        run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
      - name: Lint
        run: golangci-lint run
      - name: Test
        run: go test ./...
      - name: Terraform Validate
        run: |
          terraform init
          terraform validate
      - name: TFLint
        run: tflint
```

## Troubleshooting
- See `README.md` and `ARCHITECTURE.md` for usage and structure.
- For plugin development, see Go plugin docs and the `plugins/` directory (future).
