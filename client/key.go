package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/TykTechnologies/graphql-go-tools/pkg/graphql"
)

type AccessSpec struct {
	URL     string   `json:"url,omitempty" tfsdk:"url"`
	Methods []string `json:"methods,omitempty" tfsdk:"methods"` // HTTP methods allowed for this URL
}
type RateLimitSmoothing struct {
	// Enabled indicates if rate limit smoothing is active.
	Enabled bool `json:"enabled,omitempty" tfsdk:"enabled"`

	// Threshold is the initial rate limit beyond which smoothing will be applied. It is a count of requests during the `per` interval and should be less than the maximum configured `rate`.
	Threshold int64 `json:"threshold,omitempty" tfsdk:"threshold"`

	// Trigger is a fraction (typically in the range 0.1-1.0) of the step at which point a smoothing event will be emitted as the request rate approaches the current allowance.
	Trigger float64 `json:"trigger,omitempty" tfsdk:"trigger"`

	// Step is the increment by which the current allowance will be increased or decreased each time a smoothing event is emitted.
	Step int64 `json:"step,omitempty" tfsdk:"step"`

	// Delay is a hold-off between smoothing events and controls how frequently the current allowance will step up or down (in seconds).
	Delay int64 `json:"delay,omitempty" tfsdk:"delay"`
}
type RateLimit struct {
	// Rate is the allowed number of requests per interval.
	Rate float64 `json:"rate,omitempty" tfsdk:"rate"`
	// Per is the interval at which rate limit is enforced.
	Per float64 `json:"per,omitempty" tfsdk:"per"`

	// Smoothing contains rate limit smoothing settings.
	Smoothing *RateLimitSmoothing `json:"smoothing,omitempty" tfsdk:"smoothing"`
}
type APILimit struct {
	RateLimit
	ThrottleInterval   float64 `json:"throttle_interval,omitempty" tfsdk:"throttle_interval"`
	ThrottleRetryLimit int     `json:"throttle_retry_limit,omitempty" tfsdk:"throttle_retry_limit"`
	MaxQueryDepth      int     `json:"max_query_depth,omitempty" tfsdk:"max_query_depth"`
	QuotaMax           int64   `json:"quota_max,omitempty" tfsdk:"quota_max"`
	QuotaRenews        int64   `json:"quota_renews,omitempty" tfsdk:"quota_renews"`
	QuotaRemaining     int64   `json:"quota_remaining,omitempty" tfsdk:"quota_remaining"`
	QuotaRenewalRate   int64   `json:"quota_renewal_rate,omitempty" tfsdk:"quota_renewal_rate"`
}
type FieldAccessDefinition struct {
	TypeName  string      `json:"type_name,omitempty" tfsdk:"type_name"`
	FieldName string      `json:"field_name,omitempty" tfsdk:"field_name"`
	Limits    FieldLimits `json:"limits,omitempty" tfsdk:"limits"`
}

type FieldLimits struct {
	MaxQueryDepth int `json:"max_query_depth,omitempty" tfsdk:"max_query_depth"`
}

type BasicAuthData struct {
	Password string `json:"password,omitempty" tfsdk:"password"`
	Hash     string `json:"hash_type,omitempty" tfsdk:"hash_type"`
}

type JWTData struct {
	Secret string `json:"secret,omitempty" tfsdk:"secret"`
}

type Monitor struct {
	TriggerLimits []float64 `json:"trigger_limits,omitempty" tfsdk:"trigger_limits"` // List of limits that trigger monitoring
}
type Endpoint struct {
	Path    string          `json:"path,omitempty" tfsdk:"path"`
	Methods EndpointMethods `json:"methods,omitempty" tfsdk:"methods"`
}

type EndpointMethod struct {
	Name  string    `json:"name,omitempty" tfsdk:"name"`
	Limit RateLimit `json:"limit,omitempty" tfsdk:"limit"`
}

