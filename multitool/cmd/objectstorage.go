package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/tronicum/punchbag-cube-testsuite/multitool/pkg/mock"
	sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"
	"github.com/tronicum/punchbag-cube-testsuite/multitool/pkg/client"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var objectStorageCmd = &cobra.Command{
	Use:   "objectstorage",
	Short: "Manage S3-like object storage buckets (AWS, Azure, GCP, StackIT, Hetzner, IONOS)",
}

var policyFile string
var versioning bool
var lifecycleFile string
var hetznerToken string // CLI flag for Hetzner API token
var objectStorageOutputFormat string // "json" or "table"

var createBucketCmd = &cobra.Command{
	Use:   "create [provider] [name] [region]",
	Short: "Create a new bucket (supports --policy, --versioning, --lifecycle)",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		provider, name, region := args[0], args[1], args[2]
		bucket := &sharedmodels.ObjectStorageBucket{Name: name, Region: region, Provider: sharedmodels.CloudProvider(provider)}
		// Advanced S3 features
		if policyFile != "" {
			// Optionally parse policy file here
		}
		if lifecycleFile != "" {
			// Optionally parse lifecycle rules from file
		}
		if proxyServer != "" {
			// Proxy mode: send to cube-server/sim-server
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
		// Local/mock or real provider mode
		var store interface{ CreateBucket(*sharedmodels.ObjectStorageBucket) (*sharedmodels.ObjectStorageBucket, error); ListBuckets() ([]*sharedmodels.ObjectStorageBucket, error) }
		switch provider {
		case "aws":
			// TODO: Use AWS SDK (already implemented)
			store = mock.NewAwsObjectStorage()
		case "azure":
			// TODO: Use Azure SDK for Go for real Azure Blob Storage
			store = mock.NewAzureObjectStorage()
		case "gcp":
			// TODO: Use Google Cloud Storage SDK for Go
			store = mock.NewGcpObjectStorage()
		case "stackit":
			// TODO: Use StackIT SDK or S3-compatible if available
			store = mock.NewStackITObjectStorage()
		case "hetzner":
			// Hetzner S3: use env credentials, not token
			accessKey, secretKey := client.LoadHetznerS3Credentials()
			if accessKey != "" && secretKey != "" {
				fmt.Println("[CLI] Using real Hetzner S3 client (access key and secret found in env)")
				store = client.NewHetznerObjectStorageClientFromKeys(accessKey, secretKey)
			} else {
				fmt.Println("[CLI] ERROR: No Hetzner S3 access key/secret found in env. Using mock implementation.")
				store = mock.NewHetznerObjectStorage()
			}
		case "hetzner-rest":
			fmt.Println("[CLI] ERROR: Hetzner REST API does not support object storage management. Use the S3-compatible API and access keys.")
			os.Exit(1)
		case "ionos":
			// TODO: Use IONOS SDK or S3-compatible if available
			store = mock.NewIonosObjectStorage()
		default:
			fmt.Println("Unknown provider")
			os.Exit(1)
		}
		b, err := store.CreateBucket(bucket)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Printf("Created bucket: %+v\n", b)
	},
}

