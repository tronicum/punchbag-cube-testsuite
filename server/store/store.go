package store

import (
	"errors"
	"sync"
	"time"

	"punchbag-cube-testsuite/server/models"
)

var (
	ErrNotFound      = errors.New("resource not found")
	ErrAlreadyExists = errors.New("resource already exists")
)

// Store defines the interface for data storage operations
type Store interface {
	// AKS Cluster operations
	CreateCluster(cluster *models.AKSCluster) error
	GetCluster(id string) (*models.AKSCluster, error)
	ListClusters() ([]*models.AKSCluster, error)
	UpdateCluster(id string, cluster *models.AKSCluster) error
	DeleteCluster(id string) error

	// Test Result operations
	CreateTestResult(result *models.AKSTestResult) error
	GetTestResult(id string) (*models.AKSTestResult, error)
	ListTestResults(clusterID string) ([]*models.AKSTestResult, error)
	UpdateTestResult(id string, result *models.AKSTestResult) error

	// Node Pool operations
	CreateNodePool(nodePool *models.AKSNodePool) error
	GetNodePool(id string) (*models.AKSNodePool, error)
	ListNodePools(clusterID string) ([]*models.AKSNodePool, error)
	UpdateNodePool(id string, nodePool *models.AKSNodePool) error
	DeleteNodePool(id string) error
}

// MemoryStore implements the Store interface using in-memory storage
type MemoryStore struct {
	mu          sync.RWMutex
	clusters    map[string]*models.AKSCluster
	testResults map[string]*models.AKSTestResult
	nodePools   map[string]*models.AKSNodePool
}

// NewMemoryStore creates a new in-memory store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		clusters:    make(map[string]*models.AKSCluster),
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
