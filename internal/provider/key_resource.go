package provider

import (
	"context"
	"encoding/json"
	"terraform-provider-tykgateway/client"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ resource.Resource = &keyResource{}
var _ resource.ResourceWithConfigure = &keyResource{}

func NewKeyResource() resource.Resource {
	return &keyResource{}
}

type keyResource struct {
	client *client.Client
}

type keyResourceModel struct {
	Key                           types.String  `tfsdk:"key"`
	KeyHash                       types.String  `tfsdk:"key_hash"`
	Hashed                        types.Bool    `tfsdk:"hashed"`
	AccessRights                  types.Map     `tfsdk:"access_rights"`
	Alias                         types.String  `tfsdk:"alias"`
	Allowance                     types.Float64 `tfsdk:"allowance"`
	ApplyPolicies                 types.List    `tfsdk:"apply_policies"`
	BasicAuthData                 types.Object  `tfsdk:"basic_auth_data"`
	Certificate                   types.String  `tfsdk:"certificate"`
	DataExpires                   types.Int64   `tfsdk:"data_expires"`
	DateCreated                   types.String  `tfsdk:"date_created"`
	EnableDetailedRecording       types.Bool    `tfsdk:"enable_detailed_recording"`
	EnableHTTPSignatureValidation types.Bool    `tfsdk:"enable_http_signature_validation"`
	Expires                       types.Int64   `tfsdk:"expires"`
	HMACEnabled                   types.Bool    `tfsdk:"hmac_enabled"`
	HMACString                    types.String  `tfsdk:"hmac_string"`
	IDExtractorDeadline           types.Int64   `tfsdk:"id_extractor_deadline"`
	IsInactive                    types.Bool    `tfsdk:"is_inactive"`
	JWTData                       types.Object  `tfsdk:"jwt_data"`
	LastCheck                     types.Int64   `tfsdk:"last_check"`
	LastUpdated                   types.String  `tfsdk:"last_updated"`
	MaxQueryDepth                 types.Int64   `tfsdk:"max_query_depth"`
	MetaData                      types.String  `tfsdk:"meta_data"`
	Monitor                       types.Object  `tfsdk:"monitor"`
	OAuthClientID                 types.String  `tfsdk:"oauth_client_id"`
	OAuthKeys                     types.Map     `tfsdk:"oauth_keys"`
	OrgID                         types.String  `tfsdk:"org_id"`
	Per                           types.Float64 `tfsdk:"per"`
	QuotaMax                      types.Int64   `tfsdk:"quota_max"`
	QuotaRemaining                types.Int64   `tfsdk:"quota_remaining"`
	QuotaRenewalRate              types.Int64   `tfsdk:"quota_renewal_rate"`
	QuotaRenews                   types.Int64   `tfsdk:"quota_renews"`
	Rate                          types.Float64 `tfsdk:"rate"`
	RsaCertificateID              types.String  `tfsdk:"rsa_certificate_id"`
	SessionLifetime               types.Int64   `tfsdk:"session_lifetime"`
	Smoothing                     types.Object  `tfsdk:"smoothing"`
	Tags                          types.List    `tfsdk:"tags"`
	ThrottleInterval              types.Float64 `tfsdk:"throttle_interval"`
	ThrottleRetryLimit            types.Int64   `tfsdk:"throttle_retry_limit"`
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
			Optional:    true,
		},
	},
}

var AccessSpec = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"url": schema.StringAttribute{
			Optional:    true,
			Description: "URL that is allowed for the key.",
		},
		"methods": schema.ListAttribute{
			Description: "List of HTTP methods allowed for the URL.",
			Optional:    true,
			ElementType: types.StringType,
		},
	},
}

var RateLimitSmoothing = schema.SingleNestedAttribute{
	Description: "Smoothing configuration for the method.",
	Optional:    true,
	Attributes: map[string]schema.Attribute{
		"delay": schema.Int64Attribute{
			Description: "The delay for ratelimit smoothing",
			Optional:    true,
		},
		"enabled": schema.BoolAttribute{
			Description: "The enabled for ratelimit smoothing",
			Optional:    true,
		},
		"step": schema.Int64Attribute{
			Description: "The step for ratelimit smoothing",
			Optional:    true,
		},
		"threshold": schema.Int64Attribute{
			Description: "The threshold for ratelimit smoothing",
			Optional:    true,
		},
		"trigger": schema.Float64Attribute{
			Description: "The trigger for ratelimit smoothing",
			Optional:    true,
		},
	},
}

