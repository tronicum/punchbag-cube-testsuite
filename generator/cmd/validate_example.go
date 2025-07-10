package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var validateExampleCmd = &cobra.Command{
	Use:   "validate-example",
	Short: "Generate and validate Terraform from an example config",
	RunE: func(cmd *cobra.Command, args []string) error {
		example, _ := cmd.Flags().GetString("from-example")
		output, _ := cmd.Flags().GetString("output")
		provider, _ := cmd.Flags().GetString("provider")
		if example == "" || output == "" {
			return fmt.Errorf("--from-example and --output are required")
		}
		input := filepath.Join("examples", example)
		if _, err := os.Stat(input); err != nil {
			return fmt.Errorf("example not found: %s", input)
		}
		// Generate Terraform
		generateArgs := []string{"generate", "--input", input, "--output", output, "--provider", provider}
		if err := exec.Command(os.Args[0], generateArgs...).Run(); err != nil {
			return fmt.Errorf("generation failed: %w", err)
		}
		// Run terraform init
		initCmd := exec.Command("terraform", "init")
		initCmd.Dir = filepath.Dir(output)
		initCmd.Stdout = os.Stdout
		initCmd.Stderr = os.Stderr
		_ = initCmd.Run() // ignore errors for init
		// Run terraform validate
		validateCmd := exec.Command("terraform", "validate", output)
		validateCmd.Stdout = os.Stdout
		validateCmd.Stderr = os.Stderr
		if err := validateCmd.Run(); err != nil {
			return fmt.Errorf("terraform validate failed: %w", err)
		}
		fmt.Println("Validation successful.")
		return nil
	},
}

func init() {
	validateExampleCmd.Flags().StringP("from-example", "e", "", "Name of example YAML/JSON file in examples/")
	validateExampleCmd.Flags().StringP("output", "o", "", "Path to output Terraform file")
	validateExampleCmd.Flags().StringP("provider", "p", "azure", "Cloud provider (azure, aws, gcp)")
	rootCmd.AddCommand(validateExampleCmd)
}
