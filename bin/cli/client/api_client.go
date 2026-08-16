package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/moroz/homeosapiens-go/bin/cli/types"
)

type ApiClient struct {
	apiToken string
	baseUrl  string
}

type ApiClientConfig struct {
	APIToken string
	BaseURL  string
}

func NewApiClient(config *ApiClientConfig) *ApiClient {
	return &ApiClient{
		apiToken: config.APIToken,
		baseUrl:  config.BaseURL,
	}
}

func (c *ApiClient) GetSession() (*types.User, error) {
	endpoint, err := url.JoinPath(c.baseUrl, "/api/admin/session")
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.apiToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: want 200, got %v", resp.StatusCode)
	}
	var user types.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
