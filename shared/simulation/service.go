
package simulation

import (
	   "fmt"
	   "math/rand"
	   "os"
	   "time"
	   "github.com/tronicum/punchbag-cube-testsuite/shared/models"
)

// BucketStore returns the bucket store for direct manipulation (e.g., for dummy/test buckets)
func (s *SimulationService) BucketStore() *BucketStore {
	   return s.buckets
}

// SimulationService provides cloud provider simulation capabilities
type SimulationService struct {
	   rand         *rand.Rand
	   buckets      *BucketStore
	   persistPath  string
	   fastSimulate bool
	   debug        bool
}

// NewSimulationService creates a new simulation service
func NewSimulationService() *SimulationService {
		   persistPath := os.Getenv("CUBE_SERVER_SIM_PERSIST")
		   if persistPath == "" {
				   persistPath = "testdata/cube_server_sim_buckets.json"
		   }
	   s := &SimulationService{
			   rand: rand.New(rand.NewSource(time.Now().UnixNano())),
			   persistPath: persistPath,
			   fastSimulate: os.Getenv("FAST_SIMULATE") == "1",
			   debug: os.Getenv("CUBE_SERVER_DEBUG") == "1",
	   }
	   s.buckets = NewBucketStore(persistPath)
	   return s
}

// NewSimulationServiceWithOptions allows explicit config
func NewSimulationServiceWithOptions(fastSimulate, debug bool) *SimulationService {
	   persistPath := os.Getenv("CUBE_SERVER_SIM_PERSIST")
	   if persistPath == "" {
			   persistPath = "/tmp/cube_server_sim_buckets.json"
	   }
	   s := &SimulationService{
			   rand: rand.New(rand.NewSource(time.Now().UnixNano())),
			   persistPath: persistPath,
			   fastSimulate: fastSimulate,
			   debug: debug,
	   }
	   s.buckets = NewBucketStore(persistPath)
	   return s
}


// ProviderValidationResult represents the result of provider validation
type ProviderValidationResult struct {
	Provider  string                 `json:"provider"`
	Status    string                 `json:"status"`
	Valid     bool                   `json:"valid"`
	Regions   []string               `json:"regions,omitempty"`
	Services  map[string]interface{} `json:"services,omitempty"`
	Timestamp string                 `json:"timestamp"`
	Error     string                 `json:"error,omitempty"`
}

// SimulationRequest represents a simulation request
type SimulationRequest struct {
	Provider   string                 `json:"provider"`
	Operation  string                 `json:"operation"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// SimulationResult represents the result of a simulation
type SimulationResult struct {
	Provider  string                 `json:"provider"`
	Operation string                 `json:"operation"`
	Success   bool                   `json:"success"`
	Result    map[string]interface{} `json:"result,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Timestamp string                 `json:"timestamp"`
	Duration  time.Duration          `json:"duration"`
}

