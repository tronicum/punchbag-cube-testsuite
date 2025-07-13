package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	importpkg "github.com/tronicum/punchbag-cube-testsuite/shared/import"
	"github.com/tronicum/punchbag-cube-testsuite/werfty-transformator/transform"
)

// werfty-transformator: Convert Terraform between cloud providers or to multipass-cloud-layer
// Usage:
//   werfty-transformator --input <input.tf> --src-provider <azure|aws|gcp> --destination-provider <azure|aws|gcp|multipass-cloud-layer>
//
// Supported conversions:
//   Azure Blob <-> AWS S3
//   Any S3-like -> multipass-cloud-layer

func main() {
	inputPath := flag.String("input", "", "Input Terraform file")
	configPath := flag.String("config", "", "Optional config file (JSON/YAML)")
	srcProvider := flag.String("src-provider", "", "Source cloud provider (azure|aws|gcp)")
	destProvider := flag.String("destination-provider", "", "Destination cloud provider (azure|aws|gcp|multipass-cloud-layer)")
	terraspace := flag.Bool("terraspace", false, "Output as Terraspace project structure")
	flag.Parse()

	if *inputPath == "" || *srcProvider == "" || *destProvider == "" {
		fmt.Println("Usage: werfty-transformator --input <input.tf> --src-provider <azure|aws|gcp> --destination-provider <azure|aws|gcp|multipass-cloud-layer> [--config <config.json|yaml>] [--terraspace]")
		os.Exit(1)
	}

	// Optionally load config file for transformation context
	var config *importpkg.Config
	if *configPath != "" {
		f, err := os.Open(*configPath)
		if err != nil {
			fmt.Printf("Failed to open config file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		config, err = importpkg.LoadConfigJSON(f)
		if err != nil {
			fmt.Printf("Failed to parse config file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Loaded config: %+v\n", config)
	}

	content, err := os.ReadFile(*inputPath)
	if err != nil {
		fmt.Printf("Failed to read input: %v\n", err)
		os.Exit(1)
	}
	converted := ConvertTerraform(string(content), *srcProvider, *destProvider)

	if *terraspace {
		// Ensure terraspace is installed using multitool
		fmt.Println("Ensuring terraspace is installed via multitool...")
		_ = runMultitoolInstall("terraspace")
		// Scaffold Terraspace project structure
		writeTerraspaceProject(converted, *destProvider)
	} else {
		fmt.Println(converted)
	}
}

// ConvertTerraform maps resources from src to dest provider.
func ConvertTerraform(tf, src, dest string) string {
	if src == "azure" && dest == "aws" {
		// Convert storage first
		tf = transform.ConvertAzureBlobToAWSS3(tf)
		// Convert monitoring resources
		tf = transform.ConvertAzureMonitorToAWSCloudWatch(tf)
		// Convert budget resources
		tf = transform.ConvertAzureBudgetToAWSBudget(tf)
		return tf
	}
	if src == "aws" && dest == "azure" {
		return transform.ConvertAWSS3ToAzureBlob(tf)
	}
	if dest == "multipass-cloud-layer" {
		return transform.ConvertS3LikeToMultipassCloudLayer(tf)
	}
	return "# Conversion logic not yet implemented for this provider pair"
}

func runMultitoolInstall(pkg string) error {
	cmd := "../multitool/mt install-package " + pkg
	fmt.Println("Running:", cmd)
	return runShell(cmd)
}

func runShell(cmd string) error {
	c := exec.Command("sh", "-c", cmd)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func writeTerraspaceProject(tf, provider string) {
	projectRoot := "terraspace_project"
	modulePath := projectRoot + "/app/modules/converted"
	stackPath := projectRoot + "/app/stacks/converted"
	configPath := projectRoot + "/config"

	// Create directories
	if err := os.MkdirAll(modulePath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create module dir: %v\n", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(stackPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create stack dir: %v\n", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(configPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create config dir: %v\n", err)
		os.Exit(1)
	}

	// Write main.tf to module
	mainTf := modulePath + "/main.tf"
	if err := os.WriteFile(mainTf, []byte(tf), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write main.tf: %v\n", err)
		os.Exit(1)
	}

	// Write minimal stack that uses the module
	stackTf := stackPath + "/main.tf"
	stackContent := `module "converted" {
  source = "../../modules/converted"
}`
	if err := os.WriteFile(stackTf, []byte(stackContent), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write stack main.tf: %v\n", err)
		os.Exit(1)
	}

	// Write minimal config/app.rb
	appRb := configPath + "/app.rb"
	if err := os.WriteFile(appRb, []byte("Terraspace.configure do |config|\nend\n"), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write app.rb: %v\n", err)
		os.Exit(1)
	}
	// Write minimal config/terraform.rb
	terraformRb := configPath + "/terraform.rb"
	if err := os.WriteFile(terraformRb, []byte("Terraspace.configure do |config|\nend\n"), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write terraform.rb: %v\n", err)
		os.Exit(1)
	}
	// Write Gemfile
	gemfile := projectRoot + "/Gemfile"
	if err := os.WriteFile(gemfile, []byte("source 'https://rubygems.org'\ngem 'terraspace'\n"), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write Gemfile: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Terraspace project generated in ./terraspace_project/")
	fmt.Println("Next steps:")
	fmt.Println("  cd terraspace_project && bundle install && terraspace up converted")
	// TODO: Add support for additional providers and advanced config scaffolding
}

// Ensure transform.ConvertAWSS3ToAzureBlob and transform.ConvertS3LikeToMultipassCloudLayer are implemented in transform package.
