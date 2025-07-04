package store

import (
	"fmt"
	"sync"
	"time"

	"punchbag-cube-testsuite/models"

	"github.com/google/uuid"
)

// Store defines the interface for cluster and test result storage
type Store interface {
	// Cluster operations
	CreateCluster(cluster *models.Cluster) (*models.Cluster, error)
	GetCluster(id string) (*models.Cluster, error)
	UpdateCluster(id string, cluster *models.Cluster) (*models.Cluster, error)
	DeleteCluster(id string) error
	ListClusters() ([]*models.Cluster, error)
	ListClustersByProvider(provider string) ([]*models.Cluster, error)

	// Test result operations
	CreateTestResult(result *models.TestResult) (*models.TestResult, error)
	GetTestResult(id string) (*models.TestResult, error)
	UpdateTestResult(id string, result *models.TestResult) (*models.TestResult, error)
	ListTestResults(clusterID string) ([]*models.TestResult, error)
	
	// Legacy AKS operations (for backward compatibility)
	CreateAKSCluster(cluster *models.AKSCluster) (*models.AKSCluster, error)
	GetAKSCluster(id string) (*models.AKSCluster, error)
	UpdateAKSCluster(id string, cluster *models.AKSCluster) (*models.AKSCluster, error)
	ListAKSClusters() ([]*models.AKSCluster, error)
	
	CreateAKSTestResult(result *models.AKSTestResult) (*models.AKSTestResult, error)
	GetAKSTestResult(id string) (*models.AKSTestResult, error)
	UpdateAKSTestResult(id string, result *models.AKSTestResult) (*models.AKSTestResult, error)
	ListAKSTestResults(clusterID string) ([]*models.AKSTestResult, error)
}

// MemoryStore implements the Store interface using in-memory storage
type MemoryStore struct {
	mu          sync.RWMutex
	clusters    map[string]*models.Cluster
	testResults map[string]*models.TestResult
}

// NewMemoryStore creates a new in-memory store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		clusters:    make(map[string]*models.Cluster),
		testResults: make(map[string]*models.TestResult),
	}
}
		testResults: make(map[string]*models.AKSTestResult),
		nodePools:   make(map[string]*models.AKSNodePool),
	}
}

// AKS Cluster operations
func (s *MemoryStore) CreateCluster(cluster *models.AKSCluster) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.clusters[cluster.ID]; exists {
		return ErrAlreadyExists
	}

	cluster.CreatedAt = time.Now()
	cluster.UpdatedAt = time.Now()
	s.clusters[cluster.ID] = cluster
	return nil
}

func (s *MemoryStore) GetCluster(id string) (*models.AKSCluster, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cluster, exists := s.clusters[id]
	if !exists {
		return nil, ErrNotFound
	}
	return cluster, nil
}

func (s *MemoryStore) ListClusters() ([]*models.AKSCluster, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	clusters := make([]*models.AKSCluster, 0, len(s.clusters))
	for _, cluster := range s.clusters {
		clusters = append(clusters, cluster)
	}
	return clusters, nil
}

func (s *MemoryStore) UpdateCluster(id string, cluster *models.AKSCluster) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.clusters[id]; !exists {
		return ErrNotFound
	}

	cluster.ID = id
	cluster.UpdatedAt = time.Now()
	s.clusters[id] = cluster
	return nil
}

func (s *MemoryStore) DeleteCluster(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.clusters[id]; !exists {
		return ErrNotFound
	}

	delete(s.clusters, id)
	return nil
}

// Test Result operations
func (s *MemoryStore) CreateTestResult(result *models.AKSTestResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.testResults[result.ID]; exists {
		return ErrAlreadyExists
	}

	result.StartedAt = time.Now()
	s.testResults[result.ID] = result
	return nil
}

func (s *MemoryStore) GetTestResult(id string) (*models.AKSTestResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result, exists := s.testResults[id]
	if !exists {
		return nil, ErrNotFound
	}
	return result, nil
}

func (s *MemoryStore) ListTestResults(clusterID string) ([]*models.AKSTestResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []*models.AKSTestResult
	for _, result := range s.testResults {
		if clusterID == "" || result.ClusterID == clusterID {
			results = append(results, result)
		}
	}
	return results, nil
}

func (s *MemoryStore) UpdateTestResult(id string, result *models.AKSTestResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.testResults[id]; !exists {
		return ErrNotFound
	}

	result.ID = id
	s.testResults[id] = result
	return nil
}

