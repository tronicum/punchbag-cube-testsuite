package azure

import (
	"context"
	"fmt"

	sharederrors "github.com/tronicum/punchbag-cube-testsuite/shared/errors"
	"github.com/tronicum/punchbag-cube-testsuite/shared/log"
)

// MonitorResult represents the result of an Azure Monitor operation
type MonitorResult struct {
	ID        string   `json:"id"`
	Status    string   `json:"status"`
	Resources []string `json:"resources"`
}

// BudgetResult represents the result of an Azure Budget operation
type BudgetResult struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// ClusterResult represents the result of an AKS cluster operation
type ClusterResult struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	URL    string `json:"url,omitempty"`
}

// AzureProviderImpl implements the AzureProvider interface
type AzureProviderImpl struct {
	simulationMode bool
	subscriptionID string
	tenantID       string
	creds map[string]string
	config map[string]interface{}
	region string
	endpoints map[string]string
}
// SetCredentials sets authentication credentials
func (p *AzureProviderImpl) SetCredentials(creds map[string]string) { p.creds = creds }

// SetConfig sets provider-specific config
func (p *AzureProviderImpl) SetConfig(cfg map[string]interface{}) {
	p.config = cfg
	if r, ok := cfg["region"].(string); ok { p.region = r }
	if e, ok := cfg["endpoints"].(map[string]string); ok { p.endpoints = e }
}

// GetRegion returns the configured region
func (p *AzureProviderImpl) GetRegion() string { return p.region }

// GetEndpoint returns the endpoint for a given service
func (p *AzureProviderImpl) GetEndpoint(service string) string {
	if p.endpoints != nil {
		if ep, ok := p.endpoints[service]; ok { return ep }
	}
	// Default Azure endpoints (example)
	if service == "aks" && p.region != "" {
		return "https://management.azure.com/subscriptions/" + p.subscriptionID + "/resourceGroups/" + p.region + "/providers/Microsoft.ContainerService/managedClusters"
	}
	return ""
}

// NewAzureProvider creates a new Azure provider
func NewAzureProvider() AzureProvider {
	return &AzureProviderImpl{
		simulationMode: false,
	}
}

// GetName returns the provider name
func (p *AzureProviderImpl) GetName() string {
	return "azure"
}

// SimulationMode returns true if in simulation mode
func (p *AzureProviderImpl) SimulationMode() bool {
	return p.simulationMode
}

// SetSimulationMode enables/disables simulation mode
func (p *AzureProviderImpl) SetSimulationMode(enabled bool) {
	p.simulationMode = enabled
}

// CreateMonitor creates Azure Monitor resources
func (p *AzureProviderImpl) CreateMonitor(ctx context.Context, resourceGroup, location, workspaceName string) (*MonitorResult, error) {
	if p.simulationMode {
		log.Info("Simulating Azure Monitor creation for resourceGroup=%s, workspaceName=%s", resourceGroup, workspaceName)
		return &MonitorResult{
			ID:        fmt.Sprintf("monitor-%s-%s", resourceGroup, workspaceName),
			Status:    "created",
			Resources: []string{"log-analytics", "application-insights", "metrics"},
		}, nil
	}
	log.Warn("Direct mode not implemented for CreateMonitor")
	return nil, sharederrors.ErrNotFound
}

// CreateBudget creates an Azure budget
func (p *AzureProviderImpl) CreateBudget(ctx context.Context, name string, amount float64, resourceGroup, timeGrain string) (*BudgetResult, error) {
	if p.simulationMode {
		log.Info("Simulating Azure Budget creation for resourceGroup=%s, name=%s", resourceGroup, name)
		return &BudgetResult{
			ID:     fmt.Sprintf("budget-%s-%s", resourceGroup, name),
			Status: "created",
		}, nil
	}
	log.Warn("Direct mode not implemented for CreateBudget")
	return nil, sharederrors.ErrNotFound
}

// CreateAKSCluster creates an AKS cluster
func (p *AzureProviderImpl) CreateAKSCluster(ctx context.Context, name, resourceGroup, location string, nodeCount int) (*ClusterResult, error) {
	if p.simulationMode {
		log.Info("Simulating AKS cluster creation for resourceGroup=%s, name=%s", resourceGroup, name)
		return &ClusterResult{
			ID:     fmt.Sprintf("cluster-%s-%s", resourceGroup, name),
			Status: "created",
			URL:    fmt.Sprintf("https://%s.azmk8s.io", name),
		}, nil
	}
	log.Warn("Direct mode not implemented for CreateAKSCluster")
	return nil, sharederrors.ErrNotFound
}
