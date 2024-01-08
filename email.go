package client

import (
	"encoding/json"
	"fmt"
)

type Email struct {
	EmailAddress string `json:"email_address"`
	TemplateId   string `json:"template_id"`

	// Optional
	EmailReplyToId  string                 `json:"email_reply_to_id,omitempty"`
	Personalisation map[string]interface{} `json:"personalisation,omitempty"`
	Reference       string                 `json:"reference,omitempty"`
}

func (c Client) SendEmail(e Email) (Response, error) {
	body, err := json.Marshal(e)

	var response Response

	if err != nil {
		return response, fmt.Errorf("error marshalling body: %s", err)
	}

	resp, err := c.DoPostRequest("/v2/notifications/email", body)

	if err != nil {
		return response, fmt.Errorf("error calling email endpoint: %s", err)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		return response, fmt.Errorf("error decoding email response: %s", err)
	}

	response.StatusCode = resp.StatusCode

	return response, nil
}
