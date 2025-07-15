package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/TykTechnologies/graphql-go-tools/pkg/graphql"
)

type AccessSpec struct {
	URL     string   `json:"url"`
	Methods []string `json:"methods"`
}
type RateLimitSmoothing struct {
	// Enabled indicates if rate limit smoothing is active.
	Enabled bool `json:"enabled"`

	// Threshold is the initial rate limit beyond which smoothing will be applied. It is a count of requests during the `per` interval and should be less than the maximum configured `rate`.
	Threshold int64 `json:"threshold"`

	// Trigger is a fraction (typically in the range 0.1-1.0) of the step at which point a smoothing event will be emitted as the request rate approaches the current allowance.
	Trigger float64 `json:"trigger"`

	// Step is the increment by which the current allowance will be increased or decreased each time a smoothing event is emitted.
	Step int64 `json:"step"`

	// Delay is a hold-off between smoothing events and controls how frequently the current allowance will step up or down (in seconds).
	Delay int64 `json:"delay"`
}
type RateLimit struct {
	// Rate is the allowed number of requests per interval.
	Rate float64 `json:"rate"`
	// Per is the interval at which rate limit is enforced.
	Per float64 `json:"per"`

	// Smoothing contains rate limit smoothing settings.
	Smoothing *RateLimitSmoothing `json:"smoothing,omitempty"`
}
type APILimit struct {
	RateLimit
	ThrottleInterval   float64 `json:"throttle_interval"`
	ThrottleRetryLimit int     `json:"throttle_retry_limit"`
	MaxQueryDepth      int     `json:"max_query_depth"`
	QuotaMax           int64   `json:"quota_max"`
	QuotaRenews        int64   `json:"quota_renews"`
	QuotaRemaining     int64   `json:"quota_remaining"`
	QuotaRenewalRate   int64   `json:"quota_renewal_rate"`
}
type FieldAccessDefinition struct {
	TypeName  string      `json:"type_name"`
	FieldName string      `json:"field_name"`
	Limits    FieldLimits `json:"limits"`
}

type FieldLimits struct {
	MaxQueryDepth int `json:"max_query_depth"`
}

type BasicAuthData struct {
	Password string `json:"password"`
	Hash     string `json:"hash_type"`
}

type JWTData struct {
	Secret string `json:"secret"`
}

type Monitor struct {
	TriggerLimits []float64 `json:"trigger_limits"`
}
type Endpoint struct {
	Path    string          `json:"path,omitempty"`
	Methods EndpointMethods `json:"methods,omitempty"`
}

type EndpointMethod struct {
	Name  string    `json:"name,omitempty"`
	Limit RateLimit `json:"limit,omitempty"`
}

type EndpointMethods []EndpointMethod
type Endpoints []Endpoint
type AccessDefinition struct {
	APIName              string                  `json:"api_name"`
	APIID                string                  `json:"api_id"`
	Versions             []string                `json:"versions"`
	AllowedURLs          []AccessSpec            `json:"allowed_urls"` // mapped string MUST be a valid regex
	RestrictedTypes      []graphql.Type          `json:"restricted_types"`
	AllowedTypes         []graphql.Type          `json:"allowed_types"`
	Limit                APILimit                `json:"limit"`
	FieldAccessRights    []FieldAccessDefinition `json:"field_access_rights"`
	DisableIntrospection bool                    `json:"disable_introspection"`
	AllowanceScope       string                  `json:"allowance_scope"`
	Endpoints            Endpoints               `json:"endpoints,omitempty"`
}

type Key struct {
	LastCheck                     int64                       `json:"last_check"`
	Allowance                     float64                     `json:"allowance"`
	Rate                          float64                     `json:"rate"`
	Per                           float64                     `json:"per"`
	ThrottleInterval              float64                     `json:"throttle_interval"`
	ThrottleRetryLimit            int64                       `json:"throttle_retry_limit"`
	MaxQueryDepth                 int64                       `json:"max_query_depth"`
	DateCreated                   string                      `json:"date_created"`
	Expires                       int64                       `json:"expires"`
	QuotaMax                      int64                       `json:"quota_max"`
	QuotaRenews                   int64                       `json:"quota_renews"`
	QuotaRemaining                int64                       `json:"quota_remaining"`
	QuotaRenewalRate              int64                       `json:"quota_renewal_rate"`
	AccessRights                  map[string]AccessDefinition `json:"access_rights"`
	OrgID                         string                      `json:"org_id"`
	OauthClientID                 string                      `json:"oauth_client_id"`
	OauthKeys                     map[string]string           `json:"oauth_keys"`
	Certificate                   string                      `json:"certificate"`
	BasicAuthData                 BasicAuthData               `json:"basic_auth_data"`
	JWTData                       JWTData                     `json:"jwt_data"`
	HMACEnabled                   bool                        `json:"hmac_enabled"`
	EnableHTTPSignatureValidation bool                        `json:"enable_http_signature_validation"`
	HmacSecret                    string                      `json:"hmac_string"`
	RSACertificateId              string                      `json:"rsa_certificate_id"`
	IsInactive                    bool                        `json:"is_inactive"`
	ApplyPolicyID                 string                      `json:"apply_policy_id"`
	ApplyPolicies                 []string                    `json:"apply_policies"`
	DataExpires                   int64                       `json:"data_expires"`
	Monitor                       Monitor                     `json:"monitor"`
	EnableDetailedRecording       bool                        `json:"enable_detailed_recording"`
	MetaData                      map[string]interface{}      `json:"meta_data"`
	Tags                          []string                    `json:"tags"`
	Alias                         string                      `json:"alias"`
	LastUpdated                   string                      `json:"last_updated"`
	IdExtractorDeadline           int64                       `json:"id_extractor_deadline"`
	SessionLifetime               int64                       `json:"session_lifetime"`
	Smoothing                     RateLimitSmoothing          `json:"smoothing"`
}

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
