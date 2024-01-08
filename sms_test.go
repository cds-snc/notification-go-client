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

func TestSendSms(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request URL
		if r.URL.Path != "/v2/notifications/sms" {
			t.Errorf("Expected request to /v2/notifications/sms, got %s", r.URL.Path)
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
		var got Sms
		err = json.Unmarshal(body, &got)
		if err != nil {
			t.Errorf("Error unmarshalling request body: %s", err)
		}

		// Verify the request body
		want := Sms{
			PhoneNumber: "1234567890",
			TemplateId:  "00000000-0000-0000-0000-000000000000",
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

	s := Sms{
		PhoneNumber: "1234567890",
		TemplateId:  "00000000-0000-0000-0000-000000000000",
	}

	got, err := c.SendSms(s)

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
		t.Errorf("SendSms() = %v, want %v", got, want)
	}
}