var RateLimitType2 = schema.SingleNestedAttribute{
	Description: "Rate limit for the HTTP method.",
	Optional:    true,
	Attributes: map[string]schema.Attribute{
		"per": schema.Int64Attribute{
			Description: "Time period for the rate limit, in seconds.",
			Optional:    true,
		},
		"rate": schema.Int64Attribute{
			Description: "Rate limit for the method, in requests per second.",
			Optional:    true,
		},
		"smoothing": RateLimitSmoothing,
	},
}

var EndpointMethod = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description: "Name of the HTTP method.",
			Optional:    true,
		},
		"limit": RateLimitType2,
	},
}

var EndpointMethods = schema.ListNestedAttribute{
	Description:  "HTTP methods allowed for the endpoint.",
	Optional:     true,
	NestedObject: EndpointMethod,
}

var Endpoint = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"path": schema.StringAttribute{
			Description: "Path of the endpoint that the key has access to.",
			Optional:    true,
		},
		"methods": EndpointMethods,
	},
}

var Endpoints = schema.ListNestedAttribute{
	Description:  "List of endpoints that the key has access to.",
	Optional:     true,
	NestedObject: Endpoint,
}

var FieldLimits = schema.SingleNestedAttribute{
	Description: "Limits for the field access rights.",
	Optional:    true,
	Attributes: map[string]schema.Attribute{
		"max_query_depth": schema.Int64Attribute{
			Description: "Maximum depth of queries allowed for the field.",
			Optional:    true,
		},
	},
}

var FieldAccessDefinition = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"field_name": schema.StringAttribute{
			Description: "Name of the field.",
			Optional:    true,
		},
		"limits": FieldLimits,
	},
}

var APILimit = schema.SingleNestedAttribute{
	Description: "Rate limits for the key.",
	Optional:    true,
	Attributes: map[string]schema.Attribute{
		"max_query_depth": schema.Int64Attribute{
			Description: "Maximum depth of queries allowed for the key.",
			Optional:    true,
		},
		"rate": schema.Float64Attribute{
			Description: "Rate limit for the key, in requests per second.",
			Optional:    true,
		},
		"per": schema.Float64Attribute{
			Description: "Time period for the rate limit, in seconds.",
			Optional:    true,
		},
		"quota_max": schema.Int64Attribute{
			Description: "Maximum quota for the key, in requests.",
			Optional:    true,
		},
		"quota_remaining": schema.Int64Attribute{
			Description: "Remaining quota for the key, in requests.",
			Optional:    true,
		},
		"quota_renewal_rate": schema.Int64Attribute{
			Description: "Rate at which the quota renews, in requests per second.",
			Optional:    true,
		},
		"quota_renews": schema.Int64Attribute{
			Description: "Time when the quota renews, in Unix timestamp format.",
			Optional:    true,
		},
		"throttle_interval": schema.Float64Attribute{
			Description: "Interval for throttling requests, in seconds.",
			Optional:    true,
		},
		"throttle_retry_limit": schema.Int64Attribute{
			Description: "Number of retries allowed for throttled requests.",
			Optional:    true,
		},
		"smoothing": RateLimitSmoothing,
	},
}

