package main

import "fmt"

// generateAppInsightsTerraformBlock generates the Terraform block for Application Insights
func generateAppInsightsTerraformBlock(props map[string]interface{}, inputPath string) string {
	name := safeString(props, "name", "example-appinsights")
	location := safeString(props, "location", "West Europe")
	resourceGroup := safeString(props, "resourceGroup", "example-resource-group")
	appType := safeString(props, "applicationType", "web")
	retention := safeInt(props, "retentionInDays", 90)
	workspaceId := safeString(props, "workspaceId", "")
	dailyCap := safeInt(props, "dailyDataCapInGb", 0)
	disableIpMasking := safeBool(props, "disableIpMasking", false)
	tagsBlock := ""
	if t, ok := props["tags"].(map[string]interface{}); ok && len(t) > 0 {
		tagsBlock = "  tags = {\n"
		for k, v := range t {
			tagsBlock += fmt.Sprintf("    %q = %q\n", k, v)
		}
		tagsBlock += "  }\n"
	}
	workspaceIdLine := ""
	if workspaceId != "" {
		workspaceIdLine = fmt.Sprintf("  workspace_id = \"%s\"\n", workspaceId)
	}
	dailyCapLine := ""
	if dailyCap > 0 {
		dailyCapLine = fmt.Sprintf("  daily_data_cap_in_gb = %d\n", dailyCap)
	}
	disableIpMaskingLine := ""
	if disableIpMasking {
		disableIpMaskingLine = "  disable_ip_masking = true\n"
	}
	// TODO: Map more App Insights fields as needed
	// Example: application_id, ingestion_mode, etc.
	appId := safeString(props, "applicationId", "")
	appIdLine := ""
	if appId != "" {
		appIdLine = fmt.Sprintf("  application_id = \"%s\"\n", appId)
	}
	ingestionMode := safeString(props, "ingestionMode", "")
	ingestionModeLine := ""
	if ingestionMode != "" {
		ingestionModeLine = fmt.Sprintf("  ingestion_mode = \"%s\"\n", ingestionMode)
	}
	return fmt.Sprintf(`resource "azurerm_application_insights" "example" {
  name                = "%s"
  location            = "%s"
  resource_group_name = "%s"
  application_type    = "%s"
  retention_in_days   = %d
%s%s%s%s%s%s  // ...map more fields from JSON as needed
}`,
		name, location, resourceGroup, appType, retention, workspaceIdLine, dailyCapLine, disableIpMaskingLine, tagsBlock, appIdLine, ingestionModeLine)
}
