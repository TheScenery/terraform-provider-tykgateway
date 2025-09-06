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

const KEY = "key"
const KEY_HASH = "key_hash"

type keyResourceModel struct {
	Hashed types.Bool   `tfsdk:"hashed"`
	Key    types.String `tfsdk:"key"`
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
			"key": schema.StringAttribute{
				Description: "The key request json string",
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
	err := json.Unmarshal([]byte(data.Key.ValueString()), &key)

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

	key[KEY] = createKeyResponse.Key
	key[KEY_HASH] = createKeyResponse.KeyHash

	updatedKey, err := json.Marshal(key)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error marshalling key JSON",
			"Could not marshal key JSON, unexpected error: "+err.Error(),
		)
		return
	}

	data.Key = types.StringValue(string(updatedKey))

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

	var key map[string]any
	err := json.Unmarshal([]byte(data.Key.ValueString()), &key)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing key JSON",
			"Could not parse key JSON, unexpected error: "+err.Error(),
		)
		return
	}

	// Read API call logic
	keyId := key[KEY].(string)
	if data.Hashed.ValueBool() {
		keyId = key[KEY_HASH].(string)
	}
	keyResponse, err := r.client.GetKeyWithHashed(keyId, data.Hashed.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading key",
			"Could not read key, unexpected error: "+err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	updatedKey, err := json.Marshal(keyResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error marshalling key JSON",
			"Could not marshal key JSON, unexpected error: "+err.Error(),
		)
		return
	}

	data.Key = types.StringValue(string(updatedKey))

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
	err := json.Unmarshal([]byte(data.Key.ValueString()), &key)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing key JSON",
			"Could not parse key JSON, unexpected error: "+err.Error(),
		)
		return
	}

	// Update API call logic
	keyId := key[KEY].(string)
	if data.Hashed.ValueBool() {
		keyId = key[KEY_HASH].(string)
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
	err := json.Unmarshal([]byte(data.Key.ValueString()), &key)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing key JSON",
			"Could not parse key JSON, unexpected error: "+err.Error(),
		)
		return
	}

	// Delete API call logic
	keyId := key[KEY].(string)
	if data.Hashed.ValueBool() {
		keyId = key[KEY_HASH].(string)
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