var AccessDefinition = schema.MapNestedAttribute{
	Description: "Access rights for the key.",
	Optional:    true,
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"allowance_scope": schema.StringAttribute{
				Description: "Scope of the allowance for the key.",
				Optional:    true,
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
				Optional:    true,
			},
			"api_name": schema.StringAttribute{
				Description: "Name of the API that the key has access to.",
				Optional:    true,
			},
			"disable_introspection": schema.BoolAttribute{
				Description: "Whether introspection is disabled for the key.",
				Optional:    true,
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

var BasicAuthData = schema.SingleNestedAttribute{
	Description: "Basic authentication data for the key.",
	Optional:    true,
	Attributes: map[string]schema.Attribute{
		"hash_type": schema.StringAttribute{
			Description: "Type of hash used for the basic authentication data.",
			Optional:    true,
		},
		"password": schema.StringAttribute{
			Description: "Password for the basic authentication data.",
			Optional:    true,
		},
	},
}

var JWTData = schema.SingleNestedAttribute{
	Description: "JWT data for the key.",
	Optional:    true,
	Attributes: map[string]schema.Attribute{
		"secret": schema.StringAttribute{
			Description: "Secret used for signing the JWT.",
			Optional:    true,
		},
	},
}

var Monitor = schema.SingleNestedAttribute{
	Description: "Monitoring configuration for the key.",
	Optional:    true,
	Attributes: map[string]schema.Attribute{
		"trigger_limits": schema.ListAttribute{
			Description: "List of trigger limits for monitoring.",
			Optional:    true,
			ElementType: types.StringType,
		},
	},
}

func (r *keyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"key": schema.StringAttribute{
				Description: "The API key.",
				Computed:    true,
				Sensitive:   true,
			},
			"key_hash": schema.StringAttribute{
				Description: "The API key hash",
				Computed:    true,
				Optional:    true,
			},
			"hashed": schema.BoolAttribute{
				Description: "Indicates if the key is hashed.",
				Optional:    true,
			},
			"access_rights": AccessDefinition,
			"alias": schema.StringAttribute{
				Optional:            true,
				Description:         "Alias for the key.",
				MarkdownDescription: "Alias for the key.",
			},
			"allowance": schema.Float64Attribute{
				Optional:    true,
				Description: "The number of requests allowed for the API key.",
			},
			"apply_policies": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of policy IDs to apply to the key.",
			},
			"basic_auth_data": BasicAuthData,
			"certificate": schema.StringAttribute{
				Optional:    true,
				Description: "Certificate.",
			},
			"data_expires": schema.Int64Attribute{
				Optional:    true,
				Description: "Data expiration time.",
			},
			"date_created": schema.StringAttribute{
				Optional:    true,
				Description: "The date and time when the API key was created, in Unix timestamp format.",
			},
			"enable_detailed_recording": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable detailed recording.",
			},
			"enable_http_signature_validation": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable HTTP signature validation.",
			},
			"expires": schema.Int64Attribute{
				Optional:    true,
				Description: "The expiration time of the API key, in Unix timestamp format.",
			},
			"hmac_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether HMAC is enabled.",
			},
			"hmac_string": schema.StringAttribute{
				Optional:    true,
				Description: "HMAC secret string.",
			},
			"id_extractor_deadline": schema.Int64Attribute{
				Optional:    true,
				Description: "ID extractor deadline.",
			},
			"is_inactive": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether the key is inactive.",
			},
			"jwt_data": JWTData,
			"last_check": schema.Int64Attribute{
				Optional:    true,
				Description: "The last time the API key was checked, in Unix timestamp format.",
			},
			"last_updated": schema.StringAttribute{
				Optional:    true,
				Description: "Last updated timestamp.",
			},
			"max_query_depth": schema.Int64Attribute{
				Optional:    true,
				Description: "The maximum depth of queries allowed for the API key.",
			},
			"meta_data": schema.StringAttribute{
				Description: "Custom metadata for the key.",
				Optional:    true,
			},
			"monitor": Monitor,
			"oauth_client_id": schema.StringAttribute{
				Optional:    true,
				Description: "OAuth client ID.",
			},
			"oauth_keys": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "OAuth keys.",
			},
			"org_id": schema.StringAttribute{
				Optional:    true,
				Description: "Organization ID.",
			},
			"per": schema.Float64Attribute{
				Optional:    true,
				Description: "The time period for the rate limit, in seconds.",
			},
			"quota_max": schema.Int64Attribute{
				Optional:    true,
				Description: "The maximum quota for the API key, in requests.",
			},
			"quota_remaining": schema.Int64Attribute{
				Optional:    true,
				Description: "The remaining quota for the API key, in requests.",
			},
			"quota_renewal_rate": schema.Int64Attribute{
				Optional:    true,
				Description: "The rate at which the quota renews, in requests per second.",
			},
			"quota_renews": schema.Int64Attribute{
				Optional:    true,
				Description: "The time when the quota renews, in Unix timestamp format.",
			},
			"rate": schema.Float64Attribute{
				Optional:    true,
				Description: "The rate limit for the API key, in requests per second.",
			},
			"rsa_certificate_id": schema.StringAttribute{
				Optional:    true,
				Description: "RSA certificate ID.",
			},
			"session_lifetime": schema.Int64Attribute{
				Optional:    true,
				Description: "Session lifetime.",
			},
			"smoothing": RateLimitSmoothing,
			"tags": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Tags for the key.",
			},
			"throttle_interval": schema.Float64Attribute{
				Optional:    true,
				Description: "The interval for throttling requests, in seconds.",
			},
			"throttle_retry_limit": schema.Int64Attribute{
				Optional:    true,
				Description: "The number of retries allowed for throttled requests.",
			},
		},
	}
}

func (r *keyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected *client.Client, got something else.",
		)
		return
	}

	r.client = client

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client Not Configured",
			"The client is not configured, please check your provider configuration.",
		)
	}
}

