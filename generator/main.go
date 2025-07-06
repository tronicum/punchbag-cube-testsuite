package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// GenerateTerraformFromJSON reads a JSON file and outputs a Terraform file for supported Azure resources
// TODO: Support multiple resources in a single input file (e.g., array of resources or map of resource blocks)
// For now, only a single resource under 'properties' is supported.
func GenerateTerraformFromJSON(inputPath, outputPath string) error {
	content, err := os.ReadFile(inputPath)
	if err != nil {
		log.Printf("%s Failed to read input JSON file %s: %v", logPrefix("GenerateTerraformFromJSON"), inputPath, err)
		return fmt.Errorf("failed to read input JSON: %w", err)
	}
	var data map[string]interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		log.Printf("%s Invalid JSON in %s: %v", logPrefix("GenerateTerraformFromJSON"), inputPath, err)
		return fmt.Errorf("invalid JSON: %w", err)
	}
	// Terraform required blocks
	tfHeader := `terraform {
  required_version = ">= 1.0.0"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 3.0.0"
    }
  }
}

provider "azurerm" {
  features {}
}
`
	// Detect resource type by keys and map fields
	var tf string
	if props, ok := data["properties"].(map[string]interface{}); ok {
		// Robust resource type detection
		var resourceType string
		if _, hasNodeCount := props["nodeCount"]; hasNodeCount && strings.Contains(strings.ToLower(safeString(props, "name", "")), "aks") {
			resourceType = "aks"
		} else if strings.Contains(inputPath, "monitor") {
			resourceType = "monitor"
		} else if strings.Contains(inputPath, "loganalytics") {
			resourceType = "loganalytics"
		} else if strings.Contains(inputPath, "appinsights") {
			resourceType = "appinsights"
		} else {
			log.Printf("%s Unsupported or unrecognized resource type in %s", logPrefix("GenerateTerraformFromJSON"), inputPath)
			return fmt.Errorf("unsupported or unrecognized resource type in %s", inputPath)
		}
		// Schema validation for Azure
		if err := validateResourceProperties("azure", resourceType, props); err != nil {
			log.Printf("%s Schema validation failed: %v", logPrefix("GenerateTerraformFromJSON"), err)
			return fmt.Errorf("schema validation failed: %w", err)
		}
		if resourceType == "aks" {
			tf = generateAksTerraformBlock(props)
		} else if resourceType == "monitor" {
			tf = generateMonitorTerraformBlock(props, inputPath)
		} else if resourceType == "loganalytics" {
			tf = generateLogAnalyticsTerraformBlock(props, inputPath)
		} else if resourceType == "appinsights" {
			tf = generateAppInsightsTerraformBlock(props, inputPath)
		}
	} else {
		log.Printf("%s No 'properties' key found or not a map in %s", logPrefix("GenerateTerraformFromJSON"), inputPath)
		return fmt.Errorf("unsupported or unrecognized resource type in %s", inputPath)
	}
	if err := os.WriteFile(outputPath, []byte(tfHeader+tf), 0644); err != nil {
		log.Printf("%s Failed to write Terraform output to %s: %v", logPrefix("GenerateTerraformFromJSON"), outputPath, err)
		return err
	}
	log.Printf("%s Terraform code written to %s", logPrefix("GenerateTerraformFromJSON"), outputPath)
	return nil
}

// --- Multicloud stubs ---
// generateEksTerraformBlock generates the Terraform block for an AWS EKS cluster (stub)
func generateEksTerraformBlock(props map[string]interface{}) string {
	name := safeString(props, "name", "example-eks")
	region := safeString(props, "region", "us-west-2")
	nodeCount := safeInt(props, "nodeCount", 2)
	version := safeString(props, "eksVersion", "1.29")
	instanceType := safeString(props, "instanceType", "t3.medium")
	subnetIds := []string{}
	if s, ok := props["subnetIds"].([]interface{}); ok {
		for _, v := range s {
			subnetIds = append(subnetIds, fmt.Sprintf("\"%s\"", v))
		}
	}
	subnetLine := ""
	if len(subnetIds) > 0 {
		subnetLine = fmt.Sprintf("  subnet_ids = [%s]\n", strings.Join(subnetIds, ", "))
	}
	return fmt.Sprintf(`resource "aws_eks_cluster" "example" {
  name     = "%s"
  region   = "%s"
  version  = "%s"
  vpc_config {
    subnet_ids = [%s]
  }
  node_group {
    instance_types = ["%s"]
    desired_capacity = %d
  }
  %s// ...map more fields from JSON as needed
}`,
		name, region, version, strings.Join(subnetIds, ", "), instanceType, nodeCount, subnetLine)
}

