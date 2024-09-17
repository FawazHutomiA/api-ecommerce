package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

// SMTPConfig holds the SMTP server configuration.
type SMTPConfig struct {
	Host     string // SMTP server (e.g., smtp.mailgun.org)
	Port     string // Port (e.g., "587")
	Username string // Your SMTP username (e.g., "postmaster@mg.learnhub.id")
	Password string // Your SMTP password (e.g., the Mailgun SMTP password)
}

// SMTPData holds the email data for sending.
type SMTPData struct {
	Sender     string   // Email sender (e.g., "noreply@learnhub.id")
	Subject    string   // Email subject
	Body       string   // Plain text email body
	BodyHTML   string   // HTML email body
	Recipients []string // Recipients' email addresses
	Cc         []string // CC email addresses (optional)
	Bcc        []string // BCC email addresses (optional)
}

// func SendVerificationEmailSMTP(to, verificationLink string) error {
// 	emailFrom := helper.GetENV("EMAIL")
// 	emailPassword := helper.GetENV("PASSWORD")
// 	emailSmtpHost := helper.GetENV("SMTPHOST")
// 	emailSmtpPort := helper.GetENV("SMTPPORT")

// 	from := emailFrom

// 	// Gunakan App Password dari Google, bukan password Gmail biasa
// 	password := emailPassword // Ganti dengan App Password Anda

// 	subject := "Please verify your email"
// 	body := fmt.Sprintf("Click this link to verify your email: %s", verificationLink)

// 	message := []byte(fmt.Sprintf("Subject: %s\n\n%s", subject, body))

// 	smtpHost := emailSmtpHost
// 	smtpPort := emailSmtpPort

// 	auth := smtp.PlainAuth("", from, password, smtpHost)

// 	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

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
