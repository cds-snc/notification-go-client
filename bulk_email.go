package client

import (
	"encoding/json"
	"fmt"
)

type BulkEmail struct {
	Name       string `json:"name"`
	TemplateId string `json:"template_id"`

	// Optional
	Rows         [][]string `json:"rows,omitempty"`
	ScheduledFor string     `json:"scheduled_for,omitempty"`
	ReplyToId    string     `json:"reply_to_id,omitempty"`
	Csv          string     `json:"csv,omitempty"`
}

type bulkEmailDataResponseApiKey struct {
	Id      string `json:"id"`
	KeyType string `json:"key_type"`
	Name    string `json:"name"`
}

type bulkEmailDataResponseCreatedBy struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type bulkEmailDataResponseService struct {
	Name string `json:"name"`
}

type bulkEmailDataResponse struct {
	ApiKey             bulkEmailDataResponseApiKey    `json:"api_key"`
	Archived           bool                           `json:"archived"`
	CreatedAt          string                         `json:"created_at"`
	CreatedBy          bulkEmailDataResponseCreatedBy `json:"created_by"`
	Id                 string                         `json:"id"`
	JobStatus          string                         `json:"job_status"`
	NotificationCount  int                            `json:"notification_count"`
	OriginalFileName   string                         `json:"original_file_name"`
	ProcessingFinished string                         `json:"processing_finished"`
	ProcessingStarted  string                         `json:"processing_started"`
	SecheduledFor      string                         `json:"scheduled_for"`
	SenderId           string                         `json:"sender_id"`
	Service            string                         `json:"service"`
	ServiceName        bulkEmailDataResponseService   `json:"service_name"`
	Template           string                         `json:"template"`
	TemplateVersion    int                            `json:"template_version"`
	UpdatedAt          string                         `json:"updated_at"`
}

type BulkEmailResponse struct {
	// Valid Response
	Data bulkEmailDataResponse `json:"data"`

	// Error Response
	StatusCode int             `json:"status_code"`
	Errors     []ResponseError `json:"errors"`
}

func (c Client) SendBulkEmail(e BulkEmail) (BulkEmailResponse, error) {
	body, err := json.Marshal(e)

	var response BulkEmailResponse

	if err != nil {
		return response, fmt.Errorf("error marshalling body: %s", err)
	}

	resp, err := c.DoPostRequest("/v2/notifications/bulk", body)

	if err != nil {
		return response, fmt.Errorf("error calling bulk email endpoint: %s", err)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		return response, fmt.Errorf("error decoding bulk email response: %s", err)
	}

	response.StatusCode = resp.StatusCode

	return response, nil
}