// generateGkeTerraformBlock generates the Terraform block for a GCP GKE cluster (stub)
func generateGkeTerraformBlock(props map[string]interface{}) string {
	name := safeString(props, "name", "example-gke")
	location := safeString(props, "location", "us-central1")
	nodeCount := safeInt(props, "nodeCount", 2)
	// TODO: Map more GKE fields as needed
	return fmt.Sprintf(`resource "google_container_cluster" "example" {
  name     = "%s"
  location = "%s"
  initial_node_count = %d
  // ...map more fields from JSON as needed
}`,
		name, location, nodeCount)
}

// GenerateTerraformFromJSONMulticloud is a multicloud-ready version of GenerateTerraformFromJSON
func GenerateTerraformFromJSONMulticloud(inputPath, outputPath, provider string) error {
	content, err := os.ReadFile(inputPath)
	if err != nil {
		log.Printf("%s Failed to read input JSON file %s: %v", logPrefix("GenerateTerraformFromJSONMulticloud"), inputPath, err)
		return fmt.Errorf("failed to read input JSON: %w", err)
	}
	var data map[string]interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		log.Printf("%s Invalid JSON in %s: %v", logPrefix("GenerateTerraformFromJSONMulticloud"), inputPath, err)
		return fmt.Errorf("invalid JSON: %w", err)
	}
	// Terraform required blocks (provider-specific)
	var tfHeader string
	switch provider {
	case "aws":
		tfHeader = `terraform {\n  required_version = ">= 1.0.0"\n  required_providers {\n    aws = {\n      source  = "hashicorp/aws"\n      version = ">= 4.0.0"\n    }\n  }\n}\n\nprovider "aws" {\n  region = "us-west-2"\n}\n`
	case "gcp":
		tfHeader = `terraform {\n  required_version = ">= 1.0.0"\n  required_providers {\n    google = {\n      source  = "hashicorp/google"\n      version = ">= 4.0.0"\n    }\n  }\n}\n\nprovider "google" {\n}\n`
	case "azure":
		tfHeader = `terraform {\n  required_version = ">= 1.0.0"\n  required_providers {\n    azurerm = {\n      source  = "hashicorp/azurerm"\n      version = ">= 3.0.0"\n    }\n  }\n}\n\nprovider "azurerm" {\n  features {}\n}\n`
	default:
		return fmt.Errorf("unsupported or unrecognized provider: %s", provider)
	}
	// Detect resource type by keys and map fields
	var tf string
	if props, ok := data["properties"].(map[string]interface{}); ok {
		var resourceType string
		switch provider {
		case "aws":
			if _, hasNodeCount := props["nodeCount"]; hasNodeCount {
				resourceType = "eks"
				err := validateResourceProperties(provider, resourceType, props)
				if err != nil {
					return err
				}
				tf = generateEksTerraformBlock(props)
			} else if _, hasBucket := props["bucket"] ; hasBucket || safeString(props, "resourceType", "") == "s3" || strings.Contains(strings.ToLower(safeString(props, "name", "")), "s3") {
				resourceType = "s3"
				err := validateResourceProperties(provider, resourceType, props)
				if err != nil {
					return err
				}
				tf = generateS3TerraformBlock(props)
			} else {
				return fmt.Errorf("unsupported or unrecognized AWS resource type in %s", inputPath)
			}
		case "gcp":
			if _, hasNodeCount := props["nodeCount"]; hasNodeCount {
				resourceType = "gke"
				err := validateResourceProperties(provider, resourceType, props)
				if err != nil {
					return err
				}
				tf = generateGkeTerraformBlock(props)
			} else {
				return fmt.Errorf("unsupported or unrecognized GCP resource type in %s", inputPath)
			}
		case "azure":
			if _, hasNodeCount := props["nodeCount"]; hasNodeCount && strings.Contains(strings.ToLower(safeString(props, "name", "")), "aks") {
				resourceType = "aks"
				err := validateResourceProperties(provider, resourceType, props)
				if err != nil {
					return err
				}
				tf = generateAksTerraformBlock(props)
			} else if strings.Contains(inputPath, "monitor") {
				resourceType = "monitor"
				err := validateResourceProperties(provider, resourceType, props)
				if err != nil {
					return err
				}
				tf = generateMonitorTerraformBlock(props, inputPath)
			} else if strings.Contains(inputPath, "loganalytics") {
				resourceType = "loganalytics"
				err := validateResourceProperties(provider, resourceType, props)
				if err != nil {
					return err
				}
				tf = generateLogAnalyticsTerraformBlock(props, inputPath)
			} else if strings.Contains(inputPath, "appinsights") {
				resourceType = "appinsights"
				err := validateResourceProperties(provider, resourceType, props)
				if err != nil {
					return err
				}
				tf = generateAppInsightsTerraformBlock(props, inputPath)
			} else {
				return fmt.Errorf("unsupported or unrecognized Azure resource type in %s", inputPath)
			}
		}
	} else {
		return fmt.Errorf("unsupported or unrecognized resource type in %s", inputPath)
	}
	if err := os.WriteFile(outputPath, []byte(tfHeader+tf), 0644); err != nil {
		log.Printf("%s Failed to write Terraform output to %s: %v", logPrefix("GenerateTerraformFromJSONMulticloud"), outputPath, err)
		return err
	}
	log.Printf("%s Terraform code written to %s", logPrefix("GenerateTerraformFromJSONMulticloud"), outputPath)
	return nil
}

