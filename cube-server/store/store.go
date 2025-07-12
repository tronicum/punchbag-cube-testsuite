package store

import (
	"errors"
	"sync"

	sharedmodels "punchbag-cube-testsuite/shared/models"
)

var (
	ErrNotFound      = errors.New("resource not found")
	ErrAlreadyExists = errors.New("resource already exists")
)

// Store defines the interface for data storage operations
type Store interface {
	// Cluster operations (multi-cloud)
	CreateCluster(cluster *sharedmodels.Cluster) error
	GetCluster(id string) (*sharedmodels.Cluster, error)
	ListClusters() ([]*sharedmodels.Cluster, error)
	ListClustersByProvider(provider sharedmodels.CloudProvider) ([]*sharedmodels.Cluster, error)
	UpdateCluster(id string, cluster *sharedmodels.Cluster) error
	DeleteCluster(id string) error

	// Test Result operations
	CreateTestResult(result *sharedmodels.TestResult) error
	GetTestResult(id string) (*sharedmodels.TestResult, error)
	ListTestResults(clusterID string) ([]*sharedmodels.TestResult, error)
	UpdateTestResult(id string, result *sharedmodels.TestResult) error

	// Node Pool operations
	CreateNodePool(nodePool *sharedmodels.NodePool) error
	GetNodePool(id string) (*sharedmodels.NodePool, error)
	ListNodePools(clusterID string) ([]*sharedmodels.NodePool, error)
	UpdateNodePool(id string, nodePool *sharedmodels.NodePool) error
	DeleteNodePool(id string) error

	// Azure-specific operations
	CreateAzureMonitoring(monitoring *sharedmodels.AzureMonitoring) error
	GetAzureMonitoring(id string) (*sharedmodels.AzureMonitoring, error)
	ListAzureMonitorings() ([]*sharedmodels.AzureMonitoring, error)
	UpdateAzureMonitoring(id string, monitoring *sharedmodels.AzureMonitoring) error
	DeleteAzureMonitoring(id string) error

	CreateAzureKubernetes(kubernetes *sharedmodels.AzureKubernetes) error
	GetAzureKubernetes(id string) (*sharedmodels.AzureKubernetes, error)
	ListAzureKubernetes() ([]*sharedmodels.AzureKubernetes, error)
	UpdateAzureKubernetes(id string, kubernetes *sharedmodels.AzureKubernetes) error
	DeleteAzureKubernetes(id string) error

	CreateAzureBudget(budget *sharedmodels.AzureBudget) error
	GetAzureBudget(id string) (*sharedmodels.AzureBudget, error)
	ListAzureBudgets() ([]*sharedmodels.AzureBudget, error)
	UpdateAzureBudget(id string, budget *sharedmodels.AzureBudget) error
	DeleteAzureBudget(id string) error
}

// MemoryStore implements the Store interface using in-memory storage
type MemoryStore struct {
	mu          sync.RWMutex
	clusters    map[string]*sharedmodels.Cluster
	testResults map[string]*sharedmodels.TestResult
	nodePools   map[string]*sharedmodels.NodePool

	// Azure-specific fields
	azureMonitorings map[string]*sharedmodels.AzureMonitoring
	azureKubernetes  map[string]*sharedmodels.AzureKubernetes
	azureBudgets     map[string]*sharedmodels.AzureBudget
}

// NewMemoryStore creates a new in-memory store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		clusters:    make(map[string]*sharedmodels.Cluster),
		testResults: make(map[string]*sharedmodels.TestResult),
		nodePools:   make(map[string]*sharedmodels.NodePool),

		// Azure-specific initializations
		azureMonitorings: make(map[string]*sharedmodels.AzureMonitoring),
		azureKubernetes:  make(map[string]*sharedmodels.AzureKubernetes),
		azureBudgets:     make(map[string]*sharedmodels.AzureBudget),
	}
}

// Cluster operations

func (s *MemoryStore) CreateCluster(cluster *sharedmodels.Cluster) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.clusters[cluster.ID]; exists {
		return ErrAlreadyExists
	}
	s.clusters[cluster.ID] = cluster
	return nil
}

func (s *MemoryStore) GetCluster(id string) (*sharedmodels.Cluster, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cluster, exists := s.clusters[id]
	if !exists {
		return nil, ErrNotFound
	}
	return cluster, nil
}

func (s *MemoryStore) ListClusters() ([]*sharedmodels.Cluster, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var clusters []*sharedmodels.Cluster
	for _, cluster := range s.clusters {
		clusters = append(clusters, cluster)
	}
	return clusters, nil
}

func (s *MemoryStore) ListClustersByProvider(provider sharedmodels.CloudProvider) ([]*sharedmodels.Cluster, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var clusters []*sharedmodels.Cluster
	for _, cluster := range s.clusters {
		if cluster.Provider == provider {
			clusters = append(clusters, cluster)
		}
	}
	return clusters, nil
}

func (s *MemoryStore) UpdateCluster(id string, cluster *sharedmodels.Cluster) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.clusters[id]; !exists {
		return ErrNotFound
	}
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

func (s *MemoryStore) CreateTestResult(result *sharedmodels.TestResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.testResults[result.ID]; exists {
		return ErrAlreadyExists
	}
	s.testResults[result.ID] = result
	return nil
}

func (s *MemoryStore) GetTestResult(id string) (*sharedmodels.TestResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result, exists := s.testResults[id]
	if !exists {
		return nil, ErrNotFound
	}
	return result, nil
}

