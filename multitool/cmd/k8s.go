package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var k8sCmd = &cobra.Command{
	Use:   "k8s",
	Short: "Manage Kubernetes clusters across all providers",
}

// --- Hetzner helpers ---
// List Hetzner networks and return the first available network ID, or create one if none exist
func getOrCreateHetznerNetwork(token, location string) (int, error) {
	req, _ := http.NewRequest("GET", "https://api.hetzner.cloud/v1/networks", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error listing networks: %v", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result struct {
		Networks []struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			IPRange  string `json:"ip_range"`
			Subnets  []struct {
				Type     string `json:"type"`
				NetworkZone string `json:"network_zone"`
				IPRange  string `json:"ip_range"`
			} `json:"subnets"`
		} `json:"networks"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("error parsing networks: %v", err)
	}
	if len(result.Networks) > 0 {
		return result.Networks[0].ID, nil
	}
	// No networks found, create one
	netReq := map[string]interface{}{
		"name": "mt-demo-network",
		"ip_range": "10.0.0.0/16",
		"subnets": []map[string]interface{}{
			{
				"type": "cloud",
				"ip_range": "10.0.1.0/24",
				"network_zone": "eu-central",
			},
		},
	}
	jsonBody, _ := json.Marshal(netReq)
	createReq, _ := http.NewRequest("POST", "https://api.hetzner.cloud/v1/networks", bytes.NewBuffer(jsonBody))
	createReq.Header.Set("Authorization", "Bearer "+token)
	createReq.Header.Set("Content-Type", "application/json")
	createResp, err := http.DefaultClient.Do(createReq)
	if err != nil {
		return 0, fmt.Errorf("error creating network: %v", err)
	}
	defer createResp.Body.Close()
	createBody, _ := ioutil.ReadAll(createResp.Body)
	var createResult struct {
		Network struct {
			ID int `json:"id"`
		} `json:"network"`
	}
	if err := json.Unmarshal(createBody, &createResult); err != nil {
		return 0, fmt.Errorf("error parsing created network: %v", err)
	}
	return createResult.Network.ID, nil
}

// Fetch available Hetzner Kubernetes versions and return the latest
func getLatestHetznerK8sVersion(token string) (string, error) {
	req, _ := http.NewRequest("GET", "https://api.hetzner.cloud/v1/kubernetes_versions", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error fetching versions: %v", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if debugHetzner {
		fmt.Printf("[DEBUG] Hetzner k8s versions raw response:\n%s\n", string(body))
	}
	var result struct {
		KubernetesVersions []struct {
			Version string `json:"version"`
			Supported bool `json:"supported"`
		} `json:"kubernetes_versions"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error parsing versions: %v", err)
	}
	for i := len(result.KubernetesVersions) - 1; i >= 0; i-- {
		if result.KubernetesVersions[i].Supported {
			return result.KubernetesVersions[i].Version, nil
		}
	}
	return "", fmt.Errorf("no supported versions found")
}

var defaultHetznerK8sVersion = "1.29.2" // Update as needed

// Create a Hetzner Cloud Kubernetes cluster
func createHetznerK8sCluster(token, name, location, version string) {
	networkID, err := getOrCreateHetznerNetwork(token, location)
	if err != nil {
		fmt.Printf("Error getting or creating network: %v\n", err)
		return
	}
	if version == "" {
		version = defaultHetznerK8sVersion
	}
	clusterReq := map[string]interface{}{
		"name": name,
		"location": location,
		"network": networkID,
		"version": version,
		"network_zones": []string{"eu-central"},
		"node_pools": []map[string]interface{}{
			{
				"name": "np-cx22",
				"node_count": 1,
				"server_type": "cx22",
				"labels": map[string]string{"test": "ipv6-only"},
				"public_ipv4": false,
				"public_ipv6": true,
			},
		},
	}
	jsonBody, _ := json.Marshal(clusterReq)
	req, _ := http.NewRequest("POST", "https://api.hetzner.cloud/v1/kubernetes_clusters", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error creating Hetzner K8s cluster: %v\n", err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 201 {
		fmt.Println("Hetzner K8s cluster created.")
		fmt.Printf("Response: %s\n", string(body))
	} else {
		fmt.Printf("Cluster creation response: %s\n%s\n", resp.Status, string(body))
	}
}

// Delete a Hetzner Cloud Kubernetes cluster
func deleteHetznerK8sCluster(token string, id int) {
	fmt.Printf("Deleting Hetzner Cloud Kubernetes cluster with ID %d...\n", id)
	// TODO: Implement actual cluster deletion using Hetzner Cloud API
	fmt.Println("[stub] Would call Hetzner Cloud API to delete managed k8s cluster.")
}

// Fetch kubeconfig for a Hetzner Cloud Kubernetes cluster
func fetchHetznerKubeconfig(token string, id int) {
	fmt.Printf("Fetching kubeconfig for Hetzner Cloud Kubernetes cluster with ID %d...\n", id)
	// TODO: Implement actual kubeconfig fetching using Hetzner Cloud API
	fmt.Println("[stub] Would call Hetzner Cloud API to fetch kubeconfig for managed k8s cluster.")
}

