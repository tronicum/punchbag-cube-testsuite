package azure

import (
    "context"
    "fmt"
)

// AzureProviderImpl implements the AzureProvider interface
type AzureProviderImpl struct {
    simulationMode bool
    subscriptionID string
    tenantID       string
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
        return &MonitorResult{
            ID:        fmt.Sprintf("monitor-%s-%s", resourceGroup, workspaceName),
            Status:    "created",
            Resources: []string{"log-analytics", "application-insights", "metrics"},
        }, nil
    }
    return nil, fmt.Errorf("direct mode not implemented yet")
}

// CreateBudget creates an Azure budget
func (p *AzureProviderImpl) CreateBudget(ctx context.Context, name string, amount float64, resourceGroup, timeGrain string) (*BudgetResult, error) {
    if p.simulationMode {
        return &BudgetResult{
            ID:     fmt.Sprintf("budget-%s-%s", resourceGroup, name),
            Status: "created",
        }, nil
    }
    return nil, fmt.Errorf("direct mode not implemented yet")
}

// CreateAKSCluster creates an AKS cluster
func (p *AzureProviderImpl) CreateAKSCluster(ctx context.Context, name, resourceGroup, location string, nodeCount int) (*ClusterResult, error) {
    if p.simulationMode {
        return &ClusterResult{
            ID:     fmt.Sprintf("cluster-%s-%s", resourceGroup, name),
            Status: "created",
            URL:    fmt.Sprintf("https://%s.azmk8s.io", name),
        }, nil
    }
    return nil, fmt.Errorf("direct mode not implemented yet")
}