package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents the API client
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new API client
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Cluster represents a Kubernetes cluster across different cloud providers
type Cluster struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	CloudProvider string                 `json:"cloud_provider"` // azure, schwarz-stackit, aws, gcp
	Status        string                 `json:"status"`
	Config        map[string]interface{} `json:"config,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// AKSCluster represents an Azure Kubernetes Service cluster (for backward compatibility)
type AKSCluster struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	ResourceGroup     string            `json:"resource_group"`
	Location          string            `json:"location"`
	KubernetesVersion string            `json:"kubernetes_version"`
	Status            string            `json:"status"`
	NodeCount         int               `json:"node_count"`
	Tags              map[string]string `json:"tags,omitempty"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

// TestResult represents the result of a cluster test
type TestResult struct {
	ID          string                 `json:"id"`
	ClusterID   string                 `json:"cluster_id"`
	TestType    string                 `json:"test_type"`
	Status      string                 `json:"status"`
	Duration    time.Duration          `json:"duration"`
	Details     map[string]interface{} `json:"details,omitempty"`
	ErrorMsg    string                 `json:"error_message,omitempty"`
	StartedAt   time.Time              `json:"started_at"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
}

// AKSTestResult represents the result of an AKS test (for backward compatibility)
type AKSTestResult struct {
	ID          string                 `json:"id"`
	ClusterID   string                 `json:"cluster_id"`
	TestType    string                 `json:"test_type"`
	Status      string                 `json:"status"`
	Duration    time.Duration          `json:"duration"`
	Details     map[string]interface{} `json:"details,omitempty"`
	ErrorMsg    string                 `json:"error_message,omitempty"`
	StartedAt   time.Time              `json:"started_at"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
}

// TestRequest represents a request to run a test on a cluster
type TestRequest struct {
	ClusterID string                 `json:"cluster_id"`
	TestType  string                 `json:"test_type"`
	Config    map[string]interface{} `json:"config,omitempty"`
}

// AKSTestRequest represents a request to run a test on an AKS cluster (for backward compatibility)
type AKSTestRequest struct {
	ClusterID string                 `json:"cluster_id"`
	TestType  string                 `json:"test_type"`
	Config    map[string]interface{} `json:"config,omitempty"`
}

// ClustersResponse represents the response for listing clusters
type ClustersResponse struct {
	Clusters []*Cluster `json:"clusters"`
}

// TestResultsResponse represents the response for listing test results
type TestResultsResponse struct {
	TestResults []*TestResult `json:"test_results"`
}

// doRequest performs an HTTP request
func (c *Client) doRequest(method, path string, body interface{}) (*http.Response, error) {
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
func (c *Client) ListClusters() ([]*Cluster, error) {
	resp, err := c.doRequest("GET", "/api/v1/clusters", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	var response ClustersResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return response.Clusters, nil
}

// ListClustersByProvider lists clusters filtered by provider
func (c *Client) ListClustersByProvider(provider string) ([]*Cluster, error) {
	resp, err := c.doRequest("GET", "/api/v1/clusters?provider="+provider, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	var response ClustersResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return response.Clusters, nil
}

// ListAKSClusters lists all AKS clusters (backward compatibility)
func (c *Client) ListAKSClusters() ([]*AKSCluster, error) {
	clusters, err := c.ListClustersByProvider("azure")
	if err != nil {
		return nil, err
	}

	// Convert to AKSCluster format for backward compatibility
	aksClusters := make([]*AKSCluster, len(clusters))
	for i, cluster := range clusters {
		aksClusters[i] = c.convertToAKSCluster(cluster)
	}

	return aksClusters, nil
}

// Helper method to convert Cluster to AKSCluster
func (c *Client) convertToAKSCluster(cluster *Cluster) *AKSCluster {
	aksCluster := &AKSCluster{
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
func (c *Client) GetCluster(id string) (*Cluster, error) {
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

	var cluster Cluster
	if err := json.NewDecoder(resp.Body).Decode(&cluster); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &cluster, nil
}

// GetAKSCluster gets an AKS cluster by ID (backward compatibility)
func (c *Client) GetAKSCluster(id string) (*AKSCluster, error) {
	cluster, err := c.GetCluster(id)
	if err != nil {
		return nil, err
	}

	// Convert to AKSCluster format for backward compatibility
	aksCluster := &AKSCluster{
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
func (c *Client) CreateMultiCloudCluster(cluster *Cluster) (*Cluster, error) {
	resp, err := c.doRequest("POST", "/api/v1/clusters", cluster)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	var created Cluster
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &created, nil
}

// CreateStackITCluster creates a new StackIT cluster
func (c *Client) CreateStackITCluster(name, projectID, region string, config map[string]interface{}) (*Cluster, error) {
	cluster := &Cluster{
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
func (c *Client) CreateAzureCluster(name, resourceGroup, location string, config map[string]interface{}) (*Cluster, error) {
	cluster := &Cluster{
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
func (c *Client) CreateHetznerCluster(name, location string, config map[string]interface{}) (*Cluster, error) {
	cluster := &Cluster{
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
		}
	}
	
	return c.CreateMultiCloudCluster(cluster)
}

// CreateIONOSCluster creates a new IONOS Cloud cluster
func (c *Client) CreateIONOSCluster(name, datacenterID string, config map[string]interface{}) (*Cluster, error) {
	cluster := &Cluster{
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
func (c *Client) UpdateCluster(id string, cluster *AKSCluster) (*AKSCluster, error) {
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

	var updated AKSCluster
	if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &updated, nil
}

// DeleteCluster deletes a cluster
func (c *Client) DeleteCluster(id string) error {
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
func (c *Client) RunTest(clusterID string, testReq *TestRequest) (*TestResult, error) {
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

	var result TestResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}

// RunAKSTest runs a test on an AKS cluster (backward compatibility)
func (c *Client) RunAKSTest(clusterID string, testReq *AKSTestRequest) (*AKSTestResult, error) {
	multiTestReq := &TestRequest{
		ClusterID: testReq.ClusterID,
		TestType:  testReq.TestType,
		Config:    testReq.Config,
	}
	
	result, err := c.RunTest(clusterID, multiTestReq)
	if err != nil {
		return nil, err
	}
	
	// Convert to AKSTestResult for backward compatibility
	aksResult := &AKSTestResult{
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
func (c *Client) GetTestResult(id string) (*TestResult, error) {
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

	var result TestResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}

// GetAKSTestResult gets an AKS test result by ID (backward compatibility)
func (c *Client) GetAKSTestResult(id string) (*AKSTestResult, error) {
	result, err := c.GetTestResult(id)
	if err != nil {
		return nil, err
	}
	
	// Convert to AKSTestResult for backward compatibility
	aksResult := &AKSTestResult{
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
func (c *Client) ListTestResults(clusterID string) ([]*TestResult, error) {
	resp, err := c.doRequest("GET", "/api/v1/clusters/"+clusterID+"/tests", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	var response TestResultsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return response.TestResults, nil
}

// ListAKSTestResults lists test results for an AKS cluster (backward compatibility)
func (c *Client) ListAKSTestResults(clusterID string) ([]*AKSTestResult, error) {
	results, err := c.ListTestResults(clusterID)
	if err != nil {
		return nil, err
	}
	
	// Convert to AKSTestResult for backward compatibility
	aksResults := make([]*AKSTestResult, len(results))
	for i, result := range results {
		aksResults[i] = &AKSTestResult{
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
