package provider

import (
	"context"
	"encoding/json"
	"terraform-provider-tykgateway/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	Hashed    types.Bool   `tfsdk:"hashed"`
	KeyConfig types.String `tfsdk:"key_config"`
	Key       types.String `tfsdk:"key"`
	KeyHash   types.String `tfsdk:"key_hash"`
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
			"key_config": schema.StringAttribute{
				Description: "The key config json string",
				Required:    true,
			},
			"key": schema.StringAttribute{
				Description: "The key.",
				Computed:    true,
			},
			"key_hash": schema.StringAttribute{
				Description: "The key hash.",
				Computed:    true,
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

func (r *keyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data keyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var key map[string]any
	err := json.Unmarshal([]byte(data.KeyConfig.ValueString()), &key)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing key JSON",
			"Could not parse key JSON, unexpected error: "+err.Error(),
		)
		return
	}

	// Create API call logic
	createKeyResponse, err := r.client.CreateKeyWithHashed(key, data.Hashed.ValueBool())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating key",
			"Could not create key, unexpected error: "+err.Error(),
		)
		return
	}

	data.Key = types.StringValue(createKeyResponse.Key)
	data.KeyHash = types.StringValue(createKeyResponse.KeyHash)

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
	keyId := data.Key.ValueString()
	if data.Hashed.ValueBool() {
		keyId = data.KeyHash.ValueString()
	}
	_, err := r.client.GetKeyWithHashed(keyId, data.Hashed.ValueBool())
	if err != nil {
		// TODO: check for 404 Not Found error and remove from state
		resp.Diagnostics.AddError(
			"Error reading key",
			"Could not read key, unexpected error: "+err.Error(),
		)
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

	var key map[string]any
	err := json.Unmarshal([]byte(data.KeyConfig.ValueString()), &key)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing key JSON",
			"Could not parse key JSON, unexpected error: "+err.Error(),
		)
		return
	}

	// Update API call logic
	keyId := data.Key.ValueString()
	if data.Hashed.ValueBool() {
		keyId = data.KeyHash.ValueString()
	}
	_, err = r.client.UpdateKeyWithHashed(keyId, key, data.Hashed.ValueBool())
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

	var key map[string]any
	err := json.Unmarshal([]byte(data.KeyConfig.ValueString()), &key)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing key JSON",
			"Could not parse key JSON, unexpected error: "+err.Error(),
		)
		return
	}

	// Delete API call logic
	keyId := data.Key.ValueString()
	if data.Hashed.ValueBool() {
		keyId = data.KeyHash.ValueString()
	}
	err = r.client.DeleteKeyWithHashed(keyId, data.Hashed.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting key",
			"Could not delete key, unexpected error: "+err.Error(),
		)
		return
	}
}