// Node Pool operations
func (s *MemoryStore) CreateNodePool(nodePool *models.AKSNodePool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.nodePools[nodePool.ID]; exists {
		return ErrAlreadyExists
	}

	nodePool.CreatedAt = time.Now()
	nodePool.UpdatedAt = time.Now()
	s.nodePools[nodePool.ID] = nodePool
	return nil
}

func (s *MemoryStore) GetNodePool(id string) (*models.AKSNodePool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	nodePool, exists := s.nodePools[id]
	if !exists {
		return nil, ErrNotFound
	}
	return nodePool, nil
}

func (s *MemoryStore) ListNodePools(clusterID string) ([]*models.AKSNodePool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var nodePools []*models.AKSNodePool
	for _, nodePool := range s.nodePools {
		if clusterID == "" || nodePool.ClusterID == clusterID {
			nodePools = append(nodePools, nodePool)
		}
	}
	return nodePools, nil
}

func (s *MemoryStore) UpdateNodePool(id string, nodePool *models.AKSNodePool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.nodePools[id]; !exists {
		return ErrNotFound
	}

	nodePool.ID = id
	nodePool.UpdatedAt = time.Now()
	s.nodePools[id] = nodePool
	return nil
}

func (s *MemoryStore) DeleteNodePool(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.nodePools[id]; !exists {
		return ErrNotFound
	}

	delete(s.nodePools, id)
	return nil
}

// Multi-cloud cluster operations
func (s *MemoryStore) CreateCluster(cluster *models.Cluster) (*models.Cluster, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cluster.ID = uuid.New().String()
	cluster.CreatedAt = time.Now()
	cluster.UpdatedAt = time.Now()

	s.clusters[cluster.ID] = cluster
	return cluster, nil
}

func (s *MemoryStore) GetCluster(id string) (*models.Cluster, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cluster, exists := s.clusters[id]
	if !exists {
		return nil, fmt.Errorf("cluster not found")
	}

	return cluster, nil
}

func (s *MemoryStore) UpdateCluster(id string, cluster *models.Cluster) (*models.Cluster, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.clusters[id]
	if !exists {
		return nil, fmt.Errorf("cluster not found")
	}

	cluster.ID = existing.ID
	cluster.CreatedAt = existing.CreatedAt
	cluster.UpdatedAt = time.Now()

	s.clusters[id] = cluster
	return cluster, nil
}

func (s *MemoryStore) DeleteCluster(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.clusters[id]; !exists {
		return fmt.Errorf("cluster not found")
	}

	delete(s.clusters, id)
	return nil
}

func (s *MemoryStore) ListClusters() ([]*models.Cluster, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	clusters := make([]*models.Cluster, 0, len(s.clusters))
	for _, cluster := range s.clusters {
		clusters = append(clusters, cluster)
	}

	return clusters, nil
}

func (s *MemoryStore) ListClustersByProvider(provider string) ([]*models.Cluster, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var clusters []*models.Cluster
	for _, cluster := range s.clusters {
		if cluster.CloudProvider == provider {
			clusters = append(clusters, cluster)
		}
	}

	return clusters, nil
}

// Test result operations
func (s *MemoryStore) CreateTestResult(result *models.TestResult) (*models.TestResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	result.ID = uuid.New().String()
	result.StartedAt = time.Now()

	s.testResults[result.ID] = result
	return result, nil
}

func (s *MemoryStore) GetTestResult(id string) (*models.TestResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result, exists := s.testResults[id]
	if !exists {
		return nil, fmt.Errorf("test result not found")
	}

	return result, nil
}

func (s *MemoryStore) UpdateTestResult(id string, result *models.TestResult) (*models.TestResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.testResults[id]
	if !exists {
		return nil, fmt.Errorf("test result not found")
	}

	result.ID = existing.ID
	result.StartedAt = existing.StartedAt

	s.testResults[id] = result
	return result, nil
}

func (s *MemoryStore) ListTestResults(clusterID string) ([]*models.TestResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []*models.TestResult
	for _, result := range s.testResults {
		if result.ClusterID == clusterID {
			results = append(results, result)
		}
	}

	return results, nil
}

// Legacy AKS operations (for backward compatibility)
func (s *MemoryStore) CreateAKSCluster(cluster *models.AKSCluster) (*models.AKSCluster, error) {
	// Convert to multi-cloud cluster and create
	multiCluster := &models.Cluster{
		Name:          cluster.Name,
		CloudProvider: "azure",
		Status:        cluster.Status,
		Config: map[string]interface{}{
			"azure_config": map[string]interface{}{
				"resource_group":     cluster.ResourceGroup,
				"location":           cluster.Location,
				"kubernetes_version": cluster.KubernetesVersion,
				"node_count":         cluster.NodeCount,
				"tags":               cluster.Tags,
			},
		},
	}
	
	created, err := s.CreateCluster(multiCluster)
	if err != nil {
		return nil, err
	}
	
	// Convert back to AKS format
	return s.convertToAKSCluster(created), nil
}

