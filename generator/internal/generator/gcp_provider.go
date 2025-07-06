package generator

import "fmt"

type GCPProvider struct{}

func (g GCPProvider) GenerateTerraform(props map[string]interface{}, inputPath string) (string, string, error) {
	tfHeader := `terraform {
  required_version = ">= 1.0.0"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 4.0.0"
    }
  }
}

provider "google" {}
`
	resourceType, err := g.DetectResourceType(props, inputPath)
	if err != nil {
		return "", "", err
	}
	if err := g.ValidateResource(resourceType, props); err != nil {
		return "", "", err
	}
	var tf string
	switch resourceType {
	case "gke":
		tf = GenerateGkeTerraformBlock(props)
	default:
		return "", "", fmt.Errorf("unsupported or unrecognized GCP resource type: %s", resourceType)
	}
	return tfHeader, tf, nil
}

func (g GCPProvider) ValidateResource(resourceType string, props map[string]interface{}) error {
	return ValidateResourceProperties("gcp", resourceType, props)
}

func (g GCPProvider) DetectResourceType(props map[string]interface{}, inputPath string) (string, error) {
	if _, hasNodeCount := props["nodeCount"]; hasNodeCount {
		return "gke", nil
	}
	return "", fmt.Errorf("unsupported or unrecognized GCP resource type in %s", inputPath)
}
