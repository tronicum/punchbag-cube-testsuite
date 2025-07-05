package cmd

import (
	"fmt"

	"punchbag-cube-testsuite/client/pkg/api"
	"punchbag-cube-testsuite/client/pkg/output"

	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate [provider]",
	Short: "Validate cloud provider configuration and availability",
	Long: `Validate cloud provider configuration and check service availability.
Supported providers: azure, hetzner-hcloud, united-ionos, schwarz-stackit, aws, gcp

Examples:
  punchbag-cube-testsuite validate azure
  punchbag-cube-testsuite validate hetzner-hcloud
  punchbag-cube-testsuite validate united-ionos`,
	Args: cobra.ExactArgs(1),
	RunE: runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) error {
	provider := args[0]
	
	client := api.NewClient(apiBaseURL)
	
	result, err := client.ValidateProvider(provider)
	if err != nil {
		return fmt.Errorf("failed to validate provider %s: %w", provider, err)
	}
	
	formatter := output.NewFormatter(outputFormat)
	return formatter.FormatProviderValidation(result)
}
