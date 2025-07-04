package main

import (
	"fmt"
	"net/http"
	"io"
)

func main() {
	server := "http://localhost:8080"

	// Validate provider configuration
	provider := "aws"
	validateURL := fmt.Sprintf("%s/api/v1/validate/%s", server, provider)
	response, err := http.Get(validateURL)
	if err != nil {
		fmt.Println("Error validating provider:", err)
		return
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	fmt.Println("Validation Response:", string(body))

	// Simulate provider operation
	operation := "create-cluster"
	simulateURL := fmt.Sprintf("%s/api/v1/providers/%s/operations/%s", server, provider, operation)
	response, err = http.Post(simulateURL, "application/json", nil)
	if err != nil {
		fmt.Println("Error simulating provider operation:", err)
		return
	}
	defer response.Body.Close()
	body, _ = io.ReadAll(response.Body)
	fmt.Println("Simulation Response:", string(body))
}