// ValidateProvider simulates provider validation
func (s *SimulationService) ValidateProvider(provider string, credentials map[string]interface{}) *ProviderValidationResult {
	now := time.Now()

	result := &ProviderValidationResult{
		Provider:  provider,
		Timestamp: now.Format(time.RFC3339),
	}

	switch provider {
	case "azure":
		result.Status = "valid"
		result.Valid = true
		result.Regions = []string{"eastus", "westus2", "centralus", "westeurope", "northeurope"}
		result.Services = map[string]interface{}{
			"aks": map[string]interface{}{
				"available":           true,
				"kubernetes_versions": []string{"1.28.0", "1.27.3", "1.26.6"},
				"vm_sizes":            []string{"Standard_D2s_v3", "Standard_D4s_v3", "Standard_B2s"},
			},
			"monitoring": map[string]interface{}{
				"available":            true,
				"log_analytics":        true,
				"application_insights": true,
			},
		}
	case "aws":
		result.Status = "valid"
		result.Valid = true
		result.Regions = []string{"us-east-1", "us-west-2", "eu-west-1", "eu-central-1", "ap-southeast-1"}
		result.Services = map[string]interface{}{
			"eks": map[string]interface{}{
				"available":           true,
				"kubernetes_versions": []string{"1.28", "1.27", "1.26"},
				"instance_types":      []string{"t3.medium", "t3.large", "m5.large", "m5.xlarge"},
			},
			"cloudwatch": map[string]interface{}{
				"available": true,
				"logs":      true,
				"metrics":   true,
			},
		}
	case "gcp":
		result.Status = "valid"
		result.Valid = true
		result.Regions = []string{"us-central1", "us-west1", "europe-west1", "asia-southeast1"}
		result.Services = map[string]interface{}{
			"gke": map[string]interface{}{
				"available":           true,
				"kubernetes_versions": []string{"1.28.3-gke.1286000", "1.27.7-gke.1056000"},
				"machine_types":       []string{"e2-medium", "e2-standard-4", "n1-standard-2"},
			},
			"stackdriver": map[string]interface{}{
				"available":  true,
				"logging":    true,
				"monitoring": true,
			},
		}
	case "hetzner":
		result.Status = "valid"
		result.Valid = true
		result.Regions = []string{"nbg1", "fsn1", "hel1", "ash"}
		result.Services = map[string]interface{}{
			"kubernetes": map[string]interface{}{
				"available":           true,
				"kubernetes_versions": []string{"1.28.0", "1.27.3", "1.26.6"},
				"server_types":        []string{"cx11", "cx21", "cx31", "cx41"},
			},
		}
	case "ionos":
		result.Status = "valid"
		result.Valid = true
		result.Regions = []string{"de/fra", "de/txl", "us/las", "gb/lhr"}
		result.Services = map[string]interface{}{
			"kubernetes": map[string]interface{}{
				"available":           true,
				"kubernetes_versions": []string{"1.28.0", "1.27.3", "1.26.6"},
				"cpu_families":        []string{"AMD_OPTERON", "INTEL_XEON", "INTEL_SKYLAKE"},
			},
		}
	case "stackit":
		result.Status = "valid"
		result.Valid = true
		result.Regions = []string{"eu-central-1", "eu-west-1"}
		result.Services = map[string]interface{}{
			"ske": map[string]interface{}{
				"available":           true,
				"kubernetes_versions": []string{"1.28.0", "1.27.3", "1.26.6"},
				"machine_types":       []string{"c1.2", "c1.3", "c1.4", "c1.5"},
			},
		}
	default:
		result.Status = "invalid"
		result.Valid = false
		result.Error = fmt.Sprintf("unsupported provider: %s", provider)
	}

	return result
}

// SimulateOperation simulates a cloud provider operation

