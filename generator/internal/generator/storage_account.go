package generator

import "fmt"

// GenerateStorageAccountTerraformBlock generates Terraform for an Azure Storage Account
func GenerateStorageAccountTerraformBlock(props map[string]interface{}) string {
	name := SafeString(props, "name", "example-storage")
	resourceGroup := SafeString(props, "resourceGroup", "example-rg")
	location := SafeString(props, "location", "eastus")
	sku := SafeString(props, "sku", "Standard_LRS")
	kind := SafeString(props, "kind", "StorageV2")
	accessTier := SafeString(props, "accessTier", "Hot")

	return fmt.Sprintf(`resource "azurerm_storage_account" "%s" {
  name                     = "%s"
  resource_group_name      = "%s"
  location                 = "%s"
  account_tier             = "%s"
  account_replication_type = "%s"
  account_kind             = "%s"
  access_tier              = "%s"
}
`, name, name, resourceGroup, location, "Standard", sku, kind, accessTier)
}
