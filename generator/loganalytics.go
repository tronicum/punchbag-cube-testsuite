package main

import "fmt"

// generateLogAnalyticsTerraformBlock generates the Terraform block for a Log Analytics Workspace
func generateLogAnalyticsTerraformBlock(props map[string]interface{}, inputPath string) string {
	name := safeString(props, "name", "example-log-analytics")
	location := safeString(props, "location", "West Europe")
	resourceGroup := safeString(props, "resourceGroup", "example-resource-group")
	sku := safeString(props, "sku", "PerGB2018")
	retention := safeInt(props, "retentionInDays", 30)
	customerId := safeString(props, "customerId", "")
	workspaceCapping := safeInt(props, "workspaceCapping", 0)
	internetIngestion := safeBool(props, "internetIngestionEnabled", false)
	internetQuery := safeBool(props, "internetQueryEnabled", false)
	reservationCap := safeInt(props, "reservationCapacityInGbPerDay", 0)
	dailyQuota := safeInt(props, "dailyQuotaGb", 0)
	publicNetworkAccess := safeString(props, "publicNetworkAccessForIngestion", "")
	workspaceId := safeString(props, "workspaceId", "")
	primarySharedKey := safeString(props, "primarySharedKey", "")
	tagsBlock := ""
	if t, ok := props["tags"].(map[string]interface{}); ok && len(t) > 0 {
		tagsBlock = "  tags = {\n"
		for k, v := range t {
			tagsBlock += fmt.Sprintf("    %q = %q\n", k, v)
		}
		tagsBlock += "  }\n"
	}
	cappingBlock := ""
	if workspaceCapping > 0 {
		cappingBlock = fmt.Sprintf("  workspace_capping {\n    daily_quota_gb = %d\n  }\n", workspaceCapping)
	}
	customerIdLine := ""
	if customerId != "" {
		customerIdLine = fmt.Sprintf("  customer_id = \"%s\"\n", customerId)
	}
	internetIngestionLine := ""
	if internetIngestion {
		internetIngestionLine = "  internet_ingestion_enabled = true\n"
	}
	internetQueryLine := ""
	if internetQuery {
		internetQueryLine = "  internet_query_enabled = true\n"
	}
	reservationCapLine := ""
	if reservationCap > 0 {
		reservationCapLine = fmt.Sprintf("  reservation_capacity_in_gb_per_day = %d\n", reservationCap)
	}
	dailyQuotaLine := ""
	if dailyQuota > 0 {
		dailyQuotaLine = fmt.Sprintf("  daily_quota_gb = %d\n", dailyQuota)
	}
	publicNetworkAccessLine := ""
	if publicNetworkAccess != "" {
		publicNetworkAccessLine = fmt.Sprintf("  public_network_access_for_ingestion = \"%s\"\n", publicNetworkAccess)
	}
	workspaceIdLine := ""
	if workspaceId != "" {
		workspaceIdLine = fmt.Sprintf("  workspace_id = \"%s\"\n", workspaceId)
	}
	primarySharedKeyLine := ""
	if primarySharedKey != "" {
		primarySharedKeyLine = fmt.Sprintf("  primary_shared_key = \"%s\"\n", primarySharedKey)
	}
	// TODO: Map more Log Analytics fields as needed
	// Example: linked services, data sources, etc.
	linkedService := safeString(props, "linkedService", "")
	linkedServiceLine := ""
	if linkedService != "" {
		linkedServiceLine = fmt.Sprintf("  linked_service = \"%s\"\n", linkedService)
	}
	return fmt.Sprintf(`resource "azurerm_log_analytics_workspace" "example" {
  name                = "%s"
  location            = "%s"
  resource_group_name = "%s"
  sku                 = "%s"
  retention_in_days   = %d
%s%s%s%s%s%s%s%s%s%s%s%s  // ...map more fields from JSON as needed
}`,
		name, location, resourceGroup, sku, retention, customerIdLine, cappingBlock, internetIngestionLine, internetQueryLine, reservationCapLine, dailyQuotaLine, publicNetworkAccessLine, workspaceIdLine, primarySharedKeyLine, tagsBlock, linkedServiceLine, "")
}
