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
	_ resource.Resource              = &testResource{}
	_ resource.ResourceWithImportState = &testResource{}
)

func NewTestResource() resource.Resource {
	return &testResource{}
}

// testResource defines the resource implementation.
type testResource struct {
	client *PunchbagClient
}

// testResourceModel describes the resource data model.
type testResourceModel struct {
	ID        types.String `tfsdk:"id"`
	ClusterID types.String `tfsdk:"cluster_id"`
	TestType  types.String `tfsdk:"test_type"`
	Status    types.String `tfsdk:"status"`
	Config    types.Map    `tfsdk:"config"`
}

func (r *testResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_test"
}

func (r *testResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Test resource for running tests on AKS clusters in Punchbag Test Suite.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Test identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cluster_id": schema.StringAttribute{
				MarkdownDescription: "Cluster ID to run test on",
				Required:            true,
			},
			"test_type": schema.StringAttribute{
				MarkdownDescription: "Type of test to run (load_test, performance_test, stress_test)",
				Required:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Test status",
				Computed:            true,
			},
			"config": schema.MapAttribute{
				MarkdownDescription: "Test configuration",
				ElementType:         types.StringType,
				Optional:            true,
			},
		},
	}
}

func (r *testResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *testResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data testResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// For this example, we'll simulate test creation
	// In a real implementation, this would call the API to start a test
	testID := "test-" + data.ClusterID.ValueString() + "-" + data.TestType.ValueString()

	// Update model with test data
	data.ID = types.StringValue(testID)
	data.Status = types.StringValue("running")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *testResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data testResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// For this example, we'll simulate reading test status
	// In a real implementation, this would call the API to get test result
	data.Status = types.StringValue("completed")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *testResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data testResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Tests typically aren't updated, but we'll support it for completeness
	// In a real implementation, this might restart the test or update configuration

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *testResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data testResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// For this example, we'll simulate test deletion
	// In a real implementation, this might cancel a running test
	// Tests are typically not deleted, but marked as cancelled
}

func (r *testResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by ID
	data := testResourceModel{
		ID: types.StringValue(req.ID),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
