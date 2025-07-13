package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"
)

// APIClient represents a client for interacting with the punchbag server API.
type APIClient struct {
	   baseURL    string
	   httpClient *http.Client
}

// NewAPIClient creates a new API client. Returns nil if baseURL is empty.
func NewAPIClient(baseURL string) *APIClient {
	   if strings.TrimSpace(baseURL) == "" {
			   return nil
	   }
	   return &APIClient{
			   baseURL: strings.TrimSuffix(baseURL, "/"),
			   httpClient: &http.Client{
					   Timeout: 30 * time.Second,
			   },
	   }
}

// buildURL helps construct endpoint URLs.
func (c *APIClient) buildURL(path string) string {
	   return c.baseURL + path
}

// Ping checks if the API server is reachable and healthy.
func (c *APIClient) Ping() error {
	   url := c.buildURL("/ping")
	   resp, err := c.httpClient.Get(url)
	   if err != nil {
			   return err
	   }
	   defer resp.Body.Close()
	   if resp.StatusCode != http.StatusOK {
			   return fmt.Errorf("API ping failed: status %d", resp.StatusCode)
	   }
	   return nil
}

// ClusterClient provides methods for cluster operations.
type ClusterClient struct {
	   client *APIClient
}

// NewClusterClient creates a new cluster client. Returns nil if client is nil.
func NewClusterClient(client *APIClient) *ClusterClient {
	   if client == nil {
			   return nil
	   }
	   return &ClusterClient{client: client}
}

// CreateCluster creates a new cluster.
func (c *ClusterClient) CreateCluster(req *sharedmodels.ClusterCreateRequest) (*sharedmodels.Cluster, error) {
	   url := c.client.buildURL("/api/clusters")
	   data, err := json.Marshal(req)
	   if err != nil {
			   return nil, fmt.Errorf("failed to marshal request: %w", err)
	   }
	   resp, err := c.client.httpClient.Post(url, "application/json", strings.NewReader(string(data)))
	   if err != nil {
			   return nil, fmt.Errorf("failed to create cluster: %w", err)
	   }
	   defer resp.Body.Close()
	   if resp.StatusCode != http.StatusCreated {
			   return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	   }
	   var cluster sharedmodels.Cluster
	   if err := json.NewDecoder(resp.Body).Decode(&cluster); err != nil {
			   return nil, fmt.Errorf("failed to decode response: %w", err)
	   }
	   return &cluster, nil
}

// GetCluster retrieves a cluster by ID.
func (c *ClusterClient) GetCluster(id string) (*sharedmodels.Cluster, error) {
	   url := c.client.buildURL("/api/clusters/" + id)
	   resp, err := c.client.httpClient.Get(url)
	   if err != nil {
			   return nil, fmt.Errorf("failed to get cluster: %w", err)
	   }
	   defer resp.Body.Close()
	   if resp.StatusCode == http.StatusNotFound {
			   return nil, fmt.Errorf("cluster not found")
	   }
	   if resp.StatusCode != http.StatusOK {
			   return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	   }
	   var cluster sharedmodels.Cluster
	   if err := json.NewDecoder(resp.Body).Decode(&cluster); err != nil {
			   return nil, fmt.Errorf("failed to decode response: %w", err)
	   }
	   return &cluster, nil
}

// ListClusters retrieves all clusters.
func (c *ClusterClient) ListClusters() ([]*sharedmodels.Cluster, error) {
	   url := c.client.buildURL("/api/clusters")
	   resp, err := c.client.httpClient.Get(url)
	   if err != nil {
			   return nil, fmt.Errorf("failed to list clusters: %w", err)
	   }
	   defer resp.Body.Close()
	   if resp.StatusCode != http.StatusOK {
			   return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	   }
	   var clusters []*sharedmodels.Cluster
	   if err := json.NewDecoder(resp.Body).Decode(&clusters); err != nil {
			   return nil, fmt.Errorf("failed to decode response: %w", err)
	   }
	   return clusters, nil
}

// ListClustersByProvider retrieves clusters filtered by provider.
func (c *ClusterClient) ListClustersByProvider(provider sharedmodels.CloudProvider) ([]*sharedmodels.Cluster, error) {
	   url := c.client.buildURL("/api/clusters?provider=" + string(provider))
	   resp, err := c.client.httpClient.Get(url)
	   if err != nil {
			   return nil, fmt.Errorf("failed to list clusters: %w", err)
	   }
	   defer resp.Body.Close()
	   if resp.StatusCode != http.StatusOK {
			   return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	   }
	   var clusters []*sharedmodels.Cluster
	   if err := json.NewDecoder(resp.Body).Decode(&clusters); err != nil {
			   return nil, fmt.Errorf("failed to decode response: %w", err)
	   }
	   return clusters, nil
}

// DeleteCluster deletes a cluster by ID.
func (c *ClusterClient) DeleteCluster(id string) error {
	   url := c.client.buildURL("/api/clusters/" + id)
	   req, err := http.NewRequest(http.MethodDelete, url, nil)
	   if err != nil {
			   return fmt.Errorf("failed to create request: %w", err)
	   }
	   resp, err := c.client.httpClient.Do(req)
	   if err != nil {
			   return fmt.Errorf("failed to delete cluster: %w", err)
	   }
	   defer resp.Body.Close()
	   if resp.StatusCode == http.StatusNotFound {
			   return fmt.Errorf("cluster not found")
	   }
	   if resp.StatusCode != http.StatusNoContent {
			   return fmt.Errorf("API returned status %d", resp.StatusCode)
	   }
	   return nil
}

// TestClient provides methods for test operations.
type TestClient struct {
	   client *APIClient
}

// NewTestClient creates a new test client. Returns nil if client is nil.
func NewTestClient(client *APIClient) *TestClient {
	   if client == nil {
			   return nil
	   }
	   return &TestClient{client: client}
}

// RunTest runs a test on a cluster
func (t *TestClient) RunTest(req *sharedmodels.TestRequest) (*sharedmodels.TestResult, error) {
	url := fmt.Sprintf("%s/api/tests", t.client.baseURL)

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := t.client.httpClient.Post(url, "application/json", strings.NewReader(string(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to run test: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result sharedmodels.TestResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetTestResult retrieves a test result by ID
func (t *TestClient) GetTestResult(id string) (*sharedmodels.TestResult, error) {
	url := fmt.Sprintf("%s/api/tests/%s", t.client.baseURL, id)

	resp, err := t.client.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get test result: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("test result not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result sharedmodels.TestResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// ListTestResults retrieves test results for a cluster
func (t *TestClient) ListTestResults(clusterID string) ([]*sharedmodels.TestResult, error) {
	url := fmt.Sprintf("%s/api/clusters/%s/tests", t.client.baseURL, clusterID)

	resp, err := t.client.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to list test results: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var results []*sharedmodels.TestResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return results, nil
}
