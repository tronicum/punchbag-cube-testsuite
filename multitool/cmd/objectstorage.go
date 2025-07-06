package cmd

import (
	"fmt"
	"os"
	"punchbag-cube-testsuite/multitool/pkg/mock"
	"punchbag-cube-testsuite/multitool/pkg/models"

	"github.com/spf13/cobra"
)

var objectStorageCmd = &cobra.Command{
	Use:   "objectstorage",
	Short: "Manage S3-like object storage buckets (AWS, Azure, GCP)",
}

var createBucketCmd = &cobra.Command{
	Use:   "create [provider] [name] [region]",
	Short: "Create a new bucket",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		provider, name, region := args[0], args[1], args[2]
		bucket := &models.Bucket{Name: name, Region: region, Provider: provider}
		var store models.ObjectStorage
		switch provider {
		case "aws":
			store = mock.NewAwsObjectStorage()
		case "azure":
			store = mock.NewAzureObjectStorage()
		case "gcp":
			store = mock.NewGcpObjectStorage()
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
		var store models.ObjectStorage
		switch provider {
		case "aws":
			store = mock.NewAwsObjectStorage()
		case "azure":
			store = mock.NewAzureObjectStorage()
		case "gcp":
			store = mock.NewGcpObjectStorage()
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
		var store models.ObjectStorage
		switch provider {
		case "aws":
			store = mock.NewAwsObjectStorage()
		case "azure":
			store = mock.NewAzureObjectStorage()
		case "gcp":
			store = mock.NewGcpObjectStorage()
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
		var store models.ObjectStorage
		switch provider {
		case "aws":
			store = mock.NewAwsObjectStorage()
		case "azure":
			store = mock.NewAzureObjectStorage()
		case "gcp":
			store = mock.NewGcpObjectStorage()
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
	objectStorageCmd.AddCommand(createBucketCmd, listBucketsCmd, getBucketCmd, deleteBucketCmd)
	rootCmd.AddCommand(objectStorageCmd)
}
