package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Key struct {
}

type ApiModifyKeySuccess struct {
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
