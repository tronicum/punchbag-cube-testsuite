package mock

import (
	"time"
	"punchbag-cube-testsuite/multitool/pkg/models"
)

// In-memory stores for GCP mocks
var (
	mockGkeClusters = []*models.Cluster{}
)

// GKE
func MockCreateGke(cluster *models.Cluster) *models.Cluster {
	cluster.ID = "gke-mock-" + time.Now().Format("150405")
	cluster.Provider = models.GCP
	cluster.Status = models.ClusterStatusRunning
	cluster.CreatedAt = time.Now()
	mockGkeClusters = append(mockGkeClusters, cluster)
	return cluster
}
func MockListGke() []*models.Cluster { return mockGkeClusters }
func MockGetGke(id string) *models.Cluster {
	for _, c := range mockGkeClusters {
		if c.ID == id { return c }
	}
	return nil
}
func MockDeleteGke(id string) bool {
	for i, c := range mockGkeClusters {
		if c.ID == id {
			mockGkeClusters = append(mockGkeClusters[:i], mockGkeClusters[i+1:]...)
			return true
		}
	}
	return false
}
// Add similar mocks for GCP Logging, CloudSQL, and Budget as needed.
