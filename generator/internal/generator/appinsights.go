package generator

import "fmt"

// GenerateAppInsightsTerraformBlock generates the Terraform block for Application Insights
func GenerateAppInsightsTerraformBlock(props map[string]interface{}, inputPath string) string {
	name := SafeString(props, "name", "example-appinsights")
	location := SafeString(props, "location", "West Europe")
	resourceGroup := SafeString(props, "resourceGroup", "example-resource-group")
	appType := SafeString(props, "applicationType", "web")
	retention := SafeInt(props, "retentionInDays", 90)
	return fmt.Sprintf("resource \"azurerm_application_insights\" \"example\" {\n  name = \"%s\"\n  location = \"%s\"\n  resource_group_name = \"%s\"\n  application_type = \"%s\"\n  retention_in_days = %d\n  // ...map more fields from JSON as needed\n}", name, location, resourceGroup, appType, retention)
}
