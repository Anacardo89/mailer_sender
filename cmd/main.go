package main

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/Anacardo89/mailer_sender/internal/config"
	"github.com/Anacardo89/mailer_sender/internal/handlers"
	"github.com/Anacardo89/mailer_sender/internal/logger"
	"github.com/streadway/amqp"
)

func main() {
	logger.CreateLogger()

	// Mail Setup
	mail, err := config.LoadMailConfig()
	if err != nil {
		logger.Error.Fatal(err)
	}
	auth := smtp.PlainAuth("", mail.SmtpUser, mail.SmtpPass, mail.SmtpHost)

	smtpAddress := fmt.Sprintf("%s:%s",
		mail.SmtpHost, mail.SmtpPort)
	client, err := smtp.Dial(smtpAddress)
	if err != nil {
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
	rabbit.StartWorkers(conn, ch, msgs)

	forever := make(chan struct{})
	for msg := range msgs {
		go func(m amqp.Delivery) {
			switch msg.RoutingKey {
			case "register_email":
				handlers.SendRegisterEmail(msg, mail, client, &auth)
			}
		}(msg)
	}
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
