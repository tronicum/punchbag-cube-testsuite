package azure

import (
	"context"
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
	CreateMonitor(ctx context.Context, resourceGroup, location, workspaceName string) (*MonitorResult, error)
	// Budget operations
	CreateBudget(ctx context.Context, name string, amount float64, resourceGroup, timeGrain string) (*BudgetResult, error)
	// AKS operations
	CreateAKSCluster(ctx context.Context, name, resourceGroup, location string, nodeCount int) (*ClusterResult, error)
}

// ClusterProvider defines cluster-related operations

// type ClusterProvider interface {
//    Provider
//    CreateCluster(ctx context.Context, params ClusterParams) (*ClusterResult, error)
//    DeleteCluster(ctx context.Context, clusterID string) error
//    ListClusters(ctx context.Context) ([]*ClusterInfo, error)
// }

// MonitorProvider defines monitoring-related operations

// type MonitorProvider interface {
//    Provider
//    CreateMonitor(ctx context.Context, params MonitorParams) (*MonitorResult, error)
//    DeleteMonitor(ctx context.Context, monitorID string) error
//    ListMonitors(ctx context.Context) ([]*MonitorInfo, error)
// }

// BudgetProvider defines budget-related operations

// type BudgetProvider interface {
//    Provider
//    CreateBudget(ctx context.Context, params BudgetParams) (*BudgetResult, error)
//    DeleteBudget(ctx context.Context, budgetID string) error
//    ListBudgets(ctx context.Context) ([]*BudgetInfo, error)
// }
