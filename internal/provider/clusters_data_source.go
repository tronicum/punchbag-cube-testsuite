package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &clustersDataSource{}

func NewClustersDataSource() datasource.DataSource {
	return &clustersDataSource{}
}

// clustersDataSource defines the data source implementation.
type clustersDataSource struct {
	client *PunchbagClient
}

// clustersDataSourceModel describes the data source data model.
type clustersDataSourceModel struct {
	Clusters []clusterDataSourceModel `tfsdk:"clusters"`
	ID       types.String             `tfsdk:"id"`
}

type clusterDataSourceModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	ResourceGroup     types.String `tfsdk:"resource_group"`
	Location          types.String `tfsdk:"location"`
	KubernetesVersion types.String `tfsdk:"kubernetes_version"`
	Status            types.String `tfsdk:"status"`
	NodeCount         types.Int64  `tfsdk:"node_count"`
}

func (d *clustersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_clusters"
}

func (d *clustersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Clusters data source for retrieving all clusters.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Data source identifier",
				Computed:            true,
			},
			"clusters": schema.ListNestedAttribute{
				MarkdownDescription: "List of clusters",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Cluster identifier",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Cluster name",
							Computed:            true,
						},
						"resource_group": schema.StringAttribute{
							MarkdownDescription: "Azure resource group",
							Computed:            true,
						},
						"location": schema.StringAttribute{
							MarkdownDescription: "Azure location",
							Computed:            true,
						},
						"kubernetes_version": schema.StringAttribute{
							MarkdownDescription: "Kubernetes version",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							MarkdownDescription: "Cluster status",
							Computed:            true,
						},
						"node_count": schema.Int64Attribute{
							MarkdownDescription: "Number of nodes",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *clustersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*PunchbagClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *PunchbagClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *clustersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data clustersDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get clusters from API
	clusters, err := d.client.ListClusters()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read clusters, got error: %s", err))
		return
	}

	// Map clusters to model
	for _, cluster := range clusters {
		clusterModel := clusterDataSourceModel{
			ID:                types.StringValue(cluster.ID),
			Name:              types.StringValue(cluster.Name),
			ResourceGroup:     types.StringValue(cluster.ResourceGroup),
			Location:          types.StringValue(cluster.Location),
			KubernetesVersion: types.StringValue(cluster.KubernetesVersion),
			Status:            types.StringValue(cluster.Status),
			NodeCount:         types.Int64Value(int64(cluster.NodeCount)),
		}

		data.Clusters = append(data.Clusters, clusterModel)
	}

	// Set ID for the data source
	data.ID = types.StringValue("clusters")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
