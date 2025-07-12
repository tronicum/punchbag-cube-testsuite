package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Generate and apply Terraform from an example config (runs terraform apply)",
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
		// Run terraform apply (with user confirmation)
		fmt.Println("About to run 'terraform apply'. Press Enter to continue or Ctrl+C to abort.")
		fmt.Scanln()
		applyCmd := exec.Command("terraform", "apply", output)
		applyCmd.Stdout = os.Stdout
		applyCmd.Stderr = os.Stderr
		if err := applyCmd.Run(); err != nil {
			return fmt.Errorf("terraform apply failed: %w", err)
		}
		fmt.Println("Apply successful.")
		return nil
	},
}

func init() {
	executeCmd.Flags().StringP("from-example", "e", "", "Name of example YAML/JSON file in examples/")
	executeCmd.Flags().StringP("output", "o", "", "Path to output Terraform file")
	executeCmd.Flags().StringP("provider", "p", "azure", "Cloud provider (azure, aws, gcp)")
	rootCmd.AddCommand(executeCmd)
}
