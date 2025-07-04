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
	Provider          types.String `tfsdk:"provider"`
	ResourceGroup     types.String `tfsdk:"resource_group"`
	Location          types.String `tfsdk:"location"`
	Region            types.String `tfsdk:"region"`
	ProjectID         types.String `tfsdk:"project_id"`
	KubernetesVersion types.String `tfsdk:"kubernetes_version"`
	Status            types.String `tfsdk:"status"`
	NodeCount         types.Int64  `tfsdk:"node_count"`
	Tags              types.Map    `tfsdk:"tags"`
	
	// StackIT-specific fields
	MaintenanceTimeStart      types.String `tfsdk:"maintenance_time_start"`
	MaintenanceTimeEnd        types.String `tfsdk:"maintenance_time_end"`
	MaintenanceTimeZone       types.String `tfsdk:"maintenance_time_zone"`
	AllowPrivilegedContainers types.Bool   `tfsdk:"allow_privileged_containers"`
	
	// Hetzner Cloud-specific fields
	HetznerToken            types.String `tfsdk:"hetzner_token"`
	NetworkZone             types.String `tfsdk:"network_zone"`
	ServerType              types.String `tfsdk:"server_type"`
	SSHKeys                 types.List   `tfsdk:"ssh_keys"`
	EnablePublicNetwork     types.Bool   `tfsdk:"enable_public_network"`
	EnablePrivateNetwork    types.Bool   `tfsdk:"enable_private_network"`
	
	// IONOS Cloud-specific fields
	DatacenterID       types.String `tfsdk:"datacenter_id"`
	IONOSUsername      types.String `tfsdk:"ionos_username"`
	IONOSPassword      types.String `tfsdk:"ionos_password"`
	IONOSToken         types.String `tfsdk:"ionos_token"`
	K8sClusterName     types.String `tfsdk:"k8s_cluster_name"`
	MaintenanceDay     types.String `tfsdk:"maintenance_day"`
	MaintenanceTime    types.String `tfsdk:"maintenance_time"`
	IONOSPublic        types.Bool   `tfsdk:"ionos_public"`
	GatewayIP          types.String `tfsdk:"gateway_ip"`
}

