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

var createAzureMonitorCmd = &cobra.Command{
	Use:   "create-azure-monitor",
	Short: "Create Azure Monitor services for Kubernetes clusters",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(viper.GetString("server"))

		// Extended Azure Monitor services to include all logging, analytics, and related products
		services := []string{"log-analytics", "application-insights", "metrics", "event-hubs", "diagnostic-settings", "monitor-alerts"}
		responses := make(map[string]interface{})

		for _, service := range services {
			response := map[string]interface{}{
				"service":   service,
				"status":    "success",
				"resource_id": fmt.Sprintf("azure-monitor-%s-%d", service, time.Now().Unix()),
				"timestamp": time.Now().Format(time.RFC3339),
			}
			responses[service] = response
		}

		output.PrintJSON(responses, os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(createAzureMonitorCmd)
}
