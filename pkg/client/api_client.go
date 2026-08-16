package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type APIClient struct {
	BaseURL    string
	AuthToken  string
	HTTPClient *http.Client
}

func NewAPIClient(baseURL, authToken string) *APIClient {
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	return &APIClient{
		BaseURL:   baseURL,
		AuthToken: authToken,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *APIClient) Push(projectID, envName, base64Payload, message string) (int, error) {
	url := fmt.Sprintf("%s/api/v1/sync/push", c.BaseURL)
	bodyMap := map[string]string{
		"project_id": projectID,
		"env_name":   envName,
		"payload":    base64Payload,
		"message":    message,
	}

	jsonBytes, err := json.Marshal(bodyMap)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return 0, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to send HTTP request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("server responded with status %d: %s", resp.StatusCode, string(respBody))
	}

	var res struct {
		Version int `json:"version"`
	}
	_ = json.Unmarshal(respBody, &res)

	return res.Version, nil
}

func (c *APIClient) Pull(projectID, envName, version string) (string, int, error) {
	url := fmt.Sprintf("%s/api/v1/sync/pull?project_id=%s&env=%s", c.BaseURL, projectID, envName)
	if version != "" {
		url += "&version=" + version
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", 0, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("failed to send HTTP request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("server responded with status %d: %s", resp.StatusCode, string(respBody))
	}

	var res struct {
		Payload string `json:"payload"`
		Version int    `json:"version"`
	}
	if err := json.Unmarshal(respBody, &res); err != nil {
		return "", 0, fmt.Errorf("failed to parse server response: %w", err)
	}

	return res.Payload, res.Version, nil
}

func (c *APIClient) GetHistory(projectID, envName string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/projects/%s/history", c.BaseURL, projectID)
	if envName != "" {
		url += "?env=" + envName
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server responded with status %d: %s", resp.StatusCode, string(respBody))
	}

	var res struct {
		History []map[string]interface{} `json:"history"`
	}
	if err := json.Unmarshal(respBody, &res); err != nil {
		return nil, fmt.Errorf("failed to parse server response: %w", err)
	}

	return res.History, nil
}
