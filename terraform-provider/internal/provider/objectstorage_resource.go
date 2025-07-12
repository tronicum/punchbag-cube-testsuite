package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &objectStorageResource{}
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
	Policy       types.Object `tfsdk:"policy"`
	Versioning   types.Object `tfsdk:"versioning"`
	Lifecycle    types.List   `tfsdk:"lifecycle"`
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
			"policy": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "S3 bucket policy.",
				Attributes: map[string]schema.Attribute{
					"version": schema.StringAttribute{Optional: true},
					"statement": schema.ListNestedAttribute{
						Optional:    true,
						Description: "Policy statements.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"effect":    schema.StringAttribute{Optional: true},
								"principal": schema.MapAttribute{ElementType: types.StringType, Optional: true},
								"action":    schema.ListAttribute{ElementType: types.StringType, Optional: true},
								"resource":  schema.ListAttribute{ElementType: types.StringType, Optional: true},
								"condition": schema.MapAttribute{ElementType: types.StringType, Optional: true},
							},
						},
					},
				},
			},
			"versioning": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "S3 bucket versioning.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{Optional: true},
				},
			},
			"lifecycle": schema.ListNestedAttribute{
				Optional:    true,
				Description: "S3 bucket lifecycle rules.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":     schema.StringAttribute{Optional: true},
						"prefix": schema.StringAttribute{Optional: true},
						"status": schema.StringAttribute{Optional: true},
						"transitions": schema.ListNestedAttribute{
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"days":          schema.Int64Attribute{Optional: true},
									"storage_class": schema.StringAttribute{Optional: true},
								},
							},
						},
						"expiration_days":               schema.Int64Attribute{Optional: true},
						"noncurrent_version_expiration": schema.Int64Attribute{Optional: true},
					},
				},
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
	if r.client != nil {
		// Marshal plan to JSON for POST
		payload, err := json.Marshal(plan)
		if err != nil {
			resp.Diagnostics.AddError("Failed to marshal S3 bucket payload", err.Error())
			return
		}
		url := fmt.Sprintf("%s/api/proxy/%s/s3", r.client.HostURL, plan.Provider.ValueString())
		request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
		if err != nil {
			resp.Diagnostics.AddError("Failed to create HTTP request", err.Error())
			return
		}
		request.Header.Set("Content-Type", "application/json")
		// Add basic auth if needed
		if r.client.Username != "" && r.client.Password != "" {
			request.SetBasicAuth(r.client.Username, r.client.Password)
		}
		response, err := r.client.HTTPClient.Do(request)
		if err != nil {
			resp.Diagnostics.AddError("Failed to call S3 proxy endpoint", err.Error())
			return
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusCreated {
			resp.Diagnostics.AddError("S3 proxy returned error", fmt.Sprintf("status: %d", response.StatusCode))
			return
		}
		var created objectStorageResourceModel
		if err := json.NewDecoder(response.Body).Decode(&created); err != nil {
			resp.Diagnostics.AddError("Failed to decode S3 proxy response", err.Error())
			return
		}
		resp.State.Set(ctx, created)
		return
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
	if r.client != nil {
		url := fmt.Sprintf("%s/api/proxy/%s/s3?id=%s", r.client.HostURL, state.Provider.ValueString(), state.ID.ValueString())
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			resp.Diagnostics.AddError("Failed to create HTTP request", err.Error())
			return
		}
		if r.client.Username != "" && r.client.Password != "" {
			request.SetBasicAuth(r.client.Username, r.client.Password)
		}
		response, err := r.client.HTTPClient.Do(request)
		if err != nil {
			resp.Diagnostics.AddError("Failed to call S3 proxy endpoint", err.Error())
			return
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			resp.Diagnostics.AddError("S3 proxy returned error", fmt.Sprintf("status: %d", response.StatusCode))
			return
		}
		var bucket objectStorageResourceModel
		if err := json.NewDecoder(response.Body).Decode(&bucket); err != nil {
			resp.Diagnostics.AddError("Failed to decode S3 proxy response", err.Error())
			return
		}
		resp.State.Set(ctx, bucket)
		return
	}
	// For now, simulate read:
	resp.State.Set(ctx, state)
}

func (r *objectStorageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan objectStorageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if r.client != nil {
		payload, err := json.Marshal(plan)
		if err != nil {
			resp.Diagnostics.AddError("Failed to marshal S3 bucket payload", err.Error())
			return
		}
		url := fmt.Sprintf("%s/api/proxy/%s/s3", r.client.HostURL, plan.Provider.ValueString())
		request, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(payload))
		if err != nil {
			resp.Diagnostics.AddError("Failed to create HTTP request", err.Error())
			return
		}
		request.Header.Set("Content-Type", "application/json")
		if r.client.Username != "" && r.client.Password != "" {
			request.SetBasicAuth(r.client.Username, r.client.Password)
		}
		response, err := r.client.HTTPClient.Do(request)
		if err != nil {
			resp.Diagnostics.AddError("Failed to call S3 proxy endpoint", err.Error())
			return
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			resp.Diagnostics.AddError("S3 proxy returned error", fmt.Sprintf("status: %d", response.StatusCode))
			return
		}
		var updated objectStorageResourceModel
		if err := json.NewDecoder(response.Body).Decode(&updated); err != nil {
			resp.Diagnostics.AddError("Failed to decode S3 proxy response", err.Error())
			return
		}
		resp.State.Set(ctx, updated)
		return
	}
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
