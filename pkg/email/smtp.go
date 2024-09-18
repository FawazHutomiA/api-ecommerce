package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

// SMTPConfig holds the SMTP server configuration.
type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

// SMTPData holds the email data for sending.
type SMTPData struct {
	Sender     string
	Subject    string
	Body       string
	BodyHTML   string
	Recipients []string
	Cc         []string
	Bcc        []string
}

// SendSMTPEmail sends an email using SMTP.
func SendSMTPEmail(config SMTPConfig, data SMTPData) error {
	// Set up authentication information.
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	// Prepare the message headers and body.
	to := strings.Join(data.Recipients, ", ")
	headers := map[string]string{
		"From":         data.Sender,
		"To":           to,
		"Subject":      data.Subject,
		"MIME-Version": "1.0",
		"Content-Type": "multipart/alternative; boundary=boundary42",
	}

	// Build the email message
	message := ""
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	// Adding boundary for different content types (plain text and HTML)
	message += "\r\n--boundary42\r\n"
	message += "Content-Type: text/plain; charset=UTF-8\r\n\r\n"
	message += data.Body + "\r\n"

	message += "\r\n--boundary42\r\n"
	message += "Content-Type: text/html; charset=UTF-8\r\n\r\n"
	message += data.BodyHTML + "\r\n"

	message += "--boundary42--"

	// Combine all recipients: To, Cc, Bcc
	allRecipients := append(data.Recipients, data.Cc...)
	allRecipients = append(allRecipients, data.Bcc...)

	// Send the email
	err := smtp.SendMail(config.Host+":"+config.Port, auth, data.Sender, allRecipients, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
