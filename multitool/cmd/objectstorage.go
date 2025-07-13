
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/tronicum/punchbag-cube-testsuite/shared/models"
	awsS3 "github.com/tronicum/punchbag-cube-testsuite/shared/providers/aws"
	hetznerS3 "github.com/tronicum/punchbag-cube-testsuite/shared/providers/hetzner"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var supportedProviders = []string{"aws", "hetzner"}

var supportedProvidersCmd = &cobra.Command{
	Use:   "supported-providers",
	Short: "List supported object storage providers",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Supported object storage providers:")
		for _, p := range supportedProviders {
			fmt.Println("-", p)
		}
	},
}


var objectStorageCmd = &cobra.Command{
	Use:   "objectstorage",
   Short: "Manage S3-like object storage buckets (AWS, Hetzner)",
   Run: func(cmd *cobra.Command, args []string) {
	   fmt.Println("Supported object storage providers:")
	   for _, p := range supportedProviders {
		   fmt.Println("-", p)
	   }
	   fmt.Println("Use 'mt objectstorage [command] --help' for more info.")
   },
}

var policyFile string
var versioning bool
var lifecycleFile string
var hetznerToken string              // CLI flag for Hetzner API token
var objectStorageOutputFormat string // "json" or "table"
var forceDelete bool
var skipPrompts bool
var automationMode bool

var createBucketCmd = &cobra.Command{
	Use:   "create [provider] [name] [region]",
	Short: "Create a new bucket (supports --policy, --versioning, --lifecycle)",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		provider, name, region := args[0], args[1], args[2]
		bucket := &models.ObjectStorageBucket{Name: name, Region: region, Provider: models.CloudProvider(provider)}
		if proxyServer != "" {
			url := fmt.Sprintf("%s/api/proxy/%s/objectstorage", proxyServer, provider)
			jsonBody, err := json.Marshal(bucket)
			if err != nil {
				fmt.Println("Failed to marshal bucket:", err)
				os.Exit(1)
			}
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
			if err != nil {
				fmt.Println("Proxy request failed:", err)
				os.Exit(1)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusCreated {
				fmt.Println("Proxy server error:", resp.Status)
				os.Exit(1)
			}
			var result map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				fmt.Println("Failed to decode proxy response:", err)
				os.Exit(1)
			}
			fmt.Printf("Created bucket (proxy): %+v\n", result)
			return
		}
		var err error
		switch provider {
		case "aws":
			s3Client, e := awsS3.NewS3Client(cmd.Context(), region)
			if e != nil {
				fmt.Println("AWS S3 client error:", e)
				os.Exit(1)
			}
			err = s3Client.CreateBucket(cmd.Context(), bucket)
		case "hetzner":
			hetznerClient := hetznerS3.NewHetznerS3Client()
			err = hetznerClient.CreateBucket(cmd.Context(), bucket)
		default:
			fmt.Println("Provider not implemented in shared library yet.")
			os.Exit(1)
		}
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Printf("Created bucket: %+v\n", bucket)
	},
}

var listBucketsCmd = &cobra.Command{
	Use:   "list [provider]",
	Short: "List all buckets for a provider",
	Args:  cobra.ExactArgs(1),
	   Run: func(cmd *cobra.Command, args []string) {
			   provider := args[0]
			   var buckets []models.ObjectStorageBucket
			   var err error
			   switch provider {
			   case "aws":
					   s3Client, e := awsS3.NewS3Client(cmd.Context(), "us-east-1") // TODO: region flag
					   if e != nil {
							   fmt.Println("AWS S3 client error:", e)
							   os.Exit(1)
					   }
					   buckets, err = s3Client.ListBuckets(cmd.Context())
			   case "hetzner":
					   hetznerClient := hetznerS3.NewHetznerS3Client()
					   buckets, err = hetznerClient.ListBuckets(cmd.Context())
			   default:
					   fmt.Println("Provider not implemented in shared library yet.")
					   os.Exit(1)
			   }
			   if err != nil {
					   fmt.Println("Error:", err)
					   os.Exit(1)
			   }
			   // Output filtering
			   switch objectStorageOutputFormat {
			   case "json":
					   jsonData, err := json.MarshalIndent(buckets, "", "  ")
					   if err != nil {
							   fmt.Println("Error marshaling buckets to JSON:", err)
					   } else {
							   fmt.Println(string(jsonData))
					   }
			   case "yaml":
					   yamlData, err := yaml.Marshal(buckets)
					   if err != nil {
							   fmt.Println("Error marshaling buckets to YAML:", err)
					   } else {
							   fmt.Println(string(yamlData))
					   }
			   default:
					   fmt.Printf("[CLI] listBucketsCmd called with provider: %s\n", provider)
					   fmt.Println("[CLI] Before provider switch")
					   fmt.Printf("%-30s %-15s %-25s\n", "Name", "Region", "CreatedAt")
					   for _, b := range buckets {
							   fmt.Printf("%-30s %-15s %-25s\n", b.Name, b.Region, b.CreatedAt.Format(time.RFC3339))
					   }
			   }
	   },
}

