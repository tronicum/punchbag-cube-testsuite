package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	configloader "github.com/tronicum/punchbag-cube-testsuite/multitool/pkg"
	sharederrors "github.com/tronicum/punchbag-cube-testsuite/shared/errors"
	"github.com/tronicum/punchbag-cube-testsuite/shared/log"
	"github.com/tronicum/punchbag-cube-testsuite/shared/models"
	awsS3 "github.com/tronicum/punchbag-cube-testsuite/shared/providers/aws"
	hetznerS3 "github.com/tronicum/punchbag-cube-testsuite/shared/providers/hetzner"
	"gopkg.in/yaml.v2"
)

var supportedProviders = []string{"aws-s3", "generic-aws-s3", "hetzner"}

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
	   Short: "Manage S3-like object storage buckets (aws-s3, generic-aws-s3, hetzner)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Supported object storage providers:")
		for _, p := range supportedProviders {
			fmt.Println("-", p)
		}
		fmt.Println("Use 'mt objectstorage [command] --help' for more info.")
	},
	Annotations: map[string]string{"group": "Cloud ObjectStorage (S3) Commands"},
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
	   Use:   "create [name] [region]",
	   Short: "Create a new bucket (supports --storage-provider: aws-s3, generic-aws-s3, hetzner)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name, region := args[0], args[1]
			   provider := getProvider(cmd)
		bucket := &models.ObjectStorageBucket{Name: name, Region: region, Provider: models.CloudProvider(provider)}
		mode, _ := cmd.Flags().GetString("mode")
			   if mode == "proxy" && proxyServer != "" {
				   url := fmt.Sprintf("%s/api/v1/proxy/providers/%s/buckets", proxyServer, provider)
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
			   if mode == "simulate" && proxyServer != "" {
				   url := fmt.Sprintf("%s/api/v1/simulate/providers/%s/buckets", proxyServer, provider)
			jsonBody, err := json.Marshal(bucket)
			if err != nil {
				fmt.Println("Failed to marshal bucket:", err)
				os.Exit(1)
			}
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
			if err != nil {
				fmt.Println("Simulation request failed:", err)
				os.Exit(1)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusCreated {
				fmt.Println("Simulation server error:", resp.Status)
				os.Exit(1)
			}
			var result map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				fmt.Println("Failed to decode simulation response:", err)
				os.Exit(1)
			}
			fmt.Printf("Created bucket (simulation): %+v\n", result)
			return
		}
	   var err error
	   // In simulate mode, skip real credential checks and use dummy values
	   if mode == "simulate" && (provider == "hetzner" || provider == "aws-s3" || provider == "generic-aws-s3") {
		   if os.Getenv("SIMULATE_DUMMY_S3_CREDS") == "1" {
			   log.Info("[SIMULATE] Using dummy S3 credentials (SIMULATE_DUMMY_S3_CREDS=1)")
			   log.Info("Created bucket (simulate): %+v", bucket)
			   return
		   }
	   }
	   switch provider {
	   case "aws-s3", "generic-aws-s3":
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
		   if err == sharederrors.ErrConflict {
			   log.Warn("Bucket already exists: %v", err)
		   } else if err == sharederrors.ErrValidation {
			   log.Error("Validation failed: %v", err)
		   } else {
			   log.Error("Error creating bucket: %v", err)
		   }
		   os.Exit(1)
	   }
	   log.Info("Created bucket: %+v", bucket)
	},
}

