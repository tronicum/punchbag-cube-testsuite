package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var simulateCmd = &cobra.Command{
	Use:   "simulate",
	Short: "Dry-run: generate Terraform and run terraform validate/tflint",
	RunE: func(cmd *cobra.Command, args []string) error {
		output, _ := cmd.Flags().GetString("output")
		if output == "" {
			return fmt.Errorf("output file required for simulation")
		}
		// Run terraform validate
		validateCmd := exec.Command("terraform", "validate", output)
		validateCmd.Stdout = os.Stdout
		validateCmd.Stderr = os.Stderr
		if err := validateCmd.Run(); err != nil {
			return fmt.Errorf("terraform validate failed: %w", err)
		}
		// Run tflint
		tflintCmd := exec.Command("tflint", output)
		tflintCmd.Stdout = os.Stdout
		tflintCmd.Stderr = os.Stderr
		if err := tflintCmd.Run(); err != nil {
			return fmt.Errorf("tflint failed: %w", err)
		}
		fmt.Println("Simulation successful: terraform validate and tflint passed.")
		return nil
	},
}

func init() {
	simulateCmd.Flags().StringP("output", "o", "", "Path to output Terraform file or directory")
	simulateCmd.MarkFlagRequired("output")
	rootCmd.AddCommand(simulateCmd)
}