var getBucketCmd = &cobra.Command{
	Use:   "get [provider] [id]",
	Short: "Get a bucket by ID",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		provider, id := args[0], args[1]
		if proxyServer != "" {
			url := fmt.Sprintf("%s/api/proxy/%s/objectstorage?id=%s", proxyServer, provider, id)
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("Proxy request failed:", err)
				os.Exit(1)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
				fmt.Println("Proxy server error:", resp.Status)
				os.Exit(1)
			}
			var bucket map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&bucket); err != nil {
				fmt.Println("Failed to decode proxy response:", err)
				os.Exit(1)
			}
			fmt.Printf("Bucket (proxy): %+v\n", bucket)
			return
		}
		fmt.Println("GetBucket not implemented in shared library yet.")
		os.Exit(1)
	},
}

var deleteBucketCmd = &cobra.Command{
	Use:   "delete [provider] [id]",
	Short: "Delete a bucket by ID",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
	   provider, id := strings.ToLower(args[0]), args[1]
		// Automation/skip-prompts logic
		if automationMode {
			skipPrompts = true
			forceDelete = true
		}
		if !forceDelete && !skipPrompts {
			fmt.Printf("Are you sure you want to delete bucket '%s'? (Y/n): ", id)
			var response string
			fmt.Scanln(&response)
			if response != "Y" && response != "y" && response != "" {
				printSuccessOrError(false, automationMode, "Aborted.")
				return
			}
		}
	   if proxyServer != "" {
			   client := &http.Client{}
			   url := fmt.Sprintf("%s/api/proxy/%s/objectstorage?id=%s", proxyServer, provider, id)
			   req, err := http.NewRequest(http.MethodDelete, url, nil)
			   if err != nil {
					   printSuccessOrError(false, automationMode, "Failed to create proxy request: %v", err)
					   os.Exit(1)
			   }
			   resp, err := client.Do(req)
			   if err != nil {
					   printSuccessOrError(false, automationMode, "Proxy request failed: %v", err)
					   os.Exit(1)
			   }
			   defer resp.Body.Close()
			   if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
					   printSuccessOrError(false, automationMode, "Proxy server error: %s", resp.Status)
					   os.Exit(1)
			   }
			   printSuccessOrError(true, automationMode, "Bucket deleted successfully (proxy).")
			   return
	   }
	   var err error
	   switch provider {
	   case "aws":
			   s3Client, e := awsS3.NewS3Client(cmd.Context(), "us-east-1") // TODO: region flag
			   if e != nil {
					   printSuccessOrError(false, automationMode, "AWS S3 client error: %v", e)
					   os.Exit(1)
			   }
			   err = s3Client.DeleteBucket(cmd.Context(), id)
	   case "hetzner":
			   hetznerClient := hetznerS3.NewHetznerS3Client()
			   err = hetznerClient.DeleteBucket(cmd.Context(), id)
	   default:
			   printSuccessOrError(false, automationMode, "Provider not implemented in shared library yet.")
			   os.Exit(1)
	   }
	   if err != nil {
			   printSuccessOrError(false, automationMode, "Error deleting bucket: %v", err)
			   os.Exit(1)
	   }
	   printSuccessOrError(true, automationMode, "Bucket deleted successfully.")
	},
}

func printSuccessOrError(success, automation bool, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	if automation {
		if success {
			fmt.Printf("\033[32m%s\033[0m\n", msg) // green
		} else {
			fmt.Printf("\033[31m%s\033[0m\n", msg) // red
		}
	} else {
		fmt.Println(msg)
	}
}

func init() {
	createBucketCmd.Flags().StringVar(&policyFile, "policy", "", "Path to bucket policy JSON file")
	createBucketCmd.Flags().BoolVar(&versioning, "versioning", false, "Enable versioning")
	createBucketCmd.Flags().StringVar(&lifecycleFile, "lifecycle", "", "Path to lifecycle rules JSON file")
	createBucketCmd.Flags().StringVar(&proxyServer, "server", "", "Proxy server URL")
	createBucketCmd.Flags().StringVar(&hetznerToken, "hetzner-token", "", "Hetzner API token (overrides env/config)")
	listBucketsCmd.Flags().StringVarP(&objectStorageOutputFormat, "output", "o", "table", "Output format: json or table")
	deleteBucketCmd.Flags().BoolVarP(&forceDelete, "force", "f", false, "Force deletion without confirmation")
	objectStorageCmd.PersistentFlags().BoolVar(&skipPrompts, "skip-prompts", false, "Skip all interactive prompts (for automation)")
	objectStorageCmd.PersistentFlags().BoolVar(&automationMode, "automation-mode", false, "Enable automation/CI mode (alias for --skip-prompts --force, prints colored output)")
   objectStorageCmd.AddCommand(createBucketCmd, listBucketsCmd, getBucketCmd, deleteBucketCmd, supportedProvidersCmd)
}

func maskToken(token string) string {
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + "****" + token[len(token)-4:]
}
