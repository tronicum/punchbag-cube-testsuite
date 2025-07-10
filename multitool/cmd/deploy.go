package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Manage Terraform deployment for all clouds",
	Long:  "Plan, apply, or destroy Terraform-managed resources for Azure, AWS, and GCP.",
}

var deployPlanCmd = &cobra.Command{
	Use:   "plan [tf-file]",
	Short: "Run terraform plan on the given file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tfFile := args[0]
		fmt.Printf("Running terraform plan on %s...\n", tfFile)
		cmdExec := exec.Command("terraform", "plan", "-out=tfplan", tfFile)
		cmdExec.Stdout = os.Stdout
		cmdExec.Stderr = os.Stderr
		err := cmdExec.Run()
		if err != nil {
			fmt.Printf("terraform plan failed: %v\n", err)
			os.Exit(1)
		}
	},
}

var deployApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Run terraform apply (applies tfplan if present)",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running terraform apply...")
		cmdExec := exec.Command("terraform", "apply", "tfplan")
		cmdExec.Stdout = os.Stdout
		cmdExec.Stderr = os.Stderr
		err := cmdExec.Run()
		if err != nil {
			fmt.Printf("terraform apply failed: %v\n", err)
			os.Exit(1)
		}
	},
}

var deployDestroyCmd = &cobra.Command{
	Use:   "destroy [tf-file]",
	Short: "Run terraform destroy on the given file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tfFile := args[0]
		fmt.Printf("Running terraform destroy on %s...\n", tfFile)
		cmdExec := exec.Command("terraform", "destroy", "-auto-approve", tfFile)
		cmdExec.Stdout = os.Stdout
		cmdExec.Stderr = os.Stderr
		err := cmdExec.Run()
		if err != nil {
			fmt.Printf("terraform destroy failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	deployCmd.AddCommand(deployPlanCmd)
	deployCmd.AddCommand(deployApplyCmd)
	deployCmd.AddCommand(deployDestroyCmd)
	rootCmd.AddCommand(deployCmd)
}
