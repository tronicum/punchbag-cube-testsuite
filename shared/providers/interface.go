package providers

import "context"

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
	CreateMonitor(ctx context.Context, resourceGroup, location, workspaceName string) (string, error)
	
	// Budget operations
	CreateBudget(ctx context.Context, name string, amount float64, resourceGroup, timeGrain string) (string, error)
	
	// AKS operations
	CreateAKSCluster(ctx context.Context, name, resourceGroup, location string, nodeCount int) (string, error)
}
	Provider
	CreateMonitor(ctx context.Context, params MonitorParams) (*MonitorResult, error)
	DeleteMonitor(ctx context.Context, monitorID string) error
	ListMonitors(ctx context.Context) ([]*MonitorInfo, error)
}

// BudgetProvider defines budget-related operations
type BudgetProvider interface {
	Provider
	CreateBudget(ctx context.Context, params BudgetParams) (*BudgetResult, error)
	DeleteBudget(ctx context.Context, budgetID string) error
	ListBudgets(ctx context.Context) ([]*BudgetInfo, error)
}

// Common parameter and result types
type ClusterParams struct {
	Name          string            `json:"name"`
	ResourceGroup string            `json:"resource_group"`
	Location      string            `json:"location"`
	NodeCount     int               `json:"node_count"`
	Tags          map[string]string `json:"tags,omitempty"`
}

type ClusterResult struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	URL    string `json:"url,omitempty"`
}

type ClusterInfo struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Status        string            `json:"status"`
	Location      string            `json:"location"`
	NodeCount     int               `json:"node_count"`
	ResourceGroup string            `json:"resource_group,omitempty"`
	Tags          map[string]string `json:"tags,omitempty"`
}

type MonitorParams struct {
	Name          string `json:"name"`
	ResourceGroup string `json:"resource_group"`
	Location      string `json:"location"`
	WorkspaceName string `json:"workspace_name"`
}

type MonitorResult struct {
	ID        string   `json:"id"`
	Status    string   `json:"status"`
	Resources []string `json:"resources"`
}

type MonitorInfo struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Status        string `json:"status"`
	ResourceGroup string `json:"resource_group"`
	Location      string `json:"location"`
}

type BudgetParams struct {
	Name          string  `json:"name"`
	Amount        float64 `json:"amount"`
	ResourceGroup string  `json:"resource_group"`
	TimeGrain     string  `json:"time_grain"`
}

type BudgetResult struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type BudgetInfo struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Amount        float64 `json:"amount"`
	ResourceGroup string  `json:"resource_group"`
	TimeGrain     string  `json:"time_grain"`
}
