package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func NodeImageModifier() tfsdk.AttributePlanModifier {
	return nodeImageModifier{}
}

// nodeImageModifier is an AttributePlanModifier that sets node's image
type nodeImageModifier struct{}

// Modify fills the AttributePlanModifier interface. It sets RequiresReplace on
// the response to true if the following criteria are met:
//
// 1. The resource's state is not null; a null state indicates that we're
// creating a resource, and we never need to destroy and recreate a resource
// when we're creating it.
//
// 2. The resource's plan is not null; a null plan indicates that we're
// deleting a resource, and we never need to destroy and recreate a resource
// when we're deleting it.
//
// 3. The attribute's config is not null or the attribute is not computed; a
// computed attribute with a null config almost always means that the provider
// is changing the value, and practitioners are usually unpleasantly surprised
// when a resource is destroyed and recreated when their configuration hasn't
// changed. This has the unfortunate side effect that removing a computed field
// from the config will not trigger a destroy and recreate cycle, even when
// that is warranted. To get around this, provider developer can implement
// their own AttributePlanModifier that handles that behavior in the way that
// most makes sense for their use case.
//
// 4. The attribute's value in the plan does not match the attribute's value in
// the state.
func (r nodeImageModifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	if req.AttributeConfig == nil || req.AttributePlan == nil || req.AttributeState == nil {
		// shouldn't happen, but let's not panic if it does
		return
	}
	if req.State.Raw.IsNull() {
		// if we're creating the resource, no need to delete and
		// recreate it
		return
	}

	if req.Plan.Raw.IsNull() {
		// if we're deleting the resource, no need to delete and
		// recreate it
		return
	}
	var config onFinalityNode
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	var state onFinalityNode
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if state.Image.IsNull() {
		return
	}
	imageSlice := strings.Split(state.Image.Value, ":")
	resp.AttributePlan = types.String{Value: imageSlice[0] + ":" + config.ImageVersion.Value}

	resp.RequiresReplace = false
	if req.AttributePlan.Equal(req.AttributeState) {
		// if the plan and the state are in agreement, this attribute
		// isn't changing, don't require replace
		return
	}
}

// Description returns a human-readable description of the plan modifier.
func (r nodeImageModifier) Description(ctx context.Context) string {
	return "If the value of this attribute changes, Terraform will destroy and recreate the resource."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (r nodeImageModifier) MarkdownDescription(ctx context.Context) string {
	return "If the value of this attribute changes, Terraform will destroy and recreate the resource."
}
