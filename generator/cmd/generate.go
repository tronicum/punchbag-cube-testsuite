package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"punchbag-cube-testsuite/generator/internal/generator"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Terraform from a YAML/JSON config file or example",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, _ := cmd.Flags().GetString("input")
		fromExample, _ := cmd.Flags().GetString("from-example")
		output, _ := cmd.Flags().GetString("output")
		provider, _ := cmd.Flags().GetString("provider")
		if fromExample != "" {
			input = filepath.Join("examples", fromExample)
			if _, err := os.Stat(input); err != nil {
				return fmt.Errorf("example not found: %s", input)
			}
		}
		if provider == "" {
			provider = "azure"
		}
		if provider == "azure" {
			return generator.GenerateTerraformFromJSON(input, output)
		} else {
			return generator.GenerateTerraformFromJSONMulticloud(input, output, provider)
		}
	},
}

func init() {
	generateCmd.Flags().StringP("input", "i", "", "Path to YAML/JSON config file")
	generateCmd.Flags().StringP("from-example", "e", "", "Name of example YAML/JSON file in examples/")
	generateCmd.Flags().StringP("output", "o", "", "Path to output Terraform file")
	generateCmd.Flags().StringP("provider", "p", "azure", "Cloud provider (azure, aws, gcp)")
	generateCmd.MarkFlagRequired("output")
	rootCmd.AddCommand(generateCmd)
}
