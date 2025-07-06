package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
)

// werfty-transformator: Convert Terraform between cloud providers or to multipass-cloud-layer
// Usage:
//   werfty-transformator --input <input.tf> --src-provider <azure|aws|gcp> --destination-provider <azure|aws|gcp|multipass-cloud-layer>
//
// Supported conversions:
//   Azure Blob <-> AWS S3
//   Any S3-like -> multipass-cloud-layer

func main() {
	inputPath := flag.String("input", "", "Input Terraform file")
	srcProvider := flag.String("src-provider", "", "Source cloud provider (azure|aws|gcp)")
	destProvider := flag.String("destination-provider", "", "Destination cloud provider (azure|aws|gcp|multipass-cloud-layer)")
	flag.Parse()

	if *inputPath == "" || *srcProvider == "" || *destProvider == "" {
		fmt.Println("Usage: werfty-transformator --input <input.tf> --src-provider <azure|aws|gcp> --destination-provider <azure|aws|gcp|multipass-cloud-layer>")
		os.Exit(1)
	}
	content, err := os.ReadFile(*inputPath)
	if err != nil {
		fmt.Printf("Failed to read input: %v\n", err)
		os.Exit(1)
	}
	converted := ConvertTerraform(string(content), *srcProvider, *destProvider)
	fmt.Println(converted)
}

// ConvertTerraform maps resources from src to dest provider.
func ConvertTerraform(tf, src, dest string) string {
	if src == "azure" && dest == "aws" {
		return ConvertAzureBlobToAWSS3(tf)
	}
	if src == "aws" && dest == "azure" {
		return ConvertAWSS3ToAzureBlob(tf)
	}
	if dest == "multipass-cloud-layer" {
		return regexp.MustCompile(`resource ".*?_s3_bucket"`).ReplaceAllString(tf, "resource \"multipass_cloud_layer_bucket\"")
	}
	return "# Conversion logic not yet implemented for this provider pair"
}

// ConvertAzureBlobToAWSS3 maps Azure Blob Storage Terraform to AWS S3 Terraform
func ConvertAzureBlobToAWSS3(tf string) string {
	reAccount := regexp.MustCompile(`(?s)resource\s+\"azurerm_storage_account\"[^{]*\{.*?\}`)
	accountBlocks := reAccount.FindAllString(tf, -1)
	var awsBlocks []string
	for _, block := range accountBlocks {
		name := extractField(block, "name", "examplestorageacct")
		region := extractField(block, "location", "us-east-1")
		awsBlocks = append(awsBlocks, `resource "aws_s3_bucket" "example" {
  bucket = "`+name+`"
  region = "`+region+`"
}`)
	}

	tf = reAccount.ReplaceAllString(tf, "")
	// Remove unused variable containerBlocks
	reContainer := regexp.MustCompile(`(?s)resource\s+\"azurerm_storage_container\"[^{]*\{.*?\}`)

	tf = reContainer.ReplaceAllString(tf, "")
	if len(awsBlocks) > 0 {
		tf += "\n" + awsBlocks[0] + "\n"
	}
	return tf
}

// ConvertAWSS3ToAzureBlob maps AWS S3 Terraform to Azure Blob Storage Terraform
func ConvertAWSS3ToAzureBlob(tf string) string {
	reS3 := regexp.MustCompile(`(?s)resource\s+\"aws_s3_bucket\"[^{]*\{.*?\}`)
	s3Blocks := reS3.FindAllString(tf, -1)
	var azureBlocks []string
	for _, block := range s3Blocks {
		name := extractField(block, "bucket", "examplestorageacct")
		region := extractField(block, "region", "West Europe")
		azureBlocks = append(azureBlocks, `resource "azurerm_storage_account" "example" {
  name                     = "`+name+`"
  resource_group_name      = "example-resources"
  location                 = "`+region+`"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "content"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}`)
	}

	tf = reS3.ReplaceAllString(tf, "")
	if len(azureBlocks) > 0 {
		tf += "\n" + azureBlocks[0] + "\n"
	}
	return tf
}

// extractField extracts a string field from a Terraform resource block.
func extractField(block, field, def string) string {
	re := regexp.MustCompile(field + `\s*=\s*\"([^\"]+)\"`)
	match := re.FindStringSubmatch(block)
	if len(match) > 1 {
		return match[1]
	}
	return def
}
