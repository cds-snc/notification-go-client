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

func TestSendEmail(t *testing.T) {
	t.Parallel()

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
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Error reading request body: %s", err)
		}

		// Unmarshal the request body
		var got Email
		err = json.Unmarshal(body, &got)
		if err != nil {
			t.Errorf("Error unmarshalling request body: %s", err)
		}

		// Verify the request body
		want := Email{
			EmailAddress: "test@test.com",
			TemplateId:   "00000000-0000-0000-0000-000000000000",
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Expected request body to be %+v, got %+v", want, got)
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

	c, _ := NewClient("test")
	c.Hostname = server.URL

	e := Email{
		EmailAddress: "test@test.com",
		TemplateId:   "00000000-0000-0000-0000-000000000000",
	}

	got, err := c.SendEmail(e)

	if err != nil {
		t.Errorf("Error sending email: %s", err)
	}

	want := Response{
		Id:         "00000000-0000-0000-0000-000000000000",
		Reference:  "00000000-0000-0000-0000-000000000000",
		Uri:        "https://api.notification.canada.ca/v2/notifications/00000000-0000-0000-0000-000000000000",
		StatusCode: 201,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("SendEmail() = %v, want %v", got, want)
	}
}