func createK8sCluster(provider, name string) {
	switch provider {
	case "aws":
		fmt.Printf("[stub] Would create EKS cluster '%s' on AWS\n", name)
		// TODO: Implement AWS EKS creation logic or proxy
	case "azure":
		fmt.Printf("[stub] Would create AKS cluster '%s' on Azure\n", name)
		// TODO: Implement Azure AKS creation logic or proxy
	case "gcp":
		fmt.Printf("[stub] Would create GKE cluster '%s' on GCP\n", name)
		// TODO: Implement GCP GKE creation logic or proxy
	case "hetzner":
		token := os.Getenv("HCLOUD_TOKEN")
		if token == "" {
			fmt.Println("Hetzner Cloud API token not set. Please set HCLOUD_TOKEN environment variable.")
			return
		}
		fmt.Printf("Creating Hetzner Cloud Kubernetes cluster '%s'...\n", name)
		// TODO: Use hcloud-go or direct API call to create the cluster
		fmt.Println("[stub] Would call Hetzner Cloud API to create managed k8s cluster.")
	case "ionos":
		fmt.Printf("[stub] Would create IONOS Cloud cluster '%s'\n", name)
		// TODO: Implement IONOS cluster creation logic or proxy
	case "stackit":
		fmt.Printf("[stub] Would create STACKIT cluster '%s'\n", name)
		// TODO: Implement STACKIT cluster creation logic or proxy
	default:
		fmt.Printf("Provider '%s' not supported. Supported: aws, azure, gcp, hetzner, ionos, stackit\n", provider)
	}
}

var k8sCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		name, _ := cmd.Flags().GetString("name")
		location, _ := cmd.Flags().GetString("location")
		version, _ := cmd.Flags().GetString("version")
		if provider == "" || name == "" {
			fmt.Println("--provider and --name are required")
			os.Exit(1)
		}
		if provider == "hetzner" {
			token := os.Getenv("HCLOUD_TOKEN")
			if token == "" {
				fmt.Println("Hetzner Cloud API token not set. Please set HCLOUD_TOKEN environment variable.")
				return
			}
			if location == "" {
				location = "fsn1"
			}
			createHetznerK8sCluster(token, name, location, version)
			return
		}
		createK8sCluster(provider, name)
	},
}

var k8sGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details of a Kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		name, _ := cmd.Flags().GetString("name")
		if provider == "" || name == "" {
			fmt.Println("--provider and --name are required")
			os.Exit(1)
		}
		fmt.Printf("[stub] Would get details for cluster '%s' on provider '%s'\n", name, provider)
		// TODO: Implement provider-specific logic and proxy/simulation support
	},
}

// List Hetzner Managed Kubernetes clusters using the REST API
type hcloudK8sCluster struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}
type hcloudK8sListResp struct {
	KubernetesClusters []hcloudK8sCluster `json:"kubernetes_clusters"`
}

// Fetch and print detailed info for each Hetzner K8s cluster
func listHetznerK8sClusters(token string) {
	req, _ := http.NewRequest("GET", "https://api.hetzner.cloud/v1/kubernetes_clusters", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error listing Hetzner K8s clusters: %v\n", err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result hcloudK8sListResp
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		return
	}
	fmt.Println("Hetzner Managed Kubernetes Clusters:")
	for _, c := range result.KubernetesClusters {
		fmt.Printf("- %s (ID: %d, Status: %s)\n", c.Name, c.ID, c.Status)
		// Fetch details for each cluster
		detailReq, _ := http.NewRequest("GET", fmt.Sprintf("https://api.hetzner.cloud/v1/kubernetes_clusters/%d", c.ID), nil)
		detailReq.Header.Set("Authorization", "Bearer "+token)
		detailResp, err := http.DefaultClient.Do(detailReq)
		if err == nil {
			defer detailResp.Body.Close()
			detailBody, _ := ioutil.ReadAll(detailResp.Body)
			fmt.Printf("  Details: %s\n", string(detailBody))
		}
	}
}

