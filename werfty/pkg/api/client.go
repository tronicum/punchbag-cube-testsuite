package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"
)

// Werfty represents the API werfty
type Werfty struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewWerfty creates a new API werfty
func NewWerfty(baseURL string) *Werfty {
	return &Werfty{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doRequest performs an HTTP request
func (c *Werfty) doRequest(method, path string, body interface{}) (*http.Response, error) {
	url := c.BaseURL + path

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.HTTPClient.Do(req)
}

// ListClusters lists all clusters (multi-cloud)
func (c *Werfty) ListClusters() ([]*sharedmodels.Cluster, error) {
	resp, err := c.doRequest("GET", "/api/v1/clusters", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	var response sharedmodels.ClustersResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return response.Clusters, nil
}

// ListClustersByProvider lists clusters filtered by provider
func (c *Werfty) ListClustersByProvider(provider string) ([]*sharedmodels.Cluster, error) {
	resp, err := c.doRequest("GET", "/api/v1/clusters?provider="+provider, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	var response sharedmodels.ClustersResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return response.Clusters, nil
}

// ListAKSClusters lists all AKS clusters (backward compatibility)
func (c *Werfty) ListAKSClusters() ([]*sharedmodels.AKSCluster, error) {
	clusters, err := c.ListClustersByProvider("azure")
	if err != nil {
		return nil, err
	}

	// Convert to AKSCluster format for backward compatibility
	aksClusters := make([]*sharedmodels.AKSCluster, len(clusters))
	for i, cluster := range clusters {
		aksClusters[i] = c.convertToAKSCluster(cluster)
	}

	return aksClusters, nil
}

// Helper method to convert Cluster to AKSCluster
func (c *Werfty) convertToAKSCluster(cluster *sharedmodels.Cluster) *sharedmodels.AKSCluster {
	aksCluster := &sharedmodels.AKSCluster{
		ID:        cluster.ID,
		Name:      cluster.Name,
		Status:    cluster.Status,
		CreatedAt: cluster.CreatedAt,
		UpdatedAt: cluster.UpdatedAt,
	}

	// Extract Azure-specific fields from config
	if config, ok := cluster.Config["azure_config"].(map[string]interface{}); ok {
		if rg, ok := config["resource_group"].(string); ok {
			aksCluster.ResourceGroup = rg
		}
		if loc, ok := config["location"].(string); ok {
			aksCluster.Location = loc
		}
		if k8sVer, ok := config["kubernetes_version"].(string); ok {
			aksCluster.KubernetesVersion = k8sVer
		}
		if nodeCount, ok := config["node_count"].(float64); ok {
			aksCluster.NodeCount = int(nodeCount)
		}
	}

	return aksCluster
}

// GetCluster gets a cluster by ID (multi-cloud)
func (c *Werfty) GetCluster(id string) (*sharedmodels.Cluster, error) {
	resp, err := c.doRequest("GET", "/api/v1/clusters/"+id, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("cluster not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	var cluster sharedmodels.Cluster
	if err := json.NewDecoder(resp.Body).Decode(&cluster); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &cluster, nil
}

// GetAKSCluster gets an AKS cluster by ID (backward compatibility)
func (c *Werfty) GetAKSCluster(id string) (*sharedmodels.AKSCluster, error) {
	cluster, err := c.GetCluster(id)
	if err != nil {
		return nil, err
	}

	// Convert to AKSCluster format for backward compatibility
	aksCluster := &sharedmodels.AKSCluster{
		ID:        cluster.ID,
		Name:      cluster.Name,
		Status:    cluster.Status,
		CreatedAt: cluster.CreatedAt,
		UpdatedAt: cluster.UpdatedAt,
	}

	// Extract Azure-specific fields from config
	if config, ok := cluster.Config["azure_config"].(map[string]interface{}); ok {
		if rg, ok := config["resource_group"].(string); ok {
			aksCluster.ResourceGroup = rg
		}
		if loc, ok := config["location"].(string); ok {
			aksCluster.Location = loc
		}
		if k8sVer, ok := config["kubernetes_version"].(string); ok {
			aksCluster.KubernetesVersion = k8sVer
		}
		if nodeCount, ok := config["node_count"].(float64); ok {
			aksCluster.NodeCount = int(nodeCount)
		}
	}

	return aksCluster, nil
}

// CreateCluster creates a new cluster (multi-cloud)
func (c *Werfty) CreateMultiCloudCluster(cluster *sharedmodels.Cluster) (*sharedmodels.Cluster, error) {
	resp, err := c.doRequest("POST", "/api/v1/clusters", cluster)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	var created sharedmodels.Cluster
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &created, nil
}

// CreateStackITCluster creates a new StackIT cluster
func (c *Werfty) CreateStackITCluster(name, projectID, region string, config map[string]interface{}) (*sharedmodels.Cluster, error) {
	cluster := &sharedmodels.Cluster{
		Name:          name,
		CloudProvider: "schwarz-stackit",
		Status:        "creating",
		Config: map[string]interface{}{
			"stackit_config": map[string]interface{}{
				"project_id": projectID,
				"region":     region,
			},
		},
	}

	// Merge additional config
	if config != nil {
		if stackitConfig, ok := cluster.Config["stackit_config"].(map[string]interface{}); ok {
			for k, v := range config {
				stackitConfig[k] = v
			}
		}
	}

	return c.CreateMultiCloudCluster(cluster)
}

// CreateAzureCluster creates a new Azure cluster
func (c *Werfty) CreateAzureCluster(name, resourceGroup, location string, config map[string]interface{}) (*sharedmodels.Cluster, error) {
	cluster := &sharedmodels.Cluster{
		Name:          name,
		CloudProvider: "azure",
		Status:        "creating",
		Config: map[string]interface{}{
			"azure_config": map[string]interface{}{
				"resource_group": resourceGroup,
				"location":       location,
			},
		},
	}

	// Merge additional config
	if config != nil {
		if azureConfig, ok := cluster.Config["azure_config"].(map[string]interface{}); ok {
			for k, v := range config {
				azureConfig[k] = v
			}
		}
	}

	return c.CreateMultiCloudCluster(cluster)
}

// CreateHetznerCluster creates a new Hetzner Cloud cluster
func (c *Werfty) CreateHetznerCluster(name, location string, config map[string]interface{}) (*sharedmodels.Cluster, error) {
	cluster := &sharedmodels.Cluster{
		Name:          name,
		CloudProvider: "hetzner-hcloud",
		Status:        "creating",
		Config: map[string]interface{}{
			"hetzner_config": map[string]interface{}{
				"location": location,
			},
		},
	}
	// Merge additional config
	if config != nil {
		if hetznerConfig, ok := cluster.Config["hetzner_config"].(map[string]interface{}); ok {
			for k, v := range config {
				hetznerConfig[k] = v
			}
			// Ensure kubernetes_version is set (default to 1.28.0 if missing or empty)
			if _, ok := hetznerConfig["kubernetes_version"]; !ok || hetznerConfig["kubernetes_version"] == "" {
				hetznerConfig["kubernetes_version"] = "1.28.0"
			}
		}
	}
	return c.CreateMultiCloudCluster(cluster)
}

// CreateIONOSCluster creates a new IONOS Cloud cluster
func (c *Werfty) CreateIONOSCluster(name, datacenterID string, config map[string]interface{}) (*sharedmodels.Cluster, error) {
	cluster := &sharedmodels.Cluster{
		Name:          name,
		CloudProvider: "united-ionos",
		Status:        "creating",
		Config: map[string]interface{}{
			"ionos_config": map[string]interface{}{
				"datacenter_id": datacenterID,
			},
		},
	}

	// Merge additional config
	if config != nil {
		if ionosConfig, ok := cluster.Config["ionos_config"].(map[string]interface{}); ok {
			for k, v := range config {
				ionosConfig[k] = v
			}
		}
	}

	return c.CreateMultiCloudCluster(cluster)
}

// UpdateCluster updates an existing cluster
func (c *Werfty) UpdateCluster(id string, cluster *sharedmodels.AKSCluster) (*sharedmodels.AKSCluster, error) {
	resp, err := c.doRequest("PUT", "/api/v1/clusters/"+id, cluster)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("cluster not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	var updated sharedmodels.AKSCluster
	if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &updated, nil
}

// DeleteCluster deletes a cluster
func (c *Werfty) DeleteCluster(id string) error {
	resp, err := c.doRequest("DELETE", "/api/v1/clusters/"+id, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("cluster not found")
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	return nil
}

// RunTest runs a test on a cluster (multi-cloud)
func (c *Werfty) RunTest(clusterID string, testReq *sharedmodels.TestRequest) (*sharedmodels.TestResult, error) {
	resp, err := c.doRequest("POST", "/api/v1/clusters/"+clusterID+"/tests", testReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("cluster not found")
	}
	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	var result sharedmodels.TestResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}

// RunAKSTest runs a test on an AKS cluster (backward compatibility)
func (c *Werfty) RunAKSTest(clusterID string, testReq *sharedmodels.AKSTestRequest) (*sharedmodels.AKSTestResult, error) {
	multiTestReq := &sharedmodels.TestRequest{
		ClusterID: testReq.ClusterID,
		TestType:  testReq.TestType,
		Config:    testReq.Config,
	}

	result, err := c.RunTest(clusterID, multiTestReq)
	if err != nil {
		return nil, err
	}

	// Convert to AKSTestResult for backward compatibility
	aksResult := &sharedmodels.AKSTestResult{
		ID:          result.ID,
		ClusterID:   result.ClusterID,
		TestType:    result.TestType,
		Status:      result.Status,
		Duration:    result.Duration,
		Details:     result.Details,
		ErrorMsg:    result.ErrorMsg,
		StartedAt:   result.StartedAt,
		CompletedAt: result.CompletedAt,
	}

	return aksResult, nil
}

// GetTestResult gets a test result by ID (multi-cloud)
func (c *Werfty) GetTestResult(id string) (*sharedmodels.TestResult, error) {
	resp, err := c.doRequest("GET", "/api/v1/tests/"+id, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("test result not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	var result sharedmodels.TestResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}

// GetAKSTestResult gets an AKS test result by ID (backward compatibility)
func (c *Werfty) GetAKSTestResult(id string) (*sharedmodels.AKSTestResult, error) {
	result, err := c.GetTestResult(id)
	if err != nil {
		return nil, err
	}

	// Convert to AKSTestResult for backward compatibility
	aksResult := &sharedmodels.AKSTestResult{
		ID:          result.ID,
		ClusterID:   result.ClusterID,
		TestType:    result.TestType,
		Status:      result.Status,
		Duration:    result.Duration,
		Details:     result.Details,
		ErrorMsg:    result.ErrorMsg,
		StartedAt:   result.StartedAt,
		CompletedAt: result.CompletedAt,
	}

	return aksResult, nil
}

// ListTestResults lists test results for a cluster (multi-cloud)
func (c *Werfty) ListTestResults(clusterID string) ([]*sharedmodels.TestResult, error) {
	resp, err := c.doRequest("GET", "/api/v1/clusters/"+clusterID+"/tests", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	var response sharedmodels.TestResultsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return response.TestResults, nil
}

// ListAKSTestResults lists test results for an AKS cluster (backward compatibility)
func (c *Werfty) ListAKSTestResults(clusterID string) ([]*sharedmodels.AKSTestResult, error) {
	results, err := c.ListTestResults(clusterID)
	if err != nil {
		return nil, err
	}

	// Convert to AKSTestResult for backward compatibility
	aksResults := make([]*sharedmodels.AKSTestResult, len(results))
	for i, result := range results {
		aksResults[i] = &sharedmodels.AKSTestResult{
			ID:          result.ID,
			ClusterID:   result.ClusterID,
			TestType:    result.TestType,
			Status:      result.Status,
			Duration:    result.Duration,
			Details:     result.Details,
			ErrorMsg:    result.ErrorMsg,
			StartedAt:   result.StartedAt,
			CompletedAt: result.CompletedAt,
		}
	}

	return aksResults, nil
}

// ValidateProvider validates a cloud provider configuration
func (c *Werfty) ValidateProvider(provider string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/validate/%s", c.BaseURL, provider)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("validation failed with status: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// GetProviderInfo gets information about a cloud provider
func (c *Werfty) GetProviderInfo(provider string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/providers/%s/info", c.BaseURL, provider)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// ListProviderClusters lists clusters for a specific provider
func (c *Werfty) ListProviderClusters(provider string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/providers/%s/clusters", c.BaseURL, provider)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// ExecuteProviderOperation executes a provider-specific operation
func (c *Werfty) ExecuteProviderOperation(provider, operation, params string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/providers/%s/operations/%s", c.BaseURL, provider, operation)

	// Parse params as JSON
	var paramData map[string]interface{}
	if err := json.Unmarshal([]byte(params), &paramData); err != nil {
		return nil, fmt.Errorf("failed to parse params JSON: %w", err)
	}

	// Create request body
	body, err := json.Marshal(paramData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("operation failed with status: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}
