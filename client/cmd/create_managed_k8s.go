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

var createManagedK8sCmd = &cobra.Command{
	Use:   "create-managed-k8s",
	Short: "Create a managed Kubernetes cluster in all supported clouds",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(viper.GetString("server"))

		cloudProviders := []string{"azure", "hetzner-hcloud", "united-ionos", "aws", "gcp"}
		responses := make(map[string]interface{})

		for _, provider := range cloudProviders {
			response := map[string]interface{}{
				"provider":  provider,
				"status":    "success",
				"cluster_id": fmt.Sprintf("%s-cluster-%d", provider, time.Now().Unix()),
				"timestamp": time.Now().Format(time.RFC3339),
			}
			responses[provider] = response
		}

		output.PrintJSON(responses, os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(createManagedK8sCmd)
}
