package main

import (
	"os/exec"
	"testing"
	"os"
)

func TestTerraformValidateAndTflint(t *testing.T) {
	// This test assumes a valid Terraform file is generated at test_monitor.tf
	// and terraform/tflint are installed in the environment.
	file := "test_monitor.tf"
	if _, err := os.Stat(file); os.IsNotExist(err) {
		t.Skip("Terraform file not found, skipping validation/lint test")
	}
	cmd := exec.Command("terraform", "validate", file)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("terraform validate failed: %v\n%s", err, string(out))
	}
	tflint := exec.Command("tflint", file)
	out, err = tflint.CombinedOutput()
	if err != nil {
		t.Fatalf("tflint failed: %v\n%s", err, string(out))
	}
}
