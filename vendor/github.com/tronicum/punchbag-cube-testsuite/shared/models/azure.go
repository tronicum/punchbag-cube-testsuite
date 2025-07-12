package models

import "time"

// LogAnalyticsWorkspace represents an Azure Log Analytics workspace
type LogAnalyticsWorkspace struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	ResourceGroup string    `json:"resource_group"`
	Location      string    `json:"location"`
	CustomerID    string    `json:"customer_id"`
	Sku           string    `json:"sku"`
	RetentionDays int       `json:"retention_days"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// AzureBudget represents an Azure Budget resource
type AzureBudget struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	ResourceGroup string    `json:"resource_group"`
	Amount        float64   `json:"amount"`
	TimeGrain     string    `json:"time_grain"`
	StartDate     string    `json:"start_date"`
	EndDate       string    `json:"end_date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// AppInsightsResource represents an Azure Application Insights instance
type AppInsightsResource struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	ResourceGroup string    `json:"resource_group"`
	Location      string    `json:"location"`
	AppType       string    `json:"app_type"`
	RetentionDays int       `json:"retention_days"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
