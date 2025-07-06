package client

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	Host       string
	ApiKey     string
	HTTPClient *http.Client
}

func NewClient(host, apiKey string) *Client {
	return &Client{
		Host:       host,
		ApiKey:     apiKey,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("X-Tyk-Authorization", c.ApiKey)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
