package mock

import (
	"time"
	"punchbag-cube-testsuite/multitool/pkg/models"
)

// In-memory stores for AWS mocks
var (
	mockEksClusters = []*models.Cluster{}
	mockS3Buckets   = []*models.S3Bucket{}
)

// EKS
func MockCreateEks(cluster *models.Cluster) *models.Cluster {
	cluster.ID = "eks-mock-" + time.Now().Format("150405")
	cluster.Provider = models.AWS
	cluster.Status = models.ClusterStatusRunning
	cluster.CreatedAt = time.Now()
	mockEksClusters = append(mockEksClusters, cluster)
	return cluster
}
func MockListEks() []*models.Cluster { return mockEksClusters }
func MockGetEks(id string) *models.Cluster {
	for _, c := range mockEksClusters {
		if c.ID == id { return c }
	}
	return nil
}
func MockDeleteEks(id string) bool {
	for i, c := range mockEksClusters {
		if c.ID == id {
			mockEksClusters = append(mockEksClusters[:i], mockEksClusters[i+1:]...)
			return true
		}
	}
	return false
}

// S3
func MockCreateS3(bucket *models.S3Bucket) *models.S3Bucket {
	bucket.ID = "s3-mock-" + time.Now().Format("150405")
	bucket.CreatedAt = time.Now().Format(time.RFC3339)
	mockS3Buckets = append(mockS3Buckets, bucket)
	return bucket
}
func MockListS3() []*models.S3Bucket { return mockS3Buckets }
func MockGetS3(id string) *models.S3Bucket {
	for _, b := range mockS3Buckets {
		if b.ID == id { return b }
	}
	return nil
}
func MockDeleteS3(id string) bool {
	for i, b := range mockS3Buckets {
		if b.ID == id {
			mockS3Buckets = append(mockS3Buckets[:i], mockS3Buckets[i+1:]...)
			return true
		}
	}
	return false
}
// Add similar mocks for AWS CloudWatch and Budget as needed.
