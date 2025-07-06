package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cube-generator",
	Short: "Generate, validate, and test Terraform for Azure and multicloud resources",
	Long:  `A modular code generator and test suite for Azure and multicloud Terraform resources.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
