package provider

import (
	"context"
	"terraform-provider-tykgateway/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ provider.Provider = (*tykgatewayProvider)(nil)

// hashicupsProviderModel maps provider schema data to a Go type.
type tykgatewayProviderModel struct {
	GatewayUrl types.String `tfsdk:"gateway_url"`
	ApiKey     types.String `tfsdk:"api_key"`
}

func New() func() provider.Provider {
	return func() provider.Provider {
		return &tykgatewayProvider{}
	}
}

type tykgatewayProvider struct{}

func (p *tykgatewayProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Tyk Gateway provider for managing Tyk Gateway resources.",
		Attributes: map[string]schema.Attribute{
			"gateway_url": schema.StringAttribute{
				Description: "The URL of the Tyk Gateway instance.",
				Required:    true,
			},
			"api_key": schema.StringAttribute{
				Description: "The API key for authenticating with the Tyk Gateway instance.",
				Sensitive:   true,
				Required:    true,
			},
		},
	}
}

func (p *tykgatewayProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Debug(ctx, "Configuring TykGateway provider")
	// Retrieve provider data from configuration
	var config tykgatewayProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.GatewayUrl.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("gateway_url"),
			"Unknown Tyk Gateway url",
			"The provider cannot create the Tyk Gateway API client as there is an unknown configuration value for the Tyk Gateway url.",
		)
	}

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Tyk Gateway API Key",
			"The provider cannot create the Tyk Gateway API client as there is an unknown configuration value for the Tyk Gateway api_key. ",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	gatewayUrl := config.GatewayUrl.ValueString()
	apiKey := config.ApiKey.ValueString()

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if gatewayUrl == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("gateway_url"),
			"Unknown Tyk Gateway url",
			"The provider cannot create the Tyk Gateway API client as there is an unknown configuration value for the Tyk Gateway url.",
		)
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Tyk Gateway API Key",
			"The provider cannot create the Tyk Gateway API client as there is an unknown configuration value for the Tyk Gateway api_key. ",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "gateway url", gatewayUrl)
	ctx = tflog.SetField(ctx, "api_key", apiKey)
	tflog.MaskFieldValuesWithFieldKeys(ctx, "api_key")

	tflog.Debug(ctx, "Creating TykGateway client")

	client, err := client.NewClient(gatewayUrl, apiKey)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Tyk Gateway client",
			"Could not create Tyk Gateway client: "+err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
	tflog.Debug(ctx, "TykGateway client created successfully", map[string]interface{}{
		"gateway_url": gatewayUrl,
	})
}

func (p *tykgatewayProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "tykgateway"
}

func (p *tykgatewayProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *tykgatewayProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewKeyResource,
	}
}
