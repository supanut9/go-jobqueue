package infrastructure

import (
	"crypto/tls"
	"fmt"
	"net" // <-- IMPORTANT: ADD THIS IMPORT
	"net/smtp"
)

// SMTPGateway is the concrete implementation of the usecase.EmailService interface.
type SMTPGateway struct {
	config SMTPConfig
}

func NewSMTPGateway(config SMTPConfig) *SMTPGateway {
	return &SMTPGateway{config: config}
}

// Send implements the EmailService interface using net/smtp.
func (g *SMTPGateway) Send(to string, subject string, body string) error {
	addr := net.JoinHostPort(g.config.Host, g.config.Port)
	// The authentication mechanism requires the server name (g.config.Host)
	auth := smtp.PlainAuth("", g.config.Username, g.config.Password, g.config.Host)

	// Build the full MIME message
	msg := []byte(
		fmt.Sprintf("To: %s\r\n", to) +
			fmt.Sprintf("From: %s\r\n", g.config.From) +
			fmt.Sprintf("Subject: %s\r\n", subject) +
			"MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\r\n" +
			"\r\n" +
			body + "\r\n")

	// 1. Connect to the SMTP server (UNENCRYPTED initial connection using net.Dial)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}

	client, err := smtp.NewClient(conn, g.config.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	// 2. UPGRADE the connection using STARTTLS (Explicit TLS)
	t := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         g.config.Host,
	}
	if err = client.StartTLS(t); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	// 3. Perform authentication over the secure channel
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate with SMTP: %w", err)
	}

	// 4. Send the message (Steps 3, 4, and 5 from your original code)
	toSlice := []string{to}
	if err = client.Mail(g.config.From); err != nil {
		return err
	}
	for _, recipient := range toSlice {
		if err = client.Rcpt(recipient); err != nil {
			return err
		}
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}
	_, err = writer.Write(msg)
	if err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		return err
	}

	return client.Quit()
}
