package main

import (
	"encoding/json"
	"os"
	"testing"
	"punchbag-cube-testsuite/generator/internal/generator"
)

type MultiCloudTestCase struct {
	Provider     string
	ResourceType string
	TestDataFile string
	Expected     string // Optionally, path to expected output or inline string
}

// loadTestData loads JSON test data from a file into a map
func loadTestData(t *testing.T, filename string) map[string]interface{} {
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("failed to read test data file %s: %v", filename, err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("failed to unmarshal test data: %v", err)
	}
	return m
}

func TestMultiCloudResourceGeneration(t *testing.T) {
	testCases := []MultiCloudTestCase{
		{Provider: "azure", ResourceType: "aks", TestDataFile: "../testdata/aks.json"},
		{Provider: "aws", ResourceType: "eks", TestDataFile: "../testdata/eks.json"},
		{Provider: "gcp", ResourceType: "gke", TestDataFile: "../testdata/gke.json"},
	}

	for _, tc := range testCases {
		t.Run(tc.Provider+"-"+tc.ResourceType, func(t *testing.T) {
			props := loadTestData(t, tc.TestDataFile)
			var tfBlock string
			switch tc.Provider {
			case "azure":
				tfBlock = generateAksTerraformBlock(props)
			case "aws":
				if tc.ResourceType == "eks" {
					tfBlock = generateEksTerraformBlock(props)
				} else if tc.ResourceType == "s3" {
					tfBlock = generateS3TerraformBlock(props)
				}
			case "gcp":
				tfBlock = generateGkeTerraformBlock(props)
			default:
				t.Fatalf("unsupported provider: %s", tc.Provider)
			}
			if tfBlock == "" {
				t.Errorf("Terraform block is empty for %s/%s", tc.Provider, tc.ResourceType)
			}
			// Optionally: compare tfBlock to expected output if provided
		})
	}
}

func TestGenerateS3TerraformBlock(t *testing.T) {
	props := map[string]interface{}{
		"name":       "test-s3-bucket",
		"acl":        "private",
		"versioning": true,
	}
	tf := generateS3TerraformBlock(props)
	if len(tf) == 0 || !contains(tf, "aws_s3_bucket") {
		t.Errorf("S3 Terraform block not generated correctly: %s", tf)
	}
}

func TestGenerateTerraformFromJSONMulticloud_S3(t *testing.T) {
	inputFile := "../testdata/s3.json"
	outputFile := "test_s3.tf"
	defer os.Remove(outputFile)
	if err := generator.GenerateTerraformFromJSONMulticloud(inputFile, outputFile, "aws"); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}
	if !contains(string(content), "aws_s3_bucket") {
		t.Error("S3 Terraform not generated")
	}
}

func TestMultiCloudResourceGeneration_NegativeCases(t *testing.T) {
	cases := []struct {
		provider     string
		resourceType string
		props        map[string]interface{}
	}{
		{"aws", "eks", map[string]interface{}{"name": "eks"}}, // missing region, nodeCount
		{"gcp", "gke", map[string]interface{}{"location": "loc"}}, // missing name, nodeCount
		{"azure", "aks", map[string]interface{}{"name": "aks", "location": "loc"}}, // missing resourceGroup, nodeCount
	}
	for _, c := range cases {
		err := ValidateResourceProperties(c.provider, c.resourceType, c.props)
		if err == nil {
			t.Errorf("Expected error for %s/%s with props %v", c.provider, c.resourceType, c.props)
		}
	}
}

func TestMultiCloudResourceGeneration_UnknownProviderType(t *testing.T) {
	err := ValidateResourceProperties("foo", "bar", map[string]interface{}{"name": "n"})
	if err == nil {
		t.Error("Expected error for unknown provider/type")
	}
}