func (r *clusterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

func (r *clusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Cluster resource for managing multi-cloud Kubernetes clusters in Punchbag Test Suite.",

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
			"provider": schema.StringAttribute{
				MarkdownDescription: "Cloud provider (azure, schwarz-stackit, hetzner-hcloud, united-ionos, aws, gcp)",
				Required:            true,
			},
			"resource_group": schema.StringAttribute{
				MarkdownDescription: "Azure resource group (Azure only)",
				Optional:            true,
			},
			"location": schema.StringAttribute{
				MarkdownDescription: "Azure location (Azure only)",
				Optional:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "Cloud region (StackIT, AWS, GCP)",
				Optional:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "Project ID (StackIT, GCP)",
				Optional:            true,
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
			// StackIT-specific attributes
			"maintenance_time_start": schema.StringAttribute{
				MarkdownDescription: "Maintenance window start time (StackIT only, format: HH:MM)",
				Optional:            true,
			},
			"maintenance_time_end": schema.StringAttribute{
				MarkdownDescription: "Maintenance window end time (StackIT only, format: HH:MM)",
				Optional:            true,
			},
			"maintenance_time_zone": schema.StringAttribute{
				MarkdownDescription: "Maintenance window timezone (StackIT only)",
				Optional:            true,
			},
			"allow_privileged_containers": schema.BoolAttribute{
				MarkdownDescription: "Allow privileged containers (StackIT only)",
				Optional:            true,
			},
			// Hetzner Cloud-specific attributes
			"hetzner_token": schema.StringAttribute{
				MarkdownDescription: "Hetzner Cloud API token (Hetzner only)",
				Optional:            true,
				Sensitive:           true,
			},
			"network_zone": schema.StringAttribute{
				MarkdownDescription: "Network zone (Hetzner only)",
				Optional:            true,
			},
			"server_type": schema.StringAttribute{
				MarkdownDescription: "Server type (Hetzner only)",
				Optional:            true,
			},
			"ssh_keys": schema.ListAttribute{
				MarkdownDescription: "SSH key names (Hetzner only)",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"enable_public_network": schema.BoolAttribute{
				MarkdownDescription: "Enable public network (Hetzner only)",
				Optional:            true,
			},
			"enable_private_network": schema.BoolAttribute{
				MarkdownDescription: "Enable private network (Hetzner only)",
				Optional:            true,
			},
			// IONOS Cloud-specific attributes
			"datacenter_id": schema.StringAttribute{
				MarkdownDescription: "Datacenter ID (IONOS only)",
				Optional:            true,
			},
			"ionos_username": schema.StringAttribute{
				MarkdownDescription: "IONOS username (IONOS only)",
				Optional:            true,
			},
			"ionos_password": schema.StringAttribute{
				MarkdownDescription: "IONOS password (IONOS only)",
				Optional:            true,
				Sensitive:           true,
			},
			"ionos_token": schema.StringAttribute{
				MarkdownDescription: "IONOS API token (IONOS only)",
				Optional:            true,
				Sensitive:           true,
			},
			"k8s_cluster_name": schema.StringAttribute{
				MarkdownDescription: "Kubernetes cluster name (IONOS only)",
				Optional:            true,
			},
			"maintenance_day": schema.StringAttribute{
				MarkdownDescription: "Maintenance window day of the week (IONOS only)",
				Optional:            true,
			},
			"maintenance_time": schema.StringAttribute{
				MarkdownDescription: "Maintenance window time (IONOS only)",
				Optional:            true,
			},
			"ionos_public": schema.BoolAttribute{
				MarkdownDescription: "Enable public access (IONOS only)",
				Optional:            true,
			},
			"gateway_ip": schema.StringAttribute{
				MarkdownDescription: "Gateway IP address (IONOS only)",
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

	// Create cluster based on provider
	cluster := &Cluster{
		Name:              data.Name.ValueString(),
		Provider:          data.Provider.ValueString(),
		KubernetesVersion: data.KubernetesVersion.ValueString(),
		NodeCount:         int(data.NodeCount.ValueInt64()),
	}

	// Set provider-specific fields
	switch data.Provider.ValueString() {
	case "azure":
		cluster.ResourceGroup = data.ResourceGroup.ValueString()
		cluster.Location = data.Location.ValueString()
		if cluster.ResourceGroup == "" || cluster.Location == "" {
			resp.Diagnostics.AddError("Validation Error", "resource_group and location are required for Azure clusters")
			return
		}
	case "schwarz-stackit":
		cluster.ProjectID = data.ProjectID.ValueString()
		cluster.Region = data.Region.ValueString()
		if cluster.ProjectID == "" || cluster.Region == "" {
			resp.Diagnostics.AddError("Validation Error", "project_id and region are required for StackIT clusters")
			return
		}
		
		// Handle StackIT-specific configuration
		stackitConfig := &StackITConfig{
			ProjectID:                 cluster.ProjectID,
			MaintenanceTimeStart:      data.MaintenanceTimeStart.ValueString(),
			MaintenanceTimeEnd:        data.MaintenanceTimeEnd.ValueString(),
			MaintenanceTimeZone:       data.MaintenanceTimeZone.ValueString(),
			AllowPrivilegedContainers: data.AllowPrivilegedContainers.ValueBool(),
		}
		
		createdCluster, err := r.client.CreateStackITCluster(cluster, stackitConfig)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create StackIT cluster, got error: %s", err))
			return
		}
		cluster = createdCluster
	case "hetzner-hcloud":
		cluster.Location = data.Location.ValueString()
		if cluster.Location == "" {
			resp.Diagnostics.AddError("Validation Error", "location is required for Hetzner Cloud clusters")
			return
		}
		
		// Handle Hetzner-specific configuration
		hetznerConfig := &HetznerConfig{
			Token:                data.HetznerToken.ValueString(),
			NetworkZone:          data.NetworkZone.ValueString(),
			ServerType:           data.ServerType.ValueString(),
			Location:             cluster.Location,
			EnablePublicNetwork:  data.EnablePublicNetwork.ValueBool(),
			EnablePrivateNetwork: data.EnablePrivateNetwork.ValueBool(),
		}
		
		// Convert SSH keys list
		if !data.SSHKeys.IsNull() {
			var sshKeys []string
			for _, key := range data.SSHKeys.Elements() {
				if keyStr, ok := key.(types.String); ok {
					sshKeys = append(sshKeys, keyStr.ValueString())
				}
			}
			hetznerConfig.SSHKeys = sshKeys
		}
		
		createdCluster, err := r.client.CreateHetznerCluster(cluster, hetznerConfig)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Hetzner cluster, got error: %s", err))
			return
		}
		cluster = createdCluster
	case "united-ionos":
		datacenterID := data.DatacenterID.ValueString()
		if datacenterID == "" {
			resp.Diagnostics.AddError("Validation Error", "datacenter_id is required for IONOS Cloud clusters")
			return
		}
		
		// Handle IONOS-specific configuration
		ionosConfig := &IONOSConfig{
			DatacenterID:   datacenterID,
			Username:       data.IONOSUsername.ValueString(),
			Password:       data.IONOSPassword.ValueString(),
			Token:          data.IONOSToken.ValueString(),
			K8sClusterName: data.K8sClusterName.ValueString(),
			Public:         data.IONOSPublic.ValueBool(),
			GatewayIP:      data.GatewayIP.ValueString(),
		}
		
		// Set maintenance window if provided
		if !data.MaintenanceDay.IsNull() || !data.MaintenanceTime.IsNull() {
			ionosConfig.MaintenanceWindow.DayOfTheWeek = data.MaintenanceDay.ValueString()
			ionosConfig.MaintenanceWindow.Time = data.MaintenanceTime.ValueString()
		}
		
		createdCluster, err := r.client.CreateIONOSCluster(cluster, ionosConfig)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create IONOS cluster, got error: %s", err))
			return
		}
		cluster = createdCluster
	case "aws":
		cluster.Region = data.Region.ValueString()
		if cluster.Region == "" {
			resp.Diagnostics.AddError("Validation Error", "region is required for AWS clusters")
			return
		}
	case "gcp":
		cluster.ProjectID = data.ProjectID.ValueString()
		cluster.Region = data.Region.ValueString()
		if cluster.ProjectID == "" || cluster.Region == "" {
			resp.Diagnostics.AddError("Validation Error", "project_id and region are required for GCP clusters")
			return
		}
	default:
		resp.Diagnostics.AddError("Validation Error", "Unsupported provider. Supported providers: azure, schwarz-stackit, hetzner-hcloud, united-ionos, aws, gcp")
		return
	}

	// Set defaults
	if data.NodeCount.IsNull() {
		cluster.NodeCount = 3 // default
	}
	if data.KubernetesVersion.IsNull() {
		cluster.KubernetesVersion = "1.28.0" // default
	}

	// Create cluster if not StackIT (already created above)
	if data.Provider.ValueString() != "schwarz-stackit" {
		createdCluster, err := r.client.CreateCluster(cluster)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create cluster, got error: %s", err))
			return
		}
		cluster = createdCluster
	}

	// Update model with created cluster data
	data.ID = types.StringValue(cluster.ID)
	data.Status = types.StringValue(cluster.Status)
	data.KubernetesVersion = types.StringValue(cluster.KubernetesVersion)
	data.NodeCount = types.Int64Value(int64(cluster.NodeCount))

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
