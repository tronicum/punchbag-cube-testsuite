package provider

import (
	"context"

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
type objectStorageResource struct{}

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
				Description: "Cloud provider (aws|azure|gcp).",
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

// --- CRUD Implementation ---

func (r *objectStorageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan objectStorageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Simulate ID assignment (in real use, call cloud API or backend)
	plan.ID = types.StringValue(plan.Provider.ValueString() + ":" + plan.Name.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *objectStorageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state objectStorageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Simulate always present (in real use, fetch from cloud API/backend)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *objectStorageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan objectStorageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Simulate update (in real use, call cloud API/backend)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *objectStorageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Simulate delete (in real use, call cloud API/backend)
	resp.State.RemoveResource(ctx)
}

func (r *objectStorageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Support import by ID
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
