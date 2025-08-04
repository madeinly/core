package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/madeinly/core/internal/settings"
)

// Send dispatches an email, handling both SSL/TLS and STARTTLS connections.
func Send(to, subject, body string) error {
	appSettings := settings.GetSettings().Email

	serverAddr := fmt.Sprintf("%s:%s", appSettings.Address, appSettings.Port)
	auth := smtp.PlainAuth("", appSettings.User, appSettings.Password, appSettings.Address)

	var client *smtp.Client
	var err error

	switch appSettings.Encryption {
	case "ssl/tls":
		conn, err := tls.Dial("tcp", serverAddr, &tls.Config{ServerName: appSettings.Address})
		if err != nil {
			return fmt.Errorf("failed to dial TLS connection: %w", err)
		}
		client, err = smtp.NewClient(conn, appSettings.Address)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}

	case "starttls":
		client, err = smtp.Dial(serverAddr)
		if err != nil {
			return fmt.Errorf("failed to dial SMTP server: %w", err)
		}
		if ok, _ := client.Extension("STARTTLS"); ok {
			if err = client.StartTLS(&tls.Config{ServerName: appSettings.Address}); err != nil {
				return fmt.Errorf("failed to start TLS: %w", err)
			}
		}

	default:
		return fmt.Errorf("unsupported email encryption type: '%s'", appSettings.Encryption)
	}

	defer client.Close()

	return sendMail(client, auth, appSettings.User, to, subject, body)
}

// sendMail handles the SMTP commands to authenticate and send the email.
func sendMail(client *smtp.Client, auth smtp.Auth, from, to, subject, body string) error {
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth failed: %w", err)
	}
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("smtp mail command failed: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("smtp rcpt command failed: %w", err)
	}

	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data command failed: %w", err)
	}

	msg := []byte("From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	if _, err = wc.Write(msg); err != nil {
		return fmt.Errorf("failed to write email body: %w", err)
	}
	if err = wc.Close(); err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	return client.Quit()
}
