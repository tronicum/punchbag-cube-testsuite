package main

import (
	"os"
	"strings"
	"testing"
)

func TestGenerateTerraformFromJSON(t *testing.T) {
	input := `{"properties": {"name": "test-monitor", "resourceGroup": "test-rg", "severity": "2", "criteria": "CPU > 90%"}}`
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

func TestGenerateEksTerraformBlock(t *testing.T) {
	props := map[string]interface{}{
		"name":          "test-eks",
		"region":        "us-west-2",
		"nodeCount":     2,
		"eksVersion":    "1.29",
		"instanceType":  "t3.medium",
		"subnetIds":     []interface{}{"subnet-123", "subnet-456"},
	}
	tf := generateEksTerraformBlock(props)
	if len(tf) == 0 || !contains(tf, "aws_eks_cluster") {
		t.Errorf("EKS Terraform block not generated correctly: %s", tf)
	}
	if !contains(tf, "subnet-123") || !contains(tf, "subnet-456") {
		t.Errorf("EKS subnet IDs not mapped: %s", tf)
	}
}

func TestGenerateCloudWatchTerraformBlock(t *testing.T) {
	props := map[string]interface{}{
		"name":                 "test-alarm",
		"namespace":            "AWS/EC2",
		"metricName":           "CPUUtilization",
		"comparisonOperator":   "GreaterThanThreshold",
		"threshold":            80,
		"period":               300,
		"evaluationPeriods":    2,
		"statistic":            "Average",
		"alarmActions":         []interface{}{"arn:aws:sns:us-west-2:123456789012:my-sns"},
	}
	tf := generateCloudWatchTerraformBlock(props)
	if len(tf) == 0 || !contains(tf, "aws_cloudwatch_metric_alarm") {
		t.Errorf("CloudWatch alarm block not generated correctly: %s", tf)
	}
	if !contains(tf, "my-sns") {
		t.Errorf("CloudWatch alarm actions not mapped: %s", tf)
	}
}

func TestGenerateCloudWatchLogGroupTerraformBlock(t *testing.T) {
	props := map[string]interface{}{
		"name":            "test-log-group",
		"retentionInDays": 7,
	}
	tf := generateCloudWatchLogGroupTerraformBlock(props)
	if len(tf) == 0 || !contains(tf, "aws_cloudwatch_log_group") {
		t.Errorf("CloudWatch log group block not generated correctly: %s", tf)
	}
}

func TestGenerateAwsBudgetTerraformBlock(t *testing.T) {
	props := map[string]interface{}{
		"name":    "test-budget",
		"amount":  500,
		"period":  "MONTHLY",
	}
	tf := generateAwsBudgetTerraformBlock(props)
	if len(tf) == 0 || !contains(tf, "aws_budgets_budget") {
		t.Errorf("AWS budget block not generated correctly: %s", tf)
	}
}

func TestGenerateTerraformFromJSON_EdgeCases(t *testing.T) {
	// Missing properties
	input := `{}`
	inputFile := "test_empty.json"
	outputFile := "test_empty.tf"
	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)
	err := GenerateTerraformFromJSON(inputFile, outputFile)
	if err == nil {
		t.Error("Expected error for missing properties, got nil")
	}

	// Unknown resource type
	input = `{"properties": {"foo": "bar"}}`
	inputFile = "test_unknown.json"
	outputFile = "test_unknown.tf"
	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)
	err = GenerateTerraformFromJSON(inputFile, outputFile)
	if err == nil {
		t.Error("Expected error for unknown resource type, got nil")
	}
}

func TestGenerateTerraformFromJSONMulticloud_AWS(t *testing.T) {
	input := `{"properties": {"name": "test-eks", "region": "us-west-2", "nodeCount": 2}}`
	inputFile := "test_eks.json"
	outputFile := "test_eks.tf"
	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)
	err := GenerateTerraformFromJSONMulticloud(inputFile, outputFile, "aws")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}
	if !contains(string(content), "aws_eks_cluster") {
		t.Error("AWS EKS Terraform not generated")
	}
}

func TestGenerateTerraformFromJSONMulticloud_GCP(t *testing.T) {
	input := `{"properties": {"name": "test-gke", "location": "us-central1", "nodeCount": 1}}`
	inputFile := "test_gke.json"
	outputFile := "test_gke.tf"
	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)
	err := GenerateTerraformFromJSONMulticloud(inputFile, outputFile, "gcp")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}
	if !contains(string(content), "google_container_cluster") {
		t.Error("GCP GKE Terraform not generated")
	}
}

func TestGenerateTerraformFromJSONMulticloud_UnknownProvider(t *testing.T) {
	input := `{"properties": {"name": "test-aks", "nodeCount": 1}}`
	inputFile := "test_aks.json"
	outputFile := "test_aks.tf"
	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)
	err := GenerateTerraformFromJSONMulticloud(inputFile, outputFile, "unknown")
	if err == nil {
		t.Error("Expected error for unknown provider, got nil")
	}
}

func TestValidateResourceProperties_MissingFields(t *testing.T) {
	cases := []struct {
		provider     string
		resourceType string
		props        map[string]interface{}
		shouldFail   bool
	}{
		{"azure", "aks", map[string]interface{}{"name": "n", "location": "l"}, true}, // missing resourceGroup, nodeCount
		{"aws", "eks", map[string]interface{}{"name": "n"}, true}, // missing region, nodeCount
		{"gcp", "gke", map[string]interface{}{"name": "n", "location": "l", "nodeCount": 1}, false}, // valid
		{"aws", "s3", map[string]interface{}{}, true}, // missing name
		{"azure", "monitor", map[string]interface{}{"name": "n", "resourceGroup": "g"}, true}, // missing severity, criteria
	}
	for _, c := range cases {
		err := validateResourceProperties(c.provider, c.resourceType, c.props)
		if c.shouldFail && err == nil {
			t.Errorf("Expected failure for %s/%s, got nil", c.provider, c.resourceType)
		}
		if !c.shouldFail && err != nil {
			t.Errorf("Expected success for %s/%s, got error: %v", c.provider, c.resourceType, err)
		}
	}
}

func TestValidateResourceProperties_InvalidProviderOrType(t *testing.T) {
	err := validateResourceProperties("unknown", "aks", map[string]interface{}{"name": "n"})
	if err == nil || !contains(err.Error(), "unknown provider or resource type") {
		t.Errorf("Expected unknown provider/type error, got: %v", err)
	}
	err = validateResourceProperties("azure", "unknown", map[string]interface{}{"name": "n"})
	if err == nil || !contains(err.Error(), "unknown provider or resource type") {
		t.Errorf("Expected unknown provider/type error, got: %v", err)
	}
}

// contains is a helper for string containment
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