func (s *MemoryStore) GetAKSCluster(id string) (*models.AKSCluster, error) {
	cluster, err := s.GetCluster(id)
	if err != nil {
		return nil, err
	}
	
	if cluster.CloudProvider != "azure" {
		return nil, fmt.Errorf("cluster is not an Azure cluster")
	}
	
	return s.convertToAKSCluster(cluster), nil
}

func (s *MemoryStore) UpdateAKSCluster(id string, cluster *models.AKSCluster) (*models.AKSCluster, error) {
	// Convert to multi-cloud cluster and update
	multiCluster := &models.Cluster{
		Name:          cluster.Name,
		CloudProvider: "azure",
		Status:        cluster.Status,
		Config: map[string]interface{}{
			"azure_config": map[string]interface{}{
				"resource_group":     cluster.ResourceGroup,
				"location":           cluster.Location,
				"kubernetes_version": cluster.KubernetesVersion,
				"node_count":         cluster.NodeCount,
				"tags":               cluster.Tags,
			},
		},
	}
	
	updated, err := s.UpdateCluster(id, multiCluster)
	if err != nil {
		return nil, err
	}
	
	return s.convertToAKSCluster(updated), nil
}

func (s *MemoryStore) ListAKSClusters() ([]*models.AKSCluster, error) {
	clusters, err := s.ListClustersByProvider("azure")
	if err != nil {
		return nil, err
	}
	
	aksClusters := make([]*models.AKSCluster, len(clusters))
	for i, cluster := range clusters {
		aksClusters[i] = s.convertToAKSCluster(cluster)
	}
	
	return aksClusters, nil
}

func (s *MemoryStore) CreateAKSTestResult(result *models.AKSTestResult) (*models.AKSTestResult, error) {
	// Convert to multi-cloud test result and create
	multiResult := &models.TestResult{
		ClusterID:   result.ClusterID,
		TestType:    result.TestType,
		Status:      result.Status,
		Duration:    result.Duration,
		Details:     result.Details,
		ErrorMsg:    result.ErrorMsg,
		CompletedAt: result.CompletedAt,
	}
	
	created, err := s.CreateTestResult(multiResult)
	if err != nil {
		return nil, err
	}
	
	return s.convertToAKSTestResult(created), nil
}

func (s *MemoryStore) GetAKSTestResult(id string) (*models.AKSTestResult, error) {
	result, err := s.GetTestResult(id)
	if err != nil {
		return nil, err
	}
	
	return s.convertToAKSTestResult(result), nil
}

func (s *MemoryStore) UpdateAKSTestResult(id string, result *models.AKSTestResult) (*models.AKSTestResult, error) {
	// Convert to multi-cloud test result and update
	multiResult := &models.TestResult{
		ClusterID:   result.ClusterID,
		TestType:    result.TestType,
		Status:      result.Status,
		Duration:    result.Duration,
		Details:     result.Details,
		ErrorMsg:    result.ErrorMsg,
		CompletedAt: result.CompletedAt,
	}
	
	updated, err := s.UpdateTestResult(id, multiResult)
	if err != nil {
		return nil, err
	}
	
	return s.convertToAKSTestResult(updated), nil
}

func (s *MemoryStore) ListAKSTestResults(clusterID string) ([]*models.AKSTestResult, error) {
	results, err := s.ListTestResults(clusterID)
	if err != nil {
		return nil, err
	}
	
	aksResults := make([]*models.AKSTestResult, len(results))
	for i, result := range results {
		aksResults[i] = s.convertToAKSTestResult(result)
	}
	
	return aksResults, nil
}

// Helper methods for conversion
func (s *MemoryStore) convertToAKSCluster(cluster *models.Cluster) *models.AKSCluster {
	aksCluster := &models.AKSCluster{
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
		if nodeCountInt, ok := config["node_count"].(int); ok {
			aksCluster.NodeCount = nodeCountInt
		}
		if tags, ok := config["tags"].(map[string]string); ok {
			aksCluster.Tags = tags
		}
	}
	
	return aksCluster
}

func (s *MemoryStore) convertToAKSTestResult(result *models.TestResult) *models.AKSTestResult {
	return &models.AKSTestResult{
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
