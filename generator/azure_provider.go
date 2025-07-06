package main

import (
	"fmt"
	"strings"
)

// AzureProvider implements Provider for Azure resources

type AzureProvider struct{}

func (a AzureProvider) GenerateTerraform(props map[string]interface{}, inputPath string) (string, string, error) {
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
	resourceType, err := a.DetectResourceType(props, inputPath)
	if err != nil {
		return "", "", err
	}
	if err := a.ValidateResource(resourceType, props); err != nil {
		return "", "", err
	}
	var tf string
	switch resourceType {
	case "aks":
		tf = generateAksTerraformBlock(props)
	case "monitor":
		tf = generateMonitorTerraformBlock(props, inputPath)
	case "loganalytics":
		tf = generateLogAnalyticsTerraformBlock(props, inputPath)
	case "appinsights":
		tf = generateAppInsightsTerraformBlock(props, inputPath)
	default:
		return "", "", fmt.Errorf("unsupported or unrecognized Azure resource type: %s", resourceType)
	}
	return tfHeader, tf, nil
}

func (a AzureProvider) ValidateResource(resourceType string, props map[string]interface{}) error {
	return ValidateResourceProperties("azure", resourceType, props)
}

func (a AzureProvider) DetectResourceType(props map[string]interface{}, inputPath string) (string, error) {
	if _, hasNodeCount := props["nodeCount"]; hasNodeCount && strings.Contains(strings.ToLower(safeString(props, "name", "")), "aks") {
		return "aks", nil
	} else if strings.Contains(inputPath, "monitor") {
		return "monitor", nil
	} else if strings.Contains(inputPath, "loganalytics") {
		return "loganalytics", nil
	} else if strings.Contains(inputPath, "appinsights") {
		return "appinsights", nil
	}
	return "", fmt.Errorf("unsupported or unrecognized resource type in %s", inputPath)
}
