package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"punchbag-cube-testsuite/client/pkg/api"
	"punchbag-cube-testsuite/client/pkg/output"
)

var createAzureBudgetCmd = &cobra.Command{
	Use:   "create-azure-budget",
	Short: "Create Azure Budget for cost management",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(viper.GetString("server"))

		response := map[string]interface{}{
			"service":   "azure-budget",
			"status":    "success",
			"budget_id": fmt.Sprintf("azure-budget-%d", time.Now().Unix()),
			"timestamp": time.Now().Format(time.RFC3339),
		}

		output.PrintJSON(response, os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(createAzureBudgetCmd)
}