func (s *SimulationService) SimulateOperation(req *SimulationRequest) *SimulationResult {
	   start := time.Now()

	   if s.debug {
			   fmt.Printf("[SIM DEBUG] SimulateOperation: provider=%s, op=%s, params=%#v\n", req.Provider, req.Operation, req.Parameters)
	   }

	   result := &SimulationResult{
			   Provider:  req.Provider,
			   Operation: req.Operation,
			   Timestamp: start.Format(time.RFC3339),
	   }

	   // Simulate operation delay unless fastSimulate is enabled
	   if !s.fastSimulate {
			   delay := time.Duration(s.rand.Intn(3000)+500) * time.Millisecond
			   time.Sleep(delay)
	   }

	   switch req.Operation {
	   case "create_bucket":
			   nameVal := s.getParamOrDefault(req.Parameters, "name", "")
			   name, _ := nameVal.(string)
			   if name == "" {
				   // Generate a default name for testing if none provided
				   name = "sim-bucket-" + s.generateRandomID()
			   }
			   // Use exact name like real S3 API - no modifications
			   regionVal := s.getParamOrDefault(req.Parameters, "region", "us-west-2")
			   region, _ := regionVal.(string)
			   bucket := s.buckets.Create(req.Provider, name, region)
			   result.Success = true
			   result.Result = bucket
	   case "delete_bucket":
			   bucketName, _ := req.Parameters["bucket"].(string)
			   result.Success, result.Result = s.buckets.Delete(req.Provider, bucketName)
	   case "list_buckets":
			   buckets := s.buckets.List(req.Provider)
			   result.Success = true
			   result.Result = map[string]interface{}{ "buckets": buckets, "total": len(buckets) }
	// ...existing code for clusters and tests...
	case "create_cluster":
		result.Success = true
		result.Result = s.simulateCreateCluster(req.Provider, req.Parameters)
	case "delete_cluster":
		result.Success = true
		result.Result = map[string]interface{}{
			"cluster_id": req.Parameters["cluster_id"],
			"status":     "deleting",
			"message":    "Cluster deletion initiated",
		}
	case "list_clusters":
		result.Success = true
		result.Result = s.simulateListClusters(req.Provider)
	case "get_cluster":
		result.Success = true
		result.Result = s.simulateGetCluster(req.Provider, req.Parameters)
	case "run_test":
		result.Success = true
		result.Result = s.simulateRunTest(req.Parameters)
	case "set_bucket_policy":
		result.Success = true
		result.Result = map[string]interface{}{
			"bucket": req.Parameters["bucket"],
			"policy": req.Parameters["policy"],
			"status": "policy_set",
		}
	case "set_bucket_versioning":
		result.Success = true
		result.Result = map[string]interface{}{
			"bucket":     req.Parameters["bucket"],
			"versioning": req.Parameters["enabled"],
			"status":     "versioning_set",
		}
	case "set_bucket_lifecycle":
		result.Success = true
		result.Result = map[string]interface{}{
			"bucket":    req.Parameters["bucket"],
			"lifecycle": req.Parameters["lifecycle"],
			"status":    "lifecycle_set",
		}
	default:
		result.Success = false
		result.Error = fmt.Sprintf("unsupported operation: %s", req.Operation)
	}
	result.Duration = time.Since(start)
	return result
}


// simulateCreateBucket simulates S3/object storage bucket creation
func (s *SimulationService) simulateCreateBucket(provider string, params map[string]interface{}) map[string]interface{} {
	nameVal := s.getParamOrDefault(params, "name", "sim-bucket-")
	name, _ := nameVal.(string)
	bucketName := name + s.generateRandomID()
	regionVal := s.getParamOrDefault(params, "region", "us-west-2")
	region, _ := regionVal.(string)
	return map[string]interface{}{
		"bucket":   bucketName,
		"provider": provider,
		"region":   region,
		"status":   "created",
	}
}

// simulateCreateCluster simulates cluster creation
func (s *SimulationService) simulateCreateCluster(provider string, params map[string]interface{}) map[string]interface{} {
	clusterID := fmt.Sprintf("sim-%s-%d", provider, s.rand.Intn(10000))

	result := map[string]interface{}{
		"cluster_id": clusterID,
		"name":       params["name"],
		"provider":   provider,
		"status":     "creating",
		"node_count": s.getParamOrDefault(params, "node_count", 3),
		"created_at": time.Now().Format(time.RFC3339),
	}

	// Add provider-specific fields
	switch provider {
	case "azure":
		result["resource_group"] = s.getParamOrDefault(params, "resource_group", "rg-"+clusterID)
		result["location"] = s.getParamOrDefault(params, "location", "eastus")
		result["vm_size"] = s.getParamOrDefault(params, "vm_size", "Standard_D2s_v3")
	case "aws":
		result["region"] = s.getParamOrDefault(params, "region", "us-west-2")
		result["instance_type"] = s.getParamOrDefault(params, "instance_type", "t3.medium")
		result["vpc_id"] = "vpc-" + s.generateRandomID()
	case "gcp":
		result["project_id"] = s.getParamOrDefault(params, "project_id", "project-"+s.generateRandomID())
		result["region"] = s.getParamOrDefault(params, "region", "us-central1")
		result["machine_type"] = s.getParamOrDefault(params, "machine_type", "e2-medium")
	}

	return result
}

