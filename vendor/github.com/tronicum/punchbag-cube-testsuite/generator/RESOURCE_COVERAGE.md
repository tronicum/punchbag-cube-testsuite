# Azure Resource Terraform Generation Coverage

This generator currently supports the following Azure resources and field mappings when converting JSON to Terraform:

## Supported Resources & Fields

### Azure Monitor Metric Alert
- `name` (string)
- `resourceGroup` (string)
- `severity` (string/int)
- `criteria` (string)

### Azure Log Analytics Workspace
- `name` (string)
- `location` (string)
- `resourceGroup` (string)
- `sku` (string)
- `retentionInDays` (int)

### Azure Kubernetes Cluster (AKS)
- `name` (string)
- `location` (string)
- `resourceGroup` (string)
- `nodeCount` (int, for default node pool)

## Example JSON Input

```
{
  "properties": {
    "name": "my-aks",
    "location": "eastus",
    "resourceGroup": "my-rg",
    "nodeCount": 3
  }
}
```

## Example CLI Usage

```
go run main.go --generate-terraform --input input.json --output output.tf
```

## Extending Coverage
- To add more fields or resource types, update `GenerateTerraformFromJSON` in `main.go`.
- Use the `safeString` and `safeInt` helpers for robust field extraction.

---
For more details, see comments in `main.go`.
