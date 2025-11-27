package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	token      string
	teamID     string
	httpClient *http.Client
}

func NewClient(baseURL string, token string) *Client {
	return &Client{
		baseURL:    trimBaseURL(baseURL),
		token:      token,
		teamID:     "",
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *Client) WithHTTPClient(hc *http.Client) *Client {
	copy := *c
	copy.httpClient = hc
	return &copy
}

func (c *Client) WithTeam(teamID string) *Client {
	copy := *c
	copy.teamID = teamID
	return &copy
}

func (c *Client) do(method string, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, c.baseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	return resp, nil
}

func trimBaseURL(baseUrl string) string {
	return strings.TrimSuffix(baseUrl, "/")
}

func marshalJSON(v any) (*bytes.Buffer, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json: %w", err)
	}
	return bytes.NewBuffer(raw), nil
}
