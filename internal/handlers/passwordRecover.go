package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/Anacardo89/mailer_sender/internal/mail"
	"github.com/Anacardo89/mailer_sender/pkg/logger"
	"github.com/streadway/amqp"
)

type PasswordRecovery struct {
	Email string `json:"email"`
	User  string `json:"user"`
	Link  string `json:"link"`
}

func SendPasswordRecoveryEmail(d amqp.Delivery, m *mail.Config, c *smtp.Client, a *smtp.Auth) {
	var mailData *PasswordRecovery
	err := json.Unmarshal(d.Body, &mailData)
	if err != nil {
		logger.Error.Println(err)
		return
	}
	err = m.ValidateMail(mailData.Email)
	if err != nil {
		logger.Error.Println(err)
		return
	}
	mail_subject, mail_body := buildPasswordRecoveryEmail(mailData)
	mailAddress := fmt.Sprintf("%s:%s", m.SmtpHost, m.SmtpPort)
	var to []string
	to = append(to, mailData.Email)
	fromHeader := fmt.Sprintf("From: %s\n", m.SmtpUser)
	toHeader := fmt.Sprintf("To: %s\n", to)
	subject := fmt.Sprintf("Subject: %s\n", mail_subject)
	body := mail_body
	message := []byte(fromHeader + toHeader + subject + "\n" + body)
	err = smtp.SendMail(mailAddress,
		*a,
		m.SmtpUser,
		to,
		message)
	if err != nil {
		fmt.Println(err)
		return
	}
	d.Ack(false)
}

func buildPasswordRecoveryEmail(r *PasswordRecovery) (string, string) {
	var mbuf bytes.Buffer

	mailSubject, err := template.New("recoverPasswordSubject").Parse(recoverPasswordSubject)
	if err != nil {
		logger.Error.Println(err.Error())
	}
	err = mailSubject.Execute(&mbuf, &r)
	if err != nil {
		logger.Error.Println(err.Error())
	}
	subject := mbuf.String()

	mailBody, err := template.New("recoverPasswordBody").Parse(recoverPasswordBody)
	if err != nil {
		logger.Error.Println(err.Error())
	}
	mbuf.Reset()
	err = mailBody.Execute(&mbuf, r)
	if err != nil {
		logger.Error.Println(err.Error())
	}
	body := mbuf.String()

	return subject, body
}
