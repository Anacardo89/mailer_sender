package main

import (
	"net/smtp"
)

func (m *MailConfig) validateMail(client smtp.Client, validate string) error {
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
