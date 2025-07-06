# Azure Storage Account Example

This example demonstrates how to define an Azure Storage Account in YAML for Werfty:

```yaml
resourceType: storageaccount
properties:
  name: mystorageacct123
  resourceGroup: my-rg
  location: eastus2
  sku: Standard_LRS
  kind: StorageV2
  accessTier: Hot
```

To generate Terraform:

```sh
go run main.go generate -i examples/example_azure_storage_account.yaml -o storage_account.tf -p azure
```