// safeString returns a string value from a map or a default
func safeString(m map[string]interface{}, key, def string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return def
}

// safeInt returns an int value from a map or a default
func safeInt(m map[string]interface{}, key string, def int) int {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case int:
			return val
		case float64:
			return int(val)
		}
	}
	return def
}

// safeBool returns a bool value from a map or a default
func safeBool(m map[string]interface{}, key string, def bool) bool {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case bool:
			return val
		case string:
			return val == "true" || val == "1"
		}
	}
	return def
}

// generateAksTerraformBlock generates the Terraform block for an AKS cluster from properties
func generateAksTerraformBlock(props map[string]interface{}) string {
	name := safeString(props, "name", "example-aks")
	location := safeString(props, "location", "eastus")
	resourceGroup := safeString(props, "resourceGroup", "example-rg")
	nodeCount := safeInt(props, "nodeCount", 3)
	networkPlugin := safeString(props, "networkPlugin", "azure")
	networkPolicy := safeString(props, "networkPolicy", "azure")
	dnsPrefix := safeString(props, "dnsPrefix", "exampleaks")
	identity := safeString(props, "identity", "")
	tags := safeString(props, "tags", "")
	// AKS availability zones (if applicable)
	var zones string
	if z, ok := props["availabilityZones"].([]interface{}); ok && len(z) > 0 {
		zones = "  availability_zones = ["
		for i, v := range z {
			if i > 0 {
				zones += ", "
			}
			zones += fmt.Sprintf("\"%s\"", v)
		}
		zones += "]\n"
	}
	// AKS node pool labels and taints
	labels := safeString(props, "nodePoolLabels", "")
	tagsLine := ""
	if tags != "" {
		tagsLine = fmt.Sprintf("  tags = %s\n", tags)
	}
	// TODO: Map additional AKS fields as needed
	// Example: Enable RBAC, API server authorized IP ranges, etc.
	enableRBAC := safeBool(props, "enableRBAC", true)
	rbacLine := ""
	if enableRBAC {
		rbacLine = "  role_based_access_control {\n    enabled = true\n  }\n"
	}
	apiServerAuthorizedIPRanges := safeString(props, "apiServerAuthorizedIPRanges", "")
	apiServerIPLine := ""
	if apiServerAuthorizedIPRanges != "" {
		apiServerIPLine = fmt.Sprintf("  api_server_authorized_ip_ranges = [%s]\n", apiServerAuthorizedIPRanges)
	}
	return fmt.Sprintf(`resource "azurerm_kubernetes_cluster" "example" {
  name                = "%s"
  location            = "%s"
  resource_group_name = "%s"
  default_node_pool {
    name       = "default"
    node_count = %d
  }
  identity {
    type = "%s"
  }
  sku {
    name     = "Standard_DS2_v2"
    tier     = "Standard"
    capacity = %d
  }
  network_profile {
    network_plugin = "%s"
    network_policy = "%s"
  }
  dns_prefix          = "%s"
%s%s%s%s%s  // ...map more fields from JSON as needed
}`,
		name, location, resourceGroup, nodeCount, identity, nodeCount, networkPlugin, networkPolicy, dnsPrefix, tagsLine, zones, labels, rbacLine, apiServerIPLine)
}

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

