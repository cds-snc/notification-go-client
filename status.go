package client

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

type StatusResponse struct {
	// Valid Response
	Id                string           `json:"id"`
	Reference         string           `json:"reference"`
	EmailAddress      string           `json:"email_address"`
	PhoneNumber       string           `json:"phone_number"`
	Type              string           `json:"type"`
	Status            string           `json:"status"`
	StatusDescription string           `json:"status_description"`
	ProviderResponse  string           `json:"provider_response"`
	Template          responseTemplate `json:"template"`
	Body              string           `json:"body"`
	Subject           string           `json:"subject"`
	CreatedAt         time.Time        `json:"created_at"`
	CreatedByName     string           `json:"created_by_name"`
	SentAt            time.Time        `json:"sent_at"`
	CompletedAt       time.Time        `json:"completed_at"`

	// Error Response
	StatusCode int             `json:"status_code"`
	Errors     []ResponseError `json:"errors"`
}

type StatusResponses struct {
	// Valid Response
	Notifications []StatusResponse `json:"notifications"`
	Links         Link             `json:"links"`

	// Error Response
	StatusCode int             `json:"status_code"`
	Errors     []ResponseError `json:"errors"`
}

type Link struct {
	Current string `json:"current"`
	Next    string `json:"next"`
}

type StatusQueryOptions struct {
	OlderThan    string `url:"older_than,omitempty"`
	Reference    string `url:"reference,omitempty"`
	Status       string `url:"status,omitempty"`
	TemplateType string `url:"template_type,omitempty"`
}

func doGetStatus[T StatusResponse | StatusResponses](c Client, url string, response T) (T, int, error) {
	resp, err := c.DoGetRequest(url)

	if err != nil {
		return response, 0, fmt.Errorf("error calling status endpoint: %s", err)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		return response, 0, fmt.Errorf("error decoding status response: %s", err)
	}

	return response, resp.StatusCode, nil
}

func (c Client) GetStatus(options StatusQueryOptions) (StatusResponses, error) {
	v, _ := query.Values(options)

	response, statusCode, err := doGetStatus(c, "/v2/notifications?"+v.Encode(), StatusResponses{})

	if err != nil {
		return StatusResponses{}, err
	}

	response.StatusCode = statusCode

	return response, nil
}

func (c Client) GetStatusById(id string) (StatusResponse, error) {
	response, statusCode, err := doGetStatus(c, "/v2/notifications/"+id, StatusResponse{})

	if err != nil {
		return StatusResponse{}, err
	}

	response.StatusCode = statusCode

	return response, nil
}

func (s *StatusResponses) HasNext() bool {
	return s.Links.Next != ""
}

func (c Client) NextStatusPage(s StatusResponses) (StatusResponses, error) {
	url := strings.Replace(s.Links.Next, c.Hostname, "", 1)

	response, statusCode, err := doGetStatus(c, url, StatusResponses{})

	if err != nil {
		return StatusResponses{}, err
	}

	response.StatusCode = statusCode

	return response, nil
}
