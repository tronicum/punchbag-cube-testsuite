package generator

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// GenerateTerraformFromJSON reads a YAML/JSON file and outputs a Terraform file for supported Azure resources
func GenerateTerraformFromJSON(inputPath, outputPath string) error {
	data, err := LoadConfigFromFile(inputPath)
	if err != nil {
		log.Printf("%s Failed to load config file %s: %v", LogPrefix("GenerateTerraformFromJSON"), inputPath, err)
		return fmt.Errorf("failed to load config: %w", err)
	}
	props, ok := data["properties"].(map[string]interface{})
	if !ok {
		log.Printf("%s No 'properties' key found or not a map in %s", LogPrefix("GenerateTerraformFromJSON"), inputPath)
		return fmt.Errorf("unsupported or unrecognized resource type in %s", inputPath)
	}
	provider := AzureProvider{}
	tfHeader, tf, err := provider.GenerateTerraform(props, inputPath)
	if err != nil {
		log.Printf("%s %v", LogPrefix("GenerateTerraformFromJSON"), err)
		return err
	}
	if err := os.WriteFile(outputPath, []byte(tfHeader+tf), 0644); err != nil {
		log.Printf("%s Failed to write Terraform output to %s: %v", LogPrefix("GenerateTerraformFromJSON"), outputPath, err)
		return err
	}
	log.Printf("%s Terraform code written to %s", LogPrefix("GenerateTerraformFromJSON"), outputPath)
	return nil
}

// GenerateTerraformFromJSONMulticloud is a multicloud-ready version of GenerateTerraformFromJSON
func GenerateTerraformFromJSONMulticloud(inputPath, outputPath, provider string) error {
	data, err := LoadConfigFromFile(inputPath)
	if err != nil {
		log.Printf("%s Failed to load config file %s: %v", LogPrefix("GenerateTerraformFromJSONMulticloud"), inputPath, err)
		return fmt.Errorf("failed to load config: %w", err)
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
				err := ValidateResourceProperties(provider, resourceType, props)
				if err != nil {
					return err
				}
				tf = GenerateEksTerraformBlock(props)
			} else if _, hasBucket := props["bucket"] ; hasBucket || SafeString(props, "resourceType", "") == "s3" || strings.Contains(strings.ToLower(SafeString(props, "name", "")), "s3") {
				resourceType = "s3"
				err := ValidateResourceProperties(provider, resourceType, props)
				if err != nil {
					return err
				}
				tf = GenerateS3TerraformBlock(props)
			} else {
				return fmt.Errorf("unsupported or unrecognized AWS resource type in %s", inputPath)
			}
		case "gcp":
			if _, hasNodeCount := props["nodeCount"]; hasNodeCount {
				resourceType = "gke"
				err := ValidateResourceProperties(provider, resourceType, props)
				if err != nil {
					return err
				}
				tf = GenerateGkeTerraformBlock(props)
			} else {
				return fmt.Errorf("unsupported or unrecognized GCP resource type in %s", inputPath)
			}
		case "azure":
			if _, hasNodeCount := props["nodeCount"]; hasNodeCount && strings.Contains(strings.ToLower(SafeString(props, "name", "")), "aks") {
				resourceType = "aks"
				err := ValidateResourceProperties(provider, resourceType, props)
				if err != nil {
					return err
				}
				tf = GenerateAksTerraformBlock(props)
			} else if strings.Contains(inputPath, "monitor") {
				resourceType = "monitor"
				err := ValidateResourceProperties(provider, resourceType, props)
				if err != nil {
					return err
				}
				tf = GenerateMonitorTerraformBlock(props, inputPath)
			} else if strings.Contains(inputPath, "loganalytics") {
				resourceType = "loganalytics"
				err := ValidateResourceProperties(provider, resourceType, props)
				if err != nil {
					return err
				}
				tf = GenerateLogAnalyticsTerraformBlock(props, inputPath)
			} else if strings.Contains(inputPath, "appinsights") {
				resourceType = "appinsights"
				err := ValidateResourceProperties(provider, resourceType, props)
				if err != nil {
					return err
				}
				tf = GenerateAppInsightsTerraformBlock(props, inputPath)
			} else {
				return fmt.Errorf("unsupported or unrecognized Azure resource type in %s", inputPath)
			}
		}
	} else {
		return fmt.Errorf("unsupported or unrecognized resource type in %s", inputPath)
	}
	if err := os.WriteFile(outputPath, []byte(tfHeader+tf), 0644); err != nil {
		log.Printf("%s Failed to write Terraform output to %s: %v", LogPrefix("GenerateTerraformFromJSONMulticloud"), outputPath, err)
		return err
	}
	log.Printf("%s Terraform code written to %s", LogPrefix("GenerateTerraformFromJSONMulticloud"), outputPath)
	return nil
}