var listBucketsCmd = &cobra.Command{
	Use:   "list [provider]",
	Short: "List all buckets for a provider",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("[CLI] listBucketsCmd called with provider: %s\n", args[0])
		provider := args[0]
		fmt.Println("[CLI] Before provider switch")
		var store interface{ ListBuckets() ([]*sharedmodels.ObjectStorageBucket, error) }
		switch provider {
		case "aws":
			fmt.Println("[CLI] Using AWS mock")
			store = mock.NewAwsObjectStorage()
		case "azure":
			fmt.Println("[CLI] Using Azure mock")
			store = mock.NewAzureObjectStorage()
		case "gcp":
			fmt.Println("[CLI] Using GCP mock")
			store = mock.NewGcpObjectStorage()
		case "stackit":
			fmt.Println("[CLI] Using StackIT mock")
			store = mock.NewStackITObjectStorage()
		case "hetzner":
			fmt.Println("[CLI] Entered hetzner case")
			// Hetzner S3: use env credentials, not token
			accessKey, secretKey := client.LoadHetznerS3Credentials()
			if accessKey != "" && secretKey != "" {
				fmt.Println("[CLI] Using real Hetzner S3 client (access key and secret found in env)")
				store = client.NewHetznerObjectStorageClientFromKeys(accessKey, secretKey)
			} else {
				fmt.Println("[CLI] ERROR: No Hetzner S3 access key/secret found in env. Using mock implementation.")
				store = mock.NewHetznerObjectStorage()
			}
		case "hetzner-rest":
			fmt.Println("[CLI] Entered hetzner-rest case")
			token, _ := client.LoadHetznerAPIToken(hetznerToken)
			if token != "" {
				fmt.Println("[CLI] Using Hetzner REST API client (token found)")
				store = client.NewHetznerRESTObjectStorageClient(token)
			} else {
				fmt.Println("[CLI] ERROR: No Hetzner API token found. Using mock implementation.")
				store = mock.NewHetznerObjectStorage()
			}
		case "ionos":
			fmt.Println("[CLI] Using IONOS mock")
			store = mock.NewIonosObjectStorage()
		default:
			fmt.Println("Unknown provider")
			os.Exit(1)
		}
		fmt.Println("[CLI] After provider switch, about to call ListBuckets")
		buckets, err := store.ListBuckets()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		// Output filtering
		if objectStorageOutputFormat == "json" {
			jsonData, err := json.MarshalIndent(buckets, "", "  ")
			if err != nil {
				fmt.Println("Error marshaling buckets to JSON:", err)
			} else {
				fmt.Println(string(jsonData))
			}
		} else if objectStorageOutputFormat == "yaml" {
			yamlData, err := yaml.Marshal(buckets)
			if err != nil {
				fmt.Println("Error marshaling buckets to YAML:", err)
			} else {
				fmt.Println(string(yamlData))
			}
		} else {
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
		var store interface{ GetBucket(string) (*sharedmodels.ObjectStorageBucket, error) }
		switch provider {
		case "aws":
			store = mock.NewAwsObjectStorage()
		case "azure":
			store = mock.NewAzureObjectStorage()
		case "gcp":
			store = mock.NewGcpObjectStorage()
		case "stackit":
			store = mock.NewStackITObjectStorage()
		case "hetzner":
			store = mock.NewHetznerObjectStorage()
		case "ionos":
			store = mock.NewIonosObjectStorage()
		default:
			fmt.Println("Unknown provider")
			os.Exit(1)
		}
		b, err := store.GetBucket(id)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Printf("Bucket: %+v\n", b)
	},
}

var deleteBucketCmd = &cobra.Command{
	Use:   "delete [provider] [id]",
	Short: "Delete a bucket by ID",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		provider, id := args[0], args[1]
		if proxyServer != "" {
			client := &http.Client{}
			url := fmt.Sprintf("%s/api/proxy/%s/objectstorage?id=%s", proxyServer, provider, id)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			if err != nil {
				fmt.Println("Failed to create proxy request:", err)
				os.Exit(1)
			}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Proxy request failed:", err)
				os.Exit(1)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
				fmt.Println("Proxy server error:", resp.Status)
				os.Exit(1)
			}
			fmt.Println("Bucket deleted (proxy).")
			return
		}
		var store interface{ DeleteBucket(string) error }
		switch provider {
		case "aws":
			store = mock.NewAwsObjectStorage()
		case "azure":
			store = mock.NewAzureObjectStorage()
		case "gcp":
			store = mock.NewGcpObjectStorage()
		case "stackit":
			store = mock.NewStackITObjectStorage()
		case "hetzner":
			store = mock.NewHetznerObjectStorage()
		case "ionos":
			store = mock.NewIonosObjectStorage()
		default:
			fmt.Println("Unknown provider")
			os.Exit(1)
		}
		err := store.DeleteBucket(id)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Bucket deleted.")
	},
}

func init() {
	createBucketCmd.Flags().StringVar(&policyFile, "policy", "", "Path to bucket policy JSON file")
	createBucketCmd.Flags().BoolVar(&versioning, "versioning", false, "Enable versioning")
	createBucketCmd.Flags().StringVar(&lifecycleFile, "lifecycle", "", "Path to lifecycle rules JSON file")
	createBucketCmd.Flags().StringVar(&proxyServer, "server", "", "Proxy server URL")
	createBucketCmd.Flags().StringVar(&hetznerToken, "hetzner-token", "", "Hetzner API token (overrides env/config)")
	listBucketsCmd.Flags().StringVarP(&objectStorageOutputFormat, "output", "o", "table", "Output format: json or table")
	objectStorageCmd.AddCommand(createBucketCmd, listBucketsCmd, getBucketCmd, deleteBucketCmd)
	rootCmd.AddCommand(objectStorageCmd)
}

func maskToken(token string) string {
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + "****" + token[len(token)-4:]
}
