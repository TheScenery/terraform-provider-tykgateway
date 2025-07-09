package provider

import (
	"context"
	"terraform-provider-tykgateway/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = (*keyResource)(nil)

func NewKeyResource() resource.Resource {
	return &keyResource{}
}

type keyResource struct {
	client *client.Client
}

type keyResourceModel struct {
	Id types.String `tfsdk:"id"`
}

func (r *keyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key"
}

func (r *keyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"hashed": schema.BoolAttribute{
				Description: "Indicates if the key is hashed.",
				Optional:    true,
			},
			"access_rights": schema.MapNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"allow_urls": schema.ListNestedAttribute{
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"methods": schema.ListAttribute{
										ElementType:         types.StringType,
										Optional:            true,
										Description:         "HTTP methods allowed for the URL.",
										MarkdownDescription: "HTTP methods allowed for the URL.",
									},
									"url": schema.StringAttribute{
										Optional:            true,
										Description:         "Allowed URL.",
										MarkdownDescription: "Allowed URL.",
									},
								},
								CustomType: AllowUrlsType{
									ObjectType: types.ObjectType{
										AttrTypes: AllowUrlsValue{}.AttributeTypes(ctx),
									},
								},
							},
							Optional:            true,
							Description:         "List of allowed URLs for the key.",
							MarkdownDescription: "List of allowed URLs for the key.",
						},
						"api_id": schema.StringAttribute{
							Optional:            true,
							Description:         "API ID.",
							MarkdownDescription: "API ID.",
						},
						"api_name": schema.StringAttribute{
							Optional:            true,
							Description:         "API name.",
							MarkdownDescription: "API name.",
						},
						"limit": schema.ObjectAttribute{
							AttributeTypes: map[string]attr.Type{
								"rate":                 types.Float64Type,
								"per":                  types.Float64Type,
								"throttle_interval":    types.Float64Type,
								"throttle_retry_limit": types.Int64Type,
								"max_query_depth":      types.Int64Type,
								"quota_max":            types.Int64Type,
								"quota_renews":         types.Int64Type,
								"quota_remaining":      types.Int64Type,
								"quota_renewal_rate":   types.Int64Type,
							},
							Optional:            true,
							Description:         "Rate limiting configuration for the key.",
							MarkdownDescription: "Rate limiting configuration for the key.",
						},
						"versions": schema.ListAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							Description:         "List of API versions.",
							MarkdownDescription: "List of API versions.",
						},
					},
					CustomType: AccessRightsType{
						ObjectType: types.ObjectType{
							AttrTypes: AccessRightsValue{}.AttributeTypes(ctx),
						},
					},
				},
				Optional:            true,
				Description:         "Access rights for the key, mapping API IDs to access definitions.",
				MarkdownDescription: "Access rights for the key, mapping API IDs to access definitions.",
			},
			"alias": schema.StringAttribute{
				Optional:            true,
				Description:         "Alias for the key.",
				MarkdownDescription: "Alias for the key.",
			},
			"allowance": schema.Float64Attribute{
				Required:            true,
				Description:         "The number of requests allowed for the API key.",
				MarkdownDescription: "The number of requests allowed for the API key.",
			},
			"apply_policies": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Description:         "List of policy IDs to apply to the key.",
				MarkdownDescription: "List of policy IDs to apply to the key.",
			},
			"basic_auth_data": schema.ObjectAttribute{
				AttributeTypes: map[string]attr.Type{
					"password":  types.StringType,
					"hash_type": types.StringType,
				},
				Optional:            true,
				Description:         "Basic authentication data.",
				MarkdownDescription: "Basic authentication data.",
			},
			"certificate": schema.StringAttribute{
				Optional:            true,
				Description:         "Certificate.",
				MarkdownDescription: "Certificate.",
			},
			"data_expires": schema.Int64Attribute{
				Optional:            true,
				Description:         "Data expiration time.",
				MarkdownDescription: "Data expiration time.",
			},
			"date_created": schema.StringAttribute{
				Optional:            true,
				Description:         "The date and time when the API key was created, in Unix timestamp format.",
				MarkdownDescription: "The date and time when the API key was created, in Unix timestamp format.",
			},
			"enable_detailed_recording": schema.BoolAttribute{
				Optional:            true,
				Description:         "Enable detailed recording.",
				MarkdownDescription: "Enable detailed recording.",
			},
			"enable_http_signature_validation": schema.BoolAttribute{
				Optional:            true,
				Description:         "Enable HTTP signature validation.",
				MarkdownDescription: "Enable HTTP signature validation.",
			},
			"expires": schema.Int64Attribute{
				Optional:            true,
				Description:         "The expiration time of the API key, in Unix timestamp format.",
				MarkdownDescription: "The expiration time of the API key, in Unix timestamp format.",
			},
			"hmac_enabled": schema.BoolAttribute{
				Optional:            true,
				Description:         "Whether HMAC is enabled.",
				MarkdownDescription: "Whether HMAC is enabled.",
			},
			"hmac_string": schema.StringAttribute{
				Optional:            true,
				Description:         "HMAC secret string.",
				MarkdownDescription: "HMAC secret string.",
			},
			"id_extractor_deadline": schema.Int64Attribute{
				Optional:            true,
				Description:         "ID extractor deadline.",
				MarkdownDescription: "ID extractor deadline.",
			},
			"is_inactive": schema.BoolAttribute{
				Optional:            true,
				Description:         "Whether the key is inactive.",
				MarkdownDescription: "Whether the key is inactive.",
			},
			"jwt_data": schema.ObjectAttribute{
				AttributeTypes: map[string]attr.Type{
					"secret": types.StringType,
				},
				Optional:            true,
				Description:         "JWT data for the key.",
				MarkdownDescription: "JWT data for the key.",
			},
			"last_check": schema.Int64Attribute{
				Optional:            true,
				Description:         "The last time the API key was checked, in Unix timestamp format.",
				MarkdownDescription: "The last time the API key was checked, in Unix timestamp format.",
			},
			"last_updated": schema.StringAttribute{
				Optional:            true,
				Description:         "Last updated timestamp.",
				MarkdownDescription: "Last updated timestamp.",
			},
			"max_query_depth": schema.Int64Attribute{
				Optional:            true,
				Description:         "The maximum depth of queries allowed for the API key.",
				MarkdownDescription: "The maximum depth of queries allowed for the API key.",
			},
			"meta_data": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Description:         "Custom metadata for the key.",
				MarkdownDescription: "Custom metadata for the key.",
			},
			"monitor": schema.ObjectAttribute{
				AttributeTypes: map[string]attr.Type{
					"trigger_limits": types.ListType{
						ElemType: types.Float64Type,
					},
				},
				Optional:            true,
				Description:         "Monitoring configuration for the key.",
				MarkdownDescription: "Monitoring configuration for the key.",
			},
			"oauth_client_id": schema.StringAttribute{
				Optional:            true,
				Description:         "OAuth client ID.",
				MarkdownDescription: "OAuth client ID.",
			},
			"oauth_keys": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Description:         "OAuth keys.",
				MarkdownDescription: "OAuth keys.",
			},
			"org_id": schema.StringAttribute{
				Optional:            true,
				Description:         "Organization ID.",
				MarkdownDescription: "Organization ID.",
			},
			"per": schema.Float64Attribute{
				Optional:            true,
				Description:         "The time period for the rate limit, in seconds.",
				MarkdownDescription: "The time period for the rate limit, in seconds.",
			},
			"quota_max": schema.Int64Attribute{
				Optional:            true,
				Description:         "The maximum quota for the API key, in requests.",
				MarkdownDescription: "The maximum quota for the API key, in requests.",
			},
			"quota_remaining": schema.Int64Attribute{
				Optional:            true,
				Description:         "The remaining quota for the API key, in requests.",
				MarkdownDescription: "The remaining quota for the API key, in requests.",
			},
			"quota_renewal_rate": schema.Int64Attribute{
				Optional:            true,
				Description:         "The rate at which the quota renews, in requests per second.",
				MarkdownDescription: "The rate at which the quota renews, in requests per second.",
			},
			"quota_renews": schema.Int64Attribute{
				Optional:            true,
				Description:         "The time when the quota renews, in Unix timestamp format.",
				MarkdownDescription: "The time when the quota renews, in Unix timestamp format.",
			},
			"rate": schema.Float64Attribute{
				Optional:            true,
				Description:         "The rate limit for the API key, in requests per second.",
				MarkdownDescription: "The rate limit for the API key, in requests per second.",
			},
			"rsa_certificate_id": schema.StringAttribute{
				Optional:            true,
				Description:         "RSA certificate ID.",
				MarkdownDescription: "RSA certificate ID.",
			},
			"session_lifetime": schema.Int64Attribute{
				Optional:            true,
				Description:         "Session lifetime.",
				MarkdownDescription: "Session lifetime.",
			},
			"smoothing": schema.ObjectAttribute{
				AttributeTypes: map[string]attr.Type{
					"enabled":   types.BoolType,
					"threshold": types.Int64Type,
					"trigger":   types.Float64Type,
					"step":      types.Int64Type,
					"delay":     types.Int64Type,
				},
				Optional:            true,
				Description:         "Smoothing configuration for the key.",
				MarkdownDescription: "Smoothing configuration for the key.",
			},
			"tags": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Description:         "Tags for the key.",
				MarkdownDescription: "Tags for the key.",
			},
			"throttle_interval": schema.Float64Attribute{
				Optional:            true,
				Description:         "The interval for throttling requests, in seconds.",
				MarkdownDescription: "The interval for throttling requests, in seconds.",
			},
			"throttle_retry_limit": schema.Int64Attribute{
				Optional:            true,
				Description:         "The number of retries allowed for throttled requests.",
				MarkdownDescription: "The number of retries allowed for throttled requests.",
			},
		},
	}
}

func (r *keyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data keyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic

	// Example data value setting
	data.Id = types.StringValue("example-id")

	r.client.CreateKey(client.Key{})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *keyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data keyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *keyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data keyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *keyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data keyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
}
