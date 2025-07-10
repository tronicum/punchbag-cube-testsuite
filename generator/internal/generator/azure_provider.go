package generator

import (
	"fmt"
	"strings"
)

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
		tf = GenerateAksTerraformBlock(props)
	case "monitor":
		tf = GenerateMonitorTerraformBlock(props, inputPath)
	case "loganalytics":
		tf = GenerateLogAnalyticsTerraformBlock(props, inputPath)
	case "appinsights":
		tf = GenerateAppInsightsTerraformBlock(props, inputPath)
	case "storageaccount":
		tf = GenerateStorageAccountTerraformBlock(props)
	default:
		return "", "", fmt.Errorf("unsupported or unrecognized Azure resource type: %s", resourceType)
	}
	return tfHeader, tf, nil
}

func (a AzureProvider) ValidateResource(resourceType string, props map[string]interface{}) error {
	return ValidateResourceProperties("azure", resourceType, props)
}

func (a AzureProvider) DetectResourceType(props map[string]interface{}, inputPath string) (string, error) {
	if _, hasNodeCount := props["nodeCount"]; hasNodeCount && strings.Contains(strings.ToLower(SafeString(props, "name", "")), "aks") {
		return "aks", nil
	} else if strings.Contains(inputPath, "monitor") {
		return "monitor", nil
	} else if strings.Contains(inputPath, "loganalytics") {
		return "loganalytics", nil
	} else if strings.Contains(inputPath, "appinsights") {
		return "appinsights", nil
	} else if strings.Contains(inputPath, "storage") || SafeString(props, "resourceType", "") == "storageaccount" {
		return "storageaccount", nil
	}
	return "", fmt.Errorf("unsupported or unrecognized resource type in %s", inputPath)
}
