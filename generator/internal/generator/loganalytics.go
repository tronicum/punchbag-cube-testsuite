package generator

import "fmt"

// GenerateLogAnalyticsTerraformBlock generates the Terraform block for a Log Analytics Workspace
func GenerateLogAnalyticsTerraformBlock(props map[string]interface{}, inputPath string) string {
	name := SafeString(props, "name", "example-log-analytics")
	location := SafeString(props, "location", "West Europe")
	resourceGroup := SafeString(props, "resourceGroup", "example-resource-group")
	sku := SafeString(props, "sku", "PerGB2018")
	retention := SafeInt(props, "retentionInDays", 30)
	return fmt.Sprintf("resource \"azurerm_log_analytics_workspace\" \"example\" {\n  name = \"%s\"\n  location = \"%s\"\n  resource_group_name = \"%s\"\n  sku = \"%s\"\n  retention_in_days = %d\n  // ...map more fields from JSON as needed\n}", name, location, resourceGroup, sku, retention)
}
