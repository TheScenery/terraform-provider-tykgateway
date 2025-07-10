package provider

import (
	"context"
	"terraform-provider-tykgateway/client"

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

var GraphqlType = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"fields": schema.ListAttribute{
			Description: "List of fields that are allowed for the key.",
			Optional:    true,
			ElementType: types.StringType,
		},
		"name": schema.StringAttribute{
			Description: "Name of the allowed type.",
			Required:    true,
		},
	},
}

var AccessSpec = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"url": schema.StringAttribute{
			Required:    true,
			Description: "URL that is allowed for the key.",
		},
		"methods": schema.ListAttribute{
			Description: "List of HTTP methods allowed for the URL.",
			Optional:    true,
		},
	},
}

var RateLimitSmoothing = schema.MapNestedAttribute{
	Description: "Smoothing configuration for the method.",
	Required:    true,
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"delay": schema.Int64Attribute{
				Description: "The delay for ratelimit smoothing",
				Required:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "The enabled for ratelimit smoothing",
				Required:    true,
			},
			"step": schema.Int64Attribute{
				Description: "The step for ratelimit smoothing",
				Required:    true,
			},
			"threshold": schema.Int64Attribute{
				Description: "The threshold for ratelimit smoothing",
				Required:    true,
			},
			"trigger": schema.Float64Attribute{
				Description: "The trigger for ratelimit smoothing",
				Required:    true,
			},
		},
	},
}

var RateLimitType2 = schema.MapNestedAttribute{
	Description: "Rate limit for the HTTP method.",
	Required:    true,
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"per": schema.Int64Attribute{
				Description: "Time period for the rate limit, in seconds.",
				Required:    true,
			},
			"rate": schema.Int64Attribute{
				Description: "Rate limit for the method, in requests per second.",
				Required:    true,
			},
			"smoothing": RateLimitSmoothing,
		},
	},
}

var EndpointMethod = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description: "Name of the HTTP method.",
			Required:    true,
		},
		"limit": RateLimitType2,
	},
}

var EndpointMethods = schema.ListNestedAttribute{
	Description:  "HTTP methods allowed for the endpoint.",
	Required:     true,
	NestedObject: EndpointMethod,
}

var Endpoint = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"path": schema.StringAttribute{
			Description: "Path of the endpoint that the key has access to.",
			Required:    true,
		},
		"methods": EndpointMethods,
	},
}

var Endpoints = schema.ListNestedAttribute{
	Description:  "List of endpoints that the key has access to.",
	Required:     true,
	NestedObject: Endpoint,
}

var FieldLimits = schema.MapNestedAttribute{
	Description: "Limits for the field access rights.",
	Required:    true,
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"max_query_depth": schema.Int64Attribute{
				Description: "Maximum depth of queries allowed for the field.",
				Required:    true,
			},
		},
	},
}

var FieldAccessDefinition = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"field_name": schema.StringAttribute{
			Description: "Name of the field.",
			Required:    true,
		},
		"limits": FieldLimits,
	},
}

var APILimit = schema.MapNestedAttribute{
	Description: "Rate limits for the key.",
	Required:    true,
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"max_query_depth": schema.Int64Attribute{
				Description: "Maximum depth of queries allowed for the key.",
				Required:    true,
			},
			"rate": schema.Float64Attribute{
				Description: "Rate limit for the key, in requests per second.",
				Required:    true,
			},
			"per": schema.Float64Attribute{
				Description: "Time period for the rate limit, in seconds.",
				Required:    true,
			},
			"quota_max": schema.Int64Attribute{
				Description: "Maximum quota for the key, in requests.",
				Required:    true,
			},
			"quota_remaining": schema.Int64Attribute{
				Description: "Remaining quota for the key, in requests.",
				Required:    true,
			},
			"quota_renewal_rate": schema.Int64Attribute{
				Description: "Rate at which the quota renews, in requests per second.",
				Required:    true,
			},
			"quota_renews": schema.Int64Attribute{
				Description: "Time when the quota renews, in Unix timestamp format.",
				Required:    true,
			},
			"throttle_interval": schema.Float64Attribute{
				Description: "Interval for throttling requests, in seconds.",
				Required:    true,
			},
			"throttle_retry_limit": schema.Int64Attribute{
				Description: "Number of retries allowed for throttled requests.",
				Required:    true,
			},
			"smoothing": RateLimitSmoothing,
		},
	},
}

var AccessDefinition = schema.MapNestedAttribute{
	Description: "Access rights for the key.",
	Optional:    true,
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"allowance_scope": schema.StringAttribute{
				Description: "Scope of the allowance for the key.",
				Required:    true,
			},
			"allowed_types": schema.ListNestedAttribute{
				Description:  "List of allowed types for the key.",
				Optional:     true,
				NestedObject: GraphqlType,
			},
			"allow_urls": schema.ListNestedAttribute{
				Description:  "List of allowed URLs for the key.",
				Optional:     true,
				NestedObject: AccessSpec,
			},
			"api_id": schema.StringAttribute{
				Description: "API ID that the key has access to.",
				Required:    true,
			},
			"api_name": schema.StringAttribute{
				Description: "Name of the API that the key has access to.",
				Required:    true,
			},
			"disable_introspection": schema.BoolAttribute{
				Description: "Whether introspection is disabled for the key.",
				Required:    true,
			},
			"endpoints": Endpoints,
			"field_access_rights": schema.ListNestedAttribute{
				Description:  "The Field access righes",
				Optional:     true,
				NestedObject: FieldAccessDefinition,
			},
			"limit": APILimit,
			"restricted_types": schema.ListNestedAttribute{
				Description:  "List of restricted types for the key.",
				Optional:     true,
				NestedObject: GraphqlType,
			},
			"versions": schema.ListAttribute{
				Description: "List of API versions that the key has access to.",
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	},
}

