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
