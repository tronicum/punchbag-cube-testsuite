package main

import (
	"fmt"
	"punchbag-cube-testsuite/generator"
)

func main() {
	// Example usage of Azure Monitoring template generation
	monitoringTemplate := generator.GenerateAzureMonitoringTemplate(nil)
	fmt.Println("Generated Azure Monitoring Template:\n", monitoringTemplate)

	// Example usage of Azure Kubernetes template generation
	kubernetesTemplate := generator.GenerateAzureKubernetesTemplate(nil)
	fmt.Println("Generated Azure Kubernetes Template:\n", kubernetesTemplate)

	// Example usage of Azure Budget template generation
	budgetTemplate := generator.GenerateAzureBudgetTemplate(nil)
	fmt.Println("Generated Azure Budget Template:\n", budgetTemplate)
}
