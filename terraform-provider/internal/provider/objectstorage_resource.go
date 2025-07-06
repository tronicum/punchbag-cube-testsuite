package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource              = &objectStorageResource{}
	_ resource.ResourceWithImportState = &objectStorageResource{}
)

func NewObjectStorageResource() resource.Resource {
	return &objectStorageResource{}
}

// objectStorageResource defines the resource implementation.
type objectStorageResource struct {
	client *PunchbagClient
}

// objectStorageResourceModel describes the resource data model.
type objectStorageResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Region       types.String `tfsdk:"region"`
	Provider     types.String `tfsdk:"provider"`
	StorageClass types.String `tfsdk:"storage_class"`
	Tier         types.String `tfsdk:"tier"`
}

func (r *objectStorageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "multipass_cloud_layer_bucket"
}

func (r *objectStorageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Generic multipass-cloud-layer object storage bucket.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique bucket ID.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Bucket name.",
			},
			"region": schema.StringAttribute{
				Required:    true,
				Description: "Region.",
			},
			"provider": schema.StringAttribute{
				Required:    true,
				Description: "Cloud provider (aws|azure|gcp|ionos|stackit|hcloud).",
			},
			"storage_class": schema.StringAttribute{
				Optional:    true,
				Description: "Storage class (STANDARD, COOL, etc.).",
			},
			"tier": schema.StringAttribute{
				Optional:    true,
				Description: "Tier (for Azure Blob, etc.).",
			},
		},
	}
}

func (r *objectStorageResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*PunchbagClient); ok {
		r.client = client
	}
}

func (r *objectStorageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan objectStorageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...) 
	if resp.Diagnostics.HasError() {
		return
	}
	// --- Real API/proxy integration ---
	if r.client != nil {
		// Example: POST to /api/proxy/{provider}/s3 or similar endpoint
		// (You may want to use /api/proxy/{provider}/objectstorage for generality)
		// TODO: Marshal payload and send HTTP POST to url
		// url := r.client.HostURL + "/api/proxy/" + plan.Provider.ValueString() + "/s3"
		// payload := map[string]string{
		//   "name": plan.Name.ValueString(),
		//   "region": plan.Region.ValueString(),
		//   "storage_class": plan.StorageClass.ValueString(),
		//   "tier": plan.Tier.ValueString(),
		// }
	}
	// For now, simulate creation:
	switch plan.Provider.ValueString() {
	case "aws", "azure", "gcp", "ionos", "stackit", "hcloud":
		plan.ID = types.StringValue(plan.Name.ValueString() + "-" + plan.Provider.ValueString())
	default:
		resp.Diagnostics.AddError("Unsupported provider", "Provider '"+plan.Provider.ValueString()+"' is not supported.")
		return
	}
	resp.State.Set(ctx, plan)
}

func (r *objectStorageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state objectStorageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...) 
	if resp.Diagnostics.HasError() {
		return
	}
	// --- Real API/proxy integration scaffold ---
	// Here you would:
	// 1. Select the correct API endpoint or SDK based on state.Provider
	// 2. Authenticate
	// 3. Make the get/read bucket call
	// 4. Update state fields from the real bucket
	// For now, simulate read:
	resp.State.Set(ctx, state)
}

func (r *objectStorageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan objectStorageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...) 
	if resp.Diagnostics.HasError() {
		return
	}
	// --- Real API/proxy integration scaffold ---
	// Here you would:
	// 1. Select the correct API endpoint or SDK based on plan.Provider
	// 2. Authenticate
	// 3. Make the update bucket call
	// 4. Update plan.ID if needed
	// For now, simulate update:
	plan.ID = types.StringValue(plan.Name.ValueString() + "-" + plan.Provider.ValueString())
	resp.State.Set(ctx, plan)
}

func (r *objectStorageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state objectStorageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...) 
	if resp.Diagnostics.HasError() {
		return
	}
	// --- Real API/proxy integration scaffold ---
	// Here you would:
	// 1. Select the correct API endpoint or SDK based on state.Provider
	// 2. Authenticate
	// 3. Make the delete bucket call
	// For now, simulate delete:
	resp.State.RemoveResource(ctx)
}

func (r *objectStorageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
