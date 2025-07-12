package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	
	"punchbag-cube-testsuite/shared/providers"
	"punchbag-cube-testsuite/shared/providers/azure"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AzureMonitorResource{}

func NewAzureMonitorResource() resource.Resource {
	return &AzureMonitorResource{}
}

// AzureMonitorResource defines the resource implementation.
type AzureMonitorResource struct {
	provider *azure.Provider
}

// AzureMonitorResourceModel describes the resource data model.
type AzureMonitorResourceModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	ResourceGroup types.String `tfsdk:"resource_group"`
	Location      types.String `tfsdk:"location"`
	WorkspaceName types.String `tfsdk:"workspace_name"`
	Status        types.String `tfsdk:"status"`
}

func (r *AzureMonitorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_azure_monitor"
}

func (r *AzureMonitorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Azure Monitor resource using shared library",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Monitor identifier",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Monitor name",
			},
			"resource_group": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Azure resource group",
			},
			"location": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Azure location",
			},
			"workspace_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Log Analytics workspace name",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Monitor status",
			},
		},
	}
}

func (r *AzureMonitorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*MulticloudProvider)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *MulticloudProvider, got: %T", req.ProviderData),
		)
		return
	}

	r.provider = client.azureProvider
}

func (r *AzureMonitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AzureMonitorResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Use shared library to create monitor
	params := providers.MonitorParams{
		Name:          data.Name.ValueString(),
		ResourceGroup: data.ResourceGroup.ValueString(),
		Location:      data.Location.ValueString(),
		WorkspaceName: data.WorkspaceName.ValueString(),
	}

	result, err := r.provider.CreateMonitor(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create monitor, got error: %s", err))
		return
	}

	// Update data model with result
	data.Id = types.StringValue(result.ID)
	data.Status = types.StringValue(result.Status)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AzureMonitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AzureMonitorResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Use shared library to read monitor status
	monitors, err := r.provider.ListMonitors(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read monitor, got error: %s", err))
		return
	}

	// Find our monitor
	found := false
	for _, monitor := range monitors {
		if monitor.ID == data.Id.ValueString() {
			data.Status = types.StringValue(monitor.Status)
			found = true
			break
		}
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AzureMonitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AzureMonitorResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Implement update logic using shared library

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AzureMonitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AzureMonitorResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Use shared library to delete monitor
	err := r.provider.DeleteMonitor(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete monitor, got error: %s", err))
		return
	}
}
