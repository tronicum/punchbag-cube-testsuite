package providers

import (
    "context"
    "github.com/tronicum/punchbag-cube-testsuite/shared/types"
)

// Provider defines the common interface for all cloud providers
type Provider interface {
    // GetName returns the provider name (azure, aws, gcp, etc.)
    GetName() string
    
    // SimulationMode returns true if in simulation mode
    SimulationMode() bool
    
    // SetSimulationMode enables/disables simulation mode
    SetSimulationMode(enabled bool)
}

// AzureProvider defines Azure-specific operations
type AzureProvider interface {
    Provider
    
    // Monitor operations
    CreateMonitor(ctx context.Context, resourceGroup, location, workspaceName string) (*types.MonitorResult, error)
    
    // Budget operations
    CreateBudget(ctx context.Context, name string, amount float64, resourceGroup, timeGrain string) (*types.BudgetResult, error)
    
    // AKS operations
    CreateAKSCluster(ctx context.Context, name, resourceGroup, location string, nodeCount int) (*types.ClusterResult, error)
}