type EndpointMethods []EndpointMethod
type Endpoints []Endpoint
type AccessDefinition struct {
	APIName              string                  `json:"api_name,omitempty" tfsdk:"api_name"`
	APIID                string                  `json:"api_id,omitempty" tfsdk:"api_id"`
	Versions             []string                `json:"versions,omitempty" tfsdk:"versions"`
	AllowedURLs          []AccessSpec            `json:"allowed_urls,omitempty" tfsdk:"allowed_urls"`
	RestrictedTypes      []graphql.Type          `json:"restricted_types,omitempty" tfsdk:"restricted_types"`
	AllowedTypes         []graphql.Type          `json:"allowed_types,omitempty" tfsdk:"allowed_types"`
	Limit                APILimit                `json:"limit,omitempty" tfsdk:"limit"`
	FieldAccessRights    []FieldAccessDefinition `json:"field_access_rights,omitempty" tfsdk:"field_access_rights"`
	DisableIntrospection bool                    `json:"disable_introspection,omitempty" tfsdk:"disable_introspection"`
	AllowanceScope       string                  `json:"allowance_scope,omitempty" tfsdk:"allowance_scope"`
	Endpoints            Endpoints               `json:"endpoints,omitempty" tfsdk:"endpoints"`
}

type Key map[string]any

type ApiModifyKeySuccess struct {
	Key     string `json:"key"`
	Status  string `json:"status"`
	Action  string `json:"action"`
	KeyHash string `json:"key_hash,omitempty"`
}

func (c *Client) CreateKey(key Key) (ApiModifyKeySuccess, error) {
	return c.CreateKeyWithHashed(key, false)
}

func (c *Client) CreateKeyWithHashed(key Key, hashed bool) (ApiModifyKeySuccess, error) {
	var apiModifyKeySuccess ApiModifyKeySuccess

	rb, err := json.Marshal(key)
	if err != nil {
		return apiModifyKeySuccess, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/tyk/keys?hashed=%t", c.Host, hashed), strings.NewReader(string(rb)))
	if err != nil {
		return apiModifyKeySuccess, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return apiModifyKeySuccess, err
	}

	err = json.Unmarshal(body, &apiModifyKeySuccess)
	if err != nil {
		return ApiModifyKeySuccess{}, err
	}

	return apiModifyKeySuccess, nil
}

func (c *Client) GetKey(keyId string) (Key, error) {
	return c.GetKeyWithHashed(keyId, false)
}

func (c *Client) GetKeyWithHashed(keyId string, hashed bool) (Key, error) {
	var key Key
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/tyk/keys/%s?hashed=%t", c.Host, keyId, hashed), nil)
	if err != nil {
		return key, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return key, err
	}

	err = json.Unmarshal(body, &key)
	if err != nil {
		return key, err
	}
	return key, nil
}

func (c *Client) DeleteKey(keyId string) error {
	return c.DeleteKeyWithHashed(keyId, false)
}

func (c *Client) DeleteKeyWithHashed(keyId string, hashed bool) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/tyk/keys/%s?hashed=%t", c.Host, keyId, hashed), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateKey(keyId string, key Key) (ApiModifyKeySuccess, error) {
	return c.UpdateKeyWithHashed(keyId, key, false)
}

func (c *Client) UpdateKeyWithHashed(keyId string, key Key, hashed bool) (ApiModifyKeySuccess, error) {
	var apiModifyKeySuccess ApiModifyKeySuccess

	rb, err := json.Marshal(key)
	if err != nil {
		return apiModifyKeySuccess, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/tyk/keys/%s?hashed=%t", c.Host, keyId, hashed), strings.NewReader(string(rb)))
	if err != nil {
		return apiModifyKeySuccess, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return apiModifyKeySuccess, err
	}

	err = json.Unmarshal(body, &apiModifyKeySuccess)
	if err != nil {
		return ApiModifyKeySuccess{}, err
	}

	return apiModifyKeySuccess, nil
}
