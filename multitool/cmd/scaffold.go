package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var scaffoldCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "Scaffold infrastructure templates for various providers",
}

var scaffoldCloudFormationCmd = &cobra.Command{
	Use:   "cloudformation",
	Short: "Generate a minimal AWS CloudFormation template",
	Run: func(cmd *cobra.Command, args []string) {
		output, _ := cmd.Flags().GetString("output")
		if output == "" {
			output = "cloudformation-template.yaml"
		}
		template := `AWSTemplateFormatVersion: '2010-09-09'
Description: Minimal CloudFormation template
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties: {}
`
		err := os.WriteFile(output, []byte(template), 0644)
		if err != nil {
			fmt.Printf("Failed to write template: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("CloudFormation template written to %s\n", output)
	},
}

var scaffoldARMCmd = &cobra.Command{
	Use:   "arm",
	Short: "Generate a minimal Azure ARM template",
	Run: func(cmd *cobra.Command, args []string) {
		output, _ := cmd.Flags().GetString("output")
		if output == "" {
			output = "azure-arm-template.json"
		}
		template := `{
  "$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {},
  "variables": {},
  "resources": [
    {
      "type": "Microsoft.Storage/storageAccounts",
      "apiVersion": "2022-09-01",
      "name": "mystorageaccount",
      "location": "[resourceGroup().location]",
      "sku": { "name": "Standard_LRS" },
      "kind": "StorageV2",
      "properties": {}
    }
  ],
  "outputs": {}
}`
		err := os.WriteFile(output, []byte(template), 0644)
		if err != nil {
			fmt.Printf("Failed to write template: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Azure ARM template written to %s\n", output)
	},
}

func init() {
	scaffoldCloudFormationCmd.Flags().StringP("output", "o", "", "Output file path")
	scaffoldARMCmd.Flags().StringP("output", "o", "", "Output file path")

	scaffoldCmd.AddCommand(scaffoldCloudFormationCmd)
	scaffoldCmd.AddCommand(scaffoldARMCmd)
}
