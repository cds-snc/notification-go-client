package client_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	. "github.com/cds-snc/notification-go-client"
)

func TestSendBulkEmail(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request URL
		if r.URL.Path != "/v2/notifications/bulk" {
			t.Errorf("Expected request to /v2/notifications/bulk, got %s", r.URL.Path)
		}

		// Verify the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Error reading request body: %s", err)
		}

		// Unmarshal the request body
		var got BulkEmail
		err = json.Unmarshal(body, &got)
		if err != nil {
			t.Errorf("Error unmarshalling request body: %s", err)
		}

		// Verify the request body
		want := BulkEmail{
			Name:       "Test Bulk Email",
			TemplateId: "00000000-0000-0000-0000-000000000000",
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Expected request body to be %+v, got %+v", want, got)
		}

		// Write a response
		w.WriteHeader(http.StatusInternalServerError)
		response := BulkEmailResponse{
			Errors: []ResponseError{
				{
					Error:   "Error",
					Message: "Message",
				},
			},
		}

		json.NewEncoder(w).Encode(response)
	}))

	defer server.Close()

	c, _ := NewClient("test")
	c.Hostname = server.URL

	e := BulkEmail{
		Name:       "Test Bulk Email",
		TemplateId: "00000000-0000-0000-0000-000000000000",
	}

	got, err := c.SendBulkEmail(e)

	if err != nil {
		t.Errorf("Error sending email: %s", err)
	}

	want := BulkEmailResponse{
		Errors: []ResponseError{
			{
				Error:   "Error",
				Message: "Message",
			},
		},
		StatusCode: 500,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("SendBulkEmail() = %+v, want %+v", got, want)
	}
}
