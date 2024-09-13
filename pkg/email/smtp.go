package email

import (
	"example/pkg/helper"
	"fmt"
	"net/smtp"
)

func SendVerificationEmailSMTP(to, verificationLink string) error {
	emailFrom := helper.GetENV("EMAIL")
	emailPassword := helper.GetENV("PASSWORD")
	emailSmtpHost := helper.GetENV("SMTPHOST")
	emailSmtpPort := helper.GetENV("SMTPPORT")

	from := emailFrom

	// Gunakan App Password dari Google, bukan password Gmail biasa
	password := emailPassword // Ganti dengan App Password Anda

	subject := "Please verify your email"
	body := fmt.Sprintf("Click this link to verify your email: %s", verificationLink)

	message := []byte(fmt.Sprintf("Subject: %s\n\n%s", subject, body))

	smtpHost := emailSmtpHost
	smtpPort := emailSmtpPort

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		return err
	}
	return nil
}
