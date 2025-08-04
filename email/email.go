package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/madeinly/core/settings"
)

func Send(to, subject, body string) error {
	appSettings := settings.GetSettings()

	// SMTP server configuration.
	smtpHost := appSettings.Email.Address
	smtpPort := appSettings.Email.Port
	smtpUser := appSettings.Email.User
	smtpPass := appSettings.Email.Password
	encryption := appSettings.Email.Encryption // "ssl/tls" or "starttls"

	// Combine host and port
	serverAddr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	// Authentication.
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Message
	msg := []byte("From: " + smtpUser + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	switch encryption {
	case "ssl/tls":
		// For SSL/TLS, we dial a TLS connection from the start.
		// This is typically used with port 465.
		tlsconfig := &tls.Config{
			ServerName: smtpHost,
		}

		conn, err := tls.Dial("tcp", serverAddr, tlsconfig)
		if err != nil {
			return fmt.Errorf("failed to dial TLS connection: %w", err)
		}

		client, err := smtp.NewClient(conn, smtpHost)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}
		defer client.Close()

		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth failed: %w", err)
		}

		if err = client.Mail(smtpUser); err != nil {
			return fmt.Errorf("smtp mail command failed: %w", err)
		}
		if err = client.Rcpt(to); err != nil {
			return fmt.Errorf("smtp rcpt command failed: %w", err)
		}

		wc, err := client.Data()
		if err != nil {
			return fmt.Errorf("smtp data command failed: %w", err)
		}

		_, err = wc.Write(msg)
		if err != nil {
			return fmt.Errorf("failed to write email body: %w", err)
		}

		err = wc.Close()
		if err != nil {
			return fmt.Errorf("failed to close data writer: %w", err)
		}

		return client.Quit()

	case "starttls":
		// For STARTTLS, we connect in plain text first, then upgrade.
		// This is typically used with port 587.

		// Dial the server
		client, err := smtp.Dial(serverAddr)
		if err != nil {
			return fmt.Errorf("failed to dial SMTP server: %w", err)
		}
		defer client.Close()

		// Check for STARTTLS support and upgrade
		if ok, _ := client.Extension("STARTTLS"); ok {
			tlsconfig := &tls.Config{
				ServerName: smtpHost,
			}
			if err = client.StartTLS(tlsconfig); err != nil {
				return fmt.Errorf("failed to start TLS: %w", err)
			}
		}

		// Authenticate
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth failed: %w", err)
		}

		// Send the email
		if err = client.Mail(smtpUser); err != nil {
			return fmt.Errorf("smtp mail command failed: %w", err)
		}
		if err = client.Rcpt(to); err != nil {
			return fmt.Errorf("smtp rcpt command failed: %w", err)
		}

		wc, err := client.Data()
		if err != nil {
			return fmt.Errorf("smtp data command failed: %w", err)
		}

		_, err = wc.Write(msg)
		if err != nil {
			return fmt.Errorf("failed to write email body: %w", err)
		}

		err = wc.Close()
		if err != nil {
			return fmt.Errorf("failed to close data writer: %w", err)
		}

		return client.Quit()

	default:
		return fmt.Errorf("unsupported email encryption type: '%s'. Please use 'ssl/tls' or 'starttls'", encryption)
	}
}
