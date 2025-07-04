package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// Step 1: Create a Kubernetes cluster
	fmt.Println("Creating Kubernetes cluster...")
	cmd := exec.Command("punchbag-client", "cluster", "create", "--name", "example-cluster", "--resource-group", "example-resources", "--location", "East US")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to create cluster: %v\n", err)
		os.Exit(1)
	}

	// Step 2: Generate templates using the generator
	fmt.Println("Generating templates...")
	cmd = exec.Command("go", "run", "generator/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to generate templates: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Example completed successfully.")
}
