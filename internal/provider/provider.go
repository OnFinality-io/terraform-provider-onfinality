package provider

import (
	"context"
	"fmt"
	onf "github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"log"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.Provider = &onfinalityProvider{}

// provider satisfies the tfsdk.Provider interface and usually is included
// with all Resource and DataSource implementations.
type onfinalityProvider struct {
	// client can contain the upstream provider SDK or HTTP client used to
	// communicate with the upstream service. Resource and DataSource
	// implementations can then make calls using this client.
	//
	// TODO: If appropriate, implement upstream provider SDK or HTTP client.
	// client vendorsdk.ExampleClient

	// configured is set to true at the end of the Configure method.
	// This can be used in Resource and DataSource implementations to verify
	// that the provider was previously configured.
	configured bool

	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// providerData can be used to store data from the Terraform configuration.
type providerData struct {
	AccessKey types.String `tfsdk:"access_key"`
	SecretKey types.String `tfsdk:"secret_key"`
}

func (p *onfinalityProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data providerData
	log.Println("Ian-Configure")
	diags := req.Config.Get(ctx, &data)
	onf.Init(data.AccessKey.Value, data.SecretKey.Value, "https://api.onfinality.io/api")
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Example.Null { /* ... */ }

	// If the upstream provider SDK or HTTP client requires configuration, such
	// as authentication or logging, this is a great opportunity to do so.

	p.configured = true
}

func (p *onfinalityProvider) GetResources(ctx context.Context) (map[string]provider.ResourceType, diag.Diagnostics) {
	log.Println("Ian-GetResources")
	return map[string]provider.ResourceType{
		"onfinality_node": onFinalityNode{},
	}, nil
}

func (p *onfinalityProvider) GetDataSources(ctx context.Context) (map[string]provider.DataSourceType, diag.Diagnostics) {
	log.Println("Ian-GetDataSources")
	return map[string]provider.DataSourceType{
		"scaffolding_example": exampleDataSourceType{},
	}, nil
}

func (p *onfinalityProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	log.Println("Ian-GetSchema")
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"access_key": {
				MarkdownDescription: "Example provider attribute",
				Required:            true,
				Type:                types.StringType,
			},
			"secret_key": {
				MarkdownDescription: "Example provider attribute",
				Required:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func New(version string) func() provider.Provider {
	log.Println("Ian-New")
	return func() provider.Provider {
		log.Println("Ian-New Provider")
		return &onfinalityProvider{
			version: version,
		}
	}
}

// convertProviderType is a helper function for NewResource and NewDataSource
// implementations to associate the concrete provider type. Alternatively,
// this helper can be skipped and the provider type can be directly type
// asserted (e.g. provider: in.(*onfinalityProvider)), however using this can prevent
// potential panics.
func convertProviderType(in provider.Provider) (onfinalityProvider, diag.Diagnostics) {
	log.Println("Ian-convertProviderType")
	var diags diag.Diagnostics

	p, ok := in.(*onfinalityProvider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)
		return onfinalityProvider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)
		return onfinalityProvider{}, diags
	}

	return *p, diags
}