// simulateListClusters simulates listing clusters
func (s *SimulationService) simulateListClusters(provider string) map[string]interface{} {
	clusters := make([]map[string]interface{}, 0)

	// Generate 2-4 sample clusters
	numClusters := s.rand.Intn(3) + 2
	for i := 0; i < numClusters; i++ {
		clusterID := fmt.Sprintf("sim-%s-%d", provider, s.rand.Intn(10000))
		cluster := map[string]interface{}{
			"cluster_id": clusterID,
			"name":       fmt.Sprintf("%s-cluster-%d", provider, i+1),
			"provider":   provider,
			"status":     []string{"running", "creating", "stopped"}[s.rand.Intn(3)],
			"node_count": s.rand.Intn(5) + 1,
			"created_at": time.Now().Add(-time.Duration(s.rand.Intn(168)) * time.Hour).Format(time.RFC3339),
		}
		clusters = append(clusters, cluster)
	}

	return map[string]interface{}{
		"clusters": clusters,
		"total":    len(clusters),
	}
}

// simulateGetCluster simulates getting cluster details
func (s *SimulationService) simulateGetCluster(provider string, params map[string]interface{}) map[string]interface{} {
	clusterID := s.getParamOrDefault(params, "cluster_id", "sim-"+provider+"-1234").(string)

	return map[string]interface{}{
		"cluster_id":         clusterID,
		"name":               "sample-cluster",
		"provider":           provider,
		"status":             "running",
		"node_count":         3,
		"kubernetes_version": "1.28.0",
		"created_at":         time.Now().Add(-time.Duration(s.rand.Intn(168)) * time.Hour).Format(time.RFC3339),
		"endpoints": map[string]interface{}{
			"api_server": "https://api-" + clusterID + ".example.com",
			"dashboard":  "https://dashboard-" + clusterID + ".example.com",
		},
	}
}

// simulateRunTest simulates running a test
func (s *SimulationService) simulateRunTest(params map[string]interface{}) map[string]interface{} {
	testID := fmt.Sprintf("test-%d", s.rand.Intn(10000))
	testType := s.getParamOrDefault(params, "test_type", "connectivity").(string)

	// 90% success rate
	status := "passed"
	if s.rand.Float32() < 0.1 {
		status = "failed"
	}

	result := map[string]interface{}{
		"test_id":    testID,
		"cluster_id": params["cluster_id"],
		"test_type":  testType,
		"status":     status,
		"started_at": time.Now().Format(time.RFC3339),
		"duration":   fmt.Sprintf("%ds", s.rand.Intn(300)+30),
	}

	// Add test-specific results
	switch testType {
	case "connectivity":
		result["endpoints_tested"] = s.rand.Intn(10) + 5
		result["successful_connections"] = s.rand.Intn(15) + 10
		result["avg_response_time_ms"] = s.rand.Intn(100) + 20
	case "performance":
		result["cpu_usage_percent"] = s.rand.Intn(40) + 30
		result["memory_usage_percent"] = s.rand.Intn(50) + 25
		result["requests_per_second"] = s.rand.Intn(1000) + 500
	case "security":
		result["vulnerabilities_found"] = s.rand.Intn(3)
		result["security_score"] = s.rand.Intn(30) + 70
	case "compliance":
		result["policies_checked"] = s.rand.Intn(50) + 25
		result["compliance_score"] = s.rand.Intn(25) + 75
	}

	return result
}