var listBucketsCmd = &cobra.Command{
	   Use:   "list",
	   Short: "List all buckets for a provider (use --storage-provider: aws-s3, generic-aws-s3, hetzner)",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		provider := getProvider(cmd)
		mode, _ := cmd.Flags().GetString("mode")
		var buckets []models.ObjectStorageBucket
		var err error
		if mode == "proxy" && proxyServer != "" {
			   url := fmt.Sprintf("%s/api/v1/proxy/providers/%s/buckets", proxyServer, provider)
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("Proxy request failed:", err)
				os.Exit(1)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				fmt.Println("Proxy server error:", resp.Status)
				os.Exit(1)
			}
			if err := json.NewDecoder(resp.Body).Decode(&buckets); err != nil {
				fmt.Println("Failed to decode proxy response:", err)
				os.Exit(1)
			}
		} else if mode == "simulate" && proxyServer != "" {
			   url := fmt.Sprintf("%s/api/v1/simulate/providers/%s/buckets", proxyServer, provider)
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("Simulation request failed:", err)
				os.Exit(1)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				fmt.Println("Simulation server error:", resp.Status)
				os.Exit(1)
			}
			if err := json.NewDecoder(resp.Body).Decode(&buckets); err != nil {
				fmt.Println("Failed to decode simulation response:", err)
				os.Exit(1)
			}
	   } else {
		   // In simulate mode, skip real credential checks and use dummy values
		   if mode == "simulate" && (provider == "hetzner" || provider == "aws-s3" || provider == "generic-aws-s3") {
			   if os.Getenv("SIMULATE_DUMMY_S3_CREDS") == "1" {
				   log.Info("[SIMULATE] Using dummy S3 credentials (SIMULATE_DUMMY_S3_CREDS=1)")
				   // Return a dummy bucket list for simulation
				   buckets = []models.ObjectStorageBucket{
					   {
						   Name:      "sim-bucket-1",
						   Provider:  models.CloudProvider(provider),
						   Region:    "sim-region-1",
						   CreatedAt: time.Now().UTC(),
					   },
				   }
			   } else {
				   log.Info("[SIMULATE] SIMULATE_DUMMY_S3_CREDS not set, falling back to real credential logic.")
			   }
		   } else {
			   switch provider {
			   case "aws-s3", "generic-aws-s3":
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
				   if err.Error() == "InvalidAccessKeyId" || (err.Error() != "" && (contains(err.Error(), "InvalidAccessKeyId") || contains(err.Error(), "403"))) {
					   fmt.Println("Error: Invalid S3 access key or secret. Please check your credentials and environment variables.")
				   } else {
					   fmt.Println("Error:", err)
				   }
				   os.Exit(1)
			   }
		   }
	   }
		// contains is a helper to check if a substring exists in a string

		// ...existing code...

		// ...existing code...

		// ...existing code...
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
	   Use:   "get [id]",
	   Short: "Get a bucket by ID (use --storage-provider)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		provider := getProvider(cmd)
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
	   Use:   "delete [id]",
	   Short: "Delete a bucket by ID (use --storage-provider: aws-s3, generic-aws-s3, hetzner)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	   id := args[0]
	   provider := getProvider(cmd)
	   mode, _ := cmd.Flags().GetString("mode")
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
	   // mode already declared above
	   if proxyServer != "" {
		   client := &http.Client{}
		   var url string
		   if mode == "simulate" {
			   url = fmt.Sprintf("%s/api/v1/simulate/providers/%s/buckets/%s", proxyServer, provider, id)
		   } else if mode == "proxy" {
			   url = fmt.Sprintf("%s/api/v1/proxy/providers/%s/buckets/%s", proxyServer, provider, id)
		   } else {
			   // fallback for direct mode or error
			   url = ""
		   }
		   req, err := http.NewRequest(http.MethodDelete, url, nil)
		   if err != nil {
			   printSuccessOrError(false, automationMode, "Failed to create request: %v", err)
			   os.Exit(1)
		   }
		   resp, err := client.Do(req)
		   if err != nil {
			   printSuccessOrError(false, automationMode, "Request failed: %v", err)
			   os.Exit(1)
		   }
		   defer resp.Body.Close()
		   if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			   printSuccessOrError(false, automationMode, "Server error: %s", resp.Status)
			   os.Exit(1)
		   }
		   if mode == "simulate" {
			   printSuccessOrError(true, automationMode, "Bucket deleted successfully (simulate).")
		   } else {
			   printSuccessOrError(true, automationMode, "Bucket deleted successfully (proxy).")
		   }
		   return
	   }
	   var err error
	   // mode already declared above
	   // In simulate mode, skip real credential checks and use dummy values
	   if mode == "simulate" && (provider == "hetzner" || provider == "aws-s3" || provider == "generic-aws-s3") {
		   if os.Getenv("SIMULATE_DUMMY_S3_CREDS") == "1" {
			   log.Info("[SIMULATE] Using dummy S3 credentials (SIMULATE_DUMMY_S3_CREDS=1)")
			   printSuccessOrError(true, automationMode, "Bucket deleted successfully (simulate).")
			   return
		   }
	   }
	   switch provider {
	   case "aws-s3", "generic-aws-s3":
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
	   // Register persistent flags before adding subcommands so all subcommands inherit them
	   objectStorageCmd.PersistentFlags().StringVar(&hetznerToken, "hetzner-token", "", "Hetzner API token (overrides env/config)")
	   objectStorageCmd.PersistentFlags().BoolVar(&skipPrompts, "skip-prompts", false, "Skip all interactive prompts (for automation)")
	   objectStorageCmd.PersistentFlags().BoolVar(&automationMode, "automation-mode", false, "Enable automation/CI mode (alias for --skip-prompts --force, prints colored output)")
	   objectStorageCmd.PersistentFlags().String("storage-provider", "hetzner", "Object storage provider: aws-s3, generic-aws-s3, hetzner")

	   // Register --mode as a local flag on each subcommand
	   createBucketCmd.Flags().String("mode", "direct", "Operation mode: direct|proxy|simulate")
	   listBucketsCmd.Flags().String("mode", "direct", "Operation mode: direct|proxy|simulate")
	   getBucketCmd.Flags().String("mode", "direct", "Operation mode: direct|proxy|simulate")
	   deleteBucketCmd.Flags().String("mode", "direct", "Operation mode: direct|proxy|simulate")

	   // Subcommand-specific flags
	   createBucketCmd.Flags().StringVar(&policyFile, "policy", "", "Path to bucket policy JSON file")
	   createBucketCmd.Flags().BoolVar(&versioning, "versioning", false, "Enable versioning")
	   createBucketCmd.Flags().StringVar(&lifecycleFile, "lifecycle", "", "Path to lifecycle rules JSON file")
	   listBucketsCmd.Flags().StringVarP(&objectStorageOutputFormat, "output", "o", "table", "Output format: json or table")
	   deleteBucketCmd.Flags().BoolVarP(&forceDelete, "force", "f", false, "Force deletion without confirmation")

	   objectStorageCmd.AddCommand(createBucketCmd, listBucketsCmd, getBucketCmd, deleteBucketCmd, supportedProvidersCmd)
}

// Helper to get provider from flag or config
var defaultProvider string // TODO: load from config/env
func getProvider(cmd *cobra.Command) string {
	   p, _ := cmd.Flags().GetString("storage-provider")
	   if p != "" {
			   return p
	   }
	   prof, _ := cmd.Flags().GetString("profile")
	   cfg, err := configloader.LoadMTConfig(prof)
	   if err == nil && cfg.Provider != "" {
			   return cfg.Provider
	   }
	   if defaultProvider != "" {
			   return defaultProvider
	   }
	   return "hetzner" // fallback default
}
func getRegion(cmd *cobra.Command) string {
	prof, _ := cmd.Flags().GetString("profile")
	cfg, err := configloader.LoadMTConfig(prof)
	if err == nil && cfg.Region != "" {
		return cfg.Region
	}
	return "us-east-1"
}

func maskToken(token string) string {
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + "****" + token[len(token)-4:]
}

// Helper: contains returns true if substr is in s
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
