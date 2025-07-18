package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "werfty",
	Short: "Terraform generation and validation tool",
	Long: `A command line interface for interacting with the Punchbag Cube Test Suite API.
	
This client allows you to manage AKS clusters, run tests, and view results
from the command line.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate provider configuration",
	Run: func(cmd *cobra.Command, args []string) {
		provider := args[0]
		server := viper.GetString("server")
		url := fmt.Sprintf("%s/api/v1/validate/%s", server, provider)
		response, err := http.Get(url)
		if err != nil {
			fmt.Println("Error validating provider:", err)
			return
		}
		defer response.Body.Close()
		body, _ := io.ReadAll(response.Body)
		format := viper.GetString("format")
		fmt.Println(formatOutput(string(body), format))
	},
}

var simulateProviderCmd = &cobra.Command{
	Use:   "simulate-provider",
	Short: "Simulate provider operations",
	Run: func(cmd *cobra.Command, args []string) {
		provider := args[0]
		operation := args[1]
		server := viper.GetString("server")
		url := fmt.Sprintf("%s/api/v1/providers/%s/operations/%s", server, provider, operation)
		response, err := http.Post(url, "application/json", nil)
		if err != nil {
			fmt.Println("Error simulating provider operation:", err)
			return
		}
		defer response.Body.Close()
		body, _ := io.ReadAll(response.Body)
		format := viper.GetString("format")
		fmt.Println(formatOutput(string(body), format))
	},
}

// Add commands for generating Azure templates
var generateAzureCmd = &cobra.Command{
	Use:   "generate-template",
	Short: "Generate templates for Azure services",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please specify the Azure service: monitoring, kubernetes, or budget")
			return
		}
		format := viper.GetString("format")
		switch args[0] {
		case "monitoring":		werfty validate azure --format json
			fmt.Println(formatOutput(generator.GenerateAzureMonitoringTemplate(nil), format))
		case "kubernetes":
			fmt.Println(formatOutput(generator.GenerateAzureKubernetesTemplate(nil), format))
		case "budget":
			fmt.Println(formatOutput(generator.GenerateAzureBudgetTemplate(nil), format))
		default:
			fmt.Println("Unknown service. Please specify monitoring, kubernetes, or budget.")
		}
	},
}

// Add command for generating Azure Log Analytics template
var generateLogAnalyticsCmd = &cobra.Command{
	Use:   "generate-log-analytics",
	Short: "Generate a Terraform template for Azure Log Analytics",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(generator.GenerateAzureLogAnalyticsTemplate(nil))
	},
}

// Add missing commands to werfty
var generateTerraformCmd = &cobra.Command{
	Use:   "generate-terraform",
	Short: "Generate Terraform templates",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Generating Terraform templates...")
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.punchbag-client.yaml)")
	rootCmd.PersistentFlags().String("server", "http://localhost:8080", "Server URL")
	rootCmd.PersistentFlags().String("format", "table", "Output format (table, json, yaml)")
	rootCmd.PersistentFlags().Bool("verbose", false, "Verbose output")

	viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
	viper.BindPFlag("format", rootCmd.PersistentFlags().Lookup("format"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(simulateProviderCmd)
	rootCmd.AddCommand(generateAzureCmd)
	rootCmd.AddCommand(generateLogAnalyticsCmd)
	rootCmd.AddCommand(generateTerraformCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".punchbag-client")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// Helper function to format output
func formatOutput(data interface{}, format string) string {
	if format == "json" {
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return fmt.Sprintf("Error formatting JSON: %v", err)
		}
		return string(jsonData)
	} else if format == "table" {
		// Example table formatting (can be replaced with a library like "tablewriter")
		return fmt.Sprintf("%v", data)
	} else if format == "yaml" {
		yamlData, err := yaml.Marshal(data)
		if err != nil {
			return fmt.Sprintf("Error formatting YAML: %v", err)
		}
		return string(yamlData)
	}
	return fmt.Sprintf("Unknown format: %s", format)
}