// generateCloudWatchTerraformBlock generates the Terraform block for AWS CloudWatch Metric Alarm
func generateCloudWatchTerraformBlock(props map[string]interface{}) string {
	name := safeString(props, "name", "example-alarm")
	namespace := safeString(props, "namespace", "AWS/EC2")
	metricName := safeString(props, "metricName", "CPUUtilization")
	comparison := safeString(props, "comparisonOperator", "GreaterThanThreshold")
	threshold := safeInt(props, "threshold", 80)
	period := safeInt(props, "period", 300)
	evaluationPeriods := safeInt(props, "evaluationPeriods", 1)
	statistic := safeString(props, "statistic", "Average")
	alarmActions := []string{}
	if a, ok := props["alarmActions"].([]interface{}); ok {
		for _, v := range a {
			alarmActions = append(alarmActions, fmt.Sprintf("\"%s\"", v))
		}
	}
	alarmActionsLine := ""
	if len(alarmActions) > 0 {
		alarmActionsLine = fmt.Sprintf("  alarm_actions = [%s]\n", strings.Join(alarmActions, ", "))
	}
	return fmt.Sprintf(`resource "aws_cloudwatch_metric_alarm" "example" {
  alarm_name          = "%s"
  namespace           = "%s"
  metric_name         = "%s"
  comparison_operator = "%s"
  threshold           = %d
  period              = %d
  evaluation_periods  = %d
  statistic           = "%s"
%s  // ...map more fields from JSON as needed
}`,
		name, namespace, metricName, comparison, threshold, period, evaluationPeriods, statistic, alarmActionsLine)
}

// generateCloudWatchLogGroupTerraformBlock generates the Terraform block for AWS CloudWatch Log Group
func generateCloudWatchLogGroupTerraformBlock(props map[string]interface{}) string {
	name := safeString(props, "name", "example-log-group")
	retention := safeInt(props, "retentionInDays", 14)
	return fmt.Sprintf(`resource "aws_cloudwatch_log_group" "example" {
  name              = "%s"
  retention_in_days = %d
  // ...map more fields from JSON as needed
}`,
		name, retention)
}

// generateAwsBudgetTerraformBlock generates the Terraform block for AWS Budgets
func generateAwsBudgetTerraformBlock(props map[string]interface{}) string {
	name := safeString(props, "name", "example-budget")
	amount := safeInt(props, "amount", 1000)
	period := safeString(props, "period", "MONTHLY")
	return fmt.Sprintf(`resource "aws_budgets_budget" "example" {
  name              = "%s"
  budget_type       = "COST"
  limit_amount      = "%d"
  limit_unit        = "USD"
  time_period_start = "2023-01-01_00:00"
  time_unit         = "%s"
  // ...map more fields from JSON as needed
}`,
		name, amount, period)
}

// generateS3TerraformBlock generates the Terraform block for an AWS S3 bucket
func generateS3TerraformBlock(props map[string]interface{}) string {
	name := safeString(props, "name", "example-s3-bucket")
	acl := safeString(props, "acl", "private")
	versioning := safeBool(props, "versioning", false)
	versioningBlock := ""
	if versioning {
		versioningBlock = "  versioning {\n    enabled = true\n  }\n"
	}
	return fmt.Sprintf(`resource "aws_s3_bucket" "example" {
  bucket = "%s"
  acl    = "%s"
%s  // ...map more fields from JSON as needed
}`,
		name, acl, versioningBlock)
}

