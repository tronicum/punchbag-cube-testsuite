package cmd

import (
	"fmt"
	"punchbag-cube-testsuite/generator/internal/generator"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a YAML/JSON config file against the resource schema",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, _ := cmd.Flags().GetString("input")
		provider, _ := cmd.Flags().GetString("provider")
		if provider == "" {
			provider = "azure"
		}
		cfg, err := generator.LoadConfigFromFile(input)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		props, ok := cfg["properties"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("no 'properties' key found or not a map in %s", input)
		}
		resourceType := ""
		if v, ok := cfg["resourceType"].(string); ok {
			resourceType = v
		} else if v, ok := props["resourceType"].(string); ok {
			resourceType = v
		}
		if resourceType == "" {
			return fmt.Errorf("resourceType not specified in config")
		}
		if err := generator.ValidateResourceProperties(provider, resourceType, props); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
		fmt.Println("Validation successful.")
		return nil
	},
}

func init() {
	validateCmd.Flags().StringP("input", "i", "", "Path to YAML/JSON config file")
	validateCmd.Flags().StringP("provider", "p", "azure", "Cloud provider (azure, aws, gcp)")
	validateCmd.MarkFlagRequired("input")
	rootCmd.AddCommand(validateCmd)
}
