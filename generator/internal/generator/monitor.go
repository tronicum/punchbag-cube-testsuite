package generator

import "fmt"

// GenerateMonitorTerraformBlock generates the Terraform block for an Azure Monitor Metric Alert
func GenerateMonitorTerraformBlock(props map[string]interface{}, inputPath string) string {
	name := SafeString(props, "name", "example-monitor")
	resourceGroup := SafeString(props, "resourceGroup", "example-rg")
	severity := SafeString(props, "severity", "3")
	criteria := SafeString(props, "criteria", "")
	return fmt.Sprintf("resource \"azurerm_monitor_metric_alert\" \"example\" {\n  name = \"%s\"\n  resource_group_name = \"%s\"\n  severity = %s\n  criteria = \"%s\"\n  // ...map more fields from JSON as needed\n}", name, resourceGroup, severity, criteria)
}
