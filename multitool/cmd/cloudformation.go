package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
)

var cloudformationCmd = &cobra.Command{
	Use:   "cloudformation",
	Short: "Manage AWS CloudFormation resources",
}

var getStackCmd = &cobra.Command{
	Use:   "get-stack",
	Short: "Download a CloudFormation stack template from AWS or proxy",
	Run: func(cmd *cobra.Command, args []string) {
		stackName, _ := cmd.Flags().GetString("stack-name")
		output, _ := cmd.Flags().GetString("output")
		if stackName == "" {
			fmt.Println("--stack-name is required")
			os.Exit(1)
		}
		if output == "" {
			output = stackName + "-template.yaml"
		}
		var templateBody string
		if proxyServer != "" {
			// Fetch from proxy/simulator
			url := fmt.Sprintf("%s/api/v1/cloudformation/stack?name=%s", proxyServer, stackName)
			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("Failed to fetch from proxy: %v\n", err)
				os.Exit(1)
			}
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			templateBody = string(body)
		} else {
			// Fetch from AWS
			cfg, err := config.LoadDefaultConfig(context.Background())
			if err != nil {
				fmt.Printf("Failed to load AWS config: %v\n", err)
				os.Exit(1)
			}
			cf := cloudformation.NewFromConfig(cfg)
			res, err := cf.GetTemplate(context.Background(), &cloudformation.GetTemplateInput{
				StackName: aws.String(stackName),
			})
			if err != nil {
				fmt.Printf("Failed to get stack template: %v\n", err)
				os.Exit(1)
			}
			templateBody = aws.ToString(res.TemplateBody)
		}
		err := os.WriteFile(output, []byte(templateBody), 0644)
		if err != nil {
			fmt.Printf("Failed to write template: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("CloudFormation template written to %s\n", output)
	},
}

func init() {
	getStackCmd.Flags().String("stack-name", "", "Name of the CloudFormation stack")
	getStackCmd.Flags().StringP("output", "o", "", "Output file path")
	cloudformationCmd.AddCommand(getStackCmd)
}
