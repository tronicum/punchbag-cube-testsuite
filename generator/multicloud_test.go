package main

import (
	"encoding/json"
	"os"
	"testing"
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
				tfBlock = generateEksTerraformBlock(props)
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
