package punchbag

import (
	"os"
	"testing"
)

func TestGenerateTerraformFromJSON(t *testing.T) {
	input := `{"properties": {"name": "test-monitor"}}`
	inputFile := "test_monitor.json"
	outputFile := "test_monitor.tf"
	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)
	err := GenerateTerraformFromJSON(inputFile, outputFile)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}
	if len(content) == 0 {
		t.Error("Output Terraform file is empty")
	}
}
