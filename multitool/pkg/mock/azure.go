package mock

import (
	"time"
	"punchbag-cube-testsuite/multitool/pkg/models"
)

// In-memory stores for mocks
var (
	mockAksClusters         = []*models.Cluster{}
	mockLogAnalytics        = []*models.LogAnalyticsWorkspace{}
	mockBudgets             = []*models.AzureBudget{}
	mockAppInsights         = []*models.AppInsightsResource{}
)

// AKS
func MockCreateAks(cluster *models.Cluster) *models.Cluster {
	cluster.ID = "aks-mock-" + time.Now().Format("150405")
	cluster.Provider = models.Azure
	cluster.Status = models.ClusterStatusRunning
	cluster.CreatedAt = time.Now()
	mockAksClusters = append(mockAksClusters, cluster)
	return cluster
}
func MockListAks() []*models.Cluster { return mockAksClusters }
func MockGetAks(id string) *models.Cluster {
	for _, c := range mockAksClusters {
		if c.ID == id { return c }
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
func MockCreateLogAnalytics(ws *models.LogAnalyticsWorkspace) *models.LogAnalyticsWorkspace {
	ws.ID = "log-mock-" + time.Now().Format("150405")
	ws.CreatedAt = time.Now()
	mockLogAnalytics = append(mockLogAnalytics, ws)
	return ws
}
func MockListLogAnalytics() []*models.LogAnalyticsWorkspace { return mockLogAnalytics }
func MockGetLogAnalytics(id string) *models.LogAnalyticsWorkspace {
	for _, ws := range mockLogAnalytics {
		if ws.ID == id { return ws }
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
func MockCreateBudget(b *models.AzureBudget) *models.AzureBudget {
	b.ID = "budget-mock-" + time.Now().Format("150405")
	b.CreatedAt = time.Now()
	mockBudgets = append(mockBudgets, b)
	return b
}
func MockListBudgets() []*models.AzureBudget { return mockBudgets }
func MockGetBudget(id string) *models.AzureBudget {
	for _, b := range mockBudgets {
		if b.ID == id { return b }
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

// App Insights
func MockCreateAppInsights(a *models.AppInsightsResource) *models.AppInsightsResource {
	a.ID = "app-mock-" + time.Now().Format("150405")
	a.CreatedAt = time.Now()
	mockAppInsights = append(mockAppInsights, a)
	return a
}
func MockListAppInsights() []*models.AppInsightsResource { return mockAppInsights }
func MockGetAppInsights(id string) *models.AppInsightsResource {
	for _, a := range mockAppInsights {
		if a.ID == id { return a }
	}
	return nil
}
func MockDeleteAppInsights(id string) bool {
	for i, a := range mockAppInsights {
		if a.ID == id {
			mockAppInsights = append(mockAppInsights[:i], mockAppInsights[i+1:]...)
			return true
		}
	}
	return false
}
