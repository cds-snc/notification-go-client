package client_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	. "github.com/cds-snc/notification-go-client"
)

func TestGetStatus(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request URL
		if r.URL.Path != "/v2/notifications" {
			t.Errorf("Expected request to /v2/notifications, got %s", r.URL.Path)
		}

		// Verfiy the query string
		if r.URL.RawQuery != "template_type=email" {
			t.Errorf("Expected query string to be template_type=email, got %s", r.URL.RawQuery)
		}

		// Verify the request method
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Write a response
		w.WriteHeader(http.StatusOK)
		response := StatusResponses{
			Notifications: []StatusResponse{
				{
					Id:        "00000000-0000-0000-0000-000000000000",
					Reference: "00000000-0000-0000-0000-000000000000",
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))

	defer server.Close()

	c, _ := NewClient("test")
	c.Hostname = server.URL

	got, err := c.GetStatus(StatusQueryOptions{
		TemplateType: "email",
	})

	if err != nil {
		t.Errorf("Error calling GetStatus(): %s", err)
	}

	want := StatusResponses{
		Notifications: []StatusResponse{
			{
				Id:        "00000000-0000-0000-0000-000000000000",
				Reference: "00000000-0000-0000-0000-000000000000",
			},
		},
		StatusCode: 200,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetStatus() = %v, want %v", got, want)
	}
}

func TestGetStatusById(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request URL
		if r.URL.Path != "/v2/notifications/00000000-0000-0000-0000-000000000000" {
			t.Errorf("Expected request to /v2/notifications/00000000-0000-0000-0000-000000000000, got %s", r.URL.Path)
		}

		// Verify the request method
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Write a response
		w.WriteHeader(http.StatusOK)
		response := StatusResponse{
			Id:        "00000000-0000-0000-0000-000000000000",
			Reference: "00000000-0000-0000-0000-000000000000",
		}

		json.NewEncoder(w).Encode(response)
	}))

	defer server.Close()

	c, _ := NewClient("test")
	c.Hostname = server.URL

	got, err := c.GetStatusById("00000000-0000-0000-0000-000000000000")

	if err != nil {
		t.Errorf("Error calling GetStatusById(): %s", err)
	}

	want :=
		StatusResponse{
			Id:         "00000000-0000-0000-0000-000000000000",
			Reference:  "00000000-0000-0000-0000-000000000000",
			StatusCode: 200,
		}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetStatusById() = %v, want %v", got, want)
	}
}

func TestHasNext(t *testing.T) {
	t.Parallel()

	s := StatusResponses{
		Links: Link{
			Next: "/v2/notifications?template_type=email&page=2",
		},
	}

	if !s.HasNext() {
		t.Errorf("Expected HasNext() to be true, got false")
	}
}

func TestNextStatusPage(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request URL
		if r.URL.Path != "/v2/notifications" {
			t.Errorf("Expected request to /v2/notifications, got %s", r.URL.Path)
		}

		// Verfiy the query string
		if r.URL.RawQuery != "template_type=email&older_than=00000000-0000-0000-0000-000000000000" {
			t.Errorf("Expected query string to be template_type=email&older_than=00000000-0000-0000-0000-000000000000, got %s", r.URL.RawQuery)
		}

		// Verify the request method
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Write a response
		w.WriteHeader(http.StatusOK)
		response := StatusResponses{
			Notifications: []StatusResponse{
				{
					Id:        "00000000-0000-0000-0000-000000000000",
					Reference: "00000000-0000-0000-0000-000000000000",
				},
			},
		}

		json.NewEncoder(w).Encode(response)
	}))

	defer server.Close()

	c, _ := NewClient("test")
	c.Hostname = server.URL

	s := StatusResponses{
		Links: Link{
			Next: "/v2/notifications?template_type=email&older_than=00000000-0000-0000-0000-000000000000",
		},
	}

	got, err := c.NextStatusPage(s)

	if err != nil {
		t.Errorf("Error calling NextStatusPage(): %s", err)
	}

	want := StatusResponses{
		Notifications: []StatusResponse{
			{
				Id:        "00000000-0000-0000-0000-000000000000",
				Reference: "00000000-0000-0000-0000-000000000000",
			},
		},
		StatusCode: 200,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("NextStatusPage() = %v, want %v", got, want)
	}
}
