package main

import (
	"os"
	"testing"
)

func TestGenerateTerraformFromJSON(t *testing.T) {
	input := `{"properties": {"name": "test-monitor"}}`
	inputFile := "test_monitor.json"
	outputFile := "test_monitor.tf"
	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)
	err := GenerateTerraformFromJSON(inputFile, outputFile)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}
	if len(content) == 0 {
		t.Error("Output Terraform file is empty")
	}
}

func TestGenerateAksTerraformBlock(t *testing.T) {
	props := map[string]interface{}{
		"name":            "test-aks",
		"location":        "eastus",
		"resourceGroup":  "test-rg",
		"nodeCount":      2,
	}
	tf := generateAksTerraformBlock(props)
	if len(tf) == 0 || !contains(tf, "azurerm_kubernetes_cluster") {
		t.Errorf("AKS Terraform block not generated correctly: %s", tf)
	}
}

func TestGenerateMonitorTerraformBlock(t *testing.T) {
	props := map[string]interface{}{
		"name":            "test-monitor",
		"resourceGroup":  "test-rg",
		"severity":        "2",
		"criteria":        "CPU > 90%",
	}
	tf := generateMonitorTerraformBlock(props, "monitor.json")
	if len(tf) == 0 || !contains(tf, "azurerm_monitor_metric_alert") {
		t.Errorf("Monitor Terraform block not generated correctly: %s", tf)
	}
}

func TestGenerateLogAnalyticsTerraformBlock(t *testing.T) {
	props := map[string]interface{}{
		"name":            "test-la",
		"location":        "westeurope",
		"resourceGroup":  "test-rg",
		"sku":             "PerGB2018",
	}
	tf := generateLogAnalyticsTerraformBlock(props, "loganalytics.json")
	if len(tf) == 0 || !contains(tf, "azurerm_log_analytics_workspace") {
		t.Errorf("Log Analytics Terraform block not generated correctly: %s", tf)
	}
}

func TestGenerateAppInsightsTerraformBlock(t *testing.T) {
	props := map[string]interface{}{
		"name":            "test-ai",
		"location":        "westeurope",
		"resourceGroup":  "test-rg",
		"applicationType": "web",
	}
	tf := generateAppInsightsTerraformBlock(props, "appinsights.json")
	if len(tf) == 0 || !contains(tf, "azurerm_application_insights") {
		t.Errorf("App Insights Terraform block not generated correctly: %s", tf)
	}
}

// contains is a helper for string containment
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && (s[0:len(substr)] == substr || contains(s[1:], substr))))
}
