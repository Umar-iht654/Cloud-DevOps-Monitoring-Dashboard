package email

import (
	"fmt"
	"html"
	"log"
	"net/smtp"
	"strings"
	"time"
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

// IsConfigured checks whether SMTP settings are available.
func (s *Sender) IsConfigured() bool {
	// This returns true only when all required SMTP settings are present.
	return s.SMTPHost != "" &&
		s.SMTPPort != "" &&
		s.SMTPUsername != "" &&
		s.SMTPPassword != "" &&
		s.SMTPFrom != ""
}

// SendVerificationEmail sends or logs an email verification link.
func (s *Sender) SendVerificationEmail(toEmail string, token string) error {
	// This builds the verification URL that the user will click.
	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", s.AppBaseURL, token)

	// This checks whether SMTP is configured.
	if !s.IsConfigured() {
		// This logs the link locally so development can continue without real email credentials.
		log.Println("SMTP is not configured. Verification link:", verificationURL)

		// This returns nil because local development fallback succeeded.
		return nil
	}

	// This creates the email subject.
	subject := "Verify your StatusWatch email address"

	// This escapes the verification URL before placing it into HTML.
	safeVerificationURL := html.EscapeString(verificationURL)

	// This creates the HTML email body.
	body := fmt.Sprintf(`
<!doctype html>
<html>
  <body style="margin:0; padding:0; background:#f4f7fb; font-family:Arial, Helvetica, sans-serif; color:#1f2937;">
    <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="background:#f4f7fb; padding:32px 16px;">
      <tr>
        <td align="center">
          <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="max-width:560px; background:#ffffff; border-radius:16px; overflow:hidden; border:1px solid #e5e7eb;">
            <tr>
              <td style="padding:28px 32px; background:#0f172a;">
                <h1 style="margin:0; color:#ffffff; font-size:22px; line-height:1.3;">
                  Verify your email
                </h1>
                <p style="margin:8px 0 0; color:#cbd5e1; font-size:14px; line-height:1.6;">
                  Finish creating your StatusWatch account.
                </p>
              </td>
            </tr>

            <tr>
              <td style="padding:32px;">
                <p style="margin:0 0 16px; font-size:16px; line-height:1.6;">
                  Welcome to <strong>StatusWatch</strong>.
                </p>

                <p style="margin:0 0 24px; font-size:15px; line-height:1.6; color:#4b5563;">
                  Please verify your email address to activate your account and start monitoring your services.
                </p>

                <table role="presentation" cellspacing="0" cellpadding="0" style="margin:0 0 24px;">
                  <tr>
                    <td>
                      <a href="%s" style="display:inline-block; background:#0891b2; color:#ffffff; text-decoration:none; font-weight:700; font-size:15px; padding:13px 20px; border-radius:10px;">
                        Verify email address
                      </a>
                    </td>
                  </tr>
                </table>

                <p style="margin:0 0 8px; font-size:13px; line-height:1.6; color:#6b7280;">
                  If the button does not work, copy and paste this link into your browser:
                </p>

                <p style="margin:0; font-size:13px; line-height:1.6; word-break:break-all;">
                  <a href="%s" style="color:#0891b2;">%s</a>
                </p>

                <div style="height:1px; background:#e5e7eb; margin:28px 0;"></div>

                <p style="margin:0; font-size:13px; line-height:1.6; color:#6b7280;">
                  This verification link will expire soon. If you did not create a StatusWatch account, you can ignore this email.
                </p>
              </td>
            </tr>
          </table>

          <p style="margin:18px 0 0; font-size:12px; color:#94a3b8;">
            StatusWatch · Service monitoring and reliability alerts
          </p>
        </td>
      </tr>
    </table>
  </body>
</html>
`, safeVerificationURL, safeVerificationURL, safeVerificationURL)

	// This sends the HTML email.
	return s.sendHTMLEmail(toEmail, subject, body)
}

// SendDowntimeAlertEmail sends or safely skips a downtime notification email.
func (s *Sender) SendDowntimeAlertEmail(toEmail string, serviceName string, serviceURL string, failureReason string, alertTime time.Time) error {
	// This creates a safe fallback reason if the health checker did not provide one.
	if failureReason == "" {
		failureReason = "The service failed its health check."
	}

	// This checks whether SMTP is configured.
	if !s.IsConfigured() {
		// This logs the email that would have been sent during local development.
		log.Printf(
			"SMTP is not configured. Skipping downtime email to %s. Service: %s. URL: %s. Reason: %s.",
			toEmail,
			serviceName,
			serviceURL,
			failureReason,
		)

		// This returns nil because skipping email locally is expected behaviour.
		return nil
	}

	// This creates the email subject.
	subject := fmt.Sprintf("Service down: %s", serviceName)

	// This escapes dynamic values before placing them into HTML.
	safeServiceName := html.EscapeString(serviceName)
	safeServiceURL := html.EscapeString(serviceURL)
	safeFailureReason := html.EscapeString(failureReason)
	safeAlertTime := html.EscapeString(alertTime.Format("Mon, 02 Jan 2006 15:04:05 MST"))

	// This creates the HTML email body.
	body := fmt.Sprintf(`
<!doctype html>
<html>
  <body style="margin:0; padding:0; background:#f4f7fb; font-family:Arial, Helvetica, sans-serif; color:#1f2937;">
    <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="background:#f4f7fb; padding:32px 16px;">
      <tr>
        <td align="center">
          <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="max-width:600px; background:#ffffff; border-radius:16px; overflow:hidden; border:1px solid #e5e7eb;">
            <tr>
              <td style="padding:28px 32px; background:#991b1b;">
                <p style="margin:0 0 8px; color:#fecaca; font-size:13px; font-weight:700; letter-spacing:0.04em; text-transform:uppercase;">
                  Downtime alert
                </p>
                <h1 style="margin:0; color:#ffffff; font-size:22px; line-height:1.3;">
                  %s is down
                </h1>
              </td>
            </tr>

            <tr>
              <td style="padding:32px;">
                <p style="margin:0 0 20px; font-size:15px; line-height:1.6; color:#4b5563;">
                  StatusWatch detected that one of your monitored services has failed its health check.
                </p>

                <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="border-collapse:collapse; margin:0 0 24px;">
                  <tr>
                    <td style="padding:12px 0; border-bottom:1px solid #e5e7eb; width:140px; color:#6b7280; font-size:14px;">
                      Service
                    </td>
                    <td style="padding:12px 0; border-bottom:1px solid #e5e7eb; font-size:14px; font-weight:700;">
                      %s
                    </td>
                  </tr>
                  <tr>
                    <td style="padding:12px 0; border-bottom:1px solid #e5e7eb; color:#6b7280; font-size:14px;">
                      URL
                    </td>
                    <td style="padding:12px 0; border-bottom:1px solid #e5e7eb; font-size:14px; word-break:break-all;">
                      <a href="%s" style="color:#0891b2;">%s</a>
                    </td>
                  </tr>
                  <tr>
                    <td style="padding:12px 0; border-bottom:1px solid #e5e7eb; color:#6b7280; font-size:14px;">
                      Reason
                    </td>
                    <td style="padding:12px 0; border-bottom:1px solid #e5e7eb; font-size:14px;">
                      %s
                    </td>
                  </tr>
                  <tr>
                    <td style="padding:12px 0; color:#6b7280; font-size:14px;">
                      Alert time
                    </td>
                    <td style="padding:12px 0; font-size:14px;">
                      %s
                    </td>
                  </tr>
                </table>

                <p style="margin:0 0 24px; font-size:14px; line-height:1.6; color:#4b5563;">
                  The alert has also been saved in your dashboard. Repeated failed checks will not create duplicate downtime emails while the service remains down.
                </p>

                <div style="padding:14px 16px; background:#fef2f2; border:1px solid #fecaca; border-radius:12px;">
                  <p style="margin:0; font-size:13px; line-height:1.6; color:#7f1d1d;">
                    Check the service, confirm whether the outage is expected, and update the monitored URL or service configuration if needed.
                  </p>
                </div>
              </td>
            </tr>
          </table>

          <p style="margin:18px 0 0; font-size:12px; color:#94a3b8;">
            StatusWatch · Automated downtime notification
          </p>
        </td>
      </tr>
    </table>
  </body>
</html>
`, safeServiceName, safeServiceName, safeServiceURL, safeServiceURL, safeFailureReason, safeAlertTime)

	// This sends the HTML email.
	if err := s.sendHTMLEmail(toEmail, subject, body); err != nil {
		// This returns the error so the caller can log it without crashing the worker.
		return err
	}

	// This logs that the email was sent successfully.
	log.Printf("Sent downtime email to %s for service %s", toEmail, serviceName)

	// This returns nil because the email was sent.
	return nil
}

// sendHTMLEmail sends an HTML email through the configured SMTP server.
func (s *Sender) sendHTMLEmail(toEmail string, subject string, body string) error {
	// This removes newline characters from headers to reduce header injection risk.
	cleanFrom := sanitizeEmailHeader(s.SMTPFrom)
	cleanTo := sanitizeEmailHeader(toEmail)
	cleanSubject := sanitizeEmailHeader(subject)

	// This creates the full raw email message.
	message := []byte(
		"From: StatusWatch <" + cleanFrom + ">\r\n" +
			"To: " + cleanTo + "\r\n" +
			"Subject: " + cleanSubject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" +
			"\r\n" +
			body,
	)

	// This builds the SMTP server address.
	address := s.SMTPHost + ":" + s.SMTPPort

	// This creates SMTP authentication.
	auth := smtp.PlainAuth("", s.SMTPUsername, s.SMTPPassword, s.SMTPHost)

	// This sends the email through the configured SMTP server.
	if err := smtp.SendMail(address, auth, cleanFrom, []string{cleanTo}, message); err != nil {
		// This returns the error so the caller can decide what to do.
		return err
	}

	// This returns nil because the email was sent.
	return nil
}

// sanitizeEmailHeader removes characters that should not appear in email headers.
func sanitizeEmailHeader(value string) string {
	// This removes carriage returns and line breaks from header values.
	value = strings.ReplaceAll(value, "\r", "")
	value = strings.ReplaceAll(value, "\n", "")

	// This removes extra spaces before returning the header value.
	return strings.TrimSpace(value)
}
