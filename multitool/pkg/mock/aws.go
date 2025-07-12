package mock

import (
	"time"

	sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"
)

// In-memory stores for AWS mocks
var (
	mockEksClusters = []*sharedmodels.Cluster{}
	mockS3Buckets   = []*sharedmodels.ObjectStorageBucket{}
)

// In-memory stores for AKS, Log Analytics, and Budget mocks
var (
	mockAksClusters  = []*sharedmodels.Cluster{}
	mockLogAnalytics = []*sharedmodels.LogAnalyticsWorkspace{}
	mockBudgets      = []*sharedmodels.AzureBudget{}
)

// In-memory store for Blob Storage mocks
var (
	mockBlobBuckets = []*sharedmodels.ObjectStorageBucket{}
)

// In-memory store for App Insights mocks
var (
	mockAppInsights = []*sharedmodels.AppInsightsResource{}
)

// In-memory store for GCS mocks
var (
	mockGcsBuckets = []*sharedmodels.ObjectStorageBucket{}
)

// EKS
func MockCreateEks(cluster *sharedmodels.Cluster) *sharedmodels.Cluster {
	cluster.ID = "eks-mock-" + time.Now().Format("150405")
	cluster.Provider = sharedmodels.AWS
	cluster.Status = sharedmodels.ClusterStatusRunning
	cluster.CreatedAt = time.Now()
	mockEksClusters = append(mockEksClusters, cluster)
	return cluster
}
func MockListEks() []*sharedmodels.Cluster { return mockEksClusters }
func MockGetEks(id string) *sharedmodels.Cluster {
	for _, c := range mockEksClusters {
		if c.ID == id {
			return c
		}
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
func MockCreateS3(bucket *sharedmodels.ObjectStorageBucket) *sharedmodels.ObjectStorageBucket {
	bucket.ID = "s3-mock-" + time.Now().Format("150405")
	bucket.CreatedAt = time.Now()
	bucket.Provider = "aws"
	mockS3Buckets = append(mockS3Buckets, bucket)
	return bucket
}
func MockListS3() []*sharedmodels.ObjectStorageBucket { return mockS3Buckets }
func MockGetS3(id string) *sharedmodels.ObjectStorageBucket {
	for _, b := range mockS3Buckets {
		if b.ID == id {
			return b
		}
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

// AKS
func MockCreateAks(cluster *sharedmodels.Cluster) *sharedmodels.Cluster {
	cluster.ID = "aks-mock-" + time.Now().Format("150405")
	cluster.Provider = sharedmodels.Azure
	cluster.Status = sharedmodels.ClusterStatusRunning
	cluster.CreatedAt = time.Now()
	mockAksClusters = append(mockAksClusters, cluster)
	return cluster
}
func MockListAks() []*sharedmodels.Cluster { return mockAksClusters }
func MockGetAks(id string) *sharedmodels.Cluster {
	for _, c := range mockAksClusters {
		if c.ID == id {
			return c
		}
	}
	return nil
}
func MockDeleteAks(id string) bool {
	for i, c := range mockAksClusters {
		if c.ID == id {
			mockAksClusters = append(mockAksClusters[:i], mockAksClusters[i+1:]...)
			return true
		}
	}
	return false
}

// Log Analytics
func MockCreateLogAnalytics(ws *sharedmodels.LogAnalyticsWorkspace) *sharedmodels.LogAnalyticsWorkspace {
	ws.ID = "loganalytics-mock-" + time.Now().Format("150405")
	ws.CreatedAt = time.Now()
	mockLogAnalytics = append(mockLogAnalytics, ws)
	return ws
}
func MockListLogAnalytics() []*sharedmodels.LogAnalyticsWorkspace { return mockLogAnalytics }
func MockGetLogAnalytics(id string) *sharedmodels.LogAnalyticsWorkspace {
	for _, ws := range mockLogAnalytics {
		if ws.ID == id {
			return ws
		}
	}
	return nil
}
func MockDeleteLogAnalytics(id string) bool {
	for i, ws := range mockLogAnalytics {
		if ws.ID == id {
			mockLogAnalytics = append(mockLogAnalytics[:i], mockLogAnalytics[i+1:]...)
			return true
		}
	}
	return false
}

// Budget
func MockCreateBudget(budget *sharedmodels.AzureBudget) *sharedmodels.AzureBudget {
	budget.ID = "budget-mock-" + time.Now().Format("150405")
	budget.CreatedAt = time.Now()
	mockBudgets = append(mockBudgets, budget)
	return budget
}
func MockListBudgets() []*sharedmodels.AzureBudget { return mockBudgets }
func MockGetBudget(id string) *sharedmodels.AzureBudget {
	for _, b := range mockBudgets {
		if b.ID == id {
			return b
		}
	}
	return nil
}
func MockDeleteBudget(id string) bool {
	for i, b := range mockBudgets {
		if b.ID == id {
			mockBudgets = append(mockBudgets[:i], mockBudgets[i+1:]...)
			return true
		}
	}
	return false
}

// Blob Storage
func MockCreateBlobStorage(bucket *sharedmodels.ObjectStorageBucket) *sharedmodels.ObjectStorageBucket {
	bucket.ID = "blob-mock-" + time.Now().Format("150405")
	bucket.Provider = sharedmodels.Azure
	bucket.CreatedAt = time.Now()
	mockBlobBuckets = append(mockBlobBuckets, bucket)
	return bucket
}
func MockListBlobStorage() []*sharedmodels.ObjectStorageBucket { return mockBlobBuckets }
func MockGetBlobStorage(id string) *sharedmodels.ObjectStorageBucket {
	for _, b := range mockBlobBuckets {
		if b.ID == id {
			return b
		}
	}
	return nil
}
func MockDeleteBlobStorage(id string) bool {
	for i, b := range mockBlobBuckets {
		if b.ID == id {
			mockBlobBuckets = append(mockBlobBuckets[:i], mockBlobBuckets[i+1:]...)
			return true
		}
	}
	return false
}

// App Insights
func MockCreateAppInsights(ai *sharedmodels.AppInsightsResource) *sharedmodels.AppInsightsResource {
	ai.ID = "appinsights-mock-" + time.Now().Format("150405")
	ai.CreatedAt = time.Now()
	mockAppInsights = append(mockAppInsights, ai)
	return ai
}
func MockListAppInsights() []*sharedmodels.AppInsightsResource { return mockAppInsights }
func MockGetAppInsights(id string) *sharedmodels.AppInsightsResource {
	for _, ai := range mockAppInsights {
		if ai.ID == id {
			return ai
		}
	}
	return nil
}
func MockDeleteAppInsights(id string) bool {
	for i, ai := range mockAppInsights {
		if ai.ID == id {
			mockAppInsights = append(mockAppInsights[:i], mockAppInsights[i+1:]...)
			return true
		}
	}
	return false
}

// GCS
func MockCreateGCS(bucket *sharedmodels.ObjectStorageBucket) *sharedmodels.ObjectStorageBucket {
	bucket.ID = "gcs-mock-" + time.Now().Format("150405")
	bucket.Provider = sharedmodels.GCP
	bucket.CreatedAt = time.Now()
	mockGcsBuckets = append(mockGcsBuckets, bucket)
	return bucket
}
func MockListGCS() []*sharedmodels.ObjectStorageBucket { return mockGcsBuckets }
func MockGetGCS(id string) *sharedmodels.ObjectStorageBucket {
	for _, b := range mockGcsBuckets {
		if b.ID == id {
			return b
		}
	}
	return nil
}
func MockDeleteGCS(id string) bool {
	for i, b := range mockGcsBuckets {
		if b.ID == id {
			mockGcsBuckets = append(mockGcsBuckets[:i], mockGcsBuckets[i+1:]...)
			return true
		}
	}
	return false
}

// Add similar mocks for AWS CloudWatch and Budget as needed.
