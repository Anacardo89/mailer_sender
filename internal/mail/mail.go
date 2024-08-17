package mail

import (
	"net/smtp"
)

type Config struct {
	SmtpHost string `yaml:"smtp_host"`
	SmtpPort string `yaml:"smtp_port"`
	SmtpUser string `yaml:"smtp_user"`
	SmtpPass string `yaml:"smtp_pass"`
}

func (m *Config) ValidateMail(client *smtp.Client, validate string) error {
	err := client.Mail(m.SmtpUser)
	if err != nil {
		return err
	}
	err = client.Rcpt(validate)
	if err != nil {
		return err
	}
	return nil
}
