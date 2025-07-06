package store

import (
	"fmt"
	"sync"
	"time"

	sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"

	"github.com/google/uuid"
)

var (
	ErrNotFound      = fmt.Errorf("resource not found")
	ErrAlreadyExists = fmt.Errorf("resource already exists")
)

// Store defines the interface for cluster and test result storage
type Store interface {
	// Cluster operations
	CreateCluster(cluster *sharedmodels.Cluster) (*sharedmodels.Cluster, error)
	GetCluster(id string) (*sharedmodels.Cluster, error)
	UpdateCluster(id string, cluster *sharedmodels.Cluster) (*sharedmodels.Cluster, error)
	DeleteCluster(id string) error
	ListClusters() ([]*sharedmodels.Cluster, error)
	ListClustersByProvider(provider string) ([]*sharedmodels.Cluster, error)

	// Test result operations
	CreateTestResult(result *sharedmodels.TestResult) (*sharedmodels.TestResult, error)
	GetTestResult(id string) (*sharedmodels.TestResult, error)
	UpdateTestResult(id string, result *sharedmodels.TestResult) (*sharedmodels.TestResult, error)
	ListTestResults(clusterID string) ([]*sharedmodels.TestResult, error)
}

// MemoryStore implements the Store interface using in-memory storage
type MemoryStore struct {
	mu          sync.RWMutex
	clusters    map[string]*sharedmodels.Cluster
	testResults map[string]*sharedmodels.TestResult
}

// NewMemoryStore creates a new in-memory store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		clusters:    make(map[string]*sharedmodels.Cluster),
		testResults: make(map[string]*sharedmodels.TestResult),
	}
}

// Cluster operations
func (s *MemoryStore) CreateCluster(cluster *sharedmodels.Cluster) (*sharedmodels.Cluster, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.clusters[cluster.ID]; exists {
		return nil, ErrAlreadyExists
	}

	cluster.CreatedAt = time.Now()
	cluster.UpdatedAt = time.Now()
	s.clusters[cluster.ID] = cluster
	return cluster, nil
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

func (s *MemoryStore) UpdateCluster(id string, cluster *sharedmodels.Cluster) (*sharedmodels.Cluster, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.clusters[id]
	if !exists {
		return nil, ErrNotFound
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
		return ErrNotFound
	}

	delete(s.clusters, id)
	return nil
}

func (s *MemoryStore) ListClusters() ([]*sharedmodels.Cluster, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	clusters := make([]*sharedmodels.Cluster, 0, len(s.clusters))
	for _, cluster := range s.clusters {
		clusters = append(clusters, cluster)
	}

	return clusters, nil
}

// Test result operations
func (s *MemoryStore) CreateTestResult(result *sharedmodels.TestResult) (*sharedmodels.TestResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	result.ID = uuid.New().String()
	result.StartedAt = time.Now()

	s.testResults[result.ID] = result
	return result, nil
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

func (s *MemoryStore) UpdateTestResult(id string, result *sharedmodels.TestResult) (*sharedmodels.TestResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.testResults[id]
	if !exists {
		return nil, ErrNotFound
	}

	result.ID = existing.ID
	result.StartedAt = existing.StartedAt
	result.CompletedAt = result.CompletedAt

	s.testResults[id] = result
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

// ListClustersByProvider lists clusters by the given provider
func (s *MemoryStore) ListClustersByProvider(provider string) ([]*sharedmodels.Cluster, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var clusters []*sharedmodels.Cluster
	for _, cluster := range s.clusters {
		if string(cluster.Provider) == provider {
			clusters = append(clusters, cluster)
		}
	}
	return clusters, nil
}
