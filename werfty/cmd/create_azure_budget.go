package cmd

import (
	"fmt"
	"os"
	"time"

	"punchbag-cube-testsuite/client/pkg/api"
	"punchbag-cube-testsuite/client/pkg/output"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createAzureBudgetCmd = &cobra.Command{
	Use:   "create-azure-budget",
	Short: "Create Azure Budget for cost management",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(viper.GetString("server"))

		name, _ := cmd.Flags().GetString("name")
		amount, _ := cmd.Flags().GetFloat64("amount")
		resourceGroup, _ := cmd.Flags().GetString("resource-group")
		timeGrain, _ := cmd.Flags().GetString("time-grain")
		alertThreshold, _ := cmd.Flags().GetFloat64("alert-threshold")

		response := map[string]interface{}{
			"service":            "azure-budget",
			"status":             "success",
			"budget_id":          fmt.Sprintf("azure-budget-%d", time.Now().Unix()),
			"name":               name,
			"amount":             amount,
			"resource_group":     resourceGroup,
			"time_grain":         timeGrain,
			"alert_threshold":    alertThreshold,
			"timestamp":          time.Now().Format(time.RFC3339),
			"terraform_template": generateAzureBudgetTerraform(name, amount, resourceGroup, timeGrain, alertThreshold),
		}

		output.PrintJSON(response, os.Stdout)
	},
}

func generateAzureBudgetTerraform(name string, amount float64, resourceGroup, timeGrain string, alertThreshold float64) string {
	return fmt.Sprintf(`
resource "azurerm_consumption_budget_resource_group" "%s" {
  name              = "%s"
  resource_group_id = data.azurerm_resource_group.%s.id
  
  amount         = %.2f
  time_grain     = "%s"
  
  time_period {
    start_date = "%s"
    end_date   = "%s"
  }
  
  filter {
    dimension {
      name = "ResourceGroupName"
      values = ["%s"]
    }
  }
  
  notification {
    enabled        = true
    threshold      = %.0f
    operator       = "GreaterThan"
    threshold_type = "Actual"
    
    contact_emails = [
      # Add your email addresses here
    ]
  }
}

data "azurerm_resource_group" "%s" {
  name = "%s"
}`, name, name, resourceGroup, amount, timeGrain,
		time.Now().Format("2006-01-02"),
		time.Now().AddDate(1, 0, 0).Format("2006-01-02"),
		resourceGroup, alertThreshold, resourceGroup, resourceGroup)
}

func init() {
	rootCmd.AddCommand(createAzureBudgetCmd)

	createAzureBudgetCmd.Flags().String("name", "", "Budget name")
	createAzureBudgetCmd.Flags().Float64("amount", 1000.0, "Budget amount")
	createAzureBudgetCmd.Flags().String("resource-group", "", "Azure resource group")
	createAzureBudgetCmd.Flags().String("time-grain", "Monthly", "Budget time grain (Monthly, Quarterly, Annually)")
	createAzureBudgetCmd.Flags().Float64("alert-threshold", 80.0, "Alert threshold percentage")
	createAzureBudgetCmd.MarkFlagRequired("name")
	createAzureBudgetCmd.MarkFlagRequired("resource-group")
}
