package types

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
