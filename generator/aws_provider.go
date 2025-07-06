package main

import (
	"fmt"
	"strings"
)

// AWSProvider implements Provider for AWS resources

type AWSProvider struct{}

func (a AWSProvider) GenerateTerraform(props map[string]interface{}, inputPath string) (string, string, error) {
	tfHeader := `terraform {
  required_version = ">= 1.0.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.0.0"
    }
  }
}

provider "aws" {
  region = "us-west-2"
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
	case "eks":
		tf = generateEksTerraformBlock(props)
	case "s3":
		tf = generateS3TerraformBlock(props)
	default:
		return "", "", fmt.Errorf("unsupported or unrecognized AWS resource type: %s", resourceType)
	}
	return tfHeader, tf, nil
}

func (a AWSProvider) ValidateResource(resourceType string, props map[string]interface{}) error {
	return ValidateResourceProperties("aws", resourceType, props)
}

func (a AWSProvider) DetectResourceType(props map[string]interface{}, inputPath string) (string, error) {
	if _, hasNodeCount := props["nodeCount"]; hasNodeCount {
		return "eks", nil
	} else if _, hasBucket := props["bucket"]; hasBucket || safeString(props, "resourceType", "") == "s3" || strings.Contains(strings.ToLower(safeString(props, "name", "")), "s3") {
		return "s3", nil
	}
	return "", fmt.Errorf("unsupported or unrecognized AWS resource type in %s", inputPath)
}