// GenerateClusterFromSimulation converts simulation result to Cluster model
func (s *SimulationService) GenerateClusterFromSimulation(provider string, name string, config map[string]interface{}) *models.Cluster {
	now := time.Now()
	clusterID := fmt.Sprintf("sim-%s-%d", provider, s.rand.Intn(10000))

	cluster := &models.Cluster{
		ID:        clusterID,
		Name:      name,
		Provider:  models.CloudProvider(provider),
		Status:    models.ClusterStatusRunning,
		Config:    make(map[string]interface{}),
		CreatedAt: now.Add(-time.Duration(s.rand.Intn(24)) * time.Hour),
		UpdatedAt: now,
	}

	// Set default config
	cluster.Config["kubernetes_version"] = "1.28.0"
	cluster.Config["node_count"] = 3
	cluster.Config["auto_scaling"] = true

	// Merge provided config
	for k, v := range config {
		cluster.Config[k] = v
	}

	// Set provider-specific fields
	switch models.CloudProvider(provider) {
	case models.Azure:
		cluster.ResourceGroup = s.getParamOrDefault(config, "resource_group", "rg-"+name).(string)
		cluster.Location = s.getParamOrDefault(config, "location", "eastus").(string)
		cluster.ProviderConfig = map[string]interface{}{
			"sku":               "Standard_D2s_v3",
			"network_plugin":    "azure",
			"enable_rbac":       true,
			"enable_monitoring": true,
		}
	case models.AWS:
		cluster.Region = s.getParamOrDefault(config, "region", "us-west-2").(string)
		cluster.ProviderConfig = map[string]interface{}{
			"instance_type":    "t3.medium",
			"vpc_id":           "vpc-" + s.generateRandomID(),
			"subnet_ids":       []string{"subnet-" + s.generateRandomID(), "subnet-" + s.generateRandomID()},
			"endpoint_private": false,
		}
	case models.GCP:
		cluster.ProjectID = s.getParamOrDefault(config, "project_id", "project-"+s.generateRandomID()).(string)
		cluster.Region = s.getParamOrDefault(config, "region", "us-central1").(string)
		cluster.ProviderConfig = map[string]interface{}{
			"machine_type":     "e2-medium",
			"disk_size_gb":     100,
			"network":          "default",
			"enable_autopilot": false,
		}
	}

	return cluster
}

// GenerateTestResultFromSimulation converts simulation result to TestResult model
func (s *SimulationService) GenerateTestResultFromSimulation(clusterID, testType string) *models.TestResult {
	now := time.Now()
	testID := fmt.Sprintf("test-%d", s.rand.Intn(10000))

	// 90% success rate
	status := models.TestStatusPassed
	var errorMsg string
	if s.rand.Float32() < 0.1 {
		status = models.TestStatusFailed
		errorMsg = "Simulated test failure"
	}

	duration := time.Duration(s.rand.Intn(300)+30) * time.Second
	details := make(map[string]interface{})

	// Add test-specific details
	switch testType {
	case "connectivity":
		details["endpoints_tested"] = s.rand.Intn(10) + 5
		details["successful_connections"] = s.rand.Intn(15) + 10
		details["avg_response_time_ms"] = s.rand.Intn(100) + 20
	case "performance":
		details["cpu_usage_percent"] = s.rand.Intn(40) + 30
		details["memory_usage_percent"] = s.rand.Intn(50) + 25
		details["requests_per_second"] = s.rand.Intn(1000) + 500
		details["p95_latency_ms"] = s.rand.Intn(200) + 50
	case "security":
		details["vulnerabilities_found"] = s.rand.Intn(3)
		details["security_score"] = s.rand.Intn(30) + 70
		details["compliant_policies"] = s.rand.Intn(20) + 15
	case "compliance":
		details["policies_checked"] = s.rand.Intn(50) + 25
		details["compliant_policies"] = s.rand.Intn(45) + 20
		details["compliance_score"] = s.rand.Intn(25) + 75
	}

	completedAt := now
	return &models.TestResult{
		ID:          testID,
		ClusterID:   clusterID,
		TestType:    testType,
		Status:      status,
		Duration:    duration,
		Details:     details,
		ErrorMsg:    errorMsg,
		StartedAt:   now.Add(-duration),
		CompletedAt: &completedAt,
	}
}

// Helper functions

func (s *SimulationService) getParamOrDefault(params map[string]interface{}, key string, defaultValue interface{}) interface{} {
	if val, exists := params[key]; exists {
		return val
	}
	return defaultValue
}

func (s *SimulationService) generateRandomID() string {
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, 8)
	for i := range result {
		result[i] = chars[s.rand.Intn(len(chars))]
	}
	return string(result)
}