// validateResourceProperties checks required fields for each resource/provider
func validateResourceProperties(provider, resourceType string, props map[string]interface{}) error {
	required := []string{}
	switch provider {
	case "azure":
		switch resourceType {
		case "aks":
			required = []string{"name", "location", "resourceGroup", "nodeCount"}
		case "monitor":
			required = []string{"name", "resourceGroup", "severity", "criteria"}
		case "loganalytics":
			required = []string{"name", "location", "resourceGroup", "sku", "retentionInDays"}
		case "appinsights":
			required = []string{"name", "location", "resourceGroup", "applicationType"}
		}
	case "aws":
		switch resourceType {
		case "eks":
			required = []string{"name", "region", "nodeCount"}
		case "s3":
			required = []string{"name"}
		case "cloudwatch":
			required = []string{"name", "namespace", "metricName", "comparisonOperator", "threshold", "period", "evaluationPeriods", "statistic"}
		case "budget":
			required = []string{"name", "amount", "period"}
		}
	case "gcp":
		switch resourceType {
		case "gke":
			required = []string{"name", "location", "nodeCount"}
		}
	}
	missing := []string{}
	for _, k := range required {
		if _, ok := props[k]; !ok {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required fields for %s/%s: %v", provider, resourceType, missing)
	}
	if len(required) == 0 {
		return fmt.Errorf("unknown provider or resource type: %s/%s", provider, resourceType)
	}
	missing = []string{}
	for _, k := range required {
		if _, ok := props[k]; !ok {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required fields for %s/%s: %v", provider, resourceType, missing)
	}
	return nil
}

func main() {
	generateTerraform := flag.Bool("generate-terraform", false, "Generate Terraform from JSON input")
	input := flag.String("input", "", "Input JSON file (from multitool)")
	output := flag.String("output", "", "Output Terraform file")
	simulate := flag.Bool("simulate-import", false, "Simulate resource import (mock)")
	resourceType := flag.String("resource-type", "", "Resource type: monitor|loganalytics|aks|budget (mock)")
	provider := flag.String("provider", "azure", "Cloud provider: azure|aws|gcp") // NEW
	name := flag.String("name", "", "Resource name (mock)")
	location := flag.String("location", "", "Resource location (mock)")
	// region flag for AWS/GCP (alias for location)
	region := flag.String("region", "", "AWS/GCP region (alias for --location for AWS/GCP)")
	resourceGroup := flag.String("resource-group", "", "Resource group (mock)")
	nodeCount := flag.Int("node-count", 3, "Node count (AKS/EKS/GKE, mock)")
	flag.Parse()

	if *location == "" && *region != "" {
		*location = *region
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `werfty - Multicloud resource simulation and Terraform codegen

USAGE:
  werfty [flags]

FLAGS:
  --simulate-import         Simulate resource import (mock, outputs JSON)
  --generate-terraform      Generate Terraform from JSON input (multicloud)
  --input <file>            Input JSON file (from multitool or manual)
  --output <file>           Output Terraform file
  --provider <provider>     Cloud provider: azure|aws|gcp (required for multicloud)
  --resource-type <type>    Resource type: monitor|loganalytics|aks|eks|gke|budget (mock)
  --name <name>             Resource name (mock/simulation)
  --location <location>     Resource location (Azure/GCP) or region (AWS)
  --region <region>         AWS/GCP region (alias for --location for AWS/GCP)
  --resource-group <group>  Azure resource group (mock/simulation)
  --node-count <n>          Node count (AKS/EKS/GKE, mock)
  -h, --help                Show this help

EXAMPLES:
  Simulate AKS (Azure):
    werfty --simulate-import --resource-type aks --name my-aks --resource-group my-rg --location eastus --node-count 3

  Simulate EKS (AWS):
    werfty --simulate-import --resource-type eks --name my-eks --location us-west-2 --node-count 2

  Simulate GKE (GCP):
    werfty --simulate-import --resource-type gke --name my-gke --location us-central1 --node-count 1

  Generate Terraform for Azure:
    werfty --generate-terraform --input test_aks.json --output test_aks.tf --provider azure

  Generate Terraform for AWS:
    werfty --generate-terraform --input test_eks.json --output test_eks.tf --provider aws

  Generate Terraform for GCP:
    werfty --generate-terraform --input test_gke.json --output test_gke.tf --provider gcp

  Import existing AKS to Terraform state:
    terraform import azurerm_kubernetes_cluster.my_aks /subscriptions/${AZURE_SUBSCRIPTION_ID}/resourceGroups/my-rg/providers/Microsoft.ContainerService/managedClusters/my-aks

NOTES:
- --provider is required for multicloud workflows (azure, aws, gcp)
- --location is used for Azure/GCP, --region for AWS/GCP (interchangeable for GCP)
- --resource-type is only for simulation, not for real codegen
- Use the multitool CLI for end-to-end automation and multicloud workflows
`)
	}

	if *simulate {
		// Output a mock JSON for the requested resource type
		mock := map[string]interface{}{"properties": map[string]interface{}{}}
		switch *resourceType {
		case "monitor":
			mock["properties"].(map[string]interface{})["name"] = *name
			mock["properties"].(map[string]interface{})["resourceGroup"] = *resourceGroup
			mock["properties"].(map[string]interface{})["severity"] = "3"
			mock["properties"].(map[string]interface{})["criteria"] = "CPU > 80%"
		case "loganalytics":
			mock["properties"].(map[string]interface{})["name"] = *name
			mock["properties"].(map[string]interface{})["location"] = *location
			mock["properties"].(map[string]interface{})["resourceGroup"] = *resourceGroup
			mock["properties"].(map[string]interface{})["sku"] = "PerGB2018"
			mock["properties"].(map[string]interface{})["retentionInDays"] = 30
		case "aks":
			mock["properties"].(map[string]interface{})["name"] = *name
			mock["properties"].(map[string]interface{})["location"] = *location
			mock["properties"].(map[string]interface{})["resourceGroup"] = *resourceGroup
			mock["properties"].(map[string]interface{})["nodeCount"] = *nodeCount
		case "eks": // AWS EKS stub
			mock["properties"].(map[string]interface{})["name"] = *name
			mock["properties"].(map[string]interface{})["region"] = *location
			mock["properties"].(map[string]interface{})["nodeCount"] = *nodeCount
		case "gke": // GCP GKE stub
			mock["properties"].(map[string]interface{})["name"] = *name
			mock["properties"].(map[string]interface{})["location"] = *location
			mock["properties"].(map[string]interface{})["nodeCount"] = *nodeCount
		case "budget":
			mock["properties"].(map[string]interface{})["name"] = *name
			mock["properties"].(map[string]interface{})["amount"] = 1000
			mock["properties"].(map[string]interface{})["timeGrain"] = "Monthly"
			mock["properties"].(map[string]interface{})["resourceGroup"] = *resourceGroup
		}
		out, _ := json.MarshalIndent(mock, "", "  ")
		fmt.Println(string(out))
		return
	}

	if *generateTerraform {
		if *input == "" || *output == "" {
			fmt.Println("Usage: werfty --generate-terraform --input <input.json> --output <output.tf> --provider <provider>")
			os.Exit(1)
		}
		if *provider != "azure" && *provider != "aws" && *provider != "gcp" {
			fmt.Printf("Error: --provider must be one of: azure, aws, gcp (got '%s')\n", *provider)
			os.Exit(1)
		}
		err := GenerateTerraformFromJSONMulticloud(*input, *output, *provider)
		if err != nil {
			fmt.Printf("Terraform generation failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Terraform code written to %s\n", *output)
		// Lint the generated Terraform file if tflint is available
		if _, err := os.Stat(*output); err == nil {
			if _, err := execLookPath("tflint"); err == nil {
				fmt.Println("Running tflint on generated Terraform...")
				cmd := execCommand("tflint", *output)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				_ = cmd.Run()
			} else {
				fmt.Println("tflint not found; skipping linting.")
			}
		}
		return
	}

	// Read the Terraform template file
	filePath := "azure_services.tf"
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		os.Exit(1)
	}

	lines := strings.Split(string(content), "\n")
	generatedCode := "package azure\n\n// This file is auto-generated for resource reference only.\n// The following resources were found in azure_services.tf:\n\n"

	for _, line := range lines {
		if strings.HasPrefix(line, "resource ") {
			resourceType := strings.Split(line, " ")[1]
			resourceName := strings.Split(line, " ")[2]
			generatedCode += fmt.Sprintf("// Resource: %s %s\n", resourceType, resourceName)
		}
	}

	// Write the generated Go code to a file
	outputFilePath := "generated_resources.go"
	err = os.WriteFile(outputFilePath, []byte(generatedCode), 0644)
	if err != nil {
		fmt.Printf("Failed to write generated code: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Code generation completed. Check generated_resources.go")
}

// execLookPath is a wrapper for exec.LookPath
func execLookPath(file string) (string, error) {
	return exec.LookPath(file)
}

// execCommand is a wrapper for exec.Command
func execCommand(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}

// logPrefix returns a formatted prefix for log messages
func logPrefix(context string) string {
	return fmt.Sprintf("[werfty][%s]", context)
}

// ---
//
// # Example: Importing AKS Clusters from an Existing Azure Account (Simulated)
//
// 1. Set your Azure API credentials as environment variables:
//
//    export AZURE_CLIENT_ID=...
//    export AZURE_CLIENT_SECRET=...
//    export AZURE_TENANT_ID=...
//    export AZURE_SUBSCRIPTION_ID=...
//
// 2. Simulate AKS cluster import (no real cloud call):
//
//    go run generator/main.go --simulate-import --name my-aks --resource-group my-rg --location eastus --node-count 3
//
//    # This will call the cube-server simulation endpoint and print a simulated AKS cluster result.
//    # Use this to generate a Terraform template.
//
// 3. To generate a real Terraform template for deployment:
//
//    go run generator/main.go --generate-terraform --name my-aks --resource-group my-rg --location eastus --node-count 3
//
//    # Then apply with Terraform:
//    terraform init
//    terraform plan
//    terraform apply
//
// ---
//
// # Example: Importing Existing Resources into Terraform State
//
// 1. Use the `terraform import` command to bring an existing AKS cluster into your Terraform state:
//
//    terraform import azurerm_kubernetes_cluster.my_aks \
//      /subscriptions/${AZURE_SUBSCRIPTION_ID}/resourceGroups/my-rg/providers/Microsoft.ContainerService/managedClusters/my-aks
//
// 2. After import, run `terraform plan` to see any drift between your state and the actual resource.
//    - If your Terraform code does not match the real resource, Terraform will show a plan to change it.
//    - You may need to update your .tf files to match the imported resource's configuration.
//
// 3. The import only updates the state file; it does not update your .tf files. Always review and update your configuration after import.
//
// ---
//
// # Where to Put API Credentials
//
// - Environment variables (recommended)
// - .env file (if supported)
// - Directly in a config YAML (for local testing only)
//
// ---
//
// # Switching Between Simulation and Real Execution
//
// - Use --simulate-import for simulation only (no cloud calls)
// - Use --generate-terraform for real Terraform code
// - Use /executor endpoint (when implemented) for real cloud actions, with --dryrun or --force
//
// ---
//
// # Example CLI Usage
//
// Simulate AKS:
//   go run generator/main.go --simulate-import --name my-aks --resource-group my-rg --location eastus --node-count 3
//
// Generate Terraform:
//   go run generator/main.go --generate-terraform --name my-aks --resource-group my-rg --location eastus --node-count 3
//
// Import existing AKS:
//   terraform import azurerm_kubernetes_cluster.my_aks /subscriptions/${AZURE_SUBSCRIPTION_ID}/resourceGroups/my-rg/providers/Microsoft.ContainerService/managedClusters/my-aks
//
// ---
// # Supported Azure Resources and Field Mapping
//
// - Azure Monitor Metric Alert: name, resource_group_name, severity, criteria
// - Azure Log Analytics Workspace: name, location, resource_group_name, sku, retention_in_days
// - Azure Kubernetes Cluster: name, location, resource_group_name, node_count (default pool)
//
// More fields and resource types can be added as needed. See GenerateTerraformFromJSON for mapping logic.
// ---
