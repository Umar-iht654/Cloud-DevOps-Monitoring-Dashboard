package email

import (
	"fmt"
	"log"
	"net/smtp"
)

// Sender stores email configuration.
type Sender struct {
	// SMTPHost stores the SMTP server hostname.
	SMTPHost string

	// SMTPPort stores the SMTP server port.
	SMTPPort string

	// SMTPUsername stores the SMTP username.
	SMTPUsername string

	// SMTPPassword stores the SMTP password or API key.
	SMTPPassword string

	// SMTPFrom stores the sender email address.
	SMTPFrom string

	// AppBaseURL stores the frontend URL used in verification links.
	AppBaseURL string
}

// NewSender creates a new email sender.
func NewSender(smtpHost string, smtpPort string, smtpUsername string, smtpPassword string, smtpFrom string, appBaseURL string) *Sender {
	// This returns an email sender with environment-based configuration.
	return &Sender{
		SMTPHost:     smtpHost,
		SMTPPort:     smtpPort,
		SMTPUsername: smtpUsername,
		SMTPPassword: smtpPassword,
		SMTPFrom:     smtpFrom,
		AppBaseURL:   appBaseURL,
	}
}

// SendVerificationEmail sends or logs an email verification link.
func (s *Sender) SendVerificationEmail(toEmail string, token string) error {
	// This builds the verification URL that the user will click.
	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", s.AppBaseURL, token)

	// This checks whether SMTP is configured.
	if s.SMTPHost == "" || s.SMTPPort == "" || s.SMTPUsername == "" || s.SMTPPassword == "" || s.SMTPFrom == "" {
		// This logs the link locally so development can continue without real email credentials.
		log.Println("SMTP is not configured. Verification link:", verificationURL)

		// This returns nil because local development fallback succeeded.
		return nil
	}

	// This creates the email subject.
	subject := "Verify your Cloud Monitor email address"

	// This creates the plain-text email body.
	body := fmt.Sprintf(
		"Welcome to Cloud Monitor.\n\nVerify your email address by opening this link:\n\n%s\n\nThis link will expire soon.\n",
		verificationURL,
	)

	// This creates the full raw email message.
	message := []byte(
		"From: " + s.SMTPFrom + "\r\n" +
			"To: " + toEmail + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n" +
			"\r\n" +
			body,
	)

	// This builds the SMTP server address.
	address := s.SMTPHost + ":" + s.SMTPPort

	// This creates SMTP authentication.
	auth := smtp.PlainAuth("", s.SMTPUsername, s.SMTPPassword, s.SMTPHost)

	// This sends the email through the configured SMTP server.
	if err := smtp.SendMail(address, auth, s.SMTPFrom, []string{toEmail}, message); err != nil {
		// This returns the error so the caller can decide what to do.
		return err
	}

	// This returns nil because the email was sent.
	return nil
}
