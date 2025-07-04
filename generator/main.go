package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// Read the Terraform template file
	filePath := "azure_services.tf"
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		os.Exit(1)
	}

	// Parse the Terraform template and generate Go code
	lines := strings.Split(string(content), "\n")
	generatedCode := "package azure\n\nimport (\n\t\"fmt\"\n)\n\nfunc CreateAzureResources() {\n"

	for _, line := range lines {
		if strings.HasPrefix(line, "resource \"") {
			resourceType := strings.Split(line, " ")[1]
			resourceName := strings.Split(line, " ")[2]
			generatedCode += fmt.Sprintf("\tfmt.Println(\"Creating %s: %s\")\n", resourceType, resourceName)
		}
	}

	generatedCode += "}\n"

	// Write the generated Go code to a file
	outputFilePath := "generated_resources.go"
	err = os.WriteFile(outputFilePath, []byte(generatedCode), 0644)
	if err != nil {
		fmt.Printf("Failed to write generated code: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Code generation completed. Check generated_resources.go")
}
