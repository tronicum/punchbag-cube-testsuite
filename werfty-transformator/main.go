package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tronicum/punchbag-cube-testsuite/werfty-transformator/transform"
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
		return transform.ConvertAzureBlobToAWSS3(tf)
	}
	if src == "aws" && dest == "azure" {
		return transform.ConvertAWSS3ToAzureBlob(tf)
	}
	if dest == "multipass-cloud-layer" {
		return transform.ConvertS3LikeToMultipassCloudLayer(tf)
	}
	return "# Conversion logic not yet implemented for this provider pair"
}
