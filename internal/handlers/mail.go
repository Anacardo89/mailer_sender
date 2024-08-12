package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/Anacardo89/mailer_sender/internal/logger"
	"github.com/Anacardo89/mailer_sender/internal/mail"
	"github.com/streadway/amqp"
)

type RegisterData struct {
	Email string `json:"email"`
	User  string `json:"user"`
	Link  string `json:"link"`
}

func SendRegisterEmail(d amqp.Delivery, m *mail.Config, c *smtp.Client, a *smtp.Auth) {
	var regData *RegisterData
	err := json.Unmarshal(d.Body, &regData)
	if err != nil {
		logger.Error.Println(err)
		return
	}
	err = m.ValidateMail(*c, regData.Email)
	if err != nil {
		logger.Error.Println(err)
		return
	}
	mail_subject, mail_body := buildRegisterEmail(regData)
	mailAddress := fmt.Sprintf("%s:%s", m.SmtpHost, m.SmtpPort)
	var to []string
	to = append(to, regData.Email)
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

func buildRegisterEmail(r *RegisterData) (string, string) {
	var mbuf bytes.Buffer

	mailSubject, err := template.New("registerSubject").Parse(registerSubject)
	if err != nil {
		logger.Error.Println(err.Error())
	}
	err = mailSubject.Execute(&mbuf, &r)
	if err != nil {
		logger.Error.Println(err.Error())
	}
	subject := mbuf.String()

	mailBody, err := template.New("registerBody").Parse(registerBody)
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
