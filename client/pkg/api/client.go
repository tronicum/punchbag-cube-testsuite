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

// AKSCluster represents an Azure Kubernetes Service cluster
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

// AKSTestResult represents the result of an AKS test
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

// AKSTestRequest represents a request to run a test on an AKS cluster
type AKSTestRequest struct {
	ClusterID string                 `json:"cluster_id"`
	TestType  string                 `json:"test_type"`
	Config    map[string]interface{} `json:"config,omitempty"`
}

// ClustersResponse represents the response for listing clusters
type ClustersResponse struct {
	Clusters []*AKSCluster `json:"clusters"`
}

// TestResultsResponse represents the response for listing test results
type TestResultsResponse struct {
	TestResults []*AKSTestResult `json:"test_results"`
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

// ListClusters lists all clusters
func (c *Client) ListClusters() ([]*AKSCluster, error) {
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

// GetCluster gets a cluster by ID
func (c *Client) GetCluster(id string) (*AKSCluster, error) {
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

	var cluster AKSCluster
	if err := json.NewDecoder(resp.Body).Decode(&cluster); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &cluster, nil
}

// CreateCluster creates a new cluster
func (c *Client) CreateCluster(cluster *AKSCluster) (*AKSCluster, error) {
	resp, err := c.doRequest("POST", "/api/v1/clusters", cluster)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	var created AKSCluster
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &created, nil
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

// RunTest runs a test on a cluster
func (c *Client) RunTest(clusterID string, testReq *AKSTestRequest) (*AKSTestResult, error) {
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

	var result AKSTestResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}

// GetTestResult gets a test result by ID
func (c *Client) GetTestResult(id string) (*AKSTestResult, error) {
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

	var result AKSTestResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}

// ListTestResults lists test results for a cluster
func (c *Client) ListTestResults(clusterID string) ([]*AKSTestResult, error) {
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
