package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestEnd2End_MultitoolToTerraform(t *testing.T) {
	// Step 1: Simulate multitool downloading AKS state as JSON
	jsonFile := "test_aks.json"
	jsonContent := `{
  "properties": {
    "name": "test-aks",
    "location": "eastus",
    "resourceGroup": "test-rg",
    "nodeCount": 2
  }
}`
	if err := os.WriteFile(jsonFile, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("Failed to write mock JSON: %v", err)
	}
	defer os.Remove(jsonFile)

	// Step 2: Generate Terraform from JSON using the current binary, not go run main.go
	tfFile := "test_aks.tf"
	bin, err := os.Executable()
	if err != nil {
		t.Fatalf("Failed to get test binary: %v", err)
	}
	cmd := exec.Command(bin, "--generate-terraform", "--input", jsonFile, "--output", tfFile)
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Terraform generation failed: %v\nOutput: %s", err, string(out))
	}
	defer os.Remove(tfFile)

	// Step 3: Validate the generated Terraform file exists and is non-empty
	content, err := os.ReadFile(tfFile)
	if err != nil {
		t.Fatalf("Failed to read generated Terraform: %v", err)
	}
	if len(content) == 0 {
		t.Error("Generated Terraform file is empty")
	}

	// Step 4: Optionally, lint the Terraform file if tflint is available
	if _, err := exec.LookPath("tflint"); err == nil {
		lintCmd := exec.Command("tflint", tfFile)
		lintOut, lintErr := lintCmd.CombinedOutput()
		if lintErr != nil {
			t.Errorf("tflint failed: %v\nOutput: %s", lintErr, string(lintOut))
		}
	}
}
