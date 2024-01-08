package client_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/cds-snc/notification-go-client"
)

func TestNewClientWithInvalidApiKey(t *testing.T) {
	apiKey := "bad_key"

	_, err := NewClient(apiKey)

	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	if err.Error() != "invalid API key" {
		t.Errorf("Expected error message to be invalid API key, got %s", err.Error())
	}
}

func TestNewClientWithValidApiKey(t *testing.T) {
	apiKey := "gcntfy-testing11-00000000-0000-0000-0000-000000000000-00000000-0000-0000-0000-00000000000000"

	client, err := NewClient(apiKey)

	if err != nil {
		t.Errorf("Expected no error, got %s", err.Error())
	}

	if client.ApiKey != apiKey {
		t.Errorf("Expected API key to be %s, got %s", apiKey, client.ApiKey)
	}

	if client.Hostname != "https://api.notification.canada.ca" {
		t.Errorf("Expected hostname to be https://api.notification.canada.ca, got %s", client.Hostname)
	}
}

func TestDoGetRequestWithNoError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request URL
		if r.URL.Path != "/v2/notifications/email" {
			t.Errorf("Expected request to /v2/notifications/email, got %s", r.URL.Path)
		}

		// Verify the request method
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Read the request body
		_, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Error reading request body: %s", err)
		}

		// Write a response
		w.WriteHeader(http.StatusCreated)
		response := Response{
			Id:        "00000000-0000-0000-0000-000000000000",
			Reference: "00000000-0000-0000-0000-000000000000",
			Uri:       "https://api.notification.canada.ca/v2/notifications/00000000-0000-0000-0000-000000000000",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create a client with the mock server URL
	client := &Client{
		ApiKey:     "test",
		HttpClient: http.Client{},
		Hostname:   server.URL,
	}

	resp, errors := client.DoGetRequest("/v2/notifications/email")

	if resp.StatusCode != 201 {
		t.Errorf("Expected status code to be 201, got %d", resp.StatusCode)
	}

	if errors != nil {
		t.Errorf("Expected no errors, got %s", errors)
	}

}

func TestDoPostRequestWithNoError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request URL
		if r.URL.Path != "/v2/notifications/email" {
			t.Errorf("Expected request to /v2/notifications/email, got %s", r.URL.Path)
		}

		// Verify the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Read the request body
		_, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Error reading request body: %s", err)
		}

		// Write a response
		w.WriteHeader(http.StatusCreated)
		response := Response{
			Id:        "00000000-0000-0000-0000-000000000000",
			Reference: "00000000-0000-0000-0000-000000000000",
			Uri:       "https://api.notification.canada.ca/v2/notifications/00000000-0000-0000-0000-000000000000",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create a client with the mock server URL
	client := &Client{
		ApiKey:     "test",
		HttpClient: http.Client{},
		Hostname:   server.URL,
	}

	resp, errors := client.DoPostRequest("/v2/notifications/email", []byte("test"))

	if resp.StatusCode != 201 {
		t.Errorf("Expected status code to be 201, got %d", resp.StatusCode)
	}

	if errors != nil {
		t.Errorf("Expected no errors, got %s", errors)
	}

}
