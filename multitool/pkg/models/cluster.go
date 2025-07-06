package models

import (
	"time"
	sharedmodels "punchbag-cube-testsuite/shared/models"
)

// Cluster represents a Kubernetes cluster across different cloud providers
// Use sharedmodels.CloudProvider, sharedmodels.ClusterStatus, etc.
type Cluster struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Provider       sharedmodels.CloudProvider          `json:"provider"`
	Status         sharedmodels.ClusterStatus          `json:"status"`
	Config         map[string]interface{} `json:"config,omitempty"`
	ProviderConfig map[string]interface{} `json:"provider_config,omitempty"`
	ProjectID      string                 `json:"project_id,omitempty"`
	ResourceGroup  string                 `json:"resource_group,omitempty"`
	Location       string                 `json:"location,omitempty"`
	Region         string                 `json:"region,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// TestResult represents the result of a cluster test
type TestResult struct {
	ID          string                 `json:"id"`
	ClusterID   string                 `json:"cluster_id"`
	TestType    string                 `json:"test_type"`
	Status      sharedmodels.TestStatus             `json:"status"`
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
	ID           string    `json:"id"`
	ClusterID    string    `json:"cluster_id"`
	Name         string    `json:"name"`
	NodeCount    int       `json:"node_count"`
	MinNodes     int       `json:"min_nodes"`
	MaxNodes     int       `json:"max_nodes"`
	AutoScaling  bool      `json:"auto_scaling"`
	InstanceType string    `json:"instance_type"`
	OSType       string    `json:"os_type"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ClusterCreateRequest represents a request to create a cluster
type ClusterCreateRequest struct {
	Name           string                 `json:"name" binding:"required"`
	Provider       sharedmodels.CloudProvider          `json:"provider" binding:"required"`
	Config         map[string]interface{} `json:"config,omitempty"`
	ProviderConfig map[string]interface{} `json:"provider_config,omitempty"`
	ProjectID      string                 `json:"project_id,omitempty"`
	ResourceGroup  string                 `json:"resource_group,omitempty"`
	Location       string                 `json:"location,omitempty"`
	Region         string                 `json:"region,omitempty"`
}

// LoadTestMetrics represents metrics collected during load testing
type LoadTestMetrics struct {
	TotalRequests      int64         `json:"total_requests"`
	SuccessfulRequests int64         `json:"successful_requests"`
	FailedRequests     int64         `json:"failed_requests"`
	AverageLatency     time.Duration `json:"average_latency"`
	P95Latency         time.Duration `json:"p95_latency"`
	P99Latency         time.Duration `json:"p99_latency"`
	MinLatency         time.Duration `json:"min_latency"`
	MaxLatency         time.Duration `json:"max_latency"`
	RequestsPerSecond  float64       `json:"requests_per_second"`
	ErrorRate          float64       `json:"error_rate"`
}
