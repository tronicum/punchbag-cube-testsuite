package generator

import "fmt"

// GenerateAksTerraformBlock generates the Terraform block for an AKS cluster from properties
func GenerateAksTerraformBlock(props map[string]interface{}) string {
	name := SafeString(props, "name", "example-aks")
	location := SafeString(props, "location", "eastus")
	resourceGroup := SafeString(props, "resourceGroup", "example-rg")
	nodeCount := SafeInt(props, "nodeCount", 3)
	// ...existing code for zones, labels, tagsLine, etc...
	return fmt.Sprintf("resource \"azurerm_kubernetes_cluster\" \"example\" {\n  name = \"%s\"\n  location = \"%s\"\n  resource_group_name = \"%s\"\n  default_node_pool {\n    name = \"default\"\n    node_count = %d\n  }\n  // ...map more fields from JSON as needed\n}", name, location, resourceGroup, nodeCount)
}
