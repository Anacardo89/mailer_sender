package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"

	"github.com/Anacardo89/mailer_sender.git/logger"
	"github.com/streadway/amqp"
)

type RegisterData struct {
	Email    string `json:"email"`
	Subject  string `json:"subject"`
	MailBody string `json:"mail_body"`
}

func main() {
	logger.CreateLogger()

	// Mail Setup
	mailConfig, err := loadMailConfig()
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

	// Rabbit Setup
	rabbitConfig, err := loadRabbitConfig()
	if err != nil {
		logger.Error.Fatal(err)
	}

	url := fmt.Sprintf("amqp://%s:%s@%s%s/",
		rabbitConfig.RabbitUser, rabbitConfig.RabbitPass, rabbitConfig.RabbitHost, rabbitConfig.RabbitPort)
	conn, err := amqp.Dial(url)
	if err != nil {
		logger.Error.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logger.Error.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		rabbitConfig.QueueRegister, // name
		true,                       // durable
		false,                      // delete when unused
		false,                      // exclusive
		false,                      // no-wait
		nil,                        // arguments
	)
	if err != nil {
		logger.Error.Fatal(err)
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		logger.Error.Fatal(err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		logger.Error.Fatal(err)
	}

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			var regData *RegisterData
			err = json.Unmarshal(d.Body, &regData)
			if err != nil {
				logger.Error.Println(err)
			}
			err = mailConfig.validateMail(*client, regData.Email)
			if err != nil {
				logger.Error.Println(err)
			}
			mailAddress := fmt.Sprintf("%s:%s", mailConfig.SmtpHost, mailConfig.SmtpPort)
			var to []string
			to = append(to, regData.Email)
			message := fmt.Sprintf("%s\r\n\r\n%s\r\n",
				regData.Subject, regData.MailBody)
			err = smtp.SendMail(mailAddress,
				auth,
				mailConfig.PublicSender,
				to,
				[]byte(message))
			if err != nil {
				fmt.Println(err)
				return
			}
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
