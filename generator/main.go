package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
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
	if monitoringConfig, ok := config["azure_monitoring"].(map[string]interface{}); ok {
		fmt.Println("Generating Azure Monitoring template...")
		// Add logic to generate monitoring template
	}

	// Example: Generate Azure Kubernetes template
	if kubernetesConfig, ok := config["azure_kubernetes"].(map[string]interface{}); ok {
		fmt.Println("Generating Azure Kubernetes template...")
		// Add logic to generate Kubernetes template
	}

	// Example: Generate Azure Budgets template
	if budgetsConfig, ok := config["azure_budgets"].(map[string]interface{}); ok {
		fmt.Println("Generating Azure Budgets template...")
		// Add logic to generate budgets template
	}

	// Example: Generate Azure Log Analytics template
	if logAnalyticsConfig, ok := config["azure_log_analytics"].(map[string]interface{}); ok {
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

func main() {
	// Read the Terraform template file
	filePath := "azure_services.tf"
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		os.Exit(1)
	}

	// Parse the Terraform template and generate Go code
	lines := strings.Split(string(content), "\n")
	generatedCode := "package azure\n\nimport (\n\t\"fmt\"\n)\n\nfunc CreateAzureResources() {\n"

	for _, line := range lines {
		if strings.HasPrefix(line, "resource \"") {
			resourceType := strings.Split(line, " ")[1]
			resourceName := strings.Split(line, " ")[2]
			generatedCode += fmt.Sprintf("\tfmt.Println(\"Creating %s: %s\")\n", resourceType, resourceName)
		}
	}

	generatedCode += "}\n"

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
}