func modelToKey(ctx context.Context, data keyResourceModel) (client.Key, diag.Diagnostics) {
	var clientKey client.Key
	var accessDefinition map[string]client.AccessDefinition
	diag := data.AccessRights.ElementsAs(ctx, &accessDefinition, false)
	if diag.HasError() {
		return clientKey, diag
	}

	var oauthKeys map[string]string
	diag = data.OAuthKeys.ElementsAs(ctx, &oauthKeys, false)
	if diag.HasError() {
		return clientKey, diag
	}

	var basicAuthData client.BasicAuthData
	diag = data.BasicAuthData.As(ctx, &basicAuthData, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
	if diag.HasError() {
		return clientKey, diag
	}

	var jwtData client.JWTData
	diag = data.JWTData.As(ctx, &jwtData, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
	if diag.HasError() {
		return clientKey, diag
	}

	var applyPolicies []string
	diag = data.ApplyPolicies.ElementsAs(ctx, &applyPolicies, false)
	if diag.HasError() {
		return clientKey, diag
	}

	var monitor client.Monitor
	diag = data.Monitor.As(ctx, &monitor, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
	if diag.HasError() {
		return clientKey, diag
	}

	var metadata map[string]any
	if !data.MetaData.IsNull() && !data.MetaData.IsUnknown() {
		err := json.Unmarshal([]byte(data.MetaData.ValueString()), &metadata)
		if err != nil {
			diag.AddError(
				"Error unmarshalling metadata",
				"Could not unmarshal metadata JSON: "+err.Error())
		}
	}
	if diag.HasError() {
		return clientKey, diag
	}

	var tags []string
	diag = data.Tags.ElementsAs(ctx, &tags, false)
	if diag.HasError() {
		return clientKey, diag
	}

	var smoothing client.RateLimitSmoothing
	diag = data.Smoothing.As(ctx, &smoothing, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
	if diag.HasError() {
		return clientKey, diag
	}

	return client.Key{
		LastCheck:                     data.LastCheck.ValueInt64(),
		Allowance:                     data.Allowance.ValueFloat64(),
		Rate:                          data.Rate.ValueFloat64(),
		Per:                           data.Per.ValueFloat64(),
		ThrottleInterval:              data.ThrottleInterval.ValueFloat64(),
		ThrottleRetryLimit:            data.ThrottleRetryLimit.ValueInt64(),
		MaxQueryDepth:                 data.MaxQueryDepth.ValueInt64(),
		DateCreated:                   data.DateCreated.ValueString(),
		Expires:                       data.Expires.ValueInt64(),
		QuotaMax:                      data.QuotaMax.ValueInt64(),
		QuotaRenews:                   data.QuotaRenews.ValueInt64(),
		QuotaRemaining:                data.QuotaRemaining.ValueInt64(),
		QuotaRenewalRate:              data.QuotaRenewalRate.ValueInt64(),
		AccessRights:                  accessDefinition,
		OrgID:                         data.OrgID.ValueString(),
		OauthClientID:                 data.OAuthClientID.ValueString(),
		OauthKeys:                     oauthKeys,
		Certificate:                   data.Certificate.ValueString(),
		BasicAuthData:                 basicAuthData,
		JWTData:                       jwtData,
		HMACEnabled:                   data.HMACEnabled.ValueBool(),
		EnableHTTPSignatureValidation: data.EnableHTTPSignatureValidation.ValueBool(),
		HmacSecret:                    data.HMACString.ValueString(),
		RSACertificateId:              data.RsaCertificateID.ValueString(),
		IsInactive:                    data.IsInactive.ValueBool(),
		ApplyPolicies:                 applyPolicies,
		DataExpires:                   data.DataExpires.ValueInt64(),
		Monitor:                       monitor,
		EnableDetailedRecording:       data.EnableDetailedRecording.ValueBool(),
		MetaData:                      metadata,
		Tags:                          tags,
		Alias:                         data.Alias.ValueString(),
		LastUpdated:                   data.LastUpdated.ValueString(),
		IdExtractorDeadline:           data.IDExtractorDeadline.ValueInt64(),
		SessionLifetime:               data.SessionLifetime.ValueInt64(),
		Smoothing:                     smoothing,
	}, nil
}

func (r *keyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data keyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	key, diag := modelToKey(ctx, data)
	if diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}

	// Create API call logic
	createKeyResp, err := r.client.CreateKeyWithHashed(key, data.Hashed.ValueBool())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating key",
			"Could not create key, unexpected error: "+err.Error(),
		)
		return
	}

	data.Key = types.StringValue(createKeyResp.Key)
	data.KeyHash = types.StringValue(createKeyResp.KeyHash)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func applyKeyDataToModel(ctx context.Context, key client.Key, data *keyResourceModel) diag.Diagnostics {
	// data.AccessRights, _ = types.MapValueFrom(ctx, AccessDefinition.NestedObjectType(), key.AccessRights)
	data.Alias = types.StringValue(key.Alias)
	data.Allowance = types.Float64Value(key.Allowance)
	data.ApplyPolicies, _ = types.ListValueFrom(ctx, types.StringType, key.ApplyPolicies)
	// data.BasicAuthData, _ = types.ObjectValueFrom(ctx, BasicAuthData.NestedObjectType(), key.BasicAuthData)
	data.Certificate = types.StringValue(key.Certificate)
	data.DataExpires = types.Int64Value(key.DataExpires)
	data.DateCreated = types.StringValue(key.DateCreated)
	data.EnableDetailedRecording = types.BoolValue(key.EnableDetailedRecording)
	data.EnableHTTPSignatureValidation = types.BoolValue(key.EnableHTTPSignatureValidation)
	data.Expires = types.Int64Value(key.Expires)
	data.HMACEnabled = types.BoolValue(key.HMACEnabled)
	data.HMACString = types.StringValue(key.HmacSecret)
	data.IDExtractorDeadline = types.Int64Value(key.IdExtractorDeadline)
	data.IsInactive = types.BoolValue(key.IsInactive)
	// data.JWTData, _ = types.ObjectValueFrom(ctx, JWTData.NestedObjectType(), key.JWTData)
	data.LastCheck = types.Int64Value(key.LastCheck)
	data.LastUpdated = types.StringValue(key.LastUpdated)
	data.MaxQueryDepth = types.Int64Value(key.MaxQueryDepth)
	if key.MetaData != nil {
		metaDataJSON, _ := json.Marshal(key.MetaData)
		data.MetaData = types.StringValue(string(metaDataJSON))
	} else {
		data.MetaData = basetypes.NewStringNull()
	}
	// data.Monitor, _ = types.ObjectValueFrom(ctx, Monitor.NestedObjectType(), key.Monitor)
	data.OAuthClientID = types.StringValue(key.OauthClientID)
	// data.OAuthKeys, _ = types.MapValueFrom(ctx, basetypes.StringType, key.OauthKeys)
	data.OrgID = types.StringValue(key.OrgID)
	data.Per = types.Float64Value(key.Per)
	data.QuotaMax = types.Int64Value(key.QuotaMax)
	data.QuotaRemaining = types.Int64Value(key.QuotaRemaining)
	data.QuotaRenewalRate = types.Int64Value(key.QuotaRenewalRate)
	data.QuotaRenews = types.Int64Value(key.QuotaRenews)
	data.Rate = types.Float64Value(key.Rate)
	data.RsaCertificateID = types.StringValue(key.RSACertificateId)
	data.SessionLifetime = types.Int64Value(key.SessionLifetime)
	// data.Smoothing, _ = types.ObjectValueFrom(ctx, RateLimitSmoothing.NestedObjectType(), key.Smoothing)
	data.Tags, _ = types.ListValueFrom(ctx, types.StringType, key.Tags)
	data.ThrottleInterval = types.Float64Value(key.ThrottleInterval)
	data.ThrottleRetryLimit = types.Int64Value(key.ThrottleRetryLimit)
	return nil
}

func (r *keyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data keyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	keyId := data.Key.ValueString()
	if data.Hashed.ValueBool() {
		keyId = data.KeyHash.ValueString()
	}
	key, err := r.client.GetKeyWithHashed(keyId, data.Hashed.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading key",
			"Could not read key, unexpected error: "+err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(applyKeyDataToModel(ctx, key, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

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

	key, diag := modelToKey(ctx, data)
	if diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}
	// Update API call logic
	keyId := data.Key.ValueString()
	if data.Hashed.ValueBool() {
		keyId = data.KeyHash.ValueString()
	}
	_, err := r.client.UpdateKeyWithHashed(keyId, key, data.Hashed.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating key",
			"Could not update key, unexpected error: "+err.Error(),
		)
		return
	}
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
	keyId := data.Key.ValueString()
	if data.Hashed.ValueBool() {
		keyId = data.KeyHash.ValueString()
	}
	err := r.client.DeleteKeyWithHashed(keyId, data.Hashed.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting key",
			"Could not delete key, unexpected error: "+err.Error(),
		)
		return
	}
}
