package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"github.com/tronicum/punchbag-cube-testsuite/multitool/pkg/mock"
	sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"

	"github.com/spf13/cobra"
)

var objectStorageCmd = &cobra.Command{
	Use:   "objectstorage",
	Short: "Manage S3-like object storage buckets (AWS, Azure, GCP, StackIT, Hetzner, IONOS)",
}

var policyFile string
var versioning bool
var lifecycleFile string

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
		// Local mock mode
		var store interface{ CreateBucket(*sharedmodels.ObjectStorageBucket) (*sharedmodels.ObjectStorageBucket, error) }
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
		provider := args[0]
		if proxyServer != "" {
			url := fmt.Sprintf("%s/api/proxy/%s/objectstorage", proxyServer, provider)
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
			var buckets []map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&buckets); err != nil {
				fmt.Println("Failed to decode proxy response:", err)
				os.Exit(1)
			}
			for _, b := range buckets {
				fmt.Printf("%+v\n", b)
			}
			return
		}
		var store interface{ ListBuckets() ([]*sharedmodels.ObjectStorageBucket, error) }
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
		buckets, err := store.ListBuckets()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		for _, b := range buckets {
			fmt.Printf("%+v\n", b)
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
	objectStorageCmd.AddCommand(createBucketCmd, listBucketsCmd, getBucketCmd, deleteBucketCmd)
	rootCmd.AddCommand(objectStorageCmd)
}
