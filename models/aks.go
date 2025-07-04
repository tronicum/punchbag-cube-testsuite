package models

import (
	"time"
)

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

// TestRequest represents a request to run a test on a cluster
type TestRequest struct {
	ClusterID string                 `json:"cluster_id" binding:"required"`
	TestType  string                 `json:"test_type" binding:"required"`
	Config    map[string]interface{} `json:"config,omitempty"`
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

// AKSTestRequest represents a request to run a test on an AKS cluster (for backward compatibility)
type AKSTestRequest struct {
	ClusterID string                 `json:"cluster_id" binding:"required"`
	TestType  string                 `json:"test_type" binding:"required"`
	Config    map[string]interface{} `json:"config,omitempty"`
}

// AKSNodePool represents a node pool within an AKS cluster
type AKSNodePool struct {
	ID              string    `json:"id"`
	ClusterID       string    `json:"cluster_id"`
	Name            string    `json:"name"`
	VMSize          string    `json:"vm_size"`
	NodeCount       int       `json:"node_count"`
	MinNodes        int       `json:"min_nodes"`
	MaxNodes        int       `json:"max_nodes"`
	AutoScaling     bool      `json:"auto_scaling"`
	OSType          string    `json:"os_type"`
	KubernetesVersion string  `json:"kubernetes_version"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// PunchbagTestConfig represents configuration for punchbag testing
type PunchbagTestConfig struct {
	Duration     time.Duration `json:"duration"`
	Concurrency  int           `json:"concurrency"`
	RequestRate  int           `json:"request_rate"`
	TargetURL    string        `json:"target_url"`
	Headers      map[string]string `json:"headers,omitempty"`
	Body         string        `json:"body,omitempty"`
	Method       string        `json:"method"`
	ExpectedCode int           `json:"expected_code"`
}

// HetznerClusterConfig represents configuration for Hetzner Cloud clusters
type HetznerClusterConfig struct {
	ServerType      string            `json:"server_type"`
	Image           string            `json:"image"`
	Location        string            `json:"location"`
	Network         string            `json:"network,omitempty"`
	SSHKeys         []string          `json:"ssh_keys,omitempty"`
	Firewalls       []string          `json:"firewalls,omitempty"`
	Labels          map[string]string `json:"labels,omitempty"`
	UserData        string            `json:"user_data,omitempty"`
	Backups         bool              `json:"backups"`
	KubernetesVersion string          `json:"kubernetes_version"`
	NodeCount       int               `json:"node_count"`
	EnableIPv6      bool              `json:"enable_ipv6"`
	PrivateNetworkOnly bool           `json:"private_network_only"`
}

// IONOSClusterConfig represents configuration for IONOS Cloud clusters  
type IONOSClusterConfig struct {
	DatacenterID    string            `json:"datacenter_id"`
	CPUFamily       string            `json:"cpu_family"`
	Cores           int               `json:"cores"`
	RAM             int               `json:"ram"`
	StorageType     string            `json:"storage_type"`
	StorageSize     int               `json:"storage_size"`
	Location        string            `json:"location"`
	ImageAlias      string            `json:"image_alias"`
	SSHKeys         []string          `json:"ssh_keys,omitempty"`
	Properties      map[string]string `json:"properties,omitempty"`
	KubernetesVersion string          `json:"kubernetes_version"`
	NodeCount       int               `json:"node_count"`
	PublicLAN       bool              `json:"public_lan"`
	DHCPEnabled     bool              `json:"dhcp_enabled"`
	MaintenanceWindow string          `json:"maintenance_window,omitempty"`
}

// CloudProvider constants for multi-cloud support
const (
	CloudProviderAzure      = "azure"
	CloudProviderAWS        = "aws"
	CloudProviderGCP        = "gcp"
	CloudProviderStackit    = "schwarz-stackit"
	CloudProviderHetzner    = "hetzner-hcloud"
	CloudProviderIONOS      = "united-ionos"
)

// LoadTestMetrics represents metrics collected during load testing
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

// Structs for Azure services

// AzureMonitoring represents an Azure Monitoring resource.
type AzureMonitoring struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Config      map[string]interface{} `json:"config"`
}

// AzureKubernetes represents an Azure Kubernetes resource.
type AzureKubernetes struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	ClusterSize int                    `json:"cluster_size"`
	Config      map[string]interface{} `json:"config"`
}

// AzureBudget represents an Azure Budget resource.
type AzureBudget struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Amount      float64               `json:"amount"`
	Config      map[string]interface{} `json:"config"`
}
