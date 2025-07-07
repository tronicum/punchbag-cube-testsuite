package models

import (
	"time"
)

// CloudProvider represents the cloud provider type
type CloudProvider string

const (
	Azure   CloudProvider = "azure"
	AWS     CloudProvider = "aws"
	GCP     CloudProvider = "gcp"
	Hetzner CloudProvider = "hetzner"
	IONOS   CloudProvider = "ionos"
	StackIT CloudProvider = "stackit"
)

// Additional constants for compatibility
const (
	CloudProviderAzure   = "azure"
	CloudProviderAWS     = "aws"
	CloudProviderGCP     = "gcp"
	CloudProviderStackIT = "stackit"
	CloudProviderHetzner = "hetzner"
	CloudProviderIONOS   = "ionos"
)

// ClusterStatus represents the status of a cluster
type ClusterStatus string

const (
	ClusterStatusCreating ClusterStatus = "creating"
	ClusterStatusRunning  ClusterStatus = "running"
	ClusterStatusStopping ClusterStatus = "stopping"
	ClusterStatusStopped  ClusterStatus = "stopped"
	ClusterStatusDeleting ClusterStatus = "deleting"
	ClusterStatusFailed   ClusterStatus = "failed"
)

// TestStatus represents the status of a test
type TestStatus string

const (
	TestStatusPending TestStatus = "pending"
	TestStatusRunning TestStatus = "running"
	TestStatusPassed  TestStatus = "passed"
	TestStatusFailed  TestStatus = "failed"
)

// Cluster represents a Kubernetes cluster across different cloud providers
type Cluster struct {
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	Provider       CloudProvider `json:"provider"`
	Status         ClusterStatus `json:"status"`
	Config         map[string]interface{} `json:"config,omitempty"`
	ProviderConfig map[string]interface{} `json:"provider_config,omitempty"`
	ProjectID      string        `json:"project_id,omitempty"`
	ResourceGroup  string        `json:"resource_group,omitempty"`
	Location       string        `json:"location,omitempty"`
	Region         string        `json:"region,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

// TestResult represents the result of a cluster test
type TestResult struct {
	ID          string                 `json:"id"`
	ClusterID   string                 `json:"cluster_id"`
	TestType    string                 `json:"test_type"`
	Status      TestStatus             `json:"status"`
	Duration    time.Duration          `json:"duration"`
	Details     map[string]interface{} `json:"details,omitempty"`
	ErrorMsg    string                 `json:"error_message,omitempty"`
	StartedAt   time.Time              `json:"started_at"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
}

// TestRequest represents a request to run a test on a cluster
type TestRequest struct {
	ClusterID string                 `json:"cluster_id" binding:"required"`
	TestType  string                 `json:"test_type" binding:"required"`
	Config    map[string]interface{} `json:"config,omitempty"`
}

// NodePool represents a node pool within a cluster
type NodePool struct {
	ID              string    `json:"id"`
	ClusterID       string    `json:"cluster_id"`
	Name            string    `json:"name"`
	NodeCount       int       `json:"node_count"`
	MinNodes        int       `json:"min_nodes"`
	MaxNodes        int       `json:"max_nodes"`
	AutoScaling     bool      `json:"auto_scaling"`
	InstanceType    string    `json:"instance_type"`
	OSType          string    `json:"os_type"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
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

// Azure-specific models

type AzureMonitoring struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Config      map[string]interface{} `json:"config"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type AzureKubernetes struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	ClusterSize int                    `json:"cluster_size"`
	Config      map[string]interface{} `json:"config"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// LogAnalyticsWorkspace, AppInsightsResource, and AzureBudget are defined in azure.go
// Remove their duplicate definitions here and use the types from azure.go
// import (
//     . "github.com/tronicum/punchbag-cube-testsuite/shared/models"
// )

type LoadTestMetrics struct {
	TotalRequests     int64         `json:"total_requests"`
	SuccessfulRequests int64        `json:"successful_requests"`
	FailedRequests    int64         `json:"failed_requests"`
	AverageLatency    time.Duration `json:"average_latency"`
	P95Latency        time.Duration `json:"p95_latency"`
	P99Latency        time.Duration `json:"p99_latency"`
	MinLatency        time.Duration `json:"min_latency"`
	MaxLatency        time.Duration `json:"max_latency"`
	RequestsPerSecond float64       `json:"requests_per_second"`
	ErrorRate         float64       `json:"error_rate"`
}

type ClusterCreateRequest struct {
	Name           string                 `json:"name" binding:"required"`
	Provider       CloudProvider          `json:"provider" binding:"required"`
	Config         map[string]interface{} `json:"config,omitempty"`
	ProviderConfig map[string]interface{} `json:"provider_config,omitempty"`
	ProjectID      string                 `json:"project_id,omitempty"`
	ResourceGroup  string                 `json:"resource_group,omitempty"`
	Location       string                 `json:"location,omitempty"`
	Region         string                 `json:"region,omitempty"`
}

// ObjectStorageBucket represents a generic object storage bucket for any provider
// Provider-specific fields can be added via the ProviderConfig map.
type ObjectStorageBucket struct {
	ID            string                 `json:"id,omitempty"`
	Name          string                 `json:"name" binding:"required"`
	Provider      CloudProvider          `json:"provider" binding:"required"`
	Region        string                 `json:"region,omitempty"`
	Location      string                 `json:"location,omitempty"`
	CreatedAt     time.Time              `json:"created_at,omitempty"`
	UpdatedAt     time.Time              `json:"updated_at,omitempty"`
	Policy        *ObjectStoragePolicy   `json:"policy,omitempty"`
	Lifecycle     []ObjectStorageRule    `json:"lifecycle,omitempty"`
	ProviderConfig map[string]interface{} `json:"provider_config,omitempty"`
}

// ObjectStoragePolicy represents a generic bucket policy
// This can be extended for provider-specific policy fields.
type ObjectStoragePolicy struct {
	Version   string                   `json:"version"`
	Statement []ObjectStorageStatement `json:"statement"`
}

type ObjectStorageStatement struct {
	Effect    string                 `json:"effect"`
	Action    []string               `json:"action"`
	Resource  []string               `json:"resource"`
	Principal map[string]interface{} `json:"principal,omitempty"`
	Condition map[string]interface{} `json:"condition,omitempty"`
}

// ObjectStorageRule represents a lifecycle rule for a bucket
type ObjectStorageRule struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	// Add more fields as needed for expiration, transitions, etc.
}
