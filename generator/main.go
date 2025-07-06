package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads the YAML configuration file
func LoadConfig(filePath string) (map[string]interface{}, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}
	if err := yaml.Unmarshal(content, &config); err != nil {
		return nil, err
	}

	return config, nil
}

func GenerateAzureTemplates(config map[string]interface{}) {
	// Example: Generate Azure Monitoring template
	if _, ok := config["azure_monitoring"].(map[string]interface{}); ok {
		fmt.Println("Generating Azure Monitoring template...")
		// Add logic to generate monitoring template
	}

	// Example: Generate Azure Kubernetes template
	if _, ok := config["azure_kubernetes"].(map[string]interface{}); ok {
		fmt.Println("Generating Azure Kubernetes template...")
		// Add logic to generate Kubernetes template
	}

	// Example: Generate Azure Budgets template
	if _, ok := config["azure_budgets"].(map[string]interface{}); ok {
		fmt.Println("Generating Azure Budgets template...")
		// Add logic to generate budgets template
	}

	// Example: Generate Azure Log Analytics template
	if _, ok := config["azure_log_analytics"].(map[string]interface{}); ok {
		fmt.Println("Generating Azure Log Analytics template...")
		// Add logic to generate Log Analytics template
	}
}

// Extend template generation for Azure services

func GenerateAzureMonitoringTemplate(config map[string]interface{}) string {
	// Stub logic for generating Azure Monitoring template
	return "# Azure Monitoring Template\nresource \"azurerm_monitoring\" \"example\" {}"
}

func GenerateAzureKubernetesTemplate(config map[string]interface{}) string {
	// Stub logic for generating Azure Kubernetes template
	return "# Azure Kubernetes Template\nresource \"azurerm_kubernetes_cluster\" \"example\" {}"
}

func GenerateAzureBudgetTemplate(config map[string]interface{}) string {
	// Stub logic for generating Azure Budget template
	return "# Azure Budget Template\nresource \"azurerm_budget\" \"example\" {}"
}

// GenerateAzureLogAnalyticsTemplate generates a Terraform template for Azure Log Analytics
func GenerateAzureLogAnalyticsTemplate(config map[string]interface{}) string {
	return `
# Azure Log Analytics Template
resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-log-analytics"
  location            = "West Europe"
  resource_group_name = "example-resource-group"
  sku                 = "PerGB2018"
  retention_in_days   = 30
}
`
}

