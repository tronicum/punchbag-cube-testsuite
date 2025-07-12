package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	
	"punchbag-cube-testsuite/shared/providers"
	"punchbag-cube-testsuite/shared/providers/azure"
)

// Ensure MulticloudProvider satisfies various provider interfaces.
var _ provider.Provider = &MulticloudProvider{}

// MulticloudProvider defines the provider implementation.
type MulticloudProvider struct {
	version string
	
	// Cloud providers
	azureProvider *azure.Provider
}

// MulticloudProviderModel describes the provider data model.
type MulticloudProviderModel struct {
	CubeServerURL types.String `tfsdk:"cube_server_url"`
	SimulationMode types.Bool   `tfsdk:"simulation_mode"`
	AzureConfig   types.Object `tfsdk:"azure_config"`
	AWSConfig     types.Object `tfsdk:"aws_config"`
	GCPConfig     types.Object `tfsdk:"gcp_config"`
}

func (p *MulticloudProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "multicloud"
	resp.Version = p.version
}

func (p *MulticloudProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Multicloud provider using shared library for simulation and cross-cloud operations",
		Attributes: map[string]schema.Attribute{
			"cube_server_url": schema.StringAttribute{
				MarkdownDescription: "URL of the cube-server for simulation mode",
				Optional:            true,
			},
			"simulation_mode": schema.BoolAttribute{
				MarkdownDescription: "Enable simulation mode via cube-server",
				Optional:            true,
			},
			"azure_config": schema.SingleNestedAttribute{
				MarkdownDescription: "Azure configuration",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"subscription_id": schema.StringAttribute{
						MarkdownDescription: "Azure subscription ID",
						Optional:            true,
					},
					"tenant_id": schema.StringAttribute{
						MarkdownDescription: "Azure tenant ID",
						Optional:            true,
					},
				},
			},
			"aws_config": schema.SingleNestedAttribute{
				MarkdownDescription: "AWS configuration",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"region": schema.StringAttribute{
						MarkdownDescription: "AWS region",
						Optional:            true,
					},
				},
			},
			"gcp_config": schema.SingleNestedAttribute{
				MarkdownDescription: "GCP configuration", 
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"project": schema.StringAttribute{
						MarkdownDescription: "GCP project ID",
						Optional:            true,
					},
				},
			},
		},
	}
}

func (p *MulticloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data MulticloudProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize Azure provider using shared library
	p.azureProvider = azure.NewProvider()
	
	// Configure simulation mode
	if !data.SimulationMode.IsNull() && data.SimulationMode.ValueBool() {
		p.azureProvider.SetSimulationMode(true)
	}
	
	// Configure Azure provider
	azureConfig := make(map[string]interface{})
	// TODO: Extract Azure config from data.AzureConfig
	p.azureProvider.Configure(ctx, azureConfig)
}

func (p *MulticloudProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAzureMonitorResource,
		NewAzureAKSClusterResource,
		NewAzureBudgetResource,
	}
}

func (p *MulticloudProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAzureMonitorDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MulticloudProvider{
			version: version,
		}
	}
}
