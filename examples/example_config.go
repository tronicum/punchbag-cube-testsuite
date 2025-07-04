package main

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port int    `yaml:"port"`
		Mode string `yaml:"mode"`
	} `yaml:"server"`
	Client struct {
		ServerURL string `yaml:"server_url"`
		Format    string `yaml:"format"`
	} `yaml:"client"`
	Generator struct {
		OutputDir    string `yaml:"output_dir"`
		TemplateFile string `yaml:"template_file"`
	} `yaml:"generator"`
}

func main() {
	// Load the configuration file
	file, err := os.Open("../conf/punchy.yml")
	if err != nil {
		fmt.Printf("Failed to open config file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		fmt.Printf("Failed to parse config file: %v\n", err)
		os.Exit(1)
	}

	// Print the loaded configuration
	fmt.Println("Loaded Configuration:")
	fmt.Printf("Server Port: %d\n", config.Server.Port)
	fmt.Printf("Server Mode: %s\n", config.Server.Mode)
	fmt.Printf("Client Server URL: %s\n", config.Client.ServerURL)
	fmt.Printf("Client Format: %s\n", config.Client.Format)
	fmt.Printf("Generator Output Directory: %s\n", config.Generator.OutputDir)
	fmt.Printf("Generator Template File: %s\n", config.Generator.TemplateFile)

	// Example usage
	fmt.Println("\nExample Usage:")
	fmt.Println("1. Start the server using the specified port and mode.")
	fmt.Println("2. Use the client to interact with the server at the specified URL.")
	fmt.Println("3. Generate templates using the generator configuration.")
}
