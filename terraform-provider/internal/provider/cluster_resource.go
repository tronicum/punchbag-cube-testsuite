package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource              = &clusterResource{}
	_ resource.ResourceWithImportState = &clusterResource{}
)

func NewClusterResource() resource.Resource {
	return &clusterResource{}
}

// clusterResource defines the resource implementation.
type clusterResource struct {
	client *PunchbagClient
}

// clusterResourceModel describes the resource data model.
type clusterResourceModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	ResourceGroup     types.String `tfsdk:"resource_group"`
	Location          types.String `tfsdk:"location"`
	KubernetesVersion types.String `tfsdk:"kubernetes_version"`
	Status            types.String `tfsdk:"status"`
	NodeCount         types.Int64  `tfsdk:"node_count"`
	Tags              types.Map    `tfsdk:"tags"`
}

func (r *clusterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

func (r *clusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Cluster resource for managing AKS clusters in Punchbag Test Suite.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Cluster identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Cluster name",
				Required:            true,
			},
			"resource_group": schema.StringAttribute{
				MarkdownDescription: "Azure resource group",
				Required:            true,
			},
			"location": schema.StringAttribute{
				MarkdownDescription: "Azure location",
				Required:            true,
			},
			"kubernetes_version": schema.StringAttribute{
				MarkdownDescription: "Kubernetes version",
				Optional:            true,
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Cluster status",
				Computed:            true,
			},
			"node_count": schema.Int64Attribute{
				MarkdownDescription: "Number of nodes",
				Optional:            true,
				Computed:            true,
			},
			"tags": schema.MapAttribute{
				MarkdownDescription: "Tags",
				ElementType:         types.StringType,
				Optional:            true,
			},
		},
	}
}

func (r *clusterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*PunchbagClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *PunchbagClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *clusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data clusterResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create cluster
	cluster := &Cluster{
		Name:              data.Name.ValueString(),
		ResourceGroup:     data.ResourceGroup.ValueString(),
		Location:          data.Location.ValueString(),
		KubernetesVersion: data.KubernetesVersion.ValueString(),
		NodeCount:         int(data.NodeCount.ValueInt64()),
	}

	if data.NodeCount.IsNull() {
		cluster.NodeCount = 3 // default
	}

	if data.KubernetesVersion.IsNull() {
		cluster.KubernetesVersion = "1.28.0" // default
	}

	createdCluster, err := r.client.CreateCluster(cluster)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create cluster, got error: %s", err))
		return
	}

	// Update model with created cluster data
	data.ID = types.StringValue(createdCluster.ID)
	data.Status = types.StringValue(createdCluster.Status)
	data.KubernetesVersion = types.StringValue(createdCluster.KubernetesVersion)
	data.NodeCount = types.Int64Value(int64(createdCluster.NodeCount))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *clusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data clusterResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get cluster from API
	cluster, err := r.client.GetCluster(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read cluster, got error: %s", err))
		return
	}

	// Update model with current cluster data
	data.Name = types.StringValue(cluster.Name)
	data.ResourceGroup = types.StringValue(cluster.ResourceGroup)
	data.Location = types.StringValue(cluster.Location)
	data.KubernetesVersion = types.StringValue(cluster.KubernetesVersion)
	data.Status = types.StringValue(cluster.Status)
	data.NodeCount = types.Int64Value(int64(cluster.NodeCount))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *clusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data clusterResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update cluster
	cluster := &Cluster{
		ID:                data.ID.ValueString(),
		Name:              data.Name.ValueString(),
		ResourceGroup:     data.ResourceGroup.ValueString(),
		Location:          data.Location.ValueString(),
		KubernetesVersion: data.KubernetesVersion.ValueString(),
		NodeCount:         int(data.NodeCount.ValueInt64()),
	}

	updatedCluster, err := r.client.UpdateCluster(cluster)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update cluster, got error: %s", err))
		return
	}

	// Update model with updated cluster data
	data.Status = types.StringValue(updatedCluster.Status)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *clusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data clusterResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete cluster
	err := r.client.DeleteCluster(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete cluster, got error: %s", err))
		return
	}
}

func (r *clusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by ID
	data := clusterResourceModel{
		ID: types.StringValue(req.ID),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
