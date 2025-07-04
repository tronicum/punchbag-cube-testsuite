package cmd

import (
	"fmt"
	"os"

	"punchbag-cube-testsuite/client/pkg/api"
	"punchbag-cube-testsuite/client/pkg/output"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Manage test results",
	Long:  `Commands for managing test results and viewing test status.`,
}

var testGetCmd = &cobra.Command{
	Use:   "get [test-id]",
	Short: "Get test result details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(viper.GetString("server"))
		result, err := client.GetTestResult(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting test result: %v\n", err)
			os.Exit(1)
		}

		output.PrintTestResult(result, viper.GetString("format"))
	},
}

var testListCmd = &cobra.Command{
	Use:   "list [cluster-id]",
	Short: "List test results for a cluster",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(viper.GetString("server"))
		results, err := client.ListTestResults(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing test results: %v\n", err)
			os.Exit(1)
		}

		output.PrintTestResults(results, viper.GetString("format"))
	},
}

var testWatchCmd = &cobra.Command{
	Use:   "watch [test-id]",
	Short: "Watch test progress",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(viper.GetString("server"))
		
		fmt.Printf("Watching test %s (press Ctrl+C to stop)...\n", args[0])
		
		for {
			result, err := client.GetTestResult(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting test result: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("\rStatus: %s", result.Status)
			if result.Status == "completed" || result.Status == "failed" {
				fmt.Println()
				output.PrintTestResult(result, viper.GetString("format"))
				break
			}
			
			// Sleep for 2 seconds before checking again
			cmd.Context().Done()
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.AddCommand(testGetCmd)
	testCmd.AddCommand(testListCmd)
	testCmd.AddCommand(testWatchCmd)
}