var BasicAuthData = schema.MapNestedAttribute{
	Description: "Basic authentication data for the key.",
	Required:    true,
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"hash_type": schema.StringAttribute{
				Description: "Type of hash used for the basic authentication data.",
				Required:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password for the basic authentication data.",
				Required:    true,
			},
		},
	},
}

var JWTData = schema.MapNestedAttribute{
	Description: "JWT data for the key.",
	Required:    true,
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"secret": schema.StringAttribute{
				Description: "Secret used for signing the JWT.",
				Required:    true,
			},
		},
	},
}

var Monitor = schema.MapNestedAttribute{
	Description: "Monitoring configuration for the key.",
	Required:    true,
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"trigger_limits": schema.ListAttribute{
				Description: "List of trigger limits for monitoring.",
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	},
}

func (r *keyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"hashed": schema.BoolAttribute{
				Description: "Indicates if the key is hashed.",
				Optional:    true,
			},
			"access_rights": AccessDefinition,
			"alias": schema.StringAttribute{
				Required:            true,
				Description:         "Alias for the key.",
				MarkdownDescription: "Alias for the key.",
			},
			"allowance": schema.Float64Attribute{
				Required:    true,
				Description: "The number of requests allowed for the API key.",
			},
			"apply_policies": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of policy IDs to apply to the key.",
			},
			"basic_auth_data": BasicAuthData,
			"certificate": schema.StringAttribute{
				Required:    true,
				Description: "Certificate.",
			},
			"data_expires": schema.Int64Attribute{
				Required:    true,
				Description: "Data expiration time.",
			},
			"date_created": schema.StringAttribute{
				Required:    true,
				Description: "The date and time when the API key was created, in Unix timestamp format.",
			},
			"enable_detailed_recording": schema.BoolAttribute{
				Required:    true,
				Description: "Enable detailed recording.",
			},
			"enable_http_signature_validation": schema.BoolAttribute{
				Required:    true,
				Description: "Enable HTTP signature validation.",
			},
			"expires": schema.Int64Attribute{
				Required:    true,
				Description: "The expiration time of the API key, in Unix timestamp format.",
			},
			"hmac_enabled": schema.BoolAttribute{
				Required:    true,
				Description: "Whether HMAC is enabled.",
			},
			"hmac_string": schema.StringAttribute{
				Required:    true,
				Description: "HMAC secret string.",
			},
			"id_extractor_deadline": schema.Int64Attribute{
				Required:    true,
				Description: "ID extractor deadline.",
			},
			"is_inactive": schema.BoolAttribute{
				Required:    true,
				Description: "Whether the key is inactive.",
			},
			"jwt_data": JWTData,
			"last_check": schema.Int64Attribute{
				Required:    true,
				Description: "The last time the API key was checked, in Unix timestamp format.",
			},
			"last_updated": schema.StringAttribute{
				Required:    true,
				Description: "Last updated timestamp.",
			},
			"max_query_depth": schema.Int64Attribute{
				Required:    true,
				Description: "The maximum depth of queries allowed for the API key.",
			},
			"meta_data": schema.MapAttribute{
				ElementType: types.DynamicType,
				Optional:    true,
				Description: "Custom metadata for the key.",
			},
			"monitor": Monitor,
			"oauth_client_id": schema.StringAttribute{
				Required:    true,
				Description: "OAuth client ID.",
			},
			"oauth_keys": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "OAuth keys.",
			},
			"org_id": schema.StringAttribute{
				Required:    true,
				Description: "Organization ID.",
			},
			"per": schema.Float64Attribute{
				Required:    true,
				Description: "The time period for the rate limit, in seconds.",
			},
			"quota_max": schema.Int64Attribute{
				Required:    true,
				Description: "The maximum quota for the API key, in requests.",
			},
			"quota_remaining": schema.Int64Attribute{
				Required:    true,
				Description: "The remaining quota for the API key, in requests.",
			},
			"quota_renewal_rate": schema.Int64Attribute{
				Required:    true,
				Description: "The rate at which the quota renews, in requests per second.",
			},
			"quota_renews": schema.Int64Attribute{
				Required:    true,
				Description: "The time when the quota renews, in Unix timestamp format.",
			},
			"rate": schema.Float64Attribute{
				Required:    true,
				Description: "The rate limit for the API key, in requests per second.",
			},
			"rsa_certificate_id": schema.StringAttribute{
				Required:    true,
				Description: "RSA certificate ID.",
			},
			"session_lifetime": schema.Int64Attribute{
				Required:    true,
				Description: "Session lifetime.",
			},
			"smoothing": RateLimitSmoothing,
			"tags": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Tags for the key.",
			},
			"throttle_interval": schema.Float64Attribute{
				Required:    true,
				Description: "The interval for throttling requests, in seconds.",
			},
			"throttle_retry_limit": schema.Int64Attribute{
				Required:    true,
				Description: "The number of retries allowed for throttled requests.",
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
