// HTTP Client for Notification API

package client

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	ApiKey     string
	HttpClient http.Client
	Hostname   string
}

type ResponseError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type responseTemplate struct {
	Id      string `json:"id"`
	Uri     string `json:"uri"`
	Version int    `json:"version"`
}

type Response struct {
	// Valid Response
	Id        string            `json:"id"`
	Reference string            `json:"reference"`
	Content   map[string]string `json:"content"`
	Uri       string            `json:"uri"`
	Template  responseTemplate  `json:"template"`

	// Error Response
	StatusCode int             `json:"status_code"`
	Errors     []ResponseError `json:"errors"`
}

func NewClient(apiKey string) (Client, error) {
	// Validate API key
	if !validateApiKey(apiKey) {
		return Client{}, errors.New("invalid API key")
	}

	// Set default hostname
	client := Client{
		ApiKey:     apiKey,
		HttpClient: http.Client{Timeout: 10 * time.Second},
		Hostname:   "https://api.notification.canada.ca",
	}

	return client, nil
}

func (c Client) DoGetRequest(endpoint string) (*http.Response, error) {
	method := "GET"

	resource := fmt.Sprintf("%s%s", c.Hostname, endpoint)

	req, err := http.NewRequest(method, resource, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("ApiKey-v1 %s", c.ApiKey))

	return c.HttpClient.Do(req)
}

func (c Client) DoPostRequest(endpoint string, body []byte) (*http.Response, error) {
	method := "POST"
	contentType := "application/json"

	resource := fmt.Sprintf("%s%s", c.Hostname, endpoint)

	req, err := http.NewRequest(method, resource, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", fmt.Sprintf("ApiKey-v1 %s", c.ApiKey))

	return c.HttpClient.Do(req)
}

func validateApiKey(apiKey string) bool {
	return len(apiKey) >= 72
}
