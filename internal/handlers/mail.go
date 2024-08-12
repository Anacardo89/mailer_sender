package handlers

import (
	"encoding/json"
	"fmt"
	"net/smtp"

	"github.com/Anacardo89/mailer_sender/internal/config"
	"github.com/Anacardo89/mailer_sender/internal/logger"
	"github.com/streadway/amqp"
)

type RegisterData struct {
	Email    string `json:"email"`
	Subject  string `json:"subject"`
	MailBody string `json:"mail_body"`
}

func SendRegistrationEmail(d amqp.Delivery) {

	// Mail Setup
	mailConfig, err := config.LoadMailConfig()
	if err != nil {
		logger.Error.Fatal(err)
	}
	auth := smtp.PlainAuth("", mailConfig.SmtpUser, mailConfig.SmtpPass, mailConfig.SmtpHost)

	smtpAddress := fmt.Sprintf("%s:%s",
		mailConfig.SmtpHost, mailConfig.SmtpPort)
	client, err := smtp.Dial(smtpAddress)
	if err != nil {
		logger.Error.Fatal(err)
	}

	var regData *RegisterData
	err = json.Unmarshal(d.Body, &regData)
	if err != nil {
		logger.Error.Println(err)
	}
	err = mailConfig.ValidateMail(*client, regData.Email)
	if err != nil {
		logger.Error.Println(err)
	}
	mailAddress := fmt.Sprintf("%s:%s", mailConfig.SmtpHost, mailConfig.SmtpPort)
	var to []string
	to = append(to, regData.Email)
	fromHeader := fmt.Sprintf("From: %s\n", mailConfig.SmtpUser)
	toHeader := fmt.Sprintf("To: %s\n", to)
	subject := fmt.Sprintf("Subject: %s\n", regData.Subject)
	body := regData.MailBody
	message := []byte(fromHeader + toHeader + subject + "\n" + body)
	err = smtp.SendMail(mailAddress,
		auth,
		mailConfig.SmtpUser,
		to,
		message)
	if err != nil {
		fmt.Println(err)
		return
	}
	d.Ack(false)
}
