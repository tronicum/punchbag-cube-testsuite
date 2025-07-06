package main

import "fmt"

// generateMonitorTerraformBlock generates the Terraform block for an Azure Monitor Metric Alert
func generateMonitorTerraformBlock(props map[string]interface{}, inputPath string) string {
	name := safeString(props, "name", "example-monitor")
	resourceGroup := safeString(props, "resourceGroup", "example-rg")
	severity := safeString(props, "severity", "3")
	criteria := safeString(props, "criteria", "")
	enabled := safeBool(props, "enabled", true)
	description := safeString(props, "description", "")
	scopes := ""
	if s, ok := props["scopes"].([]interface{}); ok && len(s) > 0 {
		scopes = "  scopes = ["
		for i, v := range s {
			if i > 0 {
				scopes += ", "
			}
			scopes += fmt.Sprintf("\"%s\"", v)
		}
		scopes += "]\n"
	}
	descriptionLine := ""
	if description != "" {
		descriptionLine = fmt.Sprintf("  description = \"%s\"\n", description)
	}
	enabledLine := ""
	if !enabled {
		enabledLine = "  enabled = false\n"
	}
	// Additional Monitor fields
	evaluationFrequency := safeString(props, "evaluationFrequency", "")
	windowSize := safeString(props, "windowSize", "")
	disabled := safeBool(props, "disabled", false)
	autoMitigate := safeBool(props, "autoMitigate", true)
	alertRuleResourceId := safeString(props, "alertRuleResourceId", "")
	criteriaBlock := ""
	if c, ok := props["criteriaBlock"].(string); ok && c != "" {
		criteriaBlock = fmt.Sprintf("  criteria_block = \"%s\"\n", c)
	}
	evalFreqLine := ""
	if evaluationFrequency != "" {
		evalFreqLine = fmt.Sprintf("  evaluation_frequency = \"%s\"\n", evaluationFrequency)
	}
	windowSizeLine := ""
	if windowSize != "" {
		windowSizeLine = fmt.Sprintf("  window_size = \"%s\"\n", windowSize)
	}
	disabledLine := ""
	if disabled {
		disabledLine = "  enabled = false\n"
	}
	autoMitigateLine := ""
	if !autoMitigate {
		autoMitigateLine = "  auto_mitigate = false\n"
	}
	alertRuleResourceIdLine := ""
	if alertRuleResourceId != "" {
		alertRuleResourceIdLine = fmt.Sprintf("  alert_rule_resource_id = \"%s\"\n", alertRuleResourceId)
	}
	// TODO: Map more Monitor fields as needed
	// Example: frequency, window size, actions, etc.
	actions := safeString(props, "actions", "")
	actionsLine := ""
	if actions != "" {
		actionsLine = fmt.Sprintf("  actions = [%s]\n", actions)
	}
	return fmt.Sprintf(`resource "azurerm_monitor_metric_alert" "example" {
  name                = "%s"
  resource_group_name = "%s"
  severity            = %s
  criteria            = "%s"
%s%s%s%s%s%s%s%s%s%s%s  // ...map more fields from JSON as needed
}`,
		name, resourceGroup, severity, criteria, descriptionLine, enabledLine, scopes, evalFreqLine, windowSizeLine, disabledLine, autoMitigateLine, alertRuleResourceIdLine, criteriaBlock, actionsLine, "")
}
