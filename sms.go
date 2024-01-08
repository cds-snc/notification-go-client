package client

import (
	"encoding/json"
	"fmt"
)

type Sms struct {
	PhoneNumber string `json:"phone_number"`
	TemplateId  string `json:"template_id"`

	// Optional
	SmsSenderId     string            `json:"sms_sender_id,omitempty"`
	Personalisation map[string]string `json:"personalisation,omitempty"`
	Reference       string            `json:"reference,omitempty"`
}

func (c Client) SendSms(s Sms) (Response, error) {
	body, err := json.Marshal(s)

	var response Response

	if err != nil {
		return response, fmt.Errorf("error marshalling body: %s", err)
	}

	resp, err := c.DoPostRequest("/v2/notifications/sms", body)

	if err != nil {
		return response, fmt.Errorf("error calling sms endpoint: %s", err)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		return response, fmt.Errorf("error decoding sms response: %s", err)
	}

	response.StatusCode = resp.StatusCode

	return response, nil
}
