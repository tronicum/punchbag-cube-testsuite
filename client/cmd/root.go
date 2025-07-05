package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "punchbag-client",
	Short: "CLI client for Punchbag Cube Test Suite",
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
		fmt.Println(string(body))
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
		fmt.Println(string(body))
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
