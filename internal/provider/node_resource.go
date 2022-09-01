package provider

import (
	"context"
	"fmt"
	"github.com/OnFinality-io/onf-cli/pkg/models"
	onf "github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"log"
	"strconv"
	"strings"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.ResourceType = onFinalityNode{}
var _ resource.Resource = nodeResource{}
var _ resource.ResourceWithImportState = nodeResource{}

type nodeSpec struct {
	Key        types.String `tfsdk:"key"`
	Multiplier types.Int64  `tfsdk:"multiplier"`
}
type onFinalityNode struct {
	WorkspaceId    types.Int64  `tfsdk:"workspace_id"`
	NetworkSpecKey types.String `tfsdk:"network_spec_key"`
	NodeSpec       nodeSpec     `tfsdk:"node_spec"`
	NodeType       types.String `tfsdk:"node_type"`
	NodeName       types.String `tfsdk:"node_name"`
	ClusterHash    types.String `tfsdk:"cluster_hash"`
	Storage        types.String `tfsdk:"storage"`
	ImageVersion   types.String `tfsdk:"image_version"`
	Image          types.String `tfsdk:"image"`
	Id             types.Int64  `tfsdk:"id"`
}

func (t onFinalityNode) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example resource",

		Attributes: map[string]tfsdk.Attribute{
			"workspace_id": {
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
				Type:                types.Int64Type,
			},
			"network_spec_key": {
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
				Type:                types.StringType,
			},
			"node_spec": {
				MarkdownDescription: "",
				Required:            true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"key":        {Type: types.StringType, Required: true},
					"multiplier": {Type: types.Int64Type, Required: true},
				}),
			},
			"node_type": {
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
				Type:                types.StringType,
			},
			"node_name": {
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
				Type:                types.StringType,
			},
			"cluster_hash": {
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
				Type:                types.StringType,
			},
			"storage": {
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
				Type:                types.StringType,
			},
			"image_version": {
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
				Type:                types.StringType,
			},
			"image": {
				MarkdownDescription: "Example configurable attribute",
				Computed:            true,
				Type:                types.StringType,
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "Example identifier",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.Int64Type,
			},
		},
	}, nil
}

func (t onFinalityNode) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	log.Println("Ian-recource.go_NewResource")
	provider, diags := convertProviderType(in)

	return nodeResource{
		provider: provider,
	}, diags
}

type nodeResource struct {
	provider onfinalityProvider
}

func (r nodeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	log.Println("Ian-recource.go_Create")
	var data onFinalityNode
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// example, err := d.provider.client.CreateExample(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	//     return
	// }

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	node, err := onf.CreateNode(uint64(data.WorkspaceId.Value), &onf.CreateNodePayload{
		NetworkSpecKey: data.NetworkSpecKey.Value,
		NodeSpec:       &onf.NodeSpec{Key: data.NodeSpec.Key.Value, Multiplier: int(data.NodeSpec.Multiplier.Value)},
		NodeType:       models.NodeType(data.NodeType.Value),
		NodeName:       data.NodeName.Value,
		ClusterHash:    data.ClusterHash.Value,
		Storage:        &data.Storage.Value,
		InitFromBackup: true,
		UseApiKey:      true,
		ImageVersion:   &data.ImageVersion.Value,
		PublicPort:     true,
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create node, got error: %s", err))
		return
	}
	data.Id = types.Int64{Value: int64(node.ID)}
	data.Image = types.String{Value: node.Image}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	// tflog.Trace(ctx, "created a resource")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r nodeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	log.Println("Ian-recource.go_Read")
	var data onFinalityNode

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// example, err := d.provider.client.ReadExample(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }
	node, err := onf.GetNodeDetail(uint64(data.WorkspaceId.Value), uint64(data.Id.Value))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get node, got error: %s", err))
		return
	}
	if node.Status == "terminated" {
		tflog.Info(ctx, "Node has been terminated")
		resp.State.RemoveResource(ctx)
		return
	}
	imageSlice := strings.Split(node.Image, ":")
	diags = resp.State.Set(ctx, &onFinalityNode{
		WorkspaceId:    types.Int64{Value: int64(node.WorkspaceID)},
		NetworkSpecKey: types.String{Value: node.NetworkSpecKey},
		NodeSpec: nodeSpec{
			Key:        types.String{Value: node.NodeSpec},
			Multiplier: types.Int64{Value: int64(node.NodeSpecMultiplier)},
		},
		NodeType:     types.String{Value: node.NodeType},
		NodeName:     types.String{Value: node.Name},
		ClusterHash:  types.String{Value: node.ClusterHash},
		Storage:      types.String{Value: node.Storage},
		ImageVersion: types.String{Value: imageSlice[len(imageSlice)-1]},
		Image:        types.String{Value: node.Image},
		Id:           types.Int64{Value: int64(node.ID)},
	})
	resp.Diagnostics.Append(diags...)
}

func (r nodeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	log.Println("Ian-recource.go_Update")
	var data onFinalityNode

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// example, err := d.provider.client.UpdateExample(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r nodeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data onFinalityNode

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := onf.TerminateNode(uint64(data.WorkspaceId.Value), uint64(data.Id.Value))
	if err != nil {
		tflog.Error(ctx, "delete node error:"+err.Error())
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// example, err := d.provider.client.DeleteExample(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r nodeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	//resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	idSlice := strings.Split(req.ID, ":")
	wsId, err := strconv.ParseInt(idSlice[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Id Error", fmt.Sprintf("Unable to convert wdId to int64, got error: %s", err))
		return
	}
	id, err := strconv.ParseInt(idSlice[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Id Error", fmt.Sprintf("Unable to convert nodeId to int64, got error: %s", err))
		return
	}

	node, err := onf.GetNodeDetail(uint64(wsId), uint64(id))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get node, got error: %s", err))
		return
	}

	if node.Status == "terminated" {
		tflog.Error(ctx, "Node has been terminated")
		return
	}
	imageSlice := strings.Split(node.Image, ":")
	diags := resp.State.Set(ctx, &onFinalityNode{
		WorkspaceId:    types.Int64{Value: int64(node.WorkspaceID)},
		NetworkSpecKey: types.String{Value: node.NetworkSpecKey},
		NodeSpec: nodeSpec{
			Key:        types.String{Value: node.NodeSpec},
			Multiplier: types.Int64{Value: int64(node.NodeSpecMultiplier)},
		},
		NodeType:     types.String{Value: node.NodeType},
		NodeName:     types.String{Value: node.Name},
		ClusterHash:  types.String{Value: node.ClusterHash},
		Storage:      types.String{Value: node.Storage},
		ImageVersion: types.String{Value: imageSlice[len(imageSlice)-1]},
		Image:        types.String{Value: node.Image},
		Id:           types.Int64{Value: int64(node.ID)},
	})
	resp.Diagnostics.Append(diags...)
}
