package main

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/Anacardo89/mailer_sender/internal/config"
	"github.com/Anacardo89/mailer_sender/internal/handlers"
	"github.com/Anacardo89/mailer_sender/internal/rabbitmq"
	"github.com/Anacardo89/mailer_sender/pkg/logger"

	"github.com/streadway/amqp"
)

func main() {
	logger.CreateLogger()

	// Mail Setup
	mail := config.LoadMailConfig()

	smtpAddress := fmt.Sprintf("%s:%s",
		mail.SmtpHost, mail.SmtpPort)
	client, err := smtp.Dial(smtpAddress)
	if err != nil {
		logger.Error.Fatal(err)
	}
	defer client.Close()
	client.StartTLS(&tls.Config{
		ServerName: mail.SmtpHost,
	})
	auth := smtp.PlainAuth("", mail.SmtpUser, mail.SmtpPass, mail.SmtpHost)
	if err = client.Auth(auth); err != nil {
		logger.Error.Fatal(err)
	}

	// Rabbit Setup
	rabbit := config.LoadRabbitConfig()

	conn := rabbit.Connect()
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logger.Error.Fatal(err)
	}
	defer ch.Close()

	rabbit.DeclareQueues(ch)

	msgs := make(chan amqp.Delivery)
	rabbitmq.StartWorkers(rabbit, conn, ch, msgs)

	fmt.Println(" [*] Waiting for messages. To exit press CTRL+C")

	forever := make(chan struct{})
	for msg := range msgs {
		go func(m amqp.Delivery) {
			switch msg.RoutingKey {
			case "register_email":
				handlers.SendRegisterEmail(msg, mail, client, &auth)
			case "password_recover_email":
				handlers.SendPasswordRecoveryEmail(msg, mail, client, &auth)
			}
		}(msg)
	}
	<-forever
}
