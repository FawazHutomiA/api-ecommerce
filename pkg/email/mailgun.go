package email

import (
	"context"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

type MailgunConfig struct {
	Domain string
	ApiKey string
}

type MailgunData struct {
	Sender     string
	Subject    string
	Body       string
	BodyHTML   string
	Recipients []string
	Cc         []string
	Bcc        []string
}

type MailgunResponse struct {
	ID       string
	Response string
}

func SendMailgunEmail(config MailgunConfig, data MailgunData) (resp MailgunResponse, err error) {
	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(config.Domain, config.ApiKey)

	//When you have an EU-domain, you must specify the endpoint:
	//mg.SetAPIBase("https://api.eu.mailgun.net/v3")

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(data.Sender, data.Subject, data.Body, data.Recipients...)

	if data.BodyHTML != "" {
		message.SetHtml(data.BodyHTML)
	}

	// Set CC and BCC
	for _, c := range data.Cc {
		message.AddCC(c)
	}

	for _, b := range data.Bcc {
		message.AddBCC(b)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	responses, id, err := mg.Send(ctx, message)
	if err != nil {
		return resp, fmt.Errorf("failed to send email: %w", err)
	}

	resp = MailgunResponse{
		ID:       id,
		Response: responses,
	}

	return resp, nil
}