// Create a dummy Hetzner Object Storage bucket if not present
func createHetznerDummyBucket(token, name, location string) {
	bucketReq := map[string]string{"name": name, "location": location}
	jsonBody, _ := json.Marshal(bucketReq)
	req, _ := http.NewRequest("POST", "https://api.hetzner.cloud/v1/object_storages", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error creating Hetzner Object Storage bucket: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 201 {
		fmt.Println("Dummy bucket created.")
	} else {
		fmt.Printf("Bucket creation response: %s\n", resp.Status)
	}
}

// Fetch and print detailed info for each Hetzner Object Storage bucket
func listHetznerObjectStorage(token string) {
	bucketName := "dummy-bucket-mt"
	location := "fsn1"
	createHetznerDummyBucket(token, bucketName, location)
	req, _ := http.NewRequest("GET", "https://api.hetzner.cloud/v1/object_storages", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error listing Hetzner Object Storage buckets: %v\n", err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result hcloudObjectStorageListResp
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		return
	}
	fmt.Println("Hetzner Object Storage Buckets:")
	for _, b := range result.ObjectStorages {
		fmt.Printf("- %s (ID: %d, Location: %s)\n", b.Name, b.ID, b.Location)
		// Fetch details for each bucket
		detailReq, _ := http.NewRequest("GET", fmt.Sprintf("https://api.hetzner.cloud/v1/object_storages/%d", b.ID), nil)
		detailReq.Header.Set("Authorization", "Bearer "+token)
		detailResp, err := http.DefaultClient.Do(detailReq)
		if err == nil {
			defer detailResp.Body.Close()
			detailBody, _ := ioutil.ReadAll(detailResp.Body)
			fmt.Printf("  Details: %s\n", string(detailBody))
		}
	}
}

// Attempt to fetch pricing info (Hetzner does not provide a pricing API, so print a message)
func printHetznerPricingInfo() {
	fmt.Println("Hetzner does not provide a public pricing API. See https://www.hetzner.com/cloud for current prices.")
}

var k8sListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Kubernetes clusters",
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		if provider == "hetzner" {
			token := os.Getenv("HCLOUD_TOKEN")
			if token == "" {
				fmt.Println("Hetzner Cloud API token not set. Please set HCLOUD_TOKEN environment variable.")
				return
			}
			listHetznerK8sClusters(token)
			listHetznerObjectStorage(token)
			return
		}
		fmt.Printf("[stub] Would list clusters for provider '%s'\n", provider)
		// TODO: Implement provider-specific logic and proxy/simulation support
	},
}

var k8sDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		id, _ := cmd.Flags().GetInt("id")
		if provider == "" || id == 0 {
			fmt.Println("--provider and --id are required")
			os.Exit(1)
		}
		if provider == "hetzner" {
			token := os.Getenv("HCLOUD_TOKEN")
			if token == "" {
				fmt.Println("Hetzner Cloud API token not set. Please set HCLOUD_TOKEN environment variable.")
				return
			}
			deleteHetznerK8sCluster(token, id)
			return
		}
		fmt.Printf("[stub] Would delete cluster with ID %d for provider '%s'\n", id, provider)
	},
}

var k8sKubeconfigCmd = &cobra.Command{
	Use:   "kubeconfig",
	Short: "Fetch kubeconfig for a Kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		id, _ := cmd.Flags().GetInt("id")
		if provider == "" || id == 0 {
			fmt.Println("--provider and --id are required")
			os.Exit(1)
		}
		if provider == "hetzner" {
			token := os.Getenv("HCLOUD_TOKEN")
			if token == "" {
				fmt.Println("Hetzner Cloud API token not set. Please set HCLOUD_TOKEN environment variable.")
				return
			}
			fetchHetznerKubeconfig(token, id)
			return
		}
		fmt.Printf("[stub] Would fetch kubeconfig for cluster with ID %d for provider '%s'\n", id, provider)
	},
}

var debugHetzner bool

func init() {
	k8sCreateCmd.Flags().String("provider", "", "Cloud provider (aws|azure|gcp|hetzner|ionos|stackit)")
	k8sCreateCmd.Flags().String("name", "", "Cluster name")
	k8sCreateCmd.Flags().String("location", "", "Cluster location (for Hetzner, e.g., fsn1)")
	k8sCreateCmd.Flags().String("version", "", "Kubernetes version (Hetzner, default: latest supported)")
	k8sGetCmd.Flags().String("provider", "", "Cloud provider (aws|azure|gcp|hetzner|ionos|stackit)")
	k8sGetCmd.Flags().String("name", "", "Cluster name")
	k8sListCmd.Flags().String("provider", "", "Cloud provider (optional)")
	k8sDeleteCmd.Flags().String("provider", "", "Cloud provider (aws|azure|gcp|hetzner|ionos|stackit)")
	k8sDeleteCmd.Flags().Int("id", 0, "Cluster ID (for Hetzner)")
	k8sKubeconfigCmd.Flags().String("provider", "", "Cloud provider (aws|azure|gcp|hetzner|ionos|stackit)")
	k8sKubeconfigCmd.Flags().Int("id", 0, "Cluster ID (for Hetzner)")

	k8sCmd.AddCommand(k8sCreateCmd)
	k8sCmd.AddCommand(k8sGetCmd)
	k8sCmd.AddCommand(k8sListCmd)
	k8sCmd.AddCommand(k8sDeleteCmd)
	k8sCmd.AddCommand(k8sKubeconfigCmd)

	k8sCmd.PersistentFlags().BoolVar(&debugHetzner, "debug-hetzner", false, "Enable debug output for Hetzner API calls")
}

// Object storage types for Hetzner API
type hcloudObjectStorage struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Location string `json:"location"`
}
type hcloudObjectStorageListResp struct {
	ObjectStorages []hcloudObjectStorage `json:"object_storages"`
}
