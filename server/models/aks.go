package models

import (
	"time"
)

// CloudProvider represents the different cloud providers
type CloudProvider string

const (
	CloudProviderAzure      CloudProvider = "azure"
	CloudProviderStackIT    CloudProvider = "schwarz-stackit"
	CloudProviderHetzner    CloudProvider = "hetzner-hcloud"
	CloudProviderIONOS      CloudProvider = "united-ionos"
	CloudProviderAWS        CloudProvider = "aws"
	CloudProviderGCP        CloudProvider = "gcp"
)

// Cluster represents a Kubernetes cluster across different cloud providers
type Cluster struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Provider          CloudProvider     `json:"provider"`
	ResourceGroup     string            `json:"resource_group,omitempty"`     // Azure specific
	Location          string            `json:"location"`
	Region            string            `json:"region,omitempty"`             // Alternative to location
	ProjectID         string            `json:"project_id,omitempty"`         // StackIT/GCP specific
	KubernetesVersion string            `json:"kubernetes_version"`
	Status            string            `json:"status"`
	NodeCount         int               `json:"node_count"`
	Tags              map[string]string `json:"tags,omitempty"`
	ProviderConfig    map[string]interface{} `json:"provider_config,omitempty"` // Provider-specific config
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

// Legacy type for backward compatibility
type AKSCluster = Cluster

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

// Legacy type for backward compatibility
type AKSTestResult = TestResult

// TestRequest represents a request to run a test on a cluster
type TestRequest struct {
	ClusterID string                 `json:"cluster_id" binding:"required"`
	TestType  string                 `json:"test_type" binding:"required"`
	Config    map[string]interface{} `json:"config,omitempty"`
}

// Legacy type for backward compatibility
type AKSTestRequest = TestRequest

// NodePool represents a node pool within a cluster
type NodePool struct {
	ID              string    `json:"id"`
	ClusterID       string    `json:"cluster_id"`
	Name            string    `json:"name"`
	VMSize          string    `json:"vm_size,omitempty"`          // Azure specific
	InstanceType    string    `json:"instance_type,omitempty"`    // AWS/StackIT specific
	MachineType     string    `json:"machine_type,omitempty"`     // GCP specific
	NodeCount       int       `json:"node_count"`
	MinNodes        int       `json:"min_nodes"`
	MaxNodes        int       `json:"max_nodes"`
	AutoScaling     bool      `json:"auto_scaling"`
	OSType          string    `json:"os_type"`
	KubernetesVersion string  `json:"kubernetes_version"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Legacy type for backward compatibility
type AKSNodePool = NodePool

// StackITClusterConfig represents StackIT-specific cluster configuration
type StackITClusterConfig struct {
	ProjectID           string `json:"project_id"`
	MaintenanceTimeStart string `json:"maintenance_time_start,omitempty"`
	MaintenanceTimeEnd   string `json:"maintenance_time_end,omitempty"`
	MaintenanceTimeZone  string `json:"maintenance_time_zone,omitempty"`
	HibernationSchedules []struct {
		Start    string `json:"start"`
		End      string `json:"end"`
		Timezone string `json:"timezone"`
	} `json:"hibernation_schedules,omitempty"`
}

// AzureClusterConfig represents Azure-specific cluster configuration  
type AzureClusterConfig struct {
	ResourceGroup     string `json:"resource_group"`
	SubscriptionID    string `json:"subscription_id,omitempty"`
	NetworkProfile    string `json:"network_profile,omitempty"`
	ServicePrincipal  string `json:"service_principal,omitempty"`
}

// AWSClusterConfig represents AWS-specific cluster configuration
type AWSClusterConfig struct {
	SubnetIDs         []string `json:"subnet_ids,omitempty"`
	SecurityGroupIDs  []string `json:"security_group_ids,omitempty"`
	RoleARN           string   `json:"role_arn,omitempty"`
	VpcID             string   `json:"vpc_id,omitempty"`
}

// GCPClusterConfig represents GCP-specific cluster configuration
type GCPClusterConfig struct {
	ProjectID         string `json:"project_id"`
	Network           string `json:"network,omitempty"`
	Subnetwork        string `json:"subnetwork,omitempty"`
	ServiceAccount    string `json:"service_account,omitempty"`
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

// HetznerClusterConfig represents Hetzner Cloud-specific cluster configuration
type HetznerClusterConfig struct {
	Token                string   `json:"token,omitempty"`
	NetworkID            int      `json:"network_id,omitempty"`
	NetworkZone          string   `json:"network_zone,omitempty"`
	ServerType           string   `json:"server_type,omitempty"`
	Location             string   `json:"location,omitempty"`
	SSHKeys              []string `json:"ssh_keys,omitempty"`
	FirewallIDs          []int    `json:"firewall_ids,omitempty"`
	LoadBalancerType     string   `json:"load_balancer_type,omitempty"`
	EnablePublicNetwork  bool     `json:"enable_public_network,omitempty"`
	EnablePrivateNetwork bool     `json:"enable_private_network,omitempty"`
}

// IONOSClusterConfig represents IONOS Cloud-specific cluster configuration
type IONOSClusterConfig struct {
	DatacenterID         string `json:"datacenter_id"`
	Username             string `json:"username,omitempty"`
	Password             string `json:"password,omitempty"`
	Token                string `json:"token,omitempty"`
	Endpoint             string `json:"endpoint,omitempty"`
	K8sClusterName       string `json:"k8s_cluster_name,omitempty"`
	MaintenanceWindow    struct {
		DayOfTheWeek string `json:"day_of_the_week,omitempty"`
		Time         string `json:"time,omitempty"`
	} `json:"maintenance_window,omitempty"`
	AllowReplace         bool     `json:"allow_replace,omitempty"`
	Public               bool     `json:"public,omitempty"`
	GatewayIP            string   `json:"gateway_ip,omitempty"`
	AvailableUpgradeVersions []string `json:"available_upgrade_versions,omitempty"`
}
