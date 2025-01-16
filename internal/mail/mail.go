package mail

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type Config struct {
	SmtpHost string `yaml:"smtp_host"`
	SmtpPort string `yaml:"smtp_port"`
	SmtpUser string `yaml:"smtp_user"`
	SmtpPass string `yaml:"smtp_pass"`
}

func (m *Config) ValidateMail(validate string) error {
	addr := m.SmtpHost + ":" + m.SmtpPort

	// Connect to the SMTP server
	conn, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect to smtp server: %w", err)
	}
	defer conn.Close()

	// Initiate TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false, // Set true only if you need to skip cert verification
		ServerName:         m.SmtpHost,
	}
	if err := conn.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	// Authenticate (if needed)
	auth := smtp.PlainAuth("", m.SmtpUser, m.SmtpPass, m.SmtpHost)
	if err := conn.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Validate email
	if err := conn.Mail(m.SmtpUser); err != nil {
		return fmt.Errorf("failed to set sender address: %w", err)
	}
	if err := conn.Rcpt(validate); err != nil {
		return fmt.Errorf("failed to set recipient address: %w", err)
	}
	return nil
}