func (s *MemoryStore) ListTestResults(clusterID string) ([]*sharedmodels.TestResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []*sharedmodels.TestResult
	for _, result := range s.testResults {
		if result.ClusterID == clusterID {
			results = append(results, result)
		}
	}
	return results, nil
}

func (s *MemoryStore) UpdateTestResult(id string, result *sharedmodels.TestResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.testResults[id]; !exists {
		return ErrNotFound
	}
	s.testResults[id] = result
	return nil
}

// Node Pool operations

func (s *MemoryStore) CreateNodePool(nodePool *sharedmodels.NodePool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.nodePools[nodePool.ID]; exists {
		return ErrAlreadyExists
	}
	s.nodePools[nodePool.ID] = nodePool
	return nil
}

func (s *MemoryStore) GetNodePool(id string) (*sharedmodels.NodePool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	nodePool, exists := s.nodePools[id]
	if !exists {
		return nil, ErrNotFound
	}
	return nodePool, nil
}

func (s *MemoryStore) ListNodePools(clusterID string) ([]*sharedmodels.NodePool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var nodePools []*sharedmodels.NodePool
	for _, nodePool := range s.nodePools {
		if nodePool.ClusterID == clusterID {
			nodePools = append(nodePools, nodePool)
		}
	}
	return nodePools, nil
}

func (s *MemoryStore) UpdateNodePool(id string, nodePool *sharedmodels.NodePool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.nodePools[id]; !exists {
		return ErrNotFound
	}
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

// Azure-specific operations

func (s *MemoryStore) CreateAzureMonitoring(monitoring *sharedmodels.AzureMonitoring) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.azureMonitorings[monitoring.ID]; exists {
		return ErrAlreadyExists
	}
	s.azureMonitorings[monitoring.ID] = monitoring
	return nil
}

func (s *MemoryStore) GetAzureMonitoring(id string) (*sharedmodels.AzureMonitoring, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	monitoring, exists := s.azureMonitorings[id]
	if !exists {
		return nil, ErrNotFound
	}
	return monitoring, nil
}

func (s *MemoryStore) ListAzureMonitorings() ([]*sharedmodels.AzureMonitoring, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var monitorings []*sharedmodels.AzureMonitoring
	for _, monitoring := range s.azureMonitorings {
		monitorings = append(monitorings, monitoring)
	}
	return monitorings, nil
}

func (s *MemoryStore) UpdateAzureMonitoring(id string, monitoring *sharedmodels.AzureMonitoring) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.azureMonitorings[id]; !exists {
		return ErrNotFound
	}
	s.azureMonitorings[id] = monitoring
	return nil
}

func (s *MemoryStore) DeleteAzureMonitoring(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.azureMonitorings[id]; !exists {
		return ErrNotFound
	}
	delete(s.azureMonitorings, id)
	return nil
}

func (s *MemoryStore) CreateAzureKubernetes(kubernetes *sharedmodels.AzureKubernetes) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.azureKubernetes[kubernetes.ID]; exists {
		return ErrAlreadyExists
	}
	s.azureKubernetes[kubernetes.ID] = kubernetes
	return nil
}

func (s *MemoryStore) GetAzureKubernetes(id string) (*sharedmodels.AzureKubernetes, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	kubernetes, exists := s.azureKubernetes[id]
	if !exists {
		return nil, ErrNotFound
	}
	return kubernetes, nil
}

func (s *MemoryStore) ListAzureKubernetes() ([]*sharedmodels.AzureKubernetes, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var kubernetesList []*sharedmodels.AzureKubernetes
	for _, kubernetes := range s.azureKubernetes {
		kubernetesList = append(kubernetesList, kubernetes)
	}
	return kubernetesList, nil
}

func (s *MemoryStore) UpdateAzureKubernetes(id string, kubernetes *sharedmodels.AzureKubernetes) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.azureKubernetes[id]; !exists {
		return ErrNotFound
	}
	s.azureKubernetes[id] = kubernetes
	return nil
}

func (s *MemoryStore) DeleteAzureKubernetes(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.azureKubernetes[id]; !exists {
		return ErrNotFound
	}
	delete(s.azureKubernetes, id)
	return nil
}

func (s *MemoryStore) CreateAzureBudget(budget *sharedmodels.AzureBudget) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.azureBudgets[budget.ID]; exists {
		return ErrAlreadyExists
	}
	s.azureBudgets[budget.ID] = budget
	return nil
}

func (s *MemoryStore) GetAzureBudget(id string) (*sharedmodels.AzureBudget, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	budget, exists := s.azureBudgets[id]
	if !exists {
		return nil, ErrNotFound
	}
	return budget, nil
}

func (s *MemoryStore) ListAzureBudgets() ([]*sharedmodels.AzureBudget, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var budgets []*sharedmodels.AzureBudget
	for _, budget := range s.azureBudgets {
		budgets = append(budgets, budget)
	}
	return budgets, nil
}

func (s *MemoryStore) UpdateAzureBudget(id string, budget *sharedmodels.AzureBudget) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.azureBudgets[id]; !exists {
		return ErrNotFound
	}
	s.azureBudgets[id] = budget
	return nil
}

func (s *MemoryStore) DeleteAzureBudget(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.azureBudgets[id]; !exists {
		return ErrNotFound
	}
	delete(s.azureBudgets, id)
	return nil
}
