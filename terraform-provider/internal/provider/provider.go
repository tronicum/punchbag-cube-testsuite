package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &punchbagProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &punchbagProvider{}
}

// punchbagProvider is the provider implementation.
type punchbagProvider struct{}

// punchbagProviderModel maps provider schema data to a Go type.
type punchbagProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// Metadata returns the provider type name.
func (p *punchbagProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "punchbag"
}

// Schema defines the provider-level schema for configuration data.
func (p *punchbagProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Punchbag Cube Test Suite.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "URI for Punchbag API. May also be provided via PUNCHBAG_HOST environment variable.",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username for Punchbag API. May also be provided via PUNCHBAG_USERNAME environment variable.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password for Punchbag API. May also be provided via PUNCHBAG_PASSWORD environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares a Punchbag API client for data sources and resources.
func (p *punchbagProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config punchbagProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			"host",
			"Unknown Punchbag API Host",
			"The provider cannot create the Punchbag API client as there is an unknown configuration value for the Punchbag API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PUNCHBAG_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			"username",
			"Unknown Punchbag API Username",
			"The provider cannot create the Punchbag API client as there is an unknown configuration value for the Punchbag API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PUNCHBAG_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			"password",
			"Unknown Punchbag API Password",
			"The provider cannot create the Punchbag API client as there is an unknown configuration value for the Punchbag API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PUNCHBAG_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := "http://localhost:8080"
	username := ""
	password := ""

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			"host",
			"Missing Punchbag API Host",
			"The provider requires a host for the Punchbag API. Set the host value in the configuration or use the PUNCHBAG_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new Punchbag client using the configuration values
	client := &http.Client{}

	// Example client configuration for API requests
	punchbagClient := &PunchbagClient{
		HostURL:    host,
		Username:   username,
		Password:   password,
		HTTPClient: client,
	}

	// Make the Punchbag client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = punchbagClient
	resp.ResourceData = punchbagClient
}

// DataSources defines the data sources implemented in the provider.
func (p *punchbagProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewClustersDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *punchbagProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewClusterResource,
		NewTestResource,
	}
}

// PunchbagClient represents the API client
type PunchbagClient struct {
	HostURL    string
	Username   string
	Password   string
	HTTPClient *http.Client
}

// Authenticate performs authentication against the API
func (c *PunchbagClient) Authenticate() error {
	// In a real implementation, this would perform authentication
	// and store tokens/credentials for subsequent requests
	return nil
}

// GetCluster retrieves a cluster by ID
func (c *PunchbagClient) GetCluster(id string) (*Cluster, error) {
	// Implementation would make HTTP request to get cluster
	return &Cluster{
		ID:                id,
		Name:              "example-cluster",
		ResourceGroup:     "example-rg",
		Location:          "eastus",
		KubernetesVersion: "1.28.0",
		Status:            "running",
		NodeCount:         3,
	}, nil
}

// CreateCluster creates a new cluster
func (c *PunchbagClient) CreateCluster(cluster *Cluster) (*Cluster, error) {
	// Implementation would make HTTP request to create cluster
	cluster.ID = "generated-id"
	cluster.Status = "creating"
	return cluster, nil
}

// UpdateCluster updates an existing cluster
func (c *PunchbagClient) UpdateCluster(cluster *Cluster) (*Cluster, error) {
	// Implementation would make HTTP request to update cluster
	return cluster, nil
}

// DeleteCluster deletes a cluster
func (c *PunchbagClient) DeleteCluster(id string) error {
	// Implementation would make HTTP request to delete cluster
	return nil
}

// ListClusters retrieves all clusters
func (c *PunchbagClient) ListClusters() ([]*Cluster, error) {
	// Implementation would make HTTP request to list clusters
	return []*Cluster{
		{
			ID:                "cluster-1",
			Name:              "test-cluster-1",
			ResourceGroup:     "test-rg",
			Location:          "eastus",
			KubernetesVersion: "1.28.0",
			Status:            "running",
			NodeCount:         3,
		},
	}, nil
}

// Cluster represents a cluster resource
type Cluster struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	ResourceGroup     string            `json:"resource_group"`
	Location          string            `json:"location"`
	KubernetesVersion string            `json:"kubernetes_version"`
	Status            string            `json:"status"`
	NodeCount         int               `json:"node_count"`
	Tags              map[string]string `json:"tags,omitempty"`
}