// GenerateTerraformFromJSON reads a JSON file and outputs a Terraform file for supported Azure resources
func GenerateTerraformFromJSON(inputPath, outputPath string) error {
	content, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input JSON: %w", err)
	}
	var data map[string]interface{}
	if err := json.Unmarshal(content, &data); err != nil {
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
		if _, hasNodeCount := props["nodeCount"]; hasNodeCount && strings.Contains(strings.ToLower(safeString(props, "name", "")), "aks") {
			tf = generateAksTerraformBlock(props)
		} else if strings.Contains(inputPath, "monitor") {
			// Map common Azure Monitor fields
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
			tf = fmt.Sprintf(`resource "azurerm_monitor_metric_alert" "example" {
  name                = "%s"
  resource_group_name = "%s"
  severity            = %s
  criteria            = "%s"
%s%s%s%s%s%s%s%s%s%s  // ...map more fields from JSON as needed
}`,
				name, resourceGroup, severity, criteria, descriptionLine, enabledLine, scopes, evalFreqLine, windowSizeLine, disabledLine, autoMitigateLine, alertRuleResourceIdLine, criteriaBlock, "")
		} else if strings.Contains(inputPath, "loganalytics") {
			// Map Log Analytics fields (expanded)
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
			tf = fmt.Sprintf(`resource "azurerm_log_analytics_workspace" "example" {
  name                = "%s"
  location            = "%s"
  resource_group_name = "%s"
  sku                 = "%s"
  retention_in_days   = %d
%s%s%s%s%s%s%s%s%s%s%s  // ...map more fields from JSON as needed
}`,
				name, location, resourceGroup, sku, retention, customerIdLine, cappingBlock, internetIngestionLine, internetQueryLine, reservationCapLine, dailyQuotaLine, publicNetworkAccessLine, workspaceIdLine, primarySharedKeyLine, tagsBlock, "")
		} else if strings.Contains(inputPath, "appinsights") {
			// Map Application Insights fields
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
			tf = fmt.Sprintf(`resource "azurerm_application_insights" "example" {
  name                = "%s"
  location            = "%s"
  resource_group_name = "%s"
  application_type    = "%s"
  retention_in_days   = %d
%s%s%s%s  // ...map more fields from JSON as needed
}`,
				name, location, resourceGroup, appType, retention, workspaceIdLine, dailyCapLine, disableIpMaskingLine, tagsBlock)
		} else {
			return fmt.Errorf("unsupported or unrecognized resource type in %s", inputPath)
		}
	} else {
		return fmt.Errorf("unsupported or unrecognized resource type in %s", inputPath)
	}
	return os.WriteFile(outputPath, []byte(tfHeader+tf), 0644)
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
%s%s%s  // ...map more fields from JSON as needed
}`,
		name, location, resourceGroup, nodeCount, identity, nodeCount, networkPlugin, networkPolicy, dnsPrefix, tagsLine, zones, labels)
}

func main() {
	generateTerraform := flag.Bool("generate-terraform", false, "Generate Terraform from JSON input")
	input := flag.String("input", "", "Input JSON file (from multitool)")
	output := flag.String("output", "", "Output Terraform file")
	simulate := flag.Bool("simulate-import", false, "Simulate resource import (mock)")
	resourceType := flag.String("resource-type", "", "Resource type: monitor|loganalytics|aks|budget (mock)")
	name := flag.String("name", "", "Resource name (mock)")
	location := flag.String("location", "", "Resource location (mock)")
	resourceGroup := flag.String("resource-group", "", "Resource group (mock)")
	nodeCount := flag.Int("node-count", 3, "Node count (AKS, mock)")
	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `werfty - Azure resource simulation and Terraform codegen

USAGE:
  werfty [flags]

FLAGS:
  --simulate-import         Simulate resource import (mock, outputs JSON)
  --generate-terraform      Generate Terraform from JSON input
  --input <file>            Input JSON file (from multitool)
  --output <file>           Output Terraform file
  --resource-type <type>    Resource type: monitor|loganalytics|aks|budget (mock)
  --name <name>             Resource name (mock)
  --location <location>     Resource location (mock)
  --resource-group <group>  Resource group (mock)
  --node-count <n>          Node count (AKS, mock)
  -h, --help                Show this help

EXAMPLES:
  Simulate AKS:
    werfty --simulate-import --resource-type aks --name my-aks --resource-group my-rg --location eastus --node-count 3

  Generate Terraform from JSON:
    werfty --generate-terraform --input test_aks_expanded.json --output test_aks_expanded.tf

  Import existing AKS to Terraform state:
    terraform import azurerm_kubernetes_cluster.my_aks /subscriptions/${AZURE_SUBSCRIPTION_ID}/resourceGroups/my-rg/providers/Microsoft.ContainerService/managedClusters/my-aks
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
			fmt.Println("Usage: werfty --generate-terraform --input <input.json> --output <output.tf>")
			os.Exit(1)
		}
		err := GenerateTerraformFromJSON(*input, *output)
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

	// Parse the Terraform template and generate Go code as comments only
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

	config, err := LoadConfig("../conf/punchy.yml")
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	GenerateAzureTemplates(config)

	// Example: Use SimulatorClient to get AKS simulation result
	// client := SimulatorClient{BaseURL: "http://localhost:8080"}
	// aksParams := map[string]interface{}{
	// 	"name":            "example-aks",
	// 	"resource_group":  "example-rg",
	// 	"location":        "eastus",
	// 	"node_count":      3,
	// }
	// aksResult, err := client.SimulateAKSCluster(aksParams)
	// if err != nil {
	// 	fmt.Printf("AKS simulation failed: %v\n", err)
	// } else {
	// 	fmt.Printf("Simulated AKS result: %+v\n", aksResult)
	// 	// Use aksResult.Result to generate Terraform template
	// }
}

// execLookPath is a wrapper for exec.LookPath
func execLookPath(file string) (string, error) {
	return exec.LookPath(file)
}

// execCommand is a wrapper for exec.Command
func execCommand(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
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
