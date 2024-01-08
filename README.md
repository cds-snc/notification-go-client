# Notification Go Client

This is a Go client for the Government of Canada's [Notification API](https://documentation.notification.canada.ca/en/). You can use it for other variants of the Notify API by changing the hostname. ex: `c.Hostname = "https://api.notifications.service.gov.uk/"`

The client library is designed to be explicit and minimal without significant abstractions. As a result most structs are directly accessible without helper functions such as variardic options in constructors and setters. Additionally there is very light validation on data with a preference for the Notify API to return validation errors.

## Initializing the client
```
import client "github.com/cds-snc/notification-go-client"

func main() {
	api_key := os.Getenv("NOTIFICATION_API_KEY")

	c, err := client.NewClient(api_key)

	if err != nil {
		fmt.Printf("Error creating client: %s", err)
	}

	// Change the hostname if you are using a different environment
	c.Hostname = "https://api.notifications.service.gov.uk/"
	
	// Use the client
	...
}
```

## Sending an email without personalisation
```
	e := client.Email{
		EmailAddress: "test@test.com",
		TemplateId:   "00000000-0000-0000-0000-000000000000",
	}

	resp, err := c.SendEmail(e)

	if err != nil {
		fmt.Printf("Error sending email: %s", err)
	}

	fmt.Printf("Response: %+v", resp)
```

## Sending an email with personalisation including a file
```
	e := client.Email{
		EmailAddress: "test@test.com",
		TemplateId:   "00000000-0000-0000-0000-000000000000",,
		Personalisation: map[string]interface{}{
			"subject": "Hello, world!",
			"body":    "This is a test email"
			"application_file": {
				"file": "file as base64 encoded string",
				"filename": "your_custom_filename.pdf",
				"sending_method": "attach",
			}
		},
	}

	resp, err := c.SendEmail(e)

	if err != nil {
		fmt.Printf("Error sending email: %s", err)
	}

	fmt.Printf("Response: %+v", resp)
```

## Sending a bulk email with rows defined
```
	rows := make([][]string, 2)
	rows[0] = []string{"email_address"}
	rows[1] = []string{"test@test.com"}

	e := client.BulkEmail{
		Name:       "Test",
		TemplateId: "00000000-0000-0000-0000-000000000000",
		Rows: 	    rows,
	}

	resp, err := c.SendEmail(e)

	if err != nil {
		fmt.Printf("Error sending email: %s", err)
	}

	fmt.Printf("Response: %+v", resp)
```

## Sending a bulk email with a CSV string
```
	csv := "email_address\ntest@test.com"

	e := client.BulkEmail{
		Name:       "Test",
		TemplateId: "00000000-0000-0000-0000-000000000000",
		Csv:        csv,
	}

	resp, err := c.SendEmail(e)

	if err != nil {
		fmt.Printf("Error sending email: %s", err)
	}

	fmt.Printf("Response: %+v", resp)
```

## Sending a SMS message
```
	s := client.SMS{
		PhoneNumber: "1234567890",
		TemplateId:  "00000000-0000-0000-0000-000000000000",
	}

	resp, err := c.SendSMS(s)

	if err != nil {
		fmt.Printf("Error sending sms: %s", err)
	}

	fmt.Printf("Response: %+v", resp)
```

## Handling errors from the Notify API
```
	// Send an email with an invalid template ID
	e := client.Email{
		EmailAddress: "test@test.com",
		TemplateId:   "00000000-0000",
	}

	resp, err := c.SendEmail(e)

	// Assume there is no error with the actual request
	if err != nil {
		fmt.Printf("Error sending email: %s", err)
	}

	if resp.StatusCode == 400 {
		fmt.Printf("Error: %s", resp.Errors[0].Message)
	}
```

## Getting the status of notifications
```
	queryOptions := StatusQueryOptions{
		TemplateType: "email",
	}

	resp, err := c.GetStatus(queryOptions)

	if err != nil {
		fmt.Printf("Error getting status: %s", err)
	}

	// Example of pagination

	fmt.Printf("Page 1: received %d notifications\n", len(resp.Notifications))
	i := 1
	for resp.HasNext() {
		resp, _ = c.NextStatusPage(resp)
		fmt.Printf("Page %d: received %d notifications\n", i+1, len(resp.Notifications))
		i++
	}

```

## Getting the status of a single notification
```
	notificationId := "00000000-0000-0000-0000-000000000000"

	resp, err := c.GetStatusById(notificationId)

	if err != nil {
		fmt.Printf("Error getting status: %s", err)
	}

	fmt.Printf("Response: %+v", resp)
```

## License 
MIT License