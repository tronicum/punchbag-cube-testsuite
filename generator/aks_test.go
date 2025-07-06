package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateAksTerraformBlock_Golden(t *testing.T) {
	props := map[string]interface{}{
		"name": "test-aks",
		"location": "eastus",
		"resourceGroup": "test-rg",
		"nodeCount": 3,
	}
	actual := generateAksTerraformBlock(props)
	goldenPath := filepath.Join("testdata", "aks_block.golden")
	expected, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("failed to read golden file: %v", err)
	}
	if strings.TrimSpace(string(expected)) != strings.TrimSpace(actual) {
		t.Errorf("AKS Terraform block does not match golden file.\nExpected:\n%s\nActual:\n%s", string(expected), actual)
	}
}
